// Package handlers expone los endpoints del cri-bff-citizen.
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	scrctx "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/ctx"
	screrrors "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/errors"
	"github.com/devsebas/costaricaasservice/libs/cri-lib-shared/httpx"

	"github.com/devsebas/costaricaasservice/gateway/cri-bff-citizen/internal/service"
)

type API struct {
	Service *service.Service
}

// Register monta los handlers del BFF en el sub-router App (con realm middleware).
//
//	GET /v1/citizens/{cedula}/dashboard?year=2025  → composite (tax + audit)
//	GET /v1/citizens/{cedula}/access-log
//	GET /v1/citizens/{cedula}/tax/prefilled?year=2025
func (a *API) Register(public, app chi.Router) {
	app.Get("/v1/citizens/{cedula}/dashboard", a.dashboard)
	app.Get("/v1/citizens/{cedula}/access-log", a.accessLog)
	app.Get("/v1/citizens/{cedula}/tax/prefilled", a.taxPrefilled)
	app.Get("/v1/citizens/{cedula}/health/profile", a.healthProfile)
}

func (a *API) dashboard(w http.ResponseWriter, r *http.Request) {
	realm := scrctx.Realm(r.Context())
	cedula := chi.URLParam(r, "cedula")
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	if year == 0 {
		year = 2025
	}
	out, err := a.Service.GetDashboard(r.Context(), realm, cedula, year)
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}
	httpx.OK(w, r, out)
}

func (a *API) accessLog(w http.ResponseWriter, r *http.Request) {
	realm := scrctx.Realm(r.Context())
	cedula := chi.URLParam(r, "cedula")
	out, err := a.Service.GetAccessLog(r.Context(), realm, cedula)
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}
	// out es un envelope upstream; lo desempaca para no anidar.
	var env struct {
		Data any `json:"data"`
	}
	if err := json.Unmarshal(out, &env); err != nil {
		httpx.Fail(w, r, screrrors.Wrap(screrrors.CodeInternal, "decode upstream", err))
		return
	}
	httpx.OK(w, r, env.Data)
}

func (a *API) taxPrefilled(w http.ResponseWriter, r *http.Request) {
	realm := scrctx.Realm(r.Context())
	cedula := chi.URLParam(r, "cedula")
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	if year == 0 {
		year = 2025
	}
	out, err := a.Service.GetTaxPrefilled(r.Context(), realm, cedula, year)
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}
	var env struct {
		Data any `json:"data"`
	}
	if err := json.Unmarshal(out, &env); err != nil {
		httpx.Fail(w, r, screrrors.Wrap(screrrors.CodeInternal, "decode upstream", err))
		return
	}
	httpx.OK(w, r, env.Data)
}

func (a *API) healthProfile(w http.ResponseWriter, r *http.Request) {
	realm := scrctx.Realm(r.Context())
	cedula := chi.URLParam(r, "cedula")
	out, err := a.Service.GetHealthProfile(r.Context(), realm, cedula)
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}
	var env struct {
		Data any `json:"data"`
	}
	if err := json.Unmarshal(out, &env); err != nil {
		httpx.Fail(w, r, screrrors.Wrap(screrrors.CodeInternal, "decode upstream", err))
		return
	}
	httpx.OK(w, r, env.Data)
}
