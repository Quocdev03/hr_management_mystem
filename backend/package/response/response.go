package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response là cấu trúc JSON response thống nhất cho toàn bộ API
// Mọi endpoint đều trả về format này để frontend dễ xử lý
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Success200 - trả về 200 OK với data
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created201 - trả về 201 Created
func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// BadRequest400 - lỗi dữ liệu đầu vào không hợp lệ
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: message,
	})
}

// Unauthorized401 - chưa đăng nhập hoặc token không hợp lệ
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Success: false,
		Message: message,
	})
}

// Forbidden403 - đã đăng nhập nhưng không có quyền
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Success: false,
		Message: message,
	})
}

// NotFound404 - không tìm thấy tài nguyên
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Success: false,
		Message: message,
	})
}

// InternalServerError500 - lỗi server
func InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Message: message,
	})
}

// ValidationError - lỗi validate với chi tiết lỗi
func ValidationError(c *gin.Context, message string, errors interface{}) {
	c.JSON(http.StatusUnprocessableEntity, Response{
		Success: false,
		Message: message,
		Error:   errors,
	})
}

// NoContent204 - trả về 204 Không có nội dung
func NoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

// Accepted202 - đã chấp nhận yêu cầu nhưng đang xử lý
func Accepted(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusAccepted, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Conflict409 - tài nguyên bị xung đột (ví dụ: email đã tồn tại)
func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, Response{
		Success: false,
		Message: message,
	})
}

// TooManyRequests429 - gửi quá nhiều yêu cầu
func TooManyRequests(c *gin.Context, message string) {
	c.JSON(http.StatusTooManyRequests, Response{
		Success: false,
		Message: message,
	})
}
