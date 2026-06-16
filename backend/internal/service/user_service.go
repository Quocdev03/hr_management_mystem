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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService interface - contract cho quản lý user
type UserService interface {
	Create(req model.CreateUserRequest) (*model.User, error)
	GetUsers(query model.PaginationQuery) (*model.PaginatedResult, error)
	GetUserByID(id uint) (*model.User, error)
	UpdateUser(id uint, req model.UpdateUserRequest, requesterID uint) (*model.User, error)
	DeleteUser(id uint) error
	GetUsersWithoutEmployee() ([]model.User, error)
	GetAvailablePermissions() ([]model.Permission, error)
	UpdateUserPermissions(id uint, permissionCodes []string) ([]string, error)
}

// --- User Service Implementation ---

type userService struct {
	userRepo   repository.UserRepository
	permRepo   repository.PermissionRepository
	rdb        *redis.Client
}

func NewUserService(userRepo repository.UserRepository, permRepo repository.PermissionRepository, rdb *redis.Client) UserService {
	return &userService{
		userRepo: userRepo,
		permRepo: permRepo,
		rdb:      rdb,
	}
}

func (us *userService) Create(req model.CreateUserRequest) (*model.User, error) {
	req.UserName = strings.TrimSpace(req.UserName)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)

	if req.RoleID == 0 {
		return nil, errors.New("roleID là bắt buộc")
	}

	if existing, err := us.userRepo.FindByUsername(req.UserName); err == nil && existing != nil {
		return nil, errors.New("tên đăng nhập đã tồn tại")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("lỗi kiểm tra tên đăng nhập: %w", err)
	}

	if existing, err := us.userRepo.FindByEmail(req.Email); err == nil && existing != nil {
		return nil, errors.New("email đã tồn tại")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("lỗi kiểm tra email: %w", err)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("lỗi mã hoá mật khẩu: %w", err)
	}

	user := &model.User{
		UserName: req.UserName,
		Email:    req.Email,
		Password: string(hashed),
		RoleID:   req.RoleID,
		IsActive: req.IsActive,
	}

	if err := us.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("tạo user không thành công: %w", err)
	}

	// Invalidate dashboard stats cache
	_ = utils.InvalidateDashboardStats(context.Background(), us.rdb)

	return us.userRepo.FindByID(user.ID)
}

func (us *userService) GetUsers(query model.PaginationQuery) (*model.PaginatedResult, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 10
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	users, total, err := us.userRepo.FindAll(query)
	if err != nil {
		return nil, fmt.Errorf("lấy danh sách user bị lỗi: %w", err)
	}

	totalPage := int(math.Ceil(float64(total) / float64(query.Limit)))

	return &model.PaginatedResult{
		Items:      users,
		Total:      total,
		Page:       query.Page,
		Limit:      query.Limit,
		TotalPages: totalPage,
	}, nil
}

func (us *userService) GetUserByID(id uint) (*model.User, error) {
	if id == 0 {
		return nil, errors.New("id user phải lớn hơn 0")
	}

	user, err := us.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("không tìm thấy user")
		}
		return nil, err
	}

	return user, nil
}

func (us *userService) UpdateUser(id uint, req model.UpdateUserRequest, requesterID uint) (*model.User, error) {
	if id == 0 {
		return nil, errors.New("id user phải lớn hơn 0")
	}

	if req.UserName != nil {
		tmp := strings.TrimSpace(*req.UserName)
		req.UserName = &tmp
	}
	if req.Email != nil {
		tmp := strings.TrimSpace(strings.ToLower(*req.Email))
		req.Email = &tmp
	}
	if req.Password != nil {
		tmp := strings.TrimSpace(*req.Password)
		req.Password = &tmp
	}

	user, err := us.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("không tìm thấy user")
		}
		return nil, err
	}

	// Logic kiểm tra phân quyền an toàn.
	// 1. Ngăn tự đổi quyền bản thân HOẶC đổi quyền của Admin khác.
	if req.RoleID != nil && *req.RoleID != user.RoleID {
		if user.RoleID == 1 || user.ID == requesterID {
			return nil, errors.New("không thể tự thay đổi quyền của mình hoặc của Admin khác")
		}
	}

	// 2. Ngăn tự khoá tài khoản của chính mình.
	if user.ID == requesterID && req.IsActive != nil && !*req.IsActive {
		return nil, errors.New("không thể tự vô hiệu hoá tài khoản của chính mình")
	}

	updateData := make(map[string]interface{})

	if req.UserName != nil && *req.UserName != user.UserName {
		if existing, err := us.userRepo.FindByUsername(*req.UserName); err == nil && existing != nil && existing.ID != user.ID {
			return nil, errors.New("tên đăng nhập đã tồn tại")
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("lỗi kiểm tra tên đăng nhập: %w", err)
		}
		updateData["user_name"] = *req.UserName
	}
	if req.Email != nil && *req.Email != user.Email {
		if existing, err := us.userRepo.FindByEmail(*req.Email); err == nil && existing != nil && existing.ID != user.ID {
			return nil, errors.New("email đã tồn tại")
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("lỗi kiểm tra email: %w", err)
		}
		updateData["email"] = *req.Email
	}
	if req.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("lỗi mã hoá mật khẩu: %w", err)
		}
		updateData["password"] = string(hashed)
	}
	if req.RoleID != nil && *req.RoleID != user.RoleID {
		updateData["role_id"] = *req.RoleID
	}
	if req.IsActive != nil && *req.IsActive != user.IsActive {
		updateData["is_active"] = *req.IsActive
	}

	if len(updateData) == 0 {
		return user, nil
	}

	if err := us.userRepo.UpdateFields(id, updateData); err != nil {
		return nil, fmt.Errorf("cập nhật user bị lỗi: %w", err)
	}

	// Invalidate dashboard stats cache
	_ = utils.InvalidateDashboardStats(context.Background(), us.rdb)

	return us.userRepo.FindByID(id)
}

func (us *userService) DeleteUser(id uint) error {
	if id == 0 {
		return errors.New("id user phải lớn hơn 0")
	}

	if _, err := us.userRepo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("không tìm thấy user")
		}
		return err
	}

	err := us.userRepo.Delete(id)
	if err == nil {
		// Invalidate dashboard stats cache
		_ = utils.InvalidateDashboardStats(context.Background(), us.rdb)
	}
	return err
}

func (us *userService) GetUsersWithoutEmployee() ([]model.User, error) {
	users, err := us.userRepo.FindUsersWithoutEmployee()
	if err != nil {
		return nil, fmt.Errorf("lấy danh sách user chưa gắn employee bị lỗi: %w", err)
	}

	return users, nil
}

func (us *userService) GetAvailablePermissions() ([]model.Permission, error) {
	return us.permRepo.GetAllPermissions()
}

func (us *userService) UpdateUserPermissions(id uint, permissionCodes []string) ([]string, error) {
	if id == 0 {
		return nil, errors.New("id user phải lớn hơn 0")
	}

	if _, err := us.userRepo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("không tìm thấy user")
		}
		return nil, err
	}

	if err := us.permRepo.SetUserPermissions(id, permissionCodes); err != nil {
		return nil, fmt.Errorf("cập nhật quyền user bị lỗi: %w", err)
	}

	codes, err := us.permRepo.GetPermissionCodes(id)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy quyền sau khi cập nhật: %w", err)
	}
	return codes, nil
}
