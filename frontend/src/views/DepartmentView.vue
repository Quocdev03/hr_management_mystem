<script setup>
import {
	Building2,
	Eye,
	Pencil,
	Plus,
	Trash2,
} from "@lucide/vue";

import { useDepartmentStore } from "@/store/department";

import { useEmployeeStore } from "@/store/employee";

// ─── Component UI
import DepartmentModal from "@/components/DepartmentModal.vue";
import DepartmentDetailModal from "@/components/DepartmentDetailModal.vue";
import ConfirmationDialog from "@/components/ConfirmationDialog.vue";
import Skeleton from "@/components/Skeleton.vue";

// ─── Tiện ích ────────────────────────────────────────────────────────────────
import { ref, onMounted, computed } from "vue";
import { storeToRefs } from "pinia";
import { useToast } from "vue-toastification";
import { useModalState } from "@/helpers/useModalState";
import { usePermissions } from "@/helpers/usePermissions";
import { getInitials } from "@/helpers/formatters";

// ─── Khởi tạo ────────────────────────────────────────────────────────────────
const departmentStore = useDepartmentStore();
const employeeStore = useEmployeeStore();
const toast = useToast();

const { canCrudDepartment } = usePermissions();
const { departments, loading, pagination } = storeToRefs(departmentStore);
const employeeOptions = ref([]);

// ─── Trạng thái chi tiết phòng ban ──────────────────────────────────────────
const isDetailModalVisible = ref(false);
const selectedDepartment = ref(null);
const detailLoading = ref(false);

// ─── Modal thêm/sửa ──────────────────────────────────────────────────────────
const { isModalVisible, isEditMode, openAddModal, openEditModal, closeModal } =
	useModalState();

// ─── Tải danh sách phòng ban ─────────────────────────────────────────────────
const loadDepartments = async () => {
	try {
		await departmentStore.fetchDepartments({ limit: 100 });
	} catch (err) {
		toast.error("Không thể tải danh sách phòng ban");
		console.error("loadDepartments error:", err);
	}
};

// ─── Trạng thái local ────────────────────────────────────────────────────────
const editingDepartment = ref(null);
const formLoading = ref(false);

// ─── Trạng thái modal xoá ────────────────────────────────────────────────────
const isDeleteModalVisible = ref(false);
const deletingDepartment = ref(null);
const deleteMessage = ref("");
const deleteLoading = ref(false);

// ─── Tải danh sách nhân viên ─────────────────────────────────────────────────
const loadEmployees = async () => {
	try {
		const res = await employeeStore.fetchEmployeesForSelect({
			page: 1,
			limit: 100,
		});
		if (res.success) {
			employeeOptions.value = res.items;
		}
	} catch (err) {
		console.error("Lỗi khi tải danh sách nhân viên:", err);
	}
};

const departmentEmployees = computed(() => {
	if (!editingDepartment.value) {
		return [];
	}
	return employeeOptions.value.filter(
		(emp) =>
			emp.department_id === editingDepartment.value.id &&
			emp.status === "active",
	);
});

// ─── Mở modal thêm mới ───────────────────────────────────────────────────────
const handleAdd = () => {
	editingDepartment.value = null;
	openAddModal();
};

// ─── Mở modal sửa ────────────────────────────────────────────────────────────
const handleEdit = async (department) => {
	editingDepartment.value = { ...department };
	openEditModal();
	await loadEmployees();
};

// ─── Mở modal xoá ────────────────────────────────────────────────────────────
const handleDelete = (department) => {
	deletingDepartment.value = department;
	deleteMessage.value = `Bạn có chắc chắn muốn xoá phòng ban ${department.name}?`;
	isDeleteModalVisible.value = true;
};

// ─── Xem chi tiết phòng ban ──────────────────────────────────────────────────
const handleViewDetail = async (departmentId) => {
	detailLoading.value = true;
	try {
		const res = await departmentStore.fetchDepartmentByID(departmentId);
		if (res.success) {
			selectedDepartment.value = res.data;
			isDetailModalVisible.value = true;
		} else {
			toast.error(res.message || "Không thể tải chi tiết phòng ban");
		}
	} catch (err) {
		toast.error("Đã xảy ra lỗi khi tải chi tiết phòng ban");
		console.error("handleViewDetail error:", err);
	} finally {
		detailLoading.value = false;
	}
};

// ─── Xác nhận xoá ────────────────────────────────────────────────────────────
const confirmDelete = async () => {
	const department = deletingDepartment.value;
	if (!department) return;

	deleteLoading.value = true;
	try {
		const res = await departmentStore.deleteDepartment(department.id);

		if (!res.success) {
			toast.error(res.message);
			return;
		}

		toast.success("Xoá phòng ban thành công");
		isDeleteModalVisible.value = false;
		deletingDepartment.value = null;
		await loadDepartments();
	} catch (err) {
		toast.error(err?.message || "Đã xảy ra lỗi khi xoá phòng ban");
		console.error("confirmDelete error:", err);
	} finally {
		deleteLoading.value = false;
	}
};

// ─── Submit form thêm/sửa ─────────────────────────────────────────────────────
const handleFormSubmit = async (submittedData) => {
	formLoading.value = true;
	try {
		let res;

		if (isEditMode.value) {
			res = await departmentStore.updateDepartment(
				editingDepartment.value.id,
				{
					name: submittedData.name?.trim(),
					description: submittedData.description?.trim(),
					manager_id:
						submittedData.manager_id == null ||
						submittedData.manager_id === ""
							? 0
							: Number(submittedData.manager_id),
				},
			);
		} else {
			res = await departmentStore.createDepartment({
				name: submittedData.name,
				code: submittedData.code,
				description: submittedData.description,
				manager_id: submittedData.manager_id || null,
			});
		}

		if (!res.success) {
			toast.error(res.message);
			return;
		}

		if (isEditMode.value) {
			toast.success(res.message);
		}

		closeModal();
		await loadDepartments();
	} catch (err) {
		toast.error(err?.message || "Đã xảy ra lỗi khi lưu phòng ban");
		console.error("handleFormSubmit error:", err);
	} finally {
		formLoading.value = false;
	}
};

// ─── Khởi tạo trang ──────────────────────────────────────────────────────────
onMounted(async () => {
	await loadDepartments();
});
</script>

<template>
	<div class="department-view">
		<header class="page-header">
			<div class="header-content">
				<h1 class="page-title">Quản lý phòng ban</h1>
				<p class="page-subtitle">
					Hệ thống có tổng cộng
					<span>{{ pagination.total || departments.length }}</span>
					phòng ban
				</p>
			</div>
			<button
				v-if="canCrudDepartment"
				class="btn btn-primary"
				@click="handleAdd"
			>
				<Plus class="btn__icon" />
				Thêm phòng ban
			</button>
		</header>

		<!-- Bento Grid of Departments -->
		<main class="bento-container">
			<!-- Loading skeleton grid -->
			<div v-if="loading" class="bento-grid">
				<div
					v-for="i in 3"
					:key="'skeleton-dept-' + i"
					class="bento-card"
				>
					<div class="bento-header">
						<Skeleton type="text" width="60%" height="22px" />
						<Skeleton type="badge" width="60px" height="24px" />
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
					<div class="dept-bento-manager">
						<Skeleton type="avatar" />
						<div class="dept-skeleton-info">
							<Skeleton type="text" width="40%" />
							<Skeleton type="text" width="70%" />
						</div>
					</div>
					<div class="bento-actions dept-bento-actions">
						<Skeleton type="btn" />
						<Skeleton type="btn" v-if="canCrudDepartment" />
						<Skeleton type="btn" v-if="canCrudDepartment" />
					</div>
				</div>
			</div>

			<!-- Actual department cards -->
			<div v-else class="bento-grid">
				<div
					v-for="dept in departments"
					:key="dept.id"
					class="bento-card"
				>
					<!-- Accent light border at top -->
					<div class="bento-accent-bar"></div>

					<div class="bento-header">
						<h3 class="bento-name">{{ dept.name }}</h3>
						<span class="dept-code">{{ dept.code }}</span>
					</div>

					<div class="bento-body">
						<p class="bento-desc">
							{{
								dept.description ||
								"Không có mô tả chi tiết cho phòng ban này."
							}}
						</p>
					</div>

					<div class="dept-bento-manager">
						<div class="avatar-gradient dept-avatar">
							{{
								dept.manager
									? getInitials(
											dept.manager.first_name,
											dept.manager.last_name,
										)
									: "?"
							}}
						</div>
						<div class="manager-info">
							<span class="manager-label">Trưởng phòng</span>
							<span class="manager-name">
								{{
									dept.manager
										? dept.manager.first_name +
											" " +
											dept.manager.last_name
										: "Chưa bổ nhiệm"
								}}
							</span>
						</div>
					</div>

					<div class="bento-actions dept-bento-actions">
						<button
							class="btn-icon btn-icon--detail"
							title="Xem chi tiết"
							@click="handleViewDetail(dept.id)"
						>
							<Eye />
						</button>
						<button
							v-if="canCrudDepartment"
							class="btn-icon btn-icon--edit"
							title="Chỉnh sửa"
							@click="handleEdit(dept)"
						>
							<Pencil />
						</button>
						<button
							v-if="canCrudDepartment"
							class="btn-icon btn-icon--delete"
							title="Xoá"
							@click="handleDelete(dept)"
						>
							<Trash2 />
						</button>
					</div>
				</div>

				<!-- Empty State -->
				<div
					v-if="departments.length === 0"
					class="empty-state-container"
				>
					<div class="empty-state">
						<Building2 class="empty-state__icon-svg" />
						<p class="empty-state__text">
							Không tìm thấy phòng ban nào phù hợp.
						</p>
					</div>
				</div>
			</div>
		</main>

		<!-- Subcomponent Form Modal -->
		<DepartmentModal
			:visible="isModalVisible"
			:is-edit-mode="isEditMode"
			:editing-department="editingDepartment"
			:department-employees="departmentEmployees"
			:loading="formLoading"
			@close="closeModal"
			@submit="handleFormSubmit"
		/>

		<!-- Confirmation Delete Dialog -->
		<ConfirmationDialog
			:visible="isDeleteModalVisible"
			title="Xác nhận xoá phòng ban"
			:message="deleteMessage"
			:loading="deleteLoading"
			@confirm="confirmDelete"
			@cancel="isDeleteModalVisible = false"
		/>

		<!-- Department Detail Modal -->
		<DepartmentDetailModal
			:visible="isDetailModalVisible"
			:department="selectedDepartment"
			@close="isDetailModalVisible = false"
		/>
	</div>
</template>

<style scoped>
.department-view {
	padding-bottom: var(--space-4);
	display: flex;
	flex-direction: column;
}

.dept-avatar {
	width: 36px;
	height: 36px;
	font-size: var(--fs-xs);
}

.dept-skeleton-info {
	flex: 1;
	display: flex;
	flex-direction: column;
	gap: var(--space-1);
	margin-left: 8px;
}

.dept-bento-actions {
	margin-top: var(--space-3);
}

.dept-bento-manager {
	display: flex;
	align-items: center;
	gap: var(--space-2);
	background: rgba(255, 255, 255, 0.45);
	padding: var(--space-2);
	border-radius: var(--radius-md);
	border: 1px solid var(--border-color);
}

.manager-info {
	display: flex;
	flex-direction: column;
	overflow: hidden;
}

.manager-label {
	font-size: var(--fs-xs);
	text-transform: uppercase;
	color: var(--text-muted);
	font-weight: var(--fw-semibold);
	letter-spacing: 0.05em;
}

.manager-name {
	font-size: var(--fs-sm);
	font-weight: var(--fw-semibold);
	color: var(--text-main);
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

.empty-state-container {
	grid-column: 1 / -1;
	padding: var(--space-8) 0;
}

.empty-state__icon-svg {
	width: 48px;
	height: 48px;
	color: var(--text-light);
	margin-bottom: var(--space-2);
	display: inline-block;
}

@media (max-width: 640px) {
	.dept-bento-grid {
		grid-template-columns: 1fr;
	}
	.search-box {
		max-width: none;
	}
}
</style>
