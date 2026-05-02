package interop

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/devsebas/saascr/libs/cri-lib-shared/idgen"
)

// Client es el cliente HTTP que un svc-member usa para invocar servicios
// remotos a través del cri-svc-security-server local. El security-server se
// encarga de firmar, hacer mTLS al router, y registrar el AuditEntry.
type Client struct {
	// SecurityServerURL es el endpoint del security-server local del member.
	// Ej: "http://127.0.0.1:8090".
	SecurityServerURL string
	// SourceMember identifica el member que origina las llamadas.
	SourceMember string
	HTTPClient   *http.Client
}

// New construye un Client con timeouts saludables.
func New(securityServerURL, sourceMember string) *Client {
	return &Client{
		SecurityServerURL: securityServerURL,
		SourceMember:      sourceMember,
		HTTPClient:        &http.Client{Timeout: 30 * time.Second},
	}
}

// Call invoca un servicio remoto vía el security-server local.
func (c *Client) Call(ctx context.Context, req CallRequest) (*CallResponse, error) {
	if req.TargetMember == "" || req.Service == "" {
		return nil, errors.New("interop.Call: target_member and service are required")
	}
	if req.Timeout == 0 {
		req.Timeout = 5 * time.Second
	}

	bodyBytes, err := json.Marshal(map[string]any{
		"source_member": c.SourceMember,
		"target_member": req.TargetMember,
		"service":       req.Service,
		"version":       req.Version,
		"citizen_id":    req.CitizenID,
		"purpose":       req.Purpose,
		"request_id":    idgen.New(),
		"body":          req.Body,
		"headers":       req.Headers,
	})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.SecurityServerURL+"/v1/interop/call", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	start := time.Now()
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &CallResponse{
		Status:    resp.StatusCode,
		Body:      respBody,
		AuditID:   resp.Header.Get("X-CRI-Audit-Id"),
		EntryHash: []byte(resp.Header.Get("X-CRI-Entry-Hash")), // hex string
		Latency:   time.Since(start),
	}, nil
}
