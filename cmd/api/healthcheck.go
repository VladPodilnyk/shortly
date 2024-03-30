package main

import (
	"net/http"

	"shortly.io/internal/models"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	isOk := app.rateLimiter.Allow()
	if !isOk {
		app.rateLimiterResponse(w)
		return
	}

	response := models.SystemInfo{
		Status:      "available",
		Environment: app.config.Environment,
		Version:     app.version,
	}

	err := app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}
