package main

import "fmt"

// import "github.com/aws/aws-sdk-go/service/elasticache"
import "github.com/go-redis/redis/v7"

var client *redis.Client

func ExampleNewClient() {
	client = redis.NewClient(&redis.Options{
		//Addr: "127.0.0.1:6379",
		Addr: "chache-poc-golang.wpfjkv.ng.0001.aps1.cache.amazonaws.com:6379",
		DB: 0, // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
}

func ExampleClient() {
	err := client.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get("key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := client.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}

func main() {
	print("hello\n")
	ExampleNewClient()
	ExampleClient()
}