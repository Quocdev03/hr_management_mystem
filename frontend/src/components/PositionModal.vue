<script setup>
import { ref, watch } from "vue";
import ModalDialog from "./ModalDialog.vue";

const props = defineProps({
	visible: { type: Boolean, required: true },
	isEditMode: { type: Boolean, required: true },
	editingPosition: { type: Object, default: null },
	loading: { type: Boolean, default: false }
});

const emit = defineEmits(["close", "submit"]);

const formData = ref(buildInitialFormData());

function buildInitialFormData(data = null) {
	const d = data ?? {};
	return {
		name: d.name ?? "",
		description: d.description ?? "",
	};
}

watch(
	() => props.visible,
	(isVisible) => {
		if (isVisible) {
			formData.value = buildInitialFormData(props.editingPosition);
		}
	}
);

watch(
	() => props.editingPosition,
	(newVal) => {
		if (props.visible) {
			formData.value = buildInitialFormData(newVal);
		}
	}
);

function handleSubmit() {
	emit("submit", { ...formData.value });
}

function handleClose() {
	emit("close");
}
</script>

<template>
	<ModalDialog
		:visible="visible"
		:title="isEditMode ? 'Chỉnh sửa chức vụ' : 'Thêm chức vụ mới'"
		:subtitle="
			isEditMode
				? 'Cập nhật thông tin chức vụ hiện tại'
				: 'Nhập thông tin chức vụ mới'
		"
		size="lg"
		@close="handleClose"
	>
		<form @submit.prevent="handleSubmit" class="position-form">
			<div class="form-grid">
				<div class="form-group form-group--full">
					<label class="form-label"
						>Tên chức vụ <span class="required">*</span></label
					>
					<input
						v-model="formData.name"
						type="text"
						class="form-control"
						placeholder="Nhập tên chức vụ (ví dụ: Senior Frontend Developer)..."
						required
					/>
				</div>
				<div class="form-group form-group--full">
					<label class="form-label">Mô tả</label>
					<textarea
						v-model="formData.description"
						class="form-control"
						placeholder="Mô tả chi tiết hoặc các nhiệm vụ chính của chức vụ..."
						rows="4"
					></textarea>
				</div>
			</div>

			<div class="form-actions">
				<button
					type="button"
					class="btn btn--secondary"
					@click="handleClose"
				>
					Hủy bỏ
				</button>
				<button
					type="submit"
					class="btn btn--primary"
					:disabled="loading"
				>
					<span v-if="loading" class="spinner"></span>
					{{ isEditMode ? "Lưu thay đổi" : "Thêm chức vụ" }}
				</button>
			</div>
		</form>
	</ModalDialog>
</template>

<style scoped>
.position-form {
	display: flex;
	flex-direction: column;
	gap: var(--space-4);
}
</style>
