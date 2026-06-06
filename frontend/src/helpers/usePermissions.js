import { computed } from "vue";
import { useAuthStore } from "@/store/auth";

/**
 * Composable tập trung kiểm tra phân quyền RBAC.
 * Đọc role từ authStore.user và expose các computed boolean.
 *
 * Bảng quyền:
 * | Action              | Employee    | HR  | Admin |
 * | ------------------- | ----------- | --- | ----- |
 * | Xem danh sách NV    | ❌          | ✅  | ✅    |
 * | Xem chi tiết NV     | Chỉ của mình| ✅  | ✅    |
 * | Tạo nhân viên       | ❌          | ✅  | ✅    |
 * | Sửa nhân viên       | ❌          | ✅  | ✅    |
 * | Xóa nhân viên       | ❌          | ❌  | ✅    |
 * | Xem phòng ban       | ✅          | ✅  | ✅    |
 * | CRUD phòng ban      | ❌          | ❌  | ✅    |
 * | Quản lý User        | ❌          | ❌  | ✅    |
 */
export function usePermissions() {
	const authStore = useAuthStore();

	// role.name trả về từ backend: "admin", "hr", "employee"
	const roleName = computed(() => authStore.user?.role?.name || "");

	// === Kiểm tra role ===
	const isAdmin = computed(() => roleName.value === "admin");
	const isHR = computed(() => roleName.value === "hr");
	const isEmployee = computed(() => roleName.value === "employee");

	// === Quyền trên Nhân viên ===
	const canViewEmployeeList = computed(() => isAdmin.value || isHR.value);
	const canCreateEmployee = computed(() => isAdmin.value || isHR.value);
	const canEditEmployee = computed(() => isAdmin.value || isHR.value);
	const canDeleteEmployee = computed(() => isAdmin.value);

	/**
	 * Employee chỉ được xem chi tiết của chính mình.
	 * HR và Admin xem được tất cả.
	 * Nhận targetEmployeeId để so sánh với employee đang đăng nhập.
	 */
	const canViewEmployeeDetail = computed(() => (targetEmployeeId) => {
		if (isAdmin.value || isHR.value) return true;
		if (isEmployee.value) {
			return authStore.user?.employee?.id === targetEmployeeId;
		}
		return false;
	});

	// === Quyền trên Phòng ban ===
	const canCrudDepartment = computed(() => isAdmin.value);

	// === Quyền Quản lý User ===
	const canManageUsers = computed(() => isAdmin.value);

	// Helper: có ít nhất 1 quyền thao tác trên bảng NV không? (hiện/ẩn cột Thao tác)
	const hasAnyEmployeeAction = computed(
		() =>
			canCreateEmployee.value ||
			canEditEmployee.value ||
			canDeleteEmployee.value,
	);

	return {
		roleName,
		isAdmin,
		isHR,
		isEmployee,
		canViewEmployeeList,
		canViewEmployeeDetail,
		canCreateEmployee,
		canEditEmployee,
		canDeleteEmployee,
		canCrudDepartment,
		canManageUsers,
		hasAnyEmployeeAction,
	};
}
