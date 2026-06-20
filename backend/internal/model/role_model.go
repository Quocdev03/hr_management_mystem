package model

// ROLE Model
// Vai trò (admin, hr, employee)
type Role struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"size:50;not null;index" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	TimestampModel
}

// RoleResponse chứa thông tin role kèm permission
type RoleResponse struct {
	Role
	Permissions []string `json:"permissions"`
}

type CreateRoleRequest struct {
	Name        string   `json:"name" binding:"required,min=2,max=50"`
	Description string   `json:"description" binding:"max=255"`
	Permissions []string `json:"permissions"`
}

type UpdateRoleRequest struct {
	Name        *string  `json:"name" binding:"omitempty,min=2,max=50"`
	Description *string  `json:"description" binding:"max=255"`
	Permissions []string `json:"permissions"`
}
