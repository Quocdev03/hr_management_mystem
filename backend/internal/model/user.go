package model

// ROLE Model
// Vai trò (admin, hr, employee)
type Role struct {
	ID          uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string       `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Description string       `gorm:"size:255" json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	TimestampModel
}

// PERMISSION Model
// Quyền hạn chi tiết
type Permission struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	TimestampModel
}

// USER Model
// Tài khoản đăng nhập
type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserName string `gorm:"size:100;uniqueIndex;not null" json:"user_name"`
	Email    string `gorm:"size:150;uniqueIndex;not null" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"`
	RoleID   uint   `gorm:"not null" json:"role_id"`
	Role     Role   `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
	TimestampModel
}
