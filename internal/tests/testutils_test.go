package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/time/rate"
	"shortly.io/internal/app"
	"shortly.io/internal/config"
	"shortly.io/internal/storage"
)

type testServer struct {
	*httptest.Server
}

// return test application instance
func newTestApp() *app.AppData {
	testMongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://0.0.0.0:27017"))
	if err != nil {
		fmt.Printf("failed to connect to mongodb, err: %v\n", err)
		panic(err)
	}

	app := &app.AppData{
		Config:      config.AppConfig{Environment: "testing"},
		Version:     "1.0.0",
		Logger:      nil,
		Storage:     storage.NewMongoDbStorage(testMongoClient),
		RateLimiter: rate.NewLimiter(rate.Every(5*time.Second), 10),
	}
	return app
}

func newTestServer(routes http.Handler) *testServer {
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
