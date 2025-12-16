package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config stores all configuration for the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	Server    ServerConfig
	Snowflake SnowflakeConfig
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type SnowflakeConfig struct {
	MachineID int64 `mapstructure:"machine_id"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	// Also read from environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Find and read the config file.
	// Ignore error if config file is not found, as we can rely on env vars and defaults.
	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}
