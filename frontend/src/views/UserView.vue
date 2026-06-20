<script setup>
import {
	ChevronLeft,
	ChevronRight,
	Pencil,
	Plus,
	Search,
	TriangleAlert,
} from "@lucide/vue";

// ─── Icon SVG ────────────────────────────────────────────────────────────────

// ─── Component UI dùng chung ─────────────────────────────────────────────────
import ConfirmationDialog from "@/components/ConfirmationDialog.vue";
import UserModal from "@/components/UserModal.vue";
import Skeleton from "@/components/Skeleton.vue";

// ─── Store & tiện ích ────────────────────────────────────────────────────────
import { useUserStore } from "@/store/user";
import { useAuthStore } from "@/store/auth";
import { useToast } from "vue-toastification";
import { storeToRefs } from "pinia";
import { usePaginatedSearch } from "@/helpers/usePaginatedSearch";
import { useModalState } from "@/helpers/useModalState";
import { usePermissions } from "@/helpers/usePermissions";
import { onMounted, ref } from "vue";

// ─── Khởi tạo ────────────────────────────────────────────────────────────────

const userStore = useUserStore();
const authStore = useAuthStore();
const currentUser = authStore.user;
const toast = useToast();
const { canManageUsers } = usePermissions();

const { users, pagination, loading } = storeToRefs(userStore);

// ─── Tìm kiếm & phân trang ───────────────────────────────────────────────────

const {
	searchQuery,
	load: loadUsers,
	handlePageChange,
} = usePaginatedSearch((params) => userStore.fetchUser(params), pagination);

// ─── Trạng thái modal xoá ────────────────────────────────────────────────────

const isDeleteModalVisible = ref(false);
const deletingUser = ref(null);
const deleteMessage = ref("");
const deleteLoading = ref(false);
const originalUser = ref(null);

// ─── Trạng thái modal form thêm/sửa ──────────────────────────────────────────

const { isModalVisible, isEditMode: isEditing, openAddModal, openEditModal, closeModal } = useModalState();
const isRoleDisabled = ref(false);
const isActiveDisabled = ref(false);
const submitLoading = ref(false);
const currentUserId = ref(null);

// Form data snapshot for modal initialization
const editingUser = ref(null);

// ─── Mở modal thêm mới ───────────────────────────────────────────────────────

function handleAdd() {
	openAddModal();
	isRoleDisabled.value = false;
	isActiveDisabled.value = false;
	currentUserId.value = null;
	originalUser.value = null;
	editingUser.value = null;
}

// ─── Mở modal sửa ────────────────────────────────────────────────────────────

function handleUpdate(user) {
	openEditModal();
	// So sánh theo role.name thay vì role_id để không bị ảnh hưởng khi DB migrate lệch ID
	const isAdminUser = user.role?.name?.toLowerCase() === 'admin';
	isRoleDisabled.value = isAdminUser || currentUser?.id === user.id;
	isActiveDisabled.value = currentUser?.id === user.id;
	currentUserId.value = user.id;
	originalUser.value = { ...user };
	editingUser.value = user;
}

// ─── Submit form thêm/sửa ─────────────────────────────────────────────────────

async function handleFormSubmit(submittedData) {
	if (submittedData.password !== submittedData.password_confirm) {
		toast.error("Mật khẩu xác nhận không khớp!");
		return;
	}

	const payload = {
		user_name: submittedData.user_name.trim(),
		email: submittedData.email.trim(),
		role_id: Number(submittedData.role_id),
		is_active: Boolean(submittedData.is_active),
	};
	if (submittedData.password) {
		payload.password = submittedData.password;
	}

	submitLoading.value = true;
	let res;

	if (isEditing.value) {
		res = await userStore.updateUser(currentUserId.value, payload);
	} else {
		res = await userStore.createUser(payload);
	}
	submitLoading.value = false;

	if (res.success === false) {
		toast.error(res.message || "Có lỗi xảy ra");
		return;
	}

	toast.success(
		res.message ||
			(isEditing.value ? "Cập nhật thành công!" : "Thêm mới thành công!"),
	);

	isModalVisible.value = false;
	await loadUsers();
}

// ─── Xử lý xoá user ──────────────────────────────────────────────────────────

function handleDelete(user) {
	deletingUser.value = user;
	deleteMessage.value = `Bạn có chắc chắn muốn xoá người dùng ${user.user_name}?`;
	isDeleteModalVisible.value = true;
}

async function confirmDelete() {
	const user = deletingUser.value;
	if (!user) return;

	deleteLoading.value = true;
	const res = await userStore.deleteUser(user.id);
	deleteLoading.value = false;

	if (res.success === false) {
		toast.error(res.message || "Xoá người dùng thất bại");
		return;
	}

	toast.success(res.message || "Xoá người dùng thành công");
	isDeleteModalVisible.value = false;
	deletingUser.value = null;
	await loadUsers();
}

onMounted(async () => {
	await loadUsers();
});
</script>

<template>
	<div class="user-view">
		<header class="page-header">
			<div class="header-content">
				<h1 class="page-title">Quản lý người dùng</h1>
				<p class="page-subtitle">
					Hệ thống có tổng cộng
					<span>{{ pagination.total }}</span> người dùng
				</p>
			</div>
			<button
				v-if="canManageUsers"
				class="btn btn-primary"
				@click="handleAdd"
			>
				<Plus class="btn__icon" />
				Thêm người dùng
			</button>
		</header>

		<main class="content-card">
			<div class="toolbar">
				<div class="search-box">
					<Search class="search-box__icon" />
					<input
						v-model="searchQuery"
						class="form-control search-box__input"
						placeholder="Tìm tên hoặc email người dùng..."
					/>
				</div>
			</div>

			<div class="table-responsive responsive-table-to-cards">
				<table class="data-table">
					<thead>
						<tr>
							<th>Tên người dùng</th>
							<th>Email</th>
							<th>Trạng thái</th>
							<th class="text-right">Thao tác</th>
						</tr>
					</thead>
					<tbody>
						<!-- Loading skeleton rows -->
						<template v-if="loading">
							<tr v-for="i in 5" :key="'skeleton-' + i">
								<td
									class="text-main fw-500"
									data-label="Tên người dùng"
								>
									<Skeleton
										type="text"
										width="130px"
										height="18px"
									/>
								</td>
								<td data-label="Email">
									<Skeleton
										type="text"
										width="200px"
										height="18px"
									/>
								</td>
								<td data-label="Trạng thái">
									<Skeleton
										type="text"
										class="status-badge"
										width="90px"
										height="24px"
										style="display: inline-block"
									/>
								</td>
								<td class="text-right" data-label="Thao tác">
									<div class="action-group">
										<Skeleton type="btn" />
										<Skeleton type="btn" />
									</div>
								</td>
							</tr>
						</template>

						<!-- Actual rows when loaded -->
						<template v-else>
							<tr v-for="user in users" :key="user.id">
								<td data-label="Tên người dùng">
									<div class="user-cell">
										<div class="user-avatar">
											{{
												user.user_name
													.charAt(0)
													.toUpperCase()
											}}
										</div>
										<div class="user-details">
											<span class="user-name-txt">{{
												user.user_name
											}}</span>
											<span
									class="user-role-badge"
									:class="`role-${user.role_id}`"
								>
									{{
										user.role?.name?.toLowerCase() === 'admin'
											? "Quản trị"
											: user.role?.name?.toLowerCase() === 'hr'
												? "Quản lý"
												: "Nhân viên"
									}}
								</span>
										</div>
									</div>
								</td>
								<td data-label="Email">
									<span class="text-main fw-500">{{
										user.email
									}}</span>
								</td>
								<td data-label="Trạng thái">
									<span
										:class="[
											'status-badge',
											user.is_active
												? 'status-badge--active'
												: 'status-badge--inactive',
										]"
									>
										{{
											user.is_active
												? "Hoạt động"
												: "Ngưng"
										}}
									</span>
								</td>
								<td class="text-right" data-label="Thao tác">
									<div class="action-group">
										<button
											v-if="canManageUsers"
											class="btn-icon btn-icon--edit"
											title="Chỉnh sửa"
											@click="handleUpdate(user)"
										>
											<Pencil />
										</button>
										<button
											v-if="canManageUsers"
											class="btn-icon btn-icon--delete"
											title="Xoá"
											@click="handleDelete(user)"
										>
											<TriangleAlert />
										</button>
									</div>
								</td>
							</tr>
							<tr v-if="users.length === 0">
								<td colspan="4" class="empty-state">
									<div class="empty-state__icon">🏢</div>
									<p class="empty-state__text">
										Không có người dùng nào phù hợp.
									</p>
								</td>
							</tr>
						</template>
					</tbody>
				</table>
			</div>

			<!-- Phân trang -->
			<div class="pagination" v-if="pagination.totalPages > 0">
				<button
					class="pagination__btn"
					:disabled="pagination.page === 1"
					@click="handlePageChange(pagination.page - 1)"
				>
					<ChevronLeft />
				</button>
				<div class="pagination__info">
					Trang <span>{{ pagination.page }}</span> /
					{{ pagination.totalPages }}
				</div>
				<button
					class="pagination__btn"
					:disabled="pagination.page === pagination.totalPages"
					@click="handlePageChange(pagination.page + 1)"
				>
					<ChevronRight />
				</button>
			</div>
		</main>

		<ConfirmationDialog
			:visible="isDeleteModalVisible"
			title="Xác nhận xoá người dùng"
			:message="deleteMessage"
			:loading="deleteLoading"
			@confirm="confirmDelete"
			@cancel="isDeleteModalVisible = false"
		/>

		<!-- Form Modal Subcomponent -->
		<UserModal
			:visible="isModalVisible"
			:is-editing="isEditing"
			:editing-user="editingUser"
			:is-role-disabled="isRoleDisabled"
			:is-active-disabled="isActiveDisabled"
			:loading="submitLoading"
			@close="isModalVisible = false"
			@submit="handleFormSubmit"
		/>
	</div>
</template>

<style scoped>
.user-view {
	padding-bottom: var(--space-4);
}

.user-cell {
	display: flex;
	align-items: center;
	gap: var(--space-2);
}

.user-avatar {
	width: 40px;
	height: 40px;
	border-radius: var(--radius-md);
	background: linear-gradient(
		135deg,
		rgba(0, 192, 250, 0.12) 0%,
		rgba(66, 97, 237, 0.1) 100%
	);
	color: var(--primary-color);
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: var(--fs-sm);
	font-weight: var(--fw-bold);
	flex-shrink: 0;
}

.user-details {
	display: flex;
	flex-direction: column;
	gap: 2px;
}

.user-name-txt {
	font-size: var(--fs-sm);
	font-weight: var(--fw-semibold);
}

.user-role-badge {
	font-size: var(--fs-xs);
	font-weight: var(--fw-semibold);
	padding: 0.15rem 0.5rem;
	border-radius: var(--radius-sm);
	width: max-content;
}

.user-role-badge.role-1 {
	background: rgba(139, 92, 246, 0.1);
	color: var(--color-purple);
}

.user-role-badge.role-2 {
	background: rgba(217, 119, 6, 0.1);
	color: var(--color-amber);
}

.user-role-badge.role-3 {
	background: rgba(66, 97, 237, 0.1);
	color: var(--primary-color);
}
</style>
