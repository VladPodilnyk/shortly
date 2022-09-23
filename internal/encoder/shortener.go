package encoder

import (
	"strconv"
	"sync"
)

type UrlEncoder interface {
	Encode(givenUrl string, alias string) string
}

func MakeSimple(prefix string, counter uint64) *SimpleUrlEncoder {
	return &SimpleUrlEncoder{prefix: prefix, counter: counter}
}

type SimpleUrlEncoder struct {
	prefix  string
	counter uint64
	mu      sync.Mutex
}

func (shortener *SimpleUrlEncoder) Encode(givenUrl string, alias string) string {
	if alias == "" {
		shortener.mu.Lock()
		defer shortener.mu.Unlock()
		shortUrl := shortener.prefix + strconv.Itoa(int(shortener.counter))
		shortener.counter += 1
		return shortUrl
	}
	return shortener.prefix + alias
}
