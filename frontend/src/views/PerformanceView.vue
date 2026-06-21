<script setup>
import { ref, computed } from "vue";
import axios from "axios";
import { useToast } from "vue-toastification";

const toast = useToast();

const loading = ref(false);
const logs = ref([]);
const stats = ref({
	noCacheMaxMs: 0,
	cacheMaxMs: 0,
});

const apiBase = import.meta.env.VITE_API_URL || "http://localhost:8080/api/v1";

const testAPI = async (useCache) => {
	loading.value = true;
	try {
		const start = performance.now();
		const response = await axios.get(
			`${apiBase}/performance/test?use_cache=${useCache}`,
		);
		const end = performance.now();

		const data = response.data.data;
		const clientMs = Math.round(end - start);
		const serverMs = data.response_ms;

		const logEntry = {
			id: Date.now(),
			type: useCache ? "Redis Cache" : "MySQL",
			cacheHit: data.cache_hit,
			serverMs: serverMs,
			clientMs: clientMs,
			items: data.items_count,
			message: data.message,
			timestamp: new Date().toLocaleTimeString(),
			rawResponse: {
				status: "success",
				data: {
					cache_hit: data.cache_hit,
					response_ms: serverMs,
					items_count: data.items_count,
					message: data.message,
					provider: useCache ? "Redis Cache Store" : "MySQL RDBMS Database",
				}
			}
		};

		logs.value.unshift(logEntry);

		// Update stats for chart visualization
		if (!useCache && serverMs > stats.value.noCacheMaxMs) {
			stats.value.noCacheMaxMs = serverMs;
		}
		if (useCache && serverMs > stats.value.cacheMaxMs) {
			stats.value.cacheMaxMs = serverMs;
		}
	} catch (error) {
		toast.error("Lỗi khi gọi API: " + error.message);
	} finally {
		loading.value = false;
	}
};

const clearCache = async () => {
	loading.value = true;
	try {
		await axios.delete(`${apiBase}/performance/clear`);
		toast.success(
			"Đã xóa Redis Cache thành công. Lần gọi tiếp theo sẽ là Cold Start.",
		);
		stats.value.cacheMaxMs = 0;
	} catch (error) {
		toast.error("Lỗi khi xóa Cache");
	} finally {
		loading.value = false;
	}
};

// Toggle for JSON detail expansion
const expandedLogId = ref(null);
const toggleLogDetails = (id) => {
	expandedLogId.value = expandedLogId.value === id ? null : id;
};

// Computed KPI Stats
const mysqlLogs = computed(() => logs.value.filter(l => l.type === 'MySQL'));
const redisLogs = computed(() => logs.value.filter(l => l.type === 'Redis Cache'));

const avgMysqlMs = computed(() => {
	if (mysqlLogs.value.length === 0) return 0;
	const sum = mysqlLogs.value.reduce((acc, curr) => acc + curr.serverMs, 0);
	return Math.round(sum / mysqlLogs.value.length);
});

const avgRedisMs = computed(() => {
	if (redisLogs.value.length === 0) return 0;
	const sum = redisLogs.value.reduce((acc, curr) => acc + curr.serverMs, 0);
	return Math.round(sum / redisLogs.value.length);
});

const speedupFactor = computed(() => {
	if (avgRedisMs.value === 0 || avgMysqlMs.value === 0) return 0;
	return (avgMysqlMs.value / avgRedisMs.value).toFixed(1);
});

const hitRatio = computed(() => {
	if (logs.value.length === 0) return 0;
	const hits = logs.value.filter(l => l.cacheHit).length;
	return Math.round((hits / logs.value.length) * 100);
});
</script>

<template>
	<div class="performance-container">
		<!-- Navigation Link to main system -->
		<div class="top-nav">
			<router-link to="/dashboard" class="back-link">
				<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="19" y1="12" x2="5" y2="12"></line><polyline points="12 19 5 12 12 5"></polyline></svg>
				Quay lại Dashboard
			</router-link>
			<div class="status-indicator">
				<span class="pulse-dot"></span>
				<span class="status-text">Redis Connection: Active</span>
			</div>
		</div>

		<!-- Dashboard Header -->
		<header class="perf-header">
			<div class="header-content">
				<div class="title-group">
					<div class="logo-box">
						<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polygon points="12 2 2 7 12 12 22 7 12 2"></polygon><polyline points="2 17 12 22 22 17"></polyline><polyline points="2 12 12 17 22 12"></polyline></svg>
					</div>
					<h1>Performance Monitor</h1>
				</div>
				<p>Phân tích, so sánh tốc độ xử lý truy vấn dữ liệu trực tiếp từ cơ sở dữ liệu MySQL và bộ nhớ tạm Redis Cache.</p>
			</div>
		</header>

		<div class="dashboard-content">
			<!-- Metrics Row (Bento Grid) -->
			<div class="metrics-grid">
				<div class="metric-card">
					<div class="metric-icon db-icon">
						<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><ellipse cx="12" cy="5" rx="9" ry="3"></ellipse><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"></path><path d="M3 12c0 1.66 4 3 9 3s9-1.34 9-3"></path></svg>
					</div>
					<div class="metric-info">
						<span class="metric-label">MySQL Avg Latency</span>
						<span class="metric-value" :class="{ 'has-value': avgMysqlMs > 0 }">
							{{ avgMysqlMs > 0 ? avgMysqlMs + ' ms' : '---' }}
						</span>
						<span class="metric-sub text-db">Direct DB Queries</span>
					</div>
				</div>

				<div class="metric-card">
					<div class="metric-icon redis-icon">
						<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"></path></svg>
					</div>
					<div class="metric-info">
						<span class="metric-label">Redis Avg Latency</span>
						<span class="metric-value text-emerald" :class="{ 'has-value': avgRedisMs > 0 }">
							{{ avgRedisMs > 0 ? avgRedisMs + ' ms' : '---' }}
						</span>
						<span class="metric-sub text-redis">In-Memory Store</span>
					</div>
				</div>

				<div class="metric-card">
					<div class="metric-icon factor-icon">
						<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon></svg>
					</div>
					<div class="metric-info">
						<span class="metric-label">Performance Boost</span>
						<span class="metric-value text-purple" :class="{ 'has-value': parseFloat(speedupFactor) > 0 }">
							{{ parseFloat(speedupFactor) > 0 ? speedupFactor + 'x Faster' : '---' }}
						</span>
						<span class="metric-sub text-purple-sub">Average Speedup</span>
					</div>
				</div>

				<div class="metric-card">
					<div class="metric-icon hit-icon">
						<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
					</div>
					<div class="metric-info">
						<span class="metric-label">Cache Hit Ratio</span>
						<span class="metric-value text-amber" :class="{ 'has-value': logs.value?.length > 0 || hitRatio > 0 }">
							{{ logs.length > 0 ? hitRatio + '%' : '---' }}
						</span>
						<span class="metric-sub text-amber-sub">{{ logs.length }} total requests</span>
					</div>
				</div>
			</div>

			<!-- Main Layout Grid -->
			<div class="perf-main-grid">
				<!-- Control Panel Card -->
				<section class="card controls-card">
					<div class="card-header">
						<h3>Bảng Điều Khiển</h3>
						<p>Thực hiện các truy vấn giả lập để kiểm tra hiệu năng.</p>
					</div>

					<div class="action-buttons-group">
						<button
							class="action-btn btn-mysql"
							:disabled="loading"
							@click="testAPI(false)"
						>
							<div class="btn-visual">
								<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><ellipse cx="12" cy="5" rx="9" ry="3"></ellipse><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"></path><path d="M3 12c0 1.66 4 3 9 3s9-1.34 9-3"></path></svg>
							</div>
							<div class="btn-text">
								<span class="btn-title">MySQL Query</span>
								<span class="btn-desc">Truy vấn trực tiếp (Cold Start)</span>
							</div>
							<span class="btn-latency">~500ms</span>
						</button>

						<button
							class="action-btn btn-redis-cached"
							:disabled="loading"
							@click="testAPI(true)"
						>
							<div class="btn-visual">
								<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon></svg>
							</div>
							<div class="btn-text">
								<span class="btn-title">Redis Cached</span>
								<span class="btn-desc">Truy xuất in-memory siêu tốc</span>
							</div>
							<span class="btn-latency glow-latency">~8ms</span>
						</button>

						<div class="divider"></div>

						<button
							class="action-btn btn-clear-cache"
							:disabled="loading"
							@click="clearCache"
						>
							<div class="btn-visual">
								<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2M10 11v6M14 11v6"></path></svg>
							</div>
							<div class="btn-text">
								<span class="btn-title">Flush Redis Cache</span>
								<span class="btn-desc">Xóa mọi cache để kiểm tra lại</span>
							</div>
						</button>
					</div>

					<div class="architecture-flow">
						<h4>Cơ chế hoạt động</h4>
						<div class="flow-steps">
							<div class="flow-step">
								<div class="step-num">1</div>
								<p><strong>MySQL Direct:</strong> Đọc trực tiếp từ đĩa cứng. Độ trễ cao (~500ms do mô phỏng tải lớn).</p>
							</div>
							<div class="flow-step">
								<div class="step-num">2</div>
								<p><strong>Redis Cache:</strong> Lưu kết quả vào RAM. Các yêu cầu tiếp theo trả về tức thì (~8ms).</p>
							</div>
						</div>
					</div>
				</section>

				<!-- Visualization & Logs Card -->
				<section class="card results-card">
					<div class="card-header border-b">
						<h3>Giám Sát Thời Gian Thực</h3>
						<p>Biểu đồ so sánh và lịch sử chi tiết các yêu cầu.</p>
					</div>

					<!-- Visual Latency Chart -->
					<div class="chart-section" v-if="avgMysqlMs > 0 || avgRedisMs > 0">
						<h4 class="chart-title">So Sánh Độ Trễ Trung Bình (Ms)</h4>
						<div class="custom-chart-wrapper">
							<!-- MySQL Bar -->
							<div class="chart-row">
								<span class="chart-label">MySQL (Direct)</span>
								<div class="bar-container">
									<div class="bar-bar db-bar" :style="{ width: '100%' }">
										<span class="bar-value">{{ avgMysqlMs }}ms</span>
									</div>
								</div>
							</div>
							<!-- Redis Bar -->
							<div class="chart-row">
								<span class="chart-label">Redis (Cache)</span>
								<div class="bar-container">
									<div class="bar-bar redis-bar" :style="{ width: Math.max(2, Math.min(100, (avgRedisMs / (avgMysqlMs || 1)) * 100)) + '%' }">
										<span class="bar-value">{{ avgRedisMs || '---' }}ms</span>
									</div>
								</div>
								<span class="boost-badge" v-if="parseFloat(speedupFactor) > 0">
									{{ speedupFactor }}x Nhanh hơn
								</span>
							</div>
						</div>
					</div>

					<!-- Live Logs Table -->
					<div class="logs-wrapper">
						<div class="logs-header-meta">
							<h4>Nhật Ký Yêu Cầu (Logs)</h4>
							<span class="total-badge">{{ logs.length }} logs</span>
						</div>

						<div class="logs-table-container" v-if="logs.length > 0">
							<div class="logs-list">
								<div
									v-for="log in logs"
									:key="log.id"
									class="log-row-container"
								>
									<!-- Log Header (Clickable) -->
									<div
										class="log-row"
										@click="toggleLogDetails(log.id)"
										:class="{ 'is-expanded': expandedLogId === log.id }"
									>
										<div class="log-left">
											<span class="method-badge">GET</span>
											<span class="path-text">/performance/test</span>
										</div>

										<div class="log-mid">
											<span
												class="status-pill"
												:class="
													log.cacheHit
														? 'status-hit'
														: log.type === 'Redis Cache'
															? 'status-miss'
															: 'status-db'
												"
											>
												<span class="status-dot"></span>
												{{ log.type }}
												{{
													log.cacheHit
														? "Hit"
														: log.type === "Redis Cache"
															? "Miss"
															: "Direct"
												}}
											</span>
											<span class="size-text">{{ log.items }} items ({{ log.cacheHit ? '0.36' : '9.60' }} KB)</span>
										</div>

										<div class="log-right">
											<span class="time-text">{{ log.timestamp }}</span>
											<span class="latency-badge" :class="log.serverMs > 100 ? 'slow' : 'fast'">
												{{ log.serverMs }}ms
											</span>
											<svg class="chevron-icon" :class="{ 'rotated': expandedLogId === log.id }" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="6 9 12 15 18 9"></polyline></svg>
										</div>
									</div>

									<!-- Expanded JSON Inspector -->
									<div class="log-details" v-if="expandedLogId === log.id">
										<div class="details-inner">
											<div class="details-header">
												<span>JSON API Response Payload</span>
												<span class="mime">application/json</span>
											</div>
											<pre class="json-code"><code><span class="j-key">"status"</span>: <span class="j-val">"success"</span>,
<span class="j-key">"data"</span>: {
  <span class="j-key">"cache_hit"</span>: <span class="j-val">{{ log.rawResponse.data.cache_hit }}</span>,
  <span class="j-key">"response_ms"</span>: <span class="j-num">{{ log.rawResponse.data.response_ms }}</span>,
  <span class="j-key">"items_count"</span>: <span class="j-num">{{ log.rawResponse.data.items_count }}</span>,
  <span class="j-key">"message"</span>: <span class="j-val">"{{ log.rawResponse.data.message }}"</span>,
  <span class="j-key">"provider"</span>: <span class="j-val">"{{ log.rawResponse.data.provider }}"</span>
}</code></pre>
										</div>
									</div>
								</div>
							</div>
						</div>

						<div v-else class="empty-state">
							<div class="empty-icon">
								<svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>
							</div>
							<p class="empty-title">Chưa có nhật ký truy cập</p>
							<p class="empty-desc">Nhấp vào các nút bên bảng điều khiển để kích hoạt các yêu cầu và kiểm thử tốc độ.</p>
						</div>
					</div>
				</section>
			</div>
		</div>
	</div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700&display=swap');

.performance-container {
	min-height: 100vh;
	background-color: #f8fafc;
	font-family: "Plus Jakarta Sans", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
	color: #0f172a;
	padding-bottom: 3rem;
}

/* Top Nav bar */
.top-nav {
	max-width: 1280px;
	margin: 0 auto;
	padding: 1.25rem 2rem 0.5rem;
	display: flex;
	justify-content: space-between;
	align-items: center;
}

.back-link {
	display: flex;
	align-items: center;
	gap: 0.5rem;
	color: #64748b;
	text-decoration: none;
	font-weight: 600;
	font-size: 0.875rem;
	transition: color 0.2s ease;
}

.back-link:hover {
	color: #0f172a;
}

.status-indicator {
	display: flex;
	align-items: center;
	gap: 0.5rem;
	background: #ffffff;
	padding: 0.35rem 0.75rem;
	border-radius: 9999px;
	border: 1px solid rgba(0, 0, 0, 0.06);
	box-shadow: 0 1px 2px rgba(0, 0, 0, 0.02);
}

.pulse-dot {
	width: 8px;
	height: 8px;
	background-color: #10b981;
	border-radius: 50%;
	display: inline-block;
	box-shadow: 0 0 0 0 rgba(16, 185, 129, 0.7);
	animation: pulse 1.8s infinite;
}

@keyframes pulse {
	0% {
		transform: scale(0.95);
		box-shadow: 0 0 0 0 rgba(16, 185, 129, 0.5);
	}
	70% {
		transform: scale(1);
		box-shadow: 0 0 0 6px rgba(16, 185, 129, 0);
	}
	100% {
		transform: scale(0.95);
		box-shadow: 0 0 0 0 rgba(16, 185, 129, 0);
	}
}

.status-text {
	font-size: 0.75rem;
	font-weight: 600;
	color: #475569;
}

/* Header */
.perf-header {
	max-width: 1280px;
	margin: 0 auto;
	padding: 1.5rem 2rem 2rem;
}

.title-group {
	display: flex;
	align-items: center;
	gap: 0.75rem;
	margin-bottom: 0.5rem;
}

.logo-box {
	background: #0f172a;
	color: #ffffff;
	width: 36px;
	height: 36px;
	border-radius: 8px;
	display: flex;
	align-items: center;
	justify-content: center;
	box-shadow: 0 4px 12px rgba(15, 23, 42, 0.15);
}

.header-content h1 {
	margin: 0;
	font-size: 1.75rem;
	font-weight: 700;
	color: #0f172a;
	letter-spacing: -0.03em;
}

.header-content p {
	margin: 0;
	font-size: 0.925rem;
	line-height: 1.5;
	color: #64748b;
	max-width: 720px;
}

/* Dashboard Content Wrapper */
.dashboard-content {
	max-width: 1280px;
	margin: 0 auto;
	padding: 0 2rem;
}

/* Metrics Bento Grid */
.metrics-grid {
	display: grid;
	grid-template-columns: repeat(4, 1fr);
	gap: 1.25rem;
	margin-bottom: 2rem;
}

@media (max-width: 1024px) {
	.metrics-grid {
		grid-template-columns: repeat(2, 1fr);
	}
}

@media (max-width: 640px) {
	.metrics-grid {
		grid-template-columns: 1fr;
	}
}

.metric-card {
	background: #ffffff;
	border-radius: 12px;
	padding: 1.25rem;
	border: 1px solid rgba(0, 0, 0, 0.06);
	box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.02), 0 2px 4px -1px rgba(0, 0, 0, 0.01);
	display: flex;
	align-items: flex-start;
	gap: 1rem;
	transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.metric-card:hover {
	transform: translateY(-2px);
	box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.04), 0 4px 6px -2px rgba(0, 0, 0, 0.02);
}

.metric-icon {
	width: 42px;
	height: 42px;
	border-radius: 10px;
	display: flex;
	align-items: center;
	justify-content: center;
}

.db-icon {
	background: #eff6ff;
	color: #2563eb;
}

.redis-icon {
	background: #ecfdf5;
	color: #059669;
}

.factor-icon {
	background: #f5f3ff;
	color: #7c3aed;
}

.hit-icon {
	background: #fffbeb;
	color: #d97706;
}

.metric-info {
	display: flex;
	flex-direction: column;
}

.metric-label {
	font-size: 0.75rem;
	font-weight: 600;
	color: #64748b;
	text-transform: uppercase;
	letter-spacing: 0.05em;
	margin-bottom: 0.25rem;
}

.metric-value {
	font-size: 1.5rem;
	font-weight: 700;
	color: #94a3b8;
	line-height: 1.2;
}

.metric-value.has-value {
	color: #0f172a;
}

.metric-value.text-emerald {
	color: #10b981;
}

.metric-value.text-purple {
	color: #8b5cf6;
}

.metric-value.text-amber {
	color: #f59e0b;
}

.metric-sub {
	font-size: 0.75rem;
	font-weight: 500;
	margin-top: 0.25rem;
}

.text-db { color: #3b82f6; }
.text-redis { color: #10b981; }
.text-purple-sub { color: #a78bfa; }
.text-amber-sub { color: #fbbf24; }

/* Main Grid Layout */
.perf-main-grid {
	display: grid;
	grid-template-columns: 360px 1fr;
	gap: 1.5rem;
}

@media (max-width: 900px) {
	.perf-main-grid {
		grid-template-columns: 1fr;
	}
}

.card {
	background: #ffffff;
	border-radius: 16px;
	border: 1px solid rgba(0, 0, 0, 0.06);
	box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.02), 0 2px 4px -1px rgba(0, 0, 0, 0.01);
	overflow: hidden;
}

.card-header {
	padding: 1.5rem;
}

.card-header.border-b {
	border-bottom: 1px solid #f1f5f9;
}

.card-header h3 {
	margin: 0 0 0.25rem;
	font-size: 1.15rem;
	font-weight: 700;
	color: #0f172a;
}

.card-header p {
	margin: 0;
	font-size: 0.8rem;
	color: #64748b;
}

/* Control buttons group */
.action-buttons-group {
	padding: 0 1.5rem 1.5rem;
	display: flex;
	flex-direction: column;
	gap: 1rem;
}

.action-btn {
	display: flex;
	align-items: center;
	padding: 1rem;
	border-radius: 12px;
	border: 1px solid rgba(0, 0, 0, 0.06);
	background: #ffffff;
	cursor: pointer;
	text-align: left;
	transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
	position: relative;
	width: 100%;
}

.action-btn:disabled {
	opacity: 0.6;
	cursor: not-allowed;
}

.btn-visual {
	width: 38px;
	height: 38px;
	border-radius: 8px;
	display: flex;
	align-items: center;
	justify-content: center;
	margin-right: 0.85rem;
	transition: all 0.2s ease;
}

.btn-text {
	display: flex;
	flex-direction: column;
	flex: 1;
}

.btn-title {
	font-size: 0.9rem;
	font-weight: 700;
	color: #0f172a;
}

.btn-desc {
	font-size: 0.725rem;
	color: #64748b;
	margin-top: 0.15rem;
}

.btn-latency {
	font-size: 0.75rem;
	font-weight: 700;
	color: #94a3b8;
	background: #f1f5f9;
	padding: 0.25rem 0.5rem;
	border-radius: 6px;
}

/* MySQL Button State */
.btn-mysql:hover:not(:disabled) {
	border-color: #3b82f6;
	box-shadow: 0 4px 12px rgba(59, 130, 246, 0.08);
}
.btn-mysql:hover:not(:disabled) .btn-visual {
	background: #eff6ff;
	color: #2563eb;
}

/* Redis Button State */
.btn-redis-cached {
	border-color: rgba(16, 185, 129, 0.15);
	background: linear-gradient(to right, #ffffff, #f0fdf4);
}
.btn-redis-cached .btn-visual {
	background: #dcfce7;
	color: #15803d;
}
.btn-redis-cached:hover:not(:disabled) {
	border-color: #10b981;
	box-shadow: 0 4px 14px rgba(16, 185, 129, 0.12);
}
.btn-redis-cached:hover:not(:disabled) .btn-visual {
	background: #10b981;
	color: #ffffff;
}
.btn-redis-cached .glow-latency {
	background: #dcfce7;
	color: #166534;
}

.divider {
	height: 1px;
	background: #f1f5f9;
	margin: 0.25rem 0;
}

/* Clear Cache Button State */
.btn-clear-cache {
	border-color: rgba(239, 68, 68, 0.15);
}
.btn-clear-cache .btn-visual {
	background: #fee2e2;
	color: #b91c1c;
}
.btn-clear-cache:hover:not(:disabled) {
	border-color: #ef4444;
	background: #fdf2f2;
}

/* Architecture Flow Section */
.architecture-flow {
	background: #f8fafc;
	margin: 0 1.5rem 1.5rem;
	padding: 1.25rem;
	border-radius: 12px;
	border: 1px dashed rgba(0, 0, 0, 0.08);
}

.architecture-flow h4 {
	margin: 0 0 0.75rem;
	font-size: 0.8rem;
	font-weight: 700;
	text-transform: uppercase;
	letter-spacing: 0.05em;
	color: #475569;
}

.flow-steps {
	display: flex;
	flex-direction: column;
	gap: 0.75rem;
}

.flow-step {
	display: flex;
	align-items: flex-start;
	gap: 0.75rem;
}

.step-num {
	width: 18px;
	height: 18px;
	background: #0f172a;
	color: #ffffff;
	border-radius: 50%;
	font-size: 0.65rem;
	font-weight: 700;
	display: flex;
	align-items: center;
	justify-content: center;
	flex-shrink: 0;
	margin-top: 0.1rem;
}

.flow-step p {
	margin: 0;
	font-size: 0.75rem;
	line-height: 1.4;
	color: #475569;
}

/* Charts Section */
.chart-section {
	padding: 1.5rem;
	border-bottom: 1px solid #f1f5f9;
}

.chart-title {
	margin: 0 0 1rem;
	font-size: 0.8rem;
	font-weight: 700;
	color: #475569;
	text-transform: uppercase;
	letter-spacing: 0.05em;
}

.custom-chart-wrapper {
	display: flex;
	flex-direction: column;
	gap: 1rem;
}

.chart-row {
	display: flex;
	align-items: center;
	gap: 1rem;
}

.chart-label {
	width: 110px;
	font-size: 0.75rem;
	font-weight: 600;
	color: #64748b;
}

.bar-container {
	flex: 1;
	height: 24px;
	background: #f1f5f9;
	border-radius: 6px;
	overflow: hidden;
	position: relative;
}

.bar-bar {
	height: 100%;
	border-radius: 6px;
	display: flex;
	align-items: center;
	justify-content: flex-end;
	padding-right: 0.75rem;
	transition: width 0.6s cubic-bezier(0.4, 0, 0.2, 1);
}

.db-bar {
	background: linear-gradient(90deg, #3b82f6, #2563eb);
}

.redis-bar {
	background: linear-gradient(90deg, #10b981, #059669);
}

.bar-value {
	font-size: 0.75rem;
	font-weight: 700;
	color: #ffffff;
}

.boost-badge {
	background: #f5f3ff;
	color: #7c3aed;
	font-size: 0.7rem;
	font-weight: 700;
	padding: 0.25rem 0.5rem;
	border-radius: 9999px;
	border: 1px solid rgba(124, 58, 237, 0.15);
}

/* Logs Console Styling */
.logs-wrapper {
	padding: 1.5rem;
}

.logs-header-meta {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 1rem;
}

.logs-header-meta h4 {
	margin: 0;
	font-size: 0.85rem;
	font-weight: 700;
	color: #475569;
	text-transform: uppercase;
	letter-spacing: 0.05em;
}

.total-badge {
	font-size: 0.7rem;
	font-weight: 700;
	background: #f1f5f9;
	color: #475569;
	padding: 0.2rem 0.5rem;
	border-radius: 9999px;
}

.logs-table-container {
	max-height: 480px;
	overflow-y: auto;
	border: 1px solid #f1f5f9;
	border-radius: 12px;
}

.logs-list {
	display: flex;
	flex-direction: column;
}

.log-row-container {
	border-bottom: 1px solid #f1f5f9;
}

.log-row-container:last-child {
	border-bottom: none;
}

.log-row {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 0.85rem 1rem;
	cursor: pointer;
	transition: background-color 0.15s ease;
}

.log-row:hover, .log-row.is-expanded {
	background-color: #f8fafc;
}

.log-left {
	display: flex;
	align-items: center;
	gap: 0.5rem;
}

.method-badge {
	font-size: 0.65rem;
	font-weight: 700;
	background: #e2e8f0;
	color: #475569;
	padding: 0.15rem 0.35rem;
	border-radius: 4px;
}

.path-text {
	font-size: 0.8rem;
	font-weight: 600;
	color: #334155;
}

.log-mid {
	display: flex;
	align-items: center;
	gap: 0.75rem;
}

.status-pill {
	display: flex;
	align-items: center;
	gap: 0.35rem;
	font-size: 0.725rem;
	font-weight: 600;
	padding: 0.2rem 0.5rem;
	border-radius: 9999px;
}

.status-dot {
	width: 5px;
	height: 5px;
	border-radius: 50%;
}

.status-hit {
	background: #ecfdf5;
	color: #047857;
}
.status-hit .status-dot { background-color: #10b981; }

.status-miss {
	background: #fffbeb;
	color: #b45309;
}
.status-miss .status-dot { background-color: #fbbf24; }

.status-db {
	background: #eff6ff;
	color: #1d4ed8;
}
.status-db .status-dot { background-color: #3b82f6; }

.size-text {
	font-size: 0.7rem;
	color: #94a3b8;
}

.log-right {
	display: flex;
	align-items: center;
	gap: 0.75rem;
}

.time-text {
	font-size: 0.725rem;
	color: #94a3b8;
}

.latency-badge {
	font-size: 0.75rem;
	font-weight: 700;
	padding: 0.15rem 0.45rem;
	border-radius: 6px;
}

.latency-badge.fast {
	background: #dcfce7;
	color: #166534;
}

.latency-badge.slow {
	background: #fee2e2;
	color: #991b1b;
}

.chevron-icon {
	color: #94a3b8;
	transition: transform 0.2s ease;
}

.chevron-icon.rotated {
	transform: rotate(180deg);
}

/* Expanded Log Details */
.log-details {
	background: #0f172a;
	color: #e2e8f0;
	padding: 1.25rem;
	border-top: 1px solid #1e293b;
	border-bottom: 1px solid #1e293b;
	animation: slideDown 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

@keyframes slideDown {
	from {
		opacity: 0;
		transform: translateY(-4px);
	}
	to {
		opacity: 1;
		transform: translateY(0);
	}
}

.details-inner {
	display: flex;
	flex-direction: column;
	gap: 0.5rem;
}

.details-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	border-bottom: 1px solid #1e293b;
	padding-bottom: 0.5rem;
	font-size: 0.725rem;
	font-weight: 600;
	color: #94a3b8;
	text-transform: uppercase;
	letter-spacing: 0.05em;
}

.details-header .mime {
	color: #38bdf8;
	font-family: monospace;
	text-transform: none;
}

.json-code {
	margin: 0;
	font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
	font-size: 0.8rem;
	line-height: 1.5;
	overflow-x: auto;
}

.j-key {
	color: #a78bfa;
}

.j-val {
	color: #34d399;
}

.j-num {
	color: #fbbf24;
}

/* Empty State */
.empty-state {
	text-align: center;
	padding: 4rem 1.5rem;
}

.empty-icon {
	color: #cbd5e1;
	margin-bottom: 1rem;
	display: flex;
	justify-content: center;
}

.empty-title {
	font-weight: 700;
	font-size: 0.95rem;
	color: #334155;
	margin: 0 0 0.25rem;
}

.empty-desc {
	font-size: 0.775rem;
	color: #64748b;
	max-width: 320px;
	margin: 0 auto;
	line-height: 1.4;
}
</style>
