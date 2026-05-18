<template>
	<header class="header">
		<div class="header-left">
			<button class="menu-btn" @click="$emit('toggle')">
				<img :src="menuIcon" alt="menu" />
			</button>
			<div class="logo-section">
				<div class="logo-icon">
					<img :src="usersIcon" alt="logo" class="logo-img" />
				</div>
				<h2 class="brand-name">HR System</h2>
			</div>
		</div>

		<div class="header-right">
			<div class="user-profile">
				<div class="avatar-wrapper">
					<img :src="userIcon" alt="avatar" class="avatar-img" />
				</div>
				<div class="user-info">
					<span class="user-name">{{ userName }}</span>
					<span class="user-email">{{ userEmail }}</span>
				</div>
			</div>
			<button class="logout-btn" @click="handleLogout" title="Đăng xuất">
				<img :src="logoutIcon" alt="logout" />
			</button>
		</div>
	</header>
</template>

<script setup>
import { computed } from "vue";
import { useAuthStore } from "@/store/auth";
import menuIcon from "@/assets/svg/menu.svg";
import usersIcon from "@/assets/svg/users.svg";
import userIcon from "@/assets/svg/user.svg";
import logoutIcon from "@/assets/svg/log-out.svg";

const authStore = useAuthStore();
const userData = localStorage.getItem("user");
const userDataParse = JSON.parse(userData);
const userEmail = userDataParse.email;
const userName = userDataParse.user_name;

const handleLogout = () => {
	authStore.logout();
};

defineEmits(["toggle"]);
</script>

<style scoped>
.header {
	display: flex;
	align-items: center;
	justify-content: space-between;
	height: var(--header-height);
	padding: 0 var(--space-4);
	background: var(--bg-card);
	backdrop-filter: blur(8px);
	-webkit-backdrop-filter: blur(8px);
	border-bottom: 1px solid var(--border-color);
	z-index: 10;
	flex-shrink: 0;
}

.header-left {
	display: flex;
	align-items: center;
	gap: var(--space-3);
}

.menu-btn {
	display: none;
	background: transparent;
	border: none;
	cursor: pointer;
	padding: var(--space-1);
	border-radius: var(--radius-sm);
	transition: background-color 0.2s;
}

.menu-btn:hover {
	background: var(--bg-light);
}

.menu-btn img {
	width: 24px;
	height: 24px;
	filter: invert(45%) sepia(12%) saturate(545%) hue-rotate(182deg)
		brightness(92%) contrast(89%);
}

.logo-section {
	display: flex;
	align-items: center;
	gap: var(--space-2);
}

.logo-icon {
	display: flex;
	align-items: center;
	justify-content: center;
	width: 40px;
	height: 40px;
	background-color: #eff6ff;
	border-radius: var(--radius-md);
}

.logo-img {
	width: 24px;
	height: 24px;
	filter: invert(44%) sepia(87%) saturate(2258%) hue-rotate(200deg)
		brightness(101%) contrast(92%);
}

.brand-name {
	font-size: var(--fs-lg);
	font-weight: var(--fw-bold);
	color: var(--text-main);
	letter-spacing: var(--tracking-tight);
	margin: 0;
}

.header-right {
	display: flex;
	align-items: center;
	gap: var(--space-4);
}

.user-profile {
	display: flex;
	align-items: center;
	gap: var(--space-2);
}

.avatar-wrapper {
	width: 40px;
	height: 40px;
	border-radius: var(--radius-full);
	background: #eff6ff;
	display: flex;
	align-items: center;
	justify-content: center;
	border: 2px solid white;
	box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
	overflow: hidden;
}

.avatar-img {
	width: 20px;
	height: 20px;
	filter: invert(44%) sepia(87%) saturate(2258%) hue-rotate(200deg)
		brightness(101%) contrast(92%);
}

.user-info {
	display: flex;
	flex-direction: column;
	justify-content: center;
}

.user-name {
	font-size: var(--fs-sm);
	font-weight: var(--fw-semibold);
	color: var(--text-main);
	line-height: var(--lh-tight);
}

.user-email {
	font-size: var(--fs-xs);
	color: var(--text-muted);
}

.logout-btn {
	display: flex;
	align-items: center;
	justify-content: center;
	width: 40px;
	height: 40px;
	background: transparent;
	border: 1px solid transparent;
	border-radius: var(--radius-md);
	cursor: pointer;
	transition: all 0.2s;
}

.logout-btn:hover {
	background: #fef2f2;
	border-color: #fecaca;
}

.logout-btn img {
	width: 20px;
	height: 20px;
	filter: invert(45%) sepia(12%) saturate(545%) hue-rotate(182deg)
		brightness(92%) contrast(89%);
	transition: filter 0.2s;
}

.logout-btn:hover img {
	filter: invert(36%) sepia(82%) saturate(2462%) hue-rotate(345deg)
		brightness(96%) contrast(94%);
}

@media (max-width: 1024px) {
	.menu-btn {
		display: block;
	}
	.brand-name {
		display: none;
	}
}

@media (max-width: 640px) {
	.user-info {
		display: none;
	}
	.header {
		padding: 0 var(--space-3);
	}
	.header-right {
		gap: var(--space-2);
	}
}
</style>
