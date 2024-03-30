package main

import (
	"net/http"
)

func (app *application) logError(err error) {
	app.logger.Println(err)
}

func (app *application) errorResponse(w http.ResponseWriter, status int, message interface{}) {
	response := map[string]interface{}{"error": message}
	err := app.writeJSON(w, status, response, nil)
	if err != nil {
		app.logError(err)
		w.WriteHeader(500)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, err error) {
	app.logError(err)
	message := "server encountered a problem and couldn't proceed with your request"
	app.errorResponse(w, http.StatusNotFound, message)
}

func (app *application) notFoundResponse(w http.ResponseWriter) {
	message := "requested resource is not exists"
	app.errorResponse(w, http.StatusNotFound, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, err error) {
	app.errorResponse(w, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(w http.ResponseWriter, errors map[string]string) {
	app.errorResponse(w, http.StatusUnprocessableEntity, errors)
}

func (app *application) rateLimiterResponse(w http.ResponseWriter) {
	message := "request was rate limited"
	app.errorResponse(w, http.StatusTooManyRequests, message)
}
