<script setup>
// ─── Icon SVG ────────────────────────────────────────────────────────────────
import usersIcon from "@/assets/svg/users.svg";
import checkIcon from "@/assets/svg/check-circle.svg";
import buildingIcon from "@/assets/svg/building.svg";

// ─── Store & tiện ích ────────────────────────────────────────────────────────
import { useDashboardStore } from "@/store/dashboard";
import { onMounted, computed } from "vue";
import { storeToRefs } from "pinia";
import { useAuthStore } from "@/store/auth";
import Skeleton from "@/components/Skeleton.vue";

// ─── Khởi tạo ────────────────────────────────────────────────────────────────
const dashboardStore = useDashboardStore();
const authStore = useAuthStore();
const { stats, loading } = storeToRefs(dashboardStore);

// Lấy thông tin user
let userFromStorage = null;
try {
	const raw = localStorage.getItem("user");
	userFromStorage = raw ? JSON.parse(raw) : null;
} catch (err) {
	userFromStorage = null;
}

const userName = computed(
	() => userFromStorage?.user_name ?? authStore.user?.user_name ?? "Thành viên",
);

// ─── Tải dữ liệu dashboard ───────────────────────────────────────────────────
async function loadDashboard() {
	try {
		await dashboardStore.fetchDashboard();
	} catch (err) {
		console.error("Lỗi khi tải dashboard:", err);
	}
}

onMounted(loadDashboard);
</script>

<template>
	<div class="dashboard-container">
		<header class="page-header">
			<div class="header-content">
				<h1 class="page-title">Tổng quan Dashboard</h1>
				<p class="page-subtitle">Báo cáo tóm tắt trạng thái nhân sự hệ thống</p>
			</div>
		</header>

		<!-- Bento Grid Top Row -->
		<div class="bento-grid">
			<!-- Welcome Card (Spans 2 columns on desktop) -->
			<div class="bento-card welcome-card">
				<div class="welcome-content">
					<span class="welcome-badge">Hệ thống HRM</span>
					<h1 class="welcome-title">Xin chào, {{ userName }}! 👋</h1>
					<p class="welcome-text">
						Chào mừng bạn quay trở lại với trang quản trị nhân sự. Dưới đây là tóm tắt nhanh về tình hình nhân sự của công ty hôm nay.
					</p>
				</div>
				<div class="welcome-illustration">
					<div class="blob blob-1"></div>
					<div class="blob blob-2"></div>
				</div>
			</div>

			<!-- Tổng số nhân viên (1 column) -->
			<div class="bento-card stat-card stat-card--teal">
				<div class="stat-card-inner">
					<div class="stat-header">
						<span class="stat-label">Tổng nhân viên</span>
						<div class="stat-icon-wrapper">
							<img :src="usersIcon" class="stat-icon" alt="users" />
						</div>
					</div>
					<div class="stat-body">
						<template v-if="loading">
							<Skeleton type="text" width="80px" height="40px" />
						</template>
						<template v-else>
							<div class="stat-value">{{ stats.total_employees }}</div>
						</template>
					</div>
					<div class="stat-footer">
						<span class="trend-text">Nhân sự toàn hệ thống</span>
					</div>
				</div>
			</div>

			<!-- Nhân viên đang làm việc (1 column) -->
			<div class="bento-card stat-card stat-card--emerald">
				<div class="stat-card-inner">
					<div class="stat-header">
						<span class="stat-label">Đang hoạt động</span>
						<div class="stat-icon-wrapper">
							<img :src="checkIcon" class="stat-icon" alt="active" />
						</div>
					</div>
					<div class="stat-body">
						<template v-if="loading">
							<Skeleton type="text" width="80px" height="40px" />
						</template>
						<template v-else>
							<div class="stat-value">{{ stats.total_employees_active }}</div>
						</template>
					</div>
					<div class="stat-footer">
						<span class="trend-badge trend-badge--success">
							{{ stats.total_employees ? Math.round((stats.total_employees_active / stats.total_employees) * 100) : 0 }}% hoạt động
						</span>
					</div>
				</div>
			</div>

			<!-- Tổng số phòng ban (1 column) -->
			<div class="bento-card stat-card stat-card--amber">
				<div class="stat-card-inner">
					<div class="stat-header">
						<span class="stat-label">Phòng ban</span>
						<div class="stat-icon-wrapper">
							<img :src="buildingIcon" class="stat-icon" alt="dept" />
						</div>
					</div>
					<div class="stat-body">
						<template v-if="loading">
							<Skeleton type="text" width="80px" height="40px" />
						</template>
						<template v-else>
							<div class="stat-value">{{ stats.total_departments }}</div>
						</template>
					</div>
					<div class="stat-footer">
						<span class="trend-text">Cơ cấu tổ chức hành chính</span>
					</div>
				</div>
			</div>
		</div>

		<!-- ===== Phân bổ nhân sự theo phòng ban ===== -->
		<section class="dept-section">
			<div class="section-header">
				<h2 class="section-title">Phân bổ nhân sự theo phòng ban</h2>
				<p class="section-subtitle">Tỉ lệ phần trăm nhân viên làm việc tại từng phòng ban</p>
			</div>
			
			<div class="dept-grid">
				<!-- Skeleton khi đang tải -->
				<template v-if="loading">
					<div v-for="i in 3" :key="'skeleton-d-' + i" class="bento-card dept-card">
						<div class="dept-header">
							<Skeleton type="text" width="120px" height="18px" />
							<Skeleton type="badge" width="40px" height="22px" />
						</div>
						<div class="progress-container" style="margin-top: 1rem;">
							<div class="skeleton" style="width: 100%; height: 100%"></div>
						</div>
					</div>
				</template>

				<!-- Danh sách phòng ban thực -->
				<template v-else>
					<div
						v-for="(dept, idx) in stats.department_stats"
						:key="dept.department_name + '-' + idx"
						class="bento-card dept-card"
					>
						<div class="dept-header">
							<span class="dept-name">{{ dept.department_name }}</span>
							<span class="dept-count">{{ dept.employee_count }} nhân sự</span>
						</div>
						<div class="progress-container">
							<div
								class="progress-bar"
								:style="{
									width: (dept.employee_count / (stats.total_employees || 1)) * 100 + '%',
								}"
							></div>
						</div>
						<div class="dept-percentage">
							{{ stats.total_employees ? Math.round((dept.employee_count / stats.total_employees) * 100) : 0 }}% tổng số
						</div>
					</div>
				</template>
			</div>
		</section>
	</div>
</template>

<style scoped>
.dashboard-container {
	color: var(--text-main);
	padding: var(--space-1) 0;
	display: flex;
	flex-direction: column;
	gap: var(--space-4);
}

/* Bento Grid */
.bento-grid {
	display: grid;
	grid-template-columns: repeat(3, 1fr);
	gap: var(--space-3);
}

.bento-card {
	background: var(--bg-card);
	backdrop-filter: var(--glass-backdrop);
	-webkit-backdrop-filter: var(--glass-backdrop);
	border: var(--glass-border);
	border-radius: var(--radius-lg);
	box-shadow: var(--glass-shadow);
	padding: var(--space-3);
	transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
	position: relative;
	overflow: hidden;
}

.bento-card:hover {
	transform: translateY(-4px);
	box-shadow: var(--glass-shadow-hover);
}

/* Welcome Card */
.welcome-card {
	grid-column: span 3;
	background: linear-gradient(135deg, rgba(0, 192, 250, 0.08) 0%, rgba(66, 97, 237, 0.08) 50%, rgba(103, 23, 204, 0.08) 100%);
	display: flex;
	justify-content: space-between;
	align-items: center;
	border: 1px solid rgba(255, 255, 255, 0.8);
}

.welcome-content {
	z-index: 2;
	max-width: 70%;
}

.welcome-badge {
	display: inline-block;
	font-family: var(--font-widget);
	font-size: var(--fs-xs);
	font-weight: var(--fw-bold);
	color: var(--primary-color);
	background: rgba(66, 97, 237, 0.1);
	padding: 0.25rem 0.75rem;
	border-radius: var(--radius-full);
	margin-bottom: var(--space-2);
	text-transform: uppercase;
	letter-spacing: 0.05em;
}

.welcome-title {
	font-size: var(--fs-2xl);
	font-weight: var(--fw-bold);
	color: var(--text-main);
	margin-bottom: 0.5rem;
}

.welcome-text {
	font-size: var(--fs-sm);
	color: var(--text-muted);
	line-height: var(--lh-normal);
	margin: 0;
}

.welcome-illustration {
	position: absolute;
	right: -20px;
	top: -20px;
	width: 150px;
	height: 150px;
	z-index: 1;
	opacity: 0.8;
}

.blob {
	position: absolute;
	border-radius: 50%;
	filter: blur(25px);
}

.blob-1 {
	background: rgba(0, 192, 250, 0.25);
	width: 100px;
	height: 100px;
	right: 10px;
	top: 10px;
	animation: floatBlob 8s infinite alternate;
}

.blob-2 {
	background: rgba(103, 23, 204, 0.15);
	width: 80px;
	height: 80px;
	right: 40px;
	top: 40px;
	animation: floatBlob 6s infinite alternate-reverse;
}

@keyframes floatBlob {
	0% { transform: translate(0, 0) scale(1); }
	100% { transform: translate(10px, 10px) scale(1.1); }
}

/* Stat Cards */
.stat-card-inner {
	display: flex;
	flex-direction: column;
	height: 100%;
	justify-content: space-between;
	gap: var(--space-2);
}

.stat-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
}

.stat-label {
	font-size: var(--fs-sm);
	font-weight: var(--fw-semibold);
	color: var(--text-muted);
}

.stat-icon-wrapper {
	display: flex;
	align-items: center;
	justify-content: center;
	width: 42px;
	height: 42px;
	border-radius: var(--radius-md);
	background: rgba(255, 255, 255, 0.5);
	box-shadow: inset 0 2px 4px rgba(66, 97, 237, 0.05);
}

.stat-icon {
	width: 22px;
	height: 22px;
}

.stat-body {
	margin-top: auto;
}

.stat-value {
	font-family: var(--font-widget);
	font-size: var(--fs-2xl);
	font-weight: var(--fw-bold);
	color: var(--text-main);
	line-height: 1;
}

.stat-footer {
	font-size: var(--fs-xs);
	color: var(--text-light);
	font-weight: var(--fw-medium);
	margin-top: 2px;
}

.trend-badge {
	display: inline-block;
	padding: 0.15rem 0.5rem;
	border-radius: var(--radius-sm);
	font-weight: var(--fw-semibold);
}

.trend-badge--success {
	background: rgba(16, 185, 129, 0.1);
	color: #059669;
}

.stat-card--teal .stat-icon {
	filter: invert(59%) sepia(93%) saturate(1915%) hue-rotate(167deg) brightness(101%) contrast(101%);
}
.stat-card--emerald .stat-icon {
	filter: invert(34%) sepia(85%) saturate(3015%) hue-rotate(222deg) brightness(98%) contrast(93%);
}
.stat-card--amber .stat-icon {
	filter: invert(19%) sepia(88%) saturate(3663%) hue-rotate(265deg) brightness(85%) contrast(108%);
}

.stat-card--teal .stat-icon-wrapper {
	background: rgba(0, 192, 250, 0.1);
	border: 1px solid rgba(0, 192, 250, 0.15);
}
.stat-card--emerald .stat-icon-wrapper {
	background: rgba(66, 97, 237, 0.1);
	border: 1px solid rgba(66, 97, 237, 0.15);
}
.stat-card--amber .stat-icon-wrapper {
	background: rgba(103, 23, 204, 0.1);
	border: 1px solid rgba(103, 23, 204, 0.15);
}

/* Department Section */
.dept-section {
	margin-top: var(--space-2);
}

.section-header {
	margin-bottom: var(--space-3);
}

.section-title {
	font-size: var(--fs-lg);
	font-weight: var(--fw-bold);
	color: var(--text-main);
	margin-bottom: 2px;
}

.section-subtitle {
	font-size: var(--fs-sm);
	color: var(--text-muted);
	margin: 0;
}

.dept-grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
	gap: var(--space-3);
}

.dept-card {
	display: flex;
	flex-direction: column;
	gap: var(--space-2);
	background: var(--bg-card);
}

.dept-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
}

.dept-name {
	font-family: var(--font-title);
	font-size: var(--fs-base);
	font-weight: var(--fw-semibold);
	color: var(--text-main);
}

.dept-count {
	font-size: var(--fs-xs);
	font-weight: var(--fw-semibold);
	color: var(--primary-color);
	background: rgba(66, 97, 237, 0.08);
	padding: 0.2rem 0.6rem;
	border-radius: var(--radius-full);
}

.progress-container {
	height: 8px;
	background: rgba(66, 97, 237, 0.06);
	border-radius: var(--radius-full);
	overflow: hidden;
	margin-top: 4px;
}

.progress-bar {
	height: 100%;
	background: linear-gradient(90deg, #00C0FA 0%, #4261ED 50%, #6717CC 100%);
	border-radius: var(--radius-full);
	transition: width 1s cubic-bezier(0.4, 0, 0.2, 1);
}

.dept-percentage {
	font-size: var(--fs-xs);
	color: var(--text-light);
	align-self: flex-end;
	font-weight: var(--fw-medium);
}

/* Responsive */
@media (max-width: 1024px) {
	.bento-grid {
		grid-template-columns: repeat(2, 1fr);
	}
	.welcome-card {
		grid-column: span 2;
	}
	.stat-card--amber {
		grid-column: span 2;
	}
}

@media (max-width: 640px) {
	.bento-grid {
		grid-template-columns: 1fr;
	}
	.welcome-card {
		grid-column: span 1;
		flex-direction: column;
		align-items: flex-start;
		gap: var(--space-3);
	}
	.welcome-content {
		max-width: 100%;
	}
	.welcome-illustration {
		display: none;
	}
	.stat-card--amber {
		grid-column: span 1;
	}
	.dept-grid {
		grid-template-columns: 1fr;
	}
}
</style>
