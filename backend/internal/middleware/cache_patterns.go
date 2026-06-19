package middleware

// Cache key pattern constants dùng để invalidate cache theo nhóm endpoint.
const (
	CachePatternEmployees   = "cache:/api/v1/employees*"
	CachePatternUsers       = "cache:/api/v1/users*"
	CachePatternDepartments = "cache:/api/v1/departments*"
	CachePatternPositions   = "cache:/api/v1/positions*"
)

// Pattern groups để dùng với ClearMultipleCaches.
var (
	// EmployeeRelatedCachePatterns xóa toàn bộ cache liên quan đến nhân viên,
	// user, và phòng ban (vì chúng phụ thuộc lẫn nhau).
	EmployeeRelatedCachePatterns = []string{
		CachePatternEmployees,
		CachePatternUsers,
		CachePatternDepartments,
	}

	// UserRelatedCachePatterns xóa cache liên quan đến user accounts.
	UserRelatedCachePatterns = []string{
		CachePatternUsers,
	}

	// DepartmentRelatedCachePatterns xóa cache liên quan đến phòng ban.
	DepartmentRelatedCachePatterns = []string{
		CachePatternDepartments,
	}
)
