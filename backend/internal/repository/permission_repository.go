package repository

import (
	"chiquoc_hocgolang/internal/model"

	"gorm.io/gorm"
)

// PermissionRepository kiểm tra quyền theo permission code.
type PermissionRepository interface {
	HasPermission(userID uint, permissionCode string) (bool, error)
	GetPermissionCodes(userID uint) ([]string, error)
	GetAllPermissions() ([]model.Permission, error)
	SetUserPermissions(userID uint, permissionCodes []string) error
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) GetPermissionCodes(userID uint) ([]string, error) {
	if userID == 0 {
		return []string{}, nil
	}

	var codes []string
	if err := r.db.Table("permissions").
		Distinct().
		Joins("JOIN role_permissions rp ON rp.permission_id = permissions.id").
		Joins("JOIN users u ON u.role_id = rp.role_id").
		Where("u.id = ?", userID).
		Pluck("permissions.code", &codes).Error; err != nil {
		return nil, err
	}

	var userCodes []string
	if err := r.db.Table("permissions").
		Distinct().
		Joins("JOIN user_permissions up ON up.permission_id = permissions.id").
		Where("up.user_id = ?", userID).
		Pluck("permissions.code", &userCodes).Error; err != nil {
		return nil, err
	}

	seen := map[string]struct{}{}
	merged := make([]string, 0, len(codes)+len(userCodes))
	for _, code := range append(codes, userCodes...) {
		if _, ok := seen[code]; ok {
			continue
		}
		seen[code] = struct{}{}
		merged = append(merged, code)
	}
	return merged, nil
}

func (r *permissionRepository) GetAllPermissions() ([]model.Permission, error) {
	var perms []model.Permission
	if err := r.db.Order("code ASC").Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

func (r *permissionRepository) SetUserPermissions(userID uint, permissionCodes []string) error {
	if userID == 0 {
		return nil
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserPermission{}).Error; err != nil {
			return err
		}

		if len(permissionCodes) == 0 {
			return nil
		}

		var perms []model.Permission
		if err := tx.Where("code IN ?", permissionCodes).Find(&perms).Error; err != nil {
			return err
		}

		for _, perm := range perms {
			if err := tx.Create(&model.UserPermission{UserID: userID, PermissionID: perm.ID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *permissionRepository) HasPermission(userID uint, permissionCode string) (bool, error) {
	if userID == 0 || permissionCode == "" {
		return false, nil
	}

	var roleCount int64
	if err := r.db.Model(&model.Permission{}).
		Joins("JOIN role_permissions rp ON rp.permission_id = permissions.id").
		Joins("JOIN users u ON u.role_id = rp.role_id").
		Where("u.id = ? AND permissions.code = ?", userID, permissionCode).
		Count(&roleCount).Error; err != nil {
		return false, err
	}
	if roleCount > 0 {
		return true, nil
	}

	var userCount int64
	if err := r.db.Model(&model.Permission{}).
		Joins("JOIN user_permissions up ON up.permission_id = permissions.id").
		Where("up.user_id = ? AND permissions.code = ?", userID, permissionCode).
		Count(&userCount).Error; err != nil {
		return false, err
	}

	return userCount > 0, nil
}
