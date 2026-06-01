import axios from "axios";
// --- Cấu hình Instance ---
const api = axios.create({
	baseURL: import.meta.env.VITE_API_URL,
	timeout: 10000,
});

/**
 * Request Interceptor: Tự động gắn Token vào Header
 */

api.interceptors.request.use((config) => {
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
	(response) => response.data,
	(error) => {
		// Normalize error to an Error with a useful message so callers can use `err.message`
		let message = error.message || "Lỗi kết nối";

		if (error.response && error.response.data) {
			const d = error.response.data;
			if (d.message) {
				message = d.message;
			} else if (d.error) {
				if (typeof d.error === "string") message = d.error;
				else if (d.error.message) message = d.error.message;
				else message = JSON.stringify(d.error);
			}
		} else if (error.response && error.response.statusText) {
			message = error.response.statusText;
		}

		const status = error.response?.status;

		if (status === 401) {
			localStorage.removeItem("access_token");
			if (window.location.pathname !== "/login") {
				window.location.href = "/login";
			}
		}

		const normalized = new Error(message);
		normalized.status = status;
		normalized.response = error.response;

		return Promise.reject(normalized);
	},
);

export default api;
