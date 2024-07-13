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
	ctx := r.Context()
	userRequest, ok := getPayload[models.UserRequest](ctx)
	if !ok {
		app.serverErrorResponse(w, ErrInvalidPayload)
		return
	}

	validationErrors := validateUserRequest(userRequest, app.Config.AliasMaxSize)
	if len(validationErrors) > 0 {
		app.failedValidationResponse(w, validationErrors)
		return
	}

	shortUrl, err := getShortUrl(ctx, app, userRequest.Alias)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	finalShortUrl := makeShortUrl(app.Config.Prefix, shortUrl)
	err = app.Storage.Save(ctx, userRequest.Url, finalShortUrl)
	if err != nil {
		app.logError(err)
		app.serverErrorResponse(w, err)
	}
	app.successfulResponse(w, models.EncodedUrl{ShortUrl: finalShortUrl})
}

func (app *AppData) decodeUrlHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	encodedUrl, ok := getPayload[models.EncodedUrl](ctx)
	if !ok {
		app.serverErrorResponse(w, ErrInvalidPayload)
		return
	}

	originalUrl, err := app.Storage.Get(ctx, encodedUrl.ShortUrl)
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
		validationErrors["url"] = "invalid url"
	}

	if req.Alias != "" && len(req.Alias) > maxAliasSize {
		errorMessage := fmt.Sprintf("alias max size exceeded, must be lesser or equal to %d", maxAliasSize)
		validationErrors["alias"] = errorMessage
	}
	return validationErrors
}

func makeShortUrl(prefix string, shortUrl string) string {
	return fmt.Sprintf("%s/%s", prefix, shortUrl)
}

func getShortUrl(ctx context.Context, app *AppData, alias string) (string, error) {
	var shortUrl string
	if len(alias) > 0 {
		_, err := app.Storage.Get(ctx, alias)
		switch err {
		case storage.ErrInternalError:
			return "", err
		case nil:
			return "", ErrAliasAlreadyExists
		}
		shortUrl = alias
	} else {
		shortUrl = encoder.Base58()
	}
	return shortUrl, nil
}
