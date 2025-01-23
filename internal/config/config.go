package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	GatorConfig = ".gatorconfig.json"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getGatorConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error finding homedir: %v", err)
	}

	return filepath.Join(homeDir, GatorConfig), nil
}

func Read() (Config, error) {
	var config Config

	configPath, err := getGatorConfigPath()
	if err != nil {
		return config, fmt.Errorf("error fetching config path: %v", err)
	}

	jsonData, err := os.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("error reading gatorconfig: %v", err)
	}

	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		return config, fmt.Errorf("error parsing json config: %v", err)
	}
	return config, nil
}

func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = userName
	filePath, err := getGatorConfigPath()
	if err != nil {
		return fmt.Errorf("error getting gator config path: %v", err)
	}

	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return fmt.Errorf("error making json config: %v", err)
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("eror writing config file: %v", err)
	}

	return nil
}
