package config

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
)

var config *viper.Viper

func Init(service, env string) {
	fmt.Println("Loading config from %s\n")
	body, err := fetchConfiguration()
	if err != nil {
		fmt.Println("Couldn't load configuration, cannot start. Terminating. Error: " + err.Error())
	}
	parseConfiguration(body)
}

// Make HTTP request to fetch configuration from config server
func fetchConfiguration() ([]byte, error) {
	var bodyBytes []byte
		//panic("Couldn't load configuration, cannot start. Terminating. Error: " + err.Error())
	var err error
	bodyBytes, err = ioutil.ReadFile("config/config.json")
	if err != nil {
		fmt.Println("Couldn't read local configuration file.", err)
	} else {
		log.Print("using local config.")
	}
	return bodyBytes, err
}

// Pass JSON bytes into struct and then into Viper
func parseConfiguration(body []byte) {
	var cloudConfig springCloudConfig
	err := json.Unmarshal(body, &cloudConfig)
	if err != nil {
		fmt.Println("Cannot parse configuration, message: " + err.Error())
	}
	for key, value := range cloudConfig.PropertySources[0].Source {
		viper.Set(key, value)
		fmt.Println("Loading config property %v => %v\n", key, value)
	}
	if viper.IsSet("server_name") {
		fmt.Println("Successfully loaded configuration for service %s\n", viper.GetString("server_name"))
	}
}

// Structs having same structure as response from Spring Cloud Config
type springCloudConfig struct {
	Name            string           `json:"name"`
	Profiles        []string         `json:"profiles"`
	Label           string           `json:"label"`
	Version         string           `json:"version"`
	PropertySources []propertySource `json:"propertySources"`
}
type propertySource struct {
	Name   string                 `json:"name"`
	Source map[string]interface{} `json:"source"`
}
