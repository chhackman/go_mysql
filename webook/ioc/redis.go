package ioc

import (
	"awesomeProject/webook/config"
	"github.com/redis/go-redis/v9"
)

func InitRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr,
		Password: "HzN8m%cr!Vve",
		DB:       12,
	})
	return redisClient
}
