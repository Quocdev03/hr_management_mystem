package service

import (
	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/internal/repository"
	"chiquoc_hocgolang/internal/utils"
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
	req.BirthDate = strings.TrimSpace(req.BirthDate)
	req.Gender = strings.TrimSpace(strings.ToLower(req.Gender))

	// Validate đầu vào (defense-in-depth)
	ve := &utils.ValidationErrors{}
	if req.DepartmentID == 0 {
		ve.Add(utils.FieldDepartmentID, "Phòng ban là bắt buộc")
	}
	utils.CheckName(ve, utils.FieldFirstName, "Họ", req.FirstName, 2, 100)
	utils.CheckName(ve, utils.FieldLastName, "Tên", req.LastName, 2, 100)
	utils.CheckEmail(ve, req.Email)
	utils.CheckPhone(ve, req.Phone)
	if req.Position != "" {
		l := len([]rune(req.Position))
		if l < 2 || l > 100 {
			ve.Add(utils.FieldPosition, "Vị trí phải từ 2 đến 100 ký tự")
		}
	}
	if req.Salary < 0 {
		ve.Add(utils.FieldSalary, "Mức lương không được nhỏ hơn 0")
	}
	if req.JoinDate != "" {
		if parsed, err := time.Parse("2006-01-02", req.JoinDate); err != nil {
			ve.Add(utils.FieldJoinDate, "Ngày vào làm phải đúng định dạng YYYY-MM-DD")
		} else if parsed.After(time.Now()) {
			ve.Add(utils.FieldJoinDate, "Ngày vào làm không được là ngày trong tương lai")
		}
	}
	if req.BirthDate != "" {
		if parsed, err := time.Parse("2006-01-02", req.BirthDate); err != nil {
			ve.Add(utils.FieldBirthDate, "Ngày sinh phải đúng định dạng YYYY-MM-DD")
		} else if parsed.After(time.Now()) {
			ve.Add(utils.FieldBirthDate, "Ngày sinh không được là ngày trong tương lai")
		}
	}
	if req.Gender != "" && req.Gender != "male" && req.Gender != "female" && req.Gender != "other" {
		ve.Add(utils.FieldGender, "Giới tính không hợp lệ, chỉ nhận 'male', 'female' hoặc 'other'")
	}
	if ve.HasErrors() {
		return nil, errors.New(ve.Error())
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

	// Parse ngày sinh
	var birthDate *time.Time
	if req.BirthDate != "" {
		parsed, err := time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			return nil, errors.New("Ngày sinh không đúng định dạng, sử dụng YYYY-MM-DD")
		}
		birthDate = &parsed
	}

	if req.Gender == "" {
		req.Gender = "male"
	}

	// Tạo đối tượng employee lưu vào db
	emp := &model.Employee{
		DepartmentID: req.DepartmentID,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Phone:        req.Phone,
		Position:     req.Position,
		Salary:       req.Salary,
		JoinDate:     joinDate,
		BirthDate:    birthDate,
		Gender:       req.Gender,
		Status:       "active",
	}

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

	// Chuẩn hoá dữ liệu cho các field được truyền vào
	if req.FirstName != nil {
		tmp := strings.TrimSpace(*req.FirstName)
		req.FirstName = &tmp
	}
	if req.LastName != nil {
		tmp := strings.TrimSpace(*req.LastName)
		req.LastName = &tmp
	}
	if req.Phone != nil {
		tmp := strings.TrimSpace(*req.Phone)
		req.Phone = &tmp
	}
	if req.Position != nil {
		tmp := strings.TrimSpace(*req.Position)
		req.Position = &tmp
	}
	if req.Status != nil {
		tmp := strings.TrimSpace(strings.ToLower(*req.Status))
		req.Status = &tmp
	}
	if req.BirthDate != nil {
		tmp := strings.TrimSpace(*req.BirthDate)
		req.BirthDate = &tmp
	}
	if req.Gender != nil {
		tmp := strings.TrimSpace(strings.ToLower(*req.Gender))
		req.Gender = &tmp
	}

	// Validate đầu vào (defense-in-depth)
	ve := &utils.ValidationErrors{}
	if req.FirstName != nil {
		utils.CheckNameOptional(ve, utils.FieldFirstName, "Họ", *req.FirstName, 2, 100)
	}
	if req.LastName != nil {
		utils.CheckNameOptional(ve, utils.FieldLastName, "Tên", *req.LastName, 2, 100)
	}
	if req.Phone != nil {
		utils.CheckPhoneOptional(ve, *req.Phone)
	}
	if req.Position != nil && *req.Position != "" {
		l := len([]rune(*req.Position))
		if l < 2 || l > 100 {
			ve.Add(utils.FieldPosition, "Vị trí phải từ 2 đến 100 ký tự")
		}
	}
	if req.Salary != nil && *req.Salary < 0 {
		ve.Add(utils.FieldSalary, "Mức lương không được nhỏ hơn 0")
	}
	if req.Status != nil && *req.Status != "" && *req.Status != "active" && *req.Status != "inactive" {
		ve.Add(utils.FieldStatus, "Trạng thái chỉ có thể là 'active' hoặc 'inactive'")
	}
	if req.BirthDate != nil && *req.BirthDate != "" {
		if parsed, err := time.Parse("2006-01-02", *req.BirthDate); err != nil {
			ve.Add(utils.FieldBirthDate, "Ngày sinh phải đúng định dạng YYYY-MM-DD")
		} else if parsed.After(time.Now()) {
			ve.Add(utils.FieldBirthDate, "Ngày sinh không được là ngày trong tương lai")
		}
	}
	if req.Gender != nil && *req.Gender != "" && *req.Gender != "male" && *req.Gender != "female" && *req.Gender != "other" {
		ve.Add(utils.FieldGender, "Giới tính không hợp lệ, chỉ nhận 'male', 'female' hoặc 'other'")
	}
	if ve.HasErrors() {
		return nil, errors.New(ve.Error())
	}

	emp, err := es.empRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Không tìm thấy nhân viên")
		}
		return nil, fmt.Errorf("Lỗi khi tìm nhân viên: %w", err)
	}

	updateData := make(map[string]interface{})

	if req.DepartmentID != nil {
		if _, err := es.deptRepo.FindByID(*req.DepartmentID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("Không tìm thấy phòng ban mới")
			}
			return nil, fmt.Errorf("Lỗi kiểm tra phòng ban: %w", err)
		}
		updateData["department_id"] = *req.DepartmentID
	}
	if req.FirstName != nil {
		updateData["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updateData["last_name"] = *req.LastName
	}
	if req.Phone != nil {
		updateData["phone"] = *req.Phone
	}
	if req.Position != nil {
		updateData["position"] = *req.Position
	}
	if req.Salary != nil {
		if *req.Salary < 0 {
			return nil, errors.New("Mức lương không được nhỏ hơn 0")
		}
		updateData["salary"] = *req.Salary
	}
	if req.Status != nil {
		if *req.Status != "active" && *req.Status != "inactive" {
			return nil, errors.New("Trạng thái chỉ có thể là 'active' hoặc 'inactive'")
		}
		updateData["status"] = *req.Status
	}
	if req.BirthDate != nil {
		if *req.BirthDate == "" {
			updateData["birth_date"] = nil
		} else {
			parsed, err := time.Parse("2006-01-02", *req.BirthDate)
			if err != nil {
				return nil, errors.New("Ngày sinh không đúng định dạng, sử dụng YYYY-MM-DD")
			}
			updateData["birth_date"] = parsed
		}
	}
	if req.Gender != nil {
		if *req.Gender == "" {
			updateData["gender"] = "male"
		} else {
			updateData["gender"] = *req.Gender
		}
	}

	if len(updateData) == 0 {
		return emp, nil
	}

	if err := es.empRepo.UpdateFields(id, updateData); err != nil {
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
