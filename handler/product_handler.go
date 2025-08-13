package handler

import "github.com/gin-gonic/gin"

type ProductHandler interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	AddStock(ctx *gin.Context)
	ReduceStock(ctx *gin.Context)
	UpdateImage(ctx *gin.Context)
	PreviewImage(ctx *gin.Context)
}
