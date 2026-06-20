package repository

import (
	"chiquoc_hocgolang/internal/model"
	"errors"

	"gorm.io/gorm"
)

type RoleRepository interface {
	FindAll() ([]model.RoleResponse, error)
	FindByID(id uint) (*model.RoleResponse, error)
	Create(role *model.Role, permissionCodes []string) error
	Update(id uint, updates map[string]interface{}, permissionCodes []string) error
	Delete(id uint) error
	ExistsByName(name string) (bool, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindAll() ([]model.RoleResponse, error) {
	var roles []model.Role
	if err := r.db.Order("id ASC").Find(&roles).Error; err != nil {
		return nil, err
	}

	var results []model.RoleResponse
	for _, role := range roles {
		perms, err := r.getRolePermissions(role.ID)
		if err != nil {
			return nil, err
		}
		results = append(results, model.RoleResponse{
			Role:        role,
			Permissions: perms,
		})
	}
	return results, nil
}

func (r *roleRepository) FindByID(id uint) (*model.RoleResponse, error) {
	var role model.Role
	if err := r.db.First(&role, id).Error; err != nil {
		return nil, err
	}
	perms, err := r.getRolePermissions(id)
	if err != nil {
		return nil, err
	}
	return &model.RoleResponse{
		Role:        role,
		Permissions: perms,
	}, nil
}

func (r *roleRepository) getRolePermissions(roleID uint) ([]string, error) {
	var codes []string
	if err := r.db.Table("permissions").
		Joins("JOIN role_permissions rp ON rp.permission_id = permissions.id").
		Where("rp.role_id = ? AND permissions.deleted_at IS NULL", roleID).
		Pluck("permissions.code", &codes).Error; err != nil {
		return nil, err
	}
	if codes == nil {
		codes = []string{}
	}
	return codes, nil
}

func (r *roleRepository) Create(role *model.Role, permissionCodes []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(role).Error; err != nil {
			return err
		}
		return r.setRolePermissions(tx, role.ID, permissionCodes)
	})
}

func (r *roleRepository) Update(id uint, updates map[string]interface{}, permissionCodes []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if len(updates) > 0 {
			if err := tx.Model(&model.Role{}).Where("id = ?", id).Updates(updates).Error; err != nil {
				return err
			}
		}
		if permissionCodes != nil {
			return r.setRolePermissions(tx, id, permissionCodes)
		}
		return nil
	})
}

func (r *roleRepository) setRolePermissions(tx *gorm.DB, roleID uint, permissionCodes []string) error {
	if err := tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
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
		if err := tx.Create(&model.RolePermission{RoleID: roleID, PermissionID: perm.ID}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *roleRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Kiểm tra có user nào đang dùng role này không
		var count int64
		if err := tx.Model(&model.User{}).Where("role_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("không thể xóa vai trò đang được sử dụng bởi người dùng")
		}

		// Xóa role_permissions
		if err := tx.Where("role_id = ?", id).Delete(&model.RolePermission{}).Error; err != nil {
			return err
		}

		// Xóa role
		if err := tx.Delete(&model.Role{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *roleRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Role{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}
