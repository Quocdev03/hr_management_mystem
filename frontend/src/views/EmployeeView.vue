<script setup>
import { ref, onMounted } from 'vue';
import { storeToRefs } from 'pinia';
import { useToast } from 'vue-toastification';

import { useEmployeeStore } from '@/store/employee';
import { useDepartmentStore } from '@/store/department';
import { useUserStore } from '@/store/user';

import ModalDialog from '@/components/ModalDialog.vue';
import ConfirmationDialog from '@/components/ConfirmationDialog.vue';
import Skeleton from '@/components/Skeleton.vue';

import { getInitials, formatDate, formatStatus } from '@/helpers/formatters';
import { useModalState } from '@/helpers/useModalState';
import { usePaginatedSearch } from '@/helpers/usePaginatedSearch';
import { usePermissions } from '@/helpers/usePermissions';

import plusIcon from '@/assets/svg/plus.svg';
import searchIcon from '@/assets/svg/search.svg';
import editIcon from '@/assets/svg/edit.svg';
import deleteIcon from '@/assets/svg/delete.svg';
import prevIcon from '@/assets/svg/chevron-left.svg';
import nextIcon from '@/assets/svg/chevron-right.svg';

// Stores
const employeeStore = useEmployeeStore();
const departmentStore = useDepartmentStore();
const userStore = useUserStore();
const toast = useToast();
const {
  canCreateEmployee,
  canEditEmployee,
  canDeleteEmployee,
  hasAnyEmployeeAction,
} = usePermissions();

const { employees, pagination, loading } = storeToRefs(employeeStore);
const { departments } = storeToRefs(departmentStore);
const { usersWithoutEmp } = storeToRefs(userStore);

// Modal state
const modalState = useModalState();
const isModalVisible = modalState.isModalVisible;
const isEditMode = modalState.isEditMode;
const openAddModal = modalState.openAddModal;
const openEditModal = modalState.openEditModal;
const closeModal = modalState.closeModal;

// Search & pagination
const paginatedSearch = usePaginatedSearch(
  (params) => employeeStore.fetchEmployees(params),
  pagination,
);
const searchQuery = paginatedSearch.searchQuery;
const loadEmployees = paginatedSearch.load;
const handlePageChange = paginatedSearch.handlePageChange;

// Local state
const editingEmployee = ref(null);
const formLoading = ref(false);
const formData = ref(buildFormData());
const relationsLoaded = ref(false);

// Delete modal state
const isDeleteModalVisible = ref(false);
const deletingEmployee = ref(null);
const deleteMessage = ref('');
const deleteLoading = ref(false);

/** Tạo object formData từ data ban đầu, đảm bảo luôn có đủ các field */
function buildFormData(data = {}) {
  const d = data ?? {};
  return {
    user_id: d.user_id ?? null,
    first_name: d.first_name ?? '',
    last_name: d.last_name ?? '',
    phone: d.phone ?? '',
    department_id: d.department_id ?? '',
    position: d.position ?? '',
    salary: d.salary ?? null,
    join_date: d.join_date ?? '',
    status: d.status ?? 'active',
    gender: d.gender ?? '',
    birth_date: d.birth_date ?? '',
  };
}

async function loadFormRelations(force = false) {
  if (!force && relationsLoaded.value) {
    return;
  }

  try {
    await Promise.all([
      departmentStore.fetchDepartments(),
      userStore.fetchUsersWithoutEmployee(),
    ]);
    relationsLoaded.value = true;
  } catch (err) {
    console.error('Lỗi khi tải dữ liệu liên quan:', err);
  }
}

// Handlers
async function handleAdd() {
  editingEmployee.value = null;
  formData.value = buildFormData();
  openAddModal();
  await loadFormRelations();
}

async function handleEdit(emp) {
  const empCopy = { ...emp };

  if (empCopy.birth_date) {
    const dateParts = empCopy.birth_date.split('T');
    empCopy.birth_date = dateParts[0];
  } else {
    empCopy.birth_date = '';
  }

  if (empCopy.join_date) {
    const dateParts = empCopy.join_date.split('T');
    empCopy.join_date = dateParts[0];
  } else {
    empCopy.join_date = '';
  }

  editingEmployee.value = empCopy;
  formData.value = buildFormData(empCopy);
  openEditModal();
  await loadFormRelations();
}

function handleDelete(emp) {
  deletingEmployee.value = emp;
  deleteMessage.value = `Bạn có chắc chắn muốn xoá ${emp.first_name} ${emp.last_name}?`;
  isDeleteModalVisible.value = true;
}

async function confirmDelete() {
  const emp = deletingEmployee.value;
  if (!emp) {
    return;
  }

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
}

async function handleFormSubmit() {
  formLoading.value = true;
  let res;

  if (isEditMode.value === true) {
    // Edit mode: build a normalized dirty-check payload to avoid false diffs
    const original = buildFormData(editingEmployee.value);

    function normalizeForCompare(obj) {
      const out = {};
      Object.keys(obj).forEach((k) => {
        const v = obj[k];
        if (v === null || typeof v === 'undefined') {
          out[k] = '';
          return;
        }
        // Normalize primary numeric IDs and salary to string form
        if (k === 'user_id' || k === 'department_id' || k === 'salary') {
          out[k] = String(v);
          return;
        }
        // Dates and other primitives -> string
        out[k] = typeof v === 'object' ? JSON.stringify(v) : String(v);
      });
      return out;
    }

    const normOriginal = normalizeForCompare(original);
    const normForm = normalizeForCompare(formData.value);

    const payload = Object.fromEntries(
      Object.keys(normOriginal)
        .filter((k) => normForm[k] !== normOriginal[k])
        .map((k) => [k, formData.value[k]]),
    );

    if (Object.keys(payload).length === 0) {
      toast.info('Không có dữ liệu thay đổi');
      formLoading.value = false;
      return;
    }

    // Sanitize values for Go backend
    if ('user_id' in payload) {
      payload.user_id =
        formData.value.user_id == null || formData.value.user_id === ''
          ? 0
          : Number(formData.value.user_id);
    }
    if ('department_id' in payload) {
      payload.department_id = Number(payload.department_id);
    }

    res = await employeeStore.updateEmployee(
      editingEmployee.value.id,
      payload,
    );
  } else {
    // Create mode
    const data = { ...formData.value };
    data.user_id =
      data.user_id !== '' && data.user_id != null
        ? Number(data.user_id)
        : null;
    data.department_id =
      data.department_id !== '' ? Number(data.department_id) : undefined;

    res = await employeeStore.createEmployee(data);
  }

  formLoading.value = false;

  if (res.success === false) {
    toast.error(res.message);
    return;
  }

  relationsLoaded.value = false;

  if (isEditMode.value === true) {
    toast.success(res.message);
  }

  closeModal();
  await loadEmployees(pagination.value.page);
}

onMounted(async () => {
  await loadEmployees();
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
			<button
				v-if="canCreateEmployee"
				class="btn btn--primary"
				@click="handleAdd"
			>
				<img :src="plusIcon" alt="add" class="btn__icon" />
				Thêm nhân viên
			</button>
		</header>

		<!-- Nội dung chính -->
		<main class="content-card">
			<!-- Thanh công cụ -->
			<div class="toolbar">
				<div class="search-box">
					<img :src="searchIcon" class="search-box__icon" alt="search" />
					<input
						v-model="searchQuery"
						class="form-control search-box__input"
						placeholder="Tìm tên"
					/>
				</div>
			</div>

			<!-- Bảng dữ liệu -->
			<div class="table-responsive">
				<table class="data-table">
					<thead>
						<tr>
							<th>Nhân viên</th>
							<th>Liên hệ</th>
							<th>Phòng ban / Chức vụ</th>
							<th>Ngày vào làm</th>
							<th>Trạng thái</th>
							<th v-if="hasAnyEmployeeAction" class="text-right">
								Thao tác
							</th>
						</tr>
					</thead>
					<tbody>
						<!-- Loading skeleton rows -->
						<template v-if="loading">
							<tr v-for="i in 5" :key="'skeleton-' + i">
								<td>
									<div class="user-info">
										<Skeleton type="avatar" />
										<div
											class="user-info__details"
											style="
												width: 150px;
												display: flex;
												flex-direction: column;
												gap: var(--space-1);
											"
										>
											<Skeleton type="text" width="80%" />
											<Skeleton type="text" width="50%" />
										</div>
									</div>
								</td>
								<td>
									<Skeleton type="text" width="100px" />
								</td>
								<td>
									<div
										class="job-info"
										style="
											width: 120px;
											display: flex;
											flex-direction: column;
											gap: var(--space-1);
										"
									>
										<Skeleton type="text" width="70%" />
										<Skeleton type="text" width="50%" />
									</div>
								</td>
								<td>
									<Skeleton type="text" width="80px" />
								</td>
								<td>
									<Skeleton type="badge" />
								</td>
								<td v-if="hasAnyEmployeeAction" class="text-right">
									<div class="action-group">
										<Skeleton type="btn" />
										<Skeleton type="btn" />
									</div>
								</td>
							</tr>
						</template>

						<!-- Actual rows when loaded -->
						<template v-else>
							<tr v-for="emp in employees" :key="emp.id">
								<td>
									<div class="user-info">
										<div class="user-info__avatar">
											{{
												getInitials(emp.first_name, emp.last_name)
											}}
										</div>
										<div class="user-info__details">
											<h1 class="user-info__name">
												{{ emp.first_name }}
												{{ emp.last_name }}
											</h1>
											<span class="user-info__email" v-if="emp.user">
												{{ emp.user.email }}
											</span>
											<span class="user-info__email" v-if="emp.user">
												User: {{ emp.user.user_name }}
											</span>
											<span v-else class="user-info__email">
												Chưa có tài khoản
											</span>
										</div>
									</div>
								</td>
								<td>
									<span class="text-main fw-500">
										{{ emp.phone || "—" }}
									</span>
								</td>
								<td>
									<div class="job-info">
										<h1 class="job-info__dept">
											{{ emp.department?.name || "N/A" }}
										</h1>
										<span class="job-info__pos">
											{{ emp.position || "Nhân viên" }}
										</span>
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
								<td v-if="hasAnyEmployeeAction" class="text-right">
									<div class="action-group">
										<button
											v-if="canEditEmployee"
											class="btn-icon btn-icon--edit"
											title="Chỉnh sửa"
											@click="handleEdit(emp)"
										>
											<img :src="editIcon" alt="edit" />
										</button>
										<button
											v-if="canDeleteEmployee"
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
		<ModalDialog
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
			<form @submit.prevent="handleFormSubmit" class="employee-form">
				<div class="form-grid">
					<!-- Họ và tên -->
					<div class="form-group">
						<label class="form-label"
							>Họ <span class="required">*</span></label
						>
						<input
							v-model="formData.first_name"
							type="text"
							class="form-control"
							placeholder="Nhập họ..."
							required
						/>
					</div>
					<div class="form-group">
						<label class="form-label"
							>Tên <span class="required">*</span></label
						>
						<input
							v-model="formData.last_name"
							type="text"
							class="form-control"
							placeholder="Nhập tên..."
							required
						/>
					</div>

					<!-- Giới tính -->
					<div class="form-group">
						<label class="form-label">Giới Tính</label>
						<select v-model="formData.gender" class="form-control">
							<option value="" disabled>Chọn giới tính</option>
							<option value="male">Nam</option>
							<option value="female">Nữ</option>
						</select>
					</div>

					<!-- Ngày tháng năm sinh -->
					<div class="form-group">
						<label class="form-label">Ngày sinh</label>
						<input
							v-model="formData.birth_date"
							type="date"
							class="form-control"
						/>
					</div>

					<div class="form-group">
						<label class="form-label"
							>Số điện thoại <span class="required">*</span></label
						>
						<input
							v-model="formData.phone"
							type="text"
							class="form-control"
							placeholder="0xxxxxxxxx"
							required
						/>
					</div>

					<!-- Công việc -->
					<div class="form-group">
						<label class="form-label"
							>Phòng ban <span class="required">*</span></label
						>
						<select
							v-model="formData.department_id"
							class="form-control"
							required
						>
							<option value="" disabled>Chọn phòng ban</option>
							<option
								v-for="dept in departments"
								:key="dept.id"
								:value="dept.id"
							>
								{{ dept.name }}
							</option>
						</select>
					</div>

					<div class="form-group">
						<label class="form-label">Chức vụ</label>
						<input
							v-model="formData.position"
							type="text"
							class="form-control"
							placeholder="Ví dụ: Backend Developer"
						/>
					</div>

					<!-- Lương và Ngày vào làm -->
					<div class="form-group">
						<label class="form-label">Mức lương (VNĐ)</label>
						<input
							v-model.number="formData.salary"
							type="number"
							class="form-control"
							placeholder="0"
						/>
					</div>
					<div
						class="form-group"
						:class="{ 'form-group--disabled': isEditMode }"
					>
						<label class="form-label">Ngày vào làm</label>
						<input
							v-model="formData.join_date"
							type="date"
							class="form-control"
							:disabled="isEditMode"
						/>
					</div>
					<!-- Trạng thái và Tài khoản -->
					<div class="form-group">
						<label class="form-label">Trạng thái</label>
						<select v-model="formData.status" class="form-control">
							<option value="active">Đang làm việc</option>
							<option value="inactive">Đã nghỉ việc</option>
						</select>
					</div>
				</div>

				<!-- Gắn người dùng -->
				<div class="form-group">
					<span class="form-label">Liên kết người dùng</span>
					<select v-model="formData.user_id" class="form-control">
						<option :value="null">Không liên kết</option>

						<!-- user hiện tại (khi edit) -->
						<option
							v-if="editingEmployee?.user"
							:value="editingEmployee.user.id"
						>
							{{ editingEmployee.user.email }} (hiện tại)
						</option>

						<!-- danh sách user chưa có employee -->
						<option
							v-for="u in usersWithoutEmp"
							:key="u.id"
							:value="u.id"
						>
							{{ u.email }}
						</option>
					</select>
				</div>

				<div class="form-actions">
					<button
						type="button"
						class="btn btn--secondary"
						@click="closeModal"
					>
						Hủy bỏ
					</button>
					<button
						type="submit"
						class="btn btn--primary"
						:disabled="formLoading"
					>
						<span v-if="formLoading" class="spinner"></span>
						{{ isEditMode ? "Lưu thay đổi" : "Thêm nhân viên" }}
					</button>
				</div>
			</form>
		</ModalDialog>

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
	padding: var(--space-3) var(--space-2);
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
	display: -webkit-box;
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

/* ===== Form ===== */
.employee-form {
	display: flex;
	flex-direction: column;
	gap: var(--space-4);
}

.form-grid {
	display: grid;
	grid-template-columns: repeat(2, 1fr);
	gap: var(--space-3);
}

.form-group {
	display: flex;
	flex-direction: column;
	gap: 0.5rem;
}

.form-actions {
	display: flex;
	justify-content: flex-end;
	gap: var(--space-2);
	margin-top: var(--space-2);
	padding-top: var(--space-3);
	border-top: 1px solid var(--border-color);
}

/* Đồng bộ font chữ với main.css */
.form-label {
	font-size: var(--fs-sm);
	font-weight: var(--fw-semibold);
}

.form-group--disabled {
	opacity: 0.8;
}

.form-group--disabled .form-control {
	background-color: var(--bg-light);
	cursor: not-allowed;
	border-color: var(--border-color);
}

.form-hint {
	font-size: 11px;
	color: var(--text-light);
	margin-top: -2px;
}

.required {
	color: var(--danger-color);
}

.spinner {
	width: 16px;
	height: 16px;
	border: 2px solid rgba(255, 255, 255, 0.3);
	border-radius: 50%;
	border-top-color: #fff;
	animation: spin 0.8s linear infinite;
	display: inline-block;
	margin-right: 8px;
}

@keyframes spin {
	to {
		transform: rotate(360deg);
	}
}

@media (max-width: 640px) {
	.form-grid {
		grid-template-columns: 1fr;
	}
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
