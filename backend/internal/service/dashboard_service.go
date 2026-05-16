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

	stats := &model.Dashboards{
		TotalUsers:           s.dashRepo.CountUser(),
		TotalEmployees:       s.dashRepo.CountEmployees(),
		TotalDepartments:     s.dashRepo.CountDepartments(),
		TotalEmployeesActive: s.dashRepo.GetEmployeeActive(),
		DepartmentStats:      s.dashRepo.GetEmployeeCountByDepartment(),
	}

	return stats, nil
}
