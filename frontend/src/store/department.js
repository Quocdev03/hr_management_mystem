import api from "@/api";
import { defineStore } from "pinia";
import { ref } from "vue";

export const useDepartmentStore = defineStore("department", () => {
	// ===== State =====
	const departments = ref([]);
	const loading = ref(false);

	// ===== Actions =====
	async function fetchDepartments(params = {}) {
		loading.value = true;

		try {
			const res = await api.get("/departments", {
				params,
			});

			// Normalize response
			if (res.success) {
				departments.value = res.data.items || [];
			}

			return res;
		} catch (error) {
			console.error("Fetch departments error:", error);

			return {
				success: false,
				message: "Lỗi tải danh sách phòng ban",
			};
		} finally {
			loading.value = false;
		}
	}

	return {
		// state
		departments,
		loading,

		// actions
		fetchDepartments,
	};
});
