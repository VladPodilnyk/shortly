package storage

import (
	"context"
	"errors"
)

type Storage interface {
	Save(ctx context.Context, givenUrl string, shortUrl string) error
	Get(ctx context.Context, shortUrl string) (string, error)
}

var ErrRecordNotFound = errors.New("record not found")
var ErrInternalError = errors.New("internal error")
