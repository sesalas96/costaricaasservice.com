// Package middleware contiene middlewares chi para enforce auth en servicios
// internos. Asume que el gateway ya validó el JWT y inyectó headers.
package middleware

import (
	"net/http"

	"github.com/devsebas/costaricaasservice/libs/cri-lib-auth/auth"
	scrctx "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/ctx"
	screrrors "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/errors"
	"github.com/devsebas/costaricaasservice/libs/cri-lib-shared/httpx"
)

// Header keys inyectados por el gateway tras verificar el JWT.
const (
	HeaderSub    = "X-CRI-Sub"
	HeaderRoles  = "X-CRI-Roles"
	HeaderRealm  = "X-CRI-Realm"
	HeaderMember = "X-CRI-Member"
	HeaderJTI    = "X-CRI-Jti"
)

// FromGatewayHeaders construye un Principal desde los headers que el gateway
// inyecta y lo coloca en el context. Falla si faltan los headers críticos.
func FromGatewayHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sub := r.Header.Get(HeaderSub)
		realm := r.Header.Get(HeaderRealm)
		if sub == "" || realm == "" {
			httpx.Fail(w, r, screrrors.New(screrrors.CodeUnauthorized, "missing gateway auth headers"))
			return
		}
		p := &auth.Principal{
			Sub:    sub,
			Realm:  realm,
			Member: r.Header.Get(HeaderMember),
			JTI:    r.Header.Get(HeaderJTI),
			Roles:  parseRoles(r.Header.Get(HeaderRoles)),
		}
		ctx := auth.WithPrincipal(r.Context(), p)
		ctx = scrctx.WithRealm(ctx, realm)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole devuelve un middleware que rechaza si el principal no tiene al
// menos uno de los roles requeridos.
func RequireRole(roles ...auth.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p := auth.FromContext(r.Context())
			if p == nil {
				httpx.Fail(w, r, screrrors.New(screrrors.CodeUnauthorized, "no principal"))
				return
			}
			if !auth.HasAny(p.Roles, roles...) {
				httpx.Fail(w, r, screrrors.New(screrrors.CodeForbidden, "insufficient role"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func parseRoles(s string) []auth.Role {
	if s == "" {
		return nil
	}
	out := make([]auth.Role, 0, 4)
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			if i > start {
				out = append(out, auth.Role(s[start:i]))
			}
			start = i + 1
		}
	}
	return out
}
