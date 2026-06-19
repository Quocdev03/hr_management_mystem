<script setup>
import { computed } from "vue";
import ModalDialog from "./ModalDialog.vue";
import { getInitials } from "@/helpers/formatters";
import { Building2 } from "@lucide/vue";

const props = defineProps({
	visible: { type: Boolean, required: true },
	department: { type: Object, default: null }
});

const emit = defineEmits(["close"]);

function handleClose() {
	emit("close");
}

const employees = computed(() => props.department?.employees || []);
const managerName = computed(() => {
	if (!props.department?.manager) return "Chưa bổ nhiệm";
	const manager = props.department.manager;
	return `${manager.first_name} ${manager.last_name}`;
});
</script>

<template>
	<ModalDialog
		:visible="visible"
		title="Chi tiết phòng ban"
		size="lg"
		@close="handleClose"
	>
		<div class="detail-body" v-if="department">
			<!-- Department Header Info -->
			<div class="detail-header">
				<div class="detail-icon">
					<Building2 class="detail-icon-svg" />
				</div>
				<div class="detail-title-info">
					<h3 class="detail-name">
						{{ department.name }}
					</h3>
					<span class="detail-code">Mã phòng ban: {{ department.code }}</span>
				</div>
			</div>

			<!-- Description Section -->
			<div class="detail-section" v-if="department.description">
				<h4 class="section-title">Mô tả chi tiết</h4>
				<div class="description-box">
					{{ department.description }}
				</div>
			</div>

			<!-- Manager Card -->
			<div class="detail-section manager-section">
				<h4 class="section-title">Trưởng phòng ban</h4>
				<div class="manager-card" v-if="department.manager">
					<div class="manager-avatar">
						{{ getInitials(department.manager.first_name, department.manager.last_name) }}
					</div>
					<div class="manager-details">
						<span class="manager-name-text">{{ managerName }}</span>
						<span class="manager-subtext">{{ department.manager.position?.name || 'Trưởng phòng' }}</span>
						<span class="manager-contact" v-if="department.manager.phone">SĐT: {{ department.manager.phone }}</span>
					</div>
				</div>
				<div class="empty-manager" v-else>
					Chưa chỉ định trưởng phòng cho phòng ban này.
				</div>
			</div>


			<!-- Employee List Section -->
			<div class="detail-section">
				<h4 class="section-title">Thành viên phòng ban ({{ employees.length }})</h4>
				<div class="employee-table-wrapper" v-if="employees.length > 0">
					<table class="employee-table">
						<thead>
							<tr>
								<th>Nhân viên</th>
								<th>Chức vụ</th>
								<th>Điện thoại</th>
								<th>Trạng thái</th>
							</tr>
						</thead>
						<tbody>
							<tr v-for="emp in employees" :key="emp.id">
								<td class="emp-cell">
									<div class="emp-avatar">
										{{ getInitials(emp.first_name, emp.last_name) }}
									</div>
									<div class="emp-info">
										<span class="emp-name">{{ emp.first_name }} {{ emp.last_name }}</span>
									</div>
								</td>
								<td>{{ emp.position?.name || '—' }}</td>
								<td>{{ emp.phone || '—' }}</td>
								<td>
									<span :class="['status-dot', emp.status === 'active' ? 'active' : 'inactive']"></span>
									<span class="status-text">{{ emp.status === 'active' ? 'Đang làm việc' : 'Nghỉ việc' }}</span>
								</td>
							</tr>
						</tbody>
					</table>
				</div>
				<div class="empty-employees" v-else>
					Phòng ban này hiện chưa có nhân viên nào.
				</div>
			</div>
		</div>
		
		<template #footer>
			<button class="btn btn--secondary" @click="handleClose">Đóng</button>
		</template>
	</ModalDialog>
</template>

<style scoped>
.detail-body {
	display: flex;
	flex-direction: column;
	gap: var(--space-3);
}

.detail-header {
	display: flex;
	align-items: center;
	gap: var(--space-3);
	padding-bottom: var(--space-3);
	border-bottom: 1px solid var(--border-color);
}

.detail-icon {
	width: 56px;
	height: 56px;
	background: var(--primary-gradient);
	border-radius: var(--radius-md);
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: var(--fs-xl);
	box-shadow: 0 4px 12px rgba(66, 97, 237, 0.2);
}

.detail-icon-svg {
	width: 26px;
	height: 26px;
	color: white;
}

.detail-title-info {
	display: flex;
	flex-direction: column;
	gap: 2px;
}

.detail-name {
	font-size: var(--fs-lg);
	font-weight: 800;
	color: var(--text-main);
	margin: 0;
}

.detail-code {
	font-size: var(--fs-sm);
	color: var(--text-muted);
	font-weight: var(--fw-medium);
}

.detail-section {
	display: flex;
	flex-direction: column;
	gap: var(--space-2);
}

.section-title {
	font-size: var(--fs-xs);
	font-weight: var(--fw-bold);
	color: var(--text-muted);
	text-transform: uppercase;
	letter-spacing: 0.05em;
	margin: 0 0 2px 0;
}

.description-box {
	font-size: var(--fs-sm);
	color: var(--text-main);
	background: #f8fafc;
	padding: var(--space-2) var(--space-3);
	border-radius: var(--radius-md);
	border: 1px solid var(--border-color);
	line-height: var(--lh-normal);
}

/* Manager Card */
.manager-card {
	display: flex;
	align-items: center;
	gap: var(--space-2);
	background: rgba(66, 97, 237, 0.03);
	padding: var(--space-2) var(--space-3);
	border-radius: var(--radius-md);
	border: 1px solid rgba(66, 97, 237, 0.08);
}

.manager-avatar {
	width: 44px;
	height: 44px;
	border-radius: 8px;
	background: var(--primary-gradient);
	color: white;
	display: flex;
	align-items: center;
	justify-content: center;
	font-weight: 800;
	font-size: var(--fs-sm);
	box-shadow: 0 2px 8px rgba(66, 97, 237, 0.2);
}

.manager-details {
	display: flex;
	flex-direction: column;
	overflow: hidden;
}

.manager-name-text {
	font-size: var(--fs-sm);
	font-weight: var(--fw-bold);
	color: var(--text-main);
}

.manager-subtext {
	font-size: var(--fs-xs);
	color: var(--text-muted);
}

.manager-contact {
	font-size: var(--fs-xs);
	color: var(--primary-color);
	font-weight: var(--fw-medium);
	margin-top: 2px;
}

.empty-manager, .empty-employees {
	padding: var(--space-2) var(--space-3);
	font-size: var(--fs-sm);
	color: var(--text-muted);
	background: #f8fafc;
	border-radius: var(--radius-md);
	border: 1px dashed var(--border-color);
	text-align: center;
}

/* Employee Table */
.employee-table-wrapper {
	max-height: 250px;
	overflow-y: auto;
	border: 1px solid var(--border-color);
	border-radius: var(--radius-md);
	background: white;
}

.employee-table {
	width: 100%;
	border-collapse: collapse;
	text-align: left;
}

.employee-table th {
	padding: 8px 16px;
	background: #f8fafc;
	font-size: var(--fs-xs);
	font-weight: 800;
	color: var(--text-muted);
	text-transform: uppercase;
	letter-spacing: 0.05em;
	border-bottom: 1px solid var(--border-color);
	position: sticky;
	top: 0;
	z-index: 10;
}

.employee-table td {
	padding: 10px 16px;
	border-bottom: 1px solid var(--border-color);
	font-size: var(--fs-sm);
	vertical-align: middle;
	color: var(--text-main);
}

.employee-table tbody tr:last-child td {
	border-bottom: none;
}

.employee-table tbody tr:hover td {
	background: #f4f7ff;
}

.emp-cell {
	display: flex;
	align-items: center;
	gap: 10px;
}

.emp-avatar {
	width: 28px;
	height: 28px;
	border-radius: 6px;
	background: linear-gradient(135deg, rgba(0, 192, 250, 0.1) 0%, rgba(66, 97, 237, 0.08) 100%);
	color: var(--primary-color);
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: 11px;
	font-weight: 800;
}

.emp-info {
	display: flex;
	flex-direction: column;
}

.emp-name {
	font-weight: var(--fw-semibold);
}

/* Status dots */
.status-dot {
	display: inline-block;
	width: 6px;
	height: 6px;
	border-radius: 50%;
	margin-right: 6px;
}

.status-dot.active {
	background: var(--success-color);
	box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.2);
}

.status-dot.inactive {
	background: var(--text-light);
}

.status-text {
	font-size: var(--fs-xs);
	font-weight: var(--fw-medium);
}

</style>
