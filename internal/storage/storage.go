package storage

import (
	"errors"
)

type Storage interface {
	Save(givenUrl string, shortUrl string)
	Get(shortUrl string) (string, error)
}

var ErrRecordNotFound = errors.New("record not found")

type InMemoryStorage struct{}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{}
}

func (s *InMemoryStorage) Save(givenUrl string, shortUrl string) {
	panic("implement me")
}

func (s *InMemoryStorage) Get(shortUrl string) (string, error) {
	panic("implement me")
}
