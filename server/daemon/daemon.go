package daemon

import (
	"context"
	"log"
	"net/http"
	"woom/server"
	"woom/server/api"
	"woom/server/model"

	"github.com/caarlos0/env/v9"
	"github.com/redis/go-redis/v9"
)

func newRdbClient(url string) *redis.Client {
	opts, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opts)

	// Configure RDB and AOF backup
	ctx := context.Background()
	if err := client.ConfigSet(ctx, "save", "900 1 300 10 60 10000").Err(); err != nil {
		log.Printf("Failed to set RDB: %v\n", err)
	}
	if err := client.ConfigSet(ctx, "appendonly", "yes").Err(); err != nil {
		log.Printf("Failed to set AOF: %v\n", err)
	}
	if err := client.ConfigSet(ctx, "appendfsync", "everysec").Err(); err != nil {
		log.Printf("Failed to set AOF sync: %v\n", err)
	}

	return client
}

func Daemon(ctx context.Context) {
	cfg := server.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}

	log.Printf("%+v\n", cfg)
	rdb := newRdbClient(cfg.RedisUrl)
	log.Println(rdb)

	//model.ClearRedis(rdb)
	model.InitUserData(rdb)

	handler := api.NewApi(rdb, cfg.Secret, cfg.Live777Url, cfg.Live777Token)

	log.Println("=== started ===")
	log.Panicln(http.ListenAndServe(":"+cfg.Port, handler))
}
