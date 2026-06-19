package main

import (
	"chiquoc_hocgolang/internal/config"
	"chiquoc_hocgolang/internal/handler"
	"chiquoc_hocgolang/internal/repository"
	"chiquoc_hocgolang/internal/router"
	"chiquoc_hocgolang/internal/service"
	"chiquoc_hocgolang/internal/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load các config
	cfg := config.Load()
	utils.Info("Đang chạy HMR API - Môi trường: %s", cfg.App.Env)

	// Khởi tạo Database
	db := config.InitDB(&cfg.Database)

	// Khởi tạo Redis
	rdb := config.InitRedis(&cfg.Redis)

	// In bảng tóm tắt cấu hình Redis khi khởi động
	fmt.Printf("[REDIS-info] Địa chỉ         : %s:%s (DB: %d)\n", cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.DB)
	fmt.Printf("[REDIS-info] Dashboard cache  : TTL 1 giờ  | Key: dashboard:stats\n")
	fmt.Printf("[REDIS-info] Login rate limit  : 5 lượt / 1 phút / IP\n")
	fmt.Printf("[REDIS-info] Token blacklist   : TTL theo token | Key: blacklist:{token}\n")

	// Kết nối các layer của
	// Handler -> Service -> Repository

	// Repositories - tương tác trực tiếp với DB
	userRepo := repository.NewUserRepository(db)
	empRepo := repository.NewEmployeeRepository(db)
	deptRepo := repository.NewDepartmentRepository(db)
	dashRepo := repository.NewDashboardsRepository(db)
	permRepo := repository.NewPermissionRepository(db)
	posRepo := repository.NewPositionRepository(db)

	// Services - chứa business logic
	authScv := service.NewAuthService(userRepo, empRepo, permRepo, &cfg.JWT, rdb)
	userScv := service.NewUserService(userRepo, permRepo, rdb)
	empScv := service.NewEmployeeService(db, empRepo, deptRepo, userRepo, rdb)
	deptScv := service.NewDepartmentService(db, deptRepo, empRepo, rdb)
	dashScv := service.NewDashboardService(dashRepo, rdb)
	posScv := service.NewPositionService(db, posRepo)

	// Handlers - nhận HTTP request, gọi service, trả response
	authHandler := handler.NewAuthHandler(authScv)
	userHandler := handler.NewUserHandler(userScv)
	empHandler := handler.NewEmployeeHandler(empScv)
	deptHandler := handler.NewDepartmentHandler(deptScv)
	dashHandler := handler.NewDashboardHandler(dashScv)
	posHandler := handler.NewPositionHandler(posScv)

	// Thiết lập router
	r := router.SetupRouter(cfg, rdb, authHandler, empHandler, deptHandler, dashHandler, userHandler, posHandler, permRepo)

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
		utils.Info("Server đang chạy tại port %s", cfg.App.Port)
		utils.Info("API Base URL: http://localhost:%s/api/v1", cfg.App.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Fatal("Lỗi khi khởi động server: %v", err)
		}
	}()

	// Tạo channel để nhận tín hiệu OS (ctrl + c)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block đến khi nhận tín hiệu tắt
	sig := <-quit
	utils.Info("Đã nhận được tín hiệu: %v. Đang tắt...", sig)

	// Tạo context timeout 30s chờ request hoàn thành
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Tắt server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		utils.Error("Server forced to shutdown: %v", err)
	}

	// Đóng kết nối database
	if sqlDB, err := db.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}
	// Đóng kết nối Redis
	if rdb != nil {
		if err := rdb.Close(); err != nil {
			log.Printf("failed to close redis: %v", err)
		}
	}

	utils.Info("Server đã thoát thành công.")
}
