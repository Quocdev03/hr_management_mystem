<script setup>
import { ref, watch } from "vue";
import ModalDialog from "./ModalDialog.vue";

const props = defineProps({
	visible: { type: Boolean, required: true },
	isEditMode: { type: Boolean, required: true },
	editingDepartment: { type: Object, default: null },
	departmentEmployees: { type: Array, required: true },
	loading: { type: Boolean, default: false }
});

const emit = defineEmits(["close", "submit"]);

const buildInitialFormData = (data = null) => {
	const d = data ?? {};
	return {
		name: d.name ?? "",
		code: d.code ?? "",
		description: d.description ?? "",
		manager_id: d.manager_id ?? d.manager?.id ?? null,
	};
};

const formData = ref(buildInitialFormData());

watch(
	() => props.visible,
	(isVisible) => {
		if (isVisible) {
			formData.value = buildInitialFormData(props.editingDepartment);
		}
	}
);

watch(
	() => props.editingDepartment,
	(newVal) => {
		if (props.visible) {
			formData.value = buildInitialFormData(newVal);
		}
	}
);

const handleSubmit = () => {
	emit("submit", { ...formData.value });
};

const handleClose = () => {
	emit("close");
};
</script>

<template>
	<ModalDialog
		:visible="visible"
		:title="isEditMode ? 'Chỉnh sửa phòng ban' : 'Thêm phòng ban mới'"
		:subtitle="
			isEditMode
				? 'Cập nhật phòng ban hiện tại'
				: 'Nhập thông tin phòng ban mới'
		"
		size="md"
		@close="handleClose"
	>
		<form @submit.prevent="handleSubmit" class="form-wrapper">
			<div class="form-grid">
				<div class="form-group">
					<label class="form-label"
						>Tên phòng ban <span class="required">*</span></label
					>
					<input
						v-model="formData.name"
						type="text"
						class="form-control"
						placeholder="Nhập tên phòng ban..."
						required
					/>
				</div>
				<div class="form-group">
					<label class="form-label"
						>Mã phòng ban <span class="required">*</span></label
					>
					<input
						v-model="formData.code"
						type="text"
						class="form-control"
						placeholder="Nhập mã phòng ban..."
						required
						:disabled="isEditMode"
					/>
				</div>
				<div v-if="isEditMode" class="form-group">
					<label class="form-label">Trưởng phòng</label>
					<select
						v-model="formData.manager_id"
						class="form-control"
					>
						<option :value="null">Chọn trưởng phòng</option>
						<option
							v-for="employee in departmentEmployees"
							:key="employee.id"
							:value="employee.id"
						>
							{{ employee.first_name }} {{ employee.last_name }}
						</option>
					</select>
				</div>
				<div class="form-group form-group--full">
					<label class="form-label">Mô tả</label>
					<textarea
						v-model="formData.description"
						class="form-control"
						placeholder="Mô tả ngắn về phòng ban"
						rows="4"
					></textarea>
				</div>
			</div>

		</form>
		<template #footer>
			<button
				type="button"
				class="btn btn-secondary"
				@click="handleClose"
			>
				Hủy bỏ
			</button>
			<button
				type="button"
				class="btn btn-primary"
				:disabled="loading"
				@click="handleSubmit"
			>
				<span v-if="loading" class="btn-spinner"></span>
				{{ isEditMode ? "Lưu thay đổi" : "Thêm phòng ban" }}
			</button>
		</template>
	</ModalDialog>
</template>
