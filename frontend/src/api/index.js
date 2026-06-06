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

// Request Interceptor: Tự động gắn Access Token vào Header
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');

  if (token) {
    if (config.headers && typeof config.headers.set === 'function') {
      config.headers.set('Authorization', `Bearer ${token}`);
    } else {
      config.headers = config.headers || {};
      config.headers['Authorization'] = `Bearer ${token}`;
    }
  }

  return config;
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

    // Xử lý 401: thử làm mới Access Token trước khi redirect.
    if (status === 401 && originalRequest && !originalRequest._retry && !isRefreshEndpoint) {
      const refreshToken = localStorage.getItem('refresh_token');

      // Không có refresh token → xóa session và redirect.
      if (!refreshToken) {
        clearAuthData();
        return Promise.reject(new Error(message));
      }

      // Đang có request refresh khác chạy → xếp hàng chờ.
      if (isRefreshing) {
        try {
          const token = await new Promise((resolve, reject) => {
            failedQueue.push({ resolve, reject });
          });
          if (originalRequest.headers && typeof originalRequest.headers.set === 'function') {
            originalRequest.headers.set('Authorization', `Bearer ${token}`);
          } else {
            originalRequest.headers = originalRequest.headers || {};
            originalRequest.headers['Authorization'] = `Bearer ${token}`;
          }
          return api(originalRequest);
        } catch (err) {
          return Promise.reject(err);
        }
      }

      // Bắt đầu refresh token.
      originalRequest._retry = true;
      isRefreshing = true;

      try {
        const res = await axios.post(`${import.meta.env.VITE_API_URL}/auth/refresh`, {
          refresh_token: refreshToken,
        });

        if (!res.data?.success) {
          throw new Error('Làm mới token thất bại');
        }

        const { access_token, refresh_token: newRefreshToken, user } = res.data.data;
        localStorage.setItem('access_token', access_token);
        localStorage.setItem('refresh_token', newRefreshToken);
        localStorage.setItem('user', JSON.stringify(user));

        processQueue(null, access_token);
        if (originalRequest.headers && typeof originalRequest.headers.set === 'function') {
          originalRequest.headers.set('Authorization', `Bearer ${access_token}`);
        } else {
          originalRequest.headers = originalRequest.headers || {};
          originalRequest.headers['Authorization'] = `Bearer ${access_token}`;
        }
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
