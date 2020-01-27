package cache

import (
	"fmt"
	"github.com/spf13/viper"
)

import "github.com/go-redis/redis/v7"

var client *redis.Client

func Init() {
	redisUrl := viper.GetString("redis.url")
	client = redis.NewClient(&redis.Options{
		Addr: redisUrl,
		DB:   0, // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
}

func GetRedisClient() *redis.Client {
	return client
}
