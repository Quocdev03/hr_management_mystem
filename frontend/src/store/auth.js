import api from "@/api";
import { defineStore } from "pinia";
import { ref, computed } from "vue";

export const useAuthStore = defineStore("auth", () => {
	// ===== State =====
	const accessToken = ref(localStorage.getItem("access_token") || null);
	const user = ref(JSON.parse(localStorage.getItem("user")) || null);
	const userProfile = ref(null);
	const loading = ref(false);

	// ===== Getter =====
	const isAuthenticated = computed(() => !!accessToken.value);

	// ===== Login =====
	async function login(email, password) {
		loading.value = true;

		try {
			const res = await api.post("/auth/login", {
				email,
				password,
			});

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
				message: "Đăng nhập thất bại",
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
			console.log(error);
		} finally {
			loading.value = false;
		}
	}

	// ===== Logout =====
	function logout() {
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
