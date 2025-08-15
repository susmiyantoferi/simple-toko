package handler

import "github.com/gin-gonic/gin"

type PaymentHandler interface {
	UploadPayment(ctx *gin.Context)
	UpdateStatus(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindByOrderId(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	Delete(ctx *gin.Context)
	PreviewImage(ctx *gin.Context)
}
