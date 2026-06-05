import api from '@/api';
import { defineStore } from 'pinia';
import { ref, computed } from 'vue';

export const useAuthStore = defineStore('auth', () => {
  // State
  const accessToken = ref(localStorage.getItem('access_token') || null);
  const refreshToken = ref(localStorage.getItem('refresh_token') || null);

  // Safely parse localStorage user to avoid throwing on invalid JSON
  let savedUser = null;
  try {
    savedUser = JSON.parse(localStorage.getItem('user') || 'null');
  } catch (e) {
    savedUser = null;
  }

  const user = ref(savedUser);
  const userProfile = ref(null);
  const loading = ref(false);

  // Getter
  const isAuthenticated = computed(() => !!accessToken.value);

  // Login
  async function login(email, password) {
    loading.value = true;

    try {
      const res = await api.post('/auth/login', { email, password });

      if (res.success) {
        const { data } = res;

        accessToken.value = data.access_token;
        refreshToken.value = data.refresh_token;
        user.value = data.user;

        localStorage.setItem('access_token', data.access_token);
        localStorage.setItem('refresh_token', data.refresh_token);
        localStorage.setItem('user', JSON.stringify(data.user));
      }

      return res;
    } catch (error) {
      console.error('Login error:', error);
      return {
        success: false,
        message: error?.message || 'Đăng nhập thất bại',
      };
    } finally {
      loading.value = false;
    }
  }

  // Profile
  async function profile() {
    loading.value = true;
    try {
      const res = await api.get('/auth/me');
      if (res.success) {
        userProfile.value = res.data;
        return res;
      }
    } catch (error) {
      console.error('Profile error:', error);
      return {
        success: false,
        message: error?.message || 'Lấy thông tin thất bại',
      };
    } finally {
      loading.value = false;
    }
  }

  // Logout
  async function logout() {
    if (accessToken.value) {
      try {
        await api.post('/auth/logout', {
          refresh_token: refreshToken.value,
        });
      } catch (error) {
        console.error('Lỗi khi gọi API logout:', error);
      }
    }

    accessToken.value = null;
    refreshToken.value = null;
    user.value = null;

    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('user');

    window.location.href = '/login';
  }

  // Refresh Token
  async function refresh() {
    if (!refreshToken.value) {
      return { success: false, message: 'Không tìm thấy refresh token' };
    }
    try {
      const res = await api.post('/auth/refresh', {
        refresh_token: refreshToken.value,
      });

      if (res.success) {
        const { data } = res;
        accessToken.value = data.access_token;
        refreshToken.value = data.refresh_token;
        user.value = data.user;

        localStorage.setItem('access_token', data.access_token);
        localStorage.setItem('refresh_token', data.refresh_token);
        localStorage.setItem('user', JSON.stringify(data.user));
      } else {
        await logout();
      }

      return res;
    } catch (error) {
      console.error('Refresh token error:', error);
      await logout();
      return {
        success: false,
        message: error?.message || 'Làm mới token thất bại',
      };
    }
  }

  return {
    // state
    accessToken,
    refreshToken,
    user,
    userProfile,
    loading,

    // getter
    isAuthenticated,

    // actions
    login,
    profile,
    logout,
    refresh,
  };
});
