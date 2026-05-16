import api from "@/api";
import { defineStore } from "pinia";
import { ref } from "vue";

export const useEmployeeStore = defineStore("employee", () => {
	// ===== State =====

	const employees = ref([]);

	const pagination = ref({
		total: 0,
		page: 1,
		limit: 10,
		totalPages: 0,
	});

	const loading = ref(false);

	// ===== Fetch Employees =====

	async function fetchEmployees(params = {}) {
		loading.value = true;

		try {
			const res = await api.get("/employees", {
				params,
			});

			// DEBUG
			console.log("EMPLOYEE API:", res.data);

			if (res.success) {
				const data = res.data;

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
			console.error(error);

			return {
				success: false,
				message: "Lỗi tải danh sách nhân viên",
			};
		} finally {
			loading.value = false;
		}
	}

	// ===== Create =====

	async function createEmployee(data) {
		try {
			const res = await api.post("/employees", data);

			return res;
		} catch (error) {
			return {
				success: false,
				message: "Lỗi tạo nhân viên",
			};
		}
	}

	// ===== Update =====

	async function updateEmployee(id, data) {
		try {
			const res = await api.put(`/employees/${id}`, data);

			return res;
		} catch (error) {
			return {
				success: false,
				message: "Lỗi cập nhật nhân viên",
			};
		}
	}

	// ===== Delete =====

	async function deleteEmployee(id) {
		try {
			const res = await api.delete(`/employees/${id}`);

			return res;
		} catch (error) {
			return {
				success: false,
				message: "Lỗi xoá nhân viên",
			};
		}
	}

	return {
		employees,
		pagination,
		loading,

		fetchEmployees,
		createEmployee,
		updateEmployee,
		deleteEmployee,
	};
});
