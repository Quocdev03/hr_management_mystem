package repository

import (
	"context"
	"fmt"
	"time"

	"chiquoc_hocgolang/internal/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// PermissionRepository kiểm tra quyền theo permission code.
type PermissionRepository interface {
	HasPermission(userID uint, permissionCode string) (bool, error)
	GetPermissionCodes(userID uint) ([]string, error)
	GetAllPermissions() ([]model.Permission, error)
}

type permissionRepository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewPermissionRepository(db *gorm.DB, rdb *redis.Client) PermissionRepository {
	return &permissionRepository{db: db, rdb: rdb}
}

func (r *permissionRepository) GetPermissionCodes(userID uint) ([]string, error) {
	if userID == 0 {
		return []string{}, nil
	}

	var codes []string
	if err := r.db.Table("permissions").
		Distinct().
		Joins("JOIN role_permissions rp ON rp.permission_id = permissions.id").
		Joins("JOIN users u ON u.role_id = rp.role_id AND u.deleted_at IS NULL").
		Where("u.id = ? AND permissions.deleted_at IS NULL", userID).
		Pluck("permissions.code", &codes).Error; err != nil {
		return nil, err
	}

	return codes, nil
}

func (r *permissionRepository) GetAllPermissions() ([]model.Permission, error) {
	var perms []model.Permission
	if err := r.db.Order("code ASC").Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}


func (r *permissionRepository) HasPermission(userID uint, permissionCode string) (bool, error) {
	if userID == 0 || permissionCode == "" {
		return false, nil
	}

	cacheKey := fmt.Sprintf("permissions:user:%d", userID)

	// 1. Kiểm tra Redis Cache
	if r.rdb != nil {
		ctx := context.Background()
		exists, err := r.rdb.Exists(ctx, cacheKey).Result()
		if err == nil && exists > 0 {
			// Cache hit
			isMember, err := r.rdb.SIsMember(ctx, cacheKey, permissionCode).Result()
			if err == nil {
				return isMember, nil
			}
		}
	}

	// 2. Cache miss -> Lấy tất cả quyền của User từ DB
	userPerms, err := r.GetPermissionCodes(userID)
	if err != nil {
		return false, err
	}

	// 3. Đưa vào Redis Cache (Set)
	if r.rdb != nil && len(userPerms) > 0 {
		ctx := context.Background()
		// Convert []string to []interface{} for SAdd
		var members []interface{}
		for _, p := range userPerms {
			members = append(members, p)
		}
		
		// Pipeline để SADD và EXPIRE atomatically
		pipe := r.rdb.Pipeline()
		pipe.SAdd(ctx, cacheKey, members...)
		pipe.Expire(ctx, cacheKey, 24*time.Hour)
		_, _ = pipe.Exec(ctx) // Ignore cache write errors
	} else if r.rdb != nil && len(userPerms) == 0 {
		// Tránh cache stampede bằng cách cache mảng rỗng (với SADD 1 phần tử rác, hoặc dùng chuỗi rỗng)
		ctx := context.Background()
		r.rdb.SAdd(ctx, cacheKey, "_empty_")
		r.rdb.Expire(ctx, cacheKey, 24*time.Hour)
	}

	// 4. Kiểm tra quyền
	for _, p := range userPerms {
		if p == permissionCode {
			return true, nil
		}
	}

	return false, nil
}
