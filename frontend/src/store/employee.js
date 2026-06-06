import api from "@/api";
import { defineStore } from "pinia";
import { ref } from "vue";

export const useEmployeeStore = defineStore("employee", () => {
	// State
	const employees = ref([]);
	const pagination = ref({ total: 0, page: 1, limit: 10, totalPages: 0 });
	const loading = ref(false);

	// Fetch Employees
	async function fetchEmployees(params = {}) {
		loading.value = true;

		try {
			const res = await api.get("/employees", { params });

			if (res.success) {
				const { data } = res;

				employees.value = data.items || [];
				pagination.value = {
					total: data.total || 0,
					page: data.page || 1,
					limit: data.limit || 10,
					totalPages: data.total_pages || 0,
				};
			}

			return res;
		} catch (error) {
			console.error("Fetch employees error:", error);
			return {
				success: false,
				message: error?.message || "Lỗi tải danh sách nhân viên",
			};
		} finally {
			loading.value = false;
		}
	}

	async function fetchEmployeesForSelect(params = {}) {
		try {
			const res = await api.get("/employees", { params });

			if (res.success) {
				return {
					success: true,
					items: res.data?.items || [],
				};
			}

			return {
				success: false,
				message: res.message || "Lỗi tải danh sách nhân viên",
			};
		} catch (error) {
			console.error("Fetch employees for select error:", error);
			return {
				success: false,
				message: error?.message || "Lỗi tải danh sách nhân viên",
			};
		}
	}

	// Create
	async function createEmployee(data) {
		try {
			const res = await api.post("/employees", data);
			return res;
		} catch (error) {
			console.error("Create employee error:", error);
			return {
				success: false,
				message: error?.message || "Lỗi tạo nhân viên",
			};
		}
	}

	// Update
	async function updateEmployee(id, data) {
		try {
			const res = await api.patch(`/employees/${id}`, data);
			return res;
		} catch (error) {
			console.error("Update employee error:", error);
			return {
				success: false,
				message: error?.message || "Lỗi cập nhật nhân viên",
			};
		}
	}

	// Delete
	async function deleteEmployee(id) {
		try {
			const res = await api.delete(`/employees/${id}`);
			return res;
		} catch (error) {
			console.error("Delete employee error:", error);
			return {
				success: false,
				message: error?.message || "Lỗi xoá nhân viên",
			};
		}
	}

	return {
		employees,
		pagination,
		loading,
		fetchEmployees,
		fetchEmployeesForSelect,
		createEmployee,
		updateEmployee,
		deleteEmployee,
	};
});
