// Package handlers del cri-svc-migracion.
//
// HTTP directos (vía gateway, requieren realm):
//
//	GET /v1/immigration/{cedula}/status
//
// Dispatcher (consumido por SS local):
//
//	POST /_interop/dispatch
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	scrctx "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/ctx"
	screrrors "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/errors"
	"github.com/devsebas/costaricaasservice/libs/cri-lib-shared/httpx"

	interop "github.com/devsebas/costaricaasservice/libs/cri-lib-interop-client/interop"

	"github.com/devsebas/costaricaasservice/members/cri-svc-migracion/internal/service"
)

type API struct {
	Service  *service.Service
	Registry *interop.Registry
}

func (a *API) Register(public, app chi.Router) {
	app.Get("/v1/immigration/{cedula}/status", a.immigrationStatus)

	dispatch := interop.NewDispatcher(a.Registry)
	public.Method("POST", "/_interop/dispatch", dispatch)
}

// RegisterInteropHandlers expone los servicios de Migración al SS local.
// Otros members (ej. BCR para KYC, MEP al matricular extranjeros) los consumen.
func RegisterInteropHandlers(reg *interop.Registry, svc *service.Service) {
	reg.Register("immigration.status", "v1", func(c interop.Context, body []byte) (any, error) {
		var req struct {
			Cedula string `json:"cedula"`
		}
		if err := json.Unmarshal(body, &req); err != nil {
			return nil, screrrors.New(screrrors.CodeBadRequest, "invalid body")
		}
		return map[string]any{
			"cedula": req.Cedula,
			"status": svc.StatusForInterop(c.Realm(), req.Cedula),
		}, nil
	})
}

func (a *API) immigrationStatus(w http.ResponseWriter, r *http.Request) {
	realm := scrctx.Realm(r.Context())
	cedula := chi.URLParam(r, "cedula")
	out, err := a.Service.ImmigrationProfile(r.Context(), realm, cedula)
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}
	httpx.OK(w, r, out)
}
