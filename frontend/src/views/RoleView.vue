<script setup>
import { ShieldCheck, Plus, Pencil, Trash2, ShieldAlert } from "@lucide/vue";
import { useRoleStore } from "@/store/role";
import { storeToRefs } from "pinia";
import { ref, onMounted, computed } from "vue";
import { useToast } from "vue-toastification";
import RoleModal from "@/components/RoleModal.vue";
import ConfirmationDialog from "@/components/ConfirmationDialog.vue";
import RoleDetailModal from "@/components/RoleDetailModal.vue";
import Skeleton from "@/components/Skeleton.vue";
import { useModalState } from "@/helpers/useModalState";

const roleStore = useRoleStore();
const { roles, loading } = storeToRefs(roleStore);
const toast = useToast();

const { isModalVisible, isEditMode, openAddModal, openEditModal, closeModal } =
	useModalState();
const modalMode = computed(() => (isEditMode.value ? "edit" : "create"));
const selectedRole = ref(null);

const isDeleteVisible = ref(false);
const roleToDelete = ref(null);

const isDetailVisible = ref(false);
const roleToView = ref(null);

onMounted(async () => {
	await roleStore.fetchAvailablePermissions();
	await fetchRoles();
});

const fetchRoles = async () => {
	try {
		await roleStore.fetchRoles();
	} catch (error) {
		toast.error("Không thể tải danh sách vai trò");
	}
};

const handleViewDetail = (role) => {
	roleToView.value = role;
	isDetailVisible.value = true;
};

const openCreateModal = () => {
	openAddModal();
	selectedRole.value = null;
};

const openEditModalForRole = (role) => {
	openEditModal();
	selectedRole.value = role;
};

const handleModalSubmit = async (payload) => {
	try {
		if (!isEditMode.value) {
			await roleStore.createRole(payload);
			toast.success("Tạo vai trò thành công");
		} else {
			await roleStore.updateRole(selectedRole.value.id, payload);
			toast.success("Cập nhật vai trò thành công");
		}
		isModalVisible.value = false;
	} catch (error) {
		toast.error(error.response?.data?.message || "Có lỗi xảy ra");
	}
};

const openDeleteConfirm = (role) => {
	roleToDelete.value = role;
	isDeleteVisible.value = true;
};

const confirmDelete = async () => {
	if (!roleToDelete.value) return;
	try {
		await roleStore.deleteRole(roleToDelete.value.id);
		toast.success("Xóa vai trò thành công");
		isDeleteVisible.value = false;
		roleToDelete.value = null;
	} catch (error) {
		toast.error(error.response?.data?.message || "Không thể xóa vai trò");
	}
};

const isSystemRole = (name) => {
	const l = name.toLowerCase();
	return l === "admin" || l === "hr" || l === "employee";
};
</script>

<template>
	<div class="role-view">
		<!-- Header -->
		<header class="page-header">
			<div class="header-content">
				<h1 class="page-title">Quản lý Vai trò</h1>
				<p class="page-subtitle">
					Thiết lập quyền truy cập và chức năng cho từng nhóm người
					dùng
				</p>
			</div>

			<button @click="openCreateModal" class="btn btn-primary">
				<Plus class="btn__icon" />
				Thêm vai trò mới
			</button>
		</header>

		<!-- Bento Grid của Vai trò -->
		<main class="bento-container">
			<!-- Loading Skeleton Grid -->
			<div v-if="loading" class="bento-grid">
				<div
					v-for="i in 4"
					:key="'skeleton-role-' + i"
					class="bento-card"
				>
					<div class="bento-header">
						<Skeleton type="text" width="50%" height="22px" />
						<Skeleton type="badge" width="40px" height="24px" />
					</div>
					<div class="bento-body">
						<Skeleton
							type="text"
							width="100%"
							height="14px"
							style="margin-bottom: 8px"
						/>
						<Skeleton type="text" width="80%" height="14px" />
					</div>
					<div class="role-bento-permissions">
						<Skeleton type="badge" width="80px" height="20px" />
						<Skeleton type="badge" width="60px" height="20px" />
						<Skeleton type="badge" width="70px" height="20px" />
					</div>
					<div class="bento-actions">
						<Skeleton type="btn" />
						<Skeleton type="btn" />
					</div>
				</div>
			</div>

			<!-- Actual Role Cards -->
			<div v-else class="bento-grid">
				<div
					v-for="role in roles"
					:key="role.id"
					class="bento-card"
					@click="handleViewDetail(role)"
					style="cursor: pointer"
				>
					<!-- Accent bar -->
					<div
						class="bento-accent-bar"
						:class="{
							'bento-accent-bar--system': isSystemRole(role.name),
						}"
					></div>

					<div class="bento-header">
						<h3 class="bento-name">{{ role.name }}</h3>
						<span
							v-if="isSystemRole(role.name)"
							class="role-badge role-badge--system"
						>
							Hệ thống
						</span>
					</div>

					<div class="bento-body">
						<p class="bento-desc">
							{{
								role.description ||
								"Chưa có mô tả cho vai trò này."
							}}
						</p>
					</div>

					<div class="role-bento-permissions">
						<span class="permission-label">Quyền hạn:</span>
						<div class="permission-tags">
							<span
								v-if="role.name.toLowerCase() === 'admin'"
								class="perm-tag perm-tag--admin"
							>
								Toàn quyền
							</span>
							<template v-else>
								<span
									v-for="perm in (
										role.permissions || []
									).slice(0, 5)"
									:key="perm"
									class="perm-tag"
								>
									{{ perm }}
								</span>
								<span
									v-if="(role.permissions?.length || 0) > 5"
									class="perm-tag perm-tag--more"
									:title="
										(role.permissions || [])
											.slice(5)
											.join(', ')
									"
								>
									+{{ role.permissions.length - 5 }}
								</span>
								<span
									v-if="
										!role.permissions ||
										role.permissions.length === 0
									"
									class="text-muted"
									style="font-style: italic; font-size: 12px"
								>
									Chưa được cấp quyền
								</span>
							</template>
						</div>
					</div>

					<div class="bento-actions">
						<button
							class="btn-icon btn-icon--edit"
							title="Chỉnh sửa"
							@click.stop="openEditModalForRole(role)"
						>
							<Pencil />
						</button>
						<button
							v-if="!isSystemRole(role.name)"
							class="btn-icon btn-icon--delete"
							title="Xóa"
							@click.stop="openDeleteConfirm(role)"
						>
							<Trash2 />
						</button>
					</div>
				</div>

				<!-- Empty State -->
				<div v-if="roles.length === 0" class="empty-state-container">
					<div class="empty-state">
						<ShieldAlert class="empty-state__icon-svg" />
						<p class="empty-state__text">
							Không tìm thấy vai trò nào.
						</p>
					</div>
				</div>
			</div>
		</main>

		<!-- Modals -->
		<RoleModal
			v-model="isModalVisible"
			:mode="modalMode"
			:role-data="selectedRole"
			@submit="handleModalSubmit"
		/>

		<ConfirmationDialog
			:visible="isDeleteVisible"
			title="Xóa vai trò"
			:message="`Bạn có chắc chắn muốn xóa vai trò '${roleToDelete?.name}'? Hành động này không thể hoàn tác.`"
			confirm-text="Xóa"
			confirm-color="bg-rose-600 hover:bg-rose-700"
			@confirm="confirmDelete"
			@cancel="isDeleteVisible = false"
		/>

		<RoleDetailModal
			:visible="isDetailVisible"
			:role="roleToView"
			@close="isDetailVisible = false"
		/>
	</div>
</template>

<style scoped>
.title-icon {
	width: 28px;
	height: 28px;
	color: var(--primary-color);
}

.bento-accent-bar--system {
	background: linear-gradient(90deg, var(--danger-color), #ff8a8a);
}

.role-badge {
	font-size: 11px;
	padding: 2px 8px;
	border-radius: var(--radius-sm);
	font-weight: var(--fw-medium);
}

.role-badge--system {
	background: rgba(225, 29, 72, 0.08);
	color: var(--danger-color);
}

.role-bento-permissions {
	margin-bottom: var(--space-4);
	background: var(--bg-lighter);
	padding: var(--space-3);
	border-radius: var(--radius-md);
	border: 1px dashed var(--border-color);
}

.permission-label {
	display: block;
	font-size: 12px;
	font-weight: var(--fw-medium);
	color: var(--text-light);
	margin-bottom: 8px;
}

.permission-tags {
	display: flex;
	flex-wrap: wrap;
	gap: 6px;
}

.perm-tag {
	padding: 3px 8px;
	background: #ffffff;
	border: 1px solid var(--border-color);
	border-radius: var(--radius-sm);
	font-size: 11px;
	font-weight: var(--fw-medium);
	color: var(--text-main);
	box-shadow: 0 1px 2px rgba(0, 0, 0, 0.02);
}

.perm-tag--admin {
	background: rgba(225, 29, 72, 0.05);
	border-color: rgba(225, 29, 72, 0.2);
	color: var(--danger-color);
	font-weight: var(--fw-semibold);
}

.perm-tag--more {
	background: var(--bg-light);
	color: var(--text-muted);
	cursor: help;
}

.empty-state-container {
	grid-column: 1 / -1;
	display: flex;
	justify-content: center;
	padding: var(--space-6) 0;
}

.empty-state {
	display: flex;
	flex-direction: column;
	align-items: center;
	text-align: center;
}

.empty-state__icon-svg {
	width: 64px;
	height: 64px;
	color: var(--border-hover);
	margin-bottom: var(--space-3);
}

.empty-state__text {
	font-size: var(--fs-base);
	color: var(--text-muted);
	margin: 0;
}

@media (max-width: 768px) {
	.role-bento-grid {
		grid-template-columns: 1fr;
	}
}
</style>
