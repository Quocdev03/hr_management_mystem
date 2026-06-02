package repository

import (
	"chiquoc_hocgolang/internal/model"

	"gorm.io/gorm"
)

// EmployeeRepository interface cho các thao tác với bảng employees

type EmployeeRepository interface {
	WithTx(tx *gorm.DB) EmployeeRepository
	Create(emp *model.Employee) error
	FindAll(query model.PaginationQuery) ([]model.Employee, int64, error)
	FindByID(id uint) (*model.Employee, error)
	Update(emp *model.Employee) error
	UpdateFields(id uint, fields map[string]interface{}) error
	Delete(id uint) error
	CountByDepartment(deptID uint) (int64, error)
	FindByUserID(userID uint) (*model.Employee, error)
}

// --- EmployeeRepository Implementation ---
type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{
		db: db,
	}
}

func (r *employeeRepository) Create(emp *model.Employee) error {
	return r.db.Create(emp).Error
}

// Tìm nhân viên với phân trang, tìm kiếm, preload quan hệ
func (r *employeeRepository) FindAll(query model.PaginationQuery) ([]model.Employee, int64, error) {
	var employees []model.Employee
	var total int64

	db := r.db.Model(&model.Employee{})
	// Tìm theo tên
	if query.Search != "" {
		db = db.Where(
			"first_name LIKE ? OR last_name LIKE ? OR phone LIKE ?",
			"%"+query.Search+"%",
			"%"+query.Search+"%",
			"%"+query.Search+"%",
		)
	}

	// Lọc theo phòng ban khi cần
	if query.DepartmentID > 0 {
		db = db.Where("department_id = ?", query.DepartmentID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// offset = (page - 1) * limit
	offset := (query.Page - 1) * query.Limit

	// Preload Department đê trả ra thông tin phòng ban
	err := db.
		Preload("Department").
		Preload("User").
		Preload("User.Role").
		Offset(offset).
		Limit(query.Limit).
		Order("created_at DESC").
		Find(&employees).Error

	return employees, total, err
}

func (r *employeeRepository) FindByID(id uint) (*model.Employee, error) {
	var employee model.Employee
	err := r.db.
		Preload("Department").
		Preload("User").
		Preload("User.Role").
		First(&employee, id).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil

}

func (r *employeeRepository) FindByUserID(userID uint) (*model.Employee, error) {
	var employee model.Employee
	err := r.db.Preload("Department").Preload("User").Preload("User.Role").Where("user_id = ?", userID).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepository) Update(emp *model.Employee) error {
	return r.db.Save(emp).Error
}

func (r *employeeRepository) UpdateFields(id uint, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}

	return r.db.Model(&model.Employee{}).
		Where("id = ?", id).
		Updates(fields).Error
}

func (r *employeeRepository) WithTx(tx *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: tx}
}

func (r *employeeRepository) Delete(id uint) error {
	return r.db.Delete(&model.Employee{}, id).Error
}

func (r *employeeRepository) CountByDepartment(deptID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Employee{}).Where("department_id = ?", deptID).Count(&count).Error
	return count, err
}
