<template>
	<aside class="sidebar">sidebar</aside>
</template>
<script setup></script>
<style scoped>
.sidebar {
	grid-area: sidebar;
	background: #1e293b; /* Sử dụng text-main hoặc màu tối cụ thể */
	color: white;
	padding: var(--space-3);
	border-right: 1px solid var(--border-color);
}
</style>
