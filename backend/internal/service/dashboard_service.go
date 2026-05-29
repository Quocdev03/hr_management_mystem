package service

import (
	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/internal/repository"
	"chiquoc_hocgolang/internal/utils"
	"context"
	"fmt"
	"time"
)

type DashboardService interface {
	GetStats() (*model.Dashboards, error)
}

type dashboardService struct {
	dashRepo repository.DashboardsRepository
	cacheSvc CacheService
}

func NewDashboardService(dashRepo repository.DashboardsRepository, cacheSvc CacheService) DashboardService {
	return &dashboardService{
		dashRepo: dashRepo,
		cacheSvc: cacheSvc,
	}
}

func (s *dashboardService) GetStats() (*model.Dashboards, error) {
	ctx := context.Background()
	var stats model.Dashboards

	// Thử lấy dữ liệu từ Redis cache
	err := s.cacheSvc.Get(ctx, "dashboard:stats", &stats)
	if err == nil {
		utils.Info("Lấy dữ liệu dashboard từ Redis cache thành công!")
		return &stats, nil
	}

	utils.Info("Cache miss hoặc Redis lỗi (%v). Tiến hành lấy dữ liệu từ MySQL database...", err)

	// Lấy dữ liệu trực tiếp từ database
	totalUsers, err := s.dashRepo.CountUser()
	if err != nil {
		return nil, fmt.Errorf("lỗi lấy tổng số user: %w", err)
	}

	totalEmployees, err := s.dashRepo.CountEmployees()
	if err != nil {
		return nil, fmt.Errorf("lỗi lấy tổng số nhân viên: %w", err)
	}

	totalDepartments, err := s.dashRepo.CountDepartments()
	if err != nil {
		return nil, fmt.Errorf("lỗi lấy tổng số phòng ban: %w", err)
	}

	totalActive, err := s.dashRepo.GetEmployeeActive()
	if err != nil {
		return nil, fmt.Errorf("lỗi lấy tổng số nhân viên active: %w", err)
	}

	deptStats, err := s.dashRepo.GetEmployeeCountByDepartment()
	if err != nil {
		return nil, fmt.Errorf("lỗi lấy thống kê phòng ban: %w", err)
	}

	dbStats := &model.Dashboards{
		TotalUsers:           totalUsers,
		TotalEmployees:       totalEmployees,
		TotalDepartments:     totalDepartments,
		TotalEmployeesActive: totalActive,
		DepartmentStats:      deptStats,
	}

	// Lưu kết quả vào Redis cache, đặt TTL là 1 giờ. Bỏ qua lỗi nếu lưu cache thất bại.
	_ = s.cacheSvc.Set(ctx, "dashboard:stats", dbStats, 1*time.Hour)

	return dbStats, nil
}
