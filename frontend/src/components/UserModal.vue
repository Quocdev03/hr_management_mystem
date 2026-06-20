<script setup>
import { ref, watch } from "vue";
import ModalDialog from "./ModalDialog.vue";

const props = defineProps({
	visible: { type: Boolean, required: true },
	isEditing: { type: Boolean, required: true },
	editingUser: { type: Object, default: null },
	isRoleDisabled: { type: Boolean, default: false },
	isActiveDisabled: { type: Boolean, default: false },
	loading: { type: Boolean, default: false }
});

const emit = defineEmits(["close", "submit"]);

const formData = ref(buildInitialFormData());

const roles = [
	{ id: 1, label: "Admin (Quản trị viên)" },
	{ id: 2, label: "HR (Nhân sự)" },
	{ id: 3, label: "Employee (Nhân viên)" },
];

function buildInitialFormData(data = null) {
	const d = data ?? {};
	return {
		user_name: d.user_name ?? "",
		email: d.email ?? "",
		password: "",
		password_confirm: "",
		role_id: d.role_id ?? 3,
		is_active: d.is_active ?? true,
	};
}

watch(
	() => props.visible,
	(isVisible) => {
		if (isVisible) {
			formData.value = buildInitialFormData(props.editingUser);
		}
	}
);

watch(
	() => props.editingUser,
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
		:title="isEditing ? 'Chỉnh sửa người dùng' : 'Thêm người dùng mới'"
		:subtitle="
			isEditing
				? 'Cập nhật thông tin người dùng'
				: 'Nhập thông tin người dùng mới'
		"
		size="lg"
		@close="handleClose"
	>
		<form @submit.prevent="handleSubmit" class="form-wrapper">
			<div class="form-grid">
				<div class="form-group">
					<label class="form-label"
						>Tên đăng nhập <span class="required">*</span></label
					>
					<input
						v-model="formData.user_name"
						type="text"
						class="form-control"
						required
						placeholder="Nhập username..."
					/>
				</div>
				<div class="form-group">
					<label class="form-label"
						>Email <span class="required">*</span></label
					>
					<input
						v-model="formData.email"
						type="email"
						class="form-control"
						required
						placeholder="Nhập email..."
					/>
				</div>
				<div class="form-group">
					<label class="form-label"
						>Mật khẩu
						<span v-if="!isEditing" class="required">*</span></label
					>
					<input
						v-model="formData.password"
						type="password"
						class="form-control"
						:required="!isEditing"
						:placeholder="
							isEditing
								? 'Bỏ trống nếu không đổi'
								: 'Tối thiểu 8 ký tự'
						"
					/>
				</div>
				<div class="form-group">
					<label class="form-label"
						>Xác nhận mật khẩu
						<span
							v-if="!isEditing || formData.password"
							class="required"
							>*</span
						></label
					>
					<input
						v-model="formData.password_confirm"
						type="password"
						class="form-control"
						:required="!isEditing || !!formData.password"
						:placeholder="
							isEditing ? 'Nhập lại nếu đổi' : 'Nhập lại mật khẩu'
						"
					/>
				</div>
				<div class="form-group">
					<label class="form-label">Vai trò (Role)</label>
					<select
						v-model="formData.role_id"
						class="form-control"
						:disabled="isRoleDisabled"
					>
						<option
							v-for="role in roles"
							:key="role.id"
							:value="role.id"
						>
							{{ role.label }}
						</option>
					</select>
					<small
						v-if="isRoleDisabled"
						class="required"
						style="margin-top: 4px; display: block"
						>Không thể thay đổi quyền của bản thân/Admin khác</small
					>
				</div>
				<div class="form-group">
					<div class="form-group form-group--full">
						<label class="form-label">Trạng thái</label>
						<select
							v-model="formData.is_active"
							class="form-control"
							:disabled="isActiveDisabled"
						>
							<option :value="true">Hoạt động</option>
							<option :value="false">Ngưng</option>
						</select>
					</div>
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
				{{ isEditing ? "Lưu thay đổi" : "Thêm người dùng" }}
			</button>
		</template>
	</ModalDialog>
</template>


