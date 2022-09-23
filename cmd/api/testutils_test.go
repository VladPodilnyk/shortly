package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"golang.org/x/time/rate"
	"short.io/internal/encoder"
	"short.io/internal/storage"
)

type testServer struct {
	*httptest.Server
}

// return test application instance
func newTestApp() *application {
	app := &application{
		config:      AppConfig{Environment: "testing"},
		version:     "1.0.0",
		storage:     storage.New(),
		algorithm:   encoder.MakeSimple("https://short.est/", 0),
		ratelimiter: rate.NewLimiter(rate.Every(5*time.Second), 10),
	}
	return app
}

func newTestServer(t *testing.T, routes http.Handler) *testServer {
	ts := httptest.NewServer(routes)

	// Disable redirects
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, path string, dst interface{}) int {
	result, err := ts.Client().Get(ts.URL + path)
	if err != nil {
		t.Fatal(err)
	}

	defer result.Body.Close()
	decoder := json.NewDecoder(result.Body)
	err = decoder.Decode(dst)
	if err != nil {
		t.Fatal(err)
	}

	return result.StatusCode
}

func (ts *testServer) post(t *testing.T, path string, payload string, dst interface{}) int {
	contentType := "application/json"
	readerFromPayload := strings.NewReader(payload)
	result, err := ts.Client().Post(ts.URL+path, contentType, readerFromPayload)
	if err != nil {
		t.Fatal(err)
	}

	defer result.Body.Close()
	decoder := json.NewDecoder(result.Body)
	err = decoder.Decode(dst)
	if err != nil {
		t.Fatal(err)
	}

	return result.StatusCode
}

func checkStatusCode(t *testing.T, code int) {
	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}
}
