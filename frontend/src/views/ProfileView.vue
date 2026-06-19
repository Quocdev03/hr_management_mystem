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
						<div class="profile-stat-item">
							<Skeleton type="text" width="90px" height="16px" />
							<Skeleton type="text" width="80px" height="16px" />
						</div>
						<div class="profile-stat-item">
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
						<div class="profile-stat-item">
							<span class="profile-stat-label">Ngày tham gia:</span>
							<span class="profile-stat-value">{{
								formatDate(
									authStore?.userProfile?.employee?.join_date,
								)
							}}</span>
						</div>
						<div class="profile-stat-item">
							<span class="profile-stat-label">Phòng ban:</span>
							<span class="profile-stat-value">{{
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
											?.position?.name || "N/A"
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
	grid-template-columns: 400px 1fr;
	gap: var(--space-3);
}

.profile-card {
	position: relative;
	background: var(--bg-card);
	border-radius: var(--radius-lg);
	border: 1px solid rgba(66, 97, 237, 0.08);
	box-shadow: 0 10px 30px rgba(66, 97, 237, 0.03), 0 1px 3px rgba(66, 97, 237, 0.01);
	transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
	overflow: hidden;
}

.profile-card::before {
	content: "";
	position: absolute;
	top: 0;
	left: 0;
	right: 0;
	height: 4px;
	background: var(--primary-gradient);
	opacity: 0.9;
	z-index: 2;
}

.profile-card:hover {
	transform: translateY(-2px);
	box-shadow: 0 20px 40px rgba(66, 97, 237, 0.08);
	border-color: rgba(66, 97, 237, 0.15);
}

.sidebar-card {
	padding: var(--space-4) var(--space-3);
	display: flex;
	flex-direction: column;
	align-items: center;
	height: fit-content;
}

.avatar-section {
	text-align: center;
	margin-bottom: var(--space-2);
	width: 100%;
}

.avatar-container {
	width: 130px;
	height: 130px;
	margin: var(--space-3) auto;
}

.profile-avatar-circle-large {
	width: 130px;
	height: 130px;
	border-radius: var(--radius-xl);
	background: var(--primary-gradient);
	color: white;
	display: flex;
	align-items: center;
	justify-content: center;
	font-weight: 800;
	font-size: 4.25rem;
	box-shadow: 0 8px 24px rgba(66, 97, 237, 0.2), inset 0 4px 10px rgba(255, 255, 255, 0.25);
	margin: var(--space-3) auto;
	transition: transform 0.4s cubic-bezier(0.34, 1.56, 0.64, 1), box-shadow 0.4s ease;
	border: 4px solid rgba(255, 255, 255, 0.95);
}

.sidebar-card:hover .profile-avatar-circle-large {
	transform: scale(1.05) rotate(-3deg);
	box-shadow: 0 15px 35px rgba(66, 97, 237, 0.35), inset 0 4px 10px rgba(255, 255, 255, 0.3);
}

.user-name {
	font-family: var(--font-title);
	font-size: var(--fs-xl);
	font-weight: 800;
	color: var(--text-main);
	margin-bottom: var(--space-1);
	letter-spacing: -0.01em;
}

.user-role-tag {
	display: inline-block;
	padding: 4px 12px;
	background: linear-gradient(135deg, rgba(66, 97, 237, 0.08) 0%, rgba(0, 192, 250, 0.03) 100%);
	color: var(--primary-color);
	border: 1px solid rgba(66, 97, 237, 0.12);
	border-radius: var(--radius-full);
	font-size: var(--fs-xs);
	font-weight: var(--fw-bold);
	letter-spacing: 0.05em;
	margin-bottom: var(--space-2);
}

.status-indicator {
	display: flex;
	align-items: center;
	justify-content: center;
	gap: 6px;
	font-size: var(--fs-xs);
	color: var(--text-muted);
	background: rgba(66, 97, 237, 0.03);
	padding: 4px 12px;
	border-radius: var(--radius-full);
	border: 1px solid rgba(66, 97, 237, 0.05);
	width: fit-content;
	margin: 0 auto var(--space-1) auto;
	font-weight: var(--fw-semibold);
}

.status-dot {
	width: 8px;
	height: 8px;
	border-radius: 50%;
}
.status-dot.active {
	background: var(--success-color);
	box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.2);
}
.status-dot.inactive {
	background: var(--text-light);
}

.quick-stats {
	width: 100%;
	border-top: 1px solid rgba(66, 97, 237, 0.08);
	padding-top: var(--space-3);
	margin-top: var(--space-2);
	display: flex;
	flex-direction: column;
	gap: 6px;
}

.profile-stat-item {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 8px 12px;
	border-radius: var(--radius-md);
	background: rgba(66, 97, 237, 0.02);
	border: 1px solid transparent;
	transition: all 0.2s ease;
}

.profile-stat-item:hover {
	background: rgba(66, 97, 237, 0.05);
	border-color: rgba(66, 97, 237, 0.08);
}

.profile-stat-label {
	font-size: var(--fs-xs);
	color: var(--text-muted);
	font-weight: var(--fw-semibold);
}

.profile-stat-value {
	font-size: var(--fs-xs);
	color: var(--text-main);
	font-weight: var(--fw-bold);
	text-align: right;
}

.main-card {
	padding: var(--space-4) var(--space-4);
}

.card-header {
	display: flex;
	align-items: center;
	gap: 12px;
	margin-bottom: var(--space-3);
	border-bottom: 1px solid rgba(66, 97, 237, 0.08);
	padding-bottom: var(--space-2);
}

.card-header h3 {
	font-family: var(--font-title);
	font-size: var(--fs-lg);
	font-weight: 800;
	color: var(--text-main);
	margin: 0;
	letter-spacing: -0.01em;
}

.card-icon {
	width: 18px;
	height: 18px;
	color: var(--primary-color);
	background: rgba(66, 97, 237, 0.08);
	padding: var(--space-1);
	border-radius: var(--radius-sm);
	box-sizing: content-box;
	border: 1px solid rgba(66, 97, 237, 0.12);
}

.info-grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(270px, 1fr));
	gap: var(--space-3);
}

.info-group label {
	display: block;
	font-size: 11px;
	font-weight: 800;
	color: var(--text-light);
	text-transform: uppercase;
	letter-spacing: 0.07em;
	margin-bottom: 6px;
}

.info-group .value {
	font-size: var(--fs-sm);
	font-weight: var(--fw-semibold);
	color: var(--text-main);
	padding: 10px 16px;
	background: linear-gradient(180deg, rgba(255, 255, 255, 0.85) 0%, rgba(244, 246, 255, 0.45) 100%);
	border-radius: var(--radius-md);
	border: 1px solid rgba(66, 97, 237, 0.06);
	box-shadow: inset 0 1px 2px rgba(66, 97, 237, 0.01);
	transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.info-group:hover .value {
	border-color: rgba(66, 97, 237, 0.2);
	background: linear-gradient(180deg, rgba(255, 255, 255, 1) 0%, rgba(244, 246, 255, 0.65) 100%);
	box-shadow: 0 4px 12px rgba(66, 97, 237, 0.04);
}

.info-group .value.salary {
	color: var(--success-color);
	font-family: var(--font-widget);
	font-weight: 800;
	background: linear-gradient(180deg, rgba(255, 255, 255, 0.85) 0%, rgba(16, 185, 129, 0.03) 100%);
	border-color: rgba(16, 185, 129, 0.12);
}

.info-group:hover .value.salary {
	border-color: rgba(16, 185, 129, 0.35);
	background: linear-gradient(180deg, rgba(255, 255, 255, 1) 0%, rgba(16, 185, 129, 0.06) 100%);
	box-shadow: 0 4px 12px rgba(16, 185, 129, 0.06);
}

.mt-6 {
	margin-top: var(--space-3);
}

@media (max-width: 1023px) {
	.profile-grid {
		grid-template-columns: 1fr;
	}

	.profile-stat-item {
		display: flex;
		justify-content: flex-start;
		padding: var(--space-2) var(--space-1);
		gap: var(--space-2);
	}
}

@media (max-width: 640px) {
	.avatar-section {
		flex-direction: column;
		text-align: center;
		gap: var(--space-2);
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
