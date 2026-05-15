package common

import (
	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/package/response"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// --- Các hàm dùng chung cho tất cả handler ---

// ParseUintParam parse URL parameter thành uint
func ParseUintParam(c *gin.Context, param string) (uint, error) {
	val, err := strconv.ParseUint(c.Param(param), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}

// ParseAndValidateID parse URL param ":id" và validate > 0
// Nếu không hợp lệ, tự động trả response BadRequest và return 0, false
// Dùng chung cho tất cả endpoint có :id
func ParseAndValidateID(ctx *gin.Context, entityName string) (uint, bool) {
	id, err := ParseUintParam(ctx, "id")
	if err != nil {
		response.BadRequest(ctx, "ID "+entityName+" không hợp lệ, phải là số nguyên dương!")
		return 0, false
	}
	if id == 0 {
		response.BadRequest(ctx, "ID "+entityName+" phải lớn hơn 0!")
		return 0, false
	}
	return id, true
}

// NormalizePagination chuẩn hoá query phân trang:
// - page < 1 → 1
// - limit < 1 → 10, limit > 100 → 100
// - trim search
func NormalizePagination(query *model.PaginationQuery) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 10
	}
	if query.Limit > 100 {
		query.Limit = 100
	}
	query.Search = strings.TrimSpace(query.Search)
}
