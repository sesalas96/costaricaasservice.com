package middleware

import (
	"net/http"
	"strings"

	scrctx "github.com/devsebas/saascr/libs/cri-lib-shared/ctx"
	screrrors "github.com/devsebas/saascr/libs/cri-lib-shared/errors"
	"github.com/devsebas/saascr/libs/cri-lib-shared/httpx"
	"github.com/devsebas/saascr/libs/cri-lib-shared/realm"
)

const HeaderRealm = "X-CRI-Realm"

// ResolveRealm extrae el realm del request en el siguiente orden:
//  1. Header X-CRI-Realm (válido solo behind gateway / inter-svc).
//  2. Subdominio: {realm}.<host> con base configurable.
//
// Si baseHost es "", el subdominio NO se evalúa. El gateway llama con baseHost
// "saascr.io" y los servicios internos lo dejan vacío (solo confían en el header).
//
// Falla con 400 REALM_REQUIRED si no se puede determinar.
func ResolveRealm(baseHost string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// CORS preflight: dejar pasar OPTIONS sin exigir realm; el handler
			// CORS se ejecutará después y responderá apropiadamente.
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}
			rlm := r.Header.Get(HeaderRealm)
			if rlm == "" && baseHost != "" {
				rlm = subdomainOf(r.Host, baseHost)
			}
			if err := realm.Validate(rlm); err != nil {
				httpx.Fail(w, r, err)
				return
			}
			ctx := scrctx.WithRealm(r.Context(), rlm)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// subdomainOf retorna el primer label si host termina en "." + base.
// host = "cr-prod.saascr.io", base = "saascr.io" → "cr-prod".
func subdomainOf(host, base string) string {
	host = strings.ToLower(host)
	if i := strings.IndexByte(host, ':'); i >= 0 {
		host = host[:i]
	}
	suffix := "." + base
	if !strings.HasSuffix(host, suffix) {
		return ""
	}
	return strings.TrimSuffix(host, suffix)
}

// EnforceRealmExists evita continuar si el realm aún no fue resuelto en el ctx.
// Útil al final del stack de middlewares como sanity check.
func EnforceRealmExists(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if scrctx.Realm(r.Context()) == "" {
			httpx.Fail(w, r, screrrors.New(screrrors.CodeRealmRequired, "realm not set in context"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
