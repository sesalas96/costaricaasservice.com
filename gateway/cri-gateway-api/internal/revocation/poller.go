package revocation

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

// Poller pulle periódicamente /internal/revoked-jti/snapshot de iduc-identity
// y actualiza el Set de cada realm conocido.
//
// Diseño:
//   - Una sola goroutine por gateway con un ticker.
//   - Por cada realm en RealmsToPoll, GET /internal/revoked-jti/snapshot?realm=...
//   - Reemplaza atómicamente el bitmap.
//   - Errores se loggean pero no detienen el poller (resiliencia).
type Poller struct {
	registry      *Registry
	identityURL   string
	httpClient    *http.Client
	interval      time.Duration
	realmsToPoll  []string
}

// PollerConfig agrupa los parámetros del Poller.
type PollerConfig struct {
	Registry     *Registry
	IdentityURL  string        // ej: "http://cri-svc-iduc-identity:18081"
	Interval     time.Duration // ej: 30s
	RealmsToPoll []string      // ej: ["demo", "cr-prod"]
	HTTPTimeout  time.Duration // ej: 5s
}

// New construye un Poller listo para arrancar.
func NewPoller(cfg PollerConfig) *Poller {
	if cfg.HTTPTimeout == 0 {
		cfg.HTTPTimeout = 5 * time.Second
	}
	if cfg.Interval == 0 {
		cfg.Interval = 30 * time.Second
	}
	return &Poller{
		registry:     cfg.Registry,
		identityURL:  cfg.IdentityURL,
		httpClient:   &http.Client{Timeout: cfg.HTTPTimeout},
		interval:     cfg.Interval,
		realmsToPoll: cfg.RealmsToPoll,
	}
}

// Run bloquea hasta que ctx se cancele. Hace un poll inmediato al arrancar
// y luego cada Interval. No retorna error (los errores por realm se loggean).
func (p *Poller) Run(ctx context.Context) {
	if p.identityURL == "" || len(p.realmsToPoll) == 0 {
		slog.Warn("revocation poller disabled (no identity_url o realms)")
		<-ctx.Done()
		return
	}
	p.tick(ctx)
	t := time.NewTicker(p.interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			p.tick(ctx)
		}
	}
}

func (p *Poller) tick(ctx context.Context) {
	for _, realm := range p.realmsToPoll {
		jtis, err := p.fetch(ctx, realm)
		if err != nil {
			slog.Warn("revocation poll failed", "realm", realm, "err", err.Error())
			continue
		}
		p.registry.Get(realm).Replace(jtis)
		slog.Debug("revocation snapshot updated", "realm", realm, "count", len(jtis))
	}
}

func (p *Poller) fetch(ctx context.Context, realm string) ([]int64, error) {
	u, err := url.Parse(p.identityURL + "/internal/revoked-jti/snapshot")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("realm", realm)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-2xx: %d", resp.StatusCode)
	}
	var env struct {
		Data struct {
			Realm string  `json:"realm"`
			JTIs  []int64 `json:"jtis"`
			Count int     `json:"count"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return nil, err
	}
	return env.Data.JTIs, nil
}
