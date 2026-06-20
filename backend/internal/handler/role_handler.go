package handler

import (
	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/internal/service"
	"chiquoc_hocgolang/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleSvc service.RoleService
}

func NewRoleHandler(roleSvc service.RoleService) *RoleHandler {
	return &RoleHandler{roleSvc: roleSvc}
}

func (h *RoleHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.roleSvc.GetAllRoles()
	if err != nil {
		utils.InternalServerError(c, "Lỗi khi lấy danh sách vai trò")
		return
	}
	utils.Success(c, "Lấy danh sách vai trò thành công!", roles)
}

func (h *RoleHandler) GetRoleByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID vai trò không hợp lệ")
		return
	}

	role, err := h.roleSvc.GetRoleByID(uint(id))
	if err != nil {
		utils.NotFound(c, "Không tìm thấy vai trò")
		return
	}
	utils.Success(c, "Lấy thông tin vai trò thành công!", role)
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req model.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	role, err := h.roleSvc.CreateRole(req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Created(c, "Tạo vai trò thành công!", role)
}

func (h *RoleHandler) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID vai trò không hợp lệ")
		return
	}

	var req model.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.roleSvc.UpdateRole(uint(id), req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, "Cập nhật vai trò thành công!", nil)
}

func (h *RoleHandler) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID vai trò không hợp lệ")
		return
	}

	if err := h.roleSvc.DeleteRole(uint(id)); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, "Xóa vai trò thành công!", nil)
}
