package storage

import "go.mongodb.org/mongo-driver/mongo"

const (
	DB_NAME    = "short_refs"
	COLLECTION = "short_urls"
)

type MongoDbStorage struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoDbStorage(client *mongo.Client) *MongoDbStorage {
	col := client.Database(DB_NAME).Collection(COLLECTION)
	return &MongoDbStorage{client: client, collection: col}
}

func (s *MongoDbStorage) Save(givenUrl string, shortUrl string) {
	panic("implement me")
}

func (s *MongoDbStorage) Get(shortUrl string) (string, error) {
	panic("implement me")
}
