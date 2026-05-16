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
		<Header @toggle="toggleSidebar"></Header>
		
		<div class="body-wrapper">
			<div
				v-if="isSidebarOpen"
				class="sidebar-backdrop"
				@click="closeSidebar"
			></div>

			<Sidebar :is-open="isSidebarOpen" @close="closeSidebar"></Sidebar>
			
			<main class="content">
				<router-view />
			</main>
		</div>
	</div>
</template>

<style scoped>
.layout {
	display: flex;
	flex-direction: column;
	height: 100vh;
	width: 100vw;
	overflow: hidden;
	background: var(--bg-main);
}

.body-wrapper {
	display: flex;
	flex: 1;
	overflow: hidden;
	position: relative;
}

.content {
	flex: 1;
	padding: var(--container-padding);
	overflow-y: auto;
	background: var(--bg-main);
}

.sidebar-backdrop {
	display: none;
	position: absolute;
	inset: 0;
	background: rgba(15, 23, 42, 0.6);
	backdrop-filter: blur(4px);
	z-index: 40;
}

@media (max-width: 1024px) {
	.sidebar-backdrop {
		display: block;
	}
}
</style>
