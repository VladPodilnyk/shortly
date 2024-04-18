package main

import (
	"log"
	"os"
	"time"

	"golang.org/x/time/rate"
	"shortly.io/internal/app"
	"shortly.io/internal/config"
	"shortly.io/internal/storage"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	cfg, err := config.ReadConfig()
	if err != nil {
		logger.Fatal(err)
		return
	}

	version, err := app.GetVersion()
	if err != nil {
		logger.Fatal(err)
		return
	}

	data := &app.AppData{
		Config:      cfg,
		Version:     version,
		Logger:      logger,
		Storage:     storage.New(),
		RateLimiter: rate.NewLimiter(rate.Every(time.Minute), cfg.RequestPerMinute),
	}

	err = app.Serve(data)
	if err != nil {
		logger.Fatal(err)
	}
}
