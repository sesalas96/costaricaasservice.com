// Package clients agrupa los HTTP clients que el BFF usa para orquestar
// las llamadas a los servicios upstream (hacienda, audit, etc.).
//
// Cada cliente recibe el realm como parámetro y lo propaga vía X-CRI-Realm.
// El BFF asume que los upstreams están dentro del mismo realm y usan el
// header como source-of-truth (modelo de servicios internos).
package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Hacienda llama al cri-svc-hacienda directo.
type Hacienda struct {
	BaseURL string
	HTTP    *http.Client
}

func NewHacienda(baseURL string) *Hacienda {
	return &Hacienda{BaseURL: baseURL, HTTP: &http.Client{Timeout: 10 * time.Second}}
}

// PrefilledReturn obtiene la declaración pre-llenada (que internamente
// dispara una llamada interop a registro-civil).
func (c *Hacienda) PrefilledReturn(ctx context.Context, realm, cedula string, year int) (json.RawMessage, error) {
	u := fmt.Sprintf("%s/v1/prefilled-return/%s/%d", c.BaseURL, url.PathEscape(cedula), year)
	return c.get(ctx, realm, u)
}

func (c *Hacienda) get(ctx context.Context, realm, u string) (json.RawMessage, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-CRI-Realm", realm)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return body, fmt.Errorf("hacienda %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

// CCSS llama al cri-svc-ccss directo.
type CCSS struct {
	BaseURL string
	HTTP    *http.Client
}

func NewCCSS(baseURL string) *CCSS {
	return &CCSS{BaseURL: baseURL, HTTP: &http.Client{Timeout: 10 * time.Second}}
}

// HealthProfile retorna el perfil del paciente (nombre via interop a registro-civil + recetas + próxima cita).
func (c *CCSS) HealthProfile(ctx context.Context, realm, cedula string) (json.RawMessage, error) {
	u := fmt.Sprintf("%s/v1/health/%s/profile", c.BaseURL, url.PathEscape(cedula))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-CRI-Realm", realm)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return body, fmt.Errorf("ccss %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

// Audit llama al cri-svc-interop-audit.
type Audit struct {
	BaseURL string
	HTTP    *http.Client
}

func NewAudit(baseURL string) *Audit {
	return &Audit{BaseURL: baseURL, HTTP: &http.Client{Timeout: 10 * time.Second}}
}

// AccessLogByCitizen retorna la bitácora ciudadana.
func (c *Audit) AccessLogByCitizen(ctx context.Context, realm, cedula string) (json.RawMessage, error) {
	u := fmt.Sprintf("%s/v1/access-log/citizen/%s", c.BaseURL, url.PathEscape(cedula))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-CRI-Realm", realm)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return body, fmt.Errorf("audit %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}
