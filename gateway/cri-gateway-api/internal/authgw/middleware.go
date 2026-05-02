// Package authgw implementa el middleware de autenticación del gateway:
//
//  1. Extrae el bearer token del header Authorization.
//  2. Verifica firma RS256 con la pubkey de iduc-identity.
//  3. Verifica que el JTI no esté revocado (Roaring bitmap por realm).
//  4. Resuelve el realm del request (subdominio o header X-CRI-Realm) y
//     valida que coincida con el claim del token.
//  5. Inyecta headers internos para los servicios upstream:
//     X-CRI-Sub, X-CRI-Roles, X-CRI-Realm, X-CRI-Member, X-CRI-Jti, X-Request-Id.
//
// Diseño: si el path comienza con un prefijo público (configurable), el
// middleware no exige token (e.g., /v1/auth/login, /v1/citizens registration).
package authgw

import (
	"crypto/rsa"
	"net/http"
	"strconv"
	"strings"

	scrctx "github.com/devsebas/saascr/libs/cri-lib-shared/ctx"
	screrrors "github.com/devsebas/saascr/libs/cri-lib-shared/errors"
	"github.com/devsebas/saascr/libs/cri-lib-shared/httpx"

	authjwt "github.com/devsebas/saascr/libs/cri-lib-auth/jwt"
	authmw "github.com/devsebas/saascr/libs/cri-lib-auth/middleware"

	"github.com/devsebas/saascr/gateway/cri-gateway-api/internal/revocation"
)

// Config son los parámetros del middleware.
type Config struct {
	PublicKey      *rsa.PublicKey
	Revocations    *revocation.Registry
	PublicPrefixes []string // paths que NO requieren auth (e.g., "/v1/auth/", "/v1/citizens")
}

// New construye el middleware.
func New(cfg Config) func(http.Handler) http.Handler {
	verifier := authjwt.NewVerifier(cfg.PublicKey)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isPublicPath(r.URL.Path, cfg.PublicPrefixes) {
				next.ServeHTTP(w, r)
				return
			}

			tok, ok := bearerToken(r.Header.Get("Authorization"))
			if !ok {
				httpx.Fail(w, r, screrrors.New(screrrors.CodeUnauthorized, "missing bearer token"))
				return
			}
			principal, err := verifier.Verify(tok)
			if err != nil {
				httpx.Fail(w, r, screrrors.New(screrrors.CodeUnauthorized, "invalid token"))
				return
			}

			// Coherencia realm-en-token vs realm-en-request.
			reqRealm := scrctx.Realm(r.Context())
			if reqRealm != "" && reqRealm != principal.Realm {
				httpx.Fail(w, r, screrrors.New(screrrors.CodeRealmForbidden, "token realm does not match request realm"))
				return
			}

			// Revocación.
			jtiNum := authjwt.JTINum(principal.JTI)
			if jtiNum > 0 && cfg.Revocations.Get(principal.Realm).Contains(jtiNum) {
				httpx.Fail(w, r, screrrors.New(screrrors.CodeUnauthorized, "token revoked"))
				return
			}

			// Inyección de headers para upstream.
			r.Header.Set(authmw.HeaderSub, principal.Sub)
			r.Header.Set(authmw.HeaderRealm, principal.Realm)
			r.Header.Set(authmw.HeaderMember, principal.Member)
			r.Header.Set(authmw.HeaderJTI, strconv.FormatUint(uint64(jtiNum), 10))
			r.Header.Set(authmw.HeaderRoles, joinRoles(principal.Roles))

			next.ServeHTTP(w, r)
		})
	}
}

func bearerToken(h string) (string, bool) {
	const prefix = "Bearer "
	if !strings.HasPrefix(h, prefix) {
		return "", false
	}
	return strings.TrimSpace(h[len(prefix):]), true
}

func isPublicPath(path string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

func joinRoles[T ~string](rs []T) string {
	if len(rs) == 0 {
		return ""
	}
	out := make([]string, len(rs))
	for i, r := range rs {
		out[i] = string(r)
	}
	return strings.Join(out, ",")
}
