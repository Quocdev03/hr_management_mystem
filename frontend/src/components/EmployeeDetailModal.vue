<script setup>
import { computed } from "vue";
import ModalDialog from "./ModalDialog.vue";
import {
	getInitials,
	formatDate,
	formatStatus,
	formatCurrency,
	formatGender
} from "@/helpers/formatters";

const props = defineProps({
	visible: { type: Boolean, required: true },
	employee: { type: Object, default: null }
});

const emit = defineEmits(["close"]);

function handleClose() {
	emit("close");
}
</script>

<template>
	<ModalDialog
		:visible="visible"
		title="Hồ sơ chi tiết"
		size="md"
		@close="handleClose"
	>
		<div class="detail-body" v-if="employee">
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
						employee.position?.name || "Chưa có chức vụ"
					}}</span>
				</div>
			</div>

			<!-- Section: Personal Info -->
			<div class="detail-section">
				<h4 class="section-title">Thông tin cá nhân</h4>
				<div class="detail-list">
					<div class="detail-item">
						<span class="detail-label">Giới tính:</span>
						<span class="detail-val">{{ formatGender(employee.gender) }}</span>
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
		
		<template #footer>
			<button class="btn btn-secondary" @click="handleClose">Đóng</button>
		</template>
	</ModalDialog>
</template>

<style scoped>
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

.detail-position {
	font-size: var(--fs-sm);
	color: var(--text-muted);
}.detail-list {
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
