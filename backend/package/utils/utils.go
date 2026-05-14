package utils

import (
	"chiquoc_hocgolang/internal/model"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Tạo JWT token mới cho user
func GenerateToken(userID uint, email string, roleId uint, roleName string, secretKey string, expireHour int) (string, error) {
	claims := &model.Claims{
		UserID:   userID,
		Email:    email,
		RoleID:   roleId,
		RoleName: roleName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHour) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	// Ký token bằng HS256 (HMAC SHA-256)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// ValidateToken kiểm tra và parse JWT token
func ValidateToken(tokenString string, secretKey string) (*model.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra đúng thuật toán đã ký (HMAC SHA-256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Thuật toán ký không hợp lệ!")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	claim, ok := token.Claims.(*model.Claims)
	if !ok || !token.Valid {
		return nil, errors.New("Token không hợp lệ")
	}
	return claim, nil
}
