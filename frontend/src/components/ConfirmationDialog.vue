<template>
	<Teleport to="body">
		<Transition name="modal-fade">
			<div
				v-if="visible"
				class="confirm-overlay"
				@click.self="$emit('cancel')"
			>
				<div class="confirm-container animate__animated animate__zoomIn animate__fast">
					<div class="confirm-body">
						<div class="confirm-icon" :class="[`confirm-icon--${variant}`]">
							<TriangleAlert v-if="variant === 'danger'" />
							<HelpCircle v-else />
						</div>
						<h3 class="confirm-title">{{ title }}</h3>
						<p class="confirm-message">{{ message }}</p>
					</div>

					<div class="confirm-footer">
						<button class="btn btn-secondary" @click="$emit('cancel')">
							{{ cancelText }}
						</button>
						<button
							class="btn"
							:class="[variant === 'danger' ? 'btn-danger' : 'btn-primary']"
							@click="$emit('confirm')"
							:disabled="loading"
						>
							<span v-if="loading" class="btn-spinner"></span>
							<span>{{ confirmText }}</span>
						</button>
					</div>
				</div>
			</div>
		</Transition>
	</Teleport>
</template>

<script setup>
import { TriangleAlert, HelpCircle } from "@lucide/vue";

const props = defineProps({
	visible: Boolean,
	title: {
		type: String,
		default: "Xác nhận",
	},
	message: {
		type: String,
		default: "Bạn có chắc chắn muốn thực hiện hành động này?",
	},
	confirmText: {
		type: String,
		default: "Xác nhận",
	},
	cancelText: {
		type: String,
		default: "Hủy bỏ",
	},
	loading: Boolean,
	variant: {
		type: String,
		default: "danger", // danger | primary
	},
});

const emit = defineEmits(["confirm", "cancel"]);
</script>

<style scoped>
.confirm-overlay {
	position: fixed;
	inset: 0;
	background: rgba(15, 23, 42, 0.3);
	backdrop-filter: blur(8px);
	-webkit-backdrop-filter: blur(8px);
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 1000;
	padding: var(--space-4);
}

.confirm-container {
	background: var(--bg-card);
	border: 1px solid var(--border-color);
	border-radius: var(--radius-xl);
	width: 100%;
	max-width: 400px;
	box-shadow: var(--shadow-xl);
	overflow: hidden;
}

.confirm-body {
	padding: 24px 20px 20px;
	display: flex;
	flex-direction: column;
	align-items: center;
	text-align: center;
}

.confirm-icon {
	width: 48px;
	height: 48px;
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	margin-bottom: 14px;
	transition: transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.confirm-container:hover .confirm-icon {
	transform: scale(1.08);
}

.confirm-icon--danger {
	background: var(--danger-light);
	border: 1px solid rgba(225, 29, 72, 0.15);
	box-shadow: 0 0 0 6px rgba(225, 29, 72, 0.04);
}

.confirm-icon--danger svg {
	width: 24px;
	height: 24px;
	color: var(--danger-color);
}

.confirm-icon--primary {
	background: var(--primary-light);
	border: 1px solid rgba(66, 97, 237, 0.15);
	box-shadow: 0 0 0 6px rgba(66, 97, 237, 0.04);
}

.confirm-icon--primary svg {
	width: 24px;
	height: 24px;
	color: var(--primary-color);
}

.confirm-title {
	font-family: var(--font-title);
	font-size: var(--fs-lg);
	font-weight: var(--fw-bold);
	color: var(--text-main);
	margin: 0 0 10px 0;
}

.confirm-message {
	font-size: var(--fs-sm);
	color: var(--text-muted);
	margin: 0;
	line-height: var(--lh-normal);
}

.confirm-footer {
	background: var(--bg-card);
	padding: 12px 20px;
	display: flex;
	justify-content: space-between;
	gap: 10px;
	border-top: 1px solid var(--border-color);
	border-radius: 0 0 var(--radius-xl) var(--radius-xl);
}

.confirm-footer .btn {
	height: 36px;
	padding: 0 16px;
	font-size: var(--fs-sm);
	border-radius: 8px;
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
