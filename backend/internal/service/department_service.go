package service

import (
	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/internal/repository"
	"context"

	"errors"
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm"
)

// DepartmentService interface - contract cho quản lý phòng ban
type DepartmentService interface {
	CreateDepartment(req model.CreateDepartmentRequest) (*model.Department, error)
	GetDepartments(query model.PaginationQuery) (*model.PaginatedResult, error)
	GetDepartmentByID(id uint) (*model.Department, error)
	UpdateDepartment(id uint, req model.UpdateDepartmentRequest) (*model.Department, error)
	DeleteDepartment(id uint) error
}

// --- Department Service Implementation ---

type departmentService struct {
	db       *gorm.DB
	deptRepo repository.DepartmentRepository
	empRepo  repository.EmployeeRepository
	cacheSvc CacheService
}

func NewDepartmentService(db *gorm.DB, deptRepo repository.DepartmentRepository, empRepo repository.EmployeeRepository, cacheSvc CacheService) DepartmentService {
	return &departmentService{
		db:       db,
		deptRepo: deptRepo,
		empRepo:  empRepo,
		cacheSvc: cacheSvc,
	}
}

func (ds *departmentService) CreateDepartment(req model.CreateDepartmentRequest) (*model.Department, error) {
	// Chuẩn hoá dữ liệu
	req.Name = strings.TrimSpace(req.Name)
	req.Code = strings.TrimSpace(strings.ToUpper(req.Code))
	req.Description = strings.TrimSpace(req.Description)

	// Kiểm tra tên phòng ban đã tồn tại chưa
	if existingDepts, _, err := ds.deptRepo.FindAll(model.PaginationQuery{Page: 1, Limit: 1, Search: req.Name}); err == nil {
		for _, d := range existingDepts {
			if strings.EqualFold(d.Name, req.Name) {
				return nil, fmt.Errorf("Tên phòng ban '%s' đã tồn tại!", req.Name)
			}
		}
	}

	// Kiểm tra code phòng ban đã tồn tại chưa
	if _, err := ds.deptRepo.FindByCode(req.Code); err == nil {
		return nil, fmt.Errorf("Mã phòng ban '%s' đã tồn tại!", req.Code)
	}

	if req.ManagerID != nil {
		if *req.ManagerID == 0 {
			return nil, errors.New("ID quản lý phải lớn hơn 0")
		}
		return nil, errors.New("Không thể chỉ định trưởng phòng khi tạo phòng ban mới. Vui lòng tạo phòng ban trước rồi gán trưởng phòng sau.")
	}

	dept := &model.Department{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		ManagerID:   nil,
	}

	if err := ds.deptRepo.Create(dept); err != nil {
		return nil, fmt.Errorf("Lỗi khi tạo phòng ban: %w", err)
	}

	// Invalidate dashboard stats cache
	_ = ds.cacheSvc.Delete(context.Background(), "dashboard:stats")

	return dept, nil
}

func (ds *departmentService) GetDepartments(query model.PaginationQuery) (*model.PaginatedResult, error) {
	// Chuẩn hoá phân trang
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 10
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	depts, total, err := ds.deptRepo.FindAll(query)
	if err != nil {
		return nil, fmt.Errorf("Lỗi khi lấy thông tin phòng ban: %w", err)
	}
	totalPages := int(math.Ceil(float64(total) / float64(query.Limit)))

	return &model.PaginatedResult{
		Items:      depts,
		Total:      total,
		Page:       query.Page,
		Limit:      query.Limit,
		TotalPages: totalPages,
	}, nil
}

func (ds *departmentService) GetDepartmentByID(id uint) (*model.Department, error) {
	if id == 0 {
		return nil, errors.New("ID phòng ban phải lớn hơn 0")
	}

	dept, err := ds.deptRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Không tìm thấy phòng ban này")
		}
		return nil, err
	}
	return dept, nil
}

func (ds *departmentService) UpdateDepartment(id uint, req model.UpdateDepartmentRequest) (*model.Department, error) {
	if id == 0 {
		return nil, errors.New("ID phòng ban phải lớn hơn 0")
	}

	var updatedDept *model.Department
	if err := ds.db.Transaction(func(tx *gorm.DB) error {
		txDeptRepo := ds.deptRepo.WithTx(tx)
		txEmpRepo := ds.empRepo.WithTx(tx)

		dept, err := txDeptRepo.FindByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("Không tìm thấy phòng ban này!")
			}
			return fmt.Errorf("Lỗi khi tìm phòng ban: %w", err)
		}

		updated := false

		if req.Name != nil {
			name := strings.TrimSpace(*req.Name)
			if name == "" {
				return errors.New("Tên phòng ban không được để trống")
			}

			if !strings.EqualFold(name, dept.Name) {
				existingDepts, _, err := txDeptRepo.FindAll(model.PaginationQuery{
					Page: 1, Limit: 1, Search: name,
				})
				if err == nil {
					for _, d := range existingDepts {
						if strings.EqualFold(d.Name, name) && d.ID != id {
							return fmt.Errorf("Tên phòng ban '%s' đã tồn tại!", name)
						}
					}
				}
				dept.Name = name
				updated = true
			}
		}

		if req.Description != nil {
			dept.Description = strings.TrimSpace(*req.Description)
			updated = true
		}

		if req.ManagerID != nil {
			if *req.ManagerID == 0 {
				dept.ManagerID = nil
				dept.Manager = nil
				updated = true
			} else {
				emp, err := txEmpRepo.FindByID(*req.ManagerID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return errors.New("Không tìm thấy nhân viên được chỉ định làm quản lý!")
					}
					return fmt.Errorf("Lỗi khi tìm nhân viên quản lý: %w", err)
				}

				if emp.DepartmentID != id {
					return errors.New("Trưởng phòng phải thuộc chính phòng ban này")
				}

				existingDept, err := txDeptRepo.FindByManagerID(*req.ManagerID)
				if err == nil && existingDept.ID != id {
					return errors.New("Nhân viên này đã là trưởng phòng của phòng khác")
				}
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}

				dept.ManagerID = req.ManagerID
				updated = true
			}
		}

		if !updated {
			updatedDept = dept
			return nil
		}

		if err := txDeptRepo.Update(dept); err != nil {
			return fmt.Errorf("Lỗi khi cập nhật phòng ban: %w", err)
		}

		updatedDept = dept
		return nil
	}); err != nil {
		return nil, err
	}

	// Invalidate dashboard stats cache
	_ = ds.cacheSvc.Delete(context.Background(), "dashboard:stats")

	return updatedDept, nil
}


func (ds *departmentService) DeleteDepartment(id uint) error {
	if id == 0 {
		return errors.New("ID phòng ban phải lớn hơn 0")
	}

	err := ds.db.Transaction(func(tx *gorm.DB) error {
		txDeptRepo := ds.deptRepo.WithTx(tx)
		txEmpRepo := ds.empRepo.WithTx(tx)

		dept, err := txDeptRepo.FindByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("Không tìm thấy phòng ban này!")
			}
			return fmt.Errorf("Lỗi khi tìm phòng ban: %w", err)
		}

		count, err := txEmpRepo.CountByDepartment(id)
		if err != nil {
			return fmt.Errorf("Có lỗi khi kiểm tra nhân viên tại phòng ban này: %w", err)
		}
		if count > 0 {
			return fmt.Errorf("Không thể xoá phòng ban vì còn %d nhân viên đang hoạt động!", count)
		}

		if dept.ManagerID != nil {
			dept.ManagerID = nil
			dept.Manager = nil
			if err := txDeptRepo.Update(dept); err != nil {
				return fmt.Errorf("Lỗi khi xoá trưởng phòng trước khi xoá phòng ban: %w", err)
			}
		}

		if err := txDeptRepo.Delete(id); err != nil {
			return fmt.Errorf("Có lỗi khi xoá phòng ban này: %w", err)
		}

		return nil
	})

	if err == nil {
		// Invalidate dashboard stats cache
		_ = ds.cacheSvc.Delete(context.Background(), "dashboard:stats")
	}

	return err
}
