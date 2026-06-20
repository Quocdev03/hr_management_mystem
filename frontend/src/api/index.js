import axios from 'axios';

// Cấu hình Axios Instance
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 10000,
});

// Biến trạng thái cho cơ chế refresh token
let isRefreshing = false;
let failedQueue = [];

/**
 * Giải phóng hàng đợi các request đang chờ sau khi refresh token hoàn tất.
 * @param {Error|null} error
 * @param {string|null} token
 */
const processQueue = (error, token = null) => {
  failedQueue.forEach((p) => (error ? p.reject(error) : p.resolve(token)));
  failedQueue = [];
};

/**
 * Xóa toàn bộ dữ liệu xác thực và chuyển hướng về trang đăng nhập.
 */
const clearAuthData = () => {
  localStorage.removeItem('access_token');
  localStorage.removeItem('refresh_token');
  localStorage.removeItem('user');

  if (window.location.pathname !== '/login') {
    window.location.href = '/login';
  }
};

/**
 * Kiểm tra JWT token còn hạn hay không dựa trên payload `exp`.
 * Thêm buffer 30 giây để tránh gửi request ngay sát thời điểm hết hạn.
 * @param {string} token - JWT string
 * @returns {boolean} true nếu token đã hết hạn (hoặc không hợp lệ)
 */
const isTokenExpired = (token) => {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    if (!payload?.exp) return false; // Không có exp → coi như còn hạn
    return Date.now() >= (payload.exp - 30) * 1000; // buffer 30 giây
  } catch {
    return false; // Parse lỗi → để server xử lý
  }
};

/**
 * Thực hiện refresh token và cập nhật localStorage.
 * Trả về access_token mới hoặc throw nếu thất bại.
 */
const doRefresh = async () => {
  const refreshToken = localStorage.getItem('refresh_token');
  if (!refreshToken) throw new Error('Không có refresh token');

  const res = await axios.post(`${import.meta.env.VITE_API_URL}/auth/refresh`, {
    refresh_token: refreshToken,
  });

  if (!res.data?.success) throw new Error('Làm mới token thất bại');

  const { access_token, refresh_token: newRefreshToken, user } = res.data.data;
  localStorage.setItem('access_token', access_token);
  localStorage.setItem('refresh_token', newRefreshToken);
  localStorage.setItem('user', JSON.stringify(user));
  return access_token;
};

/**
 * Gán Authorization header cho config (hỗ trợ cả Axios AxiosHeaders lẫn plain object).
 */
const setAuthHeader = (config, token) => {
  if (config.headers && typeof config.headers.set === 'function') {
    config.headers.set('Authorization', `Bearer ${token}`);
  } else {
    config.headers = config.headers || {};
    config.headers['Authorization'] = `Bearer ${token}`;
  }
};

// Request Interceptor: Gắn token và kiểm tra hết hạn trước khi request rời đi.
// Nếu token đã hết hạn → chặn request, refresh proactively, rồi mới gửi.
// Điều này tránh "Thundering Herd": nhiều request đồng thời cùng nhận 401 từ server.
api.interceptors.request.use(async (config) => {
  const isRefreshEndpoint = config.url?.includes('/auth/refresh');
  const token = localStorage.getItem('access_token');

  if (!token || isRefreshEndpoint) {
    return config;
  }

  // Token còn hạn → gắn bình thường
  if (!isTokenExpired(token)) {
    setAuthHeader(config, token);
    return config;
  }

  // Token hết hạn → refresh proactively
  if (isRefreshing) {
    // Đã có refresh đang chạy → xếp hàng chờ token mới
    const newToken = await new Promise((resolve, reject) => {
      failedQueue.push({ resolve, reject });
    });
    setAuthHeader(config, newToken);
    return config;
  }

  isRefreshing = true;
  try {
    const newToken = await doRefresh();
    processQueue(null, newToken);
    setAuthHeader(config, newToken);
    return config;
  } catch (err) {
    processQueue(err, null);
    clearAuthData();
    return Promise.reject(err);
  } finally {
    isRefreshing = false;
  }
});

// Response Interceptor: Chuẩn hóa dữ liệu trả về và xử lý lỗi tập trung
api.interceptors.response.use(
  (response) => response.data,
  async (error) => {
    // Trích xuất message lỗi từ các cấu trúc response khác nhau.
    let message = error.message || 'Lỗi kết nối';

    if (error.response?.data) {
      const d = error.response.data;

      if (d.message) {
        message = d.message;
      } else if (d.error) {
        message =
          typeof d.error === 'string'
            ? d.error
            : d.error.message || JSON.stringify(d.error);
      }
    } else if (error.response?.statusText) {
      message = error.response.statusText;
    }

    const status = error.response?.status;
    const originalRequest = error.config;
    const isRefreshEndpoint = originalRequest?.url?.includes('/auth/refresh');

    // Xử lý 401 fallback: token hợp lệ khi gửi nhưng server từ chối (edge case).
    // Trường hợp phổ biến đã được xử lý proactively ở request interceptor.
    if (status === 401 && originalRequest && !originalRequest._retry && !isRefreshEndpoint) {
      if (!localStorage.getItem('refresh_token')) {
        clearAuthData();
        return Promise.reject(new Error(message));
      }

      // Đang có request refresh khác chạy → xếp hàng chờ.
      if (isRefreshing) {
        try {
          const token = await new Promise((resolve, reject) => {
            failedQueue.push({ resolve, reject });
          });
          setAuthHeader(originalRequest, token);
          return api(originalRequest);
        } catch (err) {
          return Promise.reject(err);
        }
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        const newToken = await doRefresh();
        processQueue(null, newToken);
        setAuthHeader(originalRequest, newToken);
        return api(originalRequest);
      } catch (err) {
        processQueue(err, null);
        clearAuthData();
        return Promise.reject(err);
      } finally {
        isRefreshing = false;
      }
    }

    // Chuẩn hóa error object để callers có thể dùng `err.message` và `err.status`.
    const normalized = new Error(message);
    normalized.status = status;
    normalized.response = error.response;

    return Promise.reject(normalized);
  },
);

export default api;
