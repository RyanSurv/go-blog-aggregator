package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg *Config) SetUser(username string) {
	cfg.CurrentUserName = username
	
	if err := write(*cfg); err != nil {
		fmt.Printf("Error writing to file path: %s", err)
		os.Exit(1)
	}
}

func (cfg *Config) PrettyPrint() {
	fmt.Printf("{\n  \"db_url\": \"%s\",\n  \"current_user_name\": \"%s\"\n}\n", cfg.DbUrl, cfg.CurrentUserName)
}

func Read() *Config {
	filePath, err := getConfigFilePath()
	if err != nil {
		fmt.Printf("Error getting config file path: %s", err)
		os.Exit(1)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading config file: %s", err)
		os.Exit(1)
	}

	var cfg Config
	if err := json.Unmarshal(content, &cfg); err != nil {
		fmt.Printf("Error unmarshal config file: %s", err)
		os.Exit(1)		
	}

	return &cfg
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", homeDir, configFileName), nil
}

func write(cfg Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	
	err = os.WriteFile(filePath, b, 0666)
	if err != nil {
		return err
	}

	return nil
}