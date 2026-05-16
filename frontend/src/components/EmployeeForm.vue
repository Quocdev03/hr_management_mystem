<template>
	<form @submit.prevent="handleSubmit" class="employee-form">
		<div class="form-grid">
			<!-- Họ và tên -->
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

			<!-- Liên hệ -->
			<div class="form-group" :class="{ 'form-group--disabled': isEdit }">
				<label class="form-label"
					>Email <span class="required" v-if="!isEdit">*</span></label
				>
				<input
					v-model="formData.email"
					type="email"
					class="form-control"
					placeholder="example@company.com"
					:required="!isEdit"
					:disabled="isEdit"
				/>
				<small v-if="isEdit" class="form-hint"
					>Không thể thay đổi email</small
				>
			</div>
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

			<!-- Công việc -->
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
			<div class="form-group">
				<label class="form-label">Chức vụ</label>
				<input
					v-model="formData.position"
					type="text"
					class="form-control"
					placeholder="Ví dụ: Backend Developer"
				/>
			</div>

			<!-- Lương và Ngày vào làm -->
			<div class="form-group">
				<label class="form-label">Mức lương (VNĐ)</label>
				<input
					v-model.number="formData.salary"
					type="number"
					class="form-control"
					placeholder="0"
				/>
			</div>
			<div class="form-group" :class="{ 'form-group--disabled': isEdit }">
				<label class="form-label">Ngày vào làm</label>
				<input
					v-model="formData.join_date"
					type="date"
					class="form-control"
					:disabled="isEdit"
				/>
			</div>
		</div>

		<div class="form-actions">
			<button
				type="button"
				class="btn btn--secondary"
				@click="$emit('cancel')"
			>
				Hủy bỏ
			</button>
			<button type="submit" class="btn btn--primary" :disabled="loading">
				<span v-if="loading" class="spinner"></span>
				{{ isEdit ? "Lưu thay đổi" : "Thêm nhân viên" }}
			</button>
		</div>
	</form>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useDepartmentStore } from "@/store/department";
import { useToast } from "vue-toastification";

const toast = useToast();
const props = defineProps({
	initialData: {
		type: Object,
		default: () => ({}),
	},
	isEdit: Boolean,
	loading: Boolean,
});

const emit = defineEmits(["submit", "cancel"]);

const deptStore = useDepartmentStore();
const departments = ref([]);

const formData = ref({
	first_name: "",
	last_name: "",
	email: "",
	phone: "",
	department_id: "",
	position: "",
	salary: null,
	join_date: "",
	...props.initialData,
});

async function fetchDepts() {
	const res = await deptStore.fetchDepartments();

	if (res.success) {
		departments.value = res.data.items || [];
	}
}

function handleSubmit() {
	// thêm mới
	if (!props.isEdit) {
		emit("submit", formData.value);
		return;
	}

	// cập nhật
	const payload = {};

	if (formData.value.first_name !== props.initialData.first_name) {
		payload.first_name = formData.value.first_name;
	}

	if (formData.value.last_name !== props.initialData.last_name) {
		payload.last_name = formData.value.last_name;
	}

	if (formData.value.phone !== props.initialData.phone) {
		payload.phone = formData.value.phone;
	}

	if (formData.value.department_id !== props.initialData.department_id) {
		payload.department_id = formData.value.department_id;
	}

	if (formData.value.position !== props.initialData.position) {
		payload.position = formData.value.position;
	}

	if (formData.value.salary !== props.initialData.salary) {
		payload.salary = formData.value.salary;
	}
	// không có thay đổi
	if (Object.keys(payload).length === 0) {
		toast.info("Không có dữ liệu thay đổi");
		return;
	}
	emit("submit", payload);
}

onMounted(fetchDepts);
</script>

<style scoped>
.employee-form {
	display: flex;
	flex-direction: column;
	gap: var(--space-4);
}

.form-grid {
	display: grid;
	grid-template-columns: repeat(2, 1fr);
	gap: var(--space-3);
}

.form-actions {
	display: flex;
	justify-content: flex-end;
	gap: var(--space-2);
	margin-top: var(--space-2);
	padding-top: var(--space-3);
	border-top: 1px solid var(--border-color);
}

/* Đồng bộ font chữ với main.css */
.form-label {
	font-size: var(--fs-sm);
	font-weight: var(--fw-semibold);
}

.form-control {
	font-size: var(--fs-base);
}

.form-group--disabled {
	opacity: 0.8;
}

.form-group--disabled .form-control {
	background-color: var(--bg-light);
	cursor: not-allowed;
	border-color: var(--border-color);
}

.form-hint {
	font-size: 11px;
	color: var(--text-light);
	margin-top: -2px;
}

@media (max-width: 640px) {
	.form-grid {
		grid-template-columns: 1fr;
	}
}

.spinner {
	width: 16px;
	height: 16px;
	border: 2px solid rgba(255, 255, 255, 0.3);
	border-radius: 50%;
	border-top-color: #fff;
	animation: spin 0.8s linear infinite;
	display: inline-block;
	margin-right: 8px;
}

@keyframes spin {
	to {
		transform: rotate(360deg);
	}
}
</style>
