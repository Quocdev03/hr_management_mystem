package service

import (
	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/internal/repository"

	"errors"
	"fmt"
	"math"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService interface - contract cho quản lý user
type UserService interface {
	Create(req model.CreateUserRequest) (*model.User, error)
	GetUsers(query model.PaginationQuery) (*model.PaginatedResult, error)
	GetUserByID(id uint) (*model.User, error)
	UpdateUser(id uint, req model.UpdateUserRequest) (*model.User, error)
	DeleteUser(id uint) error
}

// --- User Service Implementation ---

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (us *userService) Create(req model.CreateUserRequest) (*model.User, error) {
	req.UserName = strings.TrimSpace(req.UserName)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)



	if req.RoleID == 0 {
		return nil, errors.New("RoleID là bắt buộc")
	}

	if existing, err := us.userRepo.FindByUsername(req.UserName); err == nil && existing != nil {
		return nil, errors.New("Tên đăng nhập đã tồn tại")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("Lỗi kiểm tra tên đăng nhập: %w", err)
	}

	if existing, err := us.userRepo.FindByEmail(req.Email); err == nil && existing != nil {
		return nil, errors.New("Email đã tồn tại")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("Lỗi kiểm tra email: %w", err)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Lỗi mã hoá mật khẩu: %w", err)
	}

	user := &model.User{
		UserName: req.UserName,
		Email:    req.Email,
		Password: string(hashed),
		RoleID:   req.RoleID,
		IsActive: req.IsActive,
	}

	if err := us.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("Tạo user không thành công: %w", err)
	}

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
		return nil, fmt.Errorf("Lấy danh sách user bị lỗi: %w", err)
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
		return nil, errors.New("ID user phải lớn hơn 0")
	}

	user, err := us.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Không tìm thấy user")
		}
		return nil, err
	}

	return user, nil
}

func (us *userService) UpdateUser(id uint, req model.UpdateUserRequest) (*model.User, error) {
	if id == 0 {
		return nil, errors.New("ID user phải lớn hơn 0")
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
			return nil, errors.New("Không tìm thấy user")
		}
		return nil, err
	}

	updateData := make(map[string]interface{})

	if req.UserName != nil && *req.UserName != user.UserName {
		if existing, err := us.userRepo.FindByUsername(*req.UserName); err == nil && existing != nil && existing.ID != user.ID {
			return nil, errors.New("Tên đăng nhập đã tồn tại")
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Lỗi kiểm tra tên đăng nhập: %w", err)
		}
		updateData["user_name"] = *req.UserName
	}
	if req.Email != nil && *req.Email != user.Email {
		if existing, err := us.userRepo.FindByEmail(*req.Email); err == nil && existing != nil && existing.ID != user.ID {
			return nil, errors.New("Email đã tồn tại")
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Lỗi kiểm tra email: %w", err)
		}
		updateData["email"] = *req.Email
	}
	if req.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("Lỗi mã hoá mật khẩu: %w", err)
		}
		updateData["password"] = string(hashed)
	}
	if req.RoleID != nil {
		updateData["role_id"] = *req.RoleID
	}
	if req.IsActive != nil {
		updateData["is_active"] = *req.IsActive
	}

	if len(updateData) == 0 {
		return user, nil
	}

	if err := us.userRepo.UpdateFields(id, updateData); err != nil {
		return nil, fmt.Errorf("Cập nhật user bị lỗi: %w", err)
	}

	return us.userRepo.FindByID(id)
}

func (us *userService) DeleteUser(id uint) error {
	if id == 0 {
		return errors.New("ID user phải lớn hơn 0")
	}

	if _, err := us.userRepo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Không tìm thấy user")
		}
		return err
	}

	return us.userRepo.Delete(id)
}
