package middleware

import (
	"chiquoc_hocgolang/internal/config"
	"chiquoc_hocgolang/internal/utils"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Dùng constants để tránh typo khi get/set vào gin.Context
const (
	ContextKeyUserID   = "userID"
	ContextKeyEmail    = "username"
	ContextKeyRoleID   = "roleID"
	ContextKeyRoleName = "roleName"
)

// Logger Middleware
// Ghi log mỗi request: method, path, status, thời gian xử lý

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		method := ctx.Request.Method

		// Xử lý request
		ctx.Next()

		// Ghi log
		duration := time.Since(start)
		statusCode := ctx.Writer.Status()

		utils.Info("%s %s | %d | %v | %s", method, path, statusCode, duration, ctx.ClientIP())

	}
}

// Recover Bắt panic, ngăn server crash, trả về 500 thay vì terminate

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log lỗi
				utils.Error("Panic recovered: %v", err)

				// Trả về 500 thay vì crash server
				ctx.JSON(http.StatusInternalServerError, utils.Response{
					Success: false,
					Message: fmt.Sprintf("Internal server error: %v", err),
				})
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}

// Xác thực JWT token trong header Authorization
func AuthJWT(cfg *config.JWTConfig, rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Lấy token từ header trước
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(ctx, "Authorization header cần phải có")
			ctx.Abort()
			return
		}
		// Tách "Bearer " prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Unauthorized(ctx, "Định dạng Authorization không hợp lệ! Cần phải có Bearer <token>")
			ctx.Abort()
			return
		}

		tokenString := parts[1]
		// Validate và parse token
		claims, err := utils.ValidateToken(tokenString, cfg.SecretKey)
		if err != nil {
			utils.Unauthorized(ctx, "Invalid or expired token: "+err.Error())
			ctx.Abort()
			return
		}

		// Lưu thông tin user vào context để các handler sau dùng
		ctx.Set(ContextKeyUserID, claims.UserID)
		ctx.Set(ContextKeyEmail, claims.Email)
		ctx.Set(ContextKeyRoleID, claims.RoleID)
		ctx.Set(ContextKeyRoleName, claims.RoleName)

		// Set Token cho Logout
		ctx.Set("TokenString", tokenString)
		ctx.Set("TokenRemainingTime", time.Until(claims.ExpiresAt.Time))

		// Kiểm tra Token đã bị đăng xuất chưa (Blacklist)
		if rdb != nil {
			isBlacklisted, err := rdb.Exists(ctx.Request.Context(), "blacklist:"+tokenString).Result()
			if err == nil && isBlacklisted > 0 {
				utils.Unauthorized(ctx, "Token đã bị vô hiệu hoá (Đăng xuất)")
				ctx.Abort()
				return
			}
		}

		ctx.Next()

	}
}

// Rolo Kiểm tra user có role được phép không
// kiểm tra user có một trong các roles được chỉ định không

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Lấy role từ context (đã được set bởi AuthJWT)
		roleName, exists := ctx.Get(ContextKeyRoleName)
		if !exists {
			utils.Unauthorized(ctx, "User role not found in context")
			ctx.Abort()
			return
		}

		userRole := roleName.(string)

		// Kiểm tra role có trong danh sách allowed không
		for _, allowed := range allowedRoles {
			if userRole == allowed {
				ctx.Next()
				return
			}
		}

		utils.Forbidden(ctx, fmt.Sprintf("Role '%s' is not allowed to access this resource", userRole))
		ctx.Abort()
	}
}

// Cross Cho phép cross-origin requests (cần thiết khi frontend ở domain khác)
func CORS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		// Preflight request
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}
