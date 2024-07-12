package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"shortly.io/internal/encoder"
	"shortly.io/internal/models"
	"shortly.io/internal/storage"
)

func Routes(appData *AppData) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/encode",
		applyMiddleware(
			appData.encodeUrlHandler,
			rateLimiterMiddleware(appData),
			jsonParsingMiddleware[models.UserRequest](appData),
		),
	)
	mux.HandleFunc("POST /v1/decode",
		applyMiddleware(
			appData.decodeUrlHandler,
			rateLimiterMiddleware(appData),
			jsonParsingMiddleware[models.EncodedUrl](appData),
		),
	)
	mux.HandleFunc("GET /v1/status",
		applyMiddleware(
			appData.healthCheckHandler,
			rateLimiterMiddleware(appData),
		),
	)
	return mux
}

func (app *AppData) encodeUrlHandler(w http.ResponseWriter, r *http.Request) {
	userRequest, ok := getPayload[models.UserRequest](r.Context())
	if !ok {
		app.serverErrorResponse(w, errors.New("failed to parse request payload from the context"))
		return
	}

	validationErrors := validateUserRequest(userRequest, app.Config.AliasMaxSize)
	if len(validationErrors) > 0 {
		app.failedValidationResponse(w, validationErrors)
		return
	}

	var shortUrl string
	if len(userRequest.Alias) > 0 {
		_, err := app.Storage.Get(r.Context(), userRequest.Alias)
		switch err {
		case storage.ErrInternalError:
			app.serverErrorResponse(w, err)
			return
		case nil:
			app.failedValidationResponse(w, map[string]string{"alias": "already exists"})
			return
		}
		shortUrl = userRequest.Alias
	} else {
		shortUrl = encoder.Base58()
	}

	finalShortUrl := fmt.Sprintf("%s/%s", app.Config.Prefix, shortUrl)
	err := app.Storage.Save(r.Context(), userRequest.Url, finalShortUrl)
	if err != nil {
		app.logError(err)
		app.serverErrorResponse(w, err)
	}
	app.successfulResponse(w, models.EncodedUrl{ShortUrl: finalShortUrl})
}

func (app *AppData) decodeUrlHandler(w http.ResponseWriter, r *http.Request) {
	encodedUrl, ok := getPayload[models.EncodedUrl](r.Context())
	if !ok {
		app.serverErrorResponse(w, errors.New("failed to parse request payload from the context"))
		return
	}

	originalUrl, err := app.Storage.Get(r.Context(), encodedUrl.ShortUrl)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrRecordNotFound):
			app.notFoundResponse(w)
		default:
			app.serverErrorResponse(w, err)
		}
		return
	}
	app.successfulResponse(w, models.DecodedUrl{OriginalUrl: originalUrl})
}

func (app *AppData) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := models.SystemInfo{
		Status:      "available",
		Environment: app.Config.Environment,
		Version:     app.Version,
	}
	app.successfulResponse(w, response)
}

func getPayload[T any](ctx context.Context) (T, bool) {
	userRequest, ok := ctx.Value(ContextPayloadKey).(T)
	return userRequest, ok
}

func validateUserRequest(req models.UserRequest, maxAliasSize int) map[string]string {
	validationErrors := make(map[string]string)
	_, err := http.Get(req.Url)
	if err != nil {
		validationErrors["url"] = "broken"
	}

	if req.Alias != "" && len(req.Alias) > maxAliasSize {
		errorMessage := fmt.Sprintf("alias max size exceeded, must be lesser or equal to %d", maxAliasSize)
		validationErrors["alias"] = errorMessage
	}
	return validationErrors
}
