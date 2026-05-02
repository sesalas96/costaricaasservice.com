// Package proxy implementa el reverse proxy del gateway.
package proxy

import (
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	scrctx "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/ctx"
	screrrors "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/errors"
	"github.com/devsebas/costaricaasservice/libs/cri-lib-shared/httpx"

	"github.com/devsebas/costaricaasservice/gateway/cri-gateway-api/internal/upstream"
)

// Handler implementa http.Handler. Para cada request, busca la Route que
// matchea el path; si no hay match retorna 404 con envelope. Si hay match,
// crea un httputil.ReverseProxy hacia el Target y lo invoca, propagando
// X-Request-Id y los headers de auth ya inyectados por el middleware.
type Handler struct {
	routes  *upstream.Routes
	timeout time.Duration
}

// New construye un Handler.
func New(routes *upstream.Routes, timeout time.Duration) *Handler {
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &Handler{routes: routes, timeout: timeout}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := h.routes.Match(r.URL.Path)
	if route == nil {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeNotFound, "no route for path"))
		return
	}

	rp := httputil.NewSingleHostReverseProxy(route.Target)

	// Customizar Director: preservar headers de auth, possibly strip prefix.
	originalDirector := rp.Director
	rp.Director = func(req *http.Request) {
		originalDirector(req)
		if route.StripPrefix && strings.HasPrefix(req.URL.Path, route.Prefix) {
			req.URL.Path = strings.TrimPrefix(req.URL.Path, route.Prefix)
			if req.URL.Path == "" {
				req.URL.Path = "/"
			}
		}
		// Asegurarse de que el Host header sea el del upstream (no el cliente).
		req.Host = route.Target.Host
		// Propagar request id
		if rid := scrctx.RequestID(req.Context()); rid != "" {
			req.Header.Set("X-Request-Id", rid)
		}
	}

	rp.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
		httpx.Fail(w, req, screrrors.Wrap(screrrors.CodeUnavailable, "upstream unavailable", err))
	}

	rp.ServeHTTP(w, r)
}
