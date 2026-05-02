package interop

import (
	"encoding/json"
	"errors"
	"net/http"

	scrctx "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/ctx"
	screrrors "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/errors"
	"github.com/devsebas/costaricaasservice/libs/cri-lib-shared/httpx"
)

// Dispatcher es un http.Handler que recibe llamadas inter-member ya
// verificadas por el cri-svc-security-server local y las enruta al handler
// registrado en el Registry del proceso del member.
//
// Path: POST /_interop/dispatch
// Body: { "service": "persons.get", "version": "v1", "request_id": "...",
//         "realm": "demo", "requester_member": "hacienda",
//         "citizen_id": "...", "purpose": "...", "body": <opaque JSON> }
type Dispatcher struct {
	registry *Registry
}

// NewDispatcher crea un Dispatcher.
func NewDispatcher(reg *Registry) *Dispatcher { return &Dispatcher{registry: reg} }

// dispatchEnvelope es lo que el security-server envía al member tras
// verificar la firma del request entrante.
type dispatchEnvelope struct {
	Service         string          `json:"service"`
	Version         string          `json:"version"`
	RequestID       string          `json:"request_id"`
	Realm           string          `json:"realm"`
	RequesterMember string          `json:"requester_member"`
	CitizenID       string          `json:"citizen_id"`
	Purpose         string          `json:"purpose"`
	Body            json.RawMessage `json:"body"`
}

// ServeHTTP implementa http.Handler.
func (d *Dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var env dispatchEnvelope
	if err := json.NewDecoder(r.Body).Decode(&env); err != nil {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeBadRequest, "invalid dispatch envelope"))
		return
	}
	if env.Service == "" {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeBadRequest, "service required"))
		return
	}
	if env.Version == "" {
		env.Version = "v1"
	}
	handler, err := d.registry.Lookup(env.Service, env.Version)
	if err != nil {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeNotFound, "service handler not registered"))
		return
	}

	dctx := &dispatchContext{env: env, ctx: scrctx.WithRealm(scrctx.WithRequestID(r.Context(), env.RequestID), env.Realm)}
	resp, err := handler(dctx, env.Body)
	if err != nil {
		var appErr *screrrors.AppError
		if errors.As(err, &appErr) {
			httpx.Fail(w, r, appErr)
			return
		}
		httpx.Fail(w, r, screrrors.Wrap(screrrors.CodeInternal, "handler error", err))
		return
	}
	httpx.OK(w, r, resp)
}

// dispatchContext implementa la interfaz Context expuesta a los handlers.
type dispatchContext struct {
	env dispatchEnvelope
	ctx interface{} // type-erased to avoid context dependency in interface
}

func (d *dispatchContext) RequestID() string       { return d.env.RequestID }
func (d *dispatchContext) Realm() string           { return d.env.Realm }
func (d *dispatchContext) RequesterMember() string { return d.env.RequesterMember }
func (d *dispatchContext) CitizenID() string       { return d.env.CitizenID }
func (d *dispatchContext) Purpose() string         { return d.env.Purpose }
