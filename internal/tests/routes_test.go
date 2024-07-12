package tests

import (
	"fmt"
	"testing"

	"shortly.io/internal/app"
	"shortly.io/internal/models"
)

func TestAppStatusRoute(t *testing.T) {
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
}

func TestAppEnd2End(t *testing.T) {
	appData := newTestApp()
	ts := newTestServer(app.Routes(appData))
	defer ts.Close()

	// Test encode endpoint
	var encodeResHandler models.EncodedUrl
	encodeEndpointPayload := `{"url": "https://www.google.com"}`

	code := ts.post(t, "/v1/encode", encodeEndpointPayload, &encodeResHandler)

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

func TestAppAllowShortAlias(t *testing.T) {
	appData := newTestApp()
	ts := newTestServer(app.Routes(appData))
	defer ts.Close()

	var encodeResHandler models.EncodedUrl
	encodeEndpointPayload := `{"url": "https://www.google.com", "alias": "ohirok"}`

	code := ts.post(t, "/v1/encode", encodeEndpointPayload, &encodeResHandler)

	checkStatusCode(t, code)
	if encodeResHandler.ShortUrl != "https://shortly.io/ohirok" {
		t.Errorf("expected https://shortly.io/ohirok but got %s", encodeResHandler.ShortUrl)
	}

	var decodeResHandler models.DecodedUrl
	decodeEndpointPayload := fmt.Sprintf(`{"short_url": "%s"}`, encodeResHandler.ShortUrl)

	code = ts.post(t, "/v1/decode", decodeEndpointPayload, &decodeResHandler)

	checkStatusCode(t, code)
	if decodeResHandler.OriginalUrl != "https://www.google.com" {
		t.Errorf("expected https://www.google.com original url but got %s", decodeResHandler.OriginalUrl)
	}
}
