<script setup>
import { ref, watch, computed } from "vue";
import ModalDialog from "./ModalDialog.vue";
import { storeToRefs } from "pinia";
import { useRoleStore } from "@/store/role";

const props = defineProps({
	modelValue: Boolean,
	mode: {
		type: String,
		default: "create", // "create" | "edit"
	},
	roleData: {
		type: Object,
		default: () => ({}),
	},
});

const emit = defineEmits(["update:modelValue", "submit"]);

const roleStore = useRoleStore();
const { availablePermissions } = storeToRefs(roleStore);

// Form state
const formData = ref({
	name: "",
	description: "",
	permissions: [],
});

const isSubmitting = ref(false);
const errorMsg = ref("");

// Nhóm permissions theo module (employee, user, department...)
const groupedPermissions = computed(() => {
	const groups = {};
	availablePermissions.value.forEach((p) => {
		const moduleName = p.code.split(".")[0];
		if (!groups[moduleName]) {
			groups[moduleName] = [];
		}
		groups[moduleName].push(p);
	});
	return groups;
});

const translateModule = (moduleName) => {
	const map = {
		employee: "Nhân viên",
		user: "Tài khoản",
		department: "Phòng ban",
	};
	return map[moduleName] || moduleName;
};

watch(
	() => props.modelValue,
	(newVal) => {
		if (newVal) {
			errorMsg.value = "";
			if (props.mode === "edit" && props.roleData) {
				formData.value = {
					name: props.roleData.name || "",
					description: props.roleData.description || "",
					permissions: props.roleData.permissions
						? [...props.roleData.permissions]
						: [],
				};
			} else {
				formData.value = {
					name: "",
					description: "",
					permissions: [],
				};
			}
		}
	},
);

const isSystemRole = computed(() => {
	const name = formData.value.name.toLowerCase();
	return (
		props.mode === "edit" &&
		(name === "admin" || name === "hr" || name === "employee")
	);
});

const isAdminRole = computed(() => {
	return (
		props.mode === "edit" && formData.value.name.toLowerCase() === "admin"
	);
});

const handleSubmit = async () => {
	if (!formData.value.name.trim()) {
		errorMsg.value = "Tên vai trò không được để trống";
		return;
	}
	errorMsg.value = "";
	isSubmitting.value = true;
	try {
		const payload = {
			description: formData.value.description,
			permissions: formData.value.permissions,
		};
		// Chỉ gửi name nếu không phải role hệ thống (hoặc nếu là tạo mới)
		if (!isSystemRole.value) {
			payload.name = formData.value.name;
		}

		emit("submit", payload);
	} catch (error) {
		errorMsg.value = "Có lỗi xảy ra khi lưu vai trò";
	} finally {
		isSubmitting.value = false;
	}
};
</script>

<template>
	<ModalDialog
		:visible="modelValue"
		@close="$emit('update:modelValue', false)"
		:title="mode === 'create' ? 'Thêm vai trò mới' : 'Chỉnh sửa vai trò'"
		size="md"
	>
		<form @submit.prevent="handleSubmit" class="form-wrapper">
			<!-- Tên & Mô tả -->
			<div class="form-group form-group--full">
				<label class="form-label">
					Tên vai trò <span class="required">*</span>
				</label>
				<input
					v-model="formData.name"
					type="text"
					:disabled="isSystemRole"
					class="form-control"
					placeholder="Ví dụ: Quản lý chi nhánh"
				/>
				<p v-if="isSystemRole" class="helper-text helper-text--warning">
					Không thể đổi tên vai trò hệ thống.
				</p>
			</div>

			<div class="form-group form-group--full">
				<label class="form-label"> Mô tả </label>
				<textarea
					v-model="formData.description"
					rows="2"
					class="form-control"
					placeholder="Mô tả chức năng của vai trò này"
				></textarea>
			</div>

			<!-- Phân quyền -->
			<div class="permission-section form-group--full">
				<label class="section-title"> Gán Quyền hạn </label>

				<div v-if="isAdminRole" class="admin-notice">
					Vai trò <strong>Admin</strong> mặc định có toàn quyền. Bạn
					không cần phải tick thủ công.
				</div>

				<div class="permission-list custom-scrollbar">
					<div
						v-for="(perms, moduleName) in groupedPermissions"
						:key="moduleName"
						class="permission-group"
					>
						<h4 class="permission-group-title">
							{{ translateModule(moduleName) }}
						</h4>
						<div class="permission-grid">
							<label
								v-for="perm in perms"
								:key="perm.code"
								class="permission-item"
							>
								<input
									type="checkbox"
									:value="perm.code"
									v-model="formData.permissions"
									class="permission-checkbox"
								/>
								<div class="permission-info">
									<span class="permission-code">{{
										perm.code
									}}</span>
									<span class="permission-desc">{{
										perm.description
									}}</span>
								</div>
							</label>
						</div>
					</div>
				</div>
			</div>

			<!-- Error message -->
			<div v-if="errorMsg" class="error-notice form-group--full">
				{{ errorMsg }}
			</div>
		</form>
		<template #footer>
			<button
				type="button"
				@click="$emit('update:modelValue', false)"
				class="btn btn-secondary"
			>
				Hủy
			</button>
			<button
				type="button"
				@click="handleSubmit"
				:disabled="isSubmitting"
				class="btn btn-primary"
			>
				<span v-if="isSubmitting" class="btn-spinner"></span>
				{{ mode === "create" ? "Tạo mới" : "Lưu thay đổi" }}
			</button>
		</template>
	</ModalDialog>
</template>

<style scoped>
.permission-section {
	margin-top: var(--space-2);
	padding-top: var(--space-3);
	border-top: 1px solid var(--border-color);
}

.admin-notice {
	padding: var(--space-2);
	background: rgba(66, 97, 237, 0.08);
	color: var(--primary-color);
	border-radius: var(--radius-sm);
	font-size: var(--fs-sm);
	margin-bottom: var(--space-3);
}

.permission-list {
	display: flex;
	flex-direction: column;
	gap: var(--space-3);
	max-height: 40vh;
	overflow-y: auto;
	padding-right: var(--space-2);
}

.permission-group {
	background: var(--bg-lighter);
	padding: 12px 16px;
	border-radius: var(--radius-md);
	border: 1px solid var(--border-color);
}

.permission-group-title {
	font-size: var(--fs-sm);
	font-weight: var(--fw-medium);
	color: var(--text-main);
	margin-bottom: var(--space-2);
	padding-bottom: 4px;
	border-bottom: 1px solid var(--border-color);
	text-transform: capitalize;
}

.permission-grid {
	display: grid;
	grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
	gap: var(--space-2);
}

.permission-item {
	display: flex;
	align-items: flex-start;
	gap: 8px;
	cursor: pointer;
	padding: 8px 10px;
	border-radius: var(--radius-sm);
	transition: background-color 0.2s ease;
}

.permission-item:hover {
	background: var(--bg-main);
}

.permission-checkbox {
	margin-top: 2px;
	width: 16px;
	height: 16px;
	accent-color: var(--primary-color);
}

.permission-info {
	display: flex;
	flex-direction: column;
}

.permission-code {
	font-size: 12px;
	font-weight: var(--fw-medium);
	color: var(--text-main);
}

.permission-desc {
	font-size: 11px;
	color: var(--text-muted);
}

.error-notice {
	padding: var(--space-2);
	background: rgba(225, 29, 72, 0.08);
	color: var(--danger-color);
	border-radius: var(--radius-sm);
	font-size: var(--fs-sm);
}

.custom-scrollbar::-webkit-scrollbar {
	width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
	background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
	background-color: var(--border-hover);
	border-radius: 20px;
}
</style>
