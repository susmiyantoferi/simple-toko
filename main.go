package main

import (
	"net/http"
	"simple-toko/config"

	"github.com/gin-gonic/gin"
)

func main(){
	config.Database()
	
	r := gin.Default()
  r.GET("/ping", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  })
  r.Run()
}