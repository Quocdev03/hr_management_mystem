/**
 * Lấy 2 chữ cái đầu của họ và tên để hiển thị avatar.
 * @param {string} firstName
 * @param {string} lastName
 * @returns {string}
 */
export function getInitials(firstName, lastName) {
	let first = "";
	if (firstName) {
		first = firstName.charAt(0);
	}

	let last = "";
	if (lastName) {
		last = lastName.charAt(0);
	}

	return (first + last).toUpperCase();
}

/**
 * Định dạng chuỗi ngày ISO thành ngày theo locale vi-VN.
 * @param {string|null} dateString
 * @returns {string}
 */
export function formatDate(dateString) {
	if (!dateString) {
		return "—";
	}
	let date = new Date(dateString);
	return date.toLocaleDateString("vi-VN");
}

/**
 * Map trạng thái nhân viên sang tiếng Việt.
 * @param {string} status
 * @returns {string}
 */
export function formatStatus(status) {
	if (status === "active") {
		return "Đang làm việc";
	}
	if (status === "inactive") {
		return "Đã nghỉ việc";
	}
	return status;
}
export function formatCurrency(amount) {
	if (!amount) return "N/A";
	return new Intl.NumberFormat("vi-VN", {
		style: "currency",
		currency: "VND",
	}).format(amount);
}

/**
 * Map giới tính sang tiếng Việt.
 * @param {string} gender
 * @returns {string}
 */
export function formatGender(gender) {
	if (gender === "male") {
		return "Nam";
	}
	if (gender === "female") {
		return "Nữ";
	}
	return "Khác";
}
