package app

import (
	"context"
	"net/http"

	"shortly.io/internal/helpers"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc
type PayloadKey string

const ContextPayloadKey = PayloadKey("requestPayload")

func applyMiddleware(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	transformed := f
	for _, middleware := range middlewares {
		transformed = middleware(f)
	}
	return transformed
}

func rateLimiterMiddleware(app *AppData) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isOk := app.RateLimiter.Allow()
			if !isOk {
				app.rateLimitedError(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func jsonParsingMiddleware[T any](app *AppData) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var input T
			err := helpers.ReadJSON(w, r, &input)
			if err != nil {
				app.badRequestResponse(w, err)
				return
			}
			ctx := r.Context()
			ctx = context.WithValue(ctx, ContextPayloadKey, input)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func queryParsingMiddleware(app *AppData, wildcard string) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(wildcard) == 0 {
				app.serverErrorResponse(w, ErrInternalError)
				app.Logger.Printf("unexpected wildcard value, url path %s \n", r.URL.Path)
				return
			}

			token := r.PathValue(wildcard)
			if len(token) == 0 {
				app.badRequestResponse(w, ErrBrokenUrl)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, ContextPayloadKey, token)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
