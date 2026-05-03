// Binary api arranca el cri-svc-mep (CCSS member mock).
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	srv "github.com/devsebas/costaricaasservice/libs/cri-lib-http/server"
	interop "github.com/devsebas/costaricaasservice/libs/cri-lib-interop-client/interop"

	"github.com/devsebas/costaricaasservice/members/cri-svc-mep/internal/config"
	"github.com/devsebas/costaricaasservice/members/cri-svc-mep/internal/handlers"
	"github.com/devsebas/costaricaasservice/members/cri-svc-mep/internal/service"
	"github.com/devsebas/costaricaasservice/members/cri-svc-mep/internal/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", "err", err.Error())
		os.Exit(1)
	}
	setupLogging(cfg.Env)

	st := store.New()
	store.SeedDemo(st)

	icli := interop.New(cfg.SecurityServer, "mep")
	svc := service.New(st, icli)
	reg := interop.NewRegistry()
	handlers.RegisterInteropHandlers(reg, svc)

	api := &handlers.API{Service: svc, Registry: reg}

	rs, httpSrv := srv.New(srv.Config{
		Addr:         ":" + cfg.Port,
		ResolveRealm: true,
		BaseHost:     cfg.BaseHost,
	})
	api.Register(rs.Public, rs.App)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	slog.Info("starting", "service", cfg.ServiceName, "env", cfg.Env, "port", cfg.Port,
		"interop_services", len(reg.Services()), "ss", cfg.SecurityServer)
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
