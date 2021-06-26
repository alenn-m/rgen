package auth

import (
	"fmt"
	"net/http"

	"{{Root}}/api/auth/repositories/mysql"
	"{{Root}}/middleware"
	authService "{{Root}}/util/auth"
	"{{Root}}/util/req"
	"{{Root}}/util/resp"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	. "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/jmoiron/sqlx"
)

type API struct {
	svc Repository
}

const PREFIX = "auth"

func New(router chi.Router, dbClient *sqlx.DB, authSvc *authService.AuthService) {
	a := API{svc: NewController(mysql.NewAuthDB(dbClient), authSvc)}

	router = router.Route(fmt.Sprintf("/%s", PREFIX), func(r chi.Router) {
		// Login
		r.Post("/login", a.login)
		// Logout
		r.With(jwtauth.Verifier(authService.TokenAuth), middleware.AuthMiddleware).Post("/logout", a.logout)
	})
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *LoginReq) Validate() error {
	return ValidateStruct(l,
		Field(&l.Email, Required, is.Email),
		Field(&l.Password, Required),
	)
}

func (a *API) login(w http.ResponseWriter, r *http.Request) {
	var request LoginReq
	err := req.DecodeRequest(r, &request)
	if err != nil {
		resp.ReturnError(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := a.svc.Login(r.Context(), request.Email, request.Password)
	if err != nil {
		resp.ReturnError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.ReturnSuccess(w, user)
}

func (a *API) logout(w http.ResponseWriter, r *http.Request) {
	err := a.svc.Logout(r.Context())
	if err != nil {
		resp.ReturnError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.ReturnSuccess(w, nil)
}
