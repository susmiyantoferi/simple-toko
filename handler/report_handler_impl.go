package handler

import (
	"net/http"
	"simple-toko/helper"
	"simple-toko/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type reportHandlerImpl struct {
	ReportService service.ReportService
}

func NewReportHandlerImpl(reportService service.ReportService) *reportHandlerImpl {
	return &reportHandlerImpl{
		ReportService: reportService,
	}
}

func (r *reportHandlerImpl) MonthlySales(ctx *gin.Context) {
	result, err := r.ReportService.MonthlySales(ctx)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	helper.ToResponseJson(ctx, http.StatusOK, "success", result)
}

func (r *reportHandlerImpl) TopProductSales(ctx *gin.Context) {
	lim := ctx.Param("limit")
	limitParam := ctx.DefaultQuery("limit", lim)

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input type limit", nil)
		return
	}

	result, err := r.ReportService.TopProductSales(ctx, limit)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	helper.ToResponseJson(ctx, http.StatusOK, "success", result)

}

func (r *reportHandlerImpl) LessProductSales(ctx *gin.Context) {
	lim := ctx.Param("limit")
	limitParam := ctx.DefaultQuery("limit", lim)

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input type limit", nil)
		return
	}

	result, err := r.ReportService.LessProductSales(ctx, limit)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	helper.ToResponseJson(ctx, http.StatusOK, "success", result)
}
