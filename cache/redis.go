package cache

import "fmt"

// import "github.com/aws/aws-sdk-go/service/elasticache"
import "github.com/go-redis/redis/v7"

var client *redis.Client

func Init() {
	client = redis.NewClient(&redis.Options{
		//Addr: "127.0.0.1:6379",
		Addr: "chache-poc-golang.wpfjkv.ng.0001.aps1.cache.amazonaws.com:6379",
		DB:   0, // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
}

func GetRedisClient() *redis.Client {
	return client
}
