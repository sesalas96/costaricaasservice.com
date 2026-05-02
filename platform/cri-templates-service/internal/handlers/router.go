// Package handlers expone los endpoints HTTP del servicio.
// Plantilla — agregar handlers reales y registrarlos en Register.
package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/devsebas/saascr/libs/cri-lib-shared/httpx"
)

// Register monta los handlers en el router.
// `r` es el chi.Mux ya inicializado por server.New() con middlewares default.
func Register(r chi.Router) {
	r.Get("/v1/ping", ping)

	// Ejemplo: rutas internas detrás del gateway requieren middleware de auth.
	// r.Group(func(r chi.Router) {
	//     r.Use(authmw.FromGatewayHeaders)
	//     r.Get("/v1/me", me)
	// })
}

func ping(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, r, map[string]string{"pong": "true"})
}
