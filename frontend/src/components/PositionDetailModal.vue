<script setup>
import { computed } from "vue";
import ModalDialog from "./ModalDialog.vue";
import { getInitials, formatDate } from "@/helpers/formatters";
import { Briefcase, Calendar, Users } from "@lucide/vue";

const props = defineProps({
	visible: { type: Boolean, required: true },
	position: { type: Object, default: null },
	employees: { type: Array, default: () => [] }
});

const emit = defineEmits(["close"]);

const handleClose = () => {
	emit("close");
};
</script>

<template>
	<ModalDialog
		:visible="visible"
		title="Chi tiết chức vụ"
		size="md"
		@close="handleClose"
	>
		<div class="detail-body" v-if="position">
			<!-- Position Header Info -->
			<div class="detail-header">
				<div class="detail-icon">
					<Briefcase class="detail-icon-svg" />
				</div>
				<div class="detail-title-info">
					<h3 class="detail-name">
						{{ position.name }}
					</h3>
					<span class="detail-code">ID chức vụ: #{{ position.id }}</span>
				</div>
			</div>

			<!-- Bento Info Grid for metadata -->
			<div class="detail-section">
				<h4 class="section-title">
					<Calendar class="section-title-icon" />
					Thông tin hệ thống
				</h4>
				<div class="metadata-grid">
					<div class="metadata-item">
						<span class="metadata-label">Ngày tạo</span>
						<span class="metadata-value">{{ formatDate(position.created_at) || "—" }}</span>
					</div>
					<div class="metadata-item">
						<span class="metadata-label">Cập nhật lần cuối</span>
						<span class="metadata-value">{{ formatDate(position.updated_at) || "—" }}</span>
					</div>
				</div>
			</div>

			<!-- Description Section -->
			<div class="detail-section">
				<h4 class="section-title">Mô tả chi tiết</h4>
				<div class="description-box">
					{{ position.description || "Không có mô tả chi tiết cho chức vụ này." }}
				</div>
			</div>

			<!-- Employee List Section -->
			<div class="detail-section">
				<h4 class="section-title">
					<Users class="section-title-icon" />
					Nhân sự đảm nhiệm ({{ employees.length }})
				</h4>
				<div class="employee-table-wrapper" v-if="employees.length > 0">
					<table class="data-table">
						<thead>
							<tr>
								<th>Nhân viên</th>
								<th>Phòng ban</th>
								<th>Điện thoại</th>
								<th>Trạng thái</th>
							</tr>
						</thead>
						<tbody>
							<tr v-for="emp in employees" :key="emp.id">
								<td class="emp-cell">
									<div class="avatar-gradient" style="width: 28px; height: 28px; font-size: 11px;">
										{{
											getInitials(
												emp.first_name,
												emp.last_name,
											)
										}}
									</div>
									<div class="emp-info">
										<span class="emp-name"
											>{{ emp.first_name }}
											{{ emp.last_name }}</span
										>
									</div>
								</td>
								<td>{{ emp.department?.name || "—" }}</td>
								<td>{{ emp.phone || "—" }}</td>
								<td>
									<span
										:class="[
											'status-badge',
											`status-badge--${emp.status}`,
										]"
									>
										{{ emp.status === "active" ? "Đang làm việc" : "Nghỉ việc" }}
									</span>
								</td>
							</tr>
						</tbody>
					</table>
				</div>
				<div class="empty-employees" v-else>
					Chưa có nhân viên nào đảm nhiệm chức vụ này.
				</div>
			</div>
		</div>

		<template #footer>
			<button class="btn btn-secondary close-btn" @click="handleClose">Đóng</button>
		</template>
	</ModalDialog>
</template>

<style scoped>
.metadata-grid {
	display: grid;
	grid-template-columns: 1fr 1fr;
	gap: 12px;
	background: var(--bg-lighter);
	padding: 12px 14px;
	border-radius: var(--radius-md);
	border: 1px solid var(--border-color);
}

.metadata-item {
	display: flex;
	flex-direction: column;
	gap: 4px;
}

.metadata-label {
	font-size: var(--fs-xs);
	color: var(--text-muted);
	font-weight: 500;
}

.metadata-value {
	font-size: var(--fs-sm);
	color: var(--text-main);
	font-weight: 600;
}

.description-box {
	font-size: var(--fs-sm);
	color: var(--text-main);
	background: var(--bg-lighter);
	padding: 10px 14px;
	border-radius: var(--radius-md);
	border: 1px solid var(--border-color);
	line-height: var(--lh-normal);
}

.section-title-icon {
	width: 14px;
	height: 14px;
	margin-right: 4px;
	color: var(--primary-color);
	display: inline-block;
	vertical-align: middle;
}

.empty-employees {
	padding: 12px 14px;
	font-size: var(--fs-sm);
	color: var(--text-muted);
	background: var(--bg-lighter);
	border-radius: var(--radius-md);
	border: 1px dashed var(--border-color);
	text-align: center;
}

/* Employee Table Wrapper */
.employee-table-wrapper {
	max-height: 220px;
	overflow-y: auto;
	border: 1px solid var(--border-color);
	border-radius: var(--radius-md);
	background: var(--bg-card);
}

.emp-cell {
	display: flex;
	align-items: center;
	gap: 10px;
}

.emp-info {
	display: flex;
	flex-direction: column;
}

.emp-name {
	font-weight: var(--fw-semibold);
	font-size: var(--fs-sm);
}

.close-btn {
	font-size: 0.825rem;
	font-weight: 600;
	height: 36px;
	padding: 0 1rem;
	border-radius: 8px;
}

.employee-table-wrapper::-webkit-scrollbar {
	width: 6px;
}
.employee-table-wrapper::-webkit-scrollbar-track {
	background: transparent;
}
.employee-table-wrapper::-webkit-scrollbar-thumb {
	background-color: var(--border-hover);
	border-radius: 20px;
}
</style>
