package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"shortly.io/internal/models"
)

type MongoDbStorage struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoDbStorage(client *mongo.Client, dbName string, collectionName string) *MongoDbStorage {
	col := client.Database(dbName).Collection(collectionName)
	return &MongoDbStorage{client: client, collection: col}
}

func (s *MongoDbStorage) Save(ctx context.Context, givenUrl string, shortUrl string) error {
	urlData := models.UrlData{ShortUrl: shortUrl, OriginalUrl: givenUrl}
	_, err := s.collection.InsertOne(ctx, urlData)
	return err
}

func (s *MongoDbStorage) Get(ctx context.Context, shortUrl string) (string, error) {
	filter := bson.M{"short_url": shortUrl}

	var result models.UrlData
	err := s.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", ErrRecordNotFound
		}
		return "", ErrInternalError
	}
	return result.OriginalUrl, nil
}
