import axios from "axios";
// --- Cấu hình Instance ---
const api = axios.create({
	baseURL: import.meta.env.VITE_API_URL,
	timeout: 10000,
});

/**
 * Request Interceptor: Tự động gắn Token vào Header
 */

api.interceptors.request.use(config => {
	const access_token = localStorage.getItem("access_token");
	if (access_token) {
		config.headers.Authorization = `Bearer ${access_token}`;
	}
	return config;
});

/**
 * Response Interceptor: Xử lý dữ liệu trả về và bắt lỗi tập trung
 */
api.interceptors.response.use(
	response => {
		return response.data;
	},
	error => {
		if (error.response) {
			const { status, data } = error.response;

			switch (status) {
				case 401:
					// Lỗi chưa đăng nhập hoặc token hết hạn
					localStorage.removeItem("access_token");
					// Chuyển hướng về login nếu không phải đang ở trang login
					if (window.location.pathname !== "/login") {
						window.location.href = "/login";
					}
					break;
				case 403:
					console.error("Bạn không có quyền truy cập tính năng này");
					break;
				case 422:
					console.error("Dữ liệu không hợp lệ:", data.error);
					break;
				case 500:
					console.error("Lỗi hệ thống, vui lòng thử lại sau");
					break;
				default:
					console.error(data.message || "Đã có lỗi xảy ra");
			}
		} else {
			console.error("Không thể kết nối đến server");
		}

		// Vẫn reject để component gọi API có thể catch lỗi và hiển thị thông báo riêng nếu cần
		return Promise.reject(error);
	},
);

export default api;
