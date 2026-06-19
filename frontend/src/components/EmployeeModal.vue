<script setup>
import { ref, watch } from "vue";
import ModalDialog from "./ModalDialog.vue";

const props = defineProps({
	visible: { type: Boolean, required: true },
	isEditMode: { type: Boolean, required: true },
	editingEmployee: { type: Object, default: null },
	departments: { type: Array, required: true },
	usersWithoutEmp: { type: Array, required: true },
	loading: { type: Boolean, default: false }
});

const emit = defineEmits(["close", "submit"]);

// Local state for form data
const formData = ref(buildInitialFormData());

function buildInitialFormData(data = null) {
	const d = data ?? {};
	
	// Helper to format date strings for input[type="date"] (YYYY-MM-DD)
	const formatDateString = (dateStr) => {
		if (!dateStr) return "";
		return dateStr.split("T")[0];
	};

	return {
		user_id: d.user_id ?? null,
		first_name: d.first_name ?? "",
		last_name: d.last_name ?? "",
		phone: d.phone ?? "",
		department_id: d.department_id ?? "",
		position: d.position ?? "",
		salary: d.salary ?? null,
		join_date: formatDateString(d.join_date) ?? "",
		status: d.status ?? "active",
		gender: d.gender ?? "",
		birth_date: formatDateString(d.birth_date) ?? "",
	};
}

// Watch visibility and editingEmployee to sync form data
watch(
	() => props.visible,
	(isVisible) => {
		if (isVisible) {
			formData.value = buildInitialFormData(props.editingEmployee);
		}
	}
);

watch(
	() => props.editingEmployee,
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
		:title="isEditMode ? 'Chỉnh sửa nhân viên' : 'Thêm nhân viên mới'"
		:subtitle="
			isEditMode
				? 'Cập nhật thông tin chi tiết của nhân viên'
				: 'Điền thông tin để tạo nhân viên mới vào hệ thống'
		"
		size="lg"
		@close="handleClose"
	>
		<form @submit.prevent="handleSubmit" class="employee-form">
			<div class="form-grid">
				<!-- Họ & Tên -->
				<div class="form-group">
					<label class="form-label"
						>Họ <span class="required">*</span></label
					>
					<input
						v-model="formData.first_name"
						type="text"
						class="form-control"
						placeholder="Nhập họ..."
						required
					/>
				</div>
				<div class="form-group">
					<label class="form-label"
						>Tên <span class="required">*</span></label
					>
					<input
						v-model="formData.last_name"
						type="text"
						class="form-control"
						placeholder="Nhập tên..."
						required
					/>
				</div>

				<!-- Giới tính -->
				<div class="form-group">
					<label class="form-label">Giới tính</label>
					<select v-model="formData.gender" class="form-control">
						<option value="" disabled>Chọn giới tính</option>
						<option value="male">Nam</option>
						<option value="female">Nữ</option>
					</select>
				</div>

				<!-- Ngày sinh -->
				<div class="form-group">
					<label class="form-label">Ngày sinh</label>
					<input
						v-model="formData.birth_date"
						type="date"
						class="form-control"
					/>
				</div>

				<!-- Số điện thoại -->
				<div class="form-group">
					<label class="form-label"
						>Số điện thoại <span class="required">*</span></label
					>
					<input
						v-model="formData.phone"
						type="text"
						class="form-control"
						placeholder="0xxxxxxxxx"
						required
					/>
				</div>

				<!-- Phòng ban -->
				<div class="form-group">
					<label class="form-label"
						>Phòng ban <span class="required">*</span></label
					>
					<select
						v-model="formData.department_id"
						class="form-control"
						required
					>
						<option value="" disabled>Chọn phòng ban</option>
						<option
							v-for="dept in departments"
							:key="dept.id"
							:value="dept.id"
						>
							{{ dept.name }}
						</option>
					</select>
				</div>

				<!-- Chức vụ -->
				<div class="form-group">
					<label class="form-label">Chức vụ</label>
					<input
						v-model="formData.position"
						type="text"
						class="form-control"
						placeholder="Ví dụ: Backend Developer"
					/>
				</div>

				<!-- Mức lương -->
				<div class="form-group">
					<label class="form-label">Mức lương (VNĐ)</label>
					<input
						v-model.number="formData.salary"
						type="number"
						class="form-control"
						placeholder="0"
					/>
				</div>

				<!-- Ngày vào làm (khoá khi sửa) -->
				<div
					class="form-group"
					:class="{ 'form-group--disabled': isEditMode }"
				>
					<label class="form-label">Ngày vào làm</label>
					<input
						v-model="formData.join_date"
						type="date"
						class="form-control"
						:disabled="isEditMode"
					/>
				</div>

				<!-- Trạng thái -->
				<div class="form-group">
					<label class="form-label">Trạng thái</label>
					<select v-model="formData.status" class="form-control">
						<option value="active">Đang làm việc</option>
						<option value="inactive">Đã nghỉ việc</option>
					</select>
				</div>
			</div>

			<!-- Dropdown liên kết tài khoản -->
			<div class="form-group" style="margin-top: var(--space-3)">
				<span class="form-label">Liên kết người dùng</span>
				<select v-model="formData.user_id" class="form-control">
					<option :value="null">Không liên kết</option>
					<option
						v-if="editingEmployee?.user"
						:value="editingEmployee.user.id"
					>
						{{ editingEmployee.user.email }} (hiện tại)
					</option>
					<option
						v-for="u in usersWithoutEmp"
						:key="u.id"
						:value="u.id"
					>
						{{ u.email }}
					</option>
				</select>
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
					{{ isEditMode ? "Lưu thay đổi" : "Thêm nhân viên" }}
				</button>
			</div>
		</form>
	</ModalDialog>
</template>

<style scoped>
.employee-form {
	display: flex;
	flex-direction: column;
	gap: var(--space-4);
}

.form-group--disabled {
	opacity: 0.7;
}
</style>
