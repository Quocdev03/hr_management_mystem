package handler

import (
	"chiquoc_hocgolang/internal/common"
	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/internal/service"
	"chiquoc_hocgolang/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc service.UserService
}

func NewUserHandler(userSvc service.UserService) *UserHandler {
	return &UserHandler{
		userSvc: userSvc,
	}
}

func (h *UserHandler) GetUsers(ctx *gin.Context) {
	var query model.PaginationQuery

	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.BadRequest(ctx, "Tham số truy vấn không hợp lệ")
		return
	}

	common.NormalizePagination(&query)

	result, err := h.userSvc.GetUsers(query)
	if err != nil {
		utils.InternalServerError(ctx, err.Error())
		return
	}

	utils.Success(ctx, "Lấy danh sách user thành công", result)
}

func (h *UserHandler) GetUser(ctx *gin.Context) {
	id, ok := common.ParseAndValidateID(ctx, "user")
	if !ok {
		return
	}

	user, err := h.userSvc.GetUserByID(id)
	if err != nil {
		utils.NotFound(ctx, err.Error())
		return
	}

	utils.Success(ctx, "Lấy thông tin user thành công", user)
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	var req model.CreateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Dữ liệu không đúng định dạng JSON")
		return
	}

	user, err := h.userSvc.Create(req)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.Created(ctx, "Tạo user thành công", user)
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	id, ok := common.ParseAndValidateID(ctx, "user")
	if !ok {
		return
	}

	var req model.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Dữ liệu không đúng định dạng JSON")
		return
	}

	user, err := h.userSvc.UpdateUser(id, req)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.Success(ctx, "Cập nhật user thành công", user)
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	id, ok := common.ParseAndValidateID(ctx, "user")
	if !ok {
		return
	}

	if err := h.userSvc.DeleteUser(id); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.Success(ctx, "Xoá user thành công", nil)
}

func (h *UserHandler) GetUsersWithoutEmployee(ctx *gin.Context) {
	users, err := h.userSvc.GetUsersWithoutEmployee()
	if err != nil {
		utils.InternalServerError(ctx, err.Error())
		return
	}
	utils.Success(ctx, "Lấy danh sách user chưa gắn nhân viên thành công", users)
}
