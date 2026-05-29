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

// SetupRouter khởi tạo và cấu hình tất cả routes
// Nơi kết nối middleware -> handler -> service -> repository

func SetupRouter(cfg *config.Config, rdb *redis.Client, authHandler *handler.AuthHandler, empHandler *handler.EmployeeHandler, deptHandler *handler.DepartmentHandler, dashB *handler.DashboardsHanlder, userHandler *handler.UserHandler) *gin.Engine {
	// Tắt debug log trong production
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Dùng gin.New() thay vì gin.Default() để tự cấu hình middleware
	r := gin.New()

	// Áp dụng middleware chung toàn bộ routes
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// Public endpoint kiểm tra server
	r.GET("/api/v1/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "HRM API đang chạy!",
		})
	})

	// API V1 GROUP
	// Prefix: /api/v1
	v1 := r.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		// Auth không cần token nhưng áp dụng rate limit: tối đa 5 lượt/phút
		auth.POST("/login", middleware.RateLimiter(rdb, 5, time.Minute), authHandler.Login)

		// Profile Cần JWT token
		auth.GET("/profile", middleware.AuthJWT(&cfg.JWT), authHandler.GetProfile)
	}

	// Tất cả routes bên dưới đều cần JWT token
	protected := v1.Group("")
	protected.Use(middleware.AuthJWT(&cfg.JWT))

	dashboard := protected.Group("/dashboard")
	{
		dashboard.GET("/stats", dashB.GetStats)
	}

	// Employees
	employees := protected.Group("/employees")
	{
		// Xem danh sách: tất cả roles đăng nhập đều xem được
		employees.GET("", middleware.CacheResponse(rdb, 15*time.Minute), empHandler.GetEmployees)
		employees.GET("/:id", middleware.CacheResponse(rdb, 15*time.Minute), empHandler.GetEmployee)

		// / Tạo/Sửa/Xóa: chỉ admin và hr mới được
		employees.POST("", middleware.RequireRole("admin", "hr"), middleware.ClearCache(rdb, "cache:/api/v1/employees*"), empHandler.CreateEmployee)
		employees.PUT("/:id", middleware.RequireRole("admin", "hr"), middleware.ClearCache(rdb, "cache:/api/v1/employees*"), empHandler.UpdateEmployee)

		// chỉ admin xóa
		employees.DELETE("/:id", middleware.RequireRole("admin"), middleware.ClearCache(rdb, "cache:/api/v1/employees*"), empHandler.DeleteEmployee)

	}
	// Users
	users := protected.Group("/users")
	{
		// / Tạo/Sửa/Xóa: chỉ admin
		users.GET("", middleware.RequireRole("admin"), middleware.CacheResponse(rdb, 15*time.Minute), userHandler.GetUsers)
		users.GET("/:id", middleware.RequireRole("admin"), middleware.CacheResponse(rdb, 15*time.Minute), userHandler.GetUser)
		users.GET("/available", middleware.RequireRole("admin"), middleware.CacheResponse(rdb, 15*time.Minute), userHandler.GetUsersWithoutEmployee)
		
		users.POST("", middleware.RequireRole("admin"), middleware.ClearCache(rdb, "cache:/api/v1/users*"), userHandler.CreateUser)
		users.PUT("/:id", middleware.RequireRole("admin"), middleware.ClearCache(rdb, "cache:/api/v1/users*"), userHandler.UpdateUser)
		users.DELETE("/:id", middleware.RequireRole("admin"), middleware.ClearCache(rdb, "cache:/api/v1/users*"), userHandler.DeleteUser)
	}

	// Departments
	departments := protected.Group("departments")
	{
		// Tạo/Sửa/Xóa: chỉ admin
		departments.GET("", middleware.CacheResponse(rdb, 15*time.Minute), deptHandler.GetDepartments)
		departments.GET("/:id", middleware.CacheResponse(rdb, 15*time.Minute), deptHandler.GetDepartment)
		
		departments.POST("", middleware.RequireRole("admin"), middleware.ClearCache(rdb, "cache:/api/v1/departments*"), deptHandler.CreateDepartment)
		departments.PUT("/:id", middleware.RequireRole("admin"), middleware.ClearCache(rdb, "cache:/api/v1/departments*"), deptHandler.UpdateDepartment)
		departments.DELETE("/:id", middleware.RequireRole("admin"), middleware.ClearCache(rdb, "cache:/api/v1/departments*"), deptHandler.DeleteDepartment)

	}

	return r
}
