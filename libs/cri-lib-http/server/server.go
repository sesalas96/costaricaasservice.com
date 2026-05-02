// Package server centraliza el bootstrap de http.Server con timeouts saludables
// y el stack de middlewares por defecto que todos los servicios costaricaasservice usan.
package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	mw "github.com/devsebas/costaricaasservice/libs/cri-lib-http/middleware"
)

// Config define las opciones de bootstrap del server.
type Config struct {
	Addr            string        // ":8080"
	ReadTimeout     time.Duration // default 10s
	WriteTimeout    time.Duration // default 30s
	IdleTimeout     time.Duration // default 60s
	ShutdownTimeout time.Duration // default 15s
	BaseHost        string        // base FQDN para resolver subdominio → realm; vacío = no evaluar subdominio
	ResolveRealm    bool          // si true, instala el middleware de realm en el sub-router de aplicación
}

func (c Config) defaults() Config {
	if c.ReadTimeout == 0 {
		c.ReadTimeout = 10 * time.Second
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = 30 * time.Second
	}
	if c.IdleTimeout == 0 {
		c.IdleTimeout = 60 * time.Second
	}
	if c.ShutdownTimeout == 0 {
		c.ShutdownTimeout = 15 * time.Second
	}
	return c
}

// Routers expone los dos puntos de montaje del servicio:
//
//   - Public: rutas que NO requieren realm (probes, /healthz, /metrics, etc.).
//     El caller puede agregar más si necesita endpoints expuestos sin tenant.
//   - App: rutas de aplicación que requieren realm. ResolveRealm middleware
//     ya está instalado si Config.ResolveRealm = true. Aquí van /v1/* y /internal/*.
type Routers struct {
	Public chi.Router
	App    chi.Router
}

// New construye los routers públicos y de aplicación, con middlewares y un
// http.Server listo para arrancar.
func New(cfg Config) (*Routers, *http.Server) {
	cfg = cfg.defaults()

	root := chi.NewRouter()
	root.Use(mw.RequestID)
	root.Use(mw.Recover)
	root.Use(mw.Logging)

	// /healthz montado en root sin middleware de realm. Las probes lo invocan.
	root.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`ok`))
	})

	// Sub-router de aplicación. ResolveRealm aplica solo aquí.
	app := chi.NewRouter()
	if cfg.ResolveRealm {
		app.Use(mw.ResolveRealm(cfg.BaseHost))
	}
	root.Mount("/", app)

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      root,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	return &Routers{Public: root, App: app}, srv
}

// Run arranca el server y bloquea hasta que ctx se cancele, luego hace
// shutdown graceful con el timeout configurado.
func Run(ctx context.Context, srv *http.Server, shutdownTimeout time.Duration) error {
	errCh := make(chan error, 1)
	go func() {
		slog.Info("http listening", "addr", srv.Addr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		slog.Info("http shutting down")
		return srv.Shutdown(shutCtx)
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	}
}
