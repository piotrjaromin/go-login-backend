package config

import (
	"os"
	"encoding/json"
	"fmt"
)

//Config contains configuration data for modules in this project
type Config struct {
	Mongo struct {
		Server string `json:"server"`
		Database string `json:"database"`
	} `json:"mongo"`
	FrontendURL string `json:"frontendUrl"`
	Fb struct {
		ClientID string `json:"clientId"`
		ClientSecret string `json:"clientSecret"` 
	} `josn:"fb"`
	Token struct{
		SiginKey string `json:"siginKey"`
	} `json:"tokens"`
	Email struct {
		AwsRegion string `json:"awsRegion"`
		ReplyAddr string `json:"replyAddr"`
	} `json:"email"`
}

//GetConfig creates Config struct and fills it fields
//from json found under configPath
func GetConfig(configPath string) Config {

	file, errFile := os.Open(configPath)
	if errFile != nil {
		fmt.Println("error while reading config file:", errFile)
	}

	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error while unmarshalling config file:", err)
	}

	return configuration
}

//GetEnvOrDefault reads environemnt variable and returns value
//if there is no environemnt variable present then def value is returned
func GetEnvOrDefault(key string, def string) string {
	fromEnv := os.Getenv(key);
	if len(fromEnv) == 0 {
		return def
	}
	return fromEnv
}
