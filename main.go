package main

import (
	"flag"
	"flow/auth"
	"fmt"
	"os"
	"flow/db"
	"flow/cache"
	"flow/config"
	"flow/logger"
	"flow/server"
)

func main() {
	service := "flow"
	environment := os.Getenv("BOOT_CUR_ENV")
	if environment == "" {
		environment = "test"
	}
	flag.Usage = func() {
		fmt.Println("Usage: server -s {service_name} -e {environment}")
		os.Exit(1)
	}
	flag.Parse()
	config.Init(service, environment)
	logger.InitLogger()
	flag.Parse()
	auth.Init()
	db.Init()
	cache.Init()
	server.Init()
}
