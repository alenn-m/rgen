package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/aa/aa/api/user/repositories/mysql"
	"github.com/aa/aa/models"
	authService "github.com/aa/aa/util/auth"
	"github.com/aa/aa/util/paginate"
	"github.com/aa/aa/util/req"
	"github.com/aa/aa/util/resp"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	// . "github.com/go-ozzo/ozzo-validation/v4"
	// "github.com/go-ozzo/ozzo-validation/v4/is"
)

type API struct {
	svc Repository
}

const PREFIX = "users"

func New(router chi.Router, dbClient *sqlx.DB, authSvc *authService.AuthService) {
	a := API{svc: NewController(mysql.NewUserDB(dbClient), authSvc)}

	router = router.Route(fmt.Sprintf("/%s", PREFIX), func(r chi.Router) {

		// Index
		r.Get("/", a.index)

		// Create
		r.Post("/", a.store)

		// Show
		r.Get("/{id}", a.show)

		// Update
		r.Put("/{id}", a.update)

		// Delete
		r.Delete("/{id}", a.delete)

	})
}

func (a *API) index(w http.ResponseWriter, r *http.Request) {
	pReq := paginate.ParsePaginationReq(r)

	result, total, err := a.svc.Index(r.Context(), pReq.Page)
	if err != nil {
		resp.ReturnError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.ReturnPaginatedSuccess(w, pReq.Page, total, result)
}

type StoreReq struct {
	ApiToken  string `json:"api_token"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

func (u *StoreReq) Validate() error {
	return nil
}

func (a *API) store(w http.ResponseWriter, r *http.Request) {
	var storeReq StoreReq

	err := req.DecodeRequest(r, &storeReq)
	if err != nil {
		resp.ReturnError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := a.svc.Store(r.Context(), &storeReq)
	if err != nil {
		resp.ReturnError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.ReturnSuccess(w, id)
}

func (a *API) show(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		resp.ReturnError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := a.svc.Show(r.Context(), models.UserID(id))
	if err != nil {
		resp.ReturnError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.ReturnSuccess(w, result)
}

type UpdateReq struct {
	ApiToken  string `json:"api_token"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

func (u *UpdateReq) Validate() error {
	return nil
}

func (a *API) update(w http.ResponseWriter, r *http.Request) {
	var updateReq UpdateReq

	err := req.DecodeRequest(r, &updateReq)
	if err != nil {
		resp.ReturnError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		resp.ReturnError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = a.svc.Update(r.Context(), &updateReq, models.UserID(id))
	if err != nil {
		resp.ReturnError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.ReturnSuccess(w, nil)
}

func (a *API) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		resp.ReturnError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = a.svc.Delete(r.Context(), models.UserID(id))
	if err != nil {
		resp.ReturnError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.ReturnSuccess(w, nil)
}
