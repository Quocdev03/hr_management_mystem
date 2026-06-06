<script setup>
// ─── Store & tiện ích ────────────────────────────────────────────────────────
import { ref, onMounted } from "vue";
import { storeToRefs } from "pinia";
import { useToast } from "vue-toastification";

// Store quản lý nhân viên, phòng ban, người dùng
import { useEmployeeStore } from "@/store/employee";
import { useDepartmentStore } from "@/store/department";
import { useUserStore } from "@/store/user";

// Component UI dùng chung
import ModalDialog from "@/components/ModalDialog.vue";
import ConfirmationDialog from "@/components/ConfirmationDialog.vue";
import Skeleton from "@/components/Skeleton.vue";

// Helper tiện ích
import {
	getInitials, // Lấy chữ viết tắt tên
	formatDate, // Định dạng ngày
	formatStatus, // Định dạng trạng thái
	formatCurrency, // Định dạng tiền tệ
} from "@/helpers/formatters";
import { useModalState } from "@/helpers/useModalState"; // Quản lý trạng thái modal
import { usePaginatedSearch } from "@/helpers/usePaginatedSearch"; // Tìm kiếm + phân trang
import { usePermissions } from "@/helpers/usePermissions"; // Kiểm tra quyền hạn

// Icon SVG
import plusIcon from "@/assets/svg/plus.svg";
import searchIcon from "@/assets/svg/search.svg";
import editIcon from "@/assets/svg/edit.svg";
import deleteIcon from "@/assets/svg/delete.svg";
import prevIcon from "@/assets/svg/chevron-left.svg";
import nextIcon from "@/assets/svg/chevron-right.svg";
import eyeIcon from "@/assets/svg/eye.svg";

// ─── Khởi tạo store & tiện ích ───────────────────────────────────────────────

const employeeStore = useEmployeeStore();
const departmentStore = useDepartmentStore();
const userStore = useUserStore();
const toast = useToast();

// Quyền thao tác với nhân viên
const {
	canViewEmployeeDetail, // Hàm (id) => bool — Employee chỉ xem được của mình
	canCreateEmployee,
	canEditEmployee,
	canDeleteEmployee,
	hasAnyEmployeeAction, // Có ít nhất 1 quyền → hiện cột "Thao tác"
} = usePermissions();

// Reactive refs từ store
const { employees, pagination, loading } = storeToRefs(employeeStore);
const { departments } = storeToRefs(departmentStore);
const { usersWithoutEmp } = storeToRefs(userStore); // User chưa gắn nhân viên

// ─── Modal thêm/sửa ──────────────────────────────────────────────────────────

const modalState = useModalState();
const isModalVisible = modalState.isModalVisible; // Modal có đang mở không
const isEditMode = modalState.isEditMode; // Đang sửa hay thêm mới
const openAddModal = modalState.openAddModal; // Mở modal thêm
const openEditModal = modalState.openEditModal; // Mở modal sửa
const closeModal = modalState.closeModal; // Đóng modal

// ─── Tìm kiếm & phân trang ───────────────────────────────────────────────────

const paginatedSearch = usePaginatedSearch(
	(params) => employeeStore.fetchEmployees(params),
	pagination,
);
const searchQuery = paginatedSearch.searchQuery; // Từ khoá tìm kiếm
const loadEmployees = paginatedSearch.load; // Tải danh sách nhân viên
const handlePageChange = paginatedSearch.handlePageChange; // Xử lý chuyển trang

// ─── Trạng thái local ────────────────────────────────────────────────────────

const editingEmployee = ref(null); // Nhân viên đang được chỉnh sửa
const formLoading = ref(false); // Đang xử lý submit form
const formData = ref(buildFormData()); // Dữ liệu form hiện tại
const relationsLoaded = ref(false); // Đã tải dropdown (phòng ban, user) chưa

// Trạng thái modal xoá
const isDeleteModalVisible = ref(false); // Modal xoá có đang mở không
const deletingEmployee = ref(null); // Nhân viên sắp bị xoá
const deleteMessage = ref(""); // Nội dung xác nhận xoá
const deleteLoading = ref(false); // Đang xử lý xoá

// Trạng thái modal chi tiết
const isDetailModalVisible = ref(false); // Modal chi tiết có đang mở không
const selectedEmployee = ref(null); // Nhân viên đang xem chi tiết

// ─── Xem chi tiết ────────────────────────────────────────────────────────────

// Gán nhân viên được chọn và mở modal xem chi tiết
function handleViewDetails(emp) {
	selectedEmployee.value = emp;
	isDetailModalVisible.value = true;
}

// ─── Khởi tạo dữ liệu form ───────────────────────────────────────────────────

// Tạo object formData với giá trị mặc định — tránh undefined/null gây lỗi binding
function buildFormData(data = {}) {
	const d = data ?? {};
	return {
		user_id: d.user_id ?? null,
		first_name: d.first_name ?? "",
		last_name: d.last_name ?? "",
		phone: d.phone ?? "",
		department_id: d.department_id ?? "",
		position: d.position ?? "",
		salary: d.salary ?? null,
		join_date: d.join_date ?? "",
		status: d.status ?? "active",
		gender: d.gender ?? "",
		birth_date: d.birth_date ?? "",
	};
}

// ─── Tải dữ liệu liên quan cho form ─────────────────────────────────────────

// Lấy phòng ban + user chưa gắn nhân viên để đổ vào dropdown
// Có cache: bỏ qua nếu đã tải rồi, trừ khi force = true
async function loadFormRelations(force = false) {
	if (!force && relationsLoaded.value) return;

	try {
		await Promise.all([
			departmentStore.fetchDepartments(),
			userStore.fetchUsersWithoutEmployee(),
		]);
		relationsLoaded.value = true;
	} catch (err) {
		console.error("Lỗi khi tải dữ liệu liên quan:", err);
	}
}

// ─── Thêm nhân viên ──────────────────────────────────────────────────────────

// Reset form về trạng thái rỗng rồi mở modal
async function handleAdd() {
	editingEmployee.value = null;
	formData.value = buildFormData();
	openAddModal();
	await loadFormRelations();
}

// ─── Sửa nhân viên ───────────────────────────────────────────────────────────

// Chuẩn hoá ngày (cắt phần giờ ISO 8601), điền form và mở modal sửa
async function handleEdit(emp) {
	const empCopy = { ...emp };

	// API trả ISO string (2024-01-15T00:00:00Z) → input[type=date] chỉ nhận YYYY-MM-DD
	empCopy.birth_date = empCopy.birth_date
		? empCopy.birth_date.split("T")[0]
		: "";
	empCopy.join_date = empCopy.join_date ? empCopy.join_date.split("T")[0] : "";

	editingEmployee.value = empCopy;
	formData.value = buildFormData(empCopy);
	openEditModal();
	await loadFormRelations();
}

// ─── Xoá nhân viên ───────────────────────────────────────────────────────────

// Gán nhân viên cần xoá và hiển thị modal xác nhận
function handleDelete(emp) {
	deletingEmployee.value = emp;
	deleteMessage.value = `Bạn có chắc chắn muốn xoá ${emp.first_name} ${emp.last_name}?`;
	isDeleteModalVisible.value = true;
}

// Thực hiện xoá sau khi người dùng xác nhận
async function confirmDelete() {
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
	await Promise.all([
		loadEmployees(pagination.value.page),
		departmentStore.fetchDepartments({
			page: departmentStore.pagination.page,
			limit: departmentStore.pagination.limit,
		}),
	]);
}

// ─── Submit form thêm/sửa ─────────────────────────────────────────────────────

async function handleFormSubmit() {
	formLoading.value = true;
	let res;

	if (isEditMode.value === true) {
		// ── Chế độ SỬA: partial update — chỉ gửi field thực sự thay đổi ──

		const original = buildFormData(editingEmployee.value);

		/**
		 * Chuẩn hoá toàn bộ giá trị về string trước khi so sánh,
		 * tránh false-diff do khác kiểu (ví dụ: số 3 vs chuỗi "3").
		 * - null/undefined  → ""
		 * - user_id / department_id / salary → String(v)  (ID số từ API)
		 * - object → JSON.stringify
		 * - còn lại → String(v)
		 */
		function normalizeForCompare(obj) {
			const out = {};
			Object.keys(obj).forEach((k) => {
				const v = obj[k];
				if (v === null || typeof v === "undefined") {
					out[k] = "";
					return;
				}
				if (k === "user_id" || k === "department_id" || k === "salary") {
					out[k] = String(v);
					return;
				}
				out[k] = typeof v === "object" ? JSON.stringify(v) : String(v);
			});
			return out;
		}

		const normOriginal = normalizeForCompare(original);
		const normForm = normalizeForCompare(formData.value);

		// Lấy những field có giá trị khác so với ban đầu → build payload
		const payload = Object.fromEntries(
			Object.keys(normOriginal)
				.filter((k) => normForm[k] !== normOriginal[k])
				.map((k) => [k, formData.value[k]]),
		);

		// Không có gì thay đổi → báo và thoát sớm
		if (Object.keys(payload).length === 0) {
			toast.info("Không có dữ liệu thay đổi");
			formLoading.value = false;
			return;
		}

		// Ép kiểu số cho backend Go — tránh gửi string gây lỗi parse
		if ("user_id" in payload) {
			// user_id rỗng/null → gửi 0 (sentinel "huỷ liên kết tài khoản")
			payload.user_id =
				formData.value.user_id == null || formData.value.user_id === ""
					? 0
					: Number(formData.value.user_id);
		}
		if ("department_id" in payload) {
			payload.department_id = Number(payload.department_id);
		}

		res = await employeeStore.updateEmployee(
			editingEmployee.value.id,
			payload,
		);
	} else {
		// ── Chế độ THÊM MỚI: gửi toàn bộ formData ──

		const data = { ...formData.value };

		// Chuyển về number hoặc null — backend không nhận chuỗi cho ID
		data.user_id =
			data.user_id !== "" && data.user_id != null
				? Number(data.user_id)
				: null;

		// Bỏ hẳn field nếu rỗng (undefined → không có key trong JSON)
		data.department_id =
			data.department_id !== "" ? Number(data.department_id) : undefined;

		res = await employeeStore.createEmployee(data);
	}

	formLoading.value = false;

	if (res.success === false) {
		toast.error(res.message);
		return;
	}

	relationsLoaded.value = false; // Buộc tải lại dropdown lần mở form tiếp theo

	if (isEditMode.value === true) {
		toast.success(res.message);
	}

	closeModal();
	await Promise.all([
		loadEmployees(pagination.value.page),
		departmentStore.fetchDepartments({
			page: departmentStore.pagination.page,
			limit: departmentStore.pagination.limit,
		}),
	]);
}

// ─── Khởi tạo trang ──────────────────────────────────────────────────────────

onMounted(async () => {
	await loadEmployees();
});
</script>

<template>
	<div class="employee-view">
		<!-- ===== Tiêu đề trang ===== -->
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

		<!-- ===== Nội dung chính ===== -->
		<main class="content-card">
			<!-- Thanh tìm kiếm -->
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

			<!-- Bảng dữ liệu nhân viên -->
			<div class="table-responsive">
				<table class="data-table">
					<thead>
						<tr>
							<th>Nhân viên</th>
							<th>Liên hệ</th>
							<th>Phòng ban / Chức vụ</th>
							<th>Ngày vào làm</th>
							<th>Trạng thái</th>
							<!-- Chỉ hiện cột thao tác khi có ít nhất 1 quyền -->
							<th v-if="hasAnyEmployeeAction" class="text-right">
								Thao tác
							</th>
						</tr>
					</thead>
					<tbody>
						<!-- Skeleton placeholder khi đang tải -->
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
								<td><Skeleton type="text" width="100px" /></td>
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
								<td><Skeleton type="text" width="80px" /></td>
								<td><Skeleton type="badge" /></td>
								<td v-if="hasAnyEmployeeAction" class="text-right">
									<div class="action-group">
										<Skeleton type="btn" />
										<Skeleton type="btn" />
									</div>
								</td>
							</tr>
						</template>

						<!-- Dữ liệu thực sau khi tải xong -->
						<template v-else>
							<tr v-for="emp in employees" :key="emp.id">
								<!-- Thông tin nhân viên + tài khoản liên kết -->
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
											<!-- Hiển thị email + username nếu đã liên kết tài khoản -->
											<template v-if="emp.user">
												<span class="user-info__email">{{
													emp.user.email
												}}</span>
												<span class="user-info__email"
													>User: {{ emp.user.user_name }}</span
												>
											</template>
											<span v-else class="user-info__email"
												>Chưa có tài khoản</span
											>
										</div>
									</div>
								</td>

								<td>
									<span class="text-main fw-500">{{
										emp.phone || "—"
									}}</span>
								</td>

								<td>
									<div class="job-info">
										<h1 class="job-info__dept">
											{{ emp.department?.name || "N/A" }}
										</h1>
										<span class="job-info__pos">{{
											emp.position || "Nhân viên"
										}}</span>
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

								<!-- Nút thao tác: xem / sửa / xoá (tuỳ quyền) -->
								<td v-if="hasAnyEmployeeAction" class="text-right">
									<div class="action-group">
										<!-- HR/Admin xem tất cả, Employee chỉ xem được của mình -->
										<button
											v-if="canViewEmployeeDetail(emp.id)"
											class="btn-icon btn-icon--detail"
											title="Xem chi tiết"
											@click="handleViewDetails(emp)"
										>
											<img :src="eyeIcon" alt="detail" />
										</button>
										<!-- Chỉ HR và Admin mới sửa/xoá được -->
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

							<!-- Empty state khi không có kết quả -->
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

		<!-- ===== Modal Thêm / Sửa nhân viên ===== -->
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
					<!-- Họ & Tên -->
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
						<label class="form-label">Giới tính</label>
						<select v-model="formData.gender" class="form-control">
							<option value="" disabled>Chọn giới tính</option>
							<option value="male">Nam</option>
							<option value="female">Nữ</option>
						</select>
					</div>

					<!-- Ngày sinh -->
					<div class="form-group">
						<label class="form-label">Ngày sinh</label>
						<input
							v-model="formData.birth_date"
							type="date"
							class="form-control"
						/>
					</div>

					<!-- Số điện thoại -->
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

					<!-- Phòng ban -->
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

					<!-- Chức vụ -->
					<div class="form-group">
						<label class="form-label">Chức vụ</label>
						<input
							v-model="formData.position"
							type="text"
							class="form-control"
							placeholder="Ví dụ: Backend Developer"
						/>
					</div>

					<!-- Mức lương -->
					<div class="form-group">
						<label class="form-label">Mức lương (VNĐ)</label>
						<input
							v-model.number="formData.salary"
							type="number"
							class="form-control"
							placeholder="0"
						/>
					</div>

					<!--
						Ngày vào làm bị khoá khi sửa — không cho phép thay đổi
						ngày onboard sau khi đã tạo hồ sơ
					-->
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

					<!-- Trạng thái -->
					<div class="form-group">
						<label class="form-label">Trạng thái</label>
						<select v-model="formData.status" class="form-control">
							<option value="active">Đang làm việc</option>
							<option value="inactive">Đã nghỉ việc</option>
						</select>
					</div>
				</div>

				<!--
					Dropdown liên kết tài khoản:
					- Khi sửa: hiển thị thêm option user hiện tại (có thể không có trong usersWithoutEmp)
					- Khi thêm: chỉ hiện user chưa gắn nhân viên
				-->
				<div class="form-group">
					<span class="form-label">Liên kết người dùng</span>
					<select v-model="formData.user_id" class="form-control">
						<option :value="null">Không liên kết</option>
						<option
							v-if="editingEmployee?.user"
							:value="editingEmployee.user.id"
						>
							{{ editingEmployee.user.email }} (hiện tại)
						</option>
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

		<!-- ===== Modal Xác nhận xoá ===== -->
		<ConfirmationDialog
			:visible="isDeleteModalVisible"
			title="Xác nhận xoá nhân viên"
			:message="deleteMessage"
			:loading="deleteLoading"
			@confirm="confirmDelete"
			@cancel="isDeleteModalVisible = false"
		/>

		<!-- ===== Modal Chi tiết nhân viên ===== -->
		<ModalDialog
			:visible="isDetailModalVisible"
			title="Chi tiết nhân viên"
			subtitle="Thông tin hồ sơ chi tiết của nhân viên trong hệ thống"
			size="lg"
			@close="isDetailModalVisible = false"
		>
			<div v-if="selectedEmployee" class="detail-container">
				<!-- Avatar & Tên -->
				<div class="detail-header">
					<div class="detail-avatar">
						{{
							getInitials(
								selectedEmployee.first_name,
								selectedEmployee.last_name,
							)
						}}
					</div>
					<div class="detail-title-info">
						<h2 class="detail-name">
							{{ selectedEmployee.first_name }}
							{{ selectedEmployee.last_name }}
						</h2>
						<span class="detail-position">{{
							selectedEmployee.position || "Chưa có chức vụ"
						}}</span>
					</div>
				</div>

				<!-- Thông tin cá nhân -->
				<div class="detail-section">
					<h3 class="section-title">Thông tin cá nhân</h3>
					<div class="detail-grid">
						<div class="detail-item">
							<span class="detail-label">Giới tính:</span>
							<!-- Map giá trị enum về nhãn tiếng Việt -->
							<span class="detail-val">{{
								selectedEmployee.gender === "male"
									? "Nam"
									: selectedEmployee.gender === "female"
										? "Nữ"
										: "Khác"
							}}</span>
						</div>
						<div class="detail-item">
							<span class="detail-label">Ngày sinh:</span>
							<span class="detail-val">{{
								formatDate(selectedEmployee.birth_date)
							}}</span>
						</div>
						<div class="detail-item">
							<span class="detail-label">Số điện thoại:</span>
							<span class="detail-val">{{
								selectedEmployee.phone || "—"
							}}</span>
						</div>
					</div>
				</div>

				<!-- Thông tin công việc -->
				<div class="detail-section">
					<h3 class="section-title">Thông tin công việc</h3>
					<div class="detail-grid">
						<div class="detail-item">
							<span class="detail-label">Phòng ban:</span>
							<span class="detail-val">{{
								selectedEmployee.department?.name || "N/A"
							}}</span>
						</div>
						<div class="detail-item">
							<span class="detail-label">Ngày vào làm:</span>
							<span class="detail-val">{{
								formatDate(selectedEmployee.join_date)
							}}</span>
						</div>
						<div class="detail-item">
							<span class="detail-label">Mức lương:</span>
							<span class="detail-val salary-text">{{
								formatCurrency(selectedEmployee.salary)
							}}</span>
						</div>
						<div class="detail-item">
							<span class="detail-label">Trạng thái:</span>
							<span
								:class="[
									'status-badge',
									`status-badge--${selectedEmployee.status}`,
								]"
							>
								{{ formatStatus(selectedEmployee.status) }}
							</span>
						</div>
					</div>
				</div>

				<!-- Tài khoản hệ thống -->
				<div class="detail-section">
					<h3 class="section-title">Tài khoản hệ thống</h3>
					<div v-if="selectedEmployee.user" class="detail-grid">
						<div class="detail-item">
							<span class="detail-label">Tên đăng nhập:</span>
							<span class="detail-val">{{
								selectedEmployee.user.user_name
							}}</span>
						</div>
						<div class="detail-item">
							<span class="detail-label">Email liên kết:</span>
							<span class="detail-val">{{
								selectedEmployee.user.email
							}}</span>
						</div>
						<div class="detail-item">
							<span class="detail-label">Quyền hạn:</span>
							<span class="detail-val text-uppercase fw-600">
								{{ selectedEmployee.user.role?.name || "N/A" }}
							</span>
						</div>
					</div>
					<div v-else class="empty-account">
						Nhân viên này chưa liên kết với tài khoản hệ thống.
					</div>
				</div>

				<div class="form-actions">
					<button
						type="button"
						class="btn btn--secondary"
						@click="isDetailModalVisible = false"
					>
						Đóng
					</button>
				</div>
			</div>
		</ModalDialog>
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

/* ===== Detail Modal ===== */
.detail-container {
	display: flex;
	flex-direction: column;
	gap: var(--space-4);
}

.detail-header {
	display: flex;
	align-items: center;
	gap: var(--space-3);
	padding-bottom: var(--space-3);
	border-bottom: 1px solid var(--border-color);
}

.detail-avatar {
	width: 56px;
	height: 56px;
	border-radius: var(--radius-lg);
	background: #eff6ff;
	color: var(--primary-color);
	display: flex;
	align-items: center;
	justify-content: center;
	font-weight: var(--fw-bold);
	font-size: var(--fs-lg);
	border: 1px solid rgba(59, 130, 246, 0.15);
}

.detail-title-info {
	display: flex;
	flex-direction: column;
	gap: 2px;
}

.detail-name {
	font-size: var(--fs-lg);
	font-weight: var(--fw-bold);
	color: var(--text-main);
	margin: 0;
}

.detail-position {
	font-size: var(--fs-sm);
	color: var(--text-muted);
}

.detail-section {
	display: flex;
	flex-direction: column;
	gap: var(--space-2);
}

.section-title {
	font-size: var(--fs-sm);
	font-weight: var(--fw-bold);
	text-transform: uppercase;
	color: var(--text-muted);
	letter-spacing: 0.05em;
	margin: 0;
	padding-bottom: 4px;
	border-bottom: 1px dashed var(--border-color);
}

.detail-grid {
	display: grid;
	grid-template-columns: repeat(2, 1fr);
	gap: var(--space-2) var(--space-4);
}

.detail-item {
	display: flex;
	flex-direction: column;
	gap: 5px;
}

.detail-label {
	font-size: var(--fs-xs);
	color: var(--text-light);
}

.detail-val {
	font-size: var(--fs-sm);
	color: var(--text-main);
	font-weight: var(--fw-medium);
}

.salary-text {
	color: #0f766e;
	font-weight: var(--fw-semibold);
}

.empty-account {
	font-size: var(--fs-sm);
	color: var(--text-light);
	font-style: italic;
	padding: var(--space-2) 0;
}
</style>
