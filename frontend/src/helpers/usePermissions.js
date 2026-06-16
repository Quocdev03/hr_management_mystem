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

	const roleName = computed(() => authStore.user?.role?.name || "");
	const permissions = computed(() => {
		return (
			authStore.user?.permissions || authStore.userProfile?.permissions || []
		);
	});

	const hasPermission = (code) => {
		const list = permissions.value || [];
		return list.includes(code) || list.includes(`admin:${code}`);
	};

	const isAdmin = computed(
		() =>
			roleName.value === "admin" ||
			hasPermission("admin") ||
			hasPermission("role.admin"),
	);
	const isHR = computed(
		() =>
			roleName.value === "hr" ||
			hasPermission("hr") ||
			hasPermission("role.hr"),
	);
	const isEmployee = computed(() => roleName.value === "employee");

	// === Quyền trên Nhân viên ===
	const canViewEmployeeList = computed(
		() => hasPermission("employee.read") || isAdmin.value || isHR.value,
	);
	const canCreateEmployee = computed(
		() => hasPermission("employee.create") || isAdmin.value || isHR.value,
	);
	const canEditEmployee = computed(
		() => hasPermission("employee.update") || isAdmin.value || isHR.value,
	);
	const canDeleteEmployee = computed(
		() => hasPermission("employee.delete") || isAdmin.value,
	);

	/**
	 * Employee chỉ được xem chi tiết của chính mình.
	 * HR và Admin xem được tất cả.
	 * Nhận targetEmployeeId để so sánh với employee đang đăng nhập.
	 */
	const canViewEmployeeDetail = computed(() => (targetEmployeeId) => {
		if (
			hasPermission("employee.read") &&
			(isAdmin.value ||
				isHR.value ||
				authStore.user?.employee?.id === targetEmployeeId)
		) {
			return true;
		}
		return authStore.user?.employee?.id === targetEmployeeId;
	});

	// === Quyền trên Phòng ban ===
	const canCrudDepartment = computed(
		() =>
			hasPermission("department.create") ||
			hasPermission("department.update") ||
			hasPermission("department.delete") ||
			isAdmin.value,
	);

	// === Quyền Quản lý User ===
	const canManageUsers = computed(
		() =>
			hasPermission("user.read") ||
			hasPermission("user.create") ||
			hasPermission("user.update") ||
			hasPermission("user.delete") ||
			isAdmin.value,
	);

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
