package middleware

import (
	"net/http"
)

func ExampleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// TODO

		next.ServeHTTP(w, r)
	})
}
