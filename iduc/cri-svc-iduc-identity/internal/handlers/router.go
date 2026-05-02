// Package handlers expone los endpoints HTTP del cri-svc-iduc-identity.
//
// Endpoints públicos (vía gateway, requieren realm en ctx):
//
//	POST /v1/citizens          → registro
//	POST /v1/auth/login        → login (issue tokens)
//	POST /v1/auth/refresh      → rota refresh token
//	POST /v1/auth/logout       → revoca sesión
//
// Endpoint interno (consumido por el gateway, no exponer públicamente):
//
//	GET  /internal/revoked-jti/snapshot?realm=<slug>
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	scrctx "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/ctx"
	screrrors "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/errors"
	"github.com/devsebas/costaricaasservice/libs/cri-lib-shared/httpx"

	"github.com/devsebas/costaricaasservice/iduc/cri-svc-iduc-identity/internal/service"
)

// API agrupa los handlers con sus dependencias inyectadas.
type API struct {
	Service *service.Service
}

// Register monta los handlers en los routers.
//
// `app` requiere realm en ctx (rutas /v1/* expuestas públicamente vía gateway).
// `public` no requiere realm (rutas /internal/* consumidas por el gateway con
// `?realm=<slug>` query param para poder pollearlas sin header).
func (a *API) Register(public, app chi.Router) {
	app.Get("/v1/ping", ping)

	app.Post("/v1/citizens", a.register)
	app.Post("/v1/auth/login", a.login)
	app.Post("/v1/auth/refresh", a.refresh)
	app.Post("/v1/auth/logout", a.logout)

	public.Get("/internal/revoked-jti/snapshot", a.revokedSnapshot)
}

func ping(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, r, map[string]string{"pong": "true"})
}

type registerReq struct {
	Cedula   string `json:"cedula"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *API) register(w http.ResponseWriter, r *http.Request) {
	realm := scrctx.Realm(r.Context())
	if realm == "" {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeRealmRequired, "realm is required"))
		return
	}
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeBadRequest, "invalid json"))
		return
	}
	out, err := a.Service.Register(r.Context(), service.RegisterInput{
		Realm: realm, Cedula: req.Cedula, Email: req.Email, Password: req.Password,
	})
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}
	httpx.Created(w, r, out)
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *API) login(w http.ResponseWriter, r *http.Request) {
	realm := scrctx.Realm(r.Context())
	if realm == "" {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeRealmRequired, "realm is required"))
		return
	}
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeBadRequest, "invalid json"))
		return
	}
	out, err := a.Service.Login(r.Context(), service.LoginInput{
		Realm: realm, Email: req.Email, Password: req.Password,
	})
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}
	httpx.OK(w, r, out)
}

type refreshReq struct {
	RefreshToken string `json:"refreshToken"`
}

func (a *API) refresh(w http.ResponseWriter, r *http.Request) {
	realm := scrctx.Realm(r.Context())
	if realm == "" {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeRealmRequired, "realm is required"))
		return
	}
	var req refreshReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeBadRequest, "invalid json"))
		return
	}
	out, err := a.Service.Refresh(r.Context(), service.RefreshInput{
		Realm: realm, RefreshToken: req.RefreshToken,
	})
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}
	httpx.OK(w, r, out)
}

func (a *API) logout(w http.ResponseWriter, r *http.Request) {
	realm := scrctx.Realm(r.Context())
	if realm == "" {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeRealmRequired, "realm is required"))
		return
	}
	var req refreshReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeBadRequest, "invalid json"))
		return
	}
	if err := a.Service.Logout(r.Context(), realm, req.RefreshToken); err != nil {
		httpx.Fail(w, r, err)
		return
	}
	httpx.NoContent(w, r)
}

// revokedSnapshot es interno: el gateway lo consulta periódicamente con un
// `realm` query param para reconstruir su Roaring bitmap.
func (a *API) revokedSnapshot(w http.ResponseWriter, r *http.Request) {
	realmSlug := r.URL.Query().Get("realm")
	if realmSlug == "" {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeRealmRequired, "realm query param required"))
		return
	}
	jtis, err := a.Service.RevokedSnapshot(r.Context(), realmSlug)
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}
	httpx.OK(w, r, map[string]any{
		"realm": realmSlug,
		"jtis":  jtis,
		"count": len(jtis),
	})
}
