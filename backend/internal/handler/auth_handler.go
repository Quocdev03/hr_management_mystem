package handler

import (
	"chiquoc_hocgolang/internal/middleware"
	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/internal/service"
	"chiquoc_hocgolang/internal/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authSvc service.AuthService
}

func NewAuthHandler(authSvc service.AuthService) *AuthHandler {
	return &AuthHandler{
		authSvc: authSvc,
	}
}

// Login godoc
// POST /api/v1/auth/login
func (h *AuthHandler) Login(ctx *gin.Context) {
	var req model.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Dữ liệu đăng nhập không hợp lệ")
		return
	}

	// Trim whitespace
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	// Validate đầu vào
	ve := &utils.ValidationErrors{}
	utils.CheckEmail(ve, req.Email)
	utils.CheckPassword(ve, req.Password)
	if ve.HasErrors() {
		utils.ValidationError(ctx, "Dữ liệu đăng nhập không hợp lệ", ve.Errors)
		return
	}

	result, err := h.authSvc.Login(req)
	if err != nil {
		utils.Unauthorized(ctx, err.Error())
		return
	}
	utils.Success(ctx, "Đăng nhập thành công", result)
}

// GetProfile godoc
// GET /api/v1/auth/profile - cần JWT token
func (h *AuthHandler) GetProfile(ctx *gin.Context) {
	// Lấy từ context
	userIDVal, exists := ctx.Get(middleware.ContextKeyUserID)
	if !exists {
		utils.Unauthorized(ctx, "Không tìm thấy thông tin user")
		return
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		utils.Unauthorized(ctx, "UserID không hợp lệ")
		return
	}

	email, _ := ctx.Get(middleware.ContextKeyEmail)
	roleName, _ := ctx.Get(middleware.ContextKeyRoleName)

	emp, err := h.authSvc.GetProfile(userID)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.Success(ctx, "Lấy thông tin hồ sơ thành công!", gin.H{
		"user_id":   userID,
		"email":     email,
		"role_name": roleName,
		"employee":  emp,
	})
}

// Logout godoc
// POST /api/v1/auth/logout - cần JWT token
func (h *AuthHandler) Logout(ctx *gin.Context) {
	tokenStringVal, exists := ctx.Get("TokenString")
	if !exists {
		utils.Unauthorized(ctx, "Không tìm thấy token")
		return
	}
	tokenString := tokenStringVal.(string)

	remainingTimeVal, exists := ctx.Get("TokenRemainingTime")
	if !exists {
		utils.Unauthorized(ctx, "Không thể xác định thời gian hết hạn")
		return
	}
	remainingTime, ok := remainingTimeVal.(time.Duration)
	if !ok {
		utils.Unauthorized(ctx, "Dữ liệu thời gian không hợp lệ")
		return
	}

	err := h.authSvc.Logout(tokenString, remainingTime)
	if err != nil {
		utils.BadRequest(ctx, "Không thể đăng xuất: "+err.Error())
		return
	}

	utils.Success(ctx, "Đăng xuất thành công", nil)
}
