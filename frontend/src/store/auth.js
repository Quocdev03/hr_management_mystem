import api from "@/api";
import { defineStore } from "pinia";
import { ref, computed } from "vue";

export const useAuthStore = defineStore("auth", () => {
	// ===== State =====
	const accessToken = ref(localStorage.getItem("access_token") || null);
	// Safely parse localStorage user to avoid throwing on invalid JSON
	let savedUser = null;
	try {
		savedUser = JSON.parse(localStorage.getItem("user") || "null");
	} catch (e) {
		savedUser = null;
	}
	const user = ref(savedUser);
	const userProfile = ref(null);
	const loading = ref(false);

	// ===== Getter =====
	const isAuthenticated = computed(() => !!accessToken.value);

	// ===== Login =====
	async function login(email, password) {
		loading.value = true;

		try {
			const res = await api.post("/auth/login", { email, password });

			if (res.success) {
				const data = res.data;

				accessToken.value = data.access_token;
				user.value = data.user;

				localStorage.setItem("access_token", data.access_token);

				localStorage.setItem("user", JSON.stringify(data.user));
			}

			return res;
		} catch (error) {
			console.error("Login error:", error);
			return {
				success: false,
				message: error?.message || "Đăng nhập thất bại",
			};
		} finally {
			loading.value = false;
		}
	}

	// ===== Profile =====
	async function profile() {
		loading.value = true;
		try {
			const res = await api.get("/auth/profile");
			if (res.success) {
				const data = res.data;
				if (res.success) {
					userProfile.value = data;
				}
				return res;
			}
		} catch (error) {
			console.error("Profile error:", error);
			return {
				success: false,
				message: error?.message || "Lấy thông tin thất bại",
			};
		} finally {
			loading.value = false;
		}
	}

	// ===== Logout =====
	async function logout() {
		if (accessToken.value) {
			try {
				await api.post("/auth/logout");
			} catch (error) {
				console.error("Lỗi khi gọi API logout:", error);
			}
		}

		accessToken.value = null;
		user.value = null;

		localStorage.removeItem("access_token");
		localStorage.removeItem("user");

		window.location.href = "/login";
	}

	return {
		// state
		accessToken,
		user,
		userProfile,
		loading,

		// getter
		isAuthenticated,

		// actions
		login,
		profile,
		logout,
	};
});
