<script setup>
import { Building, CheckCircle, Users, TrendingUp } from "@lucide/vue";
import { useDashboardStore } from "@/store/dashboard";
import { onMounted } from "vue";
import { storeToRefs } from "pinia";
import Skeleton from "@/components/Skeleton.vue";

// ─── Khởi tạo
const dashboardStore = useDashboardStore();
const { stats, loading } = storeToRefs(dashboardStore);

// ─── Tải dữ liệu dashboard ───────────────────────────────────────────────────
const loadDashboard = async () => {
	try {
		await dashboardStore.fetchDashboard();
	} catch (err) {
		console.error("Lỗi khi tải dashboard:", err);
	}
};
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
			<div class="stat-card stat-card--teal">
				<div class="stat-accent-bar stat-accent-bar--teal"></div>
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
					<div class="stat-icon-wrapper stat-icon-wrapper--teal">
						<Users class="stat-icon" />
					</div>
				</div>
				<div class="stat-footer">
					<span class="trend-text">Nhân sự toàn hệ thống</span>
				</div>
			</div>

			<!-- Nhân viên đang làm việc -->
			<div class="stat-card stat-card--emerald">
				<div class="stat-accent-bar stat-accent-bar--emerald"></div>
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
					<div class="stat-icon-wrapper stat-icon-wrapper--emerald">
						<CheckCircle class="stat-icon" />
					</div>
				</div>
				<div class="stat-footer">
					<span class="trend-badge trend-badge--success">
						<TrendingUp class="trend-badge-icon" />
						{{ stats.total_employees ? Math.round((stats.total_employees_active / stats.total_employees) * 100) : 0 }}% hoạt động
					</span>
				</div>
			</div>

			<!-- Tổng số phòng ban -->
			<div class="stat-card stat-card--amber">
				<div class="stat-accent-bar stat-accent-bar--amber"></div>
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
					<div class="stat-icon-wrapper stat-icon-wrapper--amber">
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
						<div class="progress-container dept-skeleton-progress">
							<div class="skeleton dept-skeleton-fill"></div>
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
							<span class="dept-name">{{ dept.department_name }}</span>
							<span class="dept-count">{{ dept.employee_count }} nhân sự</span>
						</div>
						<div class="progress-track">
							<div
								class="progress-bar"
								:style="{ width: ((dept.employee_count / (stats.total_employees || 1)) * 100) + '%' }"
							></div>
						</div>
						<div class="dept-percentage">
							{{ stats.total_employees ? Math.round((dept.employee_count / stats.total_employees) * 100) : 0 }}% tổng số
						</div>
					</div>
					<div v-if="!stats.department_stats?.length" class="dept-empty">
						Chưa có dữ liệu phòng ban.
					</div>
				</template>
			</div>
		</section>
	</div>
</template>

<style scoped>
.dashboard-container {
	color: var(--text-main);
	padding-bottom: var(--space-4);
	display: flex;
	flex-direction: column;
	gap: var(--space-4);
}

/* ── Stat Grid ───────────────────────────────────────── */
.stat-grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
	gap: var(--space-3);
}

.stat-card {
	background: var(--bg-card);
	border: var(--glass-border);
	border-radius: var(--radius-lg);
	box-shadow: var(--glass-shadow);
	padding: var(--space-3);
	display: flex;
	flex-direction: column;
	gap: 0.8rem;
	position: relative;
	overflow: hidden;
	transition: box-shadow 0.25s ease, transform 0.25s ease, border-color 0.25s ease;
}

.stat-card:hover {
	transform: translateY(-3px);
}

.stat-card--teal:hover {
	border-color: rgba(0, 192, 250, 0.25);
	box-shadow: var(--glass-shadow-hover), 0 0 0 1px rgba(0, 192, 250, 0.1);
}
.stat-card--emerald:hover {
	border-color: rgba(16, 185, 129, 0.25);
	box-shadow: var(--glass-shadow-hover), 0 0 0 1px rgba(16, 185, 129, 0.1);
}
.stat-card--amber:hover {
	border-color: rgba(245, 158, 11, 0.25);
	box-shadow: var(--glass-shadow-hover), 0 0 0 1px rgba(245, 158, 11, 0.1);
}

/* top accent bar */
.stat-accent-bar {
	position: absolute;
	top: 0; left: 0; right: 0;
	height: 3px;
	border-radius: var(--radius-lg) var(--radius-lg) 0 0;
}
.stat-accent-bar--teal    { background: linear-gradient(90deg, var(--info-color), var(--primary-color)); }
.stat-accent-bar--emerald { background: linear-gradient(90deg, var(--success-color), #34d399); }
.stat-accent-bar--amber   { background: linear-gradient(90deg, var(--warning-color), #fbbf24); }

.stat-main {
	display: flex;
	align-items: flex-start;
	justify-content: space-between;
	width: 100%;
	padding-top: 6px;
}

.stat-info {
	display: flex;
	flex-direction: column;
	gap: 6px;
}

.stat-icon-wrapper {
	display: flex;
	align-items: center;
	justify-content: center;
	width: 48px;
	height: 48px;
	border-radius: var(--radius-md);
}

.stat-icon-wrapper--teal    { color: var(--info-color); background: var(--info-light); }
.stat-icon-wrapper--emerald { color: var(--success-color); background: var(--success-light); }
.stat-icon-wrapper--amber   { color: var(--warning-color); background: var(--warning-light); }

.stat-icon { width: 22px; height: 22px; }

.stat-footer {
	border-top: 1px solid var(--border-color);
	padding-top: 0.75rem;
	display: flex;
	align-items: center;
}

.trend-badge {
	display: inline-flex;
	align-items: center;
	gap: 4px;
	padding: 3px 9px;
	border-radius: var(--radius-full);
	font-weight: var(--fw-semibold);
	font-size: var(--fs-xs);
}

.trend-badge--success {
	color: var(--success-hover);
	background: var(--success-light);
}

.trend-badge-icon {
	width: 13px;
	height: 13px;
	flex-shrink: 0;
}

.trend-text {
	font-size: var(--fs-xs);
	color: var(--text-light);
	font-weight: var(--fw-medium);
}

/* ── Dept Section ────────────────────────────────────── */
.section-header {
	margin-bottom: 1.2rem;
}

.dept-grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
	gap: var(--space-3);
}

.dept-card {
	display: flex;
	flex-direction: column;
	gap: 1rem;
	background: var(--bg-card);
	border: var(--glass-border);
	border-radius: var(--radius-lg);
	box-shadow: var(--glass-shadow);
	padding: var(--space-3);
	transition: box-shadow 0.25s ease, transform 0.25s ease, border-color 0.25s ease;
}

.dept-card:hover {
	transform: translateY(-3px);
	box-shadow: var(--glass-shadow-hover);
	border-color: rgba(66, 97, 237, 0.18);
}

.dept-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	gap: 8px;
}

.dept-name {
	font-weight: var(--fw-semibold);
	color: var(--text-main);
	font-size: var(--fs-sm);
}

.dept-count {
	font-size: 11px;
	font-weight: var(--fw-semibold);
	color: var(--primary-color);
	background: rgba(66, 97, 237, 0.08);
	padding: 2px 8px;
	border-radius: var(--radius-full);
	white-space: nowrap;
	flex-shrink: 0;
}

.progress-track {
	height: 10px;
	background: rgba(66, 97, 237, 0.06);
	border-radius: var(--radius-full);
	overflow: hidden;
}

.progress-bar {
	height: 100%;
	background: var(--primary-gradient);
	border-radius: var(--radius-full);
	transition: width 1s cubic-bezier(0.4, 0, 0.2, 1);
}

.dept-percentage {
	font-size: var(--fs-xs);
	color: var(--text-light);
	align-self: flex-end;
	font-weight: var(--fw-medium);
}

.dept-empty {
	grid-column: 1 / -1;
	text-align: center;
	padding: var(--space-4) 0;
	color: var(--text-muted);
	font-size: var(--fs-sm);
}

.dept-skeleton-progress {
	margin-top: 1rem;
}

.dept-skeleton-fill {
	width: 100%;
	height: 100%;
}

@media (max-width: 640px) {
	.stat-grid { grid-template-columns: 1fr; }
	.dept-grid { grid-template-columns: 1fr; }
}
</style>
