package service

import (
	"chiquoc_hocgolang/internal/config"
	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/internal/repository"
	"chiquoc_hocgolang/package/utils"
	"errors"
	"fmt"
	"math"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Định nghĩa contract
type AuthService interface{}

type EmployeeService interface{}

type DepartmentService interface{}

// Định nghĩa các Implementations
// --- Auth Service Implementation ---
type authServive struct {
	useRepo repository.UserRepository
	jwtCfg  *config.JWTConfig
}

// Login và trả về JWT Token
func (au *authServive) Login(req *model.LoginRequest) (*model.LoginResponse, error) {
	// Tìm user theo email
	user, err := au.useRepo.FindByEmail(req.Email)
	if err != nil {
		// Trả lỗi chung không tiết lộ email hay pass
		return nil, errors.New("Email hoặc mật khẩu không hợp lệ!")
	}

	// Kiểm tra xem tài khoản có bị khoá không
	if !user.IsActive {
		return nil, errors.New("Tài khoản đã bị vô hiệu hoá!")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("Email hoặc mật khẩu không hợp lệ!")
	}

	// Tạo jwt token chứa thông tin user và role
	token, err := utils.GenerateToken(
		user.ID,
		user.Email,
		user.RoleID,
		user.Role.Name,
		au.jwtCfg.SecretKey,
		au.jwtCfg.ExpireHour,
	)
	if err != nil {
		return nil, fmt.Errorf("Có lỗi khi tạo token: %w", err)
	}

	return &model.LoginResponse{
		AccessToken: token,
		User:        *user,
	}, nil
}

// Employee Service

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

func (es *employeeService) Create(req *model.CreateEmployeeRequest) (*model.Employee, error) {
	// Kiểm tra phòng ban tồn tại
	if _, err := es.deptRepo.FindByID(req.DepartmentID); err != nil {
		return nil, errors.New("Không tìm thấy phòng ban này")
	}

	// Kiểm tra email tồn tại
	if _, err := es.empRepo.FindByEmail(req.Email); err == nil {
		return nil, errors.New("Email này đã tồn tại")
	}

	// Parse ngày vào làm, mặc định là vừa tạo
	joinDate := time.Now()
	if req.JoinDate != "" {
		parsed, err := time.Parse("2006-01-02", req.JoinDate)
		if err != nil {
			return nil, errors.New("Không đúng định dạng, sử dụng YYYY-MM-DD")
		}
		joinDate = parsed
	}
	emp := &model.Employee{
		DepartmentID: req.DepartmentID,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
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
	emp, err := es.empRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("Không tìm thấy nhân viên")
	}

	// Chỉ update các field được truyền vào
	if req.DepartmentID != 0 {
		if _, err := es.deptRepo.FindByID(req.DepartmentID); err != nil {
			return nil, errors.New("Không tìm thấy phòng ban")
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
	if req.Salary > 0 {
		emp.Salary = req.Salary
	}
	if req.Status != "" {
		emp.Status = req.Status
	}

	if err := es.empRepo.Update(emp); err != nil {
		return nil, fmt.Errorf("Cập nhật thông tin nhân viên bị lỗi: %w", err)
	}

	return es.empRepo.FindByID(id)
}

func (es *employeeService) DeleteEmployee(id uint) error {
	if _, err := es.empRepo.FindByID(id); err != nil {
		return errors.New("Không tìm thấy nhân viên này!")
	}
	if err := es.empRepo.Delete(id); err != nil {
		return fmt.Errorf("Lỗi khi xoá nhân viên này: %w", err)
	}
	return nil
}

// Department Service

type departmentService struct {
	deptRepo repository.DepartmentRepository
	empRepo  repository.EmployeeRepository
}

func NewDepartmentRepository(deptRepo repository.DepartmentRepository, empRepo repository.EmployeeRepository) departmentService {
	return departmentService{
		deptRepo: deptRepo,
		empRepo:  empRepo,
	}
}

func (ds *departmentService) CreateDepartment(req model.CreateDepartmentRequest) (*model.Department, error) {
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
	// Gọi thẳng, không cần goroutine vì đây là sequential request/response
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
	dept, err := ds.deptRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("Không tìm thấy phòng ban này!")
	}
	if req.Name != "" {
		dept.Name = req.Name
	}
	if req.Description != "" {
		dept.Description = req.Description
	}
	if req.ManagerID != nil {
		dept.ManagerID = req.ManagerID
	}

	if err := ds.deptRepo.Update(dept); err != nil {
		return nil, fmt.Errorf("Lỗi khi cập nhật phòng ban: %w", err)
	}
	return dept, nil
}
func (ds *departmentService) DeleteDepartment(id uint) error {
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
