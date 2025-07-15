package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

// loadConfig reads configuration from a file and returns a Config struct.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // json, xml, yaml, toml, etc.

	viper.AutomaticEnv() // read environment variables that match
	
	// Set defaults for testing/CI environments
	setDefaultValues()
	
	err = viper.ReadInConfig()
	if err != nil {
		// If config file is not found, use environment variables and defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; use environment variables and defaults
			err = nil
		} else {
			// Config file was found but another error was produced
			return
		}
	}
	err = viper.Unmarshal(&config)
	return
}

// setDefaultValues sets default configuration values for testing/CI environments
func setDefaultValues() {
	viper.SetDefault("DB_DRIVER", "postgres")
	viper.SetDefault("DB_SOURCE", "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable")
	viper.SetDefault("SERVER_ADDRESS", "0.0.0.0:8080")
	viper.SetDefault("TOKEN_SYMMETRIC_KEY", "12345678901234567890123456789012")
	viper.SetDefault("ACCESS_TOKEN_DURATION", "15m")
}
