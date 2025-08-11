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


	router := route.NewRouter(userHandler, inventoryHandler, addressHandler)

	port := os.Getenv("PORT_APP")
	log.Println("server run in port ", port)
	router.Run(port)

}
