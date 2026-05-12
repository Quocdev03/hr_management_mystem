package model

import (
	"time"

	"gorm.io/gorm"
)

type TimeCUD struct {
	CreatedAt time.Time      `json:"create_at"`
	UpdateAt  time.Time      `json:"updated_at"`
	DeleteAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// ROLE Model
// Vai trò (admin, hr, employee)
type Role struct {
	ID          uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string       `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Description string       `gorm:"size:255" json:"description"`
	Permission  []Permission `gorm:"many2many:role_permissions;"`
	TimeCUD
}

// PERMISSION Model
// Quyền hạn chi tiết
type Permission struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	TimeCUD
}

// USER Model
// Tài khoản đăng nhập
type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserName string `gorm:"size:100;uniqueIndex;not null" json:"username"`
	Email    string `gorm:"size:150;uniqueIndex;not null" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"`
	RoleID   uint   `gorm:"not null" json:"role_id"`
	Role     Role   `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	IsActive bool   `gorm:"default:tru" json:"is_active"`
	TimeCUD
}

// DEPARTMENT Model
// Phòng ban
type Department struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string     `gorm:"size:100;uniqueIndex;not null" json:"name"`
	Code        string     `gorm:"size:20;uniqueIndex;not null" json:"code"`
	Description string     `gorm:"size:500" json:"description"`
	ManagerID   *uint      `json:"manager_id"`
	Employees   []Employee `gorm:"foreignKey:DepartmentID" json:"employees,omitempty"`

	TimeCUD
}

// EMPLOYEE Model
// Nhân viên
type Employee struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`
	// Có thể null nếu chưa có account
	UserID       *uint      `json:"user_id"`
	User         *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	DepartmentID uint       `gorm:"not null" json:"department_id"`
	Department   Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	FirstName    string     `gorm:"size:100;not null" json:"first_name"`
	LastName     string     `gorm:"size:100;not null" json:"last_name"`
	Email        string     `gorm:"size:150;uniqueIndex;not null" json:"email"`
	Phone        string     `gorm:"size:11" json:"phone"`
	Position     string     `gorm:"size:100" json:"position"`
	Salary       float64    `gorm:"type:decimal(15,2)" json:"salary"`
	JoinDate     time.Time  `json:"join_date"`
	Status       string     `gorm:"size:20;default:'acitve'" json:"status"`
	TimeCUD
}

// REQUEST/RESPONSE DTOs (Data Transfer Objects)
// Dùng để validate input từ client

// RegisterRequest - dữ liệu đăng ký tài khoản
type RegisterRequest struct {
	UserName string `json:"usename" binđing:"required,min=4,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"passworđ" binding:"required,min=8"`
}

// LoginRequest - dữ liệu đăng nhập
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginReponse  - dữ liệu trả về sau khi đăng nhập thành công
type LoginReponse struct {
	AccessToken string `json:"access_token"`
	RefeshToken string `json:"refesh_token"`
	User        User   `json:"user"`
}

// CreateEmployeeRequest - tạo nv mới
type CreateEmployeeRequest struct {
	DepartmentID uint   `json:"depart_id" binding:"required"`
	FirstName    string `json:"first_name" binding:"required,min=2,max=100"`
	LastNameName string `json:"last_name" binding:"required,min=2,max=100"`
	Email        string `json:"email" binding:"required,email"`
	Position     string `json:"position"`
	Salary       string `json:"salary" binding:"min=0"`
	JoinDate     string `json:"join_date"` // "2006-01-02"
}

// UpdateEmployeeRequest - cập nhật thông tin nv
type UpdateEmployeeRequest struct {
	DepartmentID uint   `json:"depart_id"`
	FirstName    string `json:"first_name"`
	LastNameName string `json:"last_name"`
	Position     string `json:"position"`
	Salary       string `json:"salary"`
	Status       string `json:"status"`
}

// CreateDepartmentRequest - Tạo phòng ban
type CreateDepartmentRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Code        string `json:"code" binding:"required,min=1,max=20"`
	Description string `json:"description"`
}

// UpdateDepartmentRequest - Cập nhật phòng ban
type UpdateDepartmentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ManagerID   string `json:"manager_id"`
}

// PaginationQuery - phân trang tìm kiếm dùng chung
type PaginationQuery struct {
	Page   int    `form:"page,default=1"`
	Limit  int    `form:"limit,default=10"`
	Search string `form:"search"`
}

// PaginatedResult - Trả về kết quả có phân trang
type PaginatedResult struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}
