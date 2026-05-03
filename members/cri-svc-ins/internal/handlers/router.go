// Package handlers del cri-svc-ins.
//
// HTTP directos (vía gateway, requieren realm):
//
//	GET /v1/insurance/{cedula}/profile
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

	"github.com/devsebas/costaricaasservice/members/cri-svc-ins/internal/service"
)

type API struct {
	Service  *service.Service
	Registry *interop.Registry
}

func (a *API) Register(public, app chi.Router) {
	app.Get("/v1/insurance/{cedula}/profile", a.insuranceProfile)

	dispatch := interop.NewDispatcher(a.Registry)
	public.Method("POST", "/_interop/dispatch", dispatch)
}

// RegisterInteropHandlers expone los servicios del INS al SS local.
// Otros members (ej. BCR validando SOA antes de un préstamo de auto)
// pueden consumirlos.
func RegisterInteropHandlers(reg *interop.Registry, svc *service.Service) {
	reg.Register("insurance.policies", "v1", func(c interop.Context, body []byte) (any, error) {
		var req struct {
			Cedula string `json:"cedula"`
		}
		if err := json.Unmarshal(body, &req); err != nil {
			return nil, screrrors.New(screrrors.CodeBadRequest, "invalid body")
		}
		return map[string]any{
			"cedula":   req.Cedula,
			"policies": svc.PoliciesRawForInterop(c.Realm(), req.Cedula),
		}, nil
	})
}

func (a *API) insuranceProfile(w http.ResponseWriter, r *http.Request) {
	realm := scrctx.Realm(r.Context())
	cedula := chi.URLParam(r, "cedula")
	out, err := a.Service.InsuranceProfile(r.Context(), realm, cedula)
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}
	httpx.OK(w, r, out)
}
