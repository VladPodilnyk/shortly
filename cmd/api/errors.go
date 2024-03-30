package main

import (
	"net/http"
)

func (app *application) logError(err error) {
	app.logger.Println(err)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	response := map[string]interface{}{"error": message}
	err := app.writeJSON(w, status, response, nil)
	if err != nil {
		app.logError(err)
		w.WriteHeader(500)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(err)
	message := "server encountered a problem and couldn't proceed with your request"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "requested resource is not exists"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) rateLimiterResponse(w http.ResponseWriter, r *http.Request) {
	message := "request was rate limited"
	app.errorResponse(w, r, http.StatusTooManyRequests, message)
}
