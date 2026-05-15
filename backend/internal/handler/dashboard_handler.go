package handler

import (
	"chiquoc_hocgolang/internal/service"
	"chiquoc_hocgolang/internal/utils"

	"github.com/gin-gonic/gin"
)

type DashboardsHanlder struct {
	dashB service.DashboardService
}

func NewDashboardHanlder(dashB service.DashboardService) *DashboardsHanlder {
	return &DashboardsHanlder{
		dashB: dashB,
	}
}

func (dashB *DashboardsHanlder) GetStats(ctx *gin.Context) {
	result, err := dashB.dashB.GetStats()
	if err != nil {
		utils.InternalServerError(ctx, err.Error())
	}

	utils.Success(ctx, "Lấy số liệu dashboard thành công", result)

}
