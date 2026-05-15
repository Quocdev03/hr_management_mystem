<script setup>
import Header from "@/components/Header.vue";
import Sidebar from "@/components/Sidebar.vue";
</script>
<template>
	<div class="layout">
		<Header></Header>
		<Sidebar></Sidebar>
		<main class="content">
			<router-view />
		</main>
	</div>
</template>

<style scoped>
.layout {
	display: grid;
	grid-template-columns: var(--sidebar-width) 1fr;
	grid-template-rows: var(--header-height) 1fr;
	grid-template-areas:
		"sidebar header"
		"sidebar content";
	height: 100vh;
}

.content {
	grid-area: content;
	padding: var(--container-padding);
	overflow-y: auto;
	background: var(--bg-main);
}
</style>
