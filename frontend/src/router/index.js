import { createRouter, createWebHistory } from "vue-router";

import MainLayout from "@/layout/MainLayout.vue";
import LoginView from "@/views/LoginView.vue";
import Dashboard from "@/views/Dashboard.vue";

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
			],
		},
	],
});

export default router;
