import api from '@/api';
import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useUserStore = defineStore('user', () => {
  // State
  const users = ref([]);
  const usersWithoutEmp = ref([]);
  const pagination = ref({ total: 0, page: 1, limit: 10, totalPages: 0 });
  const loading = ref(false);

  // Fetch User
  async function fetchUser(params = {}) {
    loading.value = true;

    try {
      const res = await api.get('/users', { params });

      if (res.success) {
        const { data } = res;

        users.value = data.items || [];
        pagination.value = {
          total: data.total || 0,
          page: data.page || 1,
          limit: data.limit || 10,
          totalPages: data.total_pages || 0,
        };
      }
      return res;
    } catch (error) {
      console.error('Fetch users error:', error);
      return {
        success: false,
        message: error?.message || 'Lỗi tải danh sách người dùng',
      };
    } finally {
      loading.value = false;
    }
  }

  // UsersWithoutEmployee
  async function fetchUsersWithoutEmployee() {
    try {
      const res = await api.get('/users/available');

      if (res.success) {
        usersWithoutEmp.value = res.data;
      }
      return res;
    } catch (error) {
      console.error('Fetch users without employee error:', error);
      return {
        success: false,
        message: error?.message || 'Lỗi tải danh sách người dùng',
      };
    }
  }

  // Create Users
  async function createUser(data) {
    try {
      const res = await api.post('/users', data);
      return res;
    } catch (error) {
      console.error('Create user error:', error);
      return {
        success: false,
        message: error?.message || 'Lỗi tạo người dùng',
      };
    }
  }

  // Update Users
  async function updateUser(id, data) {
    try {
      const res = await api.put(`/users/${id}`, data);
      return res;
    } catch (error) {
      console.error('Update user error:', error);
      return {
        success: false,
        message: error?.message || 'Lỗi cập nhật người dùng',
      };
    }
  }

  // Delete Users
  async function deleteUser(id) {
    try {
      const res = await api.delete(`/users/${id}`);
      return res;
    } catch (error) {
      console.error('Delete user error:', error);
      return {
        success: false,
        message: error?.message || 'Lỗi xoá người dùng',
      };
    }
  }

  return {
    users,
    pagination,
    loading,
    usersWithoutEmp,
    fetchUser,
    createUser,
    updateUser,
    deleteUser,
    fetchUsersWithoutEmployee,
  };
});
