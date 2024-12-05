package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c Config) String() string {
	stringer := "~~ CONFIG ~~\n"
	stringer += fmt.Sprintf("-- Username: %s\n", c.CurrentUserName)
	stringer += fmt.Sprintf("-- Database URL: %s\n", c.DbUrl)

	return stringer
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configFilePath := homeDir + "/" + configFileName

	return configFilePath, nil
}

func write(cfg Config) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	// marshal Config data to write
	configData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFilePath, configData, 0666)
	if err != nil {
		return err
	}

	return nil
}

func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	// read from ~/gatorconfig.json
	configData, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, err
	}

	// unmarshal into Config struct
	var configStruct Config
	err = json.Unmarshal(configData, &configStruct)
	if err != nil {
		return Config{}, err
	}

	return configStruct, nil
}

func (c *Config) SetUser(newUserName string) error {
	c.CurrentUserName = newUserName
	err := write(*c)
	if err != nil {
		return err
	}

	return nil
}
