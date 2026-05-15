<template>
	<header class="header">header</header>
</template>
<script setup></script>
<style scoped>
.header {
	grid-area: header;
	display: flex;
	align-items: center;
	padding: 0 var(--space-3);
	border-bottom: 1px solid var(--border-color);
	background: var(--bg-card);
}
</style>
