package router

import (
	"chiquoc_hocgolang/internal/config"
	"chiquoc_hocgolang/internal/handler"
	"chiquoc_hocgolang/internal/middleware"
	"chiquoc_hocgolang/internal/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// SetupRouter khởi tạo và cấu hình tất cả routes.
// Kết nối middleware -> handler -> service -> repository.
func SetupRouter(
	cfg *config.Config,
	rdb *redis.Client,
	authHandler *handler.AuthHandler,
	empHandler *handler.EmployeeHandler,
	deptHandler *handler.DepartmentHandler,
	dashB *handler.DashboardHandler,
	userHandler *handler.UserHandler,
	permRepo repository.PermissionRepository,
) *gin.Engine {

	// Tắt debug log trong production
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Dùng gin.New() thay vì gin.Default() để tự cấu hình middleware
	r := gin.New()

	// ── Global Middleware ────────────────────────────────────────────────
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// ── Health Check (public) ───────────────────────────────────────────
	r.GET("/api/v1/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "HRM API đang chạy!",
		})
	})

	// ── API v1 Group ────────────────────────────────────────────────────
	v1 := r.Group("/api/v1")

	// ── Auth Routes (/auth) ─────────────────────────────────────────────
	// Các endpoint xác thực: login, refresh, profile, logout
	auth := v1.Group("/auth")
	{
		// POST /auth/login — Đăng nhập, rate limit 5 lượt/phút
		auth.POST("/login", middleware.RateLimiter(rdb, 5, time.Minute), authHandler.Login)

		// POST /auth/refresh — Làm mới access token, rate limit 10 lượt/phút
		auth.POST("/refresh", middleware.RateLimiter(rdb, 10, time.Minute), authHandler.Refresh)

		// GET /auth/me — Lấy hồ sơ người dùng (cần JWT)
		auth.GET("/me", middleware.AuthJWT(&cfg.JWT, rdb), authHandler.GetProfile)

		// POST /auth/logout — Đăng xuất, blacklist token (cần JWT)
		auth.POST("/logout", middleware.AuthJWT(&cfg.JWT, rdb), authHandler.Logout)
	}

	// ── Protected Routes (cần JWT) ──────────────────────────────────────
	protected := v1.Group("")
	protected.Use(middleware.AuthJWT(&cfg.JWT, rdb))

	// ── Dashboard (/dashboard) ──────────────────────────────────────────
	// Thống kê tổng quan hệ thống — cache 1 giờ (hit rate cao, key cố định)
	dashboard := protected.Group("/dashboard")
	{
		// GET /dashboard/stats — Lấy thống kê (tất cả role), cache 1 giờ
		dashboard.GET("/stats",
			middleware.CacheResponse(rdb, 1*time.Hour),
			dashB.GetStats,
		)
	}

	// ── Employees (/employees) ──────────────────────────────────────────
	// Quản lý nhân viên: CRUD + phân quyền
	employees := protected.Group("/employees")
	{
		// GET /employees — Danh sách nhân viên (cần quyền employee.read)
		employees.GET("",
			middleware.RequirePermission(permRepo, "employee.read"),
			empHandler.GetEmployees,
		)

		// GET /employees/:id — Chi tiết nhân viên (cần quyền employee.read)
		employees.GET("/:id",
			middleware.RequirePermission(permRepo, "employee.read"),
			empHandler.GetEmployee,
		)

		// POST /employees — Tạo nhân viên mới (permission-based)
		employees.POST("",
			middleware.RequirePermission(permRepo, "employee.create"),
			middleware.ClearMultipleCaches(rdb, middleware.EmployeeRelatedCachePatterns...),
			empHandler.CreateEmployee,
		)

		// PATCH /employees/:id — Cập nhật nhân viên (permission-based)
		employees.PATCH("/:id",
			middleware.RequirePermission(permRepo, "employee.update"),
			middleware.ClearMultipleCaches(rdb, middleware.EmployeeRelatedCachePatterns...),
			empHandler.UpdateEmployee,
		)

		// DELETE /employees/:id — Xóa nhân viên (permission-based)
		employees.DELETE("/:id",
			middleware.RequirePermission(permRepo, "employee.delete"),
			middleware.ClearMultipleCaches(rdb, middleware.EmployeeRelatedCachePatterns...),
			empHandler.DeleteEmployee,
		)
	}

	// ── Users (/users) ──────────────────────────────────────────────────
	// Quản lý tài khoản: chỉ admin được phép truy cập
	users := protected.Group("/users")
	{
		// GET /users — Danh sách tài khoản (permission-based)
		users.GET("",
			middleware.RequirePermission(permRepo, "user.read"),
			userHandler.GetUsers,
		)

		// GET /users/available — Tài khoản chưa liên kết nhân viên (permission-based)
		users.GET("/available",
			middleware.RequirePermission(permRepo, "user.read"),
			userHandler.GetUsersWithoutEmployee,
		)

		// GET /users/:id — Chi tiết tài khoản (permission-based)
		users.GET("/:id",
			middleware.RequirePermission(permRepo, "user.read"),
			userHandler.GetUser,
		)

		// POST /users — Tạo tài khoản mới (permission-based)
		users.POST("",
			middleware.RequirePermission(permRepo, "user.create"),
			middleware.ClearMultipleCaches(rdb, middleware.UserRelatedCachePatterns...),
			userHandler.CreateUser,
		)

		// PATCH /users/:id — Cập nhật tài khoản (permission-based)
		users.PATCH("/:id",
			middleware.RequirePermission(permRepo, "user.update"),
			middleware.ClearMultipleCaches(rdb, middleware.UserRelatedCachePatterns...),
			userHandler.UpdateUser,
		)

		// GET /users/permissions — Danh sách permission có thể gán
		users.GET("/permissions",
			middleware.RequirePermission(permRepo, "user.update"),
			userHandler.GetAvailablePermissions,
		)

		// PATCH /users/:id/permissions — Cập nhật quyền cho user
		users.PATCH("/:id/permissions",
			middleware.RequirePermission(permRepo, "user.update"),
			userHandler.UpdatePermissions,
		)

		// DELETE /users/:id — Xóa tài khoản (permission-based)
		users.DELETE("/:id",
			middleware.RequirePermission(permRepo, "user.delete"),
			middleware.ClearMultipleCaches(rdb, middleware.UserRelatedCachePatterns...),
			userHandler.DeleteUser,
		)
	}

	// ── Departments (/departments) ──────────────────────────────────────
	// Quản lý phòng ban: đọc tất cả role, ghi chỉ admin
	departments := protected.Group("/departments")
	{
		// GET /departments — Danh sách phòng ban (cần quyền department.read)
		departments.GET("",
			middleware.RequirePermission(permRepo, "department.read"),
			deptHandler.GetDepartments,
		)

		// GET /departments/:id — Chi tiết phòng ban (cần quyền department.read)
		departments.GET("/:id",
			middleware.RequirePermission(permRepo, "department.read"),
			deptHandler.GetDepartment,
		)

		// POST /departments — Tạo phòng ban mới (permission-based)
		departments.POST("",
			middleware.RequirePermission(permRepo, "department.create"),
			middleware.ClearMultipleCaches(rdb, middleware.DepartmentRelatedCachePatterns...),
			deptHandler.CreateDepartment,
		)

		// PATCH /departments/:id — Cập nhật phòng ban (permission-based)
		departments.PATCH("/:id",
			middleware.RequirePermission(permRepo, "department.update"),
			middleware.ClearMultipleCaches(rdb, middleware.DepartmentRelatedCachePatterns...),
			deptHandler.UpdateDepartment,
		)

		// DELETE /departments/:id — Xóa phòng ban (permission-based)
		departments.DELETE("/:id",
			middleware.RequirePermission(permRepo, "department.delete"),
			middleware.ClearMultipleCaches(rdb, middleware.DepartmentRelatedCachePatterns...),
			deptHandler.DeleteDepartment,
		)
	}

	return r
}
