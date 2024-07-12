package tests

import (
	"fmt"
	"testing"

	"shortly.io/internal/app"
	"shortly.io/internal/models"
)

func TestAppEndToEnd(t *testing.T) {
	appData := newTestApp()
	ts := newTestServer(app.Routes(appData))
	defer ts.Close()

	var statusRespHandler models.SystemInfo
	expectsStatusResponse := models.SystemInfo{
		Status:      "available",
		Environment: "testing",
		Version:     "1.0.0",
	}
	code := ts.get(t, "/v1/status", &statusRespHandler)
	checkStatusCode(t, code)
	if expectsStatusResponse.Show() != statusRespHandler.Show() {
		t.Errorf("expected %s, but got %s", expectsStatusResponse.Show(), statusRespHandler.Show())
	}

	// Test encode endpoint
	var encodeResHandler models.EncodedUrl
	encodeEndpointPayload := `{"url": "https://www.google.com"}`

	code = ts.post(t, "/v1/encode", encodeEndpointPayload, &encodeResHandler)

	checkStatusCode(t, code)

	// Test decode endpoint
	decodeExpectedRes := "https://www.google.com"
	decodeEndpointPayload := fmt.Sprintf(`{"short_url": "%s"}`, encodeResHandler.ShortUrl)
	var decodedResHandler models.DecodedUrl

	code = ts.post(t, "/v1/decode", decodeEndpointPayload, &decodedResHandler)

	checkStatusCode(t, code)
	if decodedResHandler.OriginalUrl != decodeExpectedRes {
		t.Errorf("expected %s original url but go %s", decodeExpectedRes, decodedResHandler.OriginalUrl)
	}
}
