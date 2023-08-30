package redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/minhnvc/go_corelib/utils"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client
var cdb *cache.Cache

func InitRedis() {
	//init redis
	clusterUrl := strings.Split(utils.GetConfig("REDIS_URL"), ",")
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:           clusterUrl,
		PoolSize:        20,
		MaxRetries:      2,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
		DialTimeout:     3 * time.Second,                   // 3000 milliseconds, timeout config
		Password:        utils.GetConfig("REDIS_PASSWORD"), // no username needed as per your config
	})

	if rdb == nil {
		panic("Can't connect redis service")
	}
	cdb = cache.New(&cache.Options{
		Redis: rdb,
	})

	fmt.Println("Redis", "Redis cluster connected")
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
