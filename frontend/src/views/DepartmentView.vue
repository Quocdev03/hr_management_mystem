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

// ─── Component UI dùng chung ─────────────────────────────────────────────────
import ModalDialog from "@/components/ModalDialog.vue";
import ConfirmationDialog from "@/components/ConfirmationDialog.vue";
import Skeleton from "@/components/Skeleton.vue";

// ─── Tiện ích ────────────────────────────────────────────────────────────────
import { ref, onMounted, computed } from "vue";
import { storeToRefs } from "pinia";
import { useToast } from "vue-toastification";
import { useModalState } from "@/helpers/useModalState";
import { usePaginatedSearch } from "@/helpers/usePaginatedSearch";
import { usePermissions } from "@/helpers/usePermissions";

// ─── Khởi tạo ────────────────────────────────────────────────────────────────

const departmentStore = useDepartmentStore();
const dashboardStore = useDashboardStore();
const employeeStore = useEmployeeStore();
const toast = useToast();

// Quyền thao tác CRUD phòng ban
const { canCrudDepartment } = usePermissions();

// Reactive refs từ store
const { departments, loading, pagination } = storeToRefs(departmentStore);
const employeeOptions = ref([]);

// ─── Modal thêm/sửa ──────────────────────────────────────────────────────────

const { isModalVisible, isEditMode, openAddModal, openEditModal, closeModal } =
	useModalState();

// ─── Tìm kiếm & phân trang ───────────────────────────────────────────────────

const {
	searchQuery,
	load: loadDepartments, // Tải danh sách phòng ban
	handlePageChange, // Xử lý chuyển trang
} = usePaginatedSearch(
	(params) => departmentStore.fetchDepartments(params),
	pagination,
);

// ─── Trạng thái local ────────────────────────────────────────────────────────

const editingDepartment = ref(null); // Phòng ban đang được chỉnh sửa
const formLoading = ref(false); // Đang xử lý submit form

// Giá trị mặc định form — dùng hàm để tránh tái sử dụng cùng một object reference
const initialFormData = () => ({
	name: "",
	code: "",
	description: "",
	manager_id: null,
});
const formData = ref(initialFormData());

// ─── Trạng thái modal xoá ────────────────────────────────────────────────────

const isDeleteModalVisible = ref(false); // Modal xoá có đang mở không
const deletingDepartment = ref(null); // Phòng ban sắp bị xoá
const deleteMessage = ref(""); // Nội dung xác nhận xoá
const deleteLoading = ref(false); // Đang xử lý xoá

// ─── Tải danh sách nhân viên ─────────────────────────────────────────────────

// Lấy toàn bộ nhân viên (tối đa 100) để đổ vào dropdown chọn trưởng phòng
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

// Reset form về trạng thái rỗng rồi mở modal
function handleAdd() {
	editingDepartment.value = null;
	formData.value = initialFormData();
	openAddModal();
}

// ─── Mở modal sửa ────────────────────────────────────────────────────────────

/**
 * Điền dữ liệu phòng ban vào form.
 * manager_id ưu tiên lấy từ field trực tiếp, fallback về manager?.id
 * (API có thể trả về dạng nested object thay vì flat id).
 */
async function handleEdit(department) {
	editingDepartment.value = { ...department };
	formData.value = {
		name: department.name ?? "",
		code: department.code ?? "",
		description: department.description ?? "",
		manager_id: department.manager_id ?? department.manager?.id ?? null,
	};
	openEditModal();
	await loadEmployees(); // Tải nhân viên để dropdown có dữ liệu
}

// ─── Mở modal xoá ────────────────────────────────────────────────────────────

// Gán phòng ban cần xoá và hiển thị modal xác nhận
function handleDelete(department) {
	deletingDepartment.value = department;
	deleteMessage.value = `Bạn có chắc chắn muốn xoá phòng ban ${department.name}?`;
	isDeleteModalVisible.value = true;
}

// ─── Xác nhận xoá ────────────────────────────────────────────────────────────

// Thực hiện xoá sau khi người dùng xác nhận
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
		deleteLoading.value = false; // Luôn tắt loading dù thành công hay lỗi
	}
}

// ─── Submit form thêm/sửa ─────────────────────────────────────────────────────

async function handleFormSubmit() {
	formLoading.value = true;
	try {
		let res;

		if (isEditMode.value) {
			// ── Chế độ SỬA: partial update — chỉ gửi field thực sự thay đổi ──

			const original = editingDepartment.value ?? {};
			const payload = {};

			if (formData.value.name !== (original.name ?? "")) {
				payload.name = formData.value.name;
			}

			if (formData.value.description !== (original.description ?? "")) {
				payload.description = formData.value.description;
			}

			/**
			 * So sánh manager_id cẩn thận vì API có thể trả dạng flat (manager_id)
			 * hoặc nested (manager.id). Chuẩn hoá cả hai về null trước khi so sánh.
			 *
			 * Khi xoá manager: gửi 0 thay vì null vì backend dùng 0 làm sentinel
			 * để phân biệt "không truyền field" vs "chủ động xoá trưởng phòng".
			 */
			const originalManagerId =
				original.manager_id ?? original.manager?.id ?? null;
			const currentManagerId = formData.value.manager_id || null;

			if (currentManagerId !== originalManagerId) {
				payload.manager_id =
					currentManagerId === null ? 0 : currentManagerId;
			}

			// Không có field nào thay đổi → báo và thoát sớm
			if (Object.keys(payload).length === 0) {
				toast.info("Không có dữ liệu thay đổi");
				return;
			}

			res = await departmentStore.updateDepartment(
				editingDepartment.value.id,
				payload,
			);
		} else {
			// ── Chế độ THÊM MỚI: gửi đủ các field ──
			res = await departmentStore.createDepartment({
				name: formData.value.name,
				code: formData.value.code,
				description: formData.value.description,
				manager_id: formData.value.manager_id || null,
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
		formLoading.value = false; // Luôn tắt loading dù thành công hay lỗi
	}
}

// ─── Khởi tạo trang ──────────────────────────────────────────────────────────

// Tải song song phòng ban + dashboard để giảm thời gian chờ
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
					<span>{{ departments.length }}</span> phòng ban
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

		<main class="content-card">
			<div class="toolbar">
				<div class="search-box">
					<img :src="searchIcon" class="search-box__icon" alt="search" />
					<input
						v-model="searchQuery"
						class="form-control search-box__input"
						placeholder="Tìm tên hoặc mã phòng ban..."
					/>
				</div>
			</div>

			<div class="table-responsive">
				<table class="data-table">
					<thead>
						<tr>
							<th>Tên phòng ban</th>
							<th>Mã</th>
							<th>Mô tả</th>
							<th class="text-center">Trưởng Phòng</th>
							<th v-if="canCrudDepartment" class="text-right">
								Thao tác
							</th>
						</tr>
					</thead>
					<tbody>
						<!-- Loading skeleton rows -->
						<template v-if="loading">
							<tr v-for="i in 5" :key="'skeleton-' + i">
								<td class="text-main fw-500">
									<Skeleton type="text" width="150px" height="18px" />
								</td>
								<td>
									<Skeleton
										type="text"
										class="dept-code"
										width="60px"
										height="22px"
									/>
								</td>
								<td class="text-muted">
									<Skeleton type="text" width="220px" height="16px" />
								</td>
								<td class="text-center">
									<Skeleton
										type="text"
										class="employee-count"
										width="120px"
										height="18px"
										style="display: inline-block"
									/>
								</td>
								<td v-if="canCrudDepartment" class="text-right">
									<div class="action-group">
										<Skeleton type="btn" />
										<Skeleton type="btn" />
									</div>
								</td>
							</tr>
						</template>

						<!-- Actual rows when loaded -->
						<template v-else>
							<tr v-for="dept in departments" :key="dept.id">
								<td class="text-main fw-500">
									{{ dept.name }}
								</td>
								<td>
									<span class="dept-code">{{ dept.code }}</span>
								</td>
								<td class="text-muted">
									{{ dept.description || "—" }}
								</td>
								<td class="text-center">
									<span class="employee-count">
										{{
											dept.manager
												? dept.manager.first_name +
													" " +
													dept.manager.last_name
												: "Chưa có"
										}}
									</span>
								</td>
								<td v-if="canCrudDepartment" class="text-right">
									<div class="action-group">
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
								</td>
							</tr>
							<tr v-if="departments.length === 0">
								<td
									:colspan="canCrudDepartment ? 5 : 4"
									class="empty-state"
								>
									<div class="empty-state__icon">🏢</div>
									<p class="empty-state__text">
										Không có phòng ban nào phù hợp.
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

		<ModalDialog
			:visible="isModalVisible"
			:title="isEditMode ? 'Chỉnh sửa phòng ban' : 'Thêm phòng ban mới'"
			:subtitle="
				isEditMode
					? 'Cập nhật phòng ban hiện tại'
					: 'Nhập thông tin phòng ban mới'
			"
			size="lg"
			@close="closeModal"
		>
			<form @submit.prevent="handleFormSubmit" class="department-form">
				<div class="form-grid">
					<div class="form-group">
						<label class="form-label"
							>Tên phòng ban <span class="required">*</span></label
						>
						<input
							v-model="formData.name"
							type="text"
							class="form-control"
							placeholder="Nhập tên phòng ban..."
							required
						/>
					</div>
					<div class="form-group">
						<label class="form-label"
							>Mã phòng ban <span class="required">*</span></label
						>
						<input
							v-model="formData.code"
							type="text"
							class="form-control"
							placeholder="Nhập mã phòng ban..."
							required
							:disabled="isEditMode"
						/>
					</div>
					<div v-if="isEditMode" class="form-group">
						<label class="form-label">Trưởng phòng</label>
						<select v-model="formData.manager_id" class="form-control">
							<option :value="null">Chọn trưởng phòng</option>
							<option
								v-for="employee in departmentEmployees"
								:key="employee.id"
								:value="employee.id"
							>
								{{ employee.first_name }} {{ employee.last_name
								}}{{
									employee.department
										? " - " + employee.department.name
										: ""
								}}
							</option>
						</select>
					</div>
					<div class="form-group form-group--full">
						<label class="form-label">Mô tả</label>
						<textarea
							v-model="formData.description"
							class="form-control"
							placeholder="Mô tả ngắn về phòng ban"
							rows="4"
						></textarea>
					</div>
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
						{{ isEditMode ? "Lưu thay đổi" : "Thêm phòng ban" }}
					</button>
				</div>
			</form>
		</ModalDialog>

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
}

/* ===== Header ===== */
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

/* ===== Card & Toolbar ===== */
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

/* ===== Table ===== */
.table-responsive {
	overflow-x: auto;
}

.data-table {
	width: 100%;
	border-collapse: collapse;
	text-align: left;
}

.data-table th {
	padding: var(--space-3) var(--space-2);
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
}

.data-table tbody tr:hover td {
	background: var(--bg-lighter);
}

/* Dept code pill */
.dept-code {
	display: inline-block;
	padding: 0.2rem 0.6rem;
	background: var(--bg-light);
	border-radius: var(--radius-sm);
	font-size: var(--fs-xs);
	font-weight: var(--fw-semibold);
	color: var(--text-muted);
	letter-spacing: 0.03em;
}

/* ===== Actions ===== */
.action-group {
	display: inline-flex;
	gap: 0.5rem;
}

.btn-icon {
	width: 34px;
	height: 34px;
	border-radius: var(--radius-sm);
	background: var(--bg-main);
	border: 1px solid var(--border-color);
	display: inline-flex;
	align-items: center;
	justify-content: center;
	cursor: pointer;
	transition: all 0.2s ease;
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

/* ===== Pagination ===== */
.pagination {
	display: flex;
	justify-content: center;
	align-items: center;
	gap: var(--space-3);
	padding: var(--space-3) var(--space-4);
	background: var(--bg-card);
	border-top: 1px solid var(--border-color);
}

.pagination__btn {
	width: 36px;
	height: 36px;
	display: inline-flex;
	align-items: center;
	justify-content: center;
	border: 1px solid var(--border-color);
	background: var(--bg-main);
	border-radius: var(--radius-md);
	cursor: pointer;
	transition: all 0.2s ease;
}

.pagination__btn:hover:not(:disabled) {
	border-color: var(--primary-color);
	background: var(--bg-light);
}

.pagination__btn:disabled {
	opacity: 0.35;
	cursor: not-allowed;
}

.pagination__btn img {
	width: 18px;
	height: 18px;
	opacity: 0.7;
}

.pagination__info {
	font-size: var(--fs-sm);
	color: var(--text-muted);
}

.pagination__info span {
	font-weight: var(--fw-bold);
	color: var(--text-main);
}

/* ===== Empty state ===== */
.empty-state {
	text-align: center;
	padding: var(--space-8) 0;
}

.empty-state__icon {
	font-size: 2.5rem;
	margin-bottom: var(--space-3);
}

.empty-state__text {
	color: var(--text-muted);
	font-size: var(--fs-sm);
}

/* ===== Misc ===== */
.btn__icon {
	width: 18px;
	height: 18px;
	filter: brightness(0) invert(1);
}

/* ===== Form ===== */
.department-form {
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

.form-group--full {
	grid-column: 1 / -1;
}

.form-label {
	font-size: var(--fs-sm);
	font-weight: var(--fw-semibold);
}

.form-actions {
	display: flex;
	justify-content: flex-end;
	gap: var(--space-2);
	margin-top: var(--space-2);
	padding-top: var(--space-3);
	border-top: 1px solid var(--border-color);
}

.required {
	color: var(--danger-color);
}

.spinner {
	width: 16px;
	height: 16px;
	border: 2px solid rgba(255, 255, 255, 0.3);
	border-right-color: transparent;
	border-radius: 50%;
	animation: spin 0.8s linear infinite;
	display: inline-block;
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
