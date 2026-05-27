import { computed } from "vue";
import { useAuthStore } from "@/store/auth";

/**
 * Composable tập trung kiểm tra phân quyền RBAC.
 * Đọc role từ authStore.user và expose các computed boolean.
 *
 * Bảng quyền:
 * | Action         | Employee | HR  | Admin |
 * | -------------- | -------- | --- | ----- |
 * | Xem nhân viên  | ✅       | ✅  | ✅    |
 * | Tạo nhân viên  | ❌       | ✅  | ✅    |
 * | Sửa nhân viên  | ❌       | ✅  | ✅    |
 * | Xóa nhân viên  | ❌       | ❌  | ✅    |
 * | Xem phòng ban  | ✅       | ✅  | ✅    |
 * | CRUD phòng ban | ❌       | ❌  | ✅    |
 * | Quản lý User   | ❌       | ❌  | ✅    |
 */
export function usePermissions() {
	const authStore = useAuthStore();

	// role.name trả về từ backend: "admin", "hr", "employee"
	const roleName = computed(() => {
		return authStore.user?.role?.name || "";
	});

	// === Kiểm tra role ===
	const isAdmin = computed(() => roleName.value === "admin");
	const isHR = computed(() => roleName.value === "hr");
	const isEmployee = computed(() => roleName.value === "employee");

	// === Quyền trên Nhân viên ===
	const canCreateEmployee = computed(() => isAdmin.value || isHR.value);
	const canEditEmployee = computed(() => isAdmin.value || isHR.value);
	const canDeleteEmployee = computed(() => isAdmin.value);

	// === Quyền trên Phòng ban ===
	const canCrudDepartment = computed(() => isAdmin.value);

	// === Quyền Quản lý User ===
	const canManageUsers = computed(() => isAdmin.value);

	// Helper: có quyền thao tác trên bảng nhân viên không? (hiện cột Thao tác)
	const hasAnyEmployeeAction = computed(() => canCreateEmployee.value || canEditEmployee.value || canDeleteEmployee.value);

	return {
		roleName,
		isAdmin,
		isHR,
		isEmployee,
		canCreateEmployee,
		canEditEmployee,
		canDeleteEmployee,
		canCrudDepartment,
		canManageUsers,
		hasAnyEmployeeAction,
	};
}
