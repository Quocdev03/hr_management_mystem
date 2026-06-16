package model

// Permission định nghĩa hành động có thể cấp cho vai trò/user.
type Permission struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Code        string `gorm:"size:100;not null;uniqueIndex" json:"code"`
	Description string `gorm:"size:255" json:"description"`
	TimestampModel
}

// RolePermission ánh xạ quyền cho vai trò.
type RolePermission struct {
	RoleID       uint `gorm:"primaryKey" json:"role_id"`
	PermissionID uint `gorm:"primaryKey" json:"permission_id"`
	TimestampModel
}

// UserPermission cho phép override quyền riêng cho user cụ thể.
type UserPermission struct {
	UserID       uint `gorm:"primaryKey" json:"user_id"`
	PermissionID uint `gorm:"primaryKey" json:"permission_id"`
	TimestampModel
}
