package app

import (
	"io/fs"
	"log"

	"golang.org/x/time/rate"

	"shortly.io/internal/config"
	"shortly.io/internal/storage"
)

type AppData struct {
	Config      config.AppConfig
	Logger      *log.Logger
	Storage     storage.Storage // app persistance
	RateLimiter *rate.Limiter   // application rate limiter;
	PublicFS    fs.FS
}
