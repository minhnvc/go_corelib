package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"

	"github.com/minhnvc/go_corelib/utils"
)

var ctx = context.Background()
var rdb *redis.Client
var cdb *cache.Cache

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     utils.GetConfig("REDIS_URL"),
		Password: utils.GetConfig("REDIS_PASSWORD"),
		DB:       0, // use default DB,
		PoolSize: 20,
	})
	if rdb == nil {
		panic("Can't connect redis service")
	}
	cdb = cache.New(&cache.Options{
		Redis: rdb,
		// LocalCache: cache.NewTinyLFU(50000, time.Minute),
	})

	fmt.Println("Redis", "Redis connected")
}

func Set(key string, value interface{}, duration time.Duration) {
	err := cdb.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   duration,
	})
	if err != nil {
		fmt.Println("Redis", err)
	}
}
func Get(key string, obj interface{}) {
	cdb.Get(ctx, key, obj)
}
