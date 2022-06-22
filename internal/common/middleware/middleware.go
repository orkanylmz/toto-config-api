package middleware

import (
	"context"
	"net/http"
)

// Key to use when setting the request ID.
type ctxKeyUseCacheID int

// RequestIDKey is the key that holds the unique request ID in a request context.
const UseCacheKey ctxKeyUseCacheID = 0

var UseCacheHeader = "X-Use-Cache"

func UseCache(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		useCache := true
		ctx := r.Context()
		useCacheHeaderValue := r.Header.Get(UseCacheHeader)
		if useCacheHeaderValue == "0" {
			useCache = false
		}
		ctx = context.WithValue(ctx, UseCacheKey, useCache)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
