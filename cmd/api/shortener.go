package main

import (
	"errors"
	"net/http"

	"short.io/internal/data"
	"short.io/internal/storage"
	"short.io/internal/validator"
)

func (app *application) encodeUrlHandler(w http.ResponseWriter, r *http.Request) {
	isOk := app.ratelimiter.Allow()
	if !isOk {
		app.rateLimitedReponse(w, r)
		return
	}

	var input data.UserRequest
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	reqValidator := validator.New(app.config.AliasMaxSize)
	reqValidator.ValidateUserRequest(input)
	if !reqValidator.Valid() {
		app.failedValidationReposnse(w, r, reqValidator.Errors)
	}

	shortUrl := app.algorithm.Encode(input.Url, input.Alias)
	app.storage.Save(input.Url, shortUrl)

	response := data.EncodedUrl{ShortUrl: shortUrl}
	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) decodeUrlHandler(w http.ResponseWriter, r *http.Request) {
	isOk := app.ratelimiter.Allow()
	if !isOk {
		app.rateLimitedReponse(w, r)
		return
	}

	var input data.EncodedUrl

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	originalUrl, err := app.storage.Get(input.ShortUrl)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	response := data.DecodedUrl{OriginalUrl: originalUrl}
	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
