<template>
	<Teleport to="body">
		<Transition name="modal-fade">
			<div v-if="visible" class="modal-overlay" @click.self="$emit('close')">
				<div class="modal-container" :class="[size]" :style="customStyle">
					<!-- Header -->
					<div class="modal-header">
						<div class="header-content">
							<h3 class="modal-title">{{ title }}</h3>
							<p v-if="subtitle" class="modal-subtitle">
								{{ subtitle }}
							</p>
						</div>
						<button class="close-btn" @click="$emit('close')">
							<img :src="closeIcon" alt="close" width="24" height="24" />
						</button>
					</div>

					<!-- Body -->
					<div class="modal-body">
						<slot></slot>
					</div>

					<!-- Footer -->
					<div v-if="$slots.footer" class="modal-footer">
						<slot name="footer"></slot>
					</div>
				</div>
			</div>
		</Transition>
	</Teleport>
</template>

<script setup>
import { computed } from "vue";
import closeIcon from "@/assets/svg/close.svg";

const props = defineProps({
	visible: Boolean,
	title: String,
	subtitle: String,
	size: {
		type: String,
		default: "", // sm, md, lg, xl
	},
	width: String,
});

defineEmits(["close"]);

const customStyle = computed(() => {
	if (props.width) {
		return { maxWidth: props.width };
	}
	return {};
});
</script>

<style scoped>
.modal-overlay {
	position: fixed;
	inset: 0;
	background: rgba(15, 23, 42, 0.55);
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 1000;
	padding: var(--space-4);
	touch-action: none;
}

.modal-container {
	background: var(--bg-card);
	border-radius: var(--radius-xl);
	width: 100%;
	max-height: 90vh;
	box-shadow: var(--shadow-modal);
	display: flex;
	flex-direction: column;
	overflow: hidden;
	animation: slideUp 0.3s ease-out;
	transform: translateZ(0);
	will-change: transform;
}

/* Kích thước chuẩn cho Modal */
.modal-container.sm {
	max-width: 440px;
}
.modal-container.md {
	max-width: 680px;
}
.modal-container.lg {
	max-width: 960px;
}
.modal-container.xl {
	max-width: 1200px;
}

.modal-header {
	padding: var(--space-3) var(--space-4);
	border-bottom: 1px solid var(--border-color);
	display: flex;
	justify-content: space-between;
	align-items: center;
}

.modal-title {
	font-size: var(--fs-xl); /* 20px-22px */
	font-weight: var(--fw-bold);
	color: var(--text-main);
	margin: 0;
}

.modal-subtitle {
	font-size: var(--fs-sm); /* 14px-15px */
	color: var(--text-muted);
	margin: 4px 0 0 0;
}

.close-btn {
	background: transparent;
	border: none;
	color: var(--text-light);
	cursor: pointer;
	padding: 8px;
	border-radius: var(--radius-md);
	transition: all 0.2s;
	display: flex;
	align-items: center;
	justify-content: center;
}

.close-btn:hover {
	background: var(--bg-light);
	color: var(--text-main);
}

.modal-body {
	padding: var(--space-4);
	overflow-y: auto;
	-webkit-overflow-scrolling: touch;
	overscroll-behavior: contain;
	touch-action: pan-y;
	font-size: var(--fs-base); /* 16px-17px */
}

.modal-footer {
	padding: var(--space-3) var(--space-4);
	border-top: 1px solid var(--border-color);
	background: var(--bg-lighter);
	display: flex;
	justify-content: flex-end;
	gap: var(--space-2);
}

/* Animations */
.modal-fade-enter-active,
.modal-fade-leave-active {
	transition: opacity 0.3s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
	opacity: 0;
}

@keyframes slideUp {
	from {
		transform: translateY(20px);
		opacity: 0;
	}
	to {
		transform: translateY(0);
		opacity: 1;
	}
}

@media (max-width: 640px) {
	.modal-overlay {
		padding: var(--space-2);
	}
	.modal-container {
		max-height: 95vh;
		border-radius: var(--radius-lg);
	}
}
</style>
