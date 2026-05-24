package service

import (
	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/internal/repository"

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
	db       *gorm.DB
	empRepo  repository.EmployeeRepository
	deptRepo repository.DepartmentRepository
	userRepo repository.UserRepository
}

func NewEmployeeService(db *gorm.DB, empRepo repository.EmployeeRepository, deptRepo repository.DepartmentRepository, userRepo repository.UserRepository) EmployeeService {
	return &employeeService{
		db:       db,
		empRepo:  empRepo,
		deptRepo: deptRepo,
		userRepo: userRepo,
	}
}

func (es *employeeService) Create(req model.CreateEmployeeRequest) (*model.Employee, error) {
	// Chuẩn hoá dữ liệu
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Position = strings.TrimSpace(req.Position)
	req.JoinDate = strings.TrimSpace(req.JoinDate)
	req.BirthDate = strings.TrimSpace(req.BirthDate)
	req.Gender = strings.TrimSpace(strings.ToLower(req.Gender))

	// Parse ngày vào làm, mặc định là vừa tạo
	joinDate := time.Now()
	if req.JoinDate != "" {
		parsed, err := time.Parse("2006-01-02", req.JoinDate)
		if err != nil {
			return nil, errors.New("Ngày vào làm không đúng định dạng, sử dụng YYYY-MM-DD")
		}
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
	if req.Status == "" {
		req.Status = "active"
	}

	var userID *uint
	if req.UserID != nil && *req.UserID > 0 {
		userID = req.UserID
	}

	// Tạo đối tượng employee lưu vào db
	emp := &model.Employee{
		DepartmentID: req.DepartmentID,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
		Position:     req.Position,
		Salary:       req.Salary,
		JoinDate:     joinDate,
		BirthDate:    birthDate,
		Gender:       req.Gender,
		Status:       req.Status,
		UserID:       userID,
	}

	if err := es.db.Transaction(func(tx *gorm.DB) error {
		txUserRepo := es.userRepo.WithTx(tx)
		txEmpRepo := es.empRepo.WithTx(tx)
		txDeptRepo := es.deptRepo.WithTx(tx)

		if userID != nil {
			user, err := txUserRepo.FindByID(*userID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("Không tìm thấy tài khoản người dùng này")
				}
				return fmt.Errorf("Lỗi kiểm tra user: %w", err)
			}

			existingEmp, err := txEmpRepo.FindByUserID(user.ID)
			if err == nil && existingEmp != nil {
				return errors.New("Tài khoản này đã được gắn cho một nhân viên khác")
			} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("Lỗi kiểm tra user: %w", err)
			}
		}

		if _, err := txDeptRepo.FindByID(req.DepartmentID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("Không tìm thấy phòng ban này")
			}
			return fmt.Errorf("Lỗi kiểm tra phòng ban: %w", err)
		}

		if err := txEmpRepo.Create(emp); err != nil {
			return fmt.Errorf("Tạo nhân viên không thành công: %w", err)
		}

		if req.IsManager {
			if err := es.setDepartmentManager(txEmpRepo, txDeptRepo, emp.DepartmentID, &emp.ID); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
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

	emp, err := es.empRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Không tìm thấy nhân viên này")
		}
		return nil, fmt.Errorf("Lỗi khi tìm nhân viên: %w", err)
	}

	oldDeptID := emp.DepartmentID
	newDeptID := emp.DepartmentID
	if req.DepartmentID != nil {
		newDeptID = *req.DepartmentID
	}

	deptChanged := oldDeptID != newDeptID
	isOldManager := emp.Department.ManagerID != nil && *emp.Department.ManagerID == emp.ID

	if err := es.db.Transaction(func(tx *gorm.DB) error {
		txEmpRepo := es.empRepo.WithTx(tx)
		txDeptRepo := es.deptRepo.WithTx(tx)
		txUserRepo := es.userRepo.WithTx(tx)

		updateData := map[string]interface{}{}

		if req.FirstName != nil {
			updateData["first_name"] = strings.TrimSpace(*req.FirstName)
		}
		if req.LastName != nil {
			updateData["last_name"] = strings.TrimSpace(*req.LastName)
		}
		if req.Phone != nil {
			updateData["phone"] = strings.TrimSpace(*req.Phone)
		}
		if req.Position != nil {
			updateData["position"] = strings.TrimSpace(*req.Position)
		}
		if req.Salary != nil {
			if *req.Salary < 0 {
				return errors.New("Lương không được âm")
			}
			updateData["salary"] = *req.Salary
		}
		if req.Status != nil {
			updateData["status"] = strings.TrimSpace(*req.Status)
		}
		if req.Gender != nil {
			updateData["gender"] = strings.TrimSpace(strings.ToLower(*req.Gender))
		}
		if req.BirthDate != nil {
			parsed, err := time.Parse("2006-01-02", *req.BirthDate)
			if err != nil {
				return errors.New("Ngày sinh không đúng định dạng, sử dụng YYYY-MM-DD")
			}
			updateData["birth_date"] = parsed
		}

		if req.UserID != nil {
			if *req.UserID == 0 {
				updateData["user_id"] = nil
			} else {
				user, err := txUserRepo.FindByID(*req.UserID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return errors.New("Không tìm thấy tài khoản người dùng này")
					}
					return err
				}

				existingEmp, err := txEmpRepo.FindByUserID(user.ID)
				if err == nil && existingEmp.ID != id {
					return errors.New("User đã gắn cho nhân viên khác")
				}
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}

				updateData["user_id"] = *req.UserID
			}
		}

		if req.DepartmentID != nil {
			if _, err := txDeptRepo.FindByID(*req.DepartmentID); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("Không tìm thấy phòng ban này")
				}
				return fmt.Errorf("Lỗi kiểm tra phòng ban: %w", err)
			}
			updateData["department_id"] = *req.DepartmentID
		}

		if len(updateData) > 0 {
			if err := txEmpRepo.UpdateFields(id, updateData); err != nil {
				return fmt.Errorf("Cập nhật nhân viên thất bại: %w", err)
			}
		}

		updatedEmp, err := txEmpRepo.FindByID(id)
		if err != nil {
			return fmt.Errorf("Reload nhân viên thất bại: %w", err)
		}

		if req.IsManager != nil {
			if *req.IsManager {
				existingDept, err := txDeptRepo.FindByManagerID(updatedEmp.ID)
				if err == nil && existingDept.ID != newDeptID {
					return errors.New("Nhân viên đã là trưởng phòng của phòng ban khác")
				}
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("Lỗi kiểm tra manager: %w", err)
				}

				if deptChanged && isOldManager {
					if err := es.setDepartmentManager(txEmpRepo, txDeptRepo, oldDeptID, nil); err != nil {
						return err
					}
				}

				if err := es.setDepartmentManager(txEmpRepo, txDeptRepo, newDeptID, &updatedEmp.ID); err != nil {
					return err
				}
			} else {
				if isOldManager {
					if err := es.setDepartmentManager(txEmpRepo, txDeptRepo, oldDeptID, nil); err != nil {
						return err
					}
				}
			}
		} else {
			if deptChanged && isOldManager {
				if err := es.setDepartmentManager(txEmpRepo, txDeptRepo, oldDeptID, nil); err != nil {
					return err
				}
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return es.empRepo.FindByID(id)
}

func (es *employeeService) setDepartmentManager(empRepo repository.EmployeeRepository, deptRepo repository.DepartmentRepository, departmentID uint, managerID *uint) error {
	if departmentID == 0 {
		return errors.New("ID phòng ban không hợp lệ khi gán trưởng phòng")
	}

	if _, err := deptRepo.FindByID(departmentID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Không tìm thấy phòng ban để gán trưởng phòng")
		}
		return fmt.Errorf("Lỗi khi tìm phòng ban: %w", err)
	}

	if managerID != nil {
		emp, err := empRepo.FindByID(*managerID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("Không tìm thấy nhân viên để gán làm trưởng phòng")
			}
			return fmt.Errorf("Lỗi khi tìm nhân viên quản lý: %w", err)
		}
		if emp.DepartmentID != departmentID {
			return errors.New("Trưởng phòng phải thuộc chính phòng ban này")
		}

		existingDept, err := deptRepo.FindByManagerID(*managerID)
		if err == nil && existingDept.ID != departmentID {
			return fmt.Errorf("Nhân viên này đã là trưởng phòng của phòng khác")
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	return deptRepo.UpdateManager(departmentID, managerID)
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
