package transport

const T_INDEX = `		// Index
		r.Get("/", a.index)`

const T_CREATE = `		// Create
		r.Post("/", a.store)`

const T_SHOW = `		// Show
		r.Get("/{id}", a.show)`
const T_UPDATE = `		// Update
		r.Put("/{id}", a.update)`
const T_DELETE = `// Delete
		r.Delete("/{id}", a.delete)`

const TRANSPORT_HEADER = `
package {{Package}}

import (
	"fmt"
	"net/http"
	"strconv"

	"{{Root}}/util/req"
	"{{Root}}/util/resp"
	"github.com/go-chi/chi"
	// . "github.com/go-ozzo/ozzo-validation/v4"
	// "github.com/go-ozzo/ozzo-validation/v4/is"
)

type API struct {
	svc Repository
}

const PREFIX = "{{Prefix}}"

func New(router chi.Router, svc Repository) {
	a := API{svc: svc}
	router = router.Route(fmt.Sprintf("/%s", PREFIX), func(r chi.Router) {
        {{TransportActions}}
	})
}`

const TRANSPORT_INDEX = `func (a *API) index(w http.ResponseWriter, r *http.Request) {
	result, err := a.svc.Index(r.Context())
	if err != nil {
		resp.ReturnError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.ReturnSuccess(w, result)
}`

const TRANSPORT_CREATE = `type StoreReq struct {
	{{Fields}}
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
}`

const TRANSPORT_SHOW = `func (a *API) show(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		resp.ReturnError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := a.svc.Show(r.Context(), id)
	if err != nil {
		resp.ReturnError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.ReturnSuccess(w, result)
}`

const TRANSPORT_UPDATE = `type UpdateReq struct {
	{{Fields}}
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

	err = a.svc.Update(r.Context(), &updateReq, id)
	if err != nil {
		resp.ReturnError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.ReturnSuccess(w, nil)
}`

const TRANSPORT_DELETE = `func (a *API) delete(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        resp.ReturnError(w, err.Error(), http.StatusInternalServerError)
        return
    }

    err = a.svc.Delete(r.Context(), id)
    if err != nil {
        resp.ReturnError(w, err.Error(), http.StatusInternalServerError)
        return
    }

    resp.ReturnSuccess(w, nil)
}`
