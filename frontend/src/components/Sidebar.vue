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
		<div class="sidebar-profile">
			<div class="profile-details">
				<div class="profile-avatar-circle">
					{{ userName ? userName.charAt(0).toUpperCase() : "U" }}
				</div>
				<div class="profile-info">
					<span class="profile-name">{{ userName }}</span>
					<span class="profile-email">{{ userEmail }}</span>
				</div>
			</div>
			<button
				class="logout-btn"
				@click="handleLogout"
				title="Đăng xuất"
				aria-label="Đăng xuất"
			>
				<LogOut class="logout-icon-svg" />
			</button>
		</div>
	</aside>
</template>

<script setup>
import { computed } from "vue";
import { useAuthStore } from "@/store/auth";
import { LayoutDashboard, Users, Building2, User, LogOut } from "@lucide/vue";
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

const handleLogout = () => {
	authStore.logout();
};
</script>

<style scoped>
.sidebar {
	width: var(--sidebar-width);
	background: var(--bg-card);
	backdrop-filter: var(--glass-backdrop);
	-webkit-backdrop-filter: var(--glass-backdrop);
	border-right: 1px solid var(--border-color);
	display: flex;
	flex-direction: column;
	transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
	z-index: 50;
	flex-shrink: 0;
	height: 100%;
}

.sidebar-brand {
	height: var(--header-height);
	display: flex;
	align-items: center;
	gap: var(--space-2);
	padding: 0 var(--space-3);
	border-bottom: 1px solid var(--border-color);
	flex-shrink: 0;
}

.logo-icon {
	display: flex;
	align-items: center;
	justify-content: center;
	width: 40px;
	height: 40px;
	background-color: var(--primary-glow);
	border-radius: var(--radius-md);
}

.logo-svg {
	width: 20px;
	height: 20px;
	color: var(--primary-color);
}

.brand-name {
	font-family: var(--font-title);
	font-size: var(--fs-lg);
	font-weight: var(--fw-bold);
	color: var(--text-main);
	letter-spacing: var(--tracking-tight);
	margin: 0;
}

.sidebar-nav {
	padding: var(--space-4) var(--space-3);
	display: flex;
	flex-direction: column;
	gap: var(--space-2);
	overflow-y: auto;
	flex: 1;
}

.nav-item {
	position: relative;
	display: flex;
	align-items: center;
	gap: var(--space-3);
	padding: var(--space-2) var(--space-3);
	color: var(--text-muted);
	border-radius: var(--radius-md);
	font-weight: var(--fw-medium);
	transition: all 0.2s ease;
	text-decoration: none;
	overflow: hidden;
}

.nav-icon {
	width: 18px;
	height: 18px;
	color: currentColor;
	transition: all 0.2s ease;
	z-index: 1;
}

.nav-item span {
	z-index: 1;
	font-size: var(--fs-sm);
}

.nav-item:hover {
	background: var(--bg-light);
	color: var(--text-main);
}

.nav-item.active {
	background: linear-gradient(
		90deg,
		rgba(66, 97, 237, 0.12) 0%,
		rgba(0, 192, 250, 0.04) 100%
	);
	color: var(--primary-color);
	font-weight: var(--fw-semibold);
}

.nav-item.active::before {
	content: "";
	position: absolute;
	left: 0;
	top: 15%;
	bottom: 15%;
	width: 4px;
	background: var(--primary-gradient);
	border-radius: 0 4px 4px 0;
	box-shadow: 2px 0 8px rgba(66, 97, 237, 0.3);
}

/* Profile section pushed to the bottom of the sidebar */
.sidebar-profile {
	margin-top: auto;
	padding: var(--space-3);
	border-top: 1px solid var(--border-color);
	display: flex;
	align-items: center;
	justify-content: space-between;
	gap: var(--space-2);
	flex-shrink: 0;
}

.profile-details {
	display: flex;
	align-items: center;
	gap: var(--space-3);
	overflow: hidden;
	flex: 1;
}

.profile-avatar-circle {
	width: 36px;
	height: 36px;
	border-radius: var(--radius-md);
	background: var(--primary-gradient);
	color: white;
	display: flex;
	align-items: center;
	justify-content: center;
	font-weight: var(--fw-bold);
	font-size: var(--fs-base);
	box-shadow: 0 4px 10px rgba(66, 97, 237, 0.2);
	flex-shrink: 0;
}

.profile-info {
	display: flex;
	flex-direction: column;
	overflow: hidden;
}

.profile-name {
	font-size: var(--fs-sm);
	font-weight: var(--fw-semibold);
	color: var(--text-main);
	line-height: var(--lh-tight);
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
	margin-bottom: 2px;
}

.profile-email {
	font-size: var(--fs-xs);
	color: var(--text-muted);
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

.logout-btn {
	display: flex;
	align-items: center;
	justify-content: center;
	width: 34px;
	height: 34px;
	background: var(--bg-card);
	border: 1px solid var(--border-color);
	border-radius: var(--radius-md);
	cursor: pointer;
	transition: all 0.2s ease;
	flex-shrink: 0;
	box-shadow: var(--shadow-sm);
}

.logout-btn:hover {
	background: #fee2e2;
	border-color: #fecaca;
	color: var(--danger-color);
}

.logout-icon-svg {
	width: 16px;
	height: 16px;
	color: currentColor;
	transition: all 0.2s ease;
}

@media (max-width: 1024px) {
	.sidebar {
		position: fixed;
		top: 0;
		bottom: 0;
		left: 0;
		transform: translateX(-100%);
		background: var(--bg-card);
		box-shadow: var(--shadow-lg);
	}

	.sidebar.is-open {
		transform: translateX(0);
	}
}
</style>
