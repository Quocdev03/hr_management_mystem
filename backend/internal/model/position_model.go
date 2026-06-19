package model

// POSITION Model
// Chức vụ thuộc phòng ban
type Position struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	TimestampModel
}

// CreatePositionRequest - Tạo chức vụ mới
type CreatePositionRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description"`
}

// UpdatePositionRequest - Cập nhật chức vụ
type UpdatePositionRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
