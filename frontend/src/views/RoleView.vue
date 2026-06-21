<script setup>
import { ShieldCheck, Plus, Pencil, Trash2, ShieldAlert, Lock, Shield, UserCog, Eye } from "@lucide/vue";
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

// Map role name to a color theme
const getRoleTheme = (name) => {
	const l = name.toLowerCase();
	if (l === "admin") return "admin";
	if (l === "hr") return "hr";
	if (l === "employee") return "employee";
	return "custom";
};

// Friendly initial(s) for role avatar
const getRoleInitial = (name) => {
	if (!name) return "?";
	return name.trim().substring(0, 2).toUpperCase();
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
					v-for="i in 3"
					:key="'skeleton-role-' + i"
					class="bento-card"
				>
					<div class="bento-accent-bar"></div>
					<div class="role-card-top">
						<Skeleton type="circle" width="52px" height="52px" />
						<div class="role-card-title-group">
							<Skeleton type="text" width="60%" height="20px" style="margin-bottom: 8px" />
							<Skeleton type="badge" width="60px" height="22px" />
						</div>
					</div>
					<div class="bento-body">
						<Skeleton type="text" width="100%" height="14px" style="margin-bottom: 8px" />
						<Skeleton type="text" width="80%" height="14px" />
					</div>
					<div class="role-perm-section">
						<Skeleton type="text" width="40%" height="12px" style="margin-bottom: 10px" />
						<div style="display: flex; gap: 6px; flex-wrap: wrap;">
							<Skeleton type="badge" width="72px" height="22px" />
							<Skeleton type="badge" width="56px" height="22px" />
							<Skeleton type="badge" width="64px" height="22px" />
						</div>
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
					class="bento-card role-card"
					:class="`role-card--${getRoleTheme(role.name)}`"
					@click="handleViewDetail(role)"
					style="cursor: pointer"
				>
					<!-- Accent bar -->
					<div
						class="bento-accent-bar"
						:class="`bento-accent-bar--${getRoleTheme(role.name)}`"
					></div>

					<!-- Card top: avatar + title -->
					<div class="role-card-top">
						<div
							class="role-avatar"
							:class="`role-avatar--${getRoleTheme(role.name)}`"
						>
							<Shield
								v-if="role.name.toLowerCase() === 'admin'"
								class="role-avatar-icon"
							/>
							<UserCog
								v-else-if="role.name.toLowerCase() === 'hr'"
								class="role-avatar-icon"
							/>
							<ShieldCheck
								v-else-if="
									role.name.toLowerCase() === 'employee'
								"
								class="role-avatar-icon"
							/>
							<span v-else class="role-avatar-initials">{{
								getRoleInitial(role.name)
							}}</span>
						</div>

						<div class="role-card-title-group">
							<h3 class="bento-name">{{ role.name }}</h3>
							<span
								v-if="isSystemRole(role.name)"
								class="role-badge role-badge--system"
							>
								<Lock class="role-badge-icon" />
								Hệ thống
							</span>
							<span v-else class="role-badge role-badge--custom">
								Tùy chỉnh
							</span>
						</div>
					</div>

					<!-- Description -->
					<div class="bento-body">
						<p class="bento-desc">
							{{
								role.description ||
								"Chưa có mô tả cho vai trò này."
							}}
						</p>
					</div>

					<!-- Permissions Section -->
					<div class="role-perm-section">
						<div class="role-perm-header">
							<span class="perm-label-text">Quyền hạn</span>
							<span
								class="perm-count-badge"
								v-if="role.name.toLowerCase() !== 'admin'"
							>
								{{
									role.permissions?.length || 0
								}}
								quyền
							</span>
							<span v-else class="perm-count-badge perm-count-badge--full">
								Toàn quyền
							</span>
						</div>

						<div class="permission-tags">
							<!-- Admin special case -->
							<span
								v-if="role.name.toLowerCase() === 'admin'"
								class="perm-tag perm-tag--admin"
							>
								✦ Tất cả quyền hạn
							</span>

							<!-- Normal permissions -->
							<template v-else>
								<span
									v-for="perm in (role.permissions || []).slice(0, 4)"
									:key="perm"
									class="perm-tag"
								>
									{{ perm }}
								</span>
								<span
									v-if="(role.permissions?.length || 0) > 4"
									class="perm-tag perm-tag--more"
									:title="
										(role.permissions || [])
											.slice(4)
											.join(', ')
									"
								>
									+{{ role.permissions.length - 4 }} khác
								</span>
								<span
									v-if="
										!role.permissions ||
										role.permissions.length === 0
									"
									class="perm-tag perm-tag--empty"
								>
									Chưa được cấp quyền
								</span>
							</template>
						</div>
					</div>

					<!-- Actions -->
					<div class="bento-actions">
						<button
							class="btn-icon btn-icon--detail"
							title="Xem chi tiết"
							@click.stop="handleViewDetail(role)"
						>
							<Eye />
						</button>
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
						<div class="empty-state-icon-wrap">
							<ShieldAlert class="empty-state__icon-svg" />
						</div>
						<h3 class="empty-state__title">Chưa có vai trò nào</h3>
						<p class="empty-state__text">
							Tạo vai trò đầu tiên để bắt đầu phân quyền cho
							người dùng.
						</p>
						<button
							class="btn btn-primary"
							@click="openCreateModal"
							style="margin-top: 1.5rem"
						>
							<Plus class="btn__icon" />
							Thêm vai trò mới
						</button>
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
/* ── View wrapper ─────────────────────────────────────────── */
.role-view {
	padding-bottom: var(--space-4);
}


/* ── Bento accent bar variants ───────────────────────────── */
.bento-accent-bar--admin {
	background: linear-gradient(90deg, var(--danger-color), #f43f5e, #fb923c);
}
.bento-accent-bar--hr {
	background: linear-gradient(90deg, var(--purple-color), #a78bfa);
}
.bento-accent-bar--employee {
	background: linear-gradient(90deg, var(--primary-color), #38bdf8);
}
.bento-accent-bar--custom {
	background: linear-gradient(90deg, var(--primary-color), var(--info-color));
}

/* ── Role Card ───────────────────────────────────────────── */
/* Subtle hover glow per theme */
.role-card--admin:hover {
	border-color: rgba(225, 29, 72, 0.18) !important;
	box-shadow:
		0 20px 40px rgba(225, 29, 72, 0.08),
		0 4px 12px rgba(225, 29, 72, 0.04) !important;
}
.role-card--hr:hover {
	border-color: rgba(124, 58, 237, 0.18) !important;
	box-shadow:
		0 20px 40px rgba(124, 58, 237, 0.08),
		0 4px 12px rgba(124, 58, 237, 0.04) !important;
}
.role-card--employee:hover {
	border-color: rgba(66, 97, 237, 0.18) !important;
	box-shadow:
		0 20px 40px rgba(66, 97, 237, 0.08),
		0 4px 12px rgba(66, 97, 237, 0.04) !important;
}

/* ── Role card top: avatar + title ───────────────────────── */
.role-card-top {
	display: flex;
	align-items: center;
	gap: var(--space-2);
	margin-bottom: var(--space-2);
}

.role-avatar {
	width: 40px;
	height: 40px;
	border-radius: var(--radius-md);
	display: flex;
	align-items: center;
	justify-content: center;
	flex-shrink: 0;
	position: relative;
}

.role-avatar--admin {
	background: var(--danger-light);
	color: var(--danger-color);
}
.role-avatar--hr {
	background: var(--purple-light);
	color: var(--purple-color);
}
.role-avatar--employee {
	background: var(--primary-light);
	color: var(--primary-color);
}
.role-avatar--custom {
	background: var(--primary-light);
	color: var(--primary-color);
}

.role-avatar-icon {
	width: 20px;
	height: 20px;
}

.role-avatar-initials {
	font-size: 13px;
	font-weight: var(--fw-bold);
	letter-spacing: 0.03em;
}

.role-card-title-group {
	flex: 1;
	min-width: 0;
}

.bento-name {
	margin: 0 0 6px 0;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

/* ── Role badges ─────────────────────────────────────────── */
.role-badge {
	display: inline-flex;
	align-items: center;
	gap: 4px;
	font-size: 11px;
	padding: 3px 8px;
	border-radius: var(--radius-full);
	font-weight: var(--fw-semibold);
	line-height: 1;
}

.role-badge-icon {
	width: 10px;
	height: 10px;
}

.role-badge--system {
	background: rgba(225, 29, 72, 0.08);
	color: var(--danger-color);
	border: 1px solid rgba(225, 29, 72, 0.15);
}

.role-badge--custom {
	background: rgba(66, 97, 237, 0.08);
	color: var(--primary-color);
	border: 1px solid rgba(66, 97, 237, 0.15);
}

/* ── Permissions section ─────────────────────────────────── */
.role-perm-section {
	margin-bottom: var(--space-3);
	background: var(--bg-lighter);
	padding: 10px 12px;
	border-radius: var(--radius-md);
	border: 1px dashed var(--border-color);
}

.role-perm-header {
	display: flex;
	align-items: center;
	justify-content: space-between;
	margin-bottom: 8px;
}

.perm-label-text {
	font-size: 11px;
	font-weight: var(--fw-semibold);
	color: var(--text-light);
	text-transform: uppercase;
	letter-spacing: 0.05em;
}

.perm-count-badge {
	font-size: 11px;
	font-weight: var(--fw-semibold);
	color: var(--primary-color);
	background: rgba(66, 97, 237, 0.08);
	padding: 2px 7px;
	border-radius: var(--radius-full);
}

.perm-count-badge--full {
	color: #e11d48;
	background: rgba(225, 29, 72, 0.08);
}

.permission-tags {
	display: flex;
	flex-wrap: wrap;
	gap: 5px;
}

.perm-tag {
	padding: 3px 9px;
	background: #ffffff;
	border: 1px solid var(--border-color);
	border-radius: var(--radius-full);
	font-size: 11px;
	font-weight: var(--fw-medium);
	color: var(--text-main);
	box-shadow: 0 1px 2px rgba(0, 0, 0, 0.03);
	white-space: nowrap;
	max-width: 120px;
	overflow: hidden;
	text-overflow: ellipsis;
}

.perm-tag--admin {
	background: rgba(225, 29, 72, 0.05);
	border-color: rgba(225, 29, 72, 0.2);
	color: var(--danger-color);
	font-weight: var(--fw-semibold);
	max-width: none;
}

.perm-tag--more {
	background: rgba(66, 97, 237, 0.06);
	border-color: rgba(66, 97, 237, 0.15);
	color: var(--primary-color);
	cursor: help;
	font-weight: var(--fw-semibold);
}

.perm-tag--empty {
	font-style: italic;
	color: var(--text-light);
	background: transparent;
	border-color: transparent;
	box-shadow: none;
	font-size: 11px;
}


/* ── Empty state ─────────────────────────────────────────── */
.empty-state-container {
	grid-column: 1 / -1;
	display: flex;
	justify-content: center;
	padding: var(--space-8) 0;
}

.empty-state-icon-wrap {
	width: 80px;
	height: 80px;
	border-radius: var(--radius-xl);
	background: rgba(66, 97, 237, 0.06);
	display: flex;
	align-items: center;
	justify-content: center;
	margin-bottom: var(--space-3);
	border: 1px dashed var(--border-hover);
}

.empty-state__icon-svg {
	width: 36px;
	height: 36px;
	color: var(--text-light);
}

.empty-state__title {
	font-size: var(--fs-lg);
	font-weight: var(--fw-bold);
	color: var(--text-main);
	margin: 0 0 8px 0;
}

</style>
