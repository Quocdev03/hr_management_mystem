package repository

import (
	"chiquoc_hocgolang/internal/model"

	"gorm.io/gorm"
)

type DashboardsRepository interface {
	CountUser() (int64, error)
	CountEmployees() (int64, error)
	CountDepartments() (int64, error)
	GetEmployeeActive() (int64, error)
	GetEmployeeCountByDepartment() ([]model.DepartmentEmployeeCount, error)
}

type dashboardsRepository struct {
	db *gorm.DB
}

func NewDashboardsRepository(db *gorm.DB) DashboardsRepository {
	return &dashboardsRepository{
		db: db,
	}
}

func (d *dashboardsRepository) CountUser() (int64, error) {
	var count int64
	err := d.db.Model(&model.User{}).Count(&count).Error
	return count, err
}

func (d *dashboardsRepository) CountEmployees() (int64, error) {
	var count int64
	err := d.db.Model(&model.Employee{}).Count(&count).Error
	return count, err
}

func (d *dashboardsRepository) CountDepartments() (int64, error) {
	var count int64
	err := d.db.Model(&model.Department{}).Count(&count).Error
	return count, err
}

func (d *dashboardsRepository) GetEmployeeCountByDepartment() ([]model.DepartmentEmployeeCount, error) {
	var stats []model.DepartmentEmployeeCount

	// Lấy tên phòng ban và đếm số nhân viên thuộc phòng ban đó
	query := `
	SELECT departments.name AS department_name, COUNT(employees.id) AS employee_count
	FROM departments
	LEFT JOIN employees ON employees.department_id = departments.id AND employees.deleted_at IS NULL
	WHERE departments.deleted_at IS NULL
	GROUP BY department_name, departments.id `
	err := d.db.Raw(query).Scan(&stats).Error
	return stats, err
}

func (d *dashboardsRepository) GetEmployeeActive() (int64, error) {
	var count int64
	err := d.db.Model(&model.Employee{}).Where("status = ?", "active").Count(&count).Error
	return count, err
}