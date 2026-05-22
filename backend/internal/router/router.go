package router

import (
	"chiquoc_hocgolang/internal/config"
	"chiquoc_hocgolang/internal/handler"
	"chiquoc_hocgolang/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter khởi tạo và cấu hình tất cả routes
// Nơi kết nối middleware -> handler -> service -> repository

func SetupRouter(cfg *config.Config, authHandler *handler.AuthHandler, empHandler *handler.EmployeeHandler, deptHandler *handler.DepartmentHandler, dashB *handler.DashboardsHanlder, userHandler *handler.UserHandler) *gin.Engine {
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
		// Auth không cần token
		auth.POST("/login", authHandler.Login)

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
		employees.GET("", empHandler.GetEmployees)
		employees.GET("/:id", empHandler.GetEmployee)

		// / Tạo/Sửa/Xóa: chỉ admin và hr mới được
		employees.POST("", middleware.RequireRole("admin", "hr"), empHandler.CreateEmployee)
		employees.PUT("/:id", middleware.RequireRole("admin", "hr"), empHandler.UpdateEmployee)
		
		

		// chỉ admin xóa
		employees.DELETE("/:id", middleware.RequireRole("admin"), empHandler.DeleteEmployee)

	}
	// Users
	users := protected.Group("/users")
	{
		users.GET("", middleware.RequireRole("admin"), userHandler.GetUsers)
		users.GET("/:id", middleware.RequireRole("admin"), userHandler.GetUser)
		users.POST("", middleware.RequireRole("admin"), userHandler.CreateUser)
		users.PUT("/:id", middleware.RequireRole("admin"), userHandler.UpdateUser)
		users.DELETE("/:id", middleware.RequireRole("admin"), userHandler.DeleteUser)
	}

	// Departments
	departments := protected.Group("departments")
	{
		departments.GET("", deptHandler.GetDepartments)
		departments.GET("/:id", deptHandler.GetDepartment)
		departments.POST("", middleware.RequireRole("admin"), deptHandler.CreateDepartment)
		departments.PUT("/:id", middleware.RequireRole("admin"), deptHandler.UpdateDepartment)
		departments.DELETE("/:id", middleware.RequireRole("admin"), deptHandler.DeleteDepartment)

	}

	return r
}
