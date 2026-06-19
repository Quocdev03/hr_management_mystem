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
					meta: { permissions: ["department.read"] },
				},
				{
					path: "employees",
					name: "employees",
					component: EmployeeView,
					meta: { permissions: ["employee.read"] },
				},
				{
					path: "users",
					name: "users",
					component: UserView,
					meta: { roles: ["admin"], permissions: ["user.read"] },
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

	let user = null;
	try {
		user = JSON.parse(localStorage.getItem("user") || "null");
	} catch (error) {
		user = null;
	}

	const roleName = user?.role?.name || "";
	const permissions = user?.permissions || [];

	const hasRequiredRole =
		!to.meta.roles ||
		to.meta.roles.length === 0 ||
		to.meta.roles.includes(roleName) ||
		to.meta.roles.includes("*");

	const hasRequiredPermission =
		!to.meta.permissions ||
		to.meta.permissions.length === 0 ||
		to.meta.permissions.some((code) => permissions.includes(code));

	// Cho phép khi thỏa mãn role HOẶC permission.
	// Điều này đảm bảo route không bị khóa bởi role-name cứng nếu quyền đã được gán riêng.
	if (!hasRequiredRole && !hasRequiredPermission) {
		return { name: "dashboard" };
	}

	return true;
});

export default router;
