package route

import (
	"simple-toko/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(UserHandler handler.UserHandler) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/")
	{
		api.POST("users", UserHandler.Create)
		api.POST("users/admin", UserHandler.CreateAdmin)
		api.PUT("users/:userId", UserHandler.Update)
		api.DELETE("users/:userId", UserHandler.Delete)
		api.GET("users/id/:userId", UserHandler.FindById)
		api.GET("users/email/:email", UserHandler.FindByEmail)
		api.GET("users", UserHandler.FindAll)

	}

	return router
}
