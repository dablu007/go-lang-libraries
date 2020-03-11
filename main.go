package main

import (
	"boiler-plate/auth"
	"boiler-plate/cache"
	"boiler-plate/config"
	"boiler-plate/db"
	"boiler-plate/logger"
	"boiler-plate/server"
	"flag"
	"fmt"
	"os"
)

func main() {
	service := "boiler-plate"
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
