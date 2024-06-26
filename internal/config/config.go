package config

import "github.com/spf13/viper"

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type AppConfig struct {
	Server           ServerConfig `mapstructure:"server"`
	MongoDbUri       string       `mapstructure:"mongo_db_uri"`
	Environment      string       `mapstructure:"env"`
	Prefix           string       `mapstructure:"prefix"`
	AliasMaxSize     int          `mapstructure:"alias_max_size"`
	RequestPerMinute int          `mapstructure:"requests_per_minute"`
}

func ReadConfig() (AppConfig, error) {
	configReader := viper.New()
	var config AppConfig

	configReader.SetConfigName(".application")
	configReader.SetConfigType("json")
	configReader.AddConfigPath(".")

	err := configReader.ReadInConfig()
	if err != nil {
		return AppConfig{}, err
	}

	err = configReader.Unmarshal(&config)
	if err != nil {
		return AppConfig{}, err
	}

	return config, nil
}
