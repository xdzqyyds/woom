package model

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

const (
	UserStorageKey      = "user_storage"
	UserOnlineStatusKey = "user_online_status"
	InvitationKey       = "invitation"
)

func ClearRedis(rdb *redis.Client) {
	ctx := context.Background()
	if err := rdb.FlushDB(ctx).Err(); err != nil {
		log.Printf("Failed to clear Redis database: %v\n", err)
	} else {
		log.Println("Redis database cleared")
	}
}

func InitUserData(rdb *redis.Client) {
	ctx := context.Background()

	// users storage with key as userId and value as password
	users := GenerateUsers()

	// add to Redis
	for user, password := range users {
		if err := rdb.HSet(ctx, UserStorageKey, user, password).Err(); err != nil {
			log.Printf("Failed to set user %s data: %v\n", user, err)
		}
	}

	// online status storage with key as userId and value as online status
	onlineStatus := GenerateOnlineStatus(rdb)

	// add to Redis
	for user, status := range onlineStatus {
		if err := rdb.HSet(ctx, UserOnlineStatusKey, user, status).Err(); err != nil {
			log.Printf("Failed to set user %s online status: %v\n", user, err)
		}
	}
}

func GenerateUsers() map[string]string {
	users := make(map[string]string)

	for i := 'a'; i <= 'z'; i++ {
		userId := string(i)
		password := userId + userId + userId
		users[userId] = password
	}

	for i := 'A'; i <= 'Z'; i++ {
		userId := string(i)
		password := userId + userId + userId
		users[userId] = password
	}

	return users
}

func GenerateOnlineStatus(rdb *redis.Client) map[string]string {
	ctx := context.Background()
	onlineStatus := make(map[string]string)

	// Get all users from user_storage
	users, err := rdb.HKeys(ctx, UserStorageKey).Result()
	if err != nil {
		log.Printf("Failed to get users from Redis: %v\n", err)
		return onlineStatus
	}

	// Initialize online status for each user
	for _, userId := range users {
		onlineStatus[userId] = "0"
	}

	return onlineStatus
}
