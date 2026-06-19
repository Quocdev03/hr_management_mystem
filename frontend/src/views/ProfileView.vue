<script setup>
// ─── Store & tiện ích ────────────────────────────────────────────────────────
import { useAuthStore } from "@/store/auth";
import { onMounted, computed } from "vue";
import { storeToRefs } from "pinia";
import Skeleton from "@/components/Skeleton.vue";
import { formatCurrency, formatDate } from "@/helpers/formatters";

// ─── Icon Lucide ─────────────────────────────────────────────────────────────
import { User, Briefcase } from "@lucide/vue";

// ─── Khởi tạo ────────────────────────────────────────────────────────────────

const authStore = useAuthStore();
const { loading } = storeToRefs(authStore);

// ─── Computed ─────────────────────────────────────────────────────────────────

const hasEmployee = computed(() => !!authStore.userProfile?.employee);

const displayName = computed(() => {
	const first = authStore.userProfile?.employee?.first_name || "";
	const last = authStore.userProfile?.employee?.last_name || "";

	if (first || last) {
		return `${first} ${last}`.trim();
	}

	return (
		authStore.userProfile?.email || authStore.user?.email || "Người dùng"
	);
});

onMounted(() => {
	authStore.profile();
});
</script>

<template>
	<div class="profile-view">
		<div class="page-header">
			<div class="header-content">
				<h1 class="page-title">Thông tin tài khoản</h1>
				<p class="page-subtitle">
					Quản lý và xem thông tin chi tiết cá nhân
				</p>
			</div>
		</div>

		<div class="profile-grid">
			<!-- Cột trái: Thông tin tổng quan -->
			<div class="profile-card sidebar-card">
				<template v-if="loading">
					<div class="avatar-section" style="width: 100%">
						<div class="avatar-container">
							<Skeleton
								type="avatar"
								width="140px"
								height="140px"
								border-radius="40px"
							/>
						</div>
						<Skeleton
							type="text"
							width="150px"
							height="24px"
							style="margin: 0 auto 0.6rem auto; display: block"
						/>
						<Skeleton
							type="text"
							width="100px"
							height="24px"
							border-radius="999px"
							style="
								margin-bottom: 0.75rem;
								display: inline-block;
							"
						/>
						<Skeleton
							type="text"
							width="120px"
							height="16px"
							style="margin: 0 auto; display: block"
						/>
					</div>
					<div class="quick-stats">
						<div class="stat-item">
							<Skeleton type="text" width="90px" height="16px" />
							<Skeleton type="text" width="80px" height="16px" />
						</div>
						<div class="stat-item">
							<Skeleton type="text" width="80px" height="16px" />
							<Skeleton type="text" width="100px" height="16px" />
						</div>
					</div>
				</template>
				<template v-else>
					<div class="avatar-section">
						<!-- Premium custom CSS initials avatar synchronized with sidebar -->
						<div class="profile-avatar-circle-large">
							{{ displayName ? displayName.charAt(0).toUpperCase() : 'U' }}
						</div>
						<h2 class="user-name">
							{{ displayName }}
						</h2>
						<p class="user-role-tag">
							{{
								authStore?.userProfile?.role_name?.toUpperCase()
							}}
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
								hasEmployee
									? authStore?.userProfile?.employee
											?.status === "active"
										? "Đang làm việc"
										: "Đã nghỉ việc"
									: "Chưa được gán nhân viên"
							}}</span>
						</div>
					</div>

					<div class="quick-stats" v-if="hasEmployee">
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
								authStore?.userProfile?.employee?.department
									?.name
							}}</span>
						</div>
					</div>
				</template>
			</div>

			<!-- Cột phải: Thông tin chi tiết -->
			<div class="profile-main">
				<template v-if="loading">
					<div class="profile-card main-card">
						<div class="card-header">
							<Skeleton
								type="text"
								width="20px"
								height="20px"
								border-radius="var(--radius-sm)"
							/>
							<Skeleton
								type="text"
								width="150px"
								height="20px"
								style="margin-bottom: 0"
							/>
						</div>
						<div class="info-grid">
							<div
								v-for="i in 5"
								:key="'skeleton-p-' + i"
								class="info-group"
							>
								<Skeleton
									type="text"
									width="80px"
									height="14px"
									style="margin-bottom: 0.5rem"
								/>
								<Skeleton
									type="text"
									width="100%"
									height="40px"
									border-radius="12px"
								/>
							</div>
						</div>
					</div>
				</template>
				<template v-else>
					<!-- Block: Thông tin cá nhân -->
					<div v-if="hasEmployee" class="profile-card main-card">
						<div class="card-header">
							<User class="card-icon" />
							<h3>Thông tin cá nhân</h3>
						</div>
						<div class="info-grid">
							<div class="info-group">
								<label>Họ và tên</label>
								<div class="value">
									{{
										authStore.userProfile?.employee
											?.first_name
									}}
									{{
										authStore.userProfile?.employee
											?.last_name
									}}
								</div>
							</div>
							<div class="info-group">
								<label>Giới tính</label>
								<div class="value">
									{{
										authStore.userProfile?.employee
											?.gender === "male"
											? "Nam"
											: authStore.userProfile?.employee
														?.gender === "female"
												? "Nữ"
												: authStore.userProfile
															?.employee
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
										authStore.userProfile?.employee
											?.phone || "N/A"
									}}
								</div>
							</div>
							<div class="info-group">
								<label>Email cá nhân</label>
								<div class="value">
									{{ authStore.user?.email || "N/A" }}
								</div>
							</div>
						</div>
					</div>

					<div v-else class="profile-card main-card">
						<div class="card-header">
							<User class="card-icon" />
							<h3>Thông tin tài khoản</h3>
						</div>
						<div class="info-grid">
							<div class="info-group">
								<label>Email</label>
								<div class="value">
									{{
										authStore.userProfile?.email ||
										authStore.user?.email ||
										"N/A"
									}}
								</div>
							</div>
							<div class="info-group">
								<label>Vai trò</label>
								<div class="value">
									{{
										authStore.userProfile?.role_name?.toUpperCase() ||
										"N/A"
									}}
								</div>
							</div>
							<div class="info-group">
								<label>Trạng thái</label>
								<div class="value">Chưa được gán nhân viên</div>
							</div>
							<div class="info-group">
								<label>Ghi chú</label>
								<div class="value">
									Tài khoản này chưa có hồ sơ nhân viên. Vui
									lòng liên hệ quản trị để cập nhật.
								</div>
							</div>
						</div>
					</div>

					<!-- Block: Thông tin công việc -->
					<div v-if="hasEmployee" class="profile-card main-card mt-6">
						<div class="card-header">
							<Briefcase class="card-icon" />
							<h3>Thông tin công việc</h3>
						</div>
						<div class="info-grid">
							<div class="info-group">
								<label>Chức vụ</label>
								<div class="value">
									{{
										authStore.userProfile?.employee
											?.position || "N/A"
									}}
								</div>
							</div>
							<div class="info-group">
								<label>Phòng ban</label>
								<div class="value">
									{{
										authStore.userProfile?.employee
											?.department?.name || "N/A"
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
											authStore.userProfile?.employee
												?.salary,
										)
									}}
								</div>
							</div>
						</div>
					</div>
				</template>
			</div>
		</div>
	</div>
</template>

<style scoped>
.profile-grid {
	display: grid;
	grid-template-columns: 340px 1fr;
	gap: var(--container-padding);
}

.profile-card {
	background: var(--bg-card);
	backdrop-filter: var(--glass-backdrop);
	-webkit-backdrop-filter: var(--glass-backdrop);
	border-radius: var(--radius-lg);
	border: var(--glass-border);
	box-shadow: var(--glass-shadow);
	overflow: hidden;
	transition: all 0.3s ease;
}

.profile-card:hover {
	box-shadow: var(--glass-shadow-hover);
}

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
	width: 100%;
}

.avatar-container {
	width: 140px;
	height: 140px;
	margin: 1.5rem auto;
}

.profile-avatar-circle-large {
	width: 140px;
	height: 140px;
	border-radius: var(--radius-xl);
	background: var(--primary-gradient);
	color: white;
	display: flex;
	align-items: center;
	justify-content: center;
	font-weight: var(--fw-bold);
	font-size: var(--fs-4xl);
	box-shadow: 0 10px 25px -5px rgba(66, 97, 237, 0.25);
	margin: 1.5rem auto;
}

.user-name {
	font-family: var(--font-title);
	font-size: var(--fs-2xl);
	font-weight: 700;
	color: var(--text-main);
	margin-bottom: 0.6rem;
}

.user-role-tag {
	display: inline-block;
	padding: 0.35rem 1rem;
	background: rgba(66, 97, 237, 0.08);
	color: var(--primary-color);
	border-radius: var(--radius-md);
	font-size: var(--fs-sm);
	font-weight: 600;
	margin-bottom: 0.75rem;
}

.status-indicator {
	display: flex;
	align-items: center;
	justify-content: center;
	gap: 8px;
	font-size: var(--fs-sm);
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
	border-top: 1px solid var(--border-color);
	padding-top: 1.5rem;
	margin-top: 0.5rem;
}

.stat-item {
	display: flex;
	justify-content: space-between;
	padding: 0.75rem 0;
}

.stat-label {
	font-size: var(--fs-sm);
	color: var(--text-muted);
}

.stat-value {
	font-size: var(--fs-sm);
	font-weight: 600;
	color: var(--text-main);
}

.main-card {
	padding: 1.75rem 2rem;
}

.card-header {
	display: flex;
	align-items: center;
	gap: 12px;
	margin-bottom: 1.5rem;
	border-bottom: 1px solid var(--border-color);
	padding-bottom: 12px;
}

.card-header h3 {
	font-family: var(--font-title);
	font-size: var(--fs-lg);
	font-weight: 700;
	color: var(--text-main);
	margin: 0;
}

.card-icon {
	width: 20px;
	height: 20px;
	color: var(--primary-color);
}

.info-grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(270px, 1fr));
	gap: 1.5rem;
}

.info-group label {
	display: block;
	font-size: var(--fs-xs);
	font-weight: 700;
	color: var(--text-muted);
	text-transform: uppercase;
	letter-spacing: 0.05em;
	margin-bottom: 0.5rem;
}

.info-group .value {
	font-size: var(--fs-base);
	font-weight: 500;
	color: var(--text-main);
	padding: 0.75rem 1rem;
	background: rgba(255, 255, 255, 0.45);
	border-radius: 12px;
	border: 1px solid rgba(66, 97, 237, 0.15);
}

.info-group .value.salary {
	color: var(--success-color);
	font-family: var(--font-widget);
	font-weight: 700;
}

.mt-6 {
	margin-top: 1.5rem;
}

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
