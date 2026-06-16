package middleware

import (
	"chiquoc_hocgolang/internal/repository"
	"chiquoc_hocgolang/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequirePermission(permRepo repository.PermissionRepository, permissionCode string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIDValue, exists := ctx.Get(ContextKeyUserID)
		if !exists {
			utils.Unauthorized(ctx, "User not found in context")
			ctx.Abort()
			return
		}

		userID, ok := userIDValue.(uint)
		if !ok {
			utils.Unauthorized(ctx, "Invalid user id in context")
			ctx.Abort()
			return
		}

		hasPerm, err := permRepo.HasPermission(userID, permissionCode)
		if err != nil {
			utils.InternalServerError(ctx, fmt.Sprintf("Permission check failed: %v", err))
			ctx.Abort()
			return
		}
		if !hasPerm {
			utils.Forbidden(ctx, "Bạn không có quyền thực hiện hành động này")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func RequireAnyPermission(permRepo repository.PermissionRepository, permissionCodes ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIDValue, exists := ctx.Get(ContextKeyUserID)
		if !exists {
			utils.Unauthorized(ctx, "User not found in context")
			ctx.Abort()
			return
		}

		userID, ok := userIDValue.(uint)
		if !ok {
			utils.Unauthorized(ctx, "Invalid user id in context")
			ctx.Abort()
			return
		}

		for _, code := range permissionCodes {
			hasPerm, err := permRepo.HasPermission(userID, code)
			if err == nil && hasPerm {
				ctx.Next()
				return
			}
		}

		utils.Forbidden(ctx, "Bạn không có quyền truy cập tài nguyên này")
		ctx.Abort()
	}
}

func PermissionDenied() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Forbidden"})
		ctx.Abort()
	}
}
