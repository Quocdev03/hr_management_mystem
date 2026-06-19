<script setup>
import { ref } from "vue";
import Header from "@/components/Header.vue";
import Sidebar from "@/components/Sidebar.vue";

const isSidebarOpen = ref(false);

const toggleSidebar = () => {
	isSidebarOpen.value = !isSidebarOpen.value;
};

const closeSidebar = () => {
	isSidebarOpen.value = false;
};
</script>

<template>
	<div class="layout" :class="{ 'sidebar-open': isSidebarOpen }">
		<!-- Sidebar on the left (full height) -->
		<Sidebar :is-open="isSidebarOpen" @close="closeSidebar"></Sidebar>
		
		<div class="main-wrapper">
			<!-- Header above the main content -->
			<Header @toggle="toggleSidebar"></Header>
			
			<main class="content">
				<router-view />
			</main>
		</div>

		<!-- Mobile Backdrop -->
		<div
			v-if="isSidebarOpen"
			class="sidebar-backdrop"
			@click="closeSidebar"
		></div>
	</div>
</template>

<style scoped>
.layout {
	display: flex;
	flex-direction: row; /* Sidebar left, Content area right */
	height: 100vh;
	width: 100vw;
	overflow: hidden;
	position: relative;
	background: transparent;
}

.main-wrapper {
	display: flex;
	flex-direction: column;
	flex: 1;
	overflow: hidden;
	height: 100%;
}

.content {
	flex: 1;
	padding: var(--container-padding);
	overflow-y: auto;
	background: var(--bg-main);
	-webkit-overflow-scrolling: touch;
	will-change: scroll-position;
}

.sidebar-backdrop {
	display: none;
	position: fixed;
	inset: 0;
	background: rgba(15, 23, 42, 0.6);
	z-index: 40;
}

@media (max-width: 1024px) {
	.sidebar-backdrop {
		display: block;
	}
}
</style>
