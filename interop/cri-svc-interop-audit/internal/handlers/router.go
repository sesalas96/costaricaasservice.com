// Package handlers expone los endpoints del cri-svc-interop-audit.
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	scrctx "github.com/devsebas/saascr/libs/cri-lib-shared/ctx"
	screrrors "github.com/devsebas/saascr/libs/cri-lib-shared/errors"
	"github.com/devsebas/saascr/libs/cri-lib-shared/httpx"

	"github.com/devsebas/saascr/interop/cri-svc-interop-audit/internal/service"
)

type API struct {
	Service *service.Service
}

// Register monta:
//
//   - app (require realm en ctx):
//     GET /v1/access-log/citizen/{cedula}  → bitácora ciudadana
//     GET /v1/audit/verify                  → integridad de la cadena
//
//   - public (sin realm middleware; toman ?realm=<slug> del query):
//     POST /internal/audit/entries          → SS-inbound emite aquí
//     GET  /internal/audit/verify           → verificador en CI
func (a *API) Register(public, app chi.Router) {
	app.Get("/v1/access-log/citizen/{cedula}", a.accessLog)
	app.Get("/v1/audit/verify", a.verify)

	public.Post("/internal/audit/entries", a.append)
	public.Get("/internal/audit/verify", a.verify)
}

type appendReq struct {
	RequesterMember string `json:"requesterMember"`
	TargetMember    string `json:"targetMember"`
	Service         string `json:"service"`
	Version         string `json:"version"`
	CitizenID       string `json:"citizenId"`
	Purpose         string `json:"purpose"`
	RequestID       string `json:"requestId"`
	Status          int    `json:"status"`
}

func (a *API) append(w http.ResponseWriter, r *http.Request) {
	realm := r.URL.Query().Get("realm")
	if realm == "" {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeRealmRequired, "realm query param required"))
		return
	}
	var req appendReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeBadRequest, "invalid json"))
		return
	}
	out, err := a.Service.Append(service.AppendInput{
		Realm:           realm,
		RequesterMember: req.RequesterMember,
		TargetMember:    req.TargetMember,
		Service:         req.Service,
		Version:         req.Version,
		CitizenID:       req.CitizenID,
		Purpose:         req.Purpose,
		RequestID:       req.RequestID,
		Status:          req.Status,
	})
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}
	httpx.Created(w, r, out)
}

func (a *API) accessLog(w http.ResponseWriter, r *http.Request) {
	realm := scrctx.Realm(r.Context())
	cedula := chi.URLParam(r, "cedula")
	out, err := a.Service.AccessLog(realm, cedula)
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}
	httpx.OK(w, r, map[string]any{
		"citizenId": cedula,
		"realm":     realm,
		"entries":   out,
		"count":     len(out),
	})
}

// verify acepta realm de ctx (rutas /v1/*) o query (rutas /internal/*).
func (a *API) verify(w http.ResponseWriter, r *http.Request) {
	realm := scrctx.Realm(r.Context())
	if realm == "" {
		realm = r.URL.Query().Get("realm")
	}
	if realm == "" {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeRealmRequired, "realm required (ctx o query)"))
		return
	}
	out := a.Service.Verify(realm)
	httpx.OK(w, r, out)
}
