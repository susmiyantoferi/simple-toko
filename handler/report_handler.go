package handler

import "github.com/gin-gonic/gin"

type ReportHandler interface {
	MonthlySales(ctx *gin.Context)
	TopProductSales(ctx *gin.Context)
	LessProductSales(ctx *gin.Context)
}
