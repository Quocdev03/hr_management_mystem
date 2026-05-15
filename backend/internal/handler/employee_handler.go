package handler

import (
	"chiquoc_hocgolang/internal/common"
	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/internal/service"
	"chiquoc_hocgolang/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	empSvc service.EmployeeService
}

func NewEmployeeHandler(empSvc service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		empSvc: empSvc,
	}
}

// GetEmployees godoc
// GET /api/v1/employees?page=1&limit=10&search=quoc
func (h *EmployeeHandler) GetEmployees(ctx *gin.Context) {
	var query model.PaginationQuery

	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.BadRequest(ctx, "Tham số truy vấn không hợp lệ")
		return
	}

	common.NormalizePagination(&query)

	result, err := h.empSvc.GetEmployees(query)
	if err != nil {
		utils.InternalServerError(ctx, err.Error())
		return
	}
	utils.Success(ctx, "Lấy danh sách nhân viên thành công!", result)
}

// GetEmployee godoc
// GET /api/v1/employees/:id
func (h *EmployeeHandler) GetEmployee(ctx *gin.Context) {
	id, ok := common.ParseAndValidateID(ctx, "nhân viên")
	if !ok {
		return
	}

	emp, err := h.empSvc.GetEmployeeByID(id)
	if err != nil {
		utils.NotFound(ctx, err.Error())
		return
	}
	utils.Success(ctx, "Lấy thông tin nhân viên thành công", emp)
}

// CreateEmployee godoc
// POST /api/v1/employees
func (h *EmployeeHandler) CreateEmployee(ctx *gin.Context) {
	var req model.CreateEmployeeRequest

	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		utils.BadRequest(ctx, "Dữ liệu không đúng định dạng JSON")
		return
	}

	// Trim whitespace cho tất cả các trường string
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Phone = strings.TrimSpace(req.Phone)
	req.Position = strings.TrimSpace(req.Position)
	req.JoinDate = strings.TrimSpace(req.JoinDate)

	// Validate đầu vào chi tiết
	if verrs := utils.ValidateCreateEmployee(
		req.DepartmentID, req.FirstName, req.LastName,
		req.Email, req.Phone, req.Position, req.JoinDate, req.Salary,
	); verrs != nil {
		utils.ValidationError(ctx, "Dữ liệu tạo nhân viên không hợp lệ", verrs.Errors)
		return
	}

	emp, err := h.empSvc.Create(req)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.Created(ctx, "Tạo nhân viên thành công!", emp)
}

// UpdateEmployee godoc
// PUT /api/v1/employees/:id
func (h *EmployeeHandler) UpdateEmployee(ctx *gin.Context) {
	id, ok := common.ParseAndValidateID(ctx, "nhân viên")
	if !ok {
		return
	}

	var req model.UpdateEmployeeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Dữ liệu không đúng định dạng JSON")
		return
	}

	// Trim whitespace
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Position = strings.TrimSpace(req.Position)
	req.Status = strings.TrimSpace(strings.ToLower(req.Status))

	// Validate đầu vào
	if verrs := utils.ValidateUpdateEmployee(
		req.DepartmentID, req.FirstName, req.LastName,
		req.Phone, req.Position, req.Status, req.Salary,
	); verrs != nil {
		utils.ValidationError(ctx, "Dữ liệu cập nhật nhân viên không hợp lệ", verrs.Errors)
		return
	}

	emp, err := h.empSvc.UpdateEmployee(id, req)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}
	utils.Success(ctx, "Cập nhật thông tin nhân viên thành công", emp)
}

// DeleteEmployee godoc
// DELETE /api/v1/employees/:id
func (h *EmployeeHandler) DeleteEmployee(ctx *gin.Context) {
	id, ok := common.ParseAndValidateID(ctx, "nhân viên")
	if !ok {
		return
	}

	if err := h.empSvc.DeleteEmployee(id); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}
	utils.Success(ctx, "Xoá nhân viên thành công!", nil)
}
