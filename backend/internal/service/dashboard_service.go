package service

import (
	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/internal/repository"
	"chiquoc_hocgolang/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type DashboardService interface {
	GetStats() (*model.Dashboards, error)
}

type dashboardService struct {
	dashRepo repository.DashboardsRepository
	rdb      *redis.Client
}

func NewDashboardService(dashRepo repository.DashboardsRepository, rdb *redis.Client) DashboardService {
	return &dashboardService{
		dashRepo: dashRepo,
		rdb:      rdb,
	}
}

func (s *dashboardService) GetStats() (*model.Dashboards, error) {
	ctx := context.Background()
	const cacheKey = "dashboard:stats"

	// Thử lấy dữ liệu từ Redis cache
	if s.rdb != nil {
		val, err := s.rdb.Get(ctx, cacheKey).Result()
		if err == nil {
			var stats model.Dashboards
			if jsonErr := json.Unmarshal([]byte(val), &stats); jsonErr == nil {
				utils.Info("lấy dữ liệu dashboard từ Redis cache thành công")
				return &stats, nil
			}
		}
		utils.Info("cache miss hoặc Redis lỗi (%v). Tiến hành lấy dữ liệu từ MySQL database...", err)
	}

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

	// Lưu kết quả vào Redis cache, TTL 1 giờ. Bỏ qua lỗi nếu Redis không khả dụng.
	if s.rdb != nil {
		bytes, marshalErr := json.Marshal(dbStats)
		if marshalErr == nil {
			if setErr := s.rdb.Set(ctx, cacheKey, bytes, time.Hour).Err(); setErr != nil {
				utils.Error("không thể lưu dashboard cache: %v", setErr)
			}
		} else {
			utils.Error("không thể marshal dashboard stats: %v", marshalErr)
		}
	}

	return dbStats, nil
}
