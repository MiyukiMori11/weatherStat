package config

import (
	"time"

	"github.com/MiyukiMori11/weatherstat/internal/domain"
	"github.com/spf13/viper"
)

type Config struct {
	Gateway Gateway `yaml:"gateway"`
	Server  Server  `yaml:"server"`
}

type Server struct {
	ListenHost         string `yaml:"listenHost"`
	ListenPort         int    `yaml:"listenPort"`
	Scheme             string `yaml:"scheme"`
	ShutdownTimeoutSec int    `yaml:"shutdownTimeoutSec"`
}

type Gateway struct {
	Services []domain.Service `yaml:"services"`
}

func Load(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("local") //TODO: исправить
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func (s *Server) ShutdownTimeout() time.Duration {
	return time.Duration(s.ShutdownTimeoutSec) * time.Second
}
