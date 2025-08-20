package route

import (
	"simple-toko/handler"
	"simple-toko/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	UserHandler handler.UserHandler,
	InventHandler handler.InventoryHandler,
	AddressHandler handler.AddressHandler,
	ProductHandler handler.ProductHandler,
	OrderHandler handler.OrderHandler,
	PaymentHandler handler.PaymentHandler,
	ReportHandler handler.ReportHandler,
) *gin.Engine {
	router := gin.Default()

	regist := router.Group("/api/v1/")
	{
		regist.GET("product", ProductHandler.FindAll)
		regist.POST("register", UserHandler.Create)
		regist.POST("login", UserHandler.Login)
		regist.POST("refresh-token", UserHandler.RefreshToken)
	}

	api := router.Group("/api/v1")
	api.Use(middleware.Authentication())
	{
		admin := api.Group("/")
		admin.Use(middleware.RoleAccessMiddleware("admin"))
		{
			//user

			admin.DELETE("users/:userId", UserHandler.Delete)
			admin.GET("users/id/:userId", UserHandler.FindById)
			admin.GET("users/email/:email", UserHandler.FindByEmail)
			regist.POST("users/admin", UserHandler.CreateAdmin)

			//address
			admin.GET("address", AddressHandler.FindAll)

			//inventory
			admin.POST("inventory", InventHandler.Create)
			admin.PUT("inventory/:invId", InventHandler.Update)
			admin.DELETE("inventory/:invId", InventHandler.Delete)
			admin.GET("inventory/:invId", InventHandler.FindById)
			admin.GET("inventory", InventHandler.FindAll)

			//product
			admin.POST("product", ProductHandler.Create)
			admin.PUT("product/:productId", ProductHandler.Update)
			admin.DELETE("product/:productId", ProductHandler.Delete)
			admin.GET("product/:productId", ProductHandler.FindById)
			//admin.GET("product", ProductHandler.FindAll)
			admin.PUT("product/:productId/add", ProductHandler.AddStock)
			admin.PUT("product/:productId/reduce", ProductHandler.ReduceStock)
			admin.PUT("product/image/:productId", ProductHandler.UpdateImage)
			admin.GET("product/image/:productId", ProductHandler.PreviewImage)

			//orders
			admin.GET("order", OrderHandler.FindAll)
			admin.PUT("order/confirm/:id", OrderHandler.ConfirmOrder)
			admin.DELETE("order/:id", OrderHandler.Delete)

			//payments
			admin.GET("payment", PaymentHandler.FindAll)
			admin.DELETE("payment/:id", PaymentHandler.Delete)
			admin.GET("payment/image/:orderId", PaymentHandler.PreviewImage)
			admin.GET("payment/:id", PaymentHandler.FindById)
			admin.PUT("payment/status/:orderId", PaymentHandler.UpdateStatus)

			//reports
			admin.GET("monthly-sales", ReportHandler.MonthlySales)
			admin.GET("top-product", ReportHandler.TopProductSales)
			admin.GET("less-product", ReportHandler.LessProductSales)
		}

		cust := api.Group("/")
		cust.Use(middleware.RoleAccessMiddleware("customer", "admin"))
		{
			cust.PUT("users/:userId", UserHandler.Update)
			cust.GET("users", UserHandler.FindAll)

			cust.POST("address", AddressHandler.Create)
			cust.PUT("address/:id", AddressHandler.Update)
			cust.DELETE("address/:id", AddressHandler.Delete)
			cust.GET("address/user/:userId", AddressHandler.FindByUserId)

			cust.POST("order", OrderHandler.CreateOrder)
			cust.PUT("order/:id", OrderHandler.UpdateAddress)
			cust.GET("order/:id", OrderHandler.FindById)

			cust.POST("payment", PaymentHandler.UploadPayment)
			cust.GET("payment/order/:orderId", PaymentHandler.FindByOrderId)
		}

	}

	return router
}
