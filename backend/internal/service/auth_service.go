package service

import (
	"chiquoc_hocgolang/internal/config"
	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/internal/repository"
	"chiquoc_hocgolang/internal/utils"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService interface - contract cho authentication
type AuthService interface {
	Login(req model.LoginRequest) (*model.LoginResponse, error)
	GetProfile(ID uint) (*model.Employee, error)
	Logout(tokenString string, remainingTime time.Duration) error
}

// --- Auth Service Implementation ---
type authServive struct {
	useRepo repository.UserRepository
	empRepo repository.EmployeeRepository
	jwtCfg  *config.JWTConfig
	rdb     *redis.Client
}

func NewAuthService(userRepo repository.UserRepository, empRepo repository.EmployeeRepository, jwtCfg *config.JWTConfig, rdb *redis.Client) AuthService {
	return &authServive{
		useRepo: userRepo,
		empRepo: empRepo,
		jwtCfg:  jwtCfg,
		rdb:     rdb,
	}
}

func (au *authServive) Logout(tokenString string, remainingTime time.Duration) error {
	if au.rdb == nil {
		return errors.New("Redis client is not initialized")
	}
	ctx := context.Background()
	return au.rdb.Set(ctx, "blacklist:"+tokenString, "true", remainingTime).Err()
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

func (au *authServive) GetProfile(id uint) (*model.Employee, error) {

	emp, err := au.empRepo.FindByUserID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.New("Không thể lấy thông tin hồ sơ nhân viên")
	}

	return emp, nil
}
