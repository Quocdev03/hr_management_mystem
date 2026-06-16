<template>
	<aside class="sidebar" :class="{ 'is-open': isOpen }">
		<nav class="sidebar-nav">
			<router-link to="/dashboard" class="nav-item" active-class="active">
				<img :src="gridIcon" alt="dashboard" class="nav-icon" />
				<span>Tổng quan</span>
			</router-link>
			<router-link to="/employees" class="nav-item" active-class="active">
				<img :src="usersIcon" alt="employees" class="nav-icon" />
				<span>Nhân viên</span>
			</router-link>
			<router-link to="/departments" class="nav-item" active-class="active">
				<img :src="departmentIcon" alt="departments" class="nav-icon" />
				<span>Phòng ban</span>
			</router-link>
			<router-link
				v-if="canManageUsers"
				to="/users"
				class="nav-item"
				active-class="active"
			>
				<img :src="usersIcon" alt="user" class="nav-icon" />
				<span>Người Dùng</span>
			</router-link>
			<router-link to="/me" class="nav-item" active-class="active">
				<img :src="profileIcon" alt="profile" class="nav-icon" />
				<span>Hồ Sơ</span>
			</router-link>
		</nav>
	</aside>
</template>

<script setup>
import gridIcon from "@/assets/svg/grid.svg";
import usersIcon from "@/assets/svg/users.svg";
import profileIcon from "@/assets/svg/user.svg";
import departmentIcon from "@/assets/svg/department.svg";
import { usePermissions } from "@/helpers/usePermissions";

defineProps({ isOpen: { type: Boolean, default: false } });

const { canManageUsers } = usePermissions();
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
	overflow-y: auto;
}

.sidebar-nav {
	padding: var(--space-4) var(--space-3);
	display: flex;
	flex-direction: column;
	gap: var(--space-2);
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
	width: 20px;
	height: 20px;
	filter: invert(45%) sepia(12%) saturate(545%) hue-rotate(182deg)
		brightness(92%) contrast(89%); /* var(--text-muted) */
	transition: all 0.2s ease;
	z-index: 1;
}

.nav-item span {
	z-index: 1;
}

.nav-item:hover {
	background: var(--bg-light);
	color: var(--text-main);
}

.nav-item:hover .nav-icon {
	filter: invert(13%) sepia(15%) saturate(1487%) hue-rotate(182deg)
		brightness(98%) contrast(92%); /* var(--text-main) */
}

.nav-item.active {
	background: linear-gradient(90deg, #eff6ff 0%, rgba(239, 246, 255, 0) 100%);
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
	background: var(--primary-color);
	border-radius: 0 4px 4px 0;
	box-shadow: 2px 0 8px rgba(59, 130, 246, 0.4);
}

.nav-item.active .nav-icon {
	filter: invert(44%) sepia(87%) saturate(2258%) hue-rotate(200deg)
		brightness(101%) contrast(92%); /* var(--primary-color) */
}

@media (max-width: 1024px) {
	.sidebar {
		position: absolute;
		top: 0;
		bottom: 0;
		left: 0;
		transform: translateX(-100%);
		background: rgba(255, 255, 255, 0.98);
		backdrop-filter: blur(8px);
		-webkit-backdrop-filter: blur(8px);
	}

	.sidebar.is-open {
		transform: translateX(0);
		box-shadow: var(--shadow-lg);
	}
}
</style>
