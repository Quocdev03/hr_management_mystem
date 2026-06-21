<script setup>
import {
	ChevronLeft,
	ChevronRight,
	Eye,
	Pencil,
	Plus,
	Search,
	Trash2,
} from "@lucide/vue";

// ─── Store & tiện ích ────────────────────────────────────────────────────────
import { ref, onMounted } from "vue";
import { storeToRefs } from "pinia";
import { useToast } from "vue-toastification";

import { useEmployeeStore } from "@/store/employee";
import { useDepartmentStore } from "@/store/department";
import { useUserStore } from "@/store/user";

// ─── Component UI dùng chung ─────────────────────────────────────────────────
import EmployeeModal from "@/components/EmployeeModal.vue";
import EmployeeDetailModal from "@/components/EmployeeDetailModal.vue";
import ConfirmationDialog from "@/components/ConfirmationDialog.vue";
import Skeleton from "@/components/Skeleton.vue";

// Helper tiện ích
import { getInitials, formatDate, formatStatus } from "@/helpers/formatters";
import { useModalState } from "@/helpers/useModalState";
import { usePaginatedSearch } from "@/helpers/usePaginatedSearch";
import { usePermissions } from "@/helpers/usePermissions";

// ─── Khởi tạo store & tiện ích ───────────────────────────────────────────────

const employeeStore = useEmployeeStore();
const departmentStore = useDepartmentStore();
const userStore = useUserStore();
const toast = useToast();

const {
	canViewEmployeeDetail,
	canCreateEmployee,
	canEditEmployee,
	canDeleteEmployee,
	hasAnyEmployeeAction,
} = usePermissions();

const { employees, pagination, loading } = storeToRefs(employeeStore);
const { departments } = storeToRefs(departmentStore);
const { usersWithoutEmp } = storeToRefs(userStore);

// ─── Modal thêm/sửa ──────────────────────────────────────────────────────────

const modalState = useModalState();
const isModalVisible = modalState.isModalVisible;
const isEditMode = modalState.isEditMode;
const openAddModal = modalState.openAddModal;
const openEditModal = modalState.openEditModal;
const closeModal = modalState.closeModal;

// ─── Tìm kiếm & phân trang ───────────────────────────────────────────────────

const paginatedSearch = usePaginatedSearch(
	(params) => employeeStore.fetchEmployees(params),
	pagination,
);
const searchQuery = paginatedSearch.searchQuery;
const loadEmployees = paginatedSearch.load;
const handlePageChange = paginatedSearch.handlePageChange;

// ─── Trạng thái local ────────────────────────────────────────────────────────

const editingEmployee = ref(null);
const formLoading = ref(false);
const relationsLoaded = ref(false);

const isDeleteModalVisible = ref(false);
const deletingEmployee = ref(null);
const deleteMessage = ref("");
const deleteLoading = ref(false);

const isDetailModalVisible = ref(false);
const selectedEmployee = ref(null);

// ─── Xem chi tiết ────────────────────────────────────────────────────────────

const handleViewDetails = (emp) => {
	selectedEmployee.value = emp;
	isDetailModalVisible.value = true;
};

// ─── Tải dữ liệu liên quan ───────────────────────────────────────────────────

const loadFormRelations = async (force = false) => {
	if (!force && relationsLoaded.value) return;

	try {
		await Promise.all([
			departmentStore.fetchDepartments({ limit: 100 }),
			userStore.fetchUsersWithoutEmployee(),
		]);
		relationsLoaded.value = true;
	} catch (err) {
		console.error("Lỗi khi tải dữ liệu liên quan:", err);
	}
};

// ─── Thao tác CRUD ───────────────────────────────────────────────────────────

const handleAdd = async () => {
	editingEmployee.value = null;
	openAddModal();
	await loadFormRelations();
};

const handleEdit = async (emp) => {
	editingEmployee.value = { ...emp };
	openEditModal();
	await loadFormRelations();
};

const handleDelete = (emp) => {
	deletingEmployee.value = emp;
	deleteMessage.value = `Bạn có chắc chắn muốn xoá ${emp.first_name} ${emp.last_name}?`;
	isDeleteModalVisible.value = true;
};

const confirmDelete = async () => {
	const emp = deletingEmployee.value;
	if (!emp) return;

	deleteLoading.value = true;
	const res = await employeeStore.deleteEmployee(emp.id);
	deleteLoading.value = false;

	if (res.success === false) {
		toast.error(res.message);
		return;
	}

	toast.success(res.message);
	isDeleteModalVisible.value = false;
	deletingEmployee.value = null;
	await loadEmployees(pagination.value.page);
};

const handleFormSubmit = async (submittedData) => {
	formLoading.value = true;

	// Chuẩn hoá chung một payload
	const payload = { ...submittedData };
	payload.department_id = payload.department_id
		? Number(payload.department_id)
		: undefined;
	payload.position_id = payload.position_id
		? Number(payload.position_id)
		: undefined;

	// user_id: khi tạo mới dùng null nếu không chọn, khi sửa dùng 0 để backend clear khoá ngoại
	if (payload.user_id) {
		payload.user_id = Number(payload.user_id);
	} else {
		payload.user_id = isEditMode.value ? 0 : null;
	}

	let res;
	if (isEditMode.value) {
		res = await employeeStore.updateEmployee(
			editingEmployee.value.id,
			payload,
		);
	} else {
		res = await employeeStore.createEmployee(payload);
	}

	formLoading.value = false;

	if (!res.success) {
		toast.error(res.message);
		return;
	}

	relationsLoaded.value = false;
	toast.success(
		res.message ||
			(isEditMode.value ? "Cập nhật thành công" : "Thêm mới thành công"),
	);

	closeModal();
	await loadEmployees(pagination.value.page);
};

onMounted(async () => {
	await loadEmployees();
});
</script>

<template>
	<div class="employee-view">
		<header class="page-header">
			<div class="header-content">
				<h1 class="page-title">Quản lý nhân sự</h1>
				<p class="page-subtitle">
					Hệ thống có tổng cộng
					<span>{{ pagination.total }}</span> nhân viên
				</p>
			</div>
			<button
				v-if="canCreateEmployee"
				class="btn btn-primary add-employee-btn"
				@click="handleAdd"
			>
				<Plus class="btn__icon" />
				Thêm nhân viên
			</button>
		</header>

		<main class="content-card">
			<div class="toolbar">
				<div class="search-box">
					<Search class="search-box__icon" />
					<input
						v-model="searchQuery"
						class="form-control search-box__input"
						placeholder="Tìm kiếm nhân viên..."
					/>
				</div>
			</div>

			<div class="table-responsive responsive-table-to-cards">
				<table class="data-table">
					<thead>
						<tr>
							<th>Nhân viên</th>
							<th>Liên hệ</th>
							<th>Phòng ban / Chức vụ</th>
							<th>Ngày vào làm</th>
							<th>Trạng thái</th>
							<th v-if="hasAnyEmployeeAction" class="text-right header-actions-col">
								Thao tác
							</th>
						</tr>
					</thead>
					<tbody>
						<!-- Skeleton placeholder -->
						<template v-if="loading">
							<tr v-for="i in 5" :key="'skeleton-' + i">
								<td data-label="Nhân viên">
									<div class="user-info">
										<Skeleton type="avatar" style="width: 36px; height: 36px; border-radius: 8px;" />
										<div
											class="user-info__details"
											style="
												width: 150px;
												display: flex;
												flex-direction: column;
												gap: 4px;
											"
										>
											<Skeleton type="text" width="80%" />
											<Skeleton type="text" width="50%" />
										</div>
									</div>
								</td>
								<td data-label="Liên hệ">
									<Skeleton type="text" width="100px" />
								</td>
								<td data-label="Phòng ban / Chức vụ">
									<div
										class="job-info"
										style="
											width: 120px;
											display: flex;
											flex-direction: column;
											gap: 4px;
										"
									>
										<Skeleton type="text" width="70%" />
										<Skeleton type="text" width="50%" />
									</div>
								</td>
								<td data-label="Ngày vào làm">
									<Skeleton type="text" width="80px" />
								</td>
								<td data-label="Trạng thái">
									<Skeleton type="badge" />
								</td>
								<td
									v-if="hasAnyEmployeeAction"
									class="text-right"
									data-label="Thao tác"
								>
									<div class="action-group">
										<Skeleton type="btn" style="width: 30px; height: 30px; border-radius: 6px;" />
										<Skeleton type="btn" style="width: 30px; height: 30px; border-radius: 6px;" />
									</div>
								</td>
							</tr>
						</template>

						<!-- Dữ liệu thực -->
						<template v-else>
							<tr v-for="emp in employees" :key="emp.id">
								<td data-label="Nhân viên">
									<div class="user-info">
										<div class="avatar-gradient">
											{{
												getInitials(
													emp.first_name,
													emp.last_name,
												)
											}}
										</div>
										<div class="user-info__details">
											<span class="user-info__name">
												{{ emp.first_name }}
												{{ emp.last_name }}
											</span>
											<span
												v-if="emp.user"
												class="user-info__email"
											>
												{{ emp.user.email }}
											</span>
											<span
												v-else
												class="user-info__email text-unlinked"
											>
												Chưa liên kết tài khoản
											</span>
										</div>
									</div>
								</td>

								<td data-label="Liên hệ">
									<span class="text-main contact-text">
										{{ emp.phone || "—" }}
									</span>
								</td>

								<td data-label="Phòng ban / Chức vụ">
									<div class="job-info">
										<span class="job-info__dept">
											{{ emp.department?.name || "N/A" }}
										</span>
										<span class="job-info__pos">
											{{
												emp.position?.name ||
												"Nhân viên"
											}}
										</span>
									</div>
								</td>

								<td
									data-label="Ngày vào làm"
									class="text-date"
								>
									{{ formatDate(emp.join_date) }}
								</td>

								<td data-label="Trạng thái">
									<span
										:class="[
											'status-badge',
											`status-badge--${emp.status}`,
										]"
									>
										{{ formatStatus(emp.status) }}
									</span>
								</td>

								<td
									v-if="hasAnyEmployeeAction"
									class="text-right"
									data-label="Thao tác"
								>
									<div class="action-group">
										<button
											v-if="canViewEmployeeDetail(emp.id)"
											class="btn-icon btn-icon--detail"
											title="Xem chi tiết"
											@click="handleViewDetails(emp)"
										>
											<Eye />
										</button>
										<button
											v-if="canEditEmployee"
											class="btn-icon btn-icon--edit"
											title="Chỉnh sửa"
											@click="handleEdit(emp)"
										>
											<Pencil />
										</button>
										<button
											v-if="canDeleteEmployee"
											class="btn-icon btn-icon--delete"
											title="Xoá"
											@click="handleDelete(emp)"
										>
											<Trash2 />
										</button>
									</div>
								</td>
							</tr>

							<tr v-if="employees.length === 0">
								<td
									:colspan="hasAnyEmployeeAction ? 6 : 5"
									class="empty-state"
								>
									<div class="empty-state__icon">👥</div>
									<p class="empty-state__text">
										Không có dữ liệu nhân viên nào phù hợp.
									</p>
								</td>
							</tr>
						</template>
					</tbody>
				</table>
			</div>

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

		<!-- Subcomponent Form Modal -->
		<EmployeeModal
			:visible="isModalVisible"
			:is-edit-mode="isEditMode"
			:editing-employee="editingEmployee"
			:departments="departments"
			:users-without-emp="usersWithoutEmp"
			:loading="formLoading"
			@close="closeModal"
			@submit="handleFormSubmit"
		/>

		<!-- Subcomponent Detail Modal -->
		<EmployeeDetailModal
			:visible="isDetailModalVisible"
			:employee="selectedEmployee"
			@close="isDetailModalVisible = false"
		/>

		<!-- Delete Dialog -->
		<ConfirmationDialog
			:visible="isDeleteModalVisible"
			title="Xác nhận xoá nhân viên"
			:message="deleteMessage"
			:loading="deleteLoading"
			@confirm="confirmDelete"
			@cancel="isDeleteModalVisible = false"
		/>
	</div>
</template>

<style scoped>
.employee-view {
	padding-bottom: var(--space-3);
}



.page-subtitle {
	font-size: 0.8rem;
	color: #64748b;
	margin-top: 0.15rem;
}

.add-employee-btn {
	font-size: 0.825rem;
	font-weight: 600;
	height: 36px;
	padding: 0 1rem;
	border-radius: 8px;
	display: inline-flex;
	align-items: center;
	gap: 0.4rem;
	box-shadow: 0 2px 4px rgba(66, 97, 237, 0.15);
}

.content-card {
	background: #ffffff;
	border-radius: 12px;
	border: 1px solid rgba(0, 0, 0, 0.05);
	box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.02), 0 2px 4px -2px rgba(0, 0, 0, 0.02);
	margin-top: 1.5rem;
	overflow: hidden;
}

.toolbar {
	padding: 10px 16px;
	border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.search-box__input {
	height: 34px;
	font-size: 0.825rem;
	padding: 8px 12px 8px 2.25rem !important;
	border-radius: 6px;
}

.search-box__icon {
	width: 14px;
	height: 14px;
	left: 0.85rem;
}

/* Avatar details */
.user-info {
	display: flex;
	align-items: center;
	gap: 0.75rem;
}

.avatar-gradient {
	width: 36px;
	height: 36px;
	border-radius: 8px;
	font-size: 0.775rem;
	font-weight: 700;
	color: #4f46e5;
	background: linear-gradient(135deg, #e0e7ff 0%, #c7d2fe 100%);
	border: 1px solid rgba(79, 70, 229, 0.08);
	display: flex;
	align-items: center;
	justify-content: center;
	flex-shrink: 0;
	box-shadow: inset 0 -1px 0 rgba(0, 0, 0, 0.05);
}

.user-info__details {
	display: flex;
	flex-direction: column;
	gap: 0.1rem;
}

.text-unlinked {
	color: #94a3b8;
	font-style: italic;
}

.contact-text {
	font-size: 0.825rem;
	font-weight: 500;
	color: #334155;
}

.job-info {
	display: flex;
	flex-direction: column;
	gap: 0.1rem;
}

.text-date {
	font-size: 0.8rem;
	color: #64748b;
}

/* Pagination container */
.pagination {
	padding: 12px 16px;
	background: #ffffff;
	border-top: 1px solid rgba(0, 0, 0, 0.05);
	gap: 1rem;
}

.pagination__btn {
	width: 30px;
	height: 30px;
	border-radius: 6px;
	border: 1px solid rgba(0, 0, 0, 0.08);
}

.pagination__btn svg {
	width: 14px;
	height: 14px;
}

.pagination__info {
	font-size: 0.775rem;
	color: #64748b;
}

.header-actions-col {
	width: 120px;
}

@media (max-width: 768px) {
	.responsive-table-to-cards .data-table tr {
		border-bottom: 6px solid #f1f5f9;
		padding: 12px 16px;
		border-radius: 8px;
		margin-bottom: 8px;
		border-top: 1px solid rgba(0, 0, 0, 0.05);
		border-left: 1px solid rgba(0, 0, 0, 0.05);
		border-right: 1px solid rgba(0, 0, 0, 0.05);
	}
	
	.responsive-table-to-cards .data-table td {
		padding: 6px 0;
	}
}
</style>
