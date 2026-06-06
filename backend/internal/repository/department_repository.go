package repository

import (
	"chiquoc_hocgolang/internal/model"

	"gorm.io/gorm"
)

// DepartmentRepository interface cho các thao tác với bảng departments
type DepartmentRepository interface {
	WithTx(tx *gorm.DB) DepartmentRepository
	Create(dept *model.Department) error
	FindAll(query model.PaginationQuery) ([]model.Department, int64, error)
	FindByID(id uint) (*model.Department, error)
	FindByCode(code string) (*model.Department, error)
	FindByManagerID(managerID uint) (*model.Department, error)
	Update(dept *model.Department) error
	UpdateManager(departmentID uint, managerID *uint) error
	Delete(id uint) error
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

	db := r.db.Model(&model.Department{}).
		Preload("Manager").
		Preload("Manager.User").
		Preload("Manager.Department")

	// Tìm kiếm theo tên hoặc mã phòng ban
	if query.Search != "" {
		db = db.Where(
			"name LIKE ? OR code LIKE ?",
			"%"+query.Search+"%",
			"%"+query.Search+"%",
		)
	}

	// Lấy tổng số record sau khi filter
	db.Count(&total)

	// Phân trang: offset = (page-1)  * limit
	offset := (query.Page - 1) * query.Limit
	err := db.
		Offset(offset).
		Limit(query.Limit).
		Find(&department).Error
	return department, total, err
}

func (r *departmentRepository) FindByID(id uint) (*model.Department, error) {
	var dept model.Department
	err := r.db.Preload("Manager").First(&dept, id).Error
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
	return &dept, nil
}

func (r *departmentRepository) FindByManagerID(managerID uint) (*model.Department, error) {
	var dept model.Department
	err := r.db.Where("manager_id = ?", managerID).First(&dept).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *departmentRepository) Update(dept *model.Department) error {
	updates := map[string]interface{}{
		"name":        dept.Name,
		"description": dept.Description,
	}

	// Phải xử lý manager_id riêng biệt:
	// - GORM Updates() bỏ qua nil values trong map → manager_id NULL sẽ không được ghi
	// - Dùng gorm.Expr("NULL") để force ghi NULL xuống DB
	if dept.ManagerID == nil {
		updates["manager_id"] = gorm.Expr("NULL")
	} else {
		updates["manager_id"] = *dept.ManagerID
	}

	// Dùng Where("id = ?") thay vì Model(struct) để GORM luôn target đúng record
	// tránh trường hợp GORM không resolve được primary key từ struct state
	return r.db.Model(&model.Department{}).
		Where("id = ?", dept.ID).
		Updates(updates).Error
}

func (r *departmentRepository) UpdateManager(departmentID uint, managerID *uint) error {
	updates := map[string]interface{}{
		"manager_id": nil,
	}
	if managerID != nil {
		updates["manager_id"] = *managerID
	}

	return r.db.Model(&model.Department{}).
		Where("id = ?", departmentID).
		Updates(updates).Error
}

func (r *departmentRepository) Delete(id uint) error {
	// Soft delete: GORM tự set deleted_at thay vì xóa thật
	return r.db.Delete(&model.Department{}, id).Error
}

func (r *departmentRepository) WithTx(tx *gorm.DB) DepartmentRepository {
	return &departmentRepository{db: tx}
}
