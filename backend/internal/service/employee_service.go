package service

import (
	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/internal/repository"
	"chiquoc_hocgolang/package/validation"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"
)

// EmployeeService interface - contract cho quản lý nhân viên
type EmployeeService interface {
	Create(req model.CreateEmployeeRequest) (*model.Employee, error)
	GetEmployees(query model.PaginationQuery) (*model.PaginatedResult, error)
	GetEmployeeByID(id uint) (*model.Employee, error)
	UpdateEmployee(id uint, req model.UpdateEmployeeRequest) (*model.Employee, error)
	DeleteEmployee(id uint) error
}

// --- Employee Service Implementation ---

type employeeService struct {
	empRepo  repository.EmployeeRepository
	deptRepo repository.DepartmentRepository
}

func NewEmployeeService(empRepo repository.EmployeeRepository, deptRepo repository.DepartmentRepository) EmployeeService {
	return &employeeService{
		empRepo:  empRepo,
		deptRepo: deptRepo,
	}
}

func (es *employeeService) Create(req model.CreateEmployeeRequest) (*model.Employee, error) {
	// Chuẩn hoá dữ liệu
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Phone = strings.TrimSpace(req.Phone)
	req.Position = strings.TrimSpace(req.Position)
	req.JoinDate = strings.TrimSpace(req.JoinDate)

	// Validate đầu vào (defense-in-depth)
	if verrs := validation.ValidateCreateEmployee(
		req.DepartmentID, req.FirstName, req.LastName,
		req.Email, req.Phone, req.Position, req.JoinDate, req.Salary,
	); verrs != nil {
		return nil, errors.New(verrs.Error())
	}

	// Kiểm tra phòng ban tồn tại
	if _, err := es.deptRepo.FindByID(req.DepartmentID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Không tìm thấy phòng ban này")
		}
		return nil, fmt.Errorf("Lỗi kiểm tra phòng ban: %w", err)
	}

	// Kiểm tra email tồn tại
	if _, err := es.empRepo.FindByEmail(req.Email); err == nil {
		return nil, errors.New("Email đã tồn tại")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("Lỗi kiểm tra email: %w", err)
	}

	// Parse ngày vào làm, mặc định là vừa tạo
	joinDate := time.Now()
	if req.JoinDate != "" {
		parsed, err := time.Parse("2006-01-02", req.JoinDate)
		if err != nil {
			return nil, errors.New("Ngày vào làm không đúng định dạng, sử dụng YYYY-MM-DD")
		}
		// Kiểm tra ngày không ở tương lai
		if parsed.After(time.Now()) {
			return nil, errors.New("Ngày vào làm không được là ngày trong tương lai")
		}
		joinDate = parsed
	}

	if req.Salary < 0 {
		return nil, errors.New("Mức lương không được nhỏ hơn 0")
	}

	emp := &model.Employee{
		DepartmentID: req.DepartmentID,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Phone:        req.Phone,
		Position:     req.Position,
		Salary:       req.Salary,
		JoinDate:     joinDate,
		Status:       "active",
	}

	// *sql.DB (GORM) đã thread-safe, không cần mutex
	if err := es.empRepo.Create(emp); err != nil {
		return nil, fmt.Errorf("Tạo nhân viên không thành công: %w", err)
	}

	return es.empRepo.FindByID(emp.ID)
}

// Danh sách nhân viên có phân trang và tìm kiếm
func (es *employeeService) GetEmployees(query model.PaginationQuery) (*model.PaginatedResult, error) {
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

	employees, total, err := es.empRepo.FindAll(query)
	if err != nil {
		return nil, fmt.Errorf("Lấy danh sách nhân viên bị lỗi: %w", err)
	}
	totalPage := int(math.Ceil(float64(total) / float64(query.Limit)))

	return &model.PaginatedResult{
		Items:      employees,
		Total:      total,
		Page:       query.Page,
		Limit:      query.Limit,
		TotalPages: totalPage,
	}, nil
}

func (es *employeeService) GetEmployeeByID(id uint) (*model.Employee, error) {
	if id == 0 {
		return nil, errors.New("ID nhân viên phải lớn hơn 0")
	}

	emp, err := es.empRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Không tìm thấy nhân viên này")
		}
		return nil, err
	}
	return emp, nil
}

func (es *employeeService) UpdateEmployee(id uint, req model.UpdateEmployeeRequest) (*model.Employee, error) {
	if id == 0 {
		return nil, errors.New("ID nhân viên phải lớn hơn 0")
	}

	// Chuẩn hoá dữ liệu
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Position = strings.TrimSpace(req.Position)
	req.Status = strings.TrimSpace(strings.ToLower(req.Status))

	// Validate đầu vào (defense-in-depth)
	if verrs := validation.ValidateUpdateEmployee(
		req.DepartmentID, req.FirstName, req.LastName,
		req.Phone, req.Position, req.Status, req.Salary,
	); verrs != nil {
		return nil, errors.New(verrs.Error())
	}

	emp, err := es.empRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Không tìm thấy nhân viên")
		}
		return nil, fmt.Errorf("Lỗi khi tìm nhân viên: %w", err)
	}

	// Chỉ update các field được truyền vào
	if req.DepartmentID != 0 {
		if _, err := es.deptRepo.FindByID(req.DepartmentID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("Không tìm thấy phòng ban mới")
			}
			return nil, fmt.Errorf("Lỗi kiểm tra phòng ban: %w", err)
		}
		emp.DepartmentID = req.DepartmentID
	}
	if req.FirstName != "" {
		emp.FirstName = req.FirstName
	}
	if req.LastName != "" {
		emp.LastName = req.LastName
	}
	if req.Phone != "" {
		emp.Phone = req.Phone
	}
	if req.Position != "" {
		emp.Position = req.Position
	}
	if req.Salary != 0 {
		if req.Salary < 0 {
			return nil, errors.New("Mức lương không được nhỏ hơn 0")
		}
		emp.Salary = req.Salary
	}
	if req.Status != "" {
		if req.Status != "active" && req.Status != "inactive" {
			return nil, errors.New("Trạng thái chỉ có thể là 'active' hoặc 'inactive'")
		}
		emp.Status = req.Status
	}

	if err := es.empRepo.Update(emp); err != nil {
		return nil, fmt.Errorf("Cập nhật thông tin nhân viên bị lỗi: %w", err)
	}

	return es.empRepo.FindByID(id)
}

func (es *employeeService) DeleteEmployee(id uint) error {
	if id == 0 {
		return errors.New("ID nhân viên phải lớn hơn 0")
	}

	if _, err := es.empRepo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Không tìm thấy nhân viên này!")
		}
		return fmt.Errorf("Lỗi khi tìm nhân viên: %w", err)
	}
	if err := es.empRepo.Delete(id); err != nil {
		return fmt.Errorf("Lỗi khi xoá nhân viên này: %w", err)
	}
	return nil
}
