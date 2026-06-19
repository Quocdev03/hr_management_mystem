<script setup>
import { Briefcase, ChevronLeft, ChevronRight, Pencil, Plus, Search, Trash2 } from '@lucide/vue';

// ─── Store ───────────────────────────────────────────────────────────────────
import { usePositionStore } from "@/store/position";

// ─── Component UI ────────────────────────────────────────────────────────────
import PositionModal from "@/components/PositionModal.vue";
import ConfirmationDialog from "@/components/ConfirmationDialog.vue";
import Skeleton from "@/components/Skeleton.vue";

// ─── Tiện ích ────────────────────────────────────────────────────────────────
import { ref, onMounted, computed } from "vue";
import { storeToRefs } from "pinia";
import { useToast } from "vue-toastification";
import { useModalState } from "@/helpers/useModalState";
import { usePermissions } from "@/helpers/usePermissions";
import { buildPatchPayload } from "@/helpers/buildPatchPayload";

// ─── Khởi tạo ────────────────────────────────────────────────────────────────
const positionStore = usePositionStore();
const toast = useToast();

const { canCrudDepartment } = usePermissions();
// Dùng chung quyền CrudDepartment cho Positions như backend config (RequirePermission("department.update"))
const canCrudPosition = computed(() => canCrudDepartment.value);

const { positions, loading } = storeToRefs(positionStore);

// ─── Modal thêm/sửa ──────────────────────────────────────────────────────────
const { isModalVisible, isEditMode, openAddModal, openEditModal, closeModal } =
	useModalState();

// ─── Tìm kiếm cục bộ (Client-side Search) ──────────────────────────────────────
const searchQuery = ref("");

const filteredPositions = computed(() => {
	const query = searchQuery.value.trim().toLowerCase();
	if (!query) return positions.value;
	return positions.value.filter(
		(pos) =>
			pos.name.toLowerCase().includes(query) ||
			(pos.description && pos.description.toLowerCase().includes(query))
	);
});

// ─── Trạng thái local ────────────────────────────────────────────────────────
const editingPosition = ref(null);
const formLoading = ref(false);

// ─── Trạng thái modal xoá ────────────────────────────────────────────────────
const isDeleteModalVisible = ref(false);
const deletingPosition = ref(null);
const deleteMessage = ref("");
const deleteLoading = ref(false);

// ─── Tải danh sách chức vụ ─────────────────────────────────────────────────
async function loadPositions() {
	try {
		await positionStore.fetchPositions();
	} catch (err) {
		toast.error("Không thể tải danh sách chức vụ");
		console.error("loadPositions error:", err);
	}
}

// ─── Mở modal thêm mới ───────────────────────────────────────────────────────
function handleAdd() {
	editingPosition.value = null;
	openAddModal();
}

// ─── Mở modal sửa ────────────────────────────────────────────────────────────
function handleEdit(position) {
	editingPosition.value = { ...position };
	openEditModal();
}

// ─── Mở modal xoá ────────────────────────────────────────────────────────────
function handleDelete(position) {
	deletingPosition.value = position;
	deleteMessage.value = `Bạn có chắc chắn muốn xoá chức vụ "${position.name}"?`;
	isDeleteModalVisible.value = true;
}

// ─── Xác nhận xoá ────────────────────────────────────────────────────────────
async function confirmDelete() {
	const position = deletingPosition.value;
	if (!position) return;

	deleteLoading.value = true;
	try {
		const res = await positionStore.deletePosition(position.id);

		if (!res.success) {
			toast.error(res.message || "Xoá chức vụ thất bại");
			return;
		}

		toast.success("Xoá chức vụ thành công");
		isDeleteModalVisible.value = false;
		deletingPosition.value = null;
		await loadPositions();
	} catch (err) {
		toast.error(err?.message || "Đã xảy ra lỗi khi xoá chức vụ");
		console.error("confirmDelete error:", err);
	} finally {
		deleteLoading.value = false;
	}
}

// ─── Submit form thêm/sửa ─────────────────────────────────────────────────────
async function handleFormSubmit(submittedData) {
	formLoading.value = true;
	try {
		let res;

		if (isEditMode.value) {
			const original = {
				name: editingPosition.value?.name ?? "",
				description: editingPosition.value?.description ?? "",
			};
			const payload = buildPatchPayload(original, submittedData, {
				fields: ["name", "description"],
			});

			if (Object.keys(payload).length === 0) {
				toast.info("Không có thay đổi nào được thực hiện");
				closeModal();
				return;
			}

			res = await positionStore.updatePosition(editingPosition.value.id, payload);
		} else {
			res = await positionStore.createPosition({
				name: submittedData.name.trim(),
				description: submittedData.description.trim(),
			});
		}

		if (!res.success) {
			toast.error(res.message || "Thao tác thất bại");
			return;
		}

		toast.success(isEditMode.value ? "Cập nhật chức vụ thành công" : "Tạo chức vụ mới thành công");
		closeModal();
		await loadPositions();
	} catch (err) {
		toast.error(err?.message || "Đã xảy ra lỗi khi lưu thông tin");
		console.error("handleFormSubmit error:", err);
	} finally {
		formLoading.value = false;
	}
}

onMounted(async () => {
	await loadPositions();
});
</script>

<template>
	<div class="position-view">
		<header class="page-header">
			<div class="header-content">
				<h1 class="page-title">Quản lý chức vụ</h1>
				<p class="page-subtitle">
					Hệ thống có tổng cộng
					<span>{{ filteredPositions.length }}</span> chức vụ
				</p>
			</div>
			<button
				v-if="canCrudPosition"
				class="btn btn--primary"
				@click="handleAdd"
			>
				<Plus class="btn__icon" />
				Thêm chức vụ
			</button>
		</header>

		<main class="content-card">
			<!-- Công cụ tìm kiếm -->
			<div class="toolbar">
				<div class="search-box">
					<Search class="search-box__icon" />
					<input
						v-model="searchQuery"
						class="form-control search-box__input"
						placeholder="Tìm tên chức vụ hoặc mô tả..."
					/>
				</div>
			</div>

			<!-- Bảng hiển thị -->
			<div class="table-responsive responsive-table-to-cards">
				<table class="data-table">
					<thead>
						<tr>
							<th>Tên chức vụ</th>
							<th>Mô tả</th>
							<th v-if="canCrudPosition" class="text-right">Thao tác</th>
						</tr>
					</thead>
					<tbody>
						<!-- Loading Skeletons -->
						<template v-if="loading">
							<tr v-for="i in 4" :key="'skeleton-' + i">
								<td data-label="Tên chức vụ">
									<div class="position-cell">
										<Skeleton type="circle" width="36px" height="36px" />
										<Skeleton type="text" width="150px" height="18px" />
									</div>
								</td>
								<td data-label="Mô tả">
									<Skeleton type="text" width="300px" height="18px" />
								</td>
								<td v-if="canCrudPosition" class="text-right" data-label="Thao tác">
									<div class="action-group">
										<Skeleton type="btn" />
										<Skeleton type="btn" />
									</div>
								</td>
							</tr>
						</template>

						<!-- Dữ liệu thực tế -->
						<template v-else>
							<tr v-for="pos in filteredPositions" :key="pos.id">
								<td data-label="Tên chức vụ">
									<div class="position-cell">
										<div class="position-icon-container">
											<Briefcase class="pos-icon" />
										</div>
										<span class="position-name-text">{{ pos.name }}</span>
									</div>
								</td>
								<td data-label="Mô tả">
									<span class="text-muted description-text">{{ pos.description || 'Chưa có mô tả' }}</span>
								</td>
								<td v-if="canCrudPosition" class="text-right" data-label="Thao tác">
									<div class="action-group">
										<button
											class="btn-icon btn-icon--edit"
											title="Chỉnh sửa"
											@click="handleEdit(pos)"
										>
											<Pencil />
										</button>
										<button
											class="btn-icon btn-icon--delete"
											title="Xoá"
											@click="handleDelete(pos)"
										>
											<Trash2 />
										</button>
									</div>
								</td>
							</tr>

							<!-- Trạng thái trống -->
							<tr v-if="filteredPositions.length === 0">
								<td colspan="3" class="empty-state">
									<div class="empty-state__icon">💼</div>
									<p class="empty-state__text">
										Không tìm thấy chức vụ nào.
									</p>
								</td>
							</tr>
						</template>
					</tbody>
				</table>
			</div>
		</main>

		<!-- Dialog xác nhận xóa -->
		<ConfirmationDialog
			:visible="isDeleteModalVisible"
			title="Xác nhận xoá chức vụ"
			:message="deleteMessage"
			:loading="deleteLoading"
			@confirm="confirmDelete"
			@cancel="isDeleteModalVisible = false"
		/>

		<!-- Modal thêm/sửa -->
		<PositionModal
			:visible="isModalVisible"
			:is-edit-mode="isEditMode"
			:editing-position="editingPosition"
			:loading="formLoading"
			@close="closeModal"
			@submit="handleFormSubmit"
		/>
	</div>
</template>

<style scoped>
.position-view {
	padding-bottom: var(--space-4);
}

.position-cell {
	display: flex;
	align-items: center;
	gap: var(--space-3);
}

.position-icon-container {
	display: flex;
	align-items: center;
	justify-content: center;
	width: 36px;
	height: 36px;
	border-radius: var(--radius-md);
	background: linear-gradient(
		135deg,
		rgba(66, 97, 237, 0.1) 0%,
		rgba(0, 192, 250, 0.05) 100%
	);
	color: var(--primary-color);
	flex-shrink: 0;
}

.pos-icon {
	width: 18px;
	height: 18px;
}

.position-name-text {
	font-size: var(--fs-sm);
	font-weight: var(--fw-semibold);
	color: var(--text-main);
}

.description-text {
	font-size: var(--fs-sm);
	display: -webkit-box;
	-webkit-line-clamp: 2;
	line-clamp: 2;
	-webkit-box-orient: vertical;
	overflow: hidden;
}
</style>
