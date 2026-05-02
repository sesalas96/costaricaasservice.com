// Package hubclient consulta el cri-svc-interop-hub para resolver claves
// públicas de members. Mantiene un cache TTL en memoria.
package hubclient

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Client consulta el hub.
type Client struct {
	hubURL     string
	httpClient *http.Client
	ttl        time.Duration

	mu    sync.RWMutex
	cache map[string]cacheEntry // key: realm+"/"+memberSlug
}

type cacheEntry struct {
	pub       any // *rsa.PublicKey | *ecdsa.PublicKey
	pemRaw    string
	expiresAt time.Time
}

// New construye un Client con cache TTL configurable.
func New(hubURL string, ttl time.Duration, httpTimeout time.Duration) *Client {
	if httpTimeout == 0 {
		httpTimeout = 5 * time.Second
	}
	return &Client{
		hubURL:     hubURL,
		httpClient: &http.Client{Timeout: httpTimeout},
		ttl:        ttl,
		cache:      make(map[string]cacheEntry),
	}
}

// PublicKey devuelve la clave pública PKIX de un member en un realm.
// Cachea por TTL. Retorna error si el hub no responde o el member no existe.
func (c *Client) PublicKey(ctx context.Context, realm, memberSlug string) (any, string, error) {
	key := realm + "/" + memberSlug

	c.mu.RLock()
	if e, ok := c.cache[key]; ok && time.Now().Before(e.expiresAt) {
		c.mu.RUnlock()
		return e.pub, e.pemRaw, nil
	}
	c.mu.RUnlock()

	pemStr, err := c.fetchPEM(ctx, realm, memberSlug)
	if err != nil {
		return nil, "", err
	}
	pub, err := parsePEM(pemStr)
	if err != nil {
		return nil, "", err
	}
	c.mu.Lock()
	c.cache[key] = cacheEntry{pub: pub, pemRaw: pemStr, expiresAt: time.Now().Add(c.ttl)}
	c.mu.Unlock()
	return pub, pemStr, nil
}

// Invalidate fuerza un fetch en la próxima consulta de un realm/member.
// Útil tras detección de firma inválida (posible rotación de clave).
func (c *Client) Invalidate(realm, memberSlug string) {
	c.mu.Lock()
	delete(c.cache, realm+"/"+memberSlug)
	c.mu.Unlock()
}

func (c *Client) fetchPEM(ctx context.Context, realm, memberSlug string) (string, error) {
	u, err := url.Parse(fmt.Sprintf("%s/internal/members/%s/public-key", c.hubURL, url.PathEscape(memberSlug)))
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("realm", realm)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("hub request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("hub returned %d", resp.StatusCode)
	}
	var env struct {
		Data struct {
			PublicKey string `json:"publicKey"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return "", err
	}
	if env.Data.PublicKey == "" {
		return "", errors.New("hub: empty publicKey")
	}
	return env.Data.PublicKey, nil
}

func parsePEM(pemStr string) (any, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("invalid PEM")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse PKIX: %w", err)
	}
	return pub, nil
}
