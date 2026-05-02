// Package handlers expone los endpoints HTTP del cri-svc-interop-hub.
//
// Endpoints públicos (vía gateway, requieren realm):
//
//	GET    /v1/catalog/{member}                       → lista servicios expuestos por un member
//	GET    /v1/catalog/{member}/{service}/{version}   → metadata de un servicio
//
// Endpoints administrativos (requieren realm + rol admin — gateway hace el RBAC):
//
//	POST   /admin/members
//	GET    /admin/members
//	GET    /admin/members/{member}
//	POST   /admin/members/{member}/services
//
// Endpoints internos (NO exponer públicamente; consumidos por security-server / router):
//
//	GET    /internal/members/{member}/public-key?realm=<slug>
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	scrctx "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/ctx"
	screrrors "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/errors"
	"github.com/devsebas/costaricaasservice/libs/cri-lib-shared/httpx"

	"github.com/devsebas/costaricaasservice/interop/cri-svc-interop-hub/internal/service"
)

type API struct {
	Service *service.Service
}

func (a *API) Register(public, app chi.Router) {
	app.Get("/v1/catalog/{member}", a.listServices)
	app.Get("/v1/catalog/{member}/{service}/{version}", a.lookupService)

	app.Post("/admin/members", a.registerMember)
	app.Get("/admin/members", a.listMembers)
	app.Get("/admin/members/{member}", a.getMember)
	app.Post("/admin/members/{member}/services", a.publishService)

	public.Get("/internal/members/{member}/public-key", a.internalMemberPubKey)
}

type registerMemberReq struct {
	MemberID    string `json:"memberId"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	PublicKey   string `json:"publicKey"`
}

func (a *API) registerMember(w http.ResponseWriter, r *http.Request) {
	realm := scrctx.Realm(r.Context())
	var req registerMemberReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeBadRequest, "invalid json"))
		return
	}
	out, err := a.Service.RegisterMember(r.Context(), service.RegisterMemberInput{
		Realm:       realm,
		MemberID:    req.MemberID,
		DisplayName: req.DisplayName,
		Description: req.Description,
		PublicKey:   req.PublicKey,
	})
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}
	httpx.Created(w, r, out)
}

func (a *API) listMembers(w http.ResponseWriter, r *http.Request) {
	realm := scrctx.Realm(r.Context())
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	out, err := a.Service.ListMembers(r.Context(), realm, limit, offset)
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}
	httpx.OK(w, r, map[string]any{"members": out, "count": len(out)})
}

func (a *API) getMember(w http.ResponseWriter, r *http.Request) {
	realm := scrctx.Realm(r.Context())
	slug := chi.URLParam(r, "member")
	out, err := a.Service.GetMember(r.Context(), realm, slug)
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}
	httpx.OK(w, r, out)
}

type publishServiceReq struct {
	ServiceID   string `json:"serviceId"`
	Version     string `json:"version"`
	Description string `json:"description"`
	SchemaURL   string `json:"schemaUrl"`
}

func (a *API) publishService(w http.ResponseWriter, r *http.Request) {
	realm := scrctx.Realm(r.Context())
	memberSlug := chi.URLParam(r, "member")
	var req publishServiceReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeBadRequest, "invalid json"))
		return
	}
	out, err := a.Service.PublishService(r.Context(), service.PublishServiceInput{
		Realm:       realm,
		MemberSlug:  memberSlug,
		ServiceID:   req.ServiceID,
		Version:     req.Version,
		Description: req.Description,
		SchemaURL:   req.SchemaURL,
	})
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}
	httpx.Created(w, r, out)
}

func (a *API) listServices(w http.ResponseWriter, r *http.Request) {
	realm := scrctx.Realm(r.Context())
	memberSlug := chi.URLParam(r, "member")
	out, err := a.Service.ListServicesByMember(r.Context(), realm, memberSlug)
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}
	httpx.OK(w, r, map[string]any{"member": memberSlug, "services": out, "count": len(out)})
}

func (a *API) lookupService(w http.ResponseWriter, r *http.Request) {
	realm := scrctx.Realm(r.Context())
	memberSlug := chi.URLParam(r, "member")
	serviceID := chi.URLParam(r, "service")
	version := chi.URLParam(r, "version")
	out, err := a.Service.LookupService(r.Context(), realm, memberSlug, serviceID, version)
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}
	httpx.OK(w, r, out)
}

// internalMemberPubKey está fuera del middleware de realm (router público).
// Lee `realm` del query param porque lo invocan security-server y router
// sin pasar por gateway.
func (a *API) internalMemberPubKey(w http.ResponseWriter, r *http.Request) {
	realmSlug := r.URL.Query().Get("realm")
	if realmSlug == "" {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeRealmRequired, "realm query param required"))
		return
	}
	memberSlug := chi.URLParam(r, "member")
	out, err := a.Service.GetMemberWithKey(r.Context(), realmSlug, memberSlug)
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}
	httpx.OK(w, r, out)
}
