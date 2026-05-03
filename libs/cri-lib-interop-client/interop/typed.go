package interop

import (
	"context"
	"encoding/json"
	"fmt"
)

// CallTyped invoca Call y desempaca el doble envelope (SS → peer dispatcher)
// devolviendo el payload decodificado en T y el audit_id del entry escrito
// en cri-svc-interop-audit.
//
// El doble envelope es:
//
//	{ data: { audit_id, body: { data: <T>, meta } } , meta }    ← SS local
//	                          └────────────────┘
//	                            └── peer dispatcher
//
// Errores:
//   - error de transporte (red, timeout, etc.)
//   - peer status >= 400, con el status incluido
//   - JSON inválido en cualquier nivel del envelope
func CallTyped[T any](ctx context.Context, c *Client, req CallRequest) (T, string, error) {
	var zero T

	resp, err := c.Call(ctx, req)
	if err != nil {
		return zero, "", err
	}
	if resp.Status >= 400 {
		return zero, "", fmt.Errorf("interop: peer %q returned status %d", req.TargetMember, resp.Status)
	}

	var ssEnv struct {
		Data struct {
			AuditID string          `json:"audit_id"`
			Body    json.RawMessage `json:"body"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &ssEnv); err != nil {
		return zero, "", fmt.Errorf("interop: decode SS envelope: %w", err)
	}

	var peerEnv struct {
		Data T `json:"data"`
	}
	if err := json.Unmarshal(ssEnv.Data.Body, &peerEnv); err != nil {
		return zero, ssEnv.Data.AuditID, fmt.Errorf("interop: decode peer envelope: %w", err)
	}
	return peerEnv.Data, ssEnv.Data.AuditID, nil
}
