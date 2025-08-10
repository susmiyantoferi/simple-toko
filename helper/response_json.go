package helper

import (
	"simple-toko/web"

	"github.com/gin-gonic/gin"
)

func ToResponseJson(ctx *gin.Context, code int, message string, data interface{}) {
	ctx.JSON(code, web.WebResponse{
		Code:    code,
		Message: message,
		Data:    data,
	},
	)
}
