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

	// Lấy thông tin user bị tác động để kiểm tra phân quyền
	targetUser, err := h.userSvc.GetUserByID(id)
	if err != nil {
		utils.NotFound(ctx, "Không tìm thấy user")
		return
	}

	var req model.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Dữ liệu không đúng định dạng JSON")
		return
	}

	// Logic kiểm tra phân quyền an toàn
	requesterID, exists := ctx.Get("userID")
	if exists {
		reqID := requesterID.(uint)
		
		// 1. Ngăn tự đổi quyền bản thân HOẶC đổi quyền của Admin khác
		if req.RoleID != nil && *req.RoleID != targetUser.RoleID {
			if targetUser.RoleID == 1 || targetUser.ID == reqID {
				utils.Forbidden(ctx, "Không thể tự thay đổi quyền của mình hoặc của Admin khác")
				return
			}
		}

		// 2. Ngăn tự khoá tài khoản của chính mình
		if targetUser.ID == reqID && req.IsActive != nil && !*req.IsActive {
			utils.Forbidden(ctx, "Không thể tự vô hiệu hoá tài khoản của chính mình")
			return
		}
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
