package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read(configFileName string) (Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	configFilePath := homeDir + "/" + configFileName
	fmt.Println(configFilePath)
	configData, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, err
	}

	var configStruct Config
	err = json.Unmarshal(configData, &configStruct)
	if err != nil {
		return Config{}, err
	}

	return configStruct, nil
}
