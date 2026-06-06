import { ref, watch } from "vue";

/**
 * Composable xử lý search + phân trang + load dữ liệu.
 * Dùng chung cho EmployeeView, DepartmentView, v.v.
 *
 * @param {(params: { page: number; limit: number; search: string }) => Promise<{ success: boolean; message?: string }>} fetchFn
 *   Hàm fetch từ store, nhận params và trả về response chuẩn { success, message }.
 * @param {import('vue').Ref<{ page: number; limit: number; totalPages: number }>} paginationRef
 *   Ref pagination từ store (storeToRefs).
 * @param {{ debounce?: number }} [options]
 *   - debounce: thời gian debounce cho search (ms, mặc định 400).
 *
 * @returns {{
 *   searchQuery: import('vue').Ref<string>,
 *   load: (page?: number) => Promise<void>,
 *   handlePageChange: (page: number) => void,
 *   errorMessage: import('vue').Ref<string | null>,
 * }}
 */
export function usePaginatedSearch(fetchFn, paginationRef, options) {
	// ─── Cấu hình ────────────────────────────────────────────────────────────

	// Lấy debounce từ options, mặc định 400ms nếu không truyền
	const debounceMs = options?.debounce !== undefined ? options.debounce : 400;

	// ─── State ───────────────────────────────────────────────────────────────

	const searchQuery = ref(""); // Từ khoá tìm kiếm hiện tại
	const errorMessage = ref(null); // Lỗi từ lần fetch gần nhất (null = không có lỗi)
	let _debounceTimer = null; // Timer debounce — giữ ref để clearTimeout khi cần

	// ─── Load dữ liệu ────────────────────────────────────────────────────────

	async function load(page = 1) {
		errorMessage.value = null;

		const params = {
			page,
			limit: paginationRef.value.limit,
			search: searchQuery.value,
		};

		const res = await fetchFn(params);

		/**
		 * Đồng bộ page về paginationRef sau khi fetch thành công.
		 * Ưu tiên page trả về từ response (res.data.page) để đảm bảo khớp
		 * với dữ liệu thực tế từ server, fallback về page đã gửi lên nếu
		 * response không trả về page (một số API không include pagination).
		 */
		if (res?.success !== false) {
			try {
				const respPage = res?.data?.page
					? Number(res.data.page)
					: Number(page);
				if (!Number.isNaN(respPage) && paginationRef?.value) {
					paginationRef.value.page = respPage;
				}
			} catch {
				// Bỏ qua nếu parse page lỗi — không ảnh hưởng đến hiển thị
			}
		}

		// Ghi nhận lỗi để component có thể hiển thị thông báo
		if (res?.success === false) {
			errorMessage.value = res.message || "Lỗi tải dữ liệu";
		}

		return res;
	}

	// ─── Chuyển trang ────────────────────────────────────────────────────────

	/**
	 * Validate page trước khi load:
	 * - Phải là số hợp lệ, >= 1, và không vượt quá totalPages
	 * - Tránh gọi API thừa khi người dùng spam nút prev/next ở biên
	 */
	function handlePageChange(page) {
		const p = Number(page);
		if (Number.isNaN(p) || p < 1) return;
		if (p > paginationRef.value.totalPages) return;
		load(p);
	}

	// ─── Debounce search ─────────────────────────────────────────────────────

	/**
	 * Mỗi lần searchQuery thay đổi, reset timer và đợi debounceMs
	 * trước khi thực sự gọi load(1).
	 * Tránh spam API khi người dùng đang gõ liên tục.
	 */
	watch(searchQuery, () => {
		clearTimeout(_debounceTimer);
		_debounceTimer = setTimeout(() => load(1), debounceMs);
	});

	// ─── Export ──────────────────────────────────────────────────────────────

	return {
		searchQuery,
		load,
		handlePageChange,
		errorMessage,
	};
}
