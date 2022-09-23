package main

import (
	"net/http"

	"short.io/internal/data"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	isOk := app.ratelimiter.Allow()
	if !isOk {
		app.rateLimitedReponse(w, r)
		return
	}

	response := data.SystemInfo{
		Status:      "available",
		Environment: app.config.Environment,
		Version:     app.version,
	}

	err := app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
