package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Load reads and loads the config inside a structure
func Load() (*Conf, error) {
	viper.SetConfigName(os.Getenv("env"))
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../../conf/") // Local and test
	viper.AddConfigPath("./conf/")        // Docker
	err := viper.ReadInConfig()
	if err != nil {
		return &Conf{}, fmt.Errorf("failed to read config: %w", err)
	}

	var config Conf
	err = viper.Unmarshal(&config)
	if err != nil {
		return &config, fmt.Errorf("failed to unmarshall config: %w", err)
	}

	return &config, nil
}

// Conf represents the configuration of the application
type Conf struct {
	Database PSQLConnConfig `mapstructure:"database"`
}

// PSQLConnConfig represents the configuration to connect to a PSQL database
type PSQLConnConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DbName   string `mapstructure:"dbname"`
}

// ToConnString generates a conn string based on the conn config
func (c *PSQLConnConfig) ToConnString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.DbName)
}
