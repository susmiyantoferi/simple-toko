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
	orderService := service.NewOrderServiceImpl(orderRepo, validate)
	orderHandler := handler.NewOrderHandlerImpl(orderService)

	router := route.NewRouter(userHandler, inventoryHandler, addressHandler, productHandler, orderHandler)

	port := os.Getenv("PORT_APP")
	log.Println("server run in port ", port)
	router.Run(port)

}
