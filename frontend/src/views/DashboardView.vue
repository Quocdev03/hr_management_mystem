<script setup>
import { Building, CheckCircle, Users } from "@lucide/vue";

// ─── Icon SVG ────────────────────────────────────────────────────────────────

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
	() =>
		userFromStorage?.user_name ?? authStore.user?.user_name ?? "Thành viên",
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
				<p class="page-subtitle">
					Báo cáo tóm tắt trạng thái nhân sự hệ thống
				</p>
			</div>
		</header>

		<!-- Thống kê tổng quan -->
		<div class="stat-grid">
			<!-- Tổng số nhân viên -->
			<div class="stat-card">
				<div class="stat-main">
					<div class="stat-info">
						<span class="stat-label">Tổng nhân viên</span>
						<template v-if="loading">
							<Skeleton type="text" width="60px" height="32px" />
						</template>
						<template v-else>
							<div class="stat-value">{{ stats.total_employees }}</div>
						</template>
					</div>
					<div class="stat-icon-wrapper text-teal">
						<Users class="stat-icon" />
					</div>
				</div>
				<div class="stat-footer">
					<span class="trend-text">Nhân sự toàn hệ thống</span>
				</div>
			</div>

			<!-- Nhân viên đang làm việc -->
			<div class="stat-card">
				<div class="stat-main">
					<div class="stat-info">
						<span class="stat-label">Đang hoạt động</span>
						<template v-if="loading">
							<Skeleton type="text" width="60px" height="32px" />
						</template>
						<template v-else>
							<div class="stat-value">{{ stats.total_employees_active }}</div>
						</template>
					</div>
					<div class="stat-icon-wrapper text-emerald">
						<CheckCircle class="stat-icon" />
					</div>
				</div>
				<div class="stat-footer">
					<span class="trend-badge trend-badge--success">
						{{ stats.total_employees ? Math.round((stats.total_employees_active / stats.total_employees) * 100) : 0 }}% hoạt động
					</span>
				</div>
			</div>

			<!-- Tổng số phòng ban -->
			<div class="stat-card">
				<div class="stat-main">
					<div class="stat-info">
						<span class="stat-label">Phòng ban</span>
						<template v-if="loading">
							<Skeleton type="text" width="60px" height="32px" />
						</template>
						<template v-else>
							<div class="stat-value">{{ stats.total_departments }}</div>
						</template>
					</div>
					<div class="stat-icon-wrapper text-amber">
						<Building class="stat-icon" />
					</div>
				</div>
				<div class="stat-footer">
					<span class="trend-text">Cơ cấu tổ chức hành chính</span>
				</div>
			</div>
		</div>

		<!-- ===== Phân bổ nhân sự theo phòng ban ===== -->
		<section class="dept-section">
			<div class="section-header">
				<h2 class="section-title">Phân bổ nhân sự theo phòng ban</h2>
				<p class="section-subtitle">
					Tỉ lệ phần trăm nhân viên làm việc tại từng phòng ban
				</p>
			</div>

			<div class="dept-grid">
				<!-- Skeleton khi đang tải -->
				<template v-if="loading">
					<div
						v-for="i in 3"
						:key="'skeleton-d-' + i"
						class="dept-card"
					>
						<div class="dept-header">
							<Skeleton type="text" width="120px" height="18px" />
							<Skeleton type="badge" width="40px" height="22px" />
						</div>
						<div
							class="progress-container"
							style="margin-top: 1rem"
						>
							<div
								class="skeleton"
								style="width: 100%; height: 100%"
							></div>
						</div>
					</div>
				</template>

				<!-- Danh sách phòng ban thực -->
				<template v-else>
					<div
						v-for="(dept, idx) in stats.department_stats"
						:key="dept.department_name + '-' + idx"
						class="dept-card"
					>
						<div class="dept-header">
							<span class="dept-name">{{
								dept.department_name
							}}</span>
							<span class="dept-count"
								>{{ dept.employee_count }} nhân sự</span
							>
						</div>
						<div class="progress-container">
							<div
								class="progress-bar"
								:style="{
									width:
										(dept.employee_count /
											(stats.total_employees || 1)) *
											100 +
										'%',
								}"
							></div>
						</div>
						<div class="dept-percentage">
							{{
								stats.total_employees
									? Math.round(
											(dept.employee_count /
												stats.total_employees) *
												100,
										)
									: 0
							}}% tổng số
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

/* Stat Grid */
.stat-grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
	gap: var(--space-3);
}

.stat-card {
	background: var(--bg-card);
	border: 1px solid var(--border-color);
	border-radius: var(--radius-lg);
	box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
	padding: 1.5rem;
	display: flex;
	flex-direction: column;
	gap: .8rem;
	transition: box-shadow 0.2s ease, transform 0.2s ease;
}

.stat-card:hover {
	transform: translateY(-2px);
	box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
}

.stat-main {
	display: flex;
	align-items: flex-start;
	justify-content: space-between;
	width: 100%;
}

.stat-info {
	display: flex;
	flex-direction: column;
	gap: 6px;
}

.stat-label {
	font-size: 0.875rem;
	font-weight: var(--fw-semibold);
	color: var(--text-muted);
}

.stat-value {
	font-family: var(--font-title);
	font-size: 1.75rem;
	font-weight: var(--fw-bold);
	color: var(--text-main);
	line-height: 1.1;
}

.stat-icon-wrapper {
	display: flex;
	align-items: center;
	justify-content: center;
	width: 48px;
	height: 48px;
	border-radius: 12px;
	flex-shrink: 0;
}

.stat-icon {
	width: 24px;
	height: 24px;
}

.text-teal { color: #00C0FA; background: rgba(0, 192, 250, 0.1); }
.text-emerald { color: #10B981; background: rgba(16, 185, 129, 0.1); }
.text-amber { color: #F59E0B; background: rgba(245, 158, 11, 0.1); }

.stat-footer {
	border-top: 1px solid var(--border-color);
	padding-top: 1rem;
	display: flex;
	align-items: center;
}

.trend-text {
	font-size: 0.8125rem;
	color: var(--text-muted);
}

.trend-badge {
	display: inline-block;
	padding: 0.2rem 0.6rem;
	border-radius: var(--radius-full);
	font-weight: var(--fw-semibold);
	font-size: 0.75rem;
}

.trend-badge--success {
	background: rgba(16, 185, 129, 0.1);
	color: #059669;
}

/* Department Section */
.dept-section {
	margin-top: var(--space-2);
}

.section-header {
	margin-bottom: var(--space-3);
}

.dept-grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
	gap: var(--space-3);
}

.dept-card {
	display: flex;
	flex-direction: column;
	gap: 1.25rem;
	background: var(--bg-card);
	border: 1px solid var(--border-color);
	border-radius: var(--radius-lg);
	box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
	padding: 1.5rem;
	transition: box-shadow 0.2s ease, transform 0.2s ease;
}

.dept-card:hover {
	transform: translateY(-2px);
	box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
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
	background: linear-gradient(90deg, #00c0fa 0%, #4261ed 100%);
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
@media (max-width: 640px) {
	.stat-grid {
		grid-template-columns: 1fr;
	}
	.dept-grid {
		grid-template-columns: 1fr;
	}
}
</style>
