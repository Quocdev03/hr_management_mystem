import { createRouter, createWebHistory } from "vue-router";

const MainLayout = () => import("@/layout/MainLayout.vue");
const LoginView = () => import("@/views/LoginView.vue");
const Dashboard = () => import("@/views/DashboardView.vue");
const EmployeeView = () => import("@/views/EmployeeView.vue");
const DepartmentView = () => import("@/views/DepartmentView.vue");
const ProfileView = () => import("@/views/ProfileView.vue");
const UserView = () => import("@/views/UserView.vue");

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
				{ path: "dashboard", name: "dashboard", component: Dashboard },
				{
					path: "departments",
					name: "departments",
					component: DepartmentView,
				},
				{
					path: "employees",
					name: "employees",
					component: EmployeeView,
				},
				{
					path: "users",
					name: "users",
					component: UserView,
					meta: { roles: ["admin"] },
				},
				{ path: "me", name: "me", component: ProfileView },
			],
		},
	],
});

router.beforeEach((to) => {
	const token = localStorage.getItem("access_token");

	// Đã đăng nhập nhưng cố vào trang login -> chặn, chuyển hướng về dashboard
	if (to.name === "login" && token) {
		return { name: "dashboard" };
	}

	// Chưa đăng nhập mà vào các trang cần quyền (không có meta.public) -> đá ra trang login
	if (!to.meta.public && !token) {
		return { name: "login" };
	}

	// Kiểm tra phân quyền theo role (nếu route có meta.roles)
	if (to.meta.roles && to.meta.roles.length > 0) {
		const user = JSON.parse(localStorage.getItem("user") || "null");
		const roleName = user?.role?.name || "";
		if (!to.meta.roles.includes(roleName)) {
			return { name: "dashboard" };
		}
	}

	// Hợp lệ -> cho phép đi tiếp
	return true;
});

export default router;
