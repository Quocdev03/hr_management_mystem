package router

import (
	"chiquoc_hocgolang/internal/config"
	"chiquoc_hocgolang/internal/handler"
	"chiquoc_hocgolang/internal/middleware"
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
		// GET /employees — Danh sách nhân viên (tất cả role, không cache — DB đủ nhanh)
		employees.GET("",
			empHandler.GetEmployees,
		)

		// GET /employees/:id — Chi tiết nhân viên (tất cả role, không cache)
		employees.GET("/:id",
			empHandler.GetEmployee,
		)

		// POST /employees — Tạo nhân viên mới (admin, hr)
		employees.POST("",
			middleware.RequireRole("admin", "hr"),
			middleware.ClearMultipleCaches(rdb, middleware.EmployeeRelatedCachePatterns...),
			empHandler.CreateEmployee,
		)

		// PUT /employees/:id — Cập nhật nhân viên đầy đủ (admin, hr)
		employees.PUT("/:id",
			middleware.RequireRole("admin", "hr"),
			middleware.ClearMultipleCaches(rdb, middleware.EmployeeRelatedCachePatterns...),
			empHandler.UpdateEmployee,
		)

		// PATCH /employees/:id — Cập nhật nhân viên (admin, hr)
		employees.PATCH("/:id",
			middleware.RequireRole("admin", "hr"),
			middleware.ClearMultipleCaches(rdb, middleware.EmployeeRelatedCachePatterns...),
			empHandler.UpdateEmployee,
		)

		// DELETE /employees/:id — Xóa nhân viên (chỉ admin)
		employees.DELETE("/:id",
			middleware.RequireRole("admin"),
			middleware.ClearMultipleCaches(rdb, middleware.EmployeeRelatedCachePatterns...),
			empHandler.DeleteEmployee,
		)
	}

	// ── Users (/users) ──────────────────────────────────────────────────
	// Quản lý tài khoản: chỉ admin được phép truy cập
	users := protected.Group("/users")
	{
		// GET /users — Danh sách tài khoản (admin, không cache)
		users.GET("",
			middleware.RequireRole("admin"),
			userHandler.GetUsers,
		)

		// GET /users/available — Tài khoản chưa liên kết nhân viên (admin, không cache)
		users.GET("/available",
			middleware.RequireRole("admin"),
			userHandler.GetUsersWithoutEmployee,
		)

		// GET /users/:id — Chi tiết tài khoản (admin, không cache)
		users.GET("/:id",
			middleware.RequireRole("admin"),
			userHandler.GetUser,
		)

		// POST /users — Tạo tài khoản mới (admin)
		users.POST("",
			middleware.RequireRole("admin"),
			middleware.ClearMultipleCaches(rdb, middleware.UserRelatedCachePatterns...),
			userHandler.CreateUser,
		)

		// PUT /users/:id — Cập nhật tài khoản đầy đủ (admin)
		users.PUT("/:id",
			middleware.RequireRole("admin"),
			middleware.ClearMultipleCaches(rdb, middleware.UserRelatedCachePatterns...),
			userHandler.UpdateUser,
		)

		// PATCH /users/:id — Cập nhật tài khoản (admin)
		users.PATCH("/:id",
			middleware.RequireRole("admin"),
			middleware.ClearMultipleCaches(rdb, middleware.UserRelatedCachePatterns...),
			userHandler.UpdateUser,
		)

		// DELETE /users/:id — Xóa tài khoản (admin)
		users.DELETE("/:id",
			middleware.RequireRole("admin"),
			middleware.ClearMultipleCaches(rdb, middleware.UserRelatedCachePatterns...),
			userHandler.DeleteUser,
		)
	}

	// ── Departments (/departments) ──────────────────────────────────────
	// Quản lý phòng ban: đọc tất cả role, ghi chỉ admin
	departments := protected.Group("/departments")
	{
		// GET /departments — Danh sách phòng ban (tất cả role, không cache)
		departments.GET("",
			deptHandler.GetDepartments,
		)

		// GET /departments/:id — Chi tiết phòng ban (tất cả role, không cache)
		departments.GET("/:id",
			deptHandler.GetDepartment,
		)

		// POST /departments — Tạo phòng ban mới (chỉ admin)
		departments.POST("",
			middleware.RequireRole("admin"),
			middleware.ClearMultipleCaches(rdb, middleware.DepartmentRelatedCachePatterns...),
			deptHandler.CreateDepartment,
		)

		// PUT /departments/:id — Cập nhật phòng ban đầy đủ (chỉ admin)
		departments.PUT("/:id",
			middleware.RequireRole("admin"),
			middleware.ClearMultipleCaches(rdb, middleware.DepartmentRelatedCachePatterns...),
			deptHandler.UpdateDepartment,
		)

		// PATCH /departments/:id — Cập nhật phòng ban (chỉ admin)
		departments.PATCH("/:id",
			middleware.RequireRole("admin"),
			middleware.ClearMultipleCaches(rdb, middleware.DepartmentRelatedCachePatterns...),
			deptHandler.UpdateDepartment,
		)

		// DELETE /departments/:id — Xóa phòng ban (chỉ admin)
		departments.DELETE("/:id",
			middleware.RequireRole("admin"),
			middleware.ClearMultipleCaches(rdb, middleware.DepartmentRelatedCachePatterns...),
			deptHandler.DeleteDepartment,
		)
	}

	return r
}
