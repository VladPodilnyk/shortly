package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/time/rate"
	"shortly.io/internal/app"
	"shortly.io/internal/config"
	"shortly.io/internal/storage"
)

//go:embed public
var public embed.FS

func logFatalAndExit(logger *log.Logger, err error) {
	if err != nil {
		logger.Fatal(err)
	}
}

func setupDb(ctx context.Context, uri string) (*mongo.Client, error) {
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return mongoClient, err
}

func closeDb(ctx context.Context, client *mongo.Client, logger *log.Logger) {
	err := client.Disconnect(ctx)
	logFatalAndExit(logger, err)
}

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := config.ReadConfig()
	logFatalAndExit(logger, err)

	publicFS, err := fs.Sub(public, "public")
	logFatalAndExit(logger, err)

	dbClient, err := setupDb(ctx, cfg.Storage.MongoDbUri)
	logFatalAndExit(logger, err)
	defer closeDb(ctx, dbClient, logger)

	mongoDb, err := storage.NewMongoDbStorage(ctx, dbClient, cfg.Storage.MongoDbName, cfg.Storage.MongoDbCollection)
	logFatalAndExit(logger, err)

	data := &app.AppData{
		Config:      cfg,
		Logger:      logger,
		Storage:     mongoDb,
		RateLimiter: rate.NewLimiter(rate.Every(time.Minute), cfg.RequestPerMinute),
		PublicFS:    publicFS,
	}

	err = app.Serve(data)
	logFatalAndExit(logger, err)
}
