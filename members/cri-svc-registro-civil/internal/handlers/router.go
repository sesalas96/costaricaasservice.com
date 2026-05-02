// Package handlers expone los endpoints HTTP del cri-svc-registro-civil.
//
// Endpoints HTTP directos (vía gateway o internos del realm; requieren realm):
//
//	GET /v1/persons/{cedula}                → persona
//	GET /v1/persons/{cedula}/vital-events   → eventos vitales
//
// Endpoint de dispatch (consumido por cri-svc-security-server local):
//
//	POST /_interop/dispatch                 → invoca el handler registrado
//	                                           en el Registry para (service, version)
package handlers

import (
	"encoding/json"

	"github.com/go-chi/chi/v5"
	"net/http"

	scrctx "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/ctx"
	screrrors "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/errors"
	"github.com/devsebas/costaricaasservice/libs/cri-lib-shared/httpx"

	interop "github.com/devsebas/costaricaasservice/libs/cri-lib-interop-client/interop"

	"github.com/devsebas/costaricaasservice/members/cri-svc-registro-civil/internal/service"
)

type API struct {
	Service  *service.Service
	Registry *interop.Registry
}

// Register monta los endpoints en los routers.
func (a *API) Register(public, app chi.Router) {
	app.Get("/v1/persons/{cedula}", a.getPerson)
	app.Get("/v1/persons/{cedula}/vital-events", a.vitalEvents)

	// Dispatcher para llamadas inter-member: el security-server local
	// envía aquí los requests ya verificados.
	dispatch := interop.NewDispatcher(a.Registry)
	public.Method("POST", "/_interop/dispatch", dispatch)
}

// RegisterInteropHandlers registra los handlers que este member expone vía
// interop. Llamado desde main al arrancar.
func RegisterInteropHandlers(reg *interop.Registry, svc *service.Service) {
	reg.Register("persons.get", "v1", func(c interop.Context, body []byte) (any, error) {
		var req struct {
			Cedula string `json:"cedula"`
		}
		if err := json.Unmarshal(body, &req); err != nil {
			return nil, screrrors.New(screrrors.CodeBadRequest, "invalid body")
		}
		return svc.GetPerson(c.Realm(), req.Cedula)
	})

	reg.Register("persons.vital-events", "v1", func(c interop.Context, body []byte) (any, error) {
		var req struct {
			Cedula string `json:"cedula"`
		}
		if err := json.Unmarshal(body, &req); err != nil {
			return nil, screrrors.New(screrrors.CodeBadRequest, "invalid body")
		}
		events, err := svc.VitalEvents(c.Realm(), req.Cedula)
		if err != nil {
			return nil, err
		}
		return map[string]any{"events": events, "count": len(events)}, nil
	})
}

func (a *API) getPerson(w http.ResponseWriter, r *http.Request) {
	realm := scrctx.Realm(r.Context())
	cedula := chi.URLParam(r, "cedula")
	out, err := a.Service.GetPerson(realm, cedula)
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}
	httpx.OK(w, r, out)
}

func (a *API) vitalEvents(w http.ResponseWriter, r *http.Request) {
	realm := scrctx.Realm(r.Context())
	cedula := chi.URLParam(r, "cedula")
	events, err := a.Service.VitalEvents(realm, cedula)
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}
	httpx.OK(w, r, map[string]any{"events": events, "count": len(events)})
}
