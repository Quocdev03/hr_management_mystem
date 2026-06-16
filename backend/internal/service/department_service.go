package service

import (
	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/internal/repository"
	"chiquoc_hocgolang/internal/utils"
	"context"

	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/redis/go-redis/v9"
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
	rdb      *redis.Client
}

func NewDepartmentService(db *gorm.DB, deptRepo repository.DepartmentRepository, empRepo repository.EmployeeRepository, rdb *redis.Client) DepartmentService {
	return &departmentService{
		db:       db,
		deptRepo: deptRepo,
		empRepo:  empRepo,
		rdb:      rdb,
	}
}

func (ds *departmentService) CreateDepartment(req model.CreateDepartmentRequest) (*model.Department, error) {
	// Chuẩn hoá dữ liệu
	req.Name = strings.TrimSpace(req.Name)
	req.Code = strings.TrimSpace(strings.ToUpper(req.Code))
	req.Description = strings.TrimSpace(req.Description)

	// Kiểm tra tên phòng ban đã tồn tại chưa
	if existingDepts, _, err := ds.deptRepo.FindAll(model.PaginationQuery{Page: 1, Limit: 100, Search: req.Name}); err == nil {
		for _, d := range existingDepts {
			if strings.EqualFold(d.Name, req.Name) {
				return nil, fmt.Errorf("tên phòng ban '%s' đã tồn tại!", req.Name)
			}
		}
	}

	// Kiểm tra code phòng ban đã tồn tại chưa
	if _, err := ds.deptRepo.FindByCode(req.Code); err == nil {
		return nil, fmt.Errorf("mã phòng ban '%s' đã tồn tại!", req.Code)
	}

	if req.ManagerID != nil {
		if *req.ManagerID == 0 {
			return nil, errors.New("id quản lý phải lớn hơn 0")
		}
		return nil, errors.New("không thể chỉ định trưởng phòng khi tạo phòng ban mới. Vui lòng tạo phòng ban trước rồi gán trưởng phòng sau.")
	}

	dept := &model.Department{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		ManagerID:   nil,
	}

	if err := ds.deptRepo.Create(dept); err != nil {
		return nil, fmt.Errorf("lỗi khi tạo phòng ban: %w", err)
	}

	// Invalidate dashboard stats cache
	_ = utils.InvalidateDashboardStats(context.Background(), ds.rdb)

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
		return nil, fmt.Errorf("lỗi khi lấy thông tin phòng ban: %w", err)
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
		return nil, errors.New("id phòng ban phải lớn hơn 0")
	}

	dept, err := ds.deptRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("không tìm thấy phòng ban này")
		}
		return nil, err
	}
	return dept, nil
}

func (ds *departmentService) UpdateDepartment(id uint, req model.UpdateDepartmentRequest) (*model.Department, error) {
	if id == 0 {
		return nil, errors.New("id phòng ban phải lớn hơn 0")
	}

	var updatedDept *model.Department
	if err := ds.db.Transaction(func(tx *gorm.DB) error {
		txDeptRepo := ds.deptRepo.WithTx(tx)
		txEmpRepo := ds.empRepo.WithTx(tx)

		dept, err := txDeptRepo.FindByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("không tìm thấy phòng ban này")
			}
			return fmt.Errorf("lỗi khi tìm phòng ban: %w", err)
		}

		updates := map[string]interface{}{}

		if req.Name != nil {
			name := strings.TrimSpace(*req.Name)
			if name == "" {
				return errors.New("tên phòng ban không được để trống")
			}

			if !strings.EqualFold(name, dept.Name) {
				existingDepts, _, err := txDeptRepo.FindAll(model.PaginationQuery{
					Page: 1, Limit: 100, Search: name,
				})
				if err == nil {
					for _, d := range existingDepts {
						if strings.EqualFold(d.Name, name) && d.ID != id {
							return fmt.Errorf("tên phòng ban '%s' đã tồn tại!", name)
						}
					}
				}
				updates["name"] = name
			}
		}

		if req.Description != nil {
			value := strings.TrimSpace(*req.Description)
			if value != dept.Description {
				updates["description"] = value
			}
		}

		if req.ManagerID != nil {
			currentManagerID := uint(0)
			if dept.ManagerID != nil {
				currentManagerID = *dept.ManagerID
			}
			newManagerID := *req.ManagerID
			if newManagerID != currentManagerID {
				if newManagerID == 0 {
					updates["manager_id"] = gorm.Expr("NULL")
				} else {
					emp, err := txEmpRepo.FindByID(*req.ManagerID)
					if err != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) {
							return errors.New("không tìm thấy nhân viên được chỉ định làm quản lý")
						}
						return fmt.Errorf("lỗi khi tìm nhân viên quản lý: %w", err)
					}

					if emp.DepartmentID != id {
						return errors.New("trưởng phòng phải thuộc chính phòng ban này")
					}

					existingDept, err := txDeptRepo.FindByManagerID(newManagerID)
					if err == nil && existingDept.ID != id {
						return errors.New("nhân viên này đã là trưởng phòng của phòng khác")
					}
					if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
						return err
					}

					updates["manager_id"] = newManagerID
				}
			}
		}

		if len(updates) == 0 {
			updatedDept = dept
			return nil
		}

		if err := txDeptRepo.Update(id, updates); err != nil {
			return fmt.Errorf("lỗi khi cập nhật phòng ban: %w", err)
		}

		// Reload lại từ DB để lấy dữ liệu mới nhất (kể cả Manager preloaded)
		reloaded, err := txDeptRepo.FindByID(id)
		if err != nil {
			return fmt.Errorf("lỗi reload phòng ban sau cập nhật: %w", err)
		}
		updatedDept = reloaded
		return nil
	}); err != nil {
		return nil, err
	}

	// Invalidate dashboard stats cache
	_ = utils.InvalidateDashboardStats(context.Background(), ds.rdb)

	return updatedDept, nil
}

func (ds *departmentService) DeleteDepartment(id uint) error {
	if id == 0 {
		return errors.New("id phòng ban phải lớn hơn 0")
	}

	err := ds.db.Transaction(func(tx *gorm.DB) error {
		txDeptRepo := ds.deptRepo.WithTx(tx)
		txEmpRepo := ds.empRepo.WithTx(tx)

		dept, err := txDeptRepo.FindByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("không tìm thấy phòng ban này")
			}
			return fmt.Errorf("lỗi khi tìm phòng ban: %w", err)
		}

		count, err := txEmpRepo.CountByDepartment(id)
		if err != nil {
			return fmt.Errorf("có lỗi khi kiểm tra nhân viên tại phòng ban này: %w", err)
		}
		if count > 0 {
			return fmt.Errorf("không thể xoá phòng ban vì còn %d nhân viên đang hoạt động", count)
		}

		if dept.ManagerID != nil {
			if err := txDeptRepo.UpdateManager(id, nil); err != nil {
				return fmt.Errorf("lỗi khi xoá trưởng phòng trước khi xoá phòng ban: %w", err)
			}
		}

		if err := txDeptRepo.Delete(id); err != nil {
			return fmt.Errorf("có lỗi khi xoá phòng ban này: %w", err)
		}

		return nil
	})

	if err == nil {
		// Invalidate dashboard stats cache
		_ = utils.InvalidateDashboardStats(context.Background(), ds.rdb)
	}

	return err
}
