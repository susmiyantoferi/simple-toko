package main

import (
	"log"
	"os"
	"simple-toko/config"
	"simple-toko/handler"
	"simple-toko/repository"
	"simple-toko/route"
	"simple-toko/service"

	"github.com/go-playground/validator/v10"
)

func main() {
	if err := os.MkdirAll("uploads/product/", os.ModePerm); err != nil {
		log.Fatal("failed make folder product")
	}

	if err := os.MkdirAll("uploads/payment/", os.ModePerm); err != nil {
		log.Fatal("failed make folder payment")
	}

	db := config.Database()
	validate := validator.New()

	userRepo := repository.NewUserRepositoryImpl(db)
	userService := service.NewUserServiceImpl(userRepo, validate)
	userHandler := handler.NewUserHandlerImpl(userService)

	inventoryRepo := repository.NewInventoryRepositoryImpl(db)
	inventoryService := service.NewInventoryServiceImpl(inventoryRepo, validate)
	inventoryHandler := handler.NewInventoryHandlerImpl(inventoryService)

	addresRepo := repository.NewAddressRepositoryImpl(db)
	addressService := service.NewAddressServiceImpl(addresRepo, userRepo, validate)
	addressHandler := handler.NewAddressHandlerImpl(addressService)

	productRepo := repository.NewProductRepositoryImpl(db)
	productService := service.NewProductServiceImpl(productRepo, inventoryRepo, validate)
	productHandler := handler.NewProductHandlerImpl(productService)

	orderRepo := repository.NewOrderRepositoryImpl(db)
	orderService := service.NewOrderServiceImpl(orderRepo, addresRepo, validate)
	orderHandler := handler.NewOrderHandlerImpl(orderService)

	payRepo := repository.NewPaymentRepositoryImpl(db)
	payService := service.NewPaymentServiceImpl(payRepo, orderRepo, validate)
	payHandler := handler.NewPaymentHandlerImpl(payService)

	reportRepo := repository.NewReportRepositoryImpl(db)
	reportService := service.NewReportServiceImpl(reportRepo)
	reportHndler := handler.NewReportHandlerImpl(reportService)

	router := route.NewRouter(
		userHandler,
		inventoryHandler,
		addressHandler,
		productHandler,
		orderHandler,
		payHandler,
		reportHndler,
	)

	port := os.Getenv("PORT_APP")
	log.Println("server run in port ", port)
	router.Run(port)

}
