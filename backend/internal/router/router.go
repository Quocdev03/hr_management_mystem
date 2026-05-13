package router

import (
	"chiquoc_hocgolang/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter khởi tạo và cấu hình tất cả routes
// Nơi kết nối middleware -> handler -> service -> repository

func SetupRouter(cfg *config.Config) *gin.Engine {
	// Tắt debug log trong production
	if cfg.App.Env == "producttion" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Dùng gin.New() thay vì gin.Default() để tự cấu hình middleware
	r := gin.New()

	// Public endpoint kiểm tra server
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "HRM API đang chạy!",
		})
	})

	return r
}
