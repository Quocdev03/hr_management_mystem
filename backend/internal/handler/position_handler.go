package handler

import (
	"chiquoc_hocgolang/internal/common"
	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/internal/service"
	"chiquoc_hocgolang/internal/utils"

	"github.com/gin-gonic/gin"
)

type PositionHandler struct {
	posSvc service.PositionService
}

func NewPositionHandler(posSvc service.PositionService) *PositionHandler {
	return &PositionHandler{
		posSvc: posSvc,
	}
}

// GetPositions godoc
// GET /api/v1/positions?department_id=1
func (h *PositionHandler) GetPositions(ctx *gin.Context) {
	positions, err := h.posSvc.GetAll()
	if err != nil {
		utils.InternalServerError(ctx, err.Error())
		return
	}
	utils.Success(ctx, "Lấy danh sách chức vụ thành công!", positions)
}

// GetPosition godoc
// GET /api/v1/positions/:id
func (h *PositionHandler) GetPosition(ctx *gin.Context) {
	id, ok := common.ParseAndValidateID(ctx, "chức vụ")
	if !ok {
		return
	}

	pos, err := h.posSvc.GetByID(id)
	if err != nil {
		utils.NotFound(ctx, err.Error())
		return
	}
	utils.Success(ctx, "Lấy thông tin chức vụ thành công!", pos)
}

// CreatePosition godoc
// POST /api/v1/positions
func (h *PositionHandler) CreatePosition(ctx *gin.Context) {
	var req model.CreatePositionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Dữ liệu không đúng định dạng JSON")
		return
	}

	pos, err := h.posSvc.Create(req)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}
	utils.Created(ctx, "Tạo chức vụ thành công!", pos)
}

// UpdatePosition godoc
// PATCH /api/v1/positions/:id
func (h *PositionHandler) UpdatePosition(ctx *gin.Context) {
	id, ok := common.ParseAndValidateID(ctx, "chức vụ")
	if !ok {
		return
	}

	var req model.UpdatePositionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Dữ liệu không đúng định dạng JSON")
		return
	}

	pos, err := h.posSvc.Update(id, req)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}
	utils.Success(ctx, "Cập nhật chức vụ thành công!", pos)
}

// DeletePosition godoc
// DELETE /api/v1/positions/:id
func (h *PositionHandler) DeletePosition(ctx *gin.Context) {
	id, ok := common.ParseAndValidateID(ctx, "chức vụ")
	if !ok {
		return
	}

	if err := h.posSvc.Delete(id); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}
	utils.Success(ctx, "Xoá chức vụ thành công!", nil)
}
