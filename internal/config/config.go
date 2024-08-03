package config

import (
	"errors"
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type StorageConfig struct {
	MongoDbUri        string `mapstructure:"uri"`
	MongoDbName       string `mapstructure:"name"`
	MongoDbCollection string `mapstructure:"collection"`
}

type AppConfig struct {
	Server           ServerConfig  `mapstructure:"server"`
	Storage          StorageConfig `mapstructure:"storage"`
	Environment      string        `mapstructure:"env"`
	Prefix           string        `mapstructure:"prefix"`
	AliasMaxSize     int           `mapstructure:"alias_max_size"`
	RequestPerMinute int           `mapstructure:"requests_per_minute"`
}

func ReadConfig() (AppConfig, error) {
	configName, err := pickConfig()
	if err != nil {
		return AppConfig{}, err
	}

	configReader := viper.New()
	var config AppConfig

	configReader.SetConfigName(configName)
	configReader.SetConfigType("json")
	configReader.AddConfigPath(".")

	err = configReader.ReadInConfig()
	if err != nil {
		return AppConfig{}, err
	}

	err = configReader.Unmarshal(&config)
	if err != nil {
		return AppConfig{}, err
	}

	return config, nil
}

func pickConfig() (string, error) {
	env := flag.String("env", "dev", "Environment (e.g., dev, prod)")
	flag.Parse()
	fmt.Println("Environment:", *env)
	switch *env {
	case "dev":
		return ".application.dev", nil
	case "prod":
		return ".application", nil
	default:
		return "", errors.New("unknown environment configuration. Available options: dev, prod")
	}
}
