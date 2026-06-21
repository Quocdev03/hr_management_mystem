import { ref } from "vue";

/**
 * Quản lý trạng thái hiển thị và chế độ (thêm/sửa) của modal.
 * Dùng chung cho EmployeeView, DepartmentView, v.v.
 */
export const useModalState = () => {
	const isModalVisible = ref(false);
	const isEditMode = ref(false);

	const openAddModal = () => {
		isEditMode.value = false;
		isModalVisible.value = true;
	};

	const openEditModal = () => {
		isEditMode.value = true;
		isModalVisible.value = true;
	};

	const closeModal = () => {
		isModalVisible.value = false;
	};

	return {
		isModalVisible,
		isEditMode,
		openAddModal,
		openEditModal,
		closeModal,
	};
};
