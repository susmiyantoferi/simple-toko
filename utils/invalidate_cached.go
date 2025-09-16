package utils

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func InvalidateCached(ctx context.Context, rds *redis.Client, id uint) {
	keyId := fmt.Sprintf("products:%d", id)
	if err := rds.Del(ctx, keyId).Err(); err != nil {
		fmt.Printf("failed delete cache find by id on key %s: %v\n", keyId, err)
	}

	i := rds.Scan(ctx, 0, "products:page*", 0).Iterator()
	for i.Next(ctx) {
		keys := i.Val()
		if err := rds.Del(ctx, keys).Err(); err != nil {
			fmt.Printf("failed delete cache find all on key %s: %v\n", keys, err)
		}
	}

	if err := i.Err(); err != nil {
		fmt.Printf("iterator err: %v\n", err)
	}
}
