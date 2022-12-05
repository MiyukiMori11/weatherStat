package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Gateway Gateway `yaml:"gateway"`
	Server  Server  `yaml:"server"`
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

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("can't read in config: %w", err)
	}

	err = viper.Unmarshal(&config)
	return
}
