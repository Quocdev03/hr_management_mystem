<script setup>
import { computed } from "vue";
import ModalDialog from "./ModalDialog.vue";
import {
	Phone,
	Calendar,
	User,
	Mail,
	Briefcase,
	DollarSign,
	Shield,
	UserCheck
} from "@lucide/vue";
import {
	getInitials,
	formatDate,
	formatStatus,
	formatCurrency,
	formatGender,
} from "@/helpers/formatters";

const props = defineProps({
	visible: { type: Boolean, required: true },
	employee: { type: Object, default: null },
});

const emit = defineEmits(["close"]);

const handleClose = () => {
	emit("close");
};
</script>

<template>
	<ModalDialog
		:visible="visible"
		title="Hồ sơ chi tiết nhân viên"
		size="md"
		@close="handleClose"
	>
		<div class="detail-body" v-if="employee">
			<div class="detail-bento-grid">
				<!-- Left Column: Dossier summary card -->
				<div class="profile-summary-card">
					<div class="profile-avatar">
						{{ getInitials(employee.first_name, employee.last_name) }}
					</div>
					<h3 class="profile-name">
						{{ employee.first_name }} {{ employee.last_name }}
					</h3>
					<span class="profile-position">{{ employee.position?.name || "Chưa có chức vụ" }}</span>
					
					<div class="profile-badges">
						<span class="badge badge--primary dept-badge">
							{{ employee.department?.name || "N/A" }}
						</span>
						<span :class="['status-badge', `status-badge--${employee.status}`]">
							{{ formatStatus(employee.status) }}
						</span>
					</div>
				</div>

				<!-- Right Column: Details grouped panels -->
				<div class="profile-details-sections">
					<!-- Personal Info Section -->
					<div class="info-section">
						<h4 class="section-title">
							<User class="section-title-icon" />
							Thông tin cá nhân
						</h4>
						<div class="detail-grid">
							<div class="detail-item">
								<span class="detail-label">Giới tính</span>
								<span class="detail-val">{{ formatGender(employee.gender) }}</span>
							</div>
							<div class="detail-item">
								<span class="detail-label">Ngày sinh</span>
								<span class="detail-val">{{ formatDate(employee.birth_date) }}</span>
							</div>
							<div class="detail-item">
								<span class="detail-label">Điện thoại</span>
								<span class="detail-val">{{ employee.phone || "—" }}</span>
							</div>
						</div>
					</div>

					<!-- Work/Employment Info Section -->
					<div class="info-section">
						<h4 class="section-title">
							<Briefcase class="section-title-icon" />
							Công việc & Thu nhập
						</h4>
						<div class="detail-grid">
							<div class="detail-item">
								<span class="detail-label">Ngày vào làm</span>
								<span class="detail-val">{{ formatDate(employee.join_date) }}</span>
							</div>
							<div class="detail-item">
								<span class="detail-label">Mức lương</span>
								<span class="detail-val salary-text">{{ formatCurrency(employee.salary) }}</span>
							</div>
						</div>
					</div>

					<!-- Account Connection Section -->
					<div class="info-section">
						<h4 class="section-title">
							<Shield class="section-title-icon" />
							Tài khoản hệ thống
						</h4>
						<div class="detail-grid" v-if="employee.user">
							<div class="detail-item">
								<span class="detail-label">Email</span>
								<span class="detail-val email-val">{{ employee.user.email }}</span>
							</div>
							<div class="detail-item">
								<span class="detail-label">Tên tài khoản</span>
								<span class="detail-val">{{ employee.user.user_name }}</span>
							</div>
							<div class="detail-item">
								<span class="detail-label">Vai trò</span>
								<span class="detail-val role-badge">{{ employee.user.role?.name || "Employee" }}</span>
							</div>
						</div>
						<div class="empty-relation-box" v-else>
							<svg class="info-icon" xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>
							Chưa liên kết tài khoản hệ thống.
						</div>
					</div>
				</div>
			</div>
		</div>

		<template #footer>
			<button class="btn btn-secondary close-modal-btn" @click="handleClose">Đóng hồ sơ</button>
		</template>
	</ModalDialog>
</template>

<style scoped>


.detail-bento-grid {
	display: grid;
	grid-template-columns: 210px 1fr;
	gap: 20px;
}

/* Profile Card Summary style */
.profile-summary-card {
	background: #f8fafc;
	border: 1px solid rgba(0, 0, 0, 0.05);
	border-radius: 12px;
	padding: 24px 16px;
	display: flex;
	flex-direction: column;
	align-items: center;
	text-align: center;
	align-self: start;
	box-shadow: inset 0 -1px 0 rgba(0, 0, 0, 0.02);
}

.profile-avatar {
	width: 68px;
	height: 68px;
	border-radius: 14px;
	background: linear-gradient(135deg, #e0e7ff 0%, #c7d2fe 100%);
	color: #4f46e5;
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: 1.5rem;
	font-weight: 700;
	border: 1px solid rgba(79, 70, 229, 0.1);
	box-shadow: 0 4px 10px rgba(79, 70, 229, 0.08);
}

.profile-name {
	margin: 14px 0 4px;
	font-size: 1rem;
	font-weight: 700;
	color: #0f172a;
	line-height: 1.3;
}

.profile-position {
	font-size: 0.775rem;
	color: #64748b;
	font-weight: 500;
}

.profile-badges {
	display: flex;
	flex-direction: column;
	gap: 8px;
	width: 100%;
	margin-top: 20px;
}

.dept-badge {
	text-transform: uppercase;
	font-size: 0.65rem;
	font-weight: 700;
	letter-spacing: 0.03em;
	padding: 4px 8px;
	border-radius: 6px;
	text-align: center;
	width: 100%;
}

.profile-badges .status-badge {
	width: 100%;
}

/* Details list style */
.profile-details-sections {
	display: flex;
	flex-direction: column;
	gap: 16px;
	flex: 1;
}

.info-section {
	display: flex;
	flex-direction: column;
	gap: 8px;
}

.section-title {
	font-size: 0.725rem;
	font-weight: 700;
	color: #4f46e5;
	text-transform: uppercase;
	letter-spacing: 0.05em;
	display: flex;
	align-items: center;
	gap: 6px;
	margin: 0;
}

.section-title-icon {
	width: 13px;
	height: 13px;
	stroke-width: 2.5;
}

.detail-grid {
	background: #ffffff;
	border: 1px solid rgba(0, 0, 0, 0.05);
	border-radius: 10px;
	padding: 10px 14px;
	display: flex;
	flex-direction: column;
	gap: 8px;
}

.detail-item {
	display: flex;
	justify-content: space-between;
	align-items: center;
	font-size: 0.825rem;
	border-bottom: 1px solid rgba(0, 0, 0, 0.03);
	padding-bottom: 6px;
}

.detail-item:last-child {
	border-bottom: none;
	padding-bottom: 0;
}

.detail-label {
	color: #64748b;
	font-size: 0.775rem;
	font-weight: 500;
}

.detail-val {
	color: #1e293b;
	font-weight: 600;
}

.salary-text {
	color: #059669;
	font-weight: 700;
	font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}

.email-val {
	word-break: break-all;
	font-size: 0.8rem;
}

.role-badge {
	text-transform: uppercase;
	background: #f1f5f9;
	padding: 2px 6px;
	border-radius: 4px;
	font-size: 0.7rem;
	color: #475569;
	border: 1px solid rgba(0, 0, 0, 0.04);
}

.empty-relation-box {
	font-size: 0.775rem;
	color: #64748b;
	font-style: italic;
	background: #f8fafc;
	padding: 10px 14px;
	border-radius: 10px;
	border: 1px dashed rgba(0, 0, 0, 0.08);
	display: flex;
	align-items: center;
	gap: 6px;
}

.info-icon {
	color: #94a3b8;
}

.close-modal-btn {
	font-size: 0.825rem;
	font-weight: 600;
	height: 36px;
	padding: 0 1rem;
	border-radius: 8px;
}

@media (max-width: 640px) {
	.detail-bento-grid {
		grid-template-columns: 1fr;
		gap: 16px;
	}
	
	.profile-summary-card {
		align-self: stretch;
	}
}
</style>
