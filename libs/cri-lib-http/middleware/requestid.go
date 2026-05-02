// Package middleware ofrece middlewares chi compartidos por todos los servicios:
// requestId, logging estructurado, recover, realm resolver, CORS.
package middleware

import (
	"net/http"

	scrctx "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/ctx"
	"github.com/devsebas/costaricaasservice/libs/cri-lib-shared/idgen"
)

const (
	HeaderRequestID = "X-Request-Id"
)

// RequestID propaga el requestId existente o genera uno nuevo (ULID).
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(HeaderRequestID)
		if !validID(id) {
			id = idgen.New()
		}
		w.Header().Set(HeaderRequestID, id)
		ctx := scrctx.WithRequestID(r.Context(), id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// validID acepta ULIDs y, defensivamente, IDs alfanuméricos cortos.
func validID(s string) bool {
	if s == "" || len(s) > 64 {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || c == '-' || c == '_') {
			return false
		}
	}
	return true
}
