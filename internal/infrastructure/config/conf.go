package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Load reads and loads the config inside a structure
func Load() (*Conf, error) {
	// Load default
	viper.SetDefault("slug.maximal-lenght", 8)
	viper.SetDefault("slug.time-to-expire", 7*24*time.Hour) // One week

	// Load from config file
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
	Database     PSQLConnConfig     `mapstructure:"database"`
	ServerDomain ServerDomainConfig `mapstructure:"server-domain"`
	Slug         SlugConfig         `mapstructure:"slug"`
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

// ServerDomainConfig represents the configuration of the server domain
type ServerDomainConfig struct {
	Scheme string `mapstructure:"scheme"`
	Domain string `mapstructure:"domain"`
	Port   int    `mapstructure:"port"`
}

// CreateBaseURL creates the base URL of the server
func (c *ServerDomainConfig) CreateBaseURL() string {
	baseURL := fmt.Sprintf("%s://%s", c.Scheme, c.Domain)
	if c.Port != 0 {
		baseURL = fmt.Sprintf("%s:%d", baseURL, c.Port)
	}
	return baseURL
}

// SlugConfig represents the configuration of the slug
type SlugConfig struct {
	MaximalLenght int           `mapstructure:"maximal-lenght"`
	TimeToExpire  time.Duration `mapstructure:"time-to-expire"`
}
