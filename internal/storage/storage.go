package storage

import (
	"errors"
)

type Storage interface {
	Save(givenUrl string, shortUrl string)
	Get(shortUrl string) (string, error)
}

var ErrRecordNotFound = errors.New("record not found")
