package model

import "time"

// EMPLOYEE Model
// Nhân viên
type Employee struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       *uint      `gorm:"index" json:"user_id"`
	User         *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	DepartmentID uint       `gorm:"not null" json:"department_id"`
	Department   Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	PositionID   uint       `gorm:"not null;index" json:"position_id"`
	Position     *Position  `gorm:"foreignKey:PositionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"position,omitempty"`
	FirstName    string     `gorm:"size:100;not null" json:"first_name"`
	LastName     string     `gorm:"size:100;not null" json:"last_name"`
	Phone        string     `gorm:"size:11" json:"phone"`
	Salary       float64    `gorm:"type:decimal(15,2)" json:"salary"`
	JoinDate     time.Time  `json:"join_date"`
	BirthDate    *time.Time `json:"birth_date"`
	Gender       string     `gorm:"size:10;default:'male'" json:"gender"`
	Status       string     `gorm:"size:20;default:'active'" json:"status"`
	TimestampModel
}

// CreateEmployeeRequest - tạo nv mới
type CreateEmployeeRequest struct {
	DepartmentID uint    `json:"department_id" binding:"required"`
	PositionID   uint    `json:"position_id" binding:"required"`
	FirstName    string  `json:"first_name" binding:"required,min=2,max=100"`
	LastName     string  `json:"last_name" binding:"required,min=2,max=100"`
	Phone        string  `json:"phone" binding:"required,startswith=0,min=10,max=11,numeric"`
	Salary       float64 `json:"salary" binding:"min=0"`
	JoinDate     string  `json:"join_date"`  // "2006-01-02"
	BirthDate    string  `json:"birth_date"` // "2006-01-02"
	Gender       string  `json:"gender"`     // "male" | "female" | "other"
	Status       string  `json:"status"`
	UserID       *uint   `json:"user_id"`
}

// UpdateEmployeeRequest - cập nhật thông tin nv
type UpdateEmployeeRequest struct {
	DepartmentID *uint    `json:"department_id"`
	PositionID   *uint    `json:"position_id"`
	FirstName    *string  `json:"first_name"`
	LastName     *string  `json:"last_name"`
	Phone        *string  `json:"phone" binding:"omitempty,startswith=0,min=10,max=11,numeric"`
	Salary       *float64 `json:"salary"`
	Status       *string  `json:"status"`
	BirthDate    *string  `json:"birth_date"`
	Gender       *string  `json:"gender"`
	UserID       *uint    `json:"user_id"`
}
