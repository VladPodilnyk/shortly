package app

import (
	"errors"
	"net/http"

	"shortly.io/internal/helpers"
)

// List of errors
var ErrInvalidPayload = errors.New("failed to parse request payload")
var ErrAliasAlreadyExists = errors.New("alias already exists")

// Response helpers
func (app *AppData) logError(err error) {
	app.Logger.Println(err)
}

func (app *AppData) errorResponse(w http.ResponseWriter, status int, message any) {
	response := map[string]any{"error": message}
	err := helpers.WriteJSON(w, status, response, nil)
	if err != nil {
		app.logError(err)
		w.WriteHeader(500)
	}
}

func (app *AppData) serverErrorResponse(w http.ResponseWriter, err error) {
	app.logError(err)
	message := "server encountered a problem and couldn't proceed with your request"
	app.errorResponse(w, http.StatusInternalServerError, message)
}

func (app *AppData) successfulResponse(w http.ResponseWriter, data any) {
	err := helpers.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}

func (app *AppData) notFoundResponse(w http.ResponseWriter) {
	message := "requested resource is not exists"
	app.errorResponse(w, http.StatusNotFound, message)
}

func (app *AppData) badRequestResponse(w http.ResponseWriter, err error) {
	app.errorResponse(w, http.StatusBadRequest, err.Error())
}

func (app *AppData) failedValidationResponse(w http.ResponseWriter, errors map[string]string) {
	app.errorResponse(w, http.StatusUnprocessableEntity, errors)
}

func (app *AppData) rateLimitedError(w http.ResponseWriter) {
	message := "request was rate limited"
	app.errorResponse(w, http.StatusTooManyRequests, message)
}
