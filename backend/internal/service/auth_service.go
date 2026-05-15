package service

import (
	"chiquoc_hocgolang/internal/config"
	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/internal/repository"
	"chiquoc_hocgolang/internal/utils"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// AuthService interface - contract cho authentication
type AuthService interface {
	Login(req model.LoginRequest) (*model.LoginResponse, error)
}

// --- Auth Service Implementation ---
type authServive struct {
	useRepo repository.UserRepository
	jwtCfg  *config.JWTConfig
}

func NewAuthService(userRepo repository.UserRepository, jwtCfg *config.JWTConfig) AuthService {
	return &authServive{
		useRepo: userRepo,
		jwtCfg:  jwtCfg,
	}
}

// Login và trả về JWT Token
func (au *authServive) Login(req model.LoginRequest) (*model.LoginResponse, error) {
	// Validate đầu vào (defense-in-depth, handler đã validate rồi)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" {
		return nil, errors.New("Email là bắt buộc")
	}
	if len(req.Password) < 8 {
		return nil, errors.New("Mật khẩu phải có ít nhất 8 ký tự")
	}

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
