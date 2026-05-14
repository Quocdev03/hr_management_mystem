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

func SetupRouter(cfg *config.Config, authHandler *handler.AuthHandler) *gin.Engine {
	// Tắt debug log trong production
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Dùng gin.New() thay vì gin.Default() để tự cấu hình middleware
	r := gin.New()

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

	// Auth - (không cần token)
	auth := v1.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)

		// Protected: cần JWT token
		auth.GET("/profile", middleware.AuthJWT(&cfg.JWT), authHandler.GetProfile)
	}

	return r
}
