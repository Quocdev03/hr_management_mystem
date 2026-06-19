import api from "@/api";
import { defineStore } from "pinia";
import { ref } from "vue";

export const usePositionStore = defineStore("position", () => {
	// State
	const positions = ref([]);
	const loading = ref(false);

	// Actions
	async function fetchPositions(params = {}) {
		loading.value = true;
		try {
			const res = await api.get("/positions", { params });
			if (res.success) {
				positions.value = res.data || [];
			}
			return res;
		} catch (error) {
			console.error("Fetch positions error:", error);
			return {
				success: false,
				message: error?.message || "Lỗi tải danh sách chức vụ",
			};
		} finally {
			loading.value = false;
		}
	}

	async function createPosition(data) {
		try {
			const res = await api.post("/positions", data);
			return res;
		} catch (error) {
			console.error("Create position error:", error);
			return {
				success: false,
				message: error?.message || "Lỗi tạo chức vụ",
			};
		}
	}

	async function updatePosition(id, data) {
		try {
			const res = await api.patch(`/positions/${id}`, data);
			return res;
		} catch (error) {
			console.error("Update position error:", error);
			return {
				success: false,
				message: error?.message || "Lỗi cập nhật chức vụ",
			};
		}
	}

	async function deletePosition(id) {
		try {
			const res = await api.delete(`/positions/${id}`);
			return res;
		} catch (error) {
			console.error("Delete position error:", error);
			return {
				success: false,
				message: error?.message || "Lỗi xoá chức vụ",
			};
		}
	}

	return {
		positions,
		loading,
		fetchPositions,
		createPosition,
		updatePosition,
		deletePosition,
	};
});
