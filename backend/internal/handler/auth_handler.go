package handler

import (
	"chiquoc_hocgolang/internal/middleware"
	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/internal/service"
	"chiquoc_hocgolang/package/response"
	"chiquoc_hocgolang/package/validation"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authSvc service.AuthService
}

func NewAuthHandler(authScv service.AuthService) *AuthHandler {
	return &AuthHandler{
		authSvc: authScv,
	}
}

// Login godoc
// POST /api/v1/auth/login
func (h *AuthHandler) Login(ctx *gin.Context) {
	var req model.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Dữ liệu đăng nhập không hợp lệ")
		return
	}

	// Trim whitespace
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	// Validate đầu vào
	if verrs := validation.ValidateLogin(req.Email, req.Password); verrs != nil {
		response.ValidationError(ctx, "Dữ liệu đăng nhập không hợp lệ", verrs.Errors)
		return
	}

	result, err := h.authSvc.Login(req)
	if err != nil {
		response.Unauthorized(ctx, err.Error())
		return
	}
	response.Success(ctx, "Đăng nhập thành công", result)
}

// GetProfile godoc
// GET /api/v1/auth/profile - cần JWT token
func (h *AuthHandler) GetProfile(ctx *gin.Context) {
	// Lấy thông tin user từ context set bởi AuthJWT middleware
	userID, _ := ctx.Get(middleware.ContextKeyUserID)
	email, _ := ctx.Get(middleware.ContextKeyEmail)
	roleNam, _ := ctx.Get(middleware.ContextKeyRoleName)

	response.Success(ctx, "Lấy thông tin hồ sơ thành công!", gin.H{
		"user_id":   userID,
		"email":     email,
		"role_name": roleNam,
	})
}
