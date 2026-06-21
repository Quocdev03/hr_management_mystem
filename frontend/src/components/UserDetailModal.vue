<script setup>
import { computed } from "vue";
import ModalDialog from "./ModalDialog.vue";
import { User, Calendar } from "@lucide/vue";
import { formatDate } from "@/helpers/formatters";

const props = defineProps({
	visible: { type: Boolean, required: true },
	user: { type: Object, default: null },
});

const emit = defineEmits(["close"]);

const handleClose = () => {
	emit("close");
};

const formatStatus = (isActive) => {
	return isActive ? "Đang hoạt động" : "Bị khóa";
};

const userInitial = computed(() => {
	if (!props.user?.user_name) return "U";
	const parts = props.user.user_name.trim().split(/[\s._-]+/);
	if (parts.length >= 2) {
		return (parts[0].charAt(0) + parts[1].charAt(0)).toUpperCase();
	}
	return props.user.user_name.substring(0, 2).toUpperCase();
});
</script>

<template>
	<ModalDialog
		:visible="visible"
		title="Chi tiết Tài khoản"
		size="md"
		@close="handleClose"
	>
		<div class="detail-body" v-if="user">
			<div class="detail-bento-grid">
				<!-- Left Column: User Summary Card -->
				<div class="profile-summary-card">
					<div class="profile-avatar">
						{{ userInitial }}
					</div>
					<h3 class="profile-name">{{ user.user_name }}</h3>
					<span class="profile-position">{{
						user.role?.name || "Chưa cấp quyền"
					}}</span>

					<div class="profile-badges">
						<span
							:class="[
								'status-badge',
								user.is_active
									? 'status-badge--active'
									: 'status-badge--inactive',
							]"
						>
							{{ formatStatus(user.is_active) }}
						</span>
					</div>
				</div>

				<!-- Right Column: User Details Panels -->
				<div class="profile-details-sections">
					<!-- Account Info Section -->
					<div class="info-section">
						<h4 class="section-title">
							<User class="section-title-icon" />
							Thông tin tài khoản
						</h4>
						<div class="detail-grid">
							<div class="detail-item">
								<span class="detail-label">Email</span>
								<span class="detail-val email-val">{{
									user.email
								}}</span>
							</div>
							<div class="detail-item">
								<span class="detail-label">Vai trò</span>
								<span class="detail-val role-badge">{{
									user.role?.name || "—"
								}}</span>
							</div>
						</div>
					</div>

					<!-- Activity Section -->
					<div class="info-section">
						<h4 class="section-title">
							<Calendar class="section-title-icon" />
							Hoạt động
						</h4>
						<div class="detail-grid">
							<div class="detail-item">
								<span class="detail-label"
									>Ngày tạo tài khoản</span
								>
								<span class="detail-val">{{
									formatDate(user.created_at) || "—"
								}}</span>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>

		<template #footer>
			<button
				class="btn btn-secondary close-modal-btn"
				@click="handleClose"
			>
				Đóng
			</button>
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
	word-break: break-all;
}

.profile-position {
	font-size: 0.775rem;
	color: #64748b;
	font-weight: 500;
	text-transform: uppercase;
}

.profile-badges {
	display: flex;
	flex-direction: column;
	gap: 8px;
	width: 100%;
	margin-top: 20px;
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
