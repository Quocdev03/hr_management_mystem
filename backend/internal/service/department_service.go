package service

import (
	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/internal/repository"
	"chiquoc_hocgolang/package/validation"
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
	deptRepo repository.DepartmentRepository
	empRepo  repository.EmployeeRepository
}

func NewDepartmentService(deptRepo repository.DepartmentRepository, empRepo repository.EmployeeRepository) DepartmentService {
	return &departmentService{
		deptRepo: deptRepo,
		empRepo:  empRepo,
	}
}

func (ds *departmentService) CreateDepartment(req model.CreateDepartmentRequest) (*model.Department, error) {
	// Chuẩn hoá dữ liệu
	req.Name = strings.TrimSpace(req.Name)
	req.Code = strings.TrimSpace(strings.ToUpper(req.Code))
	req.Description = strings.TrimSpace(req.Description)

	// Validate đầu vào (defense-in-depth)
	if verrs := validation.ValidateCreateDepartment(req.Name, req.Code, req.Description); verrs != nil {
		return nil, errors.New(verrs.Error())
	}

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

	dept := &model.Department{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
	}

	if err := ds.deptRepo.Create(dept); err != nil {
		return nil, fmt.Errorf("Lỗi khi tạo phòng ban: %w", err)
	}
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

	// Chuẩn hoá dữ liệu
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)

	// Validate đầu vào (defense-in-depth)
	if verrs := validation.ValidateUpdateDepartment(req.Name, req.Description, req.ManagerID); verrs != nil {
		return nil, errors.New(verrs.Error())
	}

	dept, err := ds.deptRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Không tìm thấy phòng ban này!")
		}
		return nil, fmt.Errorf("Lỗi khi tìm phòng ban: %w", err)
	}

	// Kiểm tra nếu đổi tên, tên mới không trùng với phòng ban khác
	if req.Name != "" && !strings.EqualFold(req.Name, dept.Name) {
		if existingDepts, _, err := ds.deptRepo.FindAll(model.PaginationQuery{Page: 1, Limit: 1, Search: req.Name}); err == nil {
			for _, d := range existingDepts {
				if strings.EqualFold(d.Name, req.Name) && d.ID != id {
					return nil, fmt.Errorf("Tên phòng ban '%s' đã tồn tại!", req.Name)
				}
			}
		}
		dept.Name = req.Name
	}

	if req.Description != "" {
		dept.Description = req.Description
	}
	if req.ManagerID != nil {
		if *req.ManagerID == 0 {
			return nil, errors.New("ID quản lý phải lớn hơn 0")
		}
		manager, err := ds.empRepo.FindByID(*req.ManagerID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("Không tìm thấy nhân viên được chỉ định làm quản lý!")
			}
			return nil, fmt.Errorf("Lỗi khi tìm nhân viên quản lý: %w", err)
		}
		if manager.DepartmentID != id {
			return nil, errors.New("Nhân viên quản lý phải thuộc phòng ban này!")
		}
		dept.ManagerID = req.ManagerID
	}

	if err := ds.deptRepo.Update(dept); err != nil {
		return nil, fmt.Errorf("Lỗi khi cập nhật phòng ban: %w", err)
	}
	return dept, nil
}

func (ds *departmentService) DeleteDepartment(id uint) error {
	if id == 0 {
		return errors.New("ID phòng ban phải lớn hơn 0")
	}

	// Kiểm tra phòng ban tồn tại
	if _, err := ds.deptRepo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Không tìm thấy phòng ban này!")
		}
		return fmt.Errorf("Lỗi khi tìm phòng ban: %w", err)
	}

	// Không cho xoá phòng ban khi còn nhân viên
	count, err := ds.empRepo.CountByDepartment(id)
	if err != nil {
		return fmt.Errorf("Có lỗi khi kiểm tra nhân viên tại phòng ban này: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("Không thể xoá phòng ban vì còn %d nhân viên đang hoạt động!", count)
	}

	if err := ds.deptRepo.Delete(id); err != nil {
		return fmt.Errorf("Có lỗi khi xoá phòng ban này")
	}
	return nil
}
