import { ref } from "vue";

/**
 * Quản lý trạng thái modal thêm/sửa.
 * Dùng chung cho EmployeeView, DepartmentView, v.v.
 *
 * @returns {{
 *   isModalVisible: import('vue').Ref<boolean>,
 *   isEditMode: import('vue').Ref<boolean>,
 *   openAddModal: () => void,
 *   openEditModal: (record: object) => void,
 *   closeModal: () => void,
 * }}
 */
export function useModalState() {
	const isModalVisible = ref(false);
	const isEditMode = ref(false);

	function openAddModal() {
		isEditMode.value = false;
		isModalVisible.value = true;
	}

	function openEditModal() {
		isEditMode.value = true;
		isModalVisible.value = true;
	}

	function closeModal() {
		isModalVisible.value = false;
	}

	return {
		isModalVisible,
		isEditMode,
		openAddModal,
		openEditModal,
		closeModal,
	};
}
