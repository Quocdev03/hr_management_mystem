import api from "@/api";
import { defineStore } from "pinia";
import { ref } from "vue";

export const useAuthStore = defineStore("auth", () => {
	const accessToken = ref(localStorage.getItem("access_token") || null);
	const user = ref(null);

	async function login(email, password) {
		const res = await api.post("/auth/login", { email, password });
		if (res.success) {
			const token = res.data.access_token;
			accessToken.value = token;
			user.value = res.data.user;
			localStorage.setItem("access_token", token);
		}
		return res;
	}

	function logout() {
		accessToken.value = null;
		user.value = null;
		localStorage.removeItem("access_token");
		window.location.href = "/login";
	}

	return { accessToken, user, login, logout };
});
