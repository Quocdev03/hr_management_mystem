package service

import (
	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/internal/repository"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

type RoleService interface {
	GetAllRoles() ([]model.RoleResponse, error)
	GetRoleByID(id uint) (*model.RoleResponse, error)
	CreateRole(req model.CreateRoleRequest) (*model.Role, error)
	UpdateRole(id uint, req model.UpdateRoleRequest) error
	DeleteRole(id uint) error
}

type roleService struct {
	roleRepo repository.RoleRepository
	userRepo repository.UserRepository
	rdb      *redis.Client
}

func NewRoleService(roleRepo repository.RoleRepository, userRepo repository.UserRepository, rdb *redis.Client) RoleService {
	return &roleService{
		roleRepo: roleRepo,
		userRepo: userRepo,
		rdb:      rdb,
	}
}

func (s *roleService) GetAllRoles() ([]model.RoleResponse, error) {
	return s.roleRepo.FindAll()
}

func (s *roleService) GetRoleByID(id uint) (*model.RoleResponse, error) {
	return s.roleRepo.FindByID(id)
}

func (s *roleService) CreateRole(req model.CreateRoleRequest) (*model.Role, error) {
	nameLower := strings.ToLower(strings.TrimSpace(req.Name))
	exists, err := s.roleRepo.ExistsByName(nameLower)
	if err != nil {
		return nil, fmt.Errorf("lỗi kiểm tra tên vai trò: %w", err)
	}
	if exists {
		return nil, errors.New("tên vai trò đã tồn tại")
	}

	role := &model.Role{
		Name:        nameLower,
		Description: req.Description,
	}

	err = s.roleRepo.Create(role, req.Permissions)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (s *roleService) UpdateRole(id uint, req model.UpdateRoleRequest) error {
	roleResp, err := s.roleRepo.FindByID(id)
	if err != nil {
		return errors.New("không tìm thấy vai trò")
	}

	// Không cho sửa role admin (chỉ bảo vệ tên và xoá, có thể cho phép sửa permission hoặc description nếu cần)
	// Để bảo vệ cốt lõi, ta chặn sửa tên role admin
	if strings.EqualFold(roleResp.Role.Name, "admin") && req.Name != nil && !strings.EqualFold(*req.Name, "admin") {
		return errors.New("không thể đổi tên vai trò hệ thống (admin)")
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		newName := strings.ToLower(strings.TrimSpace(*req.Name))
		if newName != roleResp.Role.Name {
			exists, err := s.roleRepo.ExistsByName(newName)
			if err != nil {
				return err
			}
			if exists {
				return errors.New("tên vai trò đã tồn tại")
			}
			updates["name"] = newName
		}
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}

	err = s.roleRepo.Update(id, updates, req.Permissions)
	if err != nil {
		return err
	}

	// Invalidate Cache cho tất cả user thuộc role này (nếu permission thay đổi)
	// Lấy tất cả user mang role_id = id và xoá cache permissions:user:{userID}
	if req.Permissions != nil && s.rdb != nil {
		ctx := context.Background()
		var cursor uint64
		for {
			var keys []string
			keys, cursor, err = s.rdb.Scan(ctx, cursor, "permissions:user:*", 100).Result()
			if err == nil && len(keys) > 0 {
				s.rdb.Del(ctx, keys...)
			}
			if cursor == 0 {
				break
			}
		}
	}

	return nil
}

func (s *roleService) DeleteRole(id uint) error {
	roleResp, err := s.roleRepo.FindByID(id)
	if err != nil {
		return errors.New("không tìm thấy vai trò")
	}

	if strings.EqualFold(roleResp.Role.Name, "admin") || strings.EqualFold(roleResp.Role.Name, "hr") || strings.EqualFold(roleResp.Role.Name, "employee") {
		return errors.New("không thể xóa vai trò mặc định của hệ thống")
	}

	return s.roleRepo.Delete(id)
}
