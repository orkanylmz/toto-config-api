package middleware

import (
	"context"
	"net/http"
)

// Key to use when setting the cache information.
type ctxKeyUseCacheID int

// UseCacheKey is the key that holds the cache information
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
