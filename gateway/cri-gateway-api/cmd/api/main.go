// Binary api arranca el cri-gateway-api: punto de entrada único del realm.
//
// Stack del request:
//  1. RequestID + Recover + Logging (cri-lib-http)
//  2. ResolveRealm                  — del subdominio o header X-CRI-Realm
//  3. RateLimit                     — token bucket por sub o IP
//  4. AuthGW                        — verify JWT + bitmap revocación + header injection
//  5. ReverseProxy                  → upstream según ruta del catálogo
package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	srv "github.com/devsebas/saascr/libs/cri-lib-http/server"

	"github.com/devsebas/saascr/gateway/cri-gateway-api/internal/authgw"
	"github.com/devsebas/saascr/gateway/cri-gateway-api/internal/config"
	"github.com/devsebas/saascr/gateway/cri-gateway-api/internal/proxy"
	"github.com/devsebas/saascr/gateway/cri-gateway-api/internal/ratelimit"
	"github.com/devsebas/saascr/gateway/cri-gateway-api/internal/revocation"
	"github.com/devsebas/saascr/gateway/cri-gateway-api/internal/upstream"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", "err", err.Error())
		os.Exit(1)
	}
	setupLogging(cfg.Env)

	pub, err := loadPublicKey(cfg.Auth.PublicPEMPath)
	if err != nil {
		slog.Error("public key", "path", cfg.Auth.PublicPEMPath, "err", err.Error())
		os.Exit(1)
	}

	routes, err := upstream.New(cfg.Upstreams)
	if err != nil {
		slog.Error("upstreams", "err", err.Error())
		os.Exit(1)
	}
	if len(routes.Rules()) == 0 {
		slog.Warn("no upstreams configured — gateway will only serve /healthz")
	}

	revs := revocation.NewRegistry()
	poller := revocation.NewPoller(revocation.PollerConfig{
		Registry:     revs,
		IdentityURL:  cfg.Revocation.IdentityURL,
		Interval:     cfg.Revocation.PollInterval,
		RealmsToPoll: cfg.Revocation.Realms,
	})

	rl := ratelimit.New(cfg.RateLimit.Rate, cfg.RateLimit.Burst, cfg.RateLimit.IdleTTL)
	authMW := authgw.New(authgw.Config{
		PublicKey:      pub,
		Revocations:    revs,
		PublicPrefixes: cfg.Auth.PublicPrefixes,
	})

	rs, httpSrv := srv.New(srv.Config{
		Addr:         ":" + cfg.Port,
		ResolveRealm: true,
		BaseHost:     cfg.BaseHost,
	})
	rs.App.Use(ratelimit.Middleware(rl))
	rs.App.Use(authMW)
	proxyHandler := proxy.New(routes, cfg.UpstreamTimeout)
	rs.App.Handle("/*", proxyHandler)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go poller.Run(ctx)

	slog.Info("starting", "service", cfg.ServiceName, "env", cfg.Env, "port", cfg.Port,
		"upstreams", len(routes.Rules()))
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

func loadPublicKey(path string) (*rsa.PublicKey, error) {
	pemBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read pem: %w", err)
	}
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("invalid PEM")
	}
	parsed, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		k, err2 := x509.ParsePKCS1PublicKey(block.Bytes)
		if err2 != nil {
			return nil, fmt.Errorf("parse public key: %v / %v", err, err2)
		}
		return k, nil
	}
	pub, ok := parsed.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}
	return pub, nil
}
