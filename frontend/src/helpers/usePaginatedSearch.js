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
	let debounceMs = 400;
	if (options && options.debounce !== undefined) {
		debounceMs = options.debounce;
	}

	const searchQuery = ref("");
	const errorMessage = ref(null);
	let _debounceTimer = null;

	async function load(page) {
		if (page === undefined) {
			page = 1;
		}

		errorMessage.value = null;

		let params = {
			page: page,
			limit: paginationRef.value.limit,
			search: searchQuery.value,
		};

		const res = await fetchFn(params);

		// Nếu fetch thành công, đảm bảo paginationRef có page đúng (fallback)
		if (res && res.success !== false) {
			try {
				const respPage =
					res.data && res.data.page ? Number(res.data.page) : Number(page);
				if (
					!Number.isNaN(respPage) &&
					paginationRef &&
					paginationRef.value
				) {
					paginationRef.value.page = respPage;
				}
			} catch (e) {
				// ignore
			}
		}

		if (res && res.success === false) {
			if (res.message) {
				errorMessage.value = res.message;
			} else {
				errorMessage.value = "Lỗi tải dữ liệu";
			}
		}

		return res;
	}

	function handlePageChange(page) {
		page = Number(page);
		if (Number.isNaN(page) || page < 1) {
			return;
		}
		if (page > paginationRef.value.totalPages) {
			return;
		}
		load(page);
	}

	watch(searchQuery, function () {
		clearTimeout(_debounceTimer);
		_debounceTimer = setTimeout(function () {
			load(1);
		}, debounceMs);
	});

	return {
		searchQuery: searchQuery,
		load: load,
		handlePageChange: handlePageChange,
		errorMessage: errorMessage,
	};
}
