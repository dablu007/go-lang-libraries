package main

import (
	"flag"
	"fmt"
	"go-lang/libraries/auth"
	"go-lang/libraries/cache"
	"go-lang/libraries/config"
	"go-lang/libraries/db"
	"go-lang/libraries/logger"
	"go-lang/libraries/server"
	"os"
)

func main() {
	service := "go-lang-libraries"
	environment := os.Getenv("BOOT_CUR_ENV")
	if environment == "" {
		environment = "test"
	}
	flag.Usage = func() {
		fmt.Println("Usage: server -s {service_name} -e {environment}")
		os.Exit(1)
	}
	flag.Parse()
	configUrl := "" // Put the configuration url of spring cloud config
	config.Init(configUrl, service, environment)
	logger.InitLogger()
	flag.Parse()
	auth.Init()
	db.Init()
	cache.Init()
	server.Init()
}
