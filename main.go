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


	router := route.NewRouter(userHandler, inventoryHandler)

	port := os.Getenv("PORT_APP")
	log.Println("server run in port ", port)
	router.Run(port)

}
