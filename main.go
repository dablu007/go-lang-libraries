package main

import (
	"flag"
	"fmt"
	"os"

	"flow/auth"
	"flow/cache"
	"flow/config"
	"flow/db"
	"flow/logger"
	"flow/server"
)

func main() {
	service := "flow"
	environment := os.Getenv("BOOT_CUR_ENV")
	if environment == "" {
		environment = "dev"
	}
	flag.Usage = func() {
		fmt.Println("Usage: server -s {service_name} -e {environment}")
		os.Exit(1)
	}
	flag.Parse()

	config.Init(service, environment)
	logger.InitLogger()
	auth.Init()
	flag.Parse()
	db.Init()
	cache.Init()
	server.Init()
}
