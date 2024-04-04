package main

import (
	"log"
	"os"
	"time"

	"golang.org/x/time/rate"

	"github.com/spf13/viper"
	"shortly.io/internal/storage"
)

type application struct {
	config  AppConfig
	version string
	logger  *log.Logger
	storage storage.Storage // app persistance
	// algorithm   encoder.UrlEncoder // url shortener algorithm
	rateLimiter *rate.Limiter // application rate limiter;
}

func getVersion() (string, error) {
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

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	cfg, err := readConfig()
	if err != nil {
		logger.Fatal(err)
		return
	}

	version, err := getVersion()
	if err != nil {
		logger.Fatal(err)
		return
	}

	// algorithm := encoder.MakeSimple(cfg.Prefix, 0)
	rateLimiter := rate.NewLimiter(rate.Every(time.Minute), cfg.RequestPerMinute)
	app := &application{
		config:  cfg,
		version: version,
		logger:  logger,
		storage: storage.New(),
		// algorithm:   algorithm,
		rateLimiter: rateLimiter,
	}

	err = app.serve()
	if err != nil {
		logger.Fatal(err)
	}
}
