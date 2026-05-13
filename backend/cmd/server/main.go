package main

import (
	"chiquoc_hocgolang/internal/config"
	"chiquoc_hocgolang/internal/router"
	"chiquoc_hocgolang/package/loggers"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load các config
	cfg := config.Load()
	loggers.Info("Đang chạy HMR API - Môi trường: %s", cfg.App.Env)

	// Khởi tạo Database
	db := config.InitiDB(&cfg.Database)

	// Seed data mẫu khi khởi động
	config.SeedData(db)

	// Kết nối các layer của
	// Handler -> Service -> Repository

	// Handler
	// Service
	// Repository

	// Thiết lập router
	r := router.SetupRouter(cfg)

	// Chạy server với Graceful Shutdown
	// Khi nhận SIGINT/SIGTERM, chờ các request đang xử lý hoàn thành trước khi tắt server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.App.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Chạy server trong goroutine riêng để không bị block
	go func() {
		loggers.Info("Server đang chạy tại port %s", cfg.App.Port)
		loggers.Info("API Base URL: http://localhost:%s/api/v1", cfg.App.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			loggers.Fatal("Lỗi khi khởi động server: %v", err)
		}
	}()

	// Tạo channel để nhận tín hiệu OS (ctrl + c)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block đến khi nhận tín hiệu tắt
	sig := <-quit
	loggers.Info("Đã nhận được tién hiệu: %v. Đang tắt...", sig)

	// Tạo context timeout 30s chờ request hoàn thành
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Tắt server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		loggers.Error("Server forced to shutdown: %v", err)
	}

	// Đóng kết nối database
	sqlDB, _ := db.DB()
	sqlDB.Close()

	loggers.Info("Server đã thoát thành công.")
}
