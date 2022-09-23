package storage

import (
	"errors"
)

type Storage interface {
	Save(givenUrl string, shortUrl string)
	Get(shortUrl string) (string, error)
}

var ErrRecordNotFound = errors.New("record not found")

type InMemoryStorage struct {
	longToShortForm map[string]string
	shortToLongForm map[string]string
}

func New() *InMemoryStorage {
	return &InMemoryStorage{longToShortForm: map[string]string{}, shortToLongForm: map[string]string{}}
}

func (storage *InMemoryStorage) Save(givenUrl string, shortUrl string) {
	if _, exists := storage.longToShortForm[givenUrl]; !exists {
		storage.longToShortForm[givenUrl] = shortUrl
		storage.shortToLongForm[shortUrl] = givenUrl
	}
}

func (storage *InMemoryStorage) Get(shortUrl string) (string, error) {
	value, isExist := storage.shortToLongForm[shortUrl]
	if !isExist {
		return "", ErrRecordNotFound
	}
	return value, nil
}
