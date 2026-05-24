import { createRouter, createWebHistory } from "vue-router";

import MainLayout from "@/layout/MainLayout.vue";
import LoginView from "@/views/LoginView.vue";
import Dashboard from "@/views/DashboardView.vue";
import EmployeeView from "@/views/EmployeeView.vue";
import DepartmentView from "@/views/DepartmentView.vue";
import ProfileView from "@/views/ProfileView.vue";
import UserView from "@/views/UserView.vue";
const router = createRouter({
   history: createWebHistory(),

   routes: [
      { path: "/login", name: "login", component: LoginView, meta: { public: true } },
      {
         path: "/",
         component: MainLayout,
         redirect: "/dashboard",
         children: [
            { path: "dashboard", name: "dashboard", component: Dashboard },
            { path: "employees", name: "employees", component: EmployeeView },
            { path: "departments", name: "departments", component: DepartmentView },
            { path: "users", name: "users", component: UserView },
            { path: "profile", name: "profile", component: ProfileView },
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

   // Hợp lệ -> cho phép đi tiếp
   return true;
});

export default router;
