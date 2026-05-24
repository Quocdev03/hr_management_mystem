<script setup>
	import deleteIcon from "@/assets/svg/warning.svg";
	import editIcon from "@/assets/svg/edit.svg";
	import searchIcon from "@/assets/svg/search.svg";
	import plusIcon from "@/assets/svg/plus.svg";
	import prevIcon from "@/assets/svg/chevron-left.svg";
	import nextIcon from "@/assets/svg/chevron-right.svg";
	import ConfirmationDialog from "@/components/ConfirmationDialog.vue";

	import { useUserStore } from "@/store/user";
	import { useToast } from "vue-toastification";
	import { storeToRefs } from "pinia";
	import { usePaginatedSearch } from "@/helpers/usePaginatedSearch";
	import { onMounted, ref } from "vue";
	const userStore = useUserStore();
	const toast = useToast();

	const { users, pagination, loading } = storeToRefs(userStore);

	// Search & pagination
	const {
		searchQuery,
		load: loadUsers,
		handlePageChange,
	} = usePaginatedSearch((params) => userStore.fetchUser(params), pagination);

	// Delete modal state
	const isDeleteModalVisible = ref(false);
	const deletingUser = ref(null);
	const deleteMessage = ref("");
	const deleteLoading = ref(false);

	function handleAdd() {}
	function handleUpdate() {}

	function handleDelete(user) {
		deletingUser.value = user;
		deleteMessage.value =
			"Bạn có chắc chắn muốn xoá người dùng " + user.user_name + "?";
		isDeleteModalVisible.value = true;
	}

	async function confirmDelete() {
		let user = deletingUser.value;
		if (!user) return;

		deleteLoading.value = true;
		const res = await userStore.deleteUser(user.id);
		deleteLoading.value = false;

		if (res.success === false) {
			if (res.message) {
				toast.error(res.message);
			} else {
				toast.error("Xoá người dùng thất bại");
			}
			return;
		}

		toast.success("Xoá người dùng thành công");
		isDeleteModalVisible.value = false;
		deletingUser.value = null;
		await loadUsers();
	}

	onMounted(async function () {
		await loadUsers();
	});
</script>
<template>
	<div class="user-view">
		<header class="page-header">
			<div class="header-content">
				<h1 class="page-title">Quản lý người dùng</h1>
				<p class="page-subtitle">
					Hệ thống có tổng cộng
					<span>{{ pagination.total }}</span> người dùng
				</p>
			</div>
			<button class="btn btn--primary">
				<img :src="plusIcon" alt="add" class="btn__icon" />
				Thêm người dùng
			</button>
		</header>

		<main class="content-card">
			<div class="toolbar">
				<div class="search-box">
					<img :src="searchIcon" class="search-box__icon" alt="search" />
					<input
						v-model="searchQuery"
						class="form-control search-box__input"
						placeholder="Tìm tên hoặc mã người dùng..."
					/>
				</div>
			</div>

			<div v-if="loading" class="table-loading">Đang tải dữ liệu...</div>
			<div v-else class="table-responsive">
				<table class="data-table">
					<thead>
						<tr>
							<th>Tên người dùng</th>
							<th>Email</th>
							<th>Trạng thái</th>
							<th class="text-right">Thao tác</th>
						</tr>
					</thead>
					<tbody>
						<tr v-for="user in users" :key="user.id">
							<td class="text-main fw-500">{{ user.user_name }}</td>
							<td>
								<span class="">{{ user.email }}</span>
							</td>
							<td>
								<span
									:class="[
										'status-badge',
										user.is_active
											? 'status-badge--active'
											: 'status-badge--inactive',
									]"
								>
									{{ user.is_active ? "Hoạt động" : "Ngưng" }}
								</span>
							</td>
							<td class="text-right">
								<div class="action-group">
									<button
										class="btn-icon btn-icon--edit"
										title="Chỉnh sửa"
										@click=""
									>
										<img :src="editIcon" alt="edit" />
									</button>
									<button
										class="btn-icon btn-icon--delete"
										title="Xoá"
										@click="handleDelete(user)"
									>
										<img :src="deleteIcon" alt="delete" />
									</button>
								</div>
							</td>
						</tr>
						<tr v-if="users.length === 0">
							<td colspan="4" class="empty-state">
								<div class="empty-state__icon">🏢</div>
								<p class="empty-state__text">Không có người dùng nào phù hợp.</p>
							</td>
						</tr>
					</tbody>
				</table>
			</div>

			<!-- Phân trang -->
			<div class="pagination" v-if="pagination.totalPages > 0">
				<button
					class="pagination__btn"
					:disabled="pagination.page === 1"
					@click="handlePageChange(pagination.page - 1)"
				>
					<img :src="prevIcon" alt="prev" />
				</button>
				<div class="pagination__info">
					Trang <span>{{ pagination.page }}</span> / {{ pagination.totalPages }}
				</div>
				<button
					class="pagination__btn"
					:disabled="pagination.page === pagination.totalPages"
					@click="handlePageChange(pagination.page + 1)"
				>
					<img :src="nextIcon" alt="next" />
				</button>
			</div>
		</main>

		<ConfirmationDialog
			:visible="isDeleteModalVisible"
			title="Xác nhận xoá người dùng"
			:message="deleteMessage"
			:loading="deleteLoading"
			@confirm="confirmDelete"
			@cancel="isDeleteModalVisible = false"
		/>
	</div>
</template>

<style scoped>
	.user-view {
		padding-bottom: var(--space-4);
	}

	.page-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: var(--space-4);
		flex-wrap: wrap;
		gap: var(--space-2);
	}

	.page-title {
		font-size: var(--fs-2xl);
		font-weight: var(--fw-bold);
		letter-spacing: var(--tracking-tight);
		margin: 0 0 4px 0;
		color: var(--text-main);
	}

	.page-subtitle {
		color: var(--text-muted);
		font-size: var(--fs-sm);
		margin: 0;
	}

	.page-subtitle span {
		font-weight: var(--fw-semibold);
		color: var(--primary-color);
	}

	/* ===== Card & Toolbar ===== */
	.content-card {
		background: var(--bg-card);
		border-radius: var(--radius-lg);
		border: 1px solid var(--border-color);
		box-shadow: var(--shadow-sm);
		overflow: hidden;
	}

	.toolbar {
		display: flex;
		justify-content: space-between;
		gap: var(--space-3);
		padding: var(--space-3) var(--space-4);
		border-bottom: 1px solid var(--border-color);
		flex-wrap: wrap;
		background: var(--bg-card);
		align-items: center;
	}

	.search-box {
		position: relative;
		flex: 1;
		max-width: 400px;
	}

	.search-box__icon {
		position: absolute;
		left: 1rem;
		top: 50%;
		transform: translateY(-50%);
		width: 18px;
		height: 18px;
		filter: grayscale(1) opacity(0.5);
		z-index: 1;
	}

	.search-box__input {
		padding-left: 2.75rem !important;
	}

	/* ===== Table ===== */
	.table-responsive {
		overflow-x: auto;
	}

	.data-table {
		width: 100%;
		border-collapse: collapse;
		text-align: left;
	}

	.data-table th {
		padding: var(--space-3) var(--space-4);
		background: var(--bg-lighter);
		font-size: var(--fs-xs);
		text-transform: uppercase;
		font-weight: var(--fw-bold);
		color: var(--text-muted);
		letter-spacing: 0.05em;
		border-bottom: 1px solid var(--border-color);
	}

	.data-table td {
		padding: var(--space-3) var(--space-2);
		border-bottom: 1px solid var(--border-color);
		vertical-align: middle;
	}

	.data-table tbody tr:hover td {
		background: var(--bg-lighter);
	}

	/* Dept code pill */
	.dept-code {
		display: inline-block;
		padding: 0.2rem 0.6rem;
		background: var(--bg-light);
		border-radius: var(--radius-sm);
		font-size: var(--fs-xs);
		font-weight: var(--fw-semibold);
		color: var(--text-muted);
		letter-spacing: 0.03em;
	}

	/* ===== Actions ===== */
	.action-group {
		display: inline-flex;
		gap: 0.5rem;
	}

	.btn-icon {
		width: 34px;
		height: 34px;
		border-radius: var(--radius-sm);
		background: var(--bg-main);
		border: 1px solid var(--border-color);
		display: inline-flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.btn-icon img {
		width: 16px;
		height: 16px;
		filter: grayscale(1) opacity(0.6);
	}

	.btn-icon:hover {
		border-color: var(--primary-color);
		background: #eff6ff;
	}

	.btn-icon:hover img {
		filter: none;
	}

	.btn-icon--delete:hover {
		border-color: var(--danger-color);
		background: #fee2e2;
	}

	/* ===== Pagination ===== */
	.pagination {
		display: flex;
		justify-content: center;
		align-items: center;
		gap: var(--space-3);
		padding: var(--space-3) var(--space-4);
		background: var(--bg-card);
		border-top: 1px solid var(--border-color);
	}

	.pagination__btn {
		width: 36px;
		height: 36px;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		border: 1px solid var(--border-color);
		background: var(--bg-main);
		border-radius: var(--radius-md);
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.pagination__btn:hover:not(:disabled) {
		border-color: var(--primary-color);
		background: var(--bg-light);
	}

	.pagination__btn:disabled {
		opacity: 0.35;
		cursor: not-allowed;
	}

	.pagination__btn img {
		width: 18px;
		height: 18px;
		opacity: 0.7;
	}

	.pagination__info {
		font-size: var(--fs-sm);
		color: var(--text-muted);
	}

	.pagination__info span {
		font-weight: var(--fw-bold);
		color: var(--text-main);
	}

	/* ===== Empty state ===== */
	.empty-state {
		text-align: center;
		padding: var(--space-8) 0;
	}

	.empty-state__icon {
		font-size: 2.5rem;
		margin-bottom: var(--space-3);
	}

	.empty-state__text {
		color: var(--text-muted);
		font-size: var(--fs-sm);
	}

	/* Status Badge */
	.status-badge {
		padding: 0.25rem 0.75rem;
		border-radius: var(--radius-sm);
		font-size: var(--fs-xs);
		font-weight: var(--fw-bold);
		display: inline-block;
	}

	.status-badge--active {
		background: #dcfce7;
		color: #15803d;
	}

	.status-badge--inactive {
		background: #fee2e2;
		color: #b91c1c;
	}

	/* ===== Misc ===== */
	.btn__icon {
		width: 18px;
		height: 18px;
		filter: brightness(0) invert(1);
	}

	@media (max-width: 768px) {
		.toolbar {
			flex-direction: column;
			align-items: stretch;
		}
		.search-box {
			max-width: none;
		}
	}
</style>
