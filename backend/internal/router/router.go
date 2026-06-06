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
	// Thống kê tổng quan hệ thống
	dashboard := protected.Group("/dashboard")
	{
		// GET /dashboard/stats — Lấy thống kê (tất cả role)
		dashboard.GET("/stats", dashB.GetStats)
	}

	// ── Employees (/employees) ──────────────────────────────────────────
	// Quản lý nhân viên: CRUD + phân quyền
	employees := protected.Group("/employees")
	{
		// GET /employees     — Danh sách nhân viên (tất cả role, có cache 15 phút)
		employees.GET("",
			middleware.CacheResponse(rdb, 15*time.Minute),
			empHandler.GetEmployees,
		)

		// GET /employees/:id — Chi tiết nhân viên (tất cả role, có cache 15 phút)
		employees.GET("/:id",
			middleware.CacheResponse(rdb, 15*time.Minute),
			empHandler.GetEmployee,
		)

		// POST /employees    — Tạo nhân viên mới (admin, hr)
		employees.POST("",
			middleware.RequireRole("admin", "hr"),
			middleware.ClearCache(rdb, "cache:/api/v1/employees*"),
			middleware.ClearCache(rdb, "cache:/api/v1/users*"),
			middleware.ClearCache(rdb, "cache:/api/v1/departments*"),
			empHandler.CreateEmployee,
		)

		// PUT /employees/:id — Cập nhật nhân viên đầy đủ (admin, hr)
		employees.PUT("/:id",
			middleware.RequireRole("admin", "hr"),
			middleware.ClearCache(rdb, "cache:/api/v1/employees*"),
			middleware.ClearCache(rdb, "cache:/api/v1/users*"),
			middleware.ClearCache(rdb, "cache:/api/v1/departments*"),
			empHandler.UpdateEmployee,
		)

		// PATCH /employees/:id — Cập nhật nhân viên (admin, hr)
		employees.PATCH("/:id",
			middleware.RequireRole("admin", "hr"),
			middleware.ClearCache(rdb, "cache:/api/v1/employees*"),
			middleware.ClearCache(rdb, "cache:/api/v1/users*"),
			middleware.ClearCache(rdb, "cache:/api/v1/departments*"),
			empHandler.UpdateEmployee,
		)

		// DELETE /employees/:id — Xóa nhân viên (chỉ admin)
		employees.DELETE("/:id",
			middleware.RequireRole("admin"),
			middleware.ClearCache(rdb, "cache:/api/v1/employees*"),
			middleware.ClearCache(rdb, "cache:/api/v1/users*"),
			middleware.ClearCache(rdb, "cache:/api/v1/departments*"),
			empHandler.DeleteEmployee,
		)
	}

	// ── Users (/users) ──────────────────────────────────────────────────
	// Quản lý tài khoản: chỉ admin được phép truy cập
	users := protected.Group("/users")
	{
		// GET /users           — Danh sách tài khoản (admin, có cache 15 phút)
		users.GET("",
			middleware.RequireRole("admin"),
			middleware.CacheResponse(rdb, 15*time.Minute),
			userHandler.GetUsers,
		)

		// GET /users/available — Tài khoản chưa liên kết nhân viên (admin, có cache 15 phút)
		users.GET("/available",
			middleware.RequireRole("admin"),
			middleware.CacheResponse(rdb, 15*time.Minute),
			userHandler.GetUsersWithoutEmployee,
		)

		// GET /users/:id      — Chi tiết tài khoản (admin, có cache 15 phút)
		users.GET("/:id",
			middleware.RequireRole("admin"),
			middleware.CacheResponse(rdb, 15*time.Minute),
			userHandler.GetUser,
		)

		// POST /users         — Tạo tài khoản mới (admin)
		users.POST("",
			middleware.RequireRole("admin"),
			middleware.ClearCache(rdb, "cache:/api/v1/users*"),
			userHandler.CreateUser,
		)

		// PUT /users/:id      — Cập nhật tài khoản đầy đủ (admin)
		users.PUT("/:id",
			middleware.RequireRole("admin"),
			middleware.ClearCache(rdb, "cache:/api/v1/users*"),
			userHandler.UpdateUser,
		)

		// PATCH /users/:id      — Cập nhật tài khoản (admin)
		users.PATCH("/:id",
			middleware.RequireRole("admin"),
			middleware.ClearCache(rdb, "cache:/api/v1/users*"),
			userHandler.UpdateUser,
		)

		// DELETE /users/:id   — Xóa tài khoản (admin)
		users.DELETE("/:id",
			middleware.RequireRole("admin"),
			middleware.ClearCache(rdb, "cache:/api/v1/users*"),
			userHandler.DeleteUser,
		)
	}

	// ── Departments (/departments) ──────────────────────────────────────
	// Quản lý phòng ban: đọc tất cả role, ghi chỉ admin
	departments := protected.Group("/departments")
	{
		// GET /departments     — Danh sách phòng ban (tất cả role, có cache 15 phút)
		departments.GET("",
			middleware.CacheResponse(rdb, 15*time.Minute),
			deptHandler.GetDepartments,
		)

		// GET /departments/:id — Chi tiết phòng ban (tất cả role, có cache 15 phút)
		departments.GET("/:id",
			middleware.CacheResponse(rdb, 15*time.Minute),
			deptHandler.GetDepartment,
		)

		// POST /departments    — Tạo phòng ban mới (chỉ admin)
		departments.POST("",
			middleware.RequireRole("admin"),
			middleware.ClearCache(rdb, "cache:/api/v1/departments*"),
			deptHandler.CreateDepartment,
		)

		// PUT /departments/:id — Cập nhật phòng ban đầy đủ (chỉ admin)
		departments.PUT("/:id",
			middleware.RequireRole("admin"),
			middleware.ClearCache(rdb, "cache:/api/v1/departments*"),
			deptHandler.UpdateDepartment,
		)

		// PATCH /departments/:id — Cập nhật phòng ban (chỉ admin)
		departments.PATCH("/:id",
			middleware.RequireRole("admin"),
			middleware.ClearCache(rdb, "cache:/api/v1/departments*"),
			deptHandler.UpdateDepartment,
		)

		// DELETE /departments/:id — Xóa phòng ban (chỉ admin)
		departments.DELETE("/:id",
			middleware.RequireRole("admin"),
			middleware.ClearCache(rdb, "cache:/api/v1/departments*"),
			deptHandler.DeleteDepartment,
		)
	}

	return r
}
