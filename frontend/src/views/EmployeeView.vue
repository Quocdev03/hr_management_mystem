<script setup>
import { ref, onMounted } from "vue";
import { storeToRefs } from "pinia";
import { useToast } from "vue-toastification";

import { useEmployeeStore } from "@/store/employee";
import { useDepartmentStore } from "@/store/department";

import BaseModal from "@/components/BaseModal.vue";
import EmployeeForm from "@/components/EmployeeForm.vue";
import ConfirmModal from "@/components/ConfirmModal.vue";

import { getInitials, formatDate, formatStatus } from "@/helpers/formatters";
import { useModalState } from "@/helpers/useModalState";
import { usePaginatedSearch } from "@/helpers/usePaginatedSearch";

import plusIcon from "@/assets/svg/plus.svg";
import searchIcon from "@/assets/svg/search.svg";
import editIcon from "@/assets/svg/edit.svg";
import deleteIcon from "@/assets/svg/delete.svg";
import prevIcon from "@/assets/svg/chevron-left.svg";
import nextIcon from "@/assets/svg/chevron-right.svg";

// Stores
const employeeStore = useEmployeeStore();
const departmentStore = useDepartmentStore();
const toast = useToast();

const { employees, pagination, loading } = storeToRefs(employeeStore);

// Modal state
const modalState = useModalState();
const isModalVisible = modalState.isModalVisible;
const isEditMode = modalState.isEditMode;
const openAddModal = modalState.openAddModal;
const openEditModal = modalState.openEditModal;
const closeModal = modalState.closeModal;

// Search & pagination
const paginatedSearch = usePaginatedSearch(function (params) {
	return employeeStore.fetchEmployees(params);
}, pagination);

const searchQuery = paginatedSearch.searchQuery;
const loadEmployees = paginatedSearch.load;
const handlePageChange = paginatedSearch.handlePageChange;

// Local state
const editingEmployee = ref(null);
const formLoading = ref(false);
const modalKey = ref(0);

// Delete modal state
const isDeleteModalVisible = ref(false);
const deletingEmployee = ref(null);
const deleteMessage = ref("");
const deleteLoading = ref(false);

// Handlers
function handleAdd() {
	editingEmployee.value = null;
	modalKey.value++;
	openAddModal();
}

function handleEdit(emp) {
	let empCopy = Object.assign({}, emp);

	if (empCopy.join_date) {
		let dateParts = empCopy.join_date.split("T");
		empCopy.join_date = dateParts[0];
	} else {
		empCopy.join_date = "";
	}

	editingEmployee.value = empCopy;
	modalKey.value++;
	openEditModal();
}

function handleDelete(emp) {
	deletingEmployee.value = emp;
	deleteMessage.value =
		"Bạn có chắc chắn muốn xoá " +
		emp.first_name +
		" " +
		emp.last_name +
		"?";
	isDeleteModalVisible.value = true;
}

async function confirmDelete() {
	let emp = deletingEmployee.value;
	if (!emp) {
		return;
	}

	deleteLoading.value = true;
	const res = await employeeStore.deleteEmployee(emp.id);
	deleteLoading.value = false;

	if (res.success === false) {
		if (res.message) {
			toast.error(res.message);
		} else {
			toast.error("Xoá thất bại");
		}
		return;
	}

	toast.success("Xoá nhân viên thành công");
	isDeleteModalVisible.value = false;
	deletingEmployee.value = null;
	await loadEmployees(pagination.value.page);
}

async function onFormSubmit(formData) {
	formLoading.value = true;

	let res;
	if (isEditMode.value === true) {
		res = await employeeStore.updateEmployee(
			editingEmployee.value.id,
			formData,
		);
	} else {
		res = await employeeStore.createEmployee(formData);
	}

	formLoading.value = false;

	if (res.success === false) {
		if (res.message) {
			toast.error(res.message);
		} else {
			toast.error("Có lỗi xảy ra");
		}
		return;
	}

	if (isEditMode.value === true) {
		toast.success("Cập nhật thành công");
	} else {
		toast.success("Thêm nhân viên thành công");
	}

	closeModal();
	await loadEmployees(pagination.value.page);
}

onMounted(async function () {
	await loadEmployees();
	await departmentStore.fetchDepartments();
});
</script>

<template>
	<div class="employee-view">
		<!-- Tiêu đề trang -->
		<header class="page-header">
			<div class="header-content">
				<h1 class="page-title">Quản lý nhân sự</h1>
				<p class="page-subtitle">
					Hệ thống có tổng cộng
					<span>{{ pagination.total }}</span> nhân viên
				</p>
			</div>
			<button class="btn btn--primary" @click="handleAdd">
				<img :src="plusIcon" alt="add" class="btn__icon" />
				Thêm nhân viên
			</button>
		</header>

		<!-- Nội dung chính -->
		<main class="content-card">
			<!-- Thanh công cụ -->
			<div class="toolbar">
				<div class="search-box">
					<img
						:src="searchIcon"
						class="search-box__icon"
						alt="search"
					/>
					<input
						v-model="searchQuery"
						class="form-control search-box__input"
						placeholder="Tìm tên, email..."
					/>
				</div>
			</div>

			<div v-if="loading" class="table-loading">Đang tải dữ liệu...</div>

			<!-- Bảng dữ liệu -->
			<div v-else class="table-responsive">
				<table class="data-table">
					<thead>
						<tr>
							<th>Nhân viên</th>
							<th>Liên hệ</th>
							<th>Phòng ban / Chức vụ</th>
							<th>Ngày vào làm</th>
							<th>Trạng thái</th>
							<th class="text-right">Thao tác</th>
						</tr>
					</thead>
					<tbody>
						<tr v-for="emp in employees" :key="emp.id">
							<td>
								<div class="user-info">
									<div class="user-info__avatar">
										{{
											getInitials(
												emp.first_name,
												emp.last_name,
											)
										}}
									</div>
									<div class="user-info__details">
										<div class="user-info__name">
											{{ emp.first_name }}
											{{ emp.last_name }}
										</div>
										<div class="user-info__email">
											{{ emp.email }}
										</div>
									</div>
								</div>
							</td>
							<td>
								<div class="text-main fw-500">
									{{ emp.phone || "—" }}
								</div>
							</td>
							<td>
								<div class="job-info">
									<div class="job-info__dept">
										{{ emp.department?.name || "N/A" }}
									</div>
									<div class="job-info__pos">
										{{ emp.position || "Nhân viên" }}
									</div>
								</div>
							</td>
							<td class="text-muted">
								{{ formatDate(emp.join_date) }}
							</td>
							<td>
								<span
									:class="[
										'status-badge',
										`status-badge--${emp.status}`,
									]"
								>
									{{ formatStatus(emp.status) }}
								</span>
							</td>
							<td class="text-right">
								<div class="action-group">
									<button
										class="btn-icon btn-icon--edit"
										title="Chỉnh sửa"
										@click="handleEdit(emp)"
									>
										<img :src="editIcon" alt="edit" />
									</button>
									<button
										class="btn-icon btn-icon--delete"
										title="Xoá"
										@click="handleDelete(emp)"
									>
										<img :src="deleteIcon" alt="delete" />
									</button>
								</div>
							</td>
						</tr>
						<tr v-if="employees.length === 0">
							<td colspan="6" class="empty-state">
								<div class="empty-state__icon">👥</div>
								<p class="empty-state__text">
									Không có dữ liệu nhân viên nào phù hợp.
								</p>
							</td>
						</tr>
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
		</main>

		<!-- Modal Thêm/Sửa -->
		<BaseModal
			:visible="isModalVisible"
			:title="isEditMode ? 'Chỉnh sửa nhân viên' : 'Thêm nhân viên mới'"
			:subtitle="
				isEditMode
					? 'Cập nhật thông tin chi tiết của nhân viên'
					: 'Điền thông tin để tạo nhân viên mới vào hệ thống'
			"
			size="lg"
			@close="closeModal"
		>
			<EmployeeForm
				:key="modalKey"
				:initial-data="editingEmployee"
				:is-edit="isEditMode"
				:loading="formLoading"
				@submit="onFormSubmit"
				@cancel="closeModal"
			/>
		</BaseModal>

		<ConfirmModal
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
	padding-bottom: var(--space-4);
}

/* Header */
.page-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: var(--space-4);
	flex-wrap: wrap;
	gap: var(--space-2);
}

.page-title {
	font-size: var(--fs-2xl);
	font-weight: var(--fw-bold);
	letter-spacing: var(--tracking-tight);
	margin: 0 0 4px 0;
	color: var(--text-main);
}

.page-subtitle {
	color: var(--text-muted);
	font-size: var(--fs-sm);
	margin: 0;
}

.page-subtitle span {
	font-weight: var(--fw-semibold);
	color: var(--primary-color);
}

/* Card & Toolbar */
.content-card {
	background: var(--bg-card);
	border-radius: var(--radius-lg);
	border: 1px solid var(--border-color);
	box-shadow: var(--shadow-sm);
	overflow: hidden;
}

.toolbar {
	display: flex;
	justify-content: space-between;
	gap: var(--space-3);
	padding: var(--space-3) var(--space-4);
	border-bottom: 1px solid var(--border-color);
	flex-wrap: wrap;
	background: var(--bg-card);
	align-items: center;
}

.search-box {
	position: relative;
	flex: 1;
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

/* Table */
.table-responsive {
	overflow-x: auto;
}

.data-table {
	width: 100%;
	border-collapse: collapse;
	text-align: left;
}

.data-table th {
	padding: var(--space-2) var(--space-4);
	background: var(--bg-lighter);
	font-size: var(--fs-xs);
	text-transform: uppercase;
	font-weight: var(--fw-bold);
	color: var(--text-muted);
	letter-spacing: 0.05em;
	border-bottom: 1px solid var(--border-color);
}

.data-table td {
	padding: var(--space-3) var(--space-4);
	border-bottom: 1px solid var(--border-color);
	vertical-align: middle;
	background: var(--bg-card);
}

.data-table tbody tr:hover td {
	background: var(--bg-lighter);
}

/* User Info Cell */
.user-info {
	display: flex;
	align-items: center;
	gap: var(--space-2);
}

.user-info__avatar {
	width: 40px;
	height: 40px;
	border-radius: var(--radius-md);
	background: #eff6ff;
	color: var(--primary-color);
	display: flex;
	align-items: center;
	justify-content: center;
	font-weight: var(--fw-bold);
	font-size: var(--fs-xs);
	border: 1px solid rgba(59, 130, 246, 0.1);
	flex-shrink: 0;
}

.user-info__name {
	font-weight: var(--fw-semibold);
	color: var(--text-main);
	font-size: var(--fs-sm);
}

.user-info__email {
	font-size: var(--fs-xs);
	color: var(--text-muted);
}

/* Job Info Cell */
.job-info {
	display: flex;
	flex-direction: column;
	gap: 2px;
}

.job-info__dept {
	font-weight: var(--fw-semibold);
	font-size: var(--fs-sm);
	color: var(--text-main);
}

.job-info__pos {
	font-size: var(--fs-xs);
	color: var(--text-muted);
}

/* Status Badge */
.status-badge {
	padding: 0.25rem 0.75rem;
	border-radius: var(--radius-sm);
	font-size: var(--fs-xs);
	font-weight: var(--fw-bold);
	display: inline-block;
}

.status-badge--active {
	background: #dcfce7;
	color: #15803d;
}

.status-badge--inactive {
	background: #fee2e2;
	color: #b91c1c;
}

/* Action buttons */
.action-group {
	display: flex;
	justify-content: flex-end;
	gap: 8px;
}

.btn-icon {
	width: 32px;
	height: 32px;
	border-radius: var(--radius-sm);
	border: 1px solid var(--border-color);
	background: var(--bg-card);
	display: flex;
	align-items: center;
	justify-content: center;
	cursor: pointer;
	transition: all 0.2s;
}

.btn-icon img {
	width: 16px;
	height: 16px;
	filter: grayscale(1) opacity(0.6);
}

.btn-icon:hover {
	border-color: var(--primary-color);
	background: #eff6ff;
}

.btn-icon:hover img {
	filter: none;
}

.btn-icon--delete:hover {
	border-color: var(--danger-color);
	background: #fee2e2;
}

/* Pagination */
.pagination {
	padding: var(--space-3) var(--space-4);
	display: flex;
	align-items: center;
	justify-content: center;
	gap: var(--space-3);
	background: var(--bg-card);
	border-top: 1px solid var(--border-color);
}

.pagination__btn {
	width: 36px;
	height: 36px;
	border-radius: var(--radius-md);
	border: 1px solid var(--border-color);
	background: var(--bg-main);
	display: flex;
	align-items: center;
	justify-content: center;
	cursor: pointer;
	transition: all 0.2s;
}

.pagination__btn img {
	width: 18px;
	height: 18px;
	opacity: 0.6;
}

.pagination__btn:hover:not(:disabled) {
	border-color: var(--primary-color);
	background: var(--bg-light);
}

.pagination__btn:disabled {
	opacity: 0.3;
	cursor: not-allowed;
}

.pagination__info {
	font-size: var(--fs-sm);
	color: var(--text-muted);
}

.pagination__info span {
	font-weight: var(--fw-bold);
	color: var(--text-main);
}

/* Empty state */
.empty-state {
	padding: var(--space-8) 0;
	text-align: center;
	color: var(--text-light);
}

.empty-state__icon {
	font-size: 3rem;
	margin-bottom: var(--space-2);
}

/* Misc */
.btn__icon {
	width: 18px;
	height: 18px;
	filter: brightness(0) invert(1);
}

@media (max-width: 768px) {
	.toolbar {
		flex-direction: column;
		align-items: stretch;
	}
	.search-box {
		max-width: none;
	}
}
</style>
