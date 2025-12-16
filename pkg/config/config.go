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
	DB        DBConfig `mapstructure:"db"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type SnowflakeConfig struct {
	MachineID int64 `mapstructure:"machine_id"`
}

// DBConfig holds the database connection parameters.
type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	// Also read from environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set default values
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("snowflake.machine_id", 1)
	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.port", 5432)
	viper.SetDefault("db.user", "postgres")
	viper.SetDefault("db.password", "postgres")
	viper.SetDefault("db.name", "container-manager")

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
