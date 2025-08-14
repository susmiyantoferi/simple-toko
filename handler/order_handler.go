package handler

import "github.com/gin-gonic/gin"

type OrderHandler interface {
	CreateOrder(ctx *gin.Context)
	UpdateAddress(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	ConfirmOrder(ctx *gin.Context)
}
