package handler

import "github.com/gin-gonic/gin"

type AddressHandler interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindByUserId(ctx *gin.Context)
	FindAll(ctx *gin.Context)
}
