<script setup>
import usersIcon from "@/assets/svg/users.svg";
import checkIcon from "@/assets/svg/check-circle.svg";
import buildingIcon from "@/assets/svg/building.svg";
import keyIcon from "@/assets/svg/key.svg";
import { useDashboardStore } from "@/store/dashboard";
import { computed, onMounted } from "vue";
import { storeToRefs } from "pinia";
import Skeleton from "@/components/Skeleton.vue";

const dashboardStore = useDashboardStore();
const { stats, loading } = storeToRefs(dashboardStore);
const dashboardStats = computed(() => stats.value);

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
		<div class="stats-grid">
			<!-- Stat Card 1 -->
			<div class="stat-card stat-card--blue">
				<div class="stat-icon-wrapper">
					<img :src="usersIcon" class="stat-icon" alt="users" />
				</div>
				<div class="stat-info" style="flex-grow: 1">
					<template v-if="loading">
						<Skeleton
							type="text"
							width="60px"
							height="32px"
							style="margin-bottom: 4px"
						/>
						<Skeleton type="text" width="100px" height="16px" />
					</template>
					<template v-else>
						<div class="stat-value">
							{{ dashboardStats.total_employees }}
						</div>
						<div class="stat-label">Tổng nhân viên</div>
					</template>
				</div>
			</div>

			<!-- Stat Card 2 -->
			<div class="stat-card stat-card--green">
				<div class="stat-icon-wrapper">
					<img :src="checkIcon" class="stat-icon" alt="active" />
				</div>
				<div class="stat-info" style="flex-grow: 1">
					<template v-if="loading">
						<Skeleton
							type="text"
							width="60px"
							height="32px"
							style="margin-bottom: 4px"
						/>
						<Skeleton type="text" width="100px" height="16px" />
					</template>
					<template v-else>
						<div class="stat-value">
							{{ dashboardStats.total_employees_active }}
						</div>
						<div class="stat-label">Đang làm việc</div>
					</template>
				</div>
			</div>

			<!-- Stat Card 3 -->
			<div class="stat-card stat-card--amber">
				<div class="stat-icon-wrapper">
					<img :src="buildingIcon" class="stat-icon" alt="dept" />
				</div>
				<div class="stat-info" style="flex-grow: 1">
					<template v-if="loading">
						<Skeleton
							type="text"
							width="60px"
							height="32px"
							style="margin-bottom: 4px"
						/>
						<Skeleton type="text" width="100px" height="16px" />
					</template>
					<template v-else>
						<div class="stat-value">
							{{ dashboardStats.total_departments }}
						</div>
						<div class="stat-label">Phòng ban</div>
					</template>
				</div>
			</div>
		</div>

		<!-- ===== Phân bổ nhân sự ===== -->
		<section class="dept-section">
			<h2 class="section-title">Nhân viên theo phòng ban</h2>
			<div class="dept-grid">
				<template v-if="loading">
					<div v-for="i in 3" :key="'skeleton-d-' + i" class="dept-card">
						<div class="dept-header">
							<Skeleton type="text" width="120px" height="18px" />
							<Skeleton type="badge" width="40px" height="22px" />
						</div>
						<div class="progress-container">
							<div
								class="skeleton"
								style="width: 100%; height: 100%"
							></div>
						</div>
					</div>
				</template>
				<template v-else>
					<div
						v-for="(dept, idx) in dashboardStats.department_stats"
						:key="dept.department_name + '-' + idx"
						class="dept-card"
					>
						<div class="dept-header">
							<span class="dept-name">{{ dept.department_name }}</span>
							<span class="dept-count">{{ dept.employee_count }}</span>
						</div>
						<div class="progress-container">
							<div
								class="progress-bar"
								:style="{
									width:
										(dept.employee_count /
											(dashboardStats.total_employees || 1)) *
											100 +
										'%',
								}"
							></div>
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
	padding: var(--space-2) 0;
}

.stats-grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
	gap: var(--space-3);
	margin-bottom: var(--space-6);
}

.stat-card {
	display: flex;
	align-items: center;
	gap: var(--space-3);
	padding: var(--space-3);
	background: var(--bg-card);
	border: 1px solid var(--border-color);
	border-radius: var(--radius-lg);
	box-shadow: var(--shadow-sm);
	transition:
		transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1),
		box-shadow 0.3s ease;
}

.stat-card:hover {
	transform: translateY(-4px);
	box-shadow: var(--shadow-hover);
}

.stat-icon-wrapper {
	display: flex;
	align-items: center;
	justify-content: center;
	width: 56px;
	height: 56px;
	border-radius: var(--radius-md);
	flex-shrink: 0;
	background: var(--bg-light);
}

.stat-icon {
	width: 28px;
	height: 28px;
}

.stat-info {
	display: flex;
	flex-direction: column;
}

.stat-value {
	font-size: var(--fs-3xl);
	font-weight: var(--fw-bold);
	line-height: var(--lh-tight);
	margin-bottom: 0.25rem;
}

.stat-label {
	font-size: var(--fs-sm);
	font-weight: var(--fw-medium);
	color: var(--text-muted);
}

.stat-card--blue .stat-icon {
	filter: invert(44%) sepia(87%) saturate(2258%) hue-rotate(200deg)
		brightness(101%) contrast(92%);
}
.stat-card--green .stat-icon {
	filter: invert(61%) sepia(54%) saturate(518%) hue-rotate(97deg)
		brightness(92%) contrast(85%);
}
.stat-card--amber .stat-icon {
	filter: invert(68%) sepia(67%) saturate(3062%) hue-rotate(5deg)
		brightness(103%) contrast(93%);
}
.stat-card--purple .stat-icon {
	filter: invert(41%) sepia(87%) saturate(3268%) hue-rotate(256deg)
		brightness(101%) contrast(96%);
}

.section-title {
	font-size: var(--fs-xl);
	font-weight: var(--fw-bold);
	margin-bottom: var(--space-3);
	letter-spacing: var(--tracking-tight);
}

.dept-grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
	gap: var(--space-3);
}

.dept-card {
	background: var(--bg-card);
	border: 1px solid var(--border-color);
	border-radius: var(--radius-lg);
	padding: var(--space-3);
	box-shadow: var(--shadow-sm);
	transition: all 0.2s ease;
}

.dept-card:hover {
	transform: translateY(-2px);
	box-shadow: var(--shadow-md);
}

.dept-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: var(--space-3);
}

.dept-name {
	font-weight: var(--fw-semibold);
	color: var(--text-main);
}

.dept-count {
	font-size: var(--fs-xs);
	font-weight: var(--fw-semibold);
	color: var(--primary-color);
	background: var(--bg-light);
	padding: 0.25rem 0.75rem;
	border-radius: var(--radius-full);
}

/* Progress bar styles */
.progress-container {
	height: 8px;
	background: var(--bg-light);
	border-radius: var(--radius-full);
	overflow: hidden;
}

.progress-bar {
	height: 100%;
	background: linear-gradient(90deg, var(--primary-color), #60a5fa);
	border-radius: var(--radius-full);
	transition: width 1.2s cubic-bezier(0.4, 0, 0.2, 1);
}

/* ===== Responsive ===== */
@media (max-width: 768px) {
	.stats-grid,
	.dept-grid {
		grid-template-columns: 1fr;
	}

	.stat-card {
		padding: var(--space-2);
	}
}
</style>
