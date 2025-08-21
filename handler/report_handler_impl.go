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
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "5")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err:= strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1{
		pageSize = 5
	}

	result, err := r.ReportService.MonthlySales(ctx, page, pageSize)
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
