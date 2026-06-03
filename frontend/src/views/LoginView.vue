<script setup>
	import mailIcon from "@/assets/svg/mail.svg";
	import lockIcon from "@/assets/svg/lock.svg";
	import eyeIcon from "@/assets/svg/eye.svg";
	import eyeOffIcon from "@/assets/svg/eye-off.svg";
	import usersIcon from "@/assets/svg/user.svg";
	import { reactive, ref } from "vue";
	import { useAuthStore } from "@/store/auth";
	import { useToast } from "vue-toastification";
	import { storeToRefs } from "pinia";
	import { useRouter } from "vue-router";

	const toast = useToast();
	const router = useRouter();
	const authStore = useAuthStore();
	const { loading } = storeToRefs(authStore);
	const isPasswordVisible = ref(false);
	const credentials = reactive({ email: "", password: "" });

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
	<!-- ===== Container Chính ===== -->
	<main class="login-container">
		<!-- ===== Thẻ Đăng Nhập ===== -->
		<section class="login-card">
			<!-- ===== Phần Tiêu Đề ===== -->
			<header class="logo-section">
				<div class="logo-icon">
					<img :src="usersIcon" alt="logo" class="logo-img" />
				</div>
				<h2>Welcome Back</h2>
				<p>Đăng nhập vào Hệ thống Quản lý Nhân sự</p>
			</header>

			<!-- ===== Biểu Mẫu Đăng Nhập ===== -->
			<form @submit.prevent="loginHandler" class="login-form">
				<!-- Nhập Email -->
				<div class="login-input-group">
					<label for="email">Địa chỉ Email</label>
					<div class="input-wrapper">
						<img :src="mailIcon" class="input-icon" />
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
						<img :src="lockIcon" class="input-icon" />
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
						>
							<img :src="!isPasswordVisible ? eyeIcon : eyeOffIcon" alt="toggle" />
						</button>
					</div>
				</div>

				<!-- Hành Động Phụ -->
				<div class="form-actions">
					<a href="#" class="forgot-password">Quên mật khẩu?</a>
				</div>

				<!-- Nút Gửi -->
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
		background: var(--bg-light);
		padding: var(--space-2);
	}

	.login-card {
		width: 100%;
		max-width: 450px;
		background: var(--bg-card);
		border: 1px solid var(--border-color);
		border-radius: var(--radius-lg);
		padding: var(--space-5);
		box-shadow: var(--shadow-md);
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
		width: 60px;
		height: 60px;
		background-color: #eff6ff;
		color: var(--primary-color);
		border-radius: var(--radius-md);
	}

	.logo-img {
		width: 32px;
		height: 32px;
		filter: invert(44%) sepia(87%) saturate(2258%) hue-rotate(200deg) brightness(101%)
			contrast(92%); /* #3b82f6 */
	}

	.logo-section h2 {
		color: var(--text-main);
		font-size: var(--fs-2xl);
		margin-bottom: var(--space-1);
		letter-spacing: var(--tracking-tight);
	}

	.logo-section p {
		color: var(--text-muted);
		font-size: var(--fs-sm);
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
		filter: invert(72%) sepia(10%) saturate(415%) hue-rotate(182deg) brightness(88%)
			contrast(89%); /* #94a3b8 */
		transition: all 0.2s ease;
	}

	.input-wrapper input {
		width: 100%;
		padding: 0.75rem 1rem 0.75rem 2.75rem;
		border: 1px solid var(--border-color);
		border-radius: var(--radius-sm);
		font-size: var(--fs-base);
		color: var(--text-main);
		background: white;
		transition: all 0.2s ease;
		outline: none;
	}

	.input-wrapper input::placeholder {
		color: var(--text-light);
	}

	.input-wrapper input:focus {
		border-color: var(--primary-color);
		box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
	}

	.input-wrapper input:focus ~ .input-icon {
		filter: invert(44%) sepia(87%) saturate(2258%) hue-rotate(200deg) brightness(101%)
			contrast(92%); /* #3b82f6 */
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
	}

	.toggle-password:hover {
		color: var(--text-muted);
	}

	.toggle-password img {
		width: 18px;
		height: 18px;
		filter: invert(72%) sepia(10%) saturate(415%) hue-rotate(182deg) brightness(88%)
			contrast(89%);
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
		background-color: var(--primary-color);
		color: white;
		border: none;
		border-radius: var(--radius-sm);
		padding: 0.75rem 1rem;
		font-size: var(--fs-base);
		font-weight: var(--fw-semibold);
		cursor: pointer;
		transition: all 0.2s ease;
		display: flex;
		justify-content: center;
		align-items: center;
		min-height: 48px;
		margin-top: var(--space-1);
	}

	.submit-btn:hover {
		background-color: var(--primary-hover);
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(59, 130, 246, 0.2);
	}

	.submit-btn:active {
		background-color: #1d4ed8;
		transform: translateY(0);
	}

	.submit-btn:disabled {
		opacity: 0.7;
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

	@media (max-width: 480px) {
		.login-card {
			padding: var(--space-4) var(--space-3);
			border: none;
			border-radius: 0;
			box-shadow: none;
			background: transparent;
			backdrop-filter: none;
		}
	}
</style>
