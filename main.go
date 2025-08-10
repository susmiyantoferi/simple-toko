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

func main(){
	db := config.Database()
  validate := validator.New()

	userRepo := repository.NewUserRepositoryImpl(db)
  userService := service.NewUserServiceImpl(userRepo, validate)
  userHandler := handler.NewUserHandlerImpl(userService)

  router := route.NewRouter(userHandler)

  port := os.Getenv("PORT_APP")
  log.Println("server run in port ", port)
  router.Run(port)

}