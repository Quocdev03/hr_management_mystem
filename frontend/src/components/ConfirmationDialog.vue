<template>
	<Teleport to="body">
		<Transition name="modal-fade">
			<div
				v-if="visible"
				class="confirm-overlay"
				@click.self="$emit('cancel')"
			>
				<div class="confirm-container">
					<div class="confirm-body">
						<div class="warning-icon">
							<TriangleAlert  />
						</div>
						<h3 class="confirm-title">{{ title }}</h3>
						<p class="confirm-message">{{ message }}</p>
					</div>
					
					<div class="confirm-footer">
						<button class="btn btn-cancel" @click="$emit('cancel')">Hủy bỏ</button>
						<button class="btn btn-delete" @click="$emit('confirm')" :disabled="loading">
							<span v-if="loading" class="spinner"></span>
							<span>{{ confirmText }}</span>
						</button>
					</div>
				</div>
			</div>
		</Transition>
	</Teleport>
</template>

<script setup>
import { TriangleAlert } from '@lucide/vue';

const props = defineProps({
	visible: Boolean,
	title: {
		type: String,
		default: "Xác nhận xoá",
	},
	message: {
		type: String,
		default: "Bạn có chắc chắn muốn thực hiện hành động này?",
	},
	confirmText: {
		type: String,
		default: "Xác nhận xoá",
	},
	loading: Boolean,
});

const emit = defineEmits(["confirm", "cancel"]);
</script>

<style scoped>
.confirm-overlay {
	position: fixed;
	inset: 0;
	background: rgba(15, 23, 42, 0.25);
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 1000;
	padding: var(--space-4);
}

.confirm-container {
	background: #ffffff;
	border: 1px solid rgba(66, 97, 237, 0.15);
	border-radius: var(--radius-xl);
	width: 100%;
	max-width: 400px;
	box-shadow: 0 20px 50px rgba(66, 97, 237, 0.12);
	overflow: hidden;
	animation: popIn 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.confirm-body {
	padding: 32px 24px 24px;
	display: flex;
	flex-direction: column;
	align-items: center;
	text-align: center;
}

.warning-icon {
	width: 64px;
	height: 64px;
	background-color: rgba(225, 29, 72, 0.1);
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	margin-bottom: 20px;
	padding: 14px;
	border: 8px solid rgba(225, 29, 72, 0.05);
}

.warning-icon img {
	width: 100%;
	height: 100%;
	filter: invert(24%) sepia(87%) saturate(2258%) hue-rotate(330deg) brightness(95%) contrast(92%); /* Rose filter */
}

.confirm-title {
	font-family: var(--font-title);
	font-size: var(--fs-xl);
	font-weight: 700;
	color: var(--text-main);
	margin: 0 0 12px 0;
}

.confirm-message {
	font-size: var(--fs-sm);
	color: var(--text-muted);
	line-height: 1.5;
	margin: 0;
}

.confirm-footer {
	padding: 16px 24px;
	background: #f8faf9;
	display: flex;
	justify-content: center;
	gap: 12px;
	border-top: 1px solid var(--border-color);
}

.btn {
	padding: 10px 16px;
	border-radius: 8px;
	font-weight: 600;
	font-size: var(--fs-sm);
	cursor: pointer;
	transition: background-color 0.2s ease, border-color 0.2s ease, color 0.2s ease, opacity 0.2s ease, transform 0.2s ease, box-shadow 0.2s ease;
	display: flex;
	align-items: center;
	justify-content: center;
	border: none;
}

.btn-cancel {
	background: rgba(255, 255, 255, 0.5);
	border: 1px solid var(--border-color);
	color: var(--text-muted);
	flex: 1;
	box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
}

.btn-cancel:hover {
	background: rgba(255, 255, 255, 0.8);
	border-color: var(--border-hover);
	color: var(--text-main);
}

.btn-delete {
	background: var(--danger-color);
	border: 1px solid var(--danger-color);
	color: white;
	flex: 1;
	box-shadow: 0 4px 14px rgba(225, 29, 72, 0.2);
}

.btn-delete:hover:not(:disabled) {
	background: var(--danger-hover);
	border-color: var(--danger-hover);
	box-shadow: 0 6px 20px rgba(225, 29, 72, 0.3);
}

.btn-delete:disabled {
	opacity: 0.7;
	cursor: not-allowed;
}

.spinner {
	width: 16px;
	height: 16px;
	border: 2px solid rgba(255, 255, 255, 0.3);
	border-radius: 50%;
	border-top-color: #fff;
	animation: spin 0.8s linear infinite;
	margin-right: 8px;
}

@keyframes spin {
	to { transform: rotate(360deg); }
}

@keyframes popIn {
	from {
		transform: scale(0.95) translateY(10px);
		opacity: 0;
	}
	to {
		transform: scale(1) translateY(0);
		opacity: 1;
	}
}

/* Animations overlay */
.modal-fade-enter-active,
.modal-fade-leave-active {
	transition: opacity 0.3s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
	opacity: 0;
}
</style>
