package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

func LoadConfig() error {
	// Set default values for configuration
	viper.SetDefault("db.host", "host")
	viper.SetDefault("db.port", "port")
	viper.SetDefault("db.name", "name")
	viper.SetDefault("db.username", "username")
	viper.SetDefault("db.password", "password")

	// Read configuration from YAML file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./pkg/utils")
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	return nil
}
