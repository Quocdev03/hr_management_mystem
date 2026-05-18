import api from "@/api";
import { defineStore } from "pinia";
import { ref } from "vue";

export const useDepartmentStore = defineStore("department", () => {
  // ===== State =====
  const departments = ref([]);
  const pagination = ref({
    total: 0,
    page: 1,
    limit: 10,
    totalPages: 0,
  });
  const loading = ref(false);

  // ===== Actions =====
  async function fetchDepartments(params = {}) {
    loading.value = true;

    try {
      const res = await api.get("/departments", {
        params,
      });

      if (res.success) {
        const data = res.data;
        departments.value = data.items || [];
        pagination.value = {
          total: data.total || 0,
          page: data.page || params.page || 1,
          limit: data.limit || params.limit || 10,
          totalPages: data.total_pages || 0,
        };
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

  async function createDepartment(data) {
    try {
      const res = await api.post("/departments", data);
      return res;
    } catch (error) {
      console.error("Create department error:", error);
      return {
        success: false,
        message: "Lỗi tạo phòng ban",
      };
    }
  }

  async function updateDepartment(id, data) {
    try {
      const res = await api.put(`/departments/${id}`, data);
      return res;
    } catch (error) {
      console.error("Update department error:", error);
      return {
        success: false,
        message: "Lỗi cập nhật phòng ban",
      };
    }
  }

  async function deleteDepartment(id) {
    try {
      const res = await api.delete(`/departments/${id}`);
      return res;
    } catch (error) {
      console.error("Delete department error:", error);
      return {
        success: false,
        message: "Lỗi xoá phòng ban",
      };
    }
  }

  return {
    departments,
    pagination,
    loading,

    fetchDepartments,
    createDepartment,
    updateDepartment,
    deleteDepartment,
  };
});
