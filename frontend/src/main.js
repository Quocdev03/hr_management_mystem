import "./assets/css/reset.css";
import "./assets/css/main.css";
import "vue-toastification/dist/index.css";

import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "./App.vue";
import Toast, { useToast } from "vue-toastification";
import router from "./router";

const app = createApp(App);

app.use(createPinia());
app.use(router);
app.use(Toast);

/**
 * Global Vue error handler — bắt mọi lỗi runtime từ component tree:
 * template render errors, lifecycle hook errors, v.v.
 * Thay vì để trang trắng xóa (White Screen of Death), hiện toast thông báo
 * và log chi tiết ra console để developer dễ debug.
 */
app.config.errorHandler = (err, instance, info) => {
	const isDev = import.meta.env.DEV;
	if (isDev) {
		console.error("[Vue Error]", { error: err, component: instance, info });
	}
	try {
		const toast = useToast();
		toast.error("Đã xảy ra lỗi không mong muốn. Vui lòng thử lại.");
	} catch {
		// Toast chưa khởi tạo — fallback về alert đơn giản
		console.error("Lỗi ứng dụng:", err?.message || err);
	}
};

app.mount("#app");
