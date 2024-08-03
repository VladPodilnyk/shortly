package tests

import (
	"fmt"
	"testing"

	"shortly.io/internal/app"
	"shortly.io/internal/models"
)

func setup() (*testServer, func()) {
	appData := newTestApp()
	ts := newTestServer(app.Routes(appData.TestApp))
	appData.TestApp.Config.Prefix = ts.URL
	teardown := func() {
		appData.Cleanup()
		ts.Close()
	}
	return ts, teardown
}

func TestAppStatusRoute(t *testing.T) {
	ts, teardown := setup()
	defer teardown()

	var statusRespHandler models.SystemInfo
	expectsStatusResponse := models.SystemInfo{
		Status:      "available",
		Environment: "testing",
		Version:     "1.0.0",
	}
	code := ts.get(t, "/status", &statusRespHandler)
	checkStatusCode(t, code)
	if expectsStatusResponse.Show() != statusRespHandler.Show() {
		t.Errorf("expected %s, but got %s", expectsStatusResponse.Show(), statusRespHandler.Show())
	}
}

func TestAppEnd2End(t *testing.T) {
	ts, teardown := setup()
	defer teardown()

	redirectLocation := "https://www.google.com"

	// Test encode endpoint
	var encodeResHandler models.EncodedUrl
	encodeEndpointPayload := fmt.Sprintf(`{"url": "%s"}`, redirectLocation)

	code := ts.post(t, "/encode", encodeEndpointPayload, &encodeResHandler)
	checkStatusCode(t, code)

	ts.testRedirect(t, encodeResHandler.ShortUrl, redirectLocation)
}

func TestAppAllowShortAlias(t *testing.T) {
	ts, teardown := setup()
	defer teardown()

	redirectLocation := "https://www.google.com"
	testAlias := "ohirok"

	var encodeResHandler models.EncodedUrl
	encodeEndpointPayload := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, redirectLocation, testAlias)

	code := ts.post(t, "/encode", encodeEndpointPayload, &encodeResHandler)
	checkStatusCode(t, code)
	token := getTokenFromPath(encodeResHandler.ShortUrl)

	if token != testAlias {
		t.Errorf("expected %s but got %s", testAlias, token)
	}

	ts.testRedirect(t, encodeResHandler.ShortUrl, redirectLocation)
}
