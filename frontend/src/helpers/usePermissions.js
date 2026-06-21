import { computed } from "vue";
import { useAuthStore } from "@/store/auth";

/**
 * Composable kiểm tra phân quyền RBAC dựa trên danh sách permissions.
 * Dữ liệu permissions được trả về từ API và lưu trong authStore.user.permissions
 */
export const usePermissions = () => {
	const authStore = useAuthStore();

	const roleName = computed(() => authStore.user?.role?.name?.toLowerCase() ?? "");
	const permissions = computed(() => authStore.user?.permissions ?? authStore.userProfile?.permissions ?? []);

	const isAdmin = computed(() => roleName.value === "admin");
	const isHR = computed(() => roleName.value === "hr");
	const isEmployeeRole = computed(() => roleName.value === "employee");

	const hasPermission = (code) => isAdmin.value || permissions.value.includes(code);

	// === Quyền trên Nhân viên ===
	const canViewEmployeeList = computed(() => hasPermission("employee.read"));
	const canCreateEmployee = computed(() => hasPermission("employee.create"));
	const canEditEmployee = computed(() => hasPermission("employee.update"));
	const canDeleteEmployee = computed(() => hasPermission("employee.delete"));

	/**
	 * Employee xem chi tiết:
	 * - Nếu có quyền employee.read -> xem được tất cả
	 * - Nếu không có quyền, CHỈ ĐƯỢC XEM chính mình
	 */
	const canViewEmployeeDetail = computed(() => (targetEmployeeId) => 
		hasPermission("employee.read") || authStore.user?.employee?.id === targetEmployeeId
	);

	// === Quyền trên Phòng ban ===
	const canViewDepartmentList = computed(() => hasPermission("department.read"));
	const canCreateDepartment = computed(() => hasPermission("department.create"));
	const canEditDepartment = computed(() => hasPermission("department.update"));
	const canDeleteDepartment = computed(() => hasPermission("department.delete"));
	const canCrudDepartment = computed(() => 
		["department.create", "department.update", "department.delete"].some(hasPermission)
	);

	// === Quyền trên Chức vụ ===
	const canViewPositionList = computed(() => hasPermission("position.read"));
	const canCrudPosition = computed(() => 
		["position.create", "position.update", "position.delete"].some(hasPermission)
	);

	// === Quyền Quản lý User ===
	const canManageUsers = computed(() => 
		["user.read", "user.create", "user.update", "user.delete"].some(hasPermission)
	);

	// === Quyền Quản lý Vai trò (Role) ===
	const canManageRoles = computed(() => 
		["role.read", "role.create", "role.update", "role.delete"].some(hasPermission)
	);

	// Helper: có ít nhất 1 quyền thao tác trên bảng NV không? (hiện/ẩn cột Thao tác)
	const hasAnyEmployeeAction = computed(() => 
		["employee.create", "employee.update", "employee.delete"].some(hasPermission)
	);

	return {
		roleName,
		isAdmin,
		isHR,
		isEmployee: isEmployeeRole, // alias for backward compatibility
		
		hasPermission,
		
		canViewEmployeeList,
		canViewEmployeeDetail,
		canCreateEmployee,
		canEditEmployee,
		canDeleteEmployee,
		hasAnyEmployeeAction,

		canViewDepartmentList,
		canCrudDepartment,

		canViewPositionList,
		canCrudPosition,

		canManageUsers,
		canManageRoles
	};
}
