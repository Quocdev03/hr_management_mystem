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
	background: var(--bg-card);
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
	-webkit-overflow-scrolling: touch;
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
	transition:
		background-color 0.2s ease,
		border-color 0.2s ease,
		color 0.2s ease,
		opacity 0.2s ease,
		transform 0.2s ease,
		box-shadow 0.2s ease;
	text-decoration: none;
	overflow: hidden;
}

.nav-icon {
	width: 20px;
	height: 20px;
	color: currentColor;
	transition:
		background-color 0.2s ease,
		border-color 0.2s ease,
		color 0.2s ease,
		opacity 0.2s ease,
		transform 0.2s ease,
		box-shadow 0.2s ease;
	z-index: 1;
}

.nav-item span {
	z-index: 1;
	font-size: var(--fs-base);
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

.nav-item.active .nav-icon {
	color: var(--primary-color);
	filter: drop-shadow(0 2px 8px rgba(66, 97, 237, 0.4));
	transform: scale(1.1);
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

/* ── Sidebar Footer ────────────────────────────────── */
.sidebar-footer {
	margin-top: auto;
	border-top: 1px solid var(--border-color);
	background: var(--bg-lighter);
	flex-shrink: 0;
	padding: var(--space-2) var(--space-3);
	display: flex;
	flex-direction: column;
	gap: var(--space-2);
}

.sidebar-profile {
	display: flex;
	align-items: center;
	gap: 10px;
	padding: 6px 8px;
	border-radius: var(--radius-md);
	transition: background-color 0.2s ease;
	cursor: default;
}

.sidebar-profile:hover {
	background: rgba(66, 97, 237, 0.05);
}

.profile-avatar-circle {
	width: 40px;
	height: 40px;
	border-radius: var(--radius-md);
	background: var(--primary-gradient);
	color: white;
	display: flex;
	align-items: center;
	justify-content: center;
	font-weight: var(--fw-bold);
	font-size: var(--fs-base);
	box-shadow:
		0 0 0 2px var(--bg-lighter),
		0 0 0 3px rgba(66, 97, 237, 0.25);
	flex-shrink: 0;
	transition: box-shadow 0.2s ease;
}

.sidebar-profile:hover .profile-avatar-circle {
	box-shadow:
		0 0 0 2px var(--bg-lighter),
		0 0 0 3px rgba(66, 97, 237, 0.5);
}

.profile-info {
	display: flex;
	flex-direction: column;
	overflow: hidden;
	flex: 1;
	gap: 2px;
}

.profile-name {
	font-size: var(--fs-sm);
	font-weight: var(--fw-semibold);
	color: var(--text-main);
	line-height: 1.2;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

.profile-role {
	font-size: 11px;
	font-weight: var(--fw-medium);
	color: var(--primary-color);
	background: rgba(66, 97, 237, 0.08);
	padding: 1px 6px;
	border-radius: var(--radius-full);
	width: fit-content;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
	max-width: 100%;
}

.profile-email {
	font-size: 11px;
	color: var(--text-muted);
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

.logout-btn {
	display: flex;
	align-items: center;
	justify-content: center;
	gap: 8px;
	width: 100%;
	height: 38px;
	border: 1px solid rgba(225, 29, 72, 0.2);
	border-radius: var(--radius-md);
	cursor: pointer;
	background: rgba(225, 29, 72, 0.05);
	color: var(--danger-color);
	font-size: var(--fs-sm);
	font-weight: var(--fw-medium);
	transition: all 0.2s ease;
}

.logout-btn:hover {
	background: rgba(225, 29, 72, 0.1);
	border-color: rgba(225, 29, 72, 0.35);
	transform: translateY(-1px);
	box-shadow: 0 4px 10px rgba(225, 29, 72, 0.12);
}

.logout-btn:active {
	transform: translateY(0);
}

.logout-icon-svg {
	width: 15px;
	height: 15px;
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
		background: var(--bg-card);
		box-shadow: var(--shadow-lg);
	}

	.sidebar.is-open {
		transform: translateX(0);
	}
}
</style>
