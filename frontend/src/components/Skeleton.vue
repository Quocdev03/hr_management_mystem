<script setup>
	defineProps({
		type: {
			type: String,
			default: "text", // 'text', 'avatar', 'badge', 'button'
		},
		width: {
			type: String,
			default: "100%",
		},
		height: {
			type: String,
			default: "",
		},
		borderRadius: {
			type: String,
			default: "",
		},
	});
</script>

<template>
	<div
		:class="['skeleton', `skeleton-${type}`]"
		:style="{
			width: width,
			height: height || undefined,
			borderRadius: borderRadius || undefined,
		}"
	></div>
</template>

<style scoped>
	/* Styles are defined in main.css so they can be reused globally, 
   but this component acts as a highly clean wrapper */
</style>
