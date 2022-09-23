package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	// routing some errors
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodsNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/status", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/encode", app.encodeUrlHandler)
	router.HandlerFunc(http.MethodPost, "/v1/decode", app.decodeUrlHandler)

	return router
}
