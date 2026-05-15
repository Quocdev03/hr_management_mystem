package model

import "github.com/golang-jwt/jwt/v5"

// RegisterRequest - dữ liệu đăng ký tài khoản
type RegisterRequest struct {
	UserName string `json:"user_name" binding:"required,min=4,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest - dữ liệu đăng nhập
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginResponse  - dữ liệu trả về sau khi đăng nhập thành công
type LoginResponse struct {
	AccessToken string `json:"access_token"`
	User        User   `json:"user"`
}

// Claims chứa thông tin được mã hóa trong JWT token
type Claims struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	RoleID   uint   `json:"role_id"`
	RoleName string `json:"role_name"`
	jwt.RegisteredClaims
}
