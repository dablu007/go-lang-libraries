package main

import (
	"flag"
	"fmt"
	"os"

	"flow/auth"
	"flow/config"
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
	//logger.InitLogger()
	//db.Init(config.GetConfig())
	auth.Init()
	flag.Parse()
	server.Init()
}
