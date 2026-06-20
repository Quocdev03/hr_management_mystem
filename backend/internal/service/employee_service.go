package service

import (
	"chiquoc_hocgolang/internal/common"
	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/internal/repository"
	"chiquoc_hocgolang/internal/utils"
	"context"

	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
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
	rdb      *redis.Client
}

func NewEmployeeService(db *gorm.DB, empRepo repository.EmployeeRepository, deptRepo repository.DepartmentRepository, userRepo repository.UserRepository, rdb *redis.Client) EmployeeService {
	return &employeeService{
		db:       db,
		empRepo:  empRepo,
		deptRepo: deptRepo,
		userRepo: userRepo,
		rdb:      rdb,
	}
}

func (es *employeeService) Create(req model.CreateEmployeeRequest) (*model.Employee, error) {
	// Chuẩn hoá dữ liệu
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Phone = strings.TrimSpace(req.Phone)
	req.JoinDate = strings.TrimSpace(req.JoinDate)
	req.BirthDate = strings.TrimSpace(req.BirthDate)
	req.Gender = strings.TrimSpace(strings.ToLower(req.Gender))

	// Parse ngày vào làm, mặc định là vừa tạo
	joinDate := time.Now()
	if req.JoinDate != "" {
		parsed, err := time.Parse("2006-01-02", req.JoinDate)
		if err != nil {
			return nil, errors.New("ngày vào làm không đúng định dạng, sử dụng YYYY-MM-DD")
		}
		if parsed.After(time.Now()) {
			return nil, errors.New("ngày vào làm không được là ngày trong tương lai")
		}
		joinDate = parsed
	}

	// Parse ngày sinh
	var birthDate *time.Time
	if req.BirthDate != "" {
		parsed, err := time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			return nil, errors.New("ngày sinh không đúng định dạng, sử dụng YYYY-MM-DD")
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
		PositionID:   req.PositionID,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
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
					return errors.New("không tìm thấy tài khoản người dùng này")
				}
				return fmt.Errorf("lỗi kiểm tra user: %w", err)
			}

			if !user.IsActive {
				return errors.New("không thể gắn tài khoản đang ngưng hoạt động")
			}

			existingEmp, err := txEmpRepo.FindByUserID(user.ID)
			if err == nil && existingEmp != nil {
				return errors.New("tài khoản này đã được gắn cho một nhân viên khác")
			} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("lỗi kiểm tra user: %w", err)
			}
		}

		if _, err := txDeptRepo.FindByID(req.DepartmentID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("không tìm thấy phòng ban này")
			}
			return fmt.Errorf("lỗi kiểm tra phòng ban: %w", err)
		}

		var pos model.Position
		if err := tx.Where("id = ?", req.PositionID).First(&pos).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("chức vụ không tồn tại")
			}
			return fmt.Errorf("lỗi kiểm tra chức vụ: %w", err)
		}

		if err := txEmpRepo.Create(emp); err != nil {
			return fmt.Errorf("tạo nhân viên không thành công: %w", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// Invalidate dashboard stats cache
	if err := utils.InvalidateDashboardStats(context.Background(), es.rdb); err != nil {
		utils.Error("không thể invalidate cache dashboard: %v", err)
	}

	return es.empRepo.FindByID(emp.ID)
}

// Danh sách nhân viên có phân trang và tìm kiếm
func (es *employeeService) GetEmployees(query model.PaginationQuery) (*model.PaginatedResult, error) {
	common.NormalizePagination(&query)

	employees, total, err := es.empRepo.FindAll(query)
	if err != nil {
		return nil, fmt.Errorf("lấy danh sách nhân viên bị lỗi: %w", err)
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
		return nil, errors.New("id nhân viên phải lớn hơn 0")
	}

	emp, err := es.empRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("không tìm thấy nhân viên này")
		}
		return nil, err
	}

	return emp, nil
}

func (es *employeeService) UpdateEmployee(id uint, req model.UpdateEmployeeRequest) (*model.Employee, error) {
	if id == 0 {
		return nil, errors.New("id nhân viên phải lớn hơn 0")
	}

	var result *model.Employee

	if err := es.db.Transaction(func(tx *gorm.DB) error {
		txEmpRepo := es.empRepo.WithTx(tx)
		txDeptRepo := es.deptRepo.WithTx(tx)
		txUserRepo := es.userRepo.WithTx(tx)

		// Lấy lại dữ liệu nhân viên trong transaction để tính toán trạng thái chính xác nhất
		empTx, err := txEmpRepo.FindByID(id)
		if err != nil {
			return fmt.Errorf("lỗi khi tìm nhân viên trong transaction: %w", err)
		}

		oldDeptID := empTx.DepartmentID
		newDeptID := empTx.DepartmentID
		if req.DepartmentID != nil {
			newDeptID = *req.DepartmentID
		}

		deptChanged := oldDeptID != newDeptID
		isOldManager := empTx.Department != nil && empTx.Department.ManagerID != nil && *empTx.Department.ManagerID == empTx.ID

		updateData := map[string]interface{}{}

		if req.FirstName != nil {
			value := strings.TrimSpace(*req.FirstName)
			if value != empTx.FirstName {
				updateData["first_name"] = value
			}
		}
		if req.LastName != nil {
			value := strings.TrimSpace(*req.LastName)
			if value != empTx.LastName {
				updateData["last_name"] = value
			}
		}
		if req.Phone != nil {
			value := strings.TrimSpace(*req.Phone)
			if value != empTx.Phone {
				updateData["phone"] = value
			}
		}
		if req.PositionID != nil {
			newPosID := *req.PositionID
			if newPosID != empTx.PositionID {
				var pos model.Position
				if err := tx.Where("id = ?", newPosID).First(&pos).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return errors.New("chức vụ không tồn tại")
					}
					return fmt.Errorf("lỗi kiểm tra chức vụ: %w", err)
				}
				updateData["position_id"] = newPosID
			}
		}
		if req.Salary != nil {
			if *req.Salary < 0 {
				return errors.New("lương không được âm")
			}
			if *req.Salary != empTx.Salary {
				updateData["salary"] = *req.Salary
			}
		}
		statusChangedToInactive := false
		if req.Status != nil {
			value := strings.TrimSpace(*req.Status)
			if value != empTx.Status {
				updateData["status"] = value
				if value == "inactive" {
					statusChangedToInactive = true
				}
			}
		}
		if req.Gender != nil {
			value := strings.TrimSpace(strings.ToLower(*req.Gender))
			if value != empTx.Gender {
				updateData["gender"] = value
			}
		}
		if req.BirthDate != nil {
			parsed, err := time.Parse("2006-01-02", *req.BirthDate)
			if err != nil {
				return errors.New("ngày sinh không đúng định dạng, sử dụng YYYY-MM-DD")
			}
			if empTx.BirthDate == nil || !empTx.BirthDate.Equal(parsed) {
				updateData["birth_date"] = parsed
			}
		}

		if req.UserID != nil {
			oldUserID := empTx.UserID
			newUserID := *req.UserID

			// Helper closure để validate user
			validateUser := func(uID uint) error {
				if uID == 0 {
					return nil
				}
				user, err := txUserRepo.FindByID(uID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return errors.New("không tìm thấy tài khoản người dùng này")
					}
					return err
				}
				if !user.IsActive {
					return errors.New("không thể gắn tài khoản đang ngưng hoạt động")
				}
				existingEmp, err := txEmpRepo.FindByUserID(user.ID)
				if err == nil && existingEmp.ID != id {
					return errors.New("user đã gắn cho nhân viên khác")
				}
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}
				return nil
			}

			if oldUserID == nil {
				if newUserID != 0 {
					if err := validateUser(newUserID); err != nil {
						return err
					}
					updateData["user_id"] = newUserID
				}
			} else if *oldUserID != newUserID {
				if newUserID == 0 {
					updateData["user_id"] = nil
				} else {
					if err := validateUser(newUserID); err != nil {
						return err
					}
					updateData["user_id"] = newUserID
				}
			}
		}

		if req.DepartmentID != nil {
			if *req.DepartmentID != empTx.DepartmentID {
				if _, err := txDeptRepo.FindByID(*req.DepartmentID); err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return errors.New("không tìm thấy phòng ban này")
					}
					return fmt.Errorf("lỗi kiểm tra phòng ban: %w", err)
				}
				updateData["department_id"] = *req.DepartmentID
			}
		}

		if len(updateData) > 0 {
			if err := txEmpRepo.UpdateFields(id, updateData); err != nil {
				return fmt.Errorf("cập nhật nhân viên thất bại: %w", err)
			}
		}

		updatedEmp, err := txEmpRepo.FindByID(id)
		if err != nil {
			return fmt.Errorf("reload nhân viên thất bại: %w", err)
		}

		if deptChanged && isOldManager {
			if err := txDeptRepo.UpdateManager(oldDeptID, nil); err != nil {
				return fmt.Errorf("lỗi khi gỡ trưởng phòng ở phòng cũ: %w", err)
			}
		} else if statusChangedToInactive && isOldManager {
			if err := txDeptRepo.UpdateManager(oldDeptID, nil); err != nil {
				return fmt.Errorf("lỗi khi gỡ trưởng phòng do nhân viên nghỉ việc: %w", err)
			}
		}

		result = updatedEmp
		return nil
	}); err != nil {
		return nil, err
	}

	// Invalidate dashboard stats cache
	if err := utils.InvalidateDashboardStats(context.Background(), es.rdb); err != nil {
		utils.Error("không thể invalidate cache dashboard: %v", err)
	}

	return result, nil
}

func (es *employeeService) DeleteEmployee(id uint) error {
	if id == 0 {
		return errors.New("id nhân viên phải lớn hơn 0")
	}

	if _, err := es.empRepo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("không tìm thấy nhân viên này")
		}
		return fmt.Errorf("lỗi khi tìm nhân viên: %w", err)
	}
	err := es.db.Transaction(func(tx *gorm.DB) error {
		txEmpRepo := es.empRepo.WithTx(tx)
		txDeptRepo := es.deptRepo.WithTx(tx)

		// Kiểm tra xem nhân viên này có đang làm trưởng phòng không
		if dept, err := txDeptRepo.FindByManagerID(id); err == nil {
			// Tự động gỡ quyền trưởng phòng
			if updateErr := txDeptRepo.UpdateManager(dept.ID, nil); updateErr != nil {
				return fmt.Errorf("lỗi khi gỡ quyền trưởng phòng trước khi xoá nhân viên: %w", updateErr)
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("lỗi khi kiểm tra quyền trưởng phòng: %w", err)
		}

		if err := txEmpRepo.Delete(id); err != nil {
			return fmt.Errorf("lỗi khi xoá nhân viên này: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}
	// Invalidate dashboard stats cache
	if err := utils.InvalidateDashboardStats(context.Background(), es.rdb); err != nil {
		utils.Error("không thể invalidate cache dashboard: %v", err)
	}

	return nil
}
