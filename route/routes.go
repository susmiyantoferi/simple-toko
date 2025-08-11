package route

import (
	"simple-toko/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	UserHandler handler.UserHandler,
	InventHandler handler.InventoryHandler,
	AddressHandler handler.AddressHandler,
) *gin.Engine {
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

		api.POST("inventory", InventHandler.Create)
		api.PUT("inventory/:invId", InventHandler.Update)
		api.DELETE("inventory/:invId", InventHandler.Delete)
		api.GET("inventory/:invId", InventHandler.FindById)
		api.GET("inventory", InventHandler.FindAll)

		api.POST("address", AddressHandler.Create)
		api.PUT("address/:id/user/:userId", AddressHandler.Update)
		api.DELETE("address/:id", AddressHandler.Delete)
		api.GET("address/user/:userId", AddressHandler.FindByUserId)
		api.GET("address", AddressHandler.FindAll)

	}

	return router
}
