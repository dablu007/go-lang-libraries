package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/buger/jsonparser"
	"github.com/spf13/viper"
)

var config *viper.Viper

func Init(service, env string) {
	var bodyBytes []byte
	configFile, err := http.Get("http://configurations.zestmoney.in:8888/" + service + "/" + env)
	if err != nil {
		log.Print("Error fetching configuration from server for service : " + service + " env : " + env)
		bodyBytes, err = ioutil.ReadFile("config/config.json")
		if err != nil {
			log.Fatal("Couldn't read local configuration file.", err)
		} else {
			log.Print("using local config.")
		}
	} else {
		if configFile != nil {
			bodyBytes, err = ioutil.ReadAll(configFile.Body)
			if err != nil {
				log.Fatal("Error reading configuration response body.")
			}
		}
	}

	config = viper.New()
	config.SetConfigType("json")
	config.SetConfigName(env)
	fmt.Print("body:", string(bodyBytes))
	parsedConfig, _, _, parseErr := jsonparser.Get(bodyBytes, "propertySources", "[0]", "source")
	if parseErr != nil {
		log.Fatal("Failed to parse config JSON: ", parseErr)
	}

	//err = config.ReadConfig(bytes.NewBuffer(parsedConfig))
	err = config.ReadConfig(bytes.NewReader(parsedConfig))
	if err != nil {
		log.Fatal("Failed to reading config: ", err)
	}
}

func GetConfig() *viper.Viper {
	return config
}
