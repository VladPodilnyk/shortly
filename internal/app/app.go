package app

import (
	"log"

	"github.com/spf13/viper"
	"golang.org/x/time/rate"

	"shortly.io/internal/config"
	"shortly.io/internal/storage"
)

type AppData struct {
	Config      config.AppConfig
	Version     string
	Logger      *log.Logger
	Storage     storage.Storage // app persistance
	RateLimiter *rate.Limiter   // application rate limiter;
}

func GetVersion() (string, error) {
	versionReader := viper.New()
	var appVersion struct {
		Version string `mapstructure:"version"`
	}

	versionReader.SetConfigName("version")
	versionReader.SetConfigType("json")
	versionReader.AddConfigPath(".")

	err := versionReader.ReadInConfig()
	if err != nil {
		return "", err
	}

	err = versionReader.Unmarshal(&appVersion)
	if err != nil {
		return "", err
	}

	return appVersion.Version, nil
}
