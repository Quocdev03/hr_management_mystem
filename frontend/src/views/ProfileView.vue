<script setup>
// --- Tài nguyên và Biểu tượng ---
import userBlueIcon from "@/assets/svg/user-blue.svg";
import briefcaseGreenIcon from "@/assets/svg/briefcase-green.svg";
import { useAuthStore } from "@/store/auth";
import { onMounted, computed } from "vue";
import { formatCurrency, formatDate, formatStatus } from "@/helpers/formatters";
const authStore = useAuthStore();

onMounted(() => {
	authStore.profile();
});
</script>
<template>
	<div class="profile-view">
		<!-- ===== Tiêu đề trang ===== -->
		<div class="page-header">
			<div class="header-content">
				<h1 class="page-title">Quản lý nhân sự</h1>
				<p class="page-subtitle">
					Quản lý và xem thông tin chi tiết của bạn
				</p>
			</div>
		</div>

		<!-- ===== Bố cục hồ sơ ===== -->
		<div class="profile-grid">
			<!-- ===== Cột trái: Thông tin tổng quan ===== -->
			<div class="profile-card sidebar-card">
				<div class="avatar-section">
					<div class="avatar-container">
						<img
							src="https://ui-avatars.com/api/?name=User&background=random"
							class="profile-avatar"
							v-if="!authStore.userProfile?.employee"
						/>
						<img
							:src="
								'https://ui-avatars.com/api/?name=' +
								authStore.userProfile?.employee?.first_name +
								'+' +
								authStore.userProfile?.employee?.last_name +
								'&background=random'
							"
							class="profile-avatar"
							v-else
						/>
					</div>
					<h2 class="user-name">
						{{ authStore.userProfile?.employee?.first_name }}
						{{ authStore.userProfile?.employee?.last_name }}
					</h2>
					<p class="user-role-tag">
						{{ authStore?.userProfile?.role_name?.toUpperCase() }}
					</p>
					<div class="status-indicator">
						<span
							:class="[
								'status-dot',
								authStore?.userProfile?.employee?.status ===
								'active'
									? 'active'
									: 'inactive',
							]"
						></span>
						<span>{{
							authStore?.userProfile?.employee?.status ===
							"active"
								? "Đang làm việc"
								: "Đã nghỉ việc"
						}}</span>
					</div>
				</div>

				<div class="quick-stats">
					<div class="stat-item">
						<span class="stat-label">Ngày tham gia:</span>
						<span class="stat-value">{{
							formatDate(
								authStore?.userProfile?.employee?.join_date,
							)
						}}</span>
					</div>
					<div class="stat-item">
						<span class="stat-label">Phòng ban:</span>
						<span class="stat-value">{{
							authStore?.userProfile?.employee?.department?.name
						}}</span>
					</div>
				</div>
			</div>

			<!-- ===== Cột phải: Thông tin chi tiết ===== -->
			<div class="profile-main">
				<!-- ===== Block: Thông tin cá nhân ===== -->
				<div class="profile-card main-card">
					<div class="card-header">
						<img :src="userBlueIcon" alt="user" class="card-icon" />
						<h3>Thông tin cá nhân</h3>
					</div>
					<div class="info-grid">
						<div class="info-group">
							<label>Họ và tên</label>
							<div class="value">
								{{
									authStore.userProfile?.employee?.first_name
								}}
								{{ authStore.userProfile?.employee?.last_name }}
							</div>
						</div>
						<div class="info-group">
							<label>Giới tính</label>
							<div class="value">
								{{
									authStore.userProfile?.employee?.gender ===
									"male"
										? "Nam"
										: authStore.userProfile?.employee
													?.gender === "female"
											? "Nữ"
											: authStore.userProfile?.employee
														?.gender === "other"
												? "Khác"
												: "N/A"
								}}
							</div>
						</div>
						<div class="info-group">
							<label>Ngày sinh</label>
							<div class="value">
								{{
									formatDate(
										authStore.userProfile?.employee
											?.birth_date,
									)
								}}
							</div>
						</div>
						<div class="info-group">
							<label>Số điện thoại</label>
							<div class="value">
								{{
									authStore.userProfile?.employee?.phone ||
									"N/A"
								}}
							</div>
						</div>
						<div class="info-group">
							<label>Email cá nhân</label>
							<div class="value">
								{{
									authStore.userProfile?.employee?.email ||
									"N/A"
								}}
							</div>
						</div>
					</div>
				</div>

				<!-- ===== Block: Thông tin công việc ===== -->
				<div class="profile-card main-card mt-6">
					<div class="card-header">
						<img
							:src="briefcaseGreenIcon"
							alt="work"
							class="card-icon"
						/>
						<h3>Thông tin công việc</h3>
					</div>
					<div class="info-grid">
						<div class="info-group">
							<label>Chức vụ</label>
							<div class="value">
								{{
									authStore.userProfile?.employee?.position ||
									"N/A"
								}}
							</div>
						</div>
						<div class="info-group">
							<label>Phòng ban</label>
							<div class="value">
								{{
									authStore.userProfile?.employee?.department
										?.name || "N/A"
								}}
							</div>
						</div>
						<div class="info-group">
							<label>Ngày vào làm</label>
							<div class="value">
								{{
									formatDate(
										authStore.userProfile?.employee
											?.join_date,
									)
								}}
							</div>
						</div>
						<div class="info-group">
							<label>Mức lương hiện tại</label>
							<div class="value salary">
								{{
									formatCurrency(
										authStore.userProfile?.employee?.salary,
									)
								}}
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>
<style scoped>
/* ===== Bố cục chính ===== */
.profile-grid {
	display: grid;
	grid-template-columns: 340px 1fr;
	gap: 1.5rem;
}

/* ===== Card dùng chung ===== */
.profile-card {
	background: var(--bg-card);
	border-radius: var(--border-radius-card);
	border: 1px solid var(--border-color);
	box-shadow: 0 4px 20px -2px rgba(0, 0, 0, 0.05);
	overflow: hidden;
}
.page-header {
	margin-bottom: var(--space-4);
}
.page-title {
	font-size: var(--fs-2xl);
	font-weight: var(--fw-bold);
	letter-spacing: var(--tracking-tight);
	margin: 0 0 4px 0;
	color: var(--text-main);
}

.page-subtitle {
	color: var(--text-muted);
	font-size: var(--fs-sm);
	margin: 0;
}

/* ===== Cột trái: Sidebar ===== */
.sidebar-card {
	padding: 2.5rem 1.5rem;
	display: flex;
	flex-direction: column;
	align-items: center;
	height: fit-content;
}

.avatar-section {
	text-align: center;
	margin-bottom: 2rem;
}

.avatar-container {
	width: 140px;
	height: 140px;
	margin: 1.5rem auto;
}

.profile-avatar {
	width: 100%;
	height: 100%;
	border-radius: 40px;
	object-fit: cover;
	box-shadow: 0 10px 25px -5px rgba(59, 130, 246, 0.2);
}

.user-name {
	font-size: 1.5rem;
	font-weight: 700;
	color: var(--text-main);
	margin-bottom: 0.6rem;
}

.user-role-tag {
	display: inline-block;
	padding: 0.35rem 1rem;
	background: var(--bg-light);
	color: var(--text-muted);
	border-radius: 999px;
	font-size: 0.8125rem;
	font-weight: 600;
	margin-bottom: 0.75rem;
}

.status-indicator {
	display: flex;
	align-items: center;
	justify-content: center;
	gap: 8px;
	font-size: 0.875rem;
	color: var(--text-muted);
}

.status-dot {
	width: 8px;
	height: 8px;
	border-radius: 50%;
}
.status-dot.active {
	background: var(--success-color);
	box-shadow: 0 0 0 4px rgba(16, 185, 129, 0.1);
}
.status-dot.inactive {
	background: var(--text-muted);
}

.quick-stats {
	width: 100%;
	border-top: 1px solid var(--bg-light);
	padding-top: 1.5rem;
	margin-top: 0.5rem;
}

.stat-item {
	display: flex;
	justify-content: space-between;
	padding: 0.75rem 0;
}

.stat-label {
	font-size: 0.875rem;
	color: var(--text-muted);
}

.stat-value {
	font-size: 0.875rem;
	font-weight: 600;
	color: var(--text-main);
}

/* ===== Cột phải: Thông tin chi tiết ===== */
.main-card {
	padding: 1.75rem 2rem;
}

.card-header {
	display: flex;
	align-items: center;
	gap: 12px;
	margin-bottom: 1.5rem;
}

.card-header h3 {
	font-size: 1.125rem;
	font-weight: 700;
	color: var(--text-main);
}

.card-icon {
	width: 20px;
	height: 20px;
}

.info-grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(270px, 1fr));
	gap: 1.5rem;
}

.info-group label {
	display: block;
	font-size: 0.75rem;
	font-weight: 700;
	color: var(--text-light);
	text-transform: uppercase;
	letter-spacing: 0.05em;
	margin-bottom: 0.5rem;
}

.info-group .value {
	font-size: 1rem;
	font-weight: 500;
	color: var(--text-main);
	padding: 0.75rem 1rem;
	background: var(--bg-lighter);
	border-radius: 12px;
	border: 1px solid var(--bg-light);
}

.info-group .value.salary {
	color: var(--success-color);
	font-weight: 700;
}

/* ===== Utilities ===== */
.mt-6 {
	margin-top: 1.5rem;
}

/* ===== Responsive ===== */
@media (max-width: 1023px) {
	.profile-grid {
		grid-template-columns: 1fr;
	}

	.stat-item {
		display: flex;
		justify-content: flex-start;
		padding: 0.75rem 0;
		gap: 1rem;
	}
}

@media (max-width: 640px) {
	.avatar-section {
		flex-direction: column;
		text-align: center;
		gap: 1rem;
	}
	.info-grid {
		grid-template-columns: 1fr;
	}
	.quick-stats {
		grid-template-columns: 1fr;
		gap: 0;
	}
}
</style>
