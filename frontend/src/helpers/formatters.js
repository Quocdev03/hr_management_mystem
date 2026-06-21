/**
 * Lấy 2 chữ cái đầu của họ và tên để hiển thị avatar.
 * @param {string} firstName
 * @param {string} lastName
 * @returns {string}
 */
export const getInitials = (firstName, lastName) => {
	const first = firstName?.charAt(0) ?? "";
	const last = lastName?.charAt(0) ?? "";
	return `${first}${last}`.toUpperCase();
};

/**
 * Định dạng chuỗi ngày ISO thành ngày theo locale vi-VN.
 * @param {string|null} dateString
 * @returns {string}
 */
export const formatDate = (dateString) => {
	if (!dateString) return "—";
	return new Date(dateString).toLocaleDateString("vi-VN");
};

/**
 * Map trạng thái nhân viên sang tiếng Việt.
 * @param {string} status
 * @returns {string}
 */
export const formatStatus = (status) => {
	const map = {
		active: "Đang làm việc",
		inactive: "Đã nghỉ việc",
	};
	return map[status] ?? status;
};

/**
 * Định dạng tiền tệ sang VND.
 * @param {number} amount
 * @returns {string}
 */
export const formatCurrency = (amount) => {
	if (!amount) return "N/A";
	return new Intl.NumberFormat("vi-VN", {
		style: "currency",
		currency: "VND",
	}).format(amount);
};

/**
 * Map giới tính sang tiếng Việt.
 * @param {string} gender
 * @returns {string}
 */
export const formatGender = (gender) => {
	const map = {
		male: "Nam",
		female: "Nữ",
	};
	return map[gender] ?? "Khác";
};
