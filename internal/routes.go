package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/status", app.healthCheckHandler)
	mux.HandleFunc("POST /v1/encode", app.encodeUrlHandler)
	mux.HandleFunc("POST /v1/decode", app.decodeUrlHandler)

	return mux
}
