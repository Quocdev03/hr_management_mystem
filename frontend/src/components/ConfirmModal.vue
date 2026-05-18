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
							<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
								<path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
							</svg>
						</div>
						<h3 class="confirm-title">{{ title }}</h3>
						<p class="confirm-message">{{ message }}</p>
					</div>
					
					<div class="confirm-footer">
						<button class="btn btn-cancel" @click="$emit('cancel')">Hủy bỏ</button>
						<button class="btn btn-delete" @click="$emit('confirm')" :disabled="loading">
							<span v-if="loading" class="spinner"></span>
							<span>Xác nhận xoá</span>
						</button>
					</div>
				</div>
			</div>
		</Transition>
	</Teleport>
</template>

<script setup>
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
	loading: Boolean,
});

const emit = defineEmits(["confirm", "cancel"]);
</script>

<style scoped>
.confirm-overlay {
	position: fixed;
	inset: 0;
	background: rgba(15, 23, 42, 0.65);
	backdrop-filter: blur(4px);
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 1000;
	padding: var(--space-4);
}

.confirm-container {
	background: var(--bg-card, #ffffff);
	border-radius: var(--radius-xl, 16px);
	width: 100%;
	max-width: 400px;
	box-shadow: 0 20px 25px -5px rgb(0 0 0 / 0.1), 0 8px 10px -6px rgb(0 0 0 / 0.1);
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
	color: #ef4444;
	background-color: #fef2f2;
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	margin-bottom: 20px;
	padding: 14px;
	border: 8px solid #fff1f2;
}

.warning-icon svg {
	width: 100%;
	height: 100%;
}

.confirm-title {
	font-size: 20px;
	font-weight: 700;
	color: var(--text-main, #0f172a);
	margin: 0 0 12px 0;
}

.confirm-message {
	font-size: 15px;
	color: var(--text-muted, #64748b);
	line-height: 1.5;
	margin: 0;
}

.confirm-footer {
	padding: 16px 24px;
	background: var(--bg-lighter, #f8fafc);
	display: flex;
	justify-content: center;
	gap: 12px;
	border-top: 1px solid var(--border-color, #e2e8f0);
}

.btn {
	padding: 10px 16px;
	border-radius: 8px;
	font-weight: 600;
	font-size: 14px;
	cursor: pointer;
	transition: all 0.2s ease;
	display: flex;
	align-items: center;
	justify-content: center;
	border: none;
}

.btn-cancel {
	background: white;
	border: 1px solid var(--border-color, #e2e8f0);
	color: var(--text-main, #334155);
	flex: 1;
	box-shadow: 0 1px 2px 0 rgb(0 0 0 / 0.05);
}

.btn-cancel:hover {
	background: #f1f5f9;
	color: #0f172a;
}

.btn-delete {
	background: #ef4444;
	border: 1px solid #ef4444;
	color: white;
	flex: 1;
	box-shadow: 0 1px 2px 0 rgb(0 0 0 / 0.05);
}

.btn-delete:hover:not(:disabled) {
	background: #dc2626;
	border-color: #dc2626;
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
