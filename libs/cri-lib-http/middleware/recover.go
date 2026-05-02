package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	screrrors "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/errors"
	"github.com/devsebas/costaricaasservice/libs/cri-lib-shared/httpx"
)

// Recover atrapa panics y responde 500 con envelope.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rv := recover(); rv != nil {
				slog.ErrorContext(r.Context(), "panic",
					"value", rv,
					"stack", string(debug.Stack()),
				)
				httpx.Fail(w, r, screrrors.New(screrrors.CodeInternal, "internal server error"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
