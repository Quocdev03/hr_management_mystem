<script setup>
import { computed } from "vue";
import ModalDialog from "./ModalDialog.vue";
import { User, Mail, Shield, Activity, Calendar, Clock } from "@lucide/vue";
import { formatDate } from "@/helpers/formatters";

const props = defineProps({
	visible: { type: Boolean, required: true },
	user: { type: Object, default: null },
});

const emit = defineEmits(["close"]);

function handleClose() {
	emit("close");
}

function formatStatus(isActive) {
	return isActive ? "Đang hoạt động" : "Bị khóa";
}

const getRoleBadgeClass = (roleName) => {
	if (!roleName) return "badge--primary";
	const name = roleName.toLowerCase();
	if (name === "admin") return "badge--danger";
	if (name === "hr") return "badge--purple";
	return "badge--primary";
};
</script>

<template>
	<ModalDialog
		:visible="visible"
		title="Chi tiết Tài khoản"
		size="md"
		@close="handleClose"
	>
		<div class="detail-body" v-if="user">
			<!-- Header section with avatar -->
			<div class="detail-header">
				<div class="detail-icon">
					<User class="detail-icon-svg" />
				</div>
				<div class="detail-title-info">
					<h3 class="detail-name">{{ user.user_name }}</h3>
					<span
						:class="[
							'badge',
							getRoleBadgeClass(user.role?.name),
						]"
					>
						{{ user.role?.name || "Chưa cấp quyền" }}
					</span>
				</div>
			</div>

			<!-- Section: Account Info -->
			<div class="detail-section">
				<h4 class="section-title">Thông tin tài khoản</h4>
				<div class="detail-list">
					<div class="detail-item">
						<span class="detail-label"
							><Mail class="icon-sm" /> Email:</span
						>
						<span class="detail-val">{{ user.email }}</span>
					</div>
					<div class="detail-item">
						<span class="detail-label"
							><Shield class="icon-sm" /> Quyền hạn:</span
						>
						<span class="detail-val">{{
							user.role?.name || "—"
						}}</span>
					</div>
					<div class="detail-item">
						<span class="detail-label"
							><Activity class="icon-sm" /> Trạng thái:</span
						>
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
			</div>

			<!-- Section: Activity Info -->
			<div class="detail-section">
				<h4 class="section-title">Hoạt động</h4>
				<div class="detail-list">
					<div class="detail-item">
						<span class="detail-label"
							><Calendar class="icon-sm" /> Ngày tạo:</span
						>
						<span class="detail-val">{{
							formatDate(user.created_at) || "—"
						}}</span>
					</div>
				</div>
			</div>
		</div>

		<template #footer>
			<div class="modal-footer">
				<button
					type="button"
					class="btn btn-secondary"
					@click="handleClose"
				>
					Đóng
				</button>
			</div>
		</template>
	</ModalDialog>
</template>

<style scoped>
.section-title {
	font-size: var(--fs-base);
	font-weight: var(--fw-semibold);
	color: var(--text-main);
	margin: 0;
	display: flex;
	align-items: center;
	gap: var(--space-1);
}

.detail-list {
	display: flex;
	flex-direction: column;
	gap: 0;
	background: var(--bg-lighter);
	border-radius: var(--radius-md);
	border: 1px solid var(--border-color);
	overflow: hidden;
}

.detail-item {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: var(--space-2) var(--space-3);
	border-bottom: 1px solid var(--border-color);
}

.detail-item:last-child {
	border-bottom: none;
}

.detail-label {
	display: flex;
	align-items: center;
	gap: 6px;
	font-size: 13px;
	color: var(--text-muted);
	font-weight: var(--fw-medium);
}

.detail-val {
	font-size: var(--fs-sm);
	color: var(--text-main);
	font-weight: var(--fw-semibold);
	text-align: right;
}

.icon-sm {
	width: 14px;
	height: 14px;
	color: var(--text-light);
}
</style>
