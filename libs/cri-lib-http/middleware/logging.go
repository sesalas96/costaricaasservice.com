package middleware

import (
	"log/slog"
	"net/http"
	"time"

	scrctx "github.com/devsebas/saascr/libs/cri-lib-shared/ctx"
)

// Logging emite un log por request con method, path, status, duración, requestId, realm.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r)
		slog.InfoContext(r.Context(), "http",
			"method", r.Method,
			"path", r.URL.Path,
			"status", sw.status,
			"durMs", time.Since(start).Milliseconds(),
			"requestId", scrctx.RequestID(r.Context()),
			"realm", scrctx.Realm(r.Context()),
		)
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (sw *statusWriter) WriteHeader(s int) {
	sw.status = s
	sw.ResponseWriter.WriteHeader(s)
}
