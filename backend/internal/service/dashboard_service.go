package service

import (
	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/internal/repository"
)

type DashboardService interface {
	GetStats() (*model.Dashboards, error)
}

type dashboardService struct {
	dashRepo repository.DashboardsRepository
}

func NewDashboardService(dashRepo repository.DashboardsRepository) DashboardService {
	return &dashboardService{
		dashRepo: dashRepo,
	}
}

func (s *dashboardService) GetStats() (*model.Dashboards, error) {
	// Gọi repository để lấy các con số
	stats := &model.Dashboards{
		TotalUsers:       s.dashRepo.CountUser(),
		TotalEmployees:   s.dashRepo.CountEmployees(),
		TotalDepartments: s.dashRepo.CountDepartments(),
	}

	return stats, nil
}
