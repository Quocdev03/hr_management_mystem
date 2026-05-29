package handler

import (
	"chiquoc_hocgolang/internal/service"
	"chiquoc_hocgolang/internal/utils"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashB service.DashboardService
}

func NewDashboardHandler(dashB service.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashB: dashB,
	}
}

func (dashB *DashboardHandler) GetStats(ctx *gin.Context) {
	result, err := dashB.dashB.GetStats()
	if err != nil {
		utils.InternalServerError(ctx, err.Error())
		return
	}

	utils.Success(ctx, "Lấy số liệu dashboard thành công", result)

}
