package service

import (
	"chiquoc_hocgolang/internal/common"
	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/internal/repository"
	"context"

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
	cacheSvc CacheService
}

func NewEmployeeService(db *gorm.DB, empRepo repository.EmployeeRepository, deptRepo repository.DepartmentRepository, userRepo repository.UserRepository, cacheSvc CacheService) EmployeeService {
	return &employeeService{
		db:       db,
		empRepo:  empRepo,
		deptRepo: deptRepo,
		userRepo: userRepo,
		cacheSvc: cacheSvc,
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

			if !user.IsActive {
				return errors.New("Không thể gắn tài khoản đang ngưng hoạt động")
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

		return nil
	}); err != nil {
		return nil, err
	}

	// Invalidate dashboard stats cache
	_ = es.cacheSvc.Delete(context.Background(), "dashboard:stats")

	return es.empRepo.FindByID(emp.ID)
}

// Danh sách nhân viên có phân trang và tìm kiếm
func (es *employeeService) GetEmployees(query model.PaginationQuery) (*model.PaginatedResult, error) {
	common.NormalizePagination(&query)

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

	_, err := es.empRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Không tìm thấy nhân viên này")
		}
		return nil, fmt.Errorf("Lỗi khi tìm nhân viên: %w", err)
	}

	var result *model.Employee

	if err := es.db.Transaction(func(tx *gorm.DB) error {
		txEmpRepo := es.empRepo.WithTx(tx)
		txDeptRepo := es.deptRepo.WithTx(tx)
		txUserRepo := es.userRepo.WithTx(tx)

		// Lấy lại dữ liệu nhân viên trong transaction để tính toán trạng thái chính xác nhất
		empTx, err := txEmpRepo.FindByID(id)
		if err != nil {
			return fmt.Errorf("Lỗi khi tìm nhân viên trong transaction: %w", err)
		}

		oldDeptID := empTx.DepartmentID
		newDeptID := empTx.DepartmentID
		if req.DepartmentID != nil {
			newDeptID = *req.DepartmentID
		}

		deptChanged := oldDeptID != newDeptID
		isOldManager := empTx.Department.ManagerID != nil && *empTx.Department.ManagerID == empTx.ID

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

				if !user.IsActive {
					return errors.New("Không thể gắn tài khoản đang ngưng hoạt động")
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

		if deptChanged && isOldManager {
			if err := txDeptRepo.UpdateManager(oldDeptID, nil); err != nil {
				return fmt.Errorf("Lỗi khi gỡ trưởng phòng ở phòng cũ: %w", err)
			}
		}

		result = updatedEmp
		return nil
	}); err != nil {
		return nil, err
	}

	// Invalidate dashboard stats cache
	_ = es.cacheSvc.Delete(context.Background(), "dashboard:stats")

	return result, nil
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
	err := es.db.Transaction(func(tx *gorm.DB) error {
		txEmpRepo := es.empRepo.WithTx(tx)
		txDeptRepo := es.deptRepo.WithTx(tx)

		// Kiểm tra xem nhân viên này có đang làm trưởng phòng không
		if dept, err := txDeptRepo.FindByManagerID(id); err == nil {
			// Tự động gỡ quyền trưởng phòng
			if updateErr := txDeptRepo.UpdateManager(dept.ID, nil); updateErr != nil {
				return fmt.Errorf("Lỗi khi gỡ quyền trưởng phòng trước khi xoá nhân viên: %w", updateErr)
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("Lỗi khi kiểm tra quyền trưởng phòng: %w", err)
		}

		if err := txEmpRepo.Delete(id); err != nil {
			return fmt.Errorf("Lỗi khi xoá nhân viên này: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}
	// Invalidate dashboard stats cache
	_ = es.cacheSvc.Delete(context.Background(), "dashboard:stats")

	return nil
}
