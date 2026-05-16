import { createRouter, createWebHistory } from "vue-router";

import MainLayout from "@/layout/MainLayout.vue";
import LoginView from "@/views/LoginView.vue";
import Dashboard from "@/views/Dashboard.vue";
import EmployeeView from "@/views/EmployeeView.vue";

const router = createRouter({
	history: createWebHistory(),

	routes: [
		{
			path: "/login",
			name: "login",
			component: LoginView,
			meta: { public: true },
		},
		{
			path: "/",
			component: MainLayout,
			redirect: "/dashboard",
			children: [
				{
					path: "dashboard",
					name: "dashboard",
					component: Dashboard,
				},
				{
					path: "employees",
					name: "employees",
					component: EmployeeView,
				},
			],
		},
	],
});

router.beforeEach((to, from, next) => {
	const token = localStorage.getItem("access_token");

	// Đã đăng nhập nhưng cố vào trang login -> chặn, chuyển hướng về dashboard
	if (to.name === "login" && token) {
		return next({ name: "dashboard" });
	}

	// Chưa đăng nhập mà vào các trang cần quyền (không có meta.public) -> đá ra trang login
	if (!to.meta.public && !token) {
		return next({ name: "login" });
	}

	// Hợp lệ -> cho phép đi tiếp
	next();
});

export default router;
