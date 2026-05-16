package repository

import (
	"chiquoc_hocgolang/internal/model"

	"gorm.io/gorm"
)

type DashboardsRepository interface {
	CountUser() int64
	CountEmployees() int64
	CountDepartments() int64
	GetEmployeeActive() int64
	GetEmployeeCountByDepartment() []model.DepartmentEmployeeCount
}

type dashboardsRepository struct {
	db *gorm.DB
}

func NewDashboardsRepository(db *gorm.DB) DashboardsRepository {
	return &dashboardsRepository{
		db: db,
	}
}

func (d *dashboardsRepository) CountUser() int64 {
	var count int64
	d.db.Model(&model.User{}).Count(&count)
	return count
}

func (d *dashboardsRepository) CountEmployees() int64 {
	var count int64
	d.db.Model(&model.Employee{}).Count(&count)
	return count
}

func (d *dashboardsRepository) CountDepartments() int64 {
	var count int64
	d.db.Model(&model.Department{}).Count(&count)
	return count
}

func (d *dashboardsRepository) GetEmployeeCountByDepartment() []model.DepartmentEmployeeCount {
	var stats []model.DepartmentEmployeeCount

	// Lấy tên phòng ban và đếm số nhân viên thuộc phòng ban đó
	query := `
	SELECT departments.name AS department_name, COUNT(employees.id) AS employee_count
	FROM departments
	LEFT JOIN employees ON employees.department_id = departments.id
	GROUP BY department_name, department_id `
	d.db.Raw(query).Scan(&stats)
	return stats
}

func (d *dashboardsRepository) GetEmployeeActive() int64 {
	var count int64
	query := `SELECT COUNT(*) AS active_employee
			FROM employees
			WHERE status = "active"`
	d.db.Raw(query).Scan(&count)
	return count
}
