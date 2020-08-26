package main

import (
	"flag"
	"fmt"
	"github.com/dablu007/go-lang-libraries/auth"
	"github.com/dablu007/go-lang-libraries/cache"
	"github.com/dablu007/go-lang-libraries/config"
	"github.com/dablu007/go-lang-libraries/db"
	"github.com/dablu007/go-lang-libraries/logger"
	"github.com/dablu007/go-lang-libraries/server"
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
	/* Remove the commented code the set the properties
		inside the config.json file */
	db.Init()
	cache.Init()
	server.Init()
}
