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

	// users storage with key as username and value as password
	users := map[string]string{
		"a": "aaa",
		"b": "bbb",
		"c": "ccc",
		"d": "ddd",
	}

	// add to Redis
	for user, password := range users {
		if err := rdb.HSet(ctx, userStorageKey, user, password).Err(); err != nil {
			log.Printf("Failed to set user %s data: %v\n", user, err)
		}
	}
}
