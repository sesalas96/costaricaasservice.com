// Package handlers expone los endpoints HTTP del cri-svc-hacienda.
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	scrctx "github.com/devsebas/saascr/libs/cri-lib-shared/ctx"
	screrrors "github.com/devsebas/saascr/libs/cri-lib-shared/errors"
	"github.com/devsebas/saascr/libs/cri-lib-shared/httpx"

	interop "github.com/devsebas/saascr/libs/cri-lib-interop-client/interop"

	"github.com/devsebas/saascr/members/cri-svc-hacienda/internal/service"
)

type API struct {
	Service  *service.Service
	Registry *interop.Registry
}

// Register monta los endpoints HTTP directos en `app` y el dispatcher
// inter-member en `public`.
func (a *API) Register(public, app chi.Router) {
	app.Get("/v1/tax-status/{cedula}/{year}", a.taxStatus)
	app.Get("/v1/prefilled-return/{cedula}/{year}", a.prefilledReturn)

	dispatch := interop.NewDispatcher(a.Registry)
	public.Method("POST", "/_interop/dispatch", dispatch)
}

// RegisterInteropHandlers expone los servicios de hacienda al SS local.
func RegisterInteropHandlers(reg *interop.Registry, svc *service.Service) {
	reg.Register("tax.status", "v1", func(c interop.Context, body []byte) (any, error) {
		var req struct {
			Cedula string `json:"cedula"`
			Year   int    `json:"year"`
		}
		if err := json.Unmarshal(body, &req); err != nil {
			return nil, screrrors.New(screrrors.CodeBadRequest, "invalid body")
		}
		return svc.GetTaxStatus(c.Realm(), req.Cedula, req.Year)
	})
}

func (a *API) taxStatus(w http.ResponseWriter, r *http.Request) {
	realm := scrctx.Realm(r.Context())
	cedula := chi.URLParam(r, "cedula")
	year, err := strconv.Atoi(chi.URLParam(r, "year"))
	if err != nil {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeBadRequest, "invalid year"))
		return
	}
	out, err := a.Service.GetTaxStatus(realm, cedula, year)
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}
	httpx.OK(w, r, out)
}

func (a *API) prefilledReturn(w http.ResponseWriter, r *http.Request) {
	realm := scrctx.Realm(r.Context())
	cedula := chi.URLParam(r, "cedula")
	year, err := strconv.Atoi(chi.URLParam(r, "year"))
	if err != nil {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeBadRequest, "invalid year"))
		return
	}
	out, err := a.Service.PrefilledReturn(r.Context(), realm, cedula, year)
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}
	httpx.OK(w, r, out)
}
