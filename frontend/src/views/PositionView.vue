<script setup>
import {
	Briefcase,
	ChevronLeft,
	ChevronRight,
	Pencil,
	Plus,
	Trash2,
} from "@lucide/vue";

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

// ─── Khởi tạo ────────────────────────────────────────────────────────────────
const positionStore = usePositionStore();
const toast = useToast();

const { canCrudDepartment } = usePermissions();
const canCrudPosition = computed(() => canCrudDepartment.value);

const { positions, loading } = storeToRefs(positionStore);

// ─── Modal thêm/sửa ──────────────────────────────────────────────────────────
const { isModalVisible, isEditMode, openAddModal, openEditModal, closeModal } =
	useModalState();

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
			res = await positionStore.updatePosition(editingPosition.value.id, {
				name: submittedData.name?.trim(),
				description: submittedData.description?.trim(),
			});
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

		toast.success(
			isEditMode.value
				? "Cập nhật chức vụ thành công"
				: "Tạo chức vụ mới thành công",
		);
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
					<span>{{ positions.length }}</span> chức vụ
				</p>
			</div>
			<button
				v-if="canCrudPosition"
				class="btn btn-primary"
				@click="handleAdd"
			>
				<Plus class="btn__icon" />
				Thêm chức vụ
			</button>
		</header>

		<main class="bento-container">
			<!-- Loading Skeletons -->
			<div v-if="loading" class="bento-grid">
				<div v-for="i in 3" :key="'skeleton-' + i" class="bento-card">
					<div class="bento-accent-bar"></div>
					<div class="bento-header">
						<Skeleton type="circle" width="36px" height="36px" />
						<Skeleton type="text" width="60%" height="22px" />
					</div>
					<div class="bento-body">
						<Skeleton
							type="text"
							width="100%"
							height="16px"
							style="margin-bottom: 8px"
						/>
						<Skeleton type="text" width="80%" height="16px" />
					</div>
					<div class="bento-actions" v-if="canCrudPosition">
						<Skeleton type="btn" />
						<Skeleton type="btn" />
					</div>
				</div>
			</div>

			<!-- Actual Data -->
			<div v-else class="bento-grid">
				<div v-for="pos in positions" :key="pos.id" class="bento-card">
					<div class="bento-accent-bar"></div>

					<div class="bento-header">
						<div class="avatar-gradient" style="width: 36px; height: 36px;">
							<Briefcase class="pos-icon" />
						</div>
						<h3 class="bento-name">{{ pos.name }}</h3>
					</div>

					<div class="bento-body">
						<p class="bento-desc">
							{{
								pos.description ||
								"Không có mô tả chi tiết cho chức vụ này."
							}}
						</p>
					</div>

					<div v-if="canCrudPosition" class="bento-actions">
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
				</div>

				<!-- Trạng thái trống -->
				<div
					v-if="positions.length === 0"
					class="empty-state-container"
				>
					<div class="empty-state">
						<Briefcase class="empty-state__icon-svg" />
						<p class="empty-state__text">
							Chưa có chức vụ nào trong hệ thống.
						</p>
						<button
							v-if="canCrudPosition"
							class="btn btn-primary"
							@click="handleAdd"
							style="margin-top: 1rem"
						>
							<Plus class="btn__icon" /> Thêm chức vụ
						</button>
					</div>
				</div>
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

.pos-icon {
	width: 18px;
	height: 18px;
}

.empty-state-container {
	grid-column: 1 / -1;
	display: flex;
	justify-content: center;
	padding: var(--space-6) 0;
}

.empty-state__icon-svg {
	width: 64px;
	height: 64px;
	color: var(--border-hover);
	margin-bottom: var(--space-3);
}
</style>
