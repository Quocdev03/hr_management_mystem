<script setup>
import { ref, watch } from "vue";
import ModalDialog from "./ModalDialog.vue";
import { usePositionStore } from "@/store/position";

const props = defineProps({
	visible: { type: Boolean, required: true },
	isEditMode: { type: Boolean, required: true },
	editingEmployee: { type: Object, default: null },
	departments: { type: Array, required: true },
	usersWithoutEmp: { type: Array, required: true },
	loading: { type: Boolean, default: false }
});

const emit = defineEmits(["close", "submit"]);

const positionStore = usePositionStore();
const positions = ref([]);
const loadingPositions = ref(false);

const buildInitialFormData = (data = null) => {
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
		position_id: d.position_id ?? "",
		salary: d.salary ?? null,
		join_date: formatDateString(d.join_date) ?? "",
		status: d.status ?? "active",
		gender: d.gender ?? "",
		birth_date: formatDateString(d.birth_date) ?? "",
	};
};

// Local state for form data
const formData = ref(buildInitialFormData());

// Watch visibility and editingEmployee to sync form data
watch(
	() => props.visible,
	async (isVisible) => {
		if (isVisible) {
			formData.value = buildInitialFormData(props.editingEmployee);
			if (positions.value.length === 0) {
				loadingPositions.value = true;
				try {
					const res = await positionStore.fetchPositions();
					if (res.success) {
						positions.value = res.data || [];
					}
				} catch (error) {
					console.error("Lỗi tải chức vụ:", error);
				} finally {
					loadingPositions.value = false;
				}
			}
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

const handleSubmit = () => {
	// Map position_id to number
	const payload = { ...formData.value };
	payload.position_id = payload.position_id ? Number(payload.position_id) : undefined;
	emit("submit", payload);
};

const handleClose = () => {
	emit("close");
};
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
		size="md"
		@close="handleClose"
	>
		<form @submit.prevent="handleSubmit" class="form-wrapper">
			<div class="form-grid">
				<!-- Section 1: Personal Info -->
				<h4 class="form-section-title">Thông tin cá nhân</h4>

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
				<div class="form-group form-group--full-mobile">
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

				<!-- Section 2: Work Info -->
				<h4 class="form-section-title">Công việc & Tài khoản</h4>

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
					<label class="form-label"
						>Chức vụ <span class="required">*</span></label
					>
					<select
						v-model="formData.position_id"
						class="form-control"
						:disabled="loadingPositions"
						required
					>
						<option value="" disabled>Chọn chức vụ</option>
						<option
							v-for="pos in positions"
							:key="pos.id"
							:value="pos.id"
						>
							{{ pos.name }}
						</option>
					</select>
				</div>

				<!-- Mức lương -->
				<div class="form-group">
					<label class="form-label">Mức lương (VNĐ)</label>
					<input
						v-model.number="formData.salary"
						type="number"
						class="form-control"
						placeholder="Nhập mức lương..."
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

				<!-- Ngày vào làm (khoá khi sửa) -->
				<div
					class="form-group form-group--full-mobile"
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

				<!-- Dropdown liên kết tài khoản -->
				<div class="form-group form-group--full">
					<label class="form-label">Liên kết tài khoản người dùng</label>
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
			</div>
		</form>
		<template #footer>
			<button
				type="button"
				class="btn btn-secondary cancel-btn"
				@click="handleClose"
			>
				Hủy bỏ
			</button>
			<button
				type="button"
				class="btn btn-primary submit-btn"
				:disabled="loading"
				@click="handleSubmit"
			>
				<span v-if="loading" class="btn-spinner"></span>
				{{ isEditMode ? "Lưu thay đổi" : "Thêm nhân viên" }}
			</button>
		</template>
	</ModalDialog>
</template>

<style scoped>
.form-section-title {
	grid-column: 1 / -1;
	font-size: 0.725rem;
	font-weight: 700;
	color: #4f46e5;
	text-transform: uppercase;
	letter-spacing: 0.05em;
	border-bottom: 1px dashed rgba(0, 0, 0, 0.08);
	padding-bottom: 6px;
	margin-top: 14px;
	margin-bottom: 2px;
}

.form-section-title:first-of-type {
	margin-top: 0;
}

.form-group--disabled {
	opacity: 0.65;
}

.form-group--disabled .form-control {
	background-color: #f8fafc;
	cursor: not-allowed;
	border-color: rgba(0, 0, 0, 0.05);
}

.cancel-btn {
	font-size: 0.825rem;
	font-weight: 600;
	height: 36px;
	padding: 0 1rem;
	border-radius: 8px;
}

.submit-btn {
	font-size: 0.825rem;
	font-weight: 600;
	height: 36px;
	padding: 0 1rem;
	border-radius: 8px;
	box-shadow: 0 2px 4px rgba(66, 97, 237, 0.15);
}

@media (max-width: 640px) {
	.form-group--full-mobile {
		grid-column: 1 / -1;
	}
	
	.form-section-title {
		margin-top: 10px;
	}
}
</style>
