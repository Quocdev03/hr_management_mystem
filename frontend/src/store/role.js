import { defineStore } from 'pinia';
import api from '@/api';

export const useRoleStore = defineStore('role', {
    state: () => ({
        roles: [],
        availablePermissions: [],
        loading: false,
        error: null,
    }),
    actions: {
        async fetchRoles() {
            this.loading = true;
            try {
                const response = await api.get('/roles');
                if (response.success) {
                    this.roles = response.data;
                }
            } catch (err) {
                this.error = err?.message || 'Lỗi khi tải danh sách vai trò';
                throw err;
            } finally {
                this.loading = false;
            }
        },

        async fetchAvailablePermissions() {
            try {
                const response = await api.get('/users/permissions');
                if (response.success) {
                    this.availablePermissions = response.data;
                }
            } catch (err) {
                console.error("Lỗi khi tải danh sách quyền hạn", err);
            }
        },

        async createRole(roleData) {
            try {
                const response = await api.post('/roles', roleData);
                if (response.success) {
                    await this.fetchRoles();
                    return response;
                }
                throw new Error(response.message || "Lỗi tạo vai trò");
            } catch (err) {
                throw err;
            }
        },

        async updateRole(id, roleData) {
            try {
                const response = await api.patch(`/roles/${id}`, roleData);
                if (response.success) {
                    await this.fetchRoles();
                    return response;
                }
                throw new Error(response.message || "Lỗi cập nhật vai trò");
            } catch (err) {
                throw err;
            }
        },

        async deleteRole(id) {
            try {
                const response = await api.delete(`/roles/${id}`);
                if (response.success) {
                    await this.fetchRoles();
                    return response;
                }
                throw new Error(response.message || "Lỗi xoá vai trò");
            } catch (err) {
                throw err;
            }
        }
    }
});
