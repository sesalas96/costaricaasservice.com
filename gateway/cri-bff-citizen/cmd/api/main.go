// Binary api arranca el cri-bff-citizen — orquestador BFF para MiCR.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	httpmw "github.com/devsebas/saascr/libs/cri-lib-http/middleware"
	srv "github.com/devsebas/saascr/libs/cri-lib-http/server"

	"github.com/devsebas/saascr/gateway/cri-bff-citizen/internal/clients"
	"github.com/devsebas/saascr/gateway/cri-bff-citizen/internal/config"
	"github.com/devsebas/saascr/gateway/cri-bff-citizen/internal/handlers"
	"github.com/devsebas/saascr/gateway/cri-bff-citizen/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", "err", err.Error())
		os.Exit(1)
	}
	setupLogging(cfg.Env)

	hacienda := clients.NewHacienda(cfg.Upstreams.HaciendaURL)
	ccss := clients.NewCCSS(cfg.Upstreams.CCSSURL)
	audit := clients.NewAudit(cfg.Upstreams.AuditURL)

	svc := service.New(hacienda, ccss, audit)
	api := &handlers.API{Service: svc}

	rs, httpSrv := srv.New(srv.Config{
		Addr:         ":" + cfg.Port,
		ResolveRealm: true,
		BaseHost:     cfg.BaseHost,
	})
	rs.App.Use(httpmw.CORS(cfg.CORS.AllowedOrigins))
	api.Register(rs.Public, rs.App)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	slog.Info("starting", "service", cfg.ServiceName, "env", cfg.Env, "port", cfg.Port,
		"hacienda", cfg.Upstreams.HaciendaURL, "audit", cfg.Upstreams.AuditURL)
	if err := srv.Run(ctx, httpSrv, cfg.ShutdownTimeout); err != nil {
		slog.Error("server", "err", err.Error())
		os.Exit(1)
	}
}

func setupLogging(env string) {
	var h slog.Handler
	if env == "local" {
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	} else {
		h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}
	slog.SetDefault(slog.New(h))
}
