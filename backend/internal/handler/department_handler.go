package handler

import (
	"chiquoc_hocgolang/internal/common"
	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/internal/service"
	"chiquoc_hocgolang/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
	deptSvc service.DepartmentService
}

func NewDepartmentHandler(deptSvc service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{
		deptSvc: deptSvc,
	}
}

// GetDepartments godoc
// GET /api/v1/departments?page=1&limit=10&search=IT
func (h *DepartmentHandler) GetDepartments(ctx *gin.Context) {
	var query model.PaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.BadRequest(ctx, "Tham số truy vấn không hợp lệ")
		return
	}

	common.NormalizePagination(&query)

	result, err := h.deptSvc.GetDepartments(query)
	if err != nil {
		utils.InternalServerError(ctx, err.Error())
		return
	}

	utils.Success(ctx, "Lấy danh sách phòng ban thành công!", result)
}

// GetDepartment godoc
// GET /api/v1/departments/:id
func (h *DepartmentHandler) GetDepartment(ctx *gin.Context) {
	id, ok := common.ParseAndValidateID(ctx, "phòng ban")
	if !ok {
		return
	}

	dept, err := h.deptSvc.GetDepartmentByID(id)
	if err != nil {
		utils.NotFound(ctx, err.Error())
		return
	}
	utils.Success(ctx, "Lấy thông tin phòng ban thành công!", dept)
}

// CreateDepartment godoc
// POST /api/v1/departments
func (h *DepartmentHandler) CreateDepartment(ctx *gin.Context) {
	var req model.CreateDepartmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Dữ liệu không đúng định dạng JSON")
		return
	}

	// Trim whitespace và chuẩn hoá code thành uppercase
	req.Name = strings.TrimSpace(req.Name)
	req.Code = strings.TrimSpace(strings.ToUpper(req.Code))
	req.Description = strings.TrimSpace(req.Description)

	// Validate đầu vào
	if verrs := utils.ValidateCreateDepartment(req.Name, req.Code, req.Description); verrs != nil {
		utils.ValidationError(ctx, "Dữ liệu tạo phòng ban không hợp lệ", verrs.Errors)
		return
	}

	dept, err := h.deptSvc.CreateDepartment(req)
	if err != nil {
		utils.Conflict(ctx, err.Error())
		return
	}
	utils.Created(ctx, "Tạo phòng ban thành công!", dept)
}

// UpdateDepartment godoc
// PUT /api/v1/departments/:id
func (h *DepartmentHandler) UpdateDepartment(ctx *gin.Context) {
	id, ok := common.ParseAndValidateID(ctx, "phòng ban")
	if !ok {
		return
	}

	var req model.UpdateDepartmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Dữ liệu không đúng định dạng JSON")
		return
	}

	// Trim whitespace
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)

	// Validate đầu vào
	if verrs := utils.ValidateUpdateDepartment(req.Name, req.Description, req.ManagerID); verrs != nil {
		utils.ValidationError(ctx, "Dữ liệu cập nhật phòng ban không hợp lệ", verrs.Errors)
		return
	}

	dept, err := h.deptSvc.UpdateDepartment(id, req)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}
	utils.Success(ctx, "Cập nhật phòng ban thành công!", dept)
}

// DeleteDepartment godoc
// DELETE /api/v1/departments/:id
func (h *DepartmentHandler) DeleteDepartment(ctx *gin.Context) {
	id, ok := common.ParseAndValidateID(ctx, "phòng ban")
	if !ok {
		return
	}

	if err := h.deptSvc.DeleteDepartment(id); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}
	utils.Success(ctx, "Xoá phòng ban thành công!", nil)
}
