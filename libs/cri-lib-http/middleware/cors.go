package middleware

import (
	"net/http"
	"strings"
)

// CORS devuelve un middleware con whitelist de orígenes. Si AllowedOrigins
// está vacío, no aplica restricción (devuelve `*`). En producción, configurar
// la lista explícita.
//
// Maneja preflight OPTIONS y inyecta headers en respuestas reales.
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	allowedSet := map[string]struct{}{}
	wildcard := len(allowedOrigins) == 0
	for _, o := range allowedOrigins {
		allowedSet[strings.TrimSpace(o)] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			allow := wildcard
			if !wildcard && origin != "" {
				_, allow = allowedSet[origin]
			}
			if allow {
				if wildcard {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Vary", "Origin")
				}
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-CRI-Realm, X-Request-Id")
				w.Header().Set("Access-Control-Max-Age", "300")
			}

			if r.Method == http.MethodOptions {
				if allow {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
