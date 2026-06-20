<template>
	<header class="header">
		<button class="menu-btn" @click="$emit('toggle')" aria-label="Toggle Menu">
			<Menu class="menu-icon-svg" />
		</button>
		<span class="mobile-brand">HR System</span>
	</header>
</template>

<script setup>
import { Menu } from "@lucide/vue";

defineEmits(["toggle"]);
</script>

<style scoped>
.header {
	display: none;
	align-items: center;
	gap: var(--space-3);
	height: var(--header-height-sm, 56px);
	padding:var(--space-4) var(--space-3);
	background: rgba(255, 255, 255, 0.65);
	backdrop-filter: saturate(180%) blur(16px);
	-webkit-backdrop-filter: saturate(180%) blur(16px);
	border-bottom: var(--glass-border);
	z-index: 50;
	position: sticky;
	top: 0;
}

.menu-btn {
	background: transparent;
	border: none;
	cursor: pointer;
	padding: var(--space-1);
	border-radius: var(--radius-sm);
	display: flex;
	align-items: center;
	justify-content: center;
}

.menu-icon-svg {
	width: 24px;
	height: 24px;
	color: var(--text-muted);
}

.mobile-brand {
	font-size: var(--fs-base);
	font-weight: var(--fw-bold);
	color: var(--text-main);
	letter-spacing: var(--tracking-tight);
}

@media (max-width: 1024px) {
	.header {
		display: flex; /* Shown only on tablet and mobile */
	}
}
</style>
