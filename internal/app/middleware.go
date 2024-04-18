package app

import (
	"context"
	"net/http"

	"shortly.io/internal/helpers"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

const ContextPayloadKey = "requestPayload"

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
