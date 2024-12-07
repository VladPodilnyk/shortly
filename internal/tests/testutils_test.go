package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
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

var TestContext = context.Background()

type testServer struct {
	*httptest.Server
}

type testData struct {
	TestApp *app.AppData
	Cleanup func()
}

// return test application instance
func newTestApp() testData {
	testMongoClient, err := mongo.Connect(TestContext, options.Client().ApplyURI("mongodb://0.0.0.0:27017"))
	if err != nil {
		fmt.Printf("failed to connect to mongodb, err: %v\n", err)
		panic(err)
	}

	storage, err := storage.NewMongoDbStorage(TestContext, testMongoClient, "test_refs", "test_urls")
	if err != nil {
		fmt.Printf("failed to create mongodb storage, err: %v\n", err)
		panic(err)
	}

	app := &app.AppData{
		Config:      config.AppConfig{Environment: "testing", AliasMaxSize: 10},
		Logger:      nil,
		Storage:     storage,
		RateLimiter: rate.NewLimiter(rate.Every(5*time.Second), 10),
	}

	mongoClose := func() {
		truncateDatabase(testMongoClient)
		err := testMongoClient.Disconnect(TestContext)
		if err != nil {
			fmt.Printf("failed to disconnect from mongodb, err: %v\n", err)
			panic(err)
		}
	}

	return testData{TestApp: app, Cleanup: mongoClose}
}

func newTestServer(routes http.Handler) *testServer {
	ts := httptest.NewServer(routes)
	fmt.Println("Test server started at: ", ts.URL)

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

func (ts *testServer) testRedirect(t *testing.T, path string, expectedLocation string) {
	result, err := ts.Client().Get(path)
	if err != nil {
		t.Fatal(err)
	}

	statusCode := result.StatusCode
	location := result.Header.Get("Location")

	if statusCode != http.StatusFound {
		t.Errorf("want %d; got %d", http.StatusFound, statusCode)
	}

	if location != expectedLocation {
		t.Errorf("invalid redirect location, want %s; got %s", expectedLocation, location)
	}
}

func getTokenFromPath(path string) string {
	pattern := `http://.+/([a-zA-Z0-9]+)`
	re := regexp.MustCompile(pattern)
	return re.FindStringSubmatch(path)[1]
}

func checkStatusCode(t *testing.T, code int) {
	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}
}

func truncateDatabase(client *mongo.Client) {
	databases, err := client.ListDatabaseNames(TestContext, map[string]interface{}{})
	if err != nil {
		panic(err)
	}

	// Drop each database
	for _, dbName := range databases {
		// Skip system databases
		if dbName == "admin" || dbName == "local" || dbName == "config" {
			continue
		}

		err := client.Database(dbName).Drop(context.TODO())
		if err != nil {
			fmt.Printf("Failed to drop database %s: %v", dbName, err)
			panic(err)
		}
		fmt.Printf("Dropped database: %s\n", dbName)
	}
}
