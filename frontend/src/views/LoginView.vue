<script setup>
import { Eye, EyeOff, Lock, Mail } from '@lucide/vue';

// ─── Icon SVG ────────────────────────────────────────────────────────────────

import logoPng from "@/assets/logo.png";

// ─── Store & tiện ích ────────────────────────────────────────────────────────
import { reactive, ref } from "vue";
import { useAuthStore } from "@/store/auth";
import { useToast } from "vue-toastification";
import { storeToRefs } from "pinia";
import { useRouter } from "vue-router";

// ─── Khởi tạo ────────────────────────────────────────────────────────────────
const toast = useToast();
const router = useRouter();
const authStore = useAuthStore();

const { loading } = storeToRefs(authStore);

const isPasswordVisible = ref(false);
const credentials = reactive({ email: "", password: "" });

// ─── Xử lý đăng nhập ─────────────────────────────────────────────────────────
async function loginHandler() {
	const res = await authStore.login(credentials.email, credentials.password);

	if (!res.success) {
		toast.error(res.message);
		return;
	}

	toast.success("Đăng nhập thành công!");
	router.push("/");
}

function togglePasswordVisibility() {
	isPasswordVisible.value = !isPasswordVisible.value;
}
</script>

<template>
	<main class="login-container">
		<!-- Glowing Background Blobs -->
		<div class="login-bg-blobs">
			<div class="blob blob-1"></div>
			<div class="blob blob-2"></div>
		</div>

		<section class="login-card">
			<header class="logo-section">
				<div class="logo-icon">
					<img :src="logoPng" alt="logo" class="logo-img" />
				</div>
				<h2 class="login-title">Chào mừng trở lại</h2>
				<p class="login-subtitle">
					Đăng nhập vào Hệ thống Quản lý Nhân sự
				</p>
			</header>

			<form @submit.prevent="loginHandler" class="login-form">
				<!-- Nhập Email -->
				<div class="login-input-group">
					<label for="email">Email</label>
					<div class="input-wrapper">
						<Mail class="input-icon" />
						<input
							type="email"
							id="email"
							v-model="credentials.email"
							placeholder="admin@company.com"
							required
						/>
					</div>
				</div>

				<!-- Nhập Mật Khẩu -->
				<div class="login-input-group">
					<label for="password">Mật khẩu</label>
					<div class="input-wrapper">
						<Lock class="input-icon" />
						<input
							:type="isPasswordVisible ? 'text' : 'password'"
							id="password"
							v-model="credentials.password"
							placeholder="••••••••"
							required
						/>
						<button
							type="button"
							class="toggle-password"
							@click="togglePasswordVisibility"
							aria-label="Ẩn/hiện mật khẩu"
						>
							<component :is="!isPasswordVisible ? Eye : EyeOff" />
						</button>
					</div>
				</div>

				<div class="form-actions">
					<a href="#" class="forgot-password">Quên mật khẩu?</a>
				</div>

				<button type="submit" class="submit-btn" :disabled="loading">
					<span v-if="!loading">Đăng Nhập</span>
					<span v-else class="loader"></span>
				</button>
			</form>
		</section>
	</main>
</template>

<style scoped>
.login-container {
	min-height: 100vh;
	display: flex;
	align-items: center;
	justify-content: center;
	background: var(--bg-gradient);
	padding: var(--space-2);
	position: relative;
	overflow: hidden;
}

/* Background Glowing Blobs */
.login-bg-blobs {
	position: absolute;
	inset: 0;
	z-index: 1;
	pointer-events: none;
}

.blob {
	position: absolute;
	border-radius: 50%;
	opacity: 0.6;
}

.blob-1 {
	background: radial-gradient(circle, rgba(0, 192, 250, 0.35) 0%, rgba(0, 192, 250, 0) 70%); /* Cyan */
	width: 450px;
	height: 450px;
	top: 10%;
	left: 15%;
	animation: floatBlob1 12s infinite alternate ease-in-out;
}

.blob-2 {
	background: radial-gradient(circle, rgba(103, 23, 204, 0.25) 0%, rgba(103, 23, 204, 0) 70%); /* Purple */
	width: 500px;
	height: 500px;
	bottom: 10%;
	right: 15%;
	animation: floatBlob2 10s infinite alternate ease-in-out;
}

@keyframes floatBlob1 {
	0% {
		transform: translate(0, 0) scale(1);
	}
	100% {
		transform: translate(30px, 40px) scale(1.1);
	}
}

@keyframes floatBlob2 {
	0% {
		transform: translate(0, 0) scale(1);
	}
	100% {
		transform: translate(-40px, -30px) scale(1.05);
	}
}

.login-card {
	width: 100%;
	max-width: 500px;
	background: rgba(255, 255, 255, 0.7);
	border: 1px solid rgba(255, 255, 255, 0.6);
	border-radius: var(--radius-xl);
	padding: var(--space-5);
	box-shadow: 0 20px 50px rgba(66, 97, 237, 0.15);
	z-index: 2;
	animation: fadeInCard 0.6s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes fadeInCard {
	from {
		opacity: 0;
		transform: translateY(20px);
	}
	to {
		opacity: 1;
		transform: translateY(0);
	}
}

.logo-section {
	text-align: center;
	margin-bottom: var(--space-4);
}

.logo-icon {
	margin-bottom: var(--space-2);
	display: inline-flex;
	align-items: center;
	justify-content: center;
	width: 72px;
	height: 72px;
	background-color: #ffffff;
	border-radius: var(--radius-xl);
	border: 1px solid rgba(255, 255, 255, 0.8);
	box-shadow: 0 8px 24px rgba(66, 97, 237, 0.12);
	overflow: hidden;
}

.logo-img {
	width: 100%;
	height: 100%;
	object-fit: cover;
}

.login-title {
	font-family: var(--font-title);
	color: var(--text-main);
	font-size: var(--fs-2xl);
	margin-top: 0;
	margin-bottom: var(--space-1);
	letter-spacing: var(--tracking-tight);
	font-weight: var(--fw-bold);
}

.login-subtitle {
	color: var(--text-muted);
	font-size: var(--fs-sm);
	margin: 0;
}

.login-form {
	display: flex;
	flex-direction: column;
	gap: var(--space-3);
}

.login-input-group {
	display: flex;
	flex-direction: column;
	gap: var(--space-1);
}

.login-input-group label {
	font-size: var(--fs-sm);
	font-weight: var(--fw-semibold);
	color: var(--text-muted);
}

.input-wrapper {
	position: relative;
	display: flex;
	align-items: center;
}

.input-icon {
	position: absolute;
	left: 1rem;
	width: 18px;
	height: 18px;
	color: var(--text-light);
	transition: color 0.2s ease;
	z-index: 1;
	pointer-events: none;
}

.input-wrapper input {
	width: 100%;
	padding: var(--space-2) var(--space-3) var(--space-2) 2.75rem;
	border: 1px solid var(--border-color);
	border-radius: var(--radius-md);
	font-size: var(--fs-base);
	color: var(--text-main);
	background: rgba(255, 255, 255, 0.85);
	transition: background-color 0.2s ease, border-color 0.2s ease, color 0.2s ease, box-shadow 0.2s ease;
	outline: none;
	font-family: var(--font-body);
}

.input-wrapper input::placeholder {
	color: var(--text-light);
}

.input-wrapper input:focus {
	background: white;
	border-color: var(--primary-color);
	box-shadow: 0 0 0 4px var(--primary-glow);
}

.input-wrapper input:focus ~ .input-icon {
	color: var(--primary-color);
}

.toggle-password {
	position: absolute;
	right: 0.75rem;
	background: none;
	border: none;
	color: var(--text-light);
	cursor: pointer;
	padding: 0.25rem;
	display: flex;
	align-items: center;
	justify-content: center;
	transition: color 0.2s ease;
	z-index: 10;
}

.toggle-password svg {
	width: 18px;
	height: 18px;
	color: var(--text-light);
	transition: color 0.2s ease;
}

.toggle-password:hover svg {
	color: var(--text-muted);
}

.form-actions {
	display: flex;
	align-items: center;
	justify-content: flex-end;
}

.forgot-password {
	font-size: var(--fs-sm);
	color: var(--primary-color);
	font-weight: var(--fw-medium);
	text-decoration: none;
	transition: color 0.2s ease;
}

.forgot-password:hover {
	color: var(--primary-hover);
	text-decoration: underline;
}

.submit-btn {
	background: var(--primary-gradient);
	color: white;
	border: none;
	border-radius: var(--radius-md);
	padding: var(--space-2) var(--space-3);
	font-size: var(--fs-base);
	font-weight: var(--fw-semibold);
	cursor: pointer;
	transition:
		background 0.2s cubic-bezier(0.4, 0, 0.2, 1),
		opacity 0.2s cubic-bezier(0.4, 0, 0.2, 1),
		transform 0.2s cubic-bezier(0.4, 0, 0.2, 1),
		box-shadow 0.2s cubic-bezier(0.4, 0, 0.2, 1);
	display: flex;
	justify-content: center;
	align-items: center;
	min-height: 48px;
	margin-top: var(--space-1);
	box-shadow: 0 4px 14px rgba(66, 97, 237, 0.25);
}

.submit-btn:hover:not(:disabled) {
	background: var(--primary-gradient-hover);
	transform: translateY(-1px);
	box-shadow: 0 6px 20px rgba(66, 97, 237, 0.35);
}

.submit-btn:active:not(:disabled) {
	transform: translateY(0);
}

.submit-btn:disabled {
	opacity: 0.6;
	cursor: not-allowed;
}

.loader {
	width: 20px;
	height: 20px;
	border: 2px solid rgba(255, 255, 255, 0.4);
	border-radius: var(--radius-full);
	border-top-color: white;
	animation: spin 0.8s linear infinite;
}

@keyframes spin {
	to {
		transform: rotate(360deg);
	}
}

@media (max-width: 768px) {
	.login-card {
		padding: var(--space-4) var(--space-4);
		max-width: 400px;
	}
}

@media (max-width: 480px) {
	.login-card {
		padding: var(--space-4) var(--space-3);
		border-radius: var(--radius-lg);
		background: rgba(255, 255, 255, 0.85);
		box-shadow: 0 10px 30px rgba(66, 97, 237, 0.12);
	}
}
</style>
