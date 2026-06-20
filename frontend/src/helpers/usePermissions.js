import { computed } from "vue";
import { useAuthStore } from "@/store/auth";

/**
 * Composable kiểm tra phân quyền RBAC dựa trên danh sách permissions.
 * Dữ liệu permissions được trả về từ API và lưu trong authStore.user.permissions
 */
export function usePermissions() {
    const authStore = useAuthStore();

    const roleName = computed(() => authStore.user?.role?.name?.toLowerCase() || "");
    const permissions = computed(() => {
        return authStore.user?.permissions || authStore.userProfile?.permissions || [];
    });

    const isAdmin = computed(() => roleName.value === "admin");
    const isHR = computed(() => roleName.value === "hr");
    const isEmployeeRole = computed(() => roleName.value === "employee");

    const hasPermission = (code) => {
        if (isAdmin.value) return true; // Admin có toàn quyền trên UI
        const list = permissions.value || [];
        return list.includes(code);
    };

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
    const canViewEmployeeDetail = computed(() => (targetEmployeeId) => {
        if (hasPermission("employee.read")) {
            return true;
        }
        return authStore.user?.employee?.id === targetEmployeeId;
    });

    // === Quyền trên Phòng ban ===
    const canViewDepartmentList = computed(() => hasPermission("department.read"));
    const canCreateDepartment = computed(() => hasPermission("department.create"));
    const canEditDepartment = computed(() => hasPermission("department.update"));
    const canDeleteDepartment = computed(() => hasPermission("department.delete"));
    const canCrudDepartment = computed(
        () =>
            canCreateDepartment.value ||
            canEditDepartment.value ||
            canDeleteDepartment.value
    );

    // === Quyền trên Chức vụ ===
    const canViewPositionList = computed(() => hasPermission("position.read"));
    const canCrudPosition = computed(
        () =>
            hasPermission("position.create") ||
            hasPermission("position.update") ||
            hasPermission("position.delete")
    );

    // === Quyền Quản lý User ===
    const canManageUsers = computed(
        () =>
            hasPermission("user.read") ||
            hasPermission("user.create") ||
            hasPermission("user.update") ||
            hasPermission("user.delete")
    );

    // === Quyền Quản lý Vai trò (Role) ===
    const canManageRoles = computed(
        () =>
            hasPermission("role.read") ||
            hasPermission("role.create") ||
            hasPermission("role.update") ||
            hasPermission("role.delete")
    );

    // Helper: có ít nhất 1 quyền thao tác trên bảng NV không? (hiện/ẩn cột Thao tác)
    const hasAnyEmployeeAction = computed(
        () =>
            canCreateEmployee.value ||
            canEditEmployee.value ||
            canDeleteEmployee.value
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
