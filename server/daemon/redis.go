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
	userOnlineStatusKey := "user_online_status"

	// users storage with key as userId and value as password
	users := generateUsers()

	// add to Redis
	for user, password := range users {
		if err := rdb.HSet(ctx, userStorageKey, user, password).Err(); err != nil {
			log.Printf("Failed to set user %s data: %v\n", user, err)
		}
	}

	// online status storage with key as userId and value as online status
	onlineStatus := generateOnlineStatus()

	// add to Redis
	for user, status := range onlineStatus {
		if err := rdb.HSet(ctx, userOnlineStatusKey, user, status).Err(); err != nil {
			log.Printf("Failed to set user %s online status: %v\n", user, err)
		}
	}

}

func generateUsers() map[string]string {
	users := make(map[string]string)

	// 生成从 a 到 z 的用户名和密码
	for i := 'a'; i <= 'z'; i++ {
		userId := string(i)
		password := userId + userId + userId
		users[userId] = password
	}

	// 生成从 A 到 Z 的用户名和密码
	for i := 'A'; i <= 'Z'; i++ {
		userId := string(i)
		password := userId + userId + userId
		users[userId] = password
	}

	return users
}

func generateOnlineStatus() map[string]bool {
	onlineStatus := make(map[string]bool)

	// 生成从 a 到 z 的在线状态
	for i := 'a'; i <= 'z'; i++ {
		userId := string(i)
		onlineStatus[userId] = false
	}

	// 生成从 A 到 Z 的在线状态
	for i := 'A'; i <= 'Z'; i++ {
		userId := string(i)
		onlineStatus[userId] = false
	}

	return onlineStatus
}
