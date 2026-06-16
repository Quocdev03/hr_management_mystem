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

// AuthService interface - contract cho authentication.
type AuthService interface {
	Login(req model.LoginRequest) (*model.LoginResponse, error)
	GetProfile(ID uint) (*model.Employee, error)
	GetPermissions(userID uint) ([]string, error)
	Logout(tokenString string, remainingTime time.Duration, refreshToken string) error
	RefreshToken(refreshTokenString string) (*model.LoginResponse, error)
}

// --- Auth Service Implementation ---
type authService struct {
	useRepo  repository.UserRepository
	empRepo  repository.EmployeeRepository
	permRepo repository.PermissionRepository
	jwtCfg   *config.JWTConfig
	rdb      *redis.Client
}

func NewAuthService(userRepo repository.UserRepository, empRepo repository.EmployeeRepository, permRepo repository.PermissionRepository, jwtCfg *config.JWTConfig, rdb *redis.Client) AuthService {
	return &authService{
		useRepo:  userRepo,
		empRepo:  empRepo,
		permRepo: permRepo,
		jwtCfg:   jwtCfg,
		rdb:      rdb,
	}
}

// Logout xử lý đăng xuất bằng cách blacklist access token và xóa refresh token.
func (au *authService) Logout(tokenString string, remainingTime time.Duration, refreshToken string) error {
	if au.rdb == nil {
		return errors.New("redis client is not initialized")
	}
	ctx := context.Background()

	// 1. Blacklist access token.
	if err := au.rdb.Set(ctx, "blacklist:"+tokenString, "true", remainingTime).Err(); err != nil {
		utils.Error("failed to blacklist token: %v", err)
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	// 2. Xóa refresh token CHỈ KHI refresh token hợp lệ (đã được verify).
	// Không dùng ParseUnverified fallback vì có thể bị exploit với unverified claims.
	if refreshToken != "" {
		claims, err := utils.ValidateToken(refreshToken, au.jwtCfg.SecretKey)
		if err == nil && claims.TokenType == "refresh" {
			key := fmt.Sprintf("refresh_token:%d", claims.UserID)
			if delErr := au.rdb.Del(ctx, key).Err(); delErr != nil {
				utils.Warn("failed to delete refresh token for user %d: %v", claims.UserID, delErr)
				// Không thất bại logout chỉ vì cleanup refresh token lỗi
			}
		}
	}

	return nil
}

// Login và trả về JWT Token và Refresh Token.
func (au *authService) Login(req model.LoginRequest) (*model.LoginResponse, error) {
	// Validate đầu vào (defense-in-depth, handler đã validate rồi).
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" {
		return nil, errors.New("email là bắt buộc")
	}
	if len(req.Password) < 8 {
		return nil, errors.New("mật khẩu phải có ít nhất 8 ký tự")
	}

	// Tìm user theo email.
	user, err := au.useRepo.FindByEmail(req.Email)
	if err != nil {
		// Trả lỗi chung không tiết lộ email hay pass.
		return nil, errors.New("email hoặc mật khẩu không hợp lệ")
	}

	// Kiểm tra xem tài khoản có bị khoá không.
	if !user.IsActive {
		return nil, errors.New("tài khoản đã bị vô hiệu hoá")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("email hoặc mật khẩu không hợp lệ")
	}

	permissions, err := au.permRepo.GetPermissionCodes(user.ID)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy quyền người dùng: %w", err)

	}
	user.Permissions = permissions

	// Tạo jwt token chứa thông tin user và role.
	token, err := utils.GenerateToken(
		user.ID,
		user.Email,
		user.RoleID,
		user.Role.Name,
		au.jwtCfg.SecretKey,
		au.jwtCfg.ExpireHour,
	)
	if err != nil {
		return nil, fmt.Errorf("có lỗi khi tạo token: %w", err)
	}

	// Tạo refresh token chứa thông tin user và role.
	refreshToken, err := utils.GenerateRefreshToken(
		user.ID,
		user.Email,
		user.RoleID,
		user.Role.Name,
		au.jwtCfg.SecretKey,
		au.jwtCfg.RefreshExpireDay,
	)
	if err != nil {
		return nil, fmt.Errorf("có lỗi khi tạo refresh token: %w", err)
	}

	// Lưu refresh token vào Redis.
	if au.rdb != nil {
		ctx := context.Background()
		redisKey := fmt.Sprintf("refresh_token:%d", user.ID)
		ttl := time.Duration(au.jwtCfg.RefreshExpireDay) * 24 * time.Hour
		if err := au.rdb.Set(ctx, redisKey, refreshToken, ttl).Err(); err != nil {
			return nil, fmt.Errorf("không thể lưu refresh token vào cache: %w", err)
		}
	}

	return &model.LoginResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

// RefreshToken dùng refresh token cũ để lấy access token mới và refresh token mới (rotation).
func (au *authService) RefreshToken(refreshTokenString string) (*model.LoginResponse, error) {
	if au.rdb == nil {
		return nil, errors.New("redis client is not initialized")
	}
	ctx := context.Background()

	// 1. Validate refresh token.
	claims, err := utils.ValidateToken(refreshTokenString, au.jwtCfg.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("refresh token không hợp lệ: %w", err)
	}

	// 2. Phải đúng loại token refresh.
	if claims.TokenType != "refresh" {
		return nil, errors.New("token không phải là refresh token")
	}

	// 3. Lấy refresh token trong Redis.
	redisKey := fmt.Sprintf("refresh_token:%d", claims.UserID)
	storedToken, err := au.rdb.Get(ctx, redisKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errors.New("refresh token đã hết hạn hoặc không tồn tại")
		}
		return nil, fmt.Errorf("lỗi hệ thống khi kiểm tra token: %w", err)
	}

	// 4. So khớp token gửi lên với token trong Redis.
	if storedToken != refreshTokenString {
		return nil, errors.New("refresh token không khớp hoặc đã bị vô hiệu hóa")
	}

	// 5. Tìm user để lấy thông tin mới nhất.
	user, err := au.useRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New("người dùng không tồn tại hoặc đã bị xóa")
	}
	if !user.IsActive {
		return nil, errors.New("tài khoản người dùng đã bị khóa")
	}

	permissions, err := au.permRepo.GetPermissionCodes(user.ID)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy quyền người dùng: %w", err)
	}
	user.Permissions = permissions

	// 6. Tạo access token mới.
	newAccessToken, err := utils.GenerateToken(
		user.ID,
		user.Email,
		user.RoleID,
		user.Role.Name,
		au.jwtCfg.SecretKey,
		au.jwtCfg.ExpireHour,
	)
	if err != nil {
		return nil, fmt.Errorf("không thể tạo access token mới: %w", err)
	}

	// 7. Tạo refresh token mới (Rotation để tăng bảo mật).
	newRefreshToken, err := utils.GenerateRefreshToken(
		user.ID,
		user.Email,
		user.RoleID,
		user.Role.Name,
		au.jwtCfg.SecretKey,
		au.jwtCfg.RefreshExpireDay,
	)
	if err != nil {
		return nil, fmt.Errorf("không thể tạo refresh token mới: %w", err)
	}

	// 8. Lưu refresh token mới vào Redis.
	ttl := time.Duration(au.jwtCfg.RefreshExpireDay) * 24 * time.Hour
	if err := au.rdb.Set(ctx, redisKey, newRefreshToken, ttl).Err(); err != nil {
		return nil, fmt.Errorf("không thể lưu refresh token mới: %w", err)
	}

	return &model.LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		User:         *user,
	}, nil
}

func (au *authService) GetPermissions(userID uint) ([]string, error) {
	return au.permRepo.GetPermissionCodes(userID)
}

func (au *authService) GetProfile(id uint) (*model.Employee, error) {

	emp, err := au.empRepo.FindByUserID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("không tìm thấy hồ sơ nhân viên")
		}
		return nil, errors.New("không thể lấy thông tin hồ sơ nhân viên")
	}

	return emp, nil
}
