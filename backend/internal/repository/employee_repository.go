package repository

import (
	"chiquoc_hocgolang/internal/model"

	"gorm.io/gorm"
)

// EmployeeRepository interface cho các thao tác với bảng employees

type EmployeeRepository interface {
	Create(emp *model.Employee) error
	FindAll(query model.PaginationQuery) ([]model.Employee, int64, error)
	FindByID(id uint) (*model.Employee, error)
	FindByEmail(email string) (*model.Employee, error)
	Update(emp *model.Employee) error
	UpdateFields(id uint, fields map[string]interface{}) error
	Delete(id uint) error
	CountByDepartment(deptID uint) (int64, error)
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
	return r.db.Create(&emp).Error
}

// Tìm nhân viên với phân trang, tìm kiếm, preload quan hệ
func (r *employeeRepository) FindAll(query model.PaginationQuery) ([]model.Employee, int64, error) {
	var employee []model.Employee
	var total int64

	db := r.db.Model(&model.Employee{})
	// Tìm theo tên hoặc email
	if query.Search != "" {
		db = db.Where(
			"first_name LIKE ? OR last_name LIKE ? OR email LIKE ?",
			"%"+query.Search+"%",
			"%"+query.Search+"%",
			"%"+query.Search+"%",
		)
	}

	db.Count(&total)

	// offset = (page - 1) * limit
	offset := (query.Page - 1) * query.Limit

	// Preload Department đê trả ra thông tin phòng ban
	err := db.Preload("Department").Offset(offset).Limit(query.Limit).Order("created_at DESC").Find(&employee).Error

	return employee, total, err
}

func (r *employeeRepository) FindByID(id uint) (*model.Employee, error) {
	var employee model.Employee
	err := r.db.Preload("Department").First(&employee, id).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil

}
func (r *employeeRepository) FindByEmail(email string) (*model.Employee, error) {
	var employee model.Employee
	err := r.db.Preload("Department").Where("email = ?", email).First(&employee).Error
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
	return r.db.Model(&model.Employee{}).Where("id = ?", id).Updates(fields).Error
}

func (r *employeeRepository) Delete(id uint) error {
	return r.db.Delete(&model.Employee{}, id).Error
}

func (r *employeeRepository) CountByDepartment(deptID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Employee{}).Where("department_id = ?", deptID).Count(&count).Error
	return count, err
}
