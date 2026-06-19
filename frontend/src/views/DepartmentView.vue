<script setup>
// ─── Icon SVG ────────────────────────────────────────────────────────────────
import plusIcon from "@/assets/svg/plus.svg";
import searchIcon from "@/assets/svg/search.svg";
import editIcon from "@/assets/svg/edit.svg";
import deleteIcon from "@/assets/svg/delete.svg";
import prevIcon from "@/assets/svg/chevron-left.svg";
import nextIcon from "@/assets/svg/chevron-right.svg";

// ─── Store ───────────────────────────────────────────────────────────────────
import { useDepartmentStore } from "@/store/department";
import { useDashboardStore } from "@/store/dashboard";
import { useEmployeeStore } from "@/store/employee";

// ─── Component UI ────────────────────────────────────────────────────────────
import DepartmentModal from "@/components/DepartmentModal.vue";
import ConfirmationDialog from "@/components/ConfirmationDialog.vue";
import Skeleton from "@/components/Skeleton.vue";

// ─── Tiện ích ────────────────────────────────────────────────────────────────
import { ref, onMounted, computed } from "vue";
import { storeToRefs } from "pinia";
import { useToast } from "vue-toastification";
import { useModalState } from "@/helpers/useModalState";
import { usePaginatedSearch } from "@/helpers/usePaginatedSearch";
import { usePermissions } from "@/helpers/usePermissions";
import { buildPatchPayload } from "@/helpers/buildPatchPayload";
import { getInitials } from "@/helpers/formatters";

// ─── Khởi tạo ────────────────────────────────────────────────────────────────
const departmentStore = useDepartmentStore();
const dashboardStore = useDashboardStore();
const employeeStore = useEmployeeStore();
const toast = useToast();

const { canCrudDepartment } = usePermissions();
const { departments, loading, pagination } = storeToRefs(departmentStore);
const employeeOptions = ref([]);

// ─── Modal thêm/sửa ──────────────────────────────────────────────────────────
const { isModalVisible, isEditMode, openAddModal, openEditModal, closeModal } =
	useModalState();

// ─── Tìm kiếm & phân trang ───────────────────────────────────────────────────
const {
	searchQuery,
	load: loadDepartments,
	handlePageChange,
} = usePaginatedSearch(
	(params) => departmentStore.fetchDepartments(params),
	pagination,
);

// ─── Trạng thái local ────────────────────────────────────────────────────────
const editingDepartment = ref(null);
const formLoading = ref(false);

// ─── Trạng thái modal xoá ────────────────────────────────────────────────────
const isDeleteModalVisible = ref(false);
const deletingDepartment = ref(null);
const deleteMessage = ref("");
const deleteLoading = ref(false);

// ─── Tải danh sách nhân viên ─────────────────────────────────────────────────
async function loadEmployees() {
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
}

const departmentEmployees = computed(() => {
	if (!editingDepartment.value) {
		return [];
	}
	return employeeOptions.value.filter(
		(emp) => emp.department_id === editingDepartment.value.id,
	);
});

// ─── Mở modal thêm mới ───────────────────────────────────────────────────────
function handleAdd() {
	editingDepartment.value = null;
	openAddModal();
}

// ─── Mở modal sửa ────────────────────────────────────────────────────────────
async function handleEdit(department) {
	editingDepartment.value = { ...department };
	openEditModal();
	await loadEmployees();
}

// ─── Mở modal xoá ────────────────────────────────────────────────────────────
function handleDelete(department) {
	deletingDepartment.value = department;
	deleteMessage.value = `Bạn có chắc chắn muốn xoá phòng ban ${department.name}?`;
	isDeleteModalVisible.value = true;
}

// ─── Xác nhận xoá ────────────────────────────────────────────────────────────
async function confirmDelete() {
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
}

// ─── Submit form thêm/sửa ─────────────────────────────────────────────────────
async function handleFormSubmit(submittedData) {
	formLoading.value = true;
	try {
		let res;

		if (isEditMode.value) {
			const original = {
				name: editingDepartment.value?.name ?? "",
				description: editingDepartment.value?.description ?? "",
				manager_id:
					editingDepartment.value?.manager_id ??
					editingDepartment.value?.manager?.id ??
					null,
			};
			const payload = buildPatchPayload(original, submittedData, {
				fields: ["name", "description", "manager_id"],
				transformValue: (key, value) => {
					if (key === "manager_id") {
						return value == null || value === ""
							? 0
							: Number(value);
					}
					return value;
				},
			});

			if (Object.keys(payload).length === 0) {
				toast.info("Không có dữ liệu thay đổi");
				formLoading.value = false;
				return;
			}

			res = await departmentStore.updateDepartment(
				editingDepartment.value.id,
				payload,
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
}

// ─── Khởi tạo trang ──────────────────────────────────────────────────────────
onMounted(async () => {
	await Promise.all([loadDepartments(), dashboardStore.fetchDashboard()]);
});
</script>

<template>
	<div class="department-view">
		<header class="page-header">
			<div class="header-content">
				<h1 class="page-title">Quản lý phòng ban</h1>
				<p class="page-subtitle">
					Hệ thống có tổng cộng
					<span>{{ pagination.total || departments.length }}</span> phòng ban
				</p>
			</div>
			<button
				v-if="canCrudDepartment"
				class="btn btn--primary"
				@click="handleAdd"
			>
				<img :src="plusIcon" alt="add" class="btn__icon" />
				Thêm phòng ban
			</button>
		</header>

		<!-- Toolbar Search Card -->
		<div class="search-card">
			<div class="search-box">
				<img
					:src="searchIcon"
					class="search-box__icon"
					alt="search"
				/>
				<input
					v-model="searchQuery"
					class="form-control search-box__input"
					placeholder="Tìm tên hoặc mã phòng ban..."
				/>
			</div>
		</div>

		<!-- Bento Grid of Departments -->
		<main class="dept-bento-container">
			<!-- Loading skeleton grid -->
			<div v-if="loading" class="dept-bento-grid">
				<div v-for="i in 6" :key="'skeleton-dept-' + i" class="dept-bento-card">
					<div class="dept-bento-header">
						<Skeleton type="text" width="60%" height="22px" />
						<Skeleton type="badge" width="60px" height="24px" />
					</div>
					<div class="dept-bento-body">
						<Skeleton type="text" width="100%" height="16px" style="margin-bottom: 8px" />
						<Skeleton type="text" width="80%" height="16px" />
					</div>
					<div class="dept-bento-manager">
						<Skeleton type="avatar" />
						<div style="flex: 1; display: flex; flex-direction: column; gap: var(--space-1); margin-left: 8px;">
							<Skeleton type="text" width="40%" />
							<Skeleton type="text" width="70%" />
						</div>
					</div>
					<div v-if="canCrudDepartment" class="dept-bento-actions">
						<Skeleton type="btn" />
						<Skeleton type="btn" />
					</div>
				</div>
			</div>

			<!-- Actual department cards -->
			<div v-else class="dept-bento-grid">
				<div v-for="dept in departments" :key="dept.id" class="dept-bento-card">
					<!-- Accent light border at top -->
					<div class="dept-accent-bar"></div>
					
					<div class="dept-bento-header">
						<h3 class="dept-bento-name">{{ dept.name }}</h3>
						<span class="dept-code">{{ dept.code }}</span>
					</div>
					
					<div class="dept-bento-body">
						<p class="dept-bento-desc">
							{{ dept.description || "Không có mô tả chi tiết cho phòng ban này." }}
						</p>
					</div>
					
					<div class="dept-bento-manager">
						<div class="manager-avatar">
							{{ dept.manager ? getInitials(dept.manager.first_name, dept.manager.last_name) : '?' }}
						</div>
						<div class="manager-info">
							<span class="manager-label">Trưởng phòng</span>
							<span class="manager-name">
								{{ dept.manager ? dept.manager.first_name + " " + dept.manager.last_name : "Chưa bổ nhiệm" }}
							</span>
						</div>
					</div>
					
					<div v-if="canCrudDepartment" class="dept-bento-actions">
						<button
							class="btn-icon btn-icon--edit"
							title="Chỉnh sửa"
							@click="handleEdit(dept)"
						>
							<img :src="editIcon" alt="edit" />
						</button>
						<button
							class="btn-icon btn-icon--delete"
							title="Xoá"
							@click="handleDelete(dept)"
						>
							<img :src="deleteIcon" alt="delete" />
						</button>
					</div>
				</div>

				<!-- Empty State -->
				<div v-if="departments.length === 0" class="empty-state-container">
					<div class="empty-state">
						<div class="empty-state__icon">🏢</div>
						<p class="empty-state__text">
							Không tìm thấy phòng ban nào phù hợp.
						</p>
					</div>
				</div>
			</div>

			<!-- Pagination controls (Glass styled) -->
			<div class="pagination-container" v-if="pagination.totalPages > 0">
				<div class="pagination">
					<button
						class="pagination__btn"
						:disabled="pagination.page === 1"
						@click="handlePageChange(pagination.page - 1)"
					>
						<img :src="prevIcon" alt="prev" />
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
						<img :src="nextIcon" alt="next" />
					</button>
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
	</div>
</template>

<style scoped>
.department-view {
	padding-bottom: var(--space-4);
	display: flex;
	flex-direction: column;
	gap: var(--space-3);
}

/* Search bar glass container */
.search-card {
	background: var(--bg-card);
	backdrop-filter: var(--glass-backdrop);
	-webkit-backdrop-filter: var(--glass-backdrop);
	border: var(--glass-border);
	border-radius: var(--radius-md);
	box-shadow: var(--shadow-sm);
	padding: var(--space-3);
}

.search-box {
	position: relative;
	width: 100%;
	max-width: 400px;
}

.search-box__icon {
	position: absolute;
	left: 1rem;
	top: 50%;
	transform: translateY(-50%);
	width: 18px;
	height: 18px;
	filter: grayscale(1) opacity(0.5);
	z-index: 1;
}

.search-box__input {
	padding-left: 2.75rem !important;
}

/* Bento Grid */
.dept-bento-container {
	display: flex;
	flex-direction: column;
	gap: var(--space-4);
}

.dept-bento-grid {
	display: grid;
	grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
	gap: var(--space-3);
}

.dept-bento-card {
	background: var(--bg-card);
	backdrop-filter: var(--glass-backdrop);
	-webkit-backdrop-filter: var(--glass-backdrop);
	border: var(--glass-border);
	border-radius: var(--radius-lg);
	box-shadow: var(--glass-shadow);
	padding: var(--space-3);
	transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
	position: relative;
	overflow: hidden;
	display: flex;
	flex-direction: column;
	gap: var(--space-3);
}

.dept-bento-card:hover {
	transform: translateY(-4px);
	box-shadow: var(--glass-shadow-hover);
}

.dept-accent-bar {
	position: absolute;
	top: 0;
	left: 0;
	right: 0;
	height: 4px;
	background: var(--primary-gradient);
}

.dept-bento-header {
	display: flex;
	justify-content: space-between;
	align-items: flex-start;
	gap: var(--space-2);
	margin-top: 4px;
}

.dept-bento-name {
	font-family: var(--font-title);
	font-size: var(--fs-base);
	font-weight: var(--fw-bold);
	color: var(--text-main);
	margin: 0;
	line-height: var(--lh-tight);
}

.dept-bento-body {
	flex: 1;
}

.dept-bento-desc {
	font-size: var(--fs-sm);
	color: var(--text-muted);
	line-height: var(--lh-normal);
	margin: 0;
	display: -webkit-box;
	-webkit-line-clamp: 3;
	line-clamp: 3;
	-webkit-box-orient: vertical;
	overflow: hidden;
}

/* Manager Card Sub-Component */
.dept-bento-manager {
	display: flex;
	align-items: center;
	gap: var(--space-2);
	background: rgba(255, 255, 255, 0.45);
	padding: var(--space-2);
	border-radius: var(--radius-md);
	border: 1px solid rgba(66, 97, 237, 0.08);
}

.manager-avatar {
	width: 36px;
	height: 36px;
	border-radius: var(--radius-md);
	background: linear-gradient(135deg, rgba(0, 192, 250, 0.12) 0%, rgba(66, 97, 237, 0.1) 100%);
	color: var(--primary-color);
	display: flex;
	align-items: center;
	justify-content: center;
	font-weight: var(--fw-bold);
	font-size: var(--fs-xs);
	flex-shrink: 0;
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

/* Actions at the bottom of card */
.dept-bento-actions {
	display: flex;
	justify-content: flex-end;
	gap: 0.5rem;
	border-top: 1px solid var(--border-color);
	padding-top: var(--space-2);
	margin-top: auto;
}

.empty-state-container {
	grid-column: 1 / -1;
	padding: var(--space-8) 0;
}

/* Pagination container */
.pagination-container {
	display: flex;
	justify-content: center;
	margin-top: var(--space-2);
}

.pagination {
	display: flex;
	justify-content: center;
	align-items: center;
	gap: var(--space-3);
	padding: var(--space-2) var(--space-4);
	background: var(--bg-card);
	backdrop-filter: var(--glass-backdrop);
	-webkit-backdrop-filter: var(--glass-backdrop);
	border: var(--glass-border);
	border-radius: var(--radius-md);
	box-shadow: var(--shadow-sm);
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
