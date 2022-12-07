package config

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	// Timeout for parser's iteration launch in minutes
	parseTimeoutMin = "PARSE_TIMEOUT_MIN"

	clientTimeoutSec = "CLIENT_TIMEOUT_SEC"

	weatherApiKey = "WEATHER_API_KEY"

	// DB global settings
	dbName     = "DB_NAME"
	dbUsername = "DB_USERNAME"
	dbPassword = "DB_PASSWORD"
	dbHost     = "DB_HOST"
	dbPort     = "DB_PORT"
)

type Config struct {
	Parser  *Parser
	Client  *Client `yaml:"client"`
	Storage *Storage
}

type Storage struct {
	username string
	password string
	host     string
	port     string
	name     string
}

type Client struct {
	GeoAPI     string `yaml:"geoAPI"`
	WeatherAPI string `yaml:"weatherAPI"`

	apiKey  string
	timeout int
}

type Parser struct {
	parseTimeoutMin int
}

func Load(fullPath string) (config *Config, err error) {
	path, name := filepath.Split(fullPath)
	name, format, _ := strings.Cut(name, ".")

	return load(path, name, format)
}

func load(path, name, format string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType(format)

	viper.BindEnv(parseTimeoutMin)

	viper.BindEnv(clientTimeoutSec)

	viper.BindEnv(weatherApiKey)

	viper.BindEnv(dbName, dbUsername, dbPassword, dbHost, dbPort)

	err = viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("can't read in config: %w", err)
	}

	err = viper.Unmarshal(&config)

	config.Parser = &Parser{
		parseTimeoutMin: viper.GetInt(parseTimeoutMin),
	}

	config.Client.apiKey = viper.GetString(weatherApiKey)
	config.Client.timeout = viper.GetInt(clientTimeoutSec)

	config.Storage = &Storage{
		username: viper.GetString(dbUsername),
		password: viper.GetString(dbPassword),
		host:     viper.GetString(dbHost),
		port:     viper.GetString(dbPort),
		name:     viper.GetString(dbName),
	}

	return
}

func (p *Parser) Timeout() time.Duration {
	return time.Duration(p.parseTimeoutMin) * time.Minute
}

func (c *Client) Timeout() time.Duration {
	return time.Duration(c.timeout) * time.Second
}

func (c *Client) APIKey() string {
	return c.apiKey
}

func (s *Storage) URL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", s.username, s.password, s.host, s.port, s.name)
}
