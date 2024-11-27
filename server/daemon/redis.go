package daemon

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func clearRedis(rdb *redis.Client) {
	ctx := context.Background()
	if err := rdb.FlushDB(ctx).Err(); err != nil {
		log.Printf("Failed to clear Redis database: %v\n", err)
	} else {
		log.Println("Redis database cleared")
	}
}

func initUserData(rdb *redis.Client) {
	ctx := context.Background()

	// Redis hash
	userStorageKey := "user_storage"

	// users_storage
	users := map[string]map[string]string{
		"a": {"password": "aaa"},
		"b": {"password": "bbb"},
		"c": {"password": "ccc"},
		"d": {"password": "ddd"},
	}

	// add to Redis
	for user, fields := range users {
		for key, value := range fields {
			fullKey := user + "." + key
			if err := rdb.HSet(ctx, userStorageKey, fullKey, value).Err(); err != nil {
				log.Printf("Failed to set user %s data: %v\n", user, err)
			} else {
				log.Printf("Set Redis key: %s, value: %s\n", fullKey, value)
			}
		}
	}
}
