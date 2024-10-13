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

	mux.Handle("GET /", http.FileServerFS(appData.PublicFS))

	mux.HandleFunc("POST /encode",
		applyMiddleware(
			appData.encodeUrlHandler,
			rateLimiterMiddleware(appData),
			jsonParsingMiddleware[models.UserRequest](appData),
		),
	)
	mux.HandleFunc("GET /{token}",
		applyMiddleware(
			appData.decodeUrlHandler,
			rateLimiterMiddleware(appData),
			queryParsingMiddleware(appData, "token"),
		),
	)
	mux.HandleFunc("GET /status",
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

	token, err := getToken(ctx, app, userRequest.Alias)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	err = app.Storage.Save(ctx, userRequest.Url, token)
	if err != nil {
		app.logError(err)
		app.serverErrorResponse(w, err)
	}
	shortUrl := withPrefix(app.Config.Prefix, token)
	app.successfulResponse(w, models.EncodedUrl{ShortUrl: shortUrl})
}

func (app *AppData) decodeUrlHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	encodedUrl, ok := getPayload[string](ctx)
	if !ok {
		app.serverErrorResponse(w, ErrInvalidPayload)
		return
	}

	originalUrl, err := app.Storage.Get(ctx, encodedUrl)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrRecordNotFound):
			app.notFoundResponse(w)
		default:
			app.serverErrorResponse(w, err)
		}
		return
	}
	http.Redirect(w, r, originalUrl, http.StatusFound)
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

func withPrefix(prefix string, shortUrl string) string {
	return fmt.Sprintf("%s/%s", prefix, shortUrl)
}

func getToken(ctx context.Context, app *AppData, alias string) (string, error) {
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
