import { ref } from "vue";

/**
 * Quản lý trạng thái modal thêm/sửa.
 * Dùng chung cho EmployeeView, DepartmentView, v.v.
 *
 * @returns {{
 *   isModalVisible: import('vue').Ref<boolean>,
 *   isEditMode: import('vue').Ref<boolean>,
 *   openAddModal: () => void,
 *   openEditModal: () => void,
 *   closeModal: () => void,
 * }}
 */
export function useModalState() {
	// ─── State ───────────────────────────────────────────────────────────────

	const isModalVisible = ref(false); // Modal có đang mở không
	const isEditMode = ref(false); // true = sửa, false = thêm mới

	// ─── Actions ─────────────────────────────────────────────────────────────

	// Mở modal ở chế độ thêm mới
	function openAddModal() {
		isEditMode.value = false;
		isModalVisible.value = true;
	}

	// Mở modal ở chế độ chỉnh sửa
	function openEditModal() {
		isEditMode.value = true;
		isModalVisible.value = true;
	}

	// Đóng modal — không reset isEditMode để tránh flash UI khi đóng có animation
	function closeModal() {
		isModalVisible.value = false;
	}

	// ─── Export ──────────────────────────────────────────────────────────────

	return {
		isModalVisible,
		isEditMode,
		openAddModal,
		openEditModal,
		closeModal,
	};
}
