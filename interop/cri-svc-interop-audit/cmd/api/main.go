// Binary api arranca el cri-svc-interop-audit (audit log inter-member).
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	srv "github.com/devsebas/saascr/libs/cri-lib-http/server"

	"github.com/devsebas/saascr/interop/cri-svc-interop-audit/internal/auditstore"
	"github.com/devsebas/saascr/interop/cri-svc-interop-audit/internal/config"
	"github.com/devsebas/saascr/interop/cri-svc-interop-audit/internal/handlers"
	"github.com/devsebas/saascr/interop/cri-svc-interop-audit/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", "err", err.Error())
		os.Exit(1)
	}
	setupLogging(cfg.Env)

	st := auditstore.New()
	svc := service.New(st)
	api := &handlers.API{Service: svc}

	rs, httpSrv := srv.New(srv.Config{
		Addr:         ":" + cfg.Port,
		ResolveRealm: true,
		BaseHost:     cfg.BaseHost,
	})
	api.Register(rs.Public, rs.App)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	slog.Info("starting", "service", cfg.ServiceName, "env", cfg.Env, "port", cfg.Port)
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
