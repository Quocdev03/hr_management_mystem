<script setup>
import { onMounted, onUnmounted, watch } from "vue";
import {
	getInitials,
	formatDate,
	formatStatus,
	formatCurrency,
} from "@/helpers/formatters";

const props = defineProps({
	visible: { type: Boolean, required: true },
	employee: { type: Object, default: null }
});

const emit = defineEmits(["close"]);

function handleClose() {
	emit("close");
}

// Close drawer on escape key
function handleKeyDown(e) {
	if (e.key === "Escape" && props.visible) {
		handleClose();
	}
}

onMounted(() => {
	window.addEventListener("keydown", handleKeyDown);
});

onUnmounted(() => {
	window.removeEventListener("keydown", handleKeyDown);
});
</script>

<template>
	<div>
		<!-- Backdrop -->
		<div
			class="drawer-backdrop"
			:class="{ 'is-active': visible }"
			@click="handleClose"
		></div>

		<!-- Slide Content -->
		<div class="drawer-content" :class="{ 'is-active': visible }">
			<header class="drawer-header">
				<h2 class="drawer-title">Hồ sơ chi tiết</h2>
				<button class="drawer-close" @click="handleClose" aria-label="Đóng panel">
					&times;
				</button>
			</header>

			<div class="drawer-body" v-if="employee">
				<!-- Header section with avatar -->
				<div class="detail-header">
					<div class="detail-avatar">
						{{ getInitials(employee.first_name, employee.last_name) }}
					</div>
					<div class="detail-title-info">
						<h3 class="detail-name">
							{{ employee.first_name }} {{ employee.last_name }}
						</h3>
						<span class="detail-position">{{
							employee.position || "Chưa có chức vụ"
						}}</span>
					</div>
				</div>

				<!-- Section: Personal Info -->
				<div class="detail-section">
					<h4 class="section-title">Thông tin cá nhân</h4>
					<div class="detail-list">
						<div class="detail-item">
							<span class="detail-label">Giới tính:</span>
							<span class="detail-val">{{
								employee.gender === "male"
									? "Nam"
									: employee.gender === "female"
										? "Nữ"
										: "Khác"
							}}</span>
						</div>
						<div class="detail-item">
							<span class="detail-label">Ngày sinh:</span>
							<span class="detail-val">{{
								formatDate(employee.birth_date)
							}}</span>
						</div>
						<div class="detail-item">
							<span class="detail-label">Số điện thoại:</span>
							<span class="detail-val">{{
								employee.phone || "—"
							}}</span>
						</div>
					</div>
				</div>

				<!-- Section: Work Info -->
				<div class="detail-section">
					<h4 class="section-title">Thông tin công việc</h4>
					<div class="detail-list">
						<div class="detail-item">
							<span class="detail-label">Phòng ban:</span>
							<span class="detail-val">{{
								employee.department?.name || "N/A"
							}}</span>
						</div>
						<div class="detail-item">
							<span class="detail-label">Ngày vào làm:</span>
							<span class="detail-val">{{
								formatDate(employee.join_date)
							}}</span>
						</div>
						<div class="detail-item">
							<span class="detail-label">Mức lương:</span>
							<span class="detail-val salary-text">{{
								formatCurrency(employee.salary)
							}}</span>
						</div>
						<div class="detail-item">
							<span class="detail-label">Trạng thái:</span>
							<span
								:class="[
									'status-badge',
									`status-badge--${employee.status}`,
								]"
							>
								{{ formatStatus(employee.status) }}
							</span>
						</div>
					</div>
				</div>

				<!-- Section: Linked Account -->
				<div class="detail-section">
					<h4 class="section-title">Tài khoản liên kết</h4>
					<div class="detail-list" v-if="employee.user">
						<div class="detail-item">
							<span class="detail-label">Email:</span>
							<span class="detail-val">{{
								employee.user.email
							}}</span>
						</div>
						<div class="detail-item">
							<span class="detail-label">Username:</span>
							<span class="detail-val">{{
								employee.user.user_name
							}}</span>
						</div>
						<div class="detail-item">
							<span class="detail-label">Vai trò:</span>
							<span class="detail-val text-uppercase">{{
								employee.user.role?.name || "Employee"
							}}</span>
						</div>
					</div>
					<div class="empty-relation" v-else>
						Nhân viên này chưa liên kết tài khoản hệ thống.
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
.detail-header {
	display: flex;
	align-items: center;
	gap: var(--space-3);
	margin-bottom: var(--space-4);
	padding-bottom: var(--space-3);
	border-bottom: 1px solid var(--border-color);
}

.detail-avatar {
	width: 64px;
	height: 64px;
	border-radius: var(--radius-lg);
	background: var(--primary-gradient);
	color: white;
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: var(--fs-xl);
	font-weight: var(--fw-bold);
	box-shadow: 0 4px 10px rgba(66, 97, 237, 0.2);
}

.detail-title-info {
	display: flex;
	flex-direction: column;
}

.detail-name {
	font-size: var(--fs-lg);
	font-weight: var(--fw-bold);
	margin: 0;
}

.detail-position {
	font-size: var(--fs-sm);
	color: var(--text-muted);
}

.detail-section {
	margin-bottom: var(--space-4);
}

.section-title {
	font-size: var(--fs-sm);
	font-weight: var(--fw-bold);
	text-transform: uppercase;
	color: var(--text-muted);
	border-bottom: 1px solid var(--border-color);
	padding-bottom: 4px;
	margin-top: 0;
	margin-bottom: var(--space-2);
}

.detail-list {
	display: flex;
	flex-direction: column;
	gap: var(--space-2);
}

.detail-item {
	display: flex;
	justify-content: space-between;
	align-items: center;
	font-size: var(--fs-base);
}

.detail-label {
	color: var(--text-muted);
}

.detail-val {
	font-weight: var(--fw-medium);
	color: var(--text-main);
}

.salary-text {
	color: var(--success-color);
	font-weight: var(--fw-bold);
}

.text-uppercase {
	text-transform: uppercase;
}

.empty-relation {
	font-size: var(--fs-sm);
	color: var(--text-light);
	font-style: italic;
	text-align: center;
	padding: var(--space-2) 0;
}
</style>
