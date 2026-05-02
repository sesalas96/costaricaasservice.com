// Package handlers del cri-svc-bcr.
//
// Endpoints (admin-only en producción; abierto en MVP demo):
//
//	POST /admin/kyc-check?cedula=X  → KYC completo (consulta RC + Hacienda firmadas)
//
// BCR no expone servicios al fabric, solo consume; por eso no monta dispatcher.
package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	scrctx "github.com/devsebas/saascr/libs/cri-lib-shared/ctx"
	screrrors "github.com/devsebas/saascr/libs/cri-lib-shared/errors"
	"github.com/devsebas/saascr/libs/cri-lib-shared/httpx"

	"github.com/devsebas/saascr/members/cri-svc-bcr/internal/service"
)

type API struct {
	Service *service.Service
}

func (a *API) Register(public, app chi.Router) {
	app.Post("/admin/kyc-check", a.kycCheck)
}

func (a *API) kycCheck(w http.ResponseWriter, r *http.Request) {
	realm := scrctx.Realm(r.Context())
	cedula := r.URL.Query().Get("cedula")
	if cedula == "" {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeBadRequest, "cedula query param required"))
		return
	}
	out, err := a.Service.KYCCheck(r.Context(), realm, cedula)
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}
	httpx.OK(w, r, out)
}
