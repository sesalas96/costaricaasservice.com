// Binary api arranca el cri-svc-security-server: el daemon que cada member
// despliega para firmar y verificar requests inter-member.
//
// Endpoints:
//
//	POST /v1/interop/call   — recibido del SDK del member local (outbound)
//	POST /v1/interop/inbox  — recibido de otros security-servers (inbound)
//
// Ningún endpoint usa el middleware de realm: el realm está fijo en config
// (un SS pertenece a un solo realm/member).
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	srv "github.com/devsebas/saascr/libs/cri-lib-http/server"

	"github.com/devsebas/saascr/interop/cri-svc-security-server/internal/config"
	"github.com/devsebas/saascr/interop/cri-svc-security-server/internal/hubclient"
	"github.com/devsebas/saascr/interop/cri-svc-security-server/internal/inbound"
	"github.com/devsebas/saascr/interop/cri-svc-security-server/internal/outbound"
	"github.com/devsebas/saascr/interop/cri-svc-security-server/internal/signing"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", "err", err.Error())
		os.Exit(1)
	}
	setupLogging(cfg.Env)

	if cfg.Member.Slug == "" || cfg.Member.Realm == "" {
		slog.Error("member.slug and member.realm must be set in config")
		os.Exit(1)
	}

	signer, err := signing.NewSignerFromPEM(cfg.Member.PrivatePEMPath, cfg.Member.Slug)
	if err != nil {
		slog.Error("signer", "path", cfg.Member.PrivatePEMPath, "err", err.Error())
		os.Exit(1)
	}

	hub := hubclient.New(cfg.Hub.URL, cfg.Hub.CacheTTL, cfg.CallTimeout)
	out := outbound.New(cfg, signer)
	in := inbound.New(cfg, hub)

	rs, httpSrv := srv.New(srv.Config{
		Addr:         ":" + cfg.Port,
		ResolveRealm: false, // SS no usa realm middleware: realm fijo por config
	})
	rs.Public.Handle("/v1/interop/call", out)
	rs.Public.Handle("/v1/interop/inbox", in)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	slog.Info("starting", "service", cfg.ServiceName, "env", cfg.Env, "port", cfg.Port,
		"member", cfg.Member.Slug, "realm", cfg.Member.Realm, "peers", len(cfg.Peers))
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
