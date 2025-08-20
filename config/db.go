package config

import (
	"fmt"
	"log"
	"os"
	"simple-toko/entity"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	// "github.com/joho/godotenv"
)

func Database() *gorm.DB{
	// err := godotenv.Load()
	// if err != nil{
	// 	log.Fatal("error load env")
	// }

	user := os.Getenv("DB_USERNAME")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbName)

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(
		&entity.User{},
		&entity.Address{},
		&entity.Inventory{},
		&entity.Product{}, 
		&entity.Order{}, 
		&entity.OrderProduct{},
		&entity.Payment{}, 
	)
	if err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	return db
}