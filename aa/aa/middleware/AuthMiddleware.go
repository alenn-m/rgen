package middleware

import (
	"context"
	"net/http"

	authService "github.com/aa/aa/util/auth"
	"github.com/go-chi/jwtauth"
)

var AuthSvc *authService.AuthService

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())
		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		if token == nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		if !token.Valid {
			AuthSvc.Logout(token.Raw)

			http.Error(w, http.StatusText(401), 401)
			return
		}

		// user is taken from cache
		user, err := AuthSvc.GetAuthUser(token.Raw)
		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		// and inserted into context
		c := context.WithValue(r.Context(), authService.AuthUser, user)
		r = r.WithContext(c)

		next.ServeHTTP(w, r)
	})
}
