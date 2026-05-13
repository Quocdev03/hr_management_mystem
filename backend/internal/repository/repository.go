package repository

import (
	"chiquoc_hocgolang/internal/model"

	"gorm.io/gorm"
)

// Định nghĩa contract
// UserRepository interface cho các thao tác với bảng users
type UserRepository interface {
	Create(user *model.User) error
	FindByID(id uint) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Update(user *model.User) error
}

// DepartmentRepository interface cho các thao tác với bảng departments
type DepartmentRepository interface {
	Create(dept *model.Department) error
	FindAll(query model.PaginationQuery) ([]model.Department, int64, error)
	FindByID(id uint) (*model.Department, error)
	FindByCode(code string) (*model.Department, error)
	Update(dept *model.Department) error
	Delete(id uint) error
}

// EmployeeRepository interface cho các thao tác với bảng employees

type EmployeeRepository interface {
	Create(emp *model.Employee) error
	FindAll(query model.PaginationQuery) ([]model.Employee, int64, error)
	FindByID(id uint) (*model.Employee, error)
	FindByEmail(email string) (*model.Employee, error)
	Update(emp *model.Employee) error
	Delete(id uint) error
	CountByDepartment(deptID uint) (int64, error)
}

// Định nghĩa các Implementations
// --- UserRepository Implementation ---
type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.
		Preload("Role").
		Where("id = ?", id).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.
		Preload("Role").
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Role").Where("email = ?  ", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// --- DepartmentRepository Implementation ---

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{
		db: db,
	}
}

func (r *departmentRepository) Create(dept *model.Department) error {
	return r.db.Create(dept).Error
}

// Tìm all phòng ban, phân trang, tìm kiếm
func (r *departmentRepository) FindAll(query model.PaginationQuery) ([]model.Department, int64, error) {
	var department []model.Department
	var total int64

	db := r.db.Model(&model.Department{})

	// Tìm kiếm theo tên hoặc mã phòng ban
	if query.Search != "" {
		db.Where(
			"name LIKE ? OR code LIKE ?",
			"%"+query.Search+"%",
			"%"+query.Search+"%",
		)
	}

	// Lấy tổng trang
	db.Count(&total)

	// Phân trang: offset = (page-1)  * limit
	offset := (query.Page - 1) * query.Limit
	err := db.Offset(offset).Limit(query.Limit).Find(&department).Error
	return department, total, err
}

func (r *departmentRepository) FindByID(id uint) (*model.Department, error) {
	var dept model.Department
	err := r.db.First(&dept, id).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *departmentRepository) FindByCode(code string) (*model.Department, error) {
	var dept model.Department
	err := r.db.Where("code = ?", code).First(&dept).Error
	if err != nil {
		return nil, err
	}
	return &dept, err
}

func (r *departmentRepository) Update(dept *model.Department) error {
	return r.db.Save(&dept).Error
}

func (r *departmentRepository) Delete(id uint) error {
	// Soft delete: GORM tự set deleted_at thay vì xóa thật
	return r.db.Delete(&model.Department{}, id).Error
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
	err := r.db.Where("email = ?", email).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepository) Update(emp *model.Employee) error {
	return r.db.Save(emp).Error
}

func (r *employeeRepository) Delete(id uint) error {
	return r.db.Delete(&model.Employee{}, id).Error
}

func (r *employeeRepository) CountByDepartment(deptID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Employee{}).Where("department_id = ?", deptID).Count(&count).Error
	return count, err
}
