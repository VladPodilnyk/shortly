package main

import (
	"errors"
	"net/http"

	"shortly.io/internal/models"
	"shortly.io/internal/storage"
	"shortly.io/internal/validator"
)

func (app *application) encodeUrlHandler(w http.ResponseWriter, r *http.Request) {
	isOk := app.rateLimiter.Allow()
	if !isOk {
		app.rateLimiterResponse(w)
		return
	}

	var input models.UserRequest
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	reqValidator := validator.New(app.config.AliasMaxSize)
	reqValidator.ValidateUserRequest(input)
	if !reqValidator.Valid() {
		app.failedValidationResponse(w, reqValidator.Errors)
	}

	shortUrl := app.algorithm.Encode(input.Url, input.Alias)
	app.storage.Save(input.Url, shortUrl)

	response := models.EncodedUrl{ShortUrl: shortUrl}
	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}

func (app *application) decodeUrlHandler(w http.ResponseWriter, r *http.Request) {
	isOk := app.rateLimiter.Allow()
	if !isOk {
		app.rateLimiterResponse(w)
		return
	}

	var input models.EncodedUrl

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	originalUrl, err := app.storage.Get(input.ShortUrl)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrRecordNotFound):
			app.notFoundResponse(w)
		default:
			app.serverErrorResponse(w, err)
		}
		return
	}

	response := models.DecodedUrl{OriginalUrl: originalUrl}
	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}
