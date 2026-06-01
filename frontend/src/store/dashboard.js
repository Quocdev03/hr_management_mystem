import api from "@/api";
import { defineStore } from "pinia";
import { ref } from "vue";

export const useDashboardStore = defineStore("dashboard", () => {
	// ===== State =====
	const stats = ref({
		total_users: 0,
		total_departments: 0,
		total_employees: 0,
		total_employees_active: 0,
		department_stats: [],
	});

	const loading = ref(false);

	// ===== Actions =====
	async function fetchDashboard() {
		loading.value = true;

		try {
			const res = await api.get("/dashboard/stats");

			// Normalize response
			if (res.success) {
				stats.value = res.data;
			}

			return res;
		} catch (error) {
			console.error("Fetch dashboard error:", error);
			return {
				success: false,
				message: error?.message || "Lỗi tải dashboard",
			};
		} finally {
			loading.value = false;
		}
	}

	return {
		// state
		stats,
		loading,

		// actions
		fetchDashboard,
	};
});
