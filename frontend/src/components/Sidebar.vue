<template>
	<aside class="sidebar" :class="{ 'is-open': isOpen }">
		<!-- Brand Logo Section -->
		<div class="sidebar-brand">
			<div class="logo-icon">
				<Users class="logo-svg" />
			</div>
			<h2 class="brand-name">HR System</h2>
		</div>

		<!-- Nav links -->
		<nav class="sidebar-nav">
			<router-link to="/dashboard" class="nav-item" active-class="active">
				<LayoutDashboard class="nav-icon" />
				<span>Tổng quan</span>
			</router-link>
			<router-link
				v-if="canViewEmployeeList"
				to="/employees"
				class="nav-item"
				active-class="active"
			>
				<Users class="nav-icon" />
				<span>Nhân viên</span>
			</router-link>
			<router-link
				v-if="canViewDepartmentList"
				to="/departments"
				class="nav-item"
				active-class="active"
			>
				<Building2 class="nav-icon" />
				<span>Phòng ban</span>
			</router-link>
			<router-link
				v-if="canViewDepartmentList"
				to="/positions"
				class="nav-item"
				active-class="active"
			>
				<Briefcase class="nav-icon" />
				<span>Chức vụ</span>
			</router-link>
			<router-link
				v-if="canManageUsers"
				to="/roles"
				class="nav-item"
				active-class="active"
			>
				<ShieldCheck class="nav-icon" />
				<span>Vai trò</span>
			</router-link>
			<router-link
				v-if="canManageUsers"
				to="/users"
				class="nav-item"
				active-class="active"
			>
				<Users class="nav-icon" />
				<span>Người Dùng</span>
			</router-link>
			<router-link to="/me" class="nav-item" active-class="active">
				<User class="nav-icon" />
				<span>Hồ Sơ</span>
			</router-link>
		</nav>

		<!-- Profile Section at bottom of Sidebar -->
		<div class="sidebar-footer">
			<div class="sidebar-profile">
				<div class="profile-avatar-circle">
					{{ userName ? userName.charAt(0).toUpperCase() : "U" }}
				</div>
				<div class="profile-info">
					<span class="profile-name" :title="userName">{{ userName }}</span>
					<span v-if="userRole" class="profile-role">{{ userRole }}</span>
					<span v-else class="profile-email" :title="userEmail">{{ userEmail }}</span>
				</div>
			</div>
			<button
				class="logout-btn"
				@click="handleLogout"
				aria-label="Đăng xuất"
			>
				<LogOut class="logout-icon-svg" />
				<span>Đăng xuất</span>
			</button>
		</div>
	</aside>
</template>

<script setup>
import { computed, onMounted } from "vue";
import { useAuthStore } from "@/store/auth";
import {
	LayoutDashboard,
	Users,
	Building2,
	User,
	LogOut,
	Briefcase,
	ShieldCheck,
} from "@lucide/vue";
import { usePermissions } from "@/helpers/usePermissions";

defineProps({ isOpen: { type: Boolean, default: false } });

const { canManageUsers, canViewEmployeeList, canViewDepartmentList } =
	usePermissions();
const authStore = useAuthStore();

// Safely parse local storage user data
let userFromStorage = null;
try {
	const raw = localStorage.getItem("user");
	userFromStorage = raw ? JSON.parse(raw) : null;
} catch (err) {
	userFromStorage = null;
}

const userEmail = computed(
	() => userFromStorage?.email ?? authStore.user?.email ?? "",
);
const userName = computed(
	() => userFromStorage?.user_name ?? authStore.user?.user_name ?? "",
);
const userRole = computed(() => authStore.userProfile?.role_name ?? "");

onMounted(() => {
	if (!authStore.userProfile) {
		authStore.profile();
	}
});

const handleLogout = () => {
	authStore.logout();
};
</script>

<style scoped>
.sidebar {
	width: var(--sidebar-width);
	background: #ffffff;
	border-right: 1px solid rgba(0, 0, 0, 0.05);
	display: flex;
	flex-direction: column;
	transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
	z-index: 50;
	flex-shrink: 0;
	height: 100%;
}

.sidebar-brand {
	height: 64px;
	display: flex;
	align-items: center;
	gap: 10px;
	padding: 0 16px;
	border-bottom: 1px solid rgba(0, 0, 0, 0.04);
	flex-shrink: 0;
}

.logo-icon {
	display: flex;
	align-items: center;
	justify-content: center;
	width: 32px;
	height: 32px;
	background: var(--primary-gradient);
	color: #ffffff;
	border-radius: 8px;
	box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08);
	transition: transform 0.3s ease;
}

.sidebar-brand:hover .logo-icon {
	transform: rotate(5deg) scale(1.05);
}

.logo-svg {
	width: 16px;
	height: 16px;
	color: #ffffff;
}

.brand-name {
	font-family: var(--font-title);
	font-size: 15px;
	font-weight: 700;
	color: #0f172a;
	letter-spacing: -0.02em;
	margin: 0;
}

.sidebar-nav {
	padding: 16px 12px;
	display: flex;
	flex-direction: column;
	gap: 4px;
	overflow-y: auto;
	flex: 1;
	-webkit-overflow-scrolling: touch;
}

/* Custom scrollbar for sidebar nav */
.sidebar-nav::-webkit-scrollbar {
	width: 4px;
}
.sidebar-nav::-webkit-scrollbar-track {
	background: transparent;
}
.sidebar-nav::-webkit-scrollbar-thumb {
	background: rgba(0, 0, 0, 0.1);
	border-radius: 99px;
}
.sidebar-nav::-webkit-scrollbar-thumb:hover {
	background: rgba(0, 0, 0, 0.2);
}

.nav-item {
	position: relative;
	display: flex;
	align-items: center;
	gap: 10px;
	padding: 15px 16px;
	color: #64748b;
	border-radius: 6px;
	font-weight: 500;
	transition: all 0.15s ease;
	text-decoration: none;
	overflow: hidden;
}

.nav-icon {
	width: 18px;
	height: 18px;
	color: currentColor;
	z-index: 1;
	transition: transform 0.15s ease;
}

.nav-item span {
	z-index: 1;
	font-size: 13.5px;
}

.nav-item:hover {
	background: #f1f5f9;
	color: #0f172a;
}

.nav-item:hover .nav-icon {
	transform: translateX(1px);
}

.nav-item.active {
	background: #eff6ff;
	color: #2563eb;
	font-weight: 600;
}

.nav-item.active .nav-icon {
	color: #2563eb;
}

/* Active indicator line */
.nav-item.active::before {
	content: '';
	position: absolute;
	left: 0;
	top: 50%;
	transform: translateY(-50%);
	width: 3px;
	height: 16px;
	background: #2563eb;
	border-radius: 0 4px 4px 0;
}

/* ── Sidebar Footer ────────────────────────────────── */
.sidebar-footer {
	margin-top: auto;
	border-top: 1px solid rgba(0, 0, 0, 0.04);
	background: #ffffff;
	flex-shrink: 0;
	padding: 16px 14px;
	display: flex;
	flex-direction: column;
	gap: 12px;
}

.sidebar-profile {
	display: flex;
	align-items: center;
	gap: 10px;
	padding: 6px;
	border-radius: 8px;
	transition: background-color 0.15s ease;
	cursor: default;
}

.sidebar-profile:hover {
	background: #f8fafc;
}

.profile-avatar-circle {
	width: 32px;
	height: 32px;
	border-radius: 8px;
	background: linear-gradient(135deg, #e0e7ff 0%, #c7d2fe 100%);
	color: #4f46e5;
	display: flex;
	align-items: center;
	justify-content: center;
	font-weight: 700;
	font-size: 12px;
	border: 1px solid rgba(79, 70, 229, 0.1);
	box-shadow: 0 2px 4px rgba(0, 0, 0, 0.02);
	flex-shrink: 0;
}

.profile-info {
	display: flex;
	flex-direction: column;
	overflow: hidden;
	flex: 1;
	gap: 2px;
}

.profile-name {
	font-size: 13px;
	font-weight: 600;
	color: #0f172a;
	line-height: 1.2;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

.profile-role {
	font-size: 10px;
	font-weight: 600;
	color: #2563eb;
	background: #eff6ff;
	padding: 1px 6px;
	border-radius: 9999px;
	width: fit-content;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
	max-width: 100%;
}

.profile-email {
	font-size: 10px;
	color: #64748b;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

.logout-btn {
	display: flex;
	align-items: center;
	justify-content: center;
	gap: 6px;
	width: 100%;
	height: 32px;
	border: 1px solid #fee2e2;
	border-radius: 6px;
	cursor: pointer;
	background: #fff5f5;
	color: #ef4444;
	font-size: 12px;
	font-weight: 600;
	transition: all 0.15s ease;
}

.logout-btn:hover {
	background: #fee2e2;
	border-color: #fca5a5;
	transform: translateY(-1px);
	box-shadow: 0 2px 6px rgba(225, 29, 72, 0.08);
}

.logout-btn:active {
	transform: translateY(0);
}

.logout-icon-svg {
	width: 13px;
	height: 13px;
	color: currentColor;
	flex-shrink: 0;
}

@media (max-width: 1024px) {
	.sidebar {
		position: fixed;
		top: 0;
		bottom: 0;
		left: 0;
		transform: translateX(-100%);
		background: #ffffff;
		box-shadow: var(--shadow-lg);
	}

	.sidebar.is-open {
		transform: translateX(0);
	}
}
</style>
