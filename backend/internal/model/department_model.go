package model

// DEPARTMENT Model
// Phòng ban
type Department struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string     `gorm:"size:100;uniqueIndex;not null" json:"name"`
	Code        string     `gorm:"size:20;uniqueIndex;not null" json:"code"`
	Description string     `gorm:"size:500" json:"description"`
	ManagerID   *uint      `json:"manager_id"`
	Manager     *Employee  `gorm:"foreignKey:ManagerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"manager,omitempty"`
	Employees   []Employee `gorm:"foreignKey:DepartmentID" json:"employees,omitempty"`

	TimestampModel
}

// CreateDepartmentRequest - Tạo phòng ban
type CreateDepartmentRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Code        string `json:"code" binding:"required,min=1,max=20"`
	Description string `json:"description"`
	ManagerID   *uint  `json:"manager_id"`
}

// UpdateDepartmentRequest - Cập nhật phòng ban
type UpdateDepartmentRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	// ManagerID: nil = không thay đổi, 0 = xoá manager, >0 = set manager mới
	ManagerID *uint `json:"manager_id"`
}
