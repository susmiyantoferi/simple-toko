package config

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	err := godotenv.Load()
	if err != nil{
		log.Fatal("error load env")
	}

	rdsAddr := os.Getenv("REDIS_ADDRS")
	rdsPwd := os.Getenv("REDIS_PWD")

	Redis := redis.NewClient(&redis.Options{
		Addr:     rdsAddr,
		Password: rdsPwd,
		DB:       0,
	})

	_, err = Redis.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("failed to connect redis %v", err)
	}

	return Redis

}
