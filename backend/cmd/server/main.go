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
	db := config.InitiDB(&cfg.Database)

	// Seed data mẫu - chỉ chạy khi APP_SEED=true trong .env
	// Đặt APP_SEED=true để chèn data mẫu, APP_SEED=false (hoặc bỏ trống) để bỏ qua
	// config.SeedData(db)

	// Kết nối các layer của
	// Handler -> Service -> Repository

	// Repositories - tương tác trực tiếp với DB
	userRepo := repository.NewUserRepository(db)
	empRepo := repository.NewEmployeeRepository(db)
	deptRepo := repository.NewDepartmentRepository(db)
	dashRepo := repository.NewDashboardsRepository(db)

	// Services - chứa business logic
	authScv := service.NewAuthService(userRepo, &cfg.JWT)
	userScv := service.NewUserService(userRepo)
	empScv := service.NewEmployeeService(empRepo, deptRepo)
	deptScv := service.NewDepartmentService(deptRepo, empRepo)
	dashScv := service.NewDashboardService(dashRepo)

	// Handlers - nhận HTTP request, gọi service, trả response
	authHandler := handler.NewAuthHandler(authScv)
	userHandler := handler.NewUserHandler(userScv)
	empHandler := handler.NewEmployeeHandler(empScv)
	deptHandler := handler.NewDepartmentHandler(deptScv)
	dashHandler := handler.NewDashboardHanlder(dashScv)

	// Thiết lập router
	r := router.SetupRouter(cfg, authHandler, empHandler, deptHandler, dashHandler, userHandler)

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
	utils.Info("Đã nhận được tién hiệu: %v. Đang tắt...", sig)

	// Tạo context timeout 30s chờ request hoàn thành
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Tắt server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		utils.Error("Server forced to shutdown: %v", err)
	}

	// Đóng kết nối database
	sqlDB, _ := db.DB()
	sqlDB.Close()

	utils.Info("Server đã thoát thành công.")
}
