// Package outbound expone POST /v1/interop/call: el endpoint que el SDK
// del member (cri-lib-interop-client) usa para enviar requests inter-member.
//
// Flujo:
//
//  1. Decodificar el envelope del SDK ({source_member, target_member, service,
//     version, citizen_id, purpose, request_id, body, headers}).
//  2. Asegurar source_member == this SS member (defensa).
//  3. Resolver inbox URL del target_member desde la config (peers).
//  4. Firmar un Envelope con la priv key del member (signing.Sign).
//  5. POST al inbox del peer.
//  6. Devolver al SDK el body + metadata (audit_id stub, latency, status).
package outbound

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	scrctx "github.com/devsebas/saascr/libs/cri-lib-shared/ctx"
	screrrors "github.com/devsebas/saascr/libs/cri-lib-shared/errors"
	"github.com/devsebas/saascr/libs/cri-lib-shared/httpx"
	"github.com/devsebas/saascr/libs/cri-lib-shared/idgen"

	"github.com/devsebas/saascr/interop/cri-svc-security-server/internal/config"
	"github.com/devsebas/saascr/interop/cri-svc-security-server/internal/signing"
)

// Handler maneja /v1/interop/call.
type Handler struct {
	Cfg       *config.Config
	Signer    *signing.Signer
	HTTP      *http.Client
}

// New construye el Handler con un http client con timeout configurado.
func New(cfg *config.Config, signer *signing.Signer) *Handler {
	return &Handler{
		Cfg:    cfg,
		Signer: signer,
		HTTP:   &http.Client{Timeout: cfg.CallTimeout},
	}
}

// callRequest es lo que el SDK envía al SS local.
type callRequest struct {
	SourceMember string          `json:"source_member"`
	TargetMember string          `json:"target_member"`
	Service      string          `json:"service"`
	Version      string          `json:"version"`
	CitizenID    string          `json:"citizen_id"`
	Purpose      string          `json:"purpose"`
	RequestID    string          `json:"request_id"`
	Body         json.RawMessage `json:"body"`
	Headers      map[string]any  `json:"headers"`
}

// ServeHTTP implementa http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req callRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeBadRequest, "invalid call request"))
		return
	}
	if req.TargetMember == "" || req.Service == "" {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeBadRequest, "target_member and service are required"))
		return
	}
	if req.SourceMember != "" && req.SourceMember != h.Cfg.Member.Slug {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeForbidden, "source_member mismatch"))
		return
	}
	if req.RequestID == "" {
		req.RequestID = idgen.New()
	}
	if req.Version == "" {
		req.Version = "v1"
	}

	inboxURL, ok := h.Cfg.PeerByMember(req.TargetMember)
	if !ok {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeNotFound, "peer not configured: "+req.TargetMember))
		return
	}

	env, err := h.Signer.Sign(signing.Header{
		RequestID: req.RequestID,
		Realm:     h.Cfg.Member.Realm,
		SrcMember: h.Cfg.Member.Slug,
		DstMember: req.TargetMember,
		Service:   req.Service,
		Version:   req.Version,
	}, req.Body)
	if err != nil {
		httpx.Fail(w, r, screrrors.Wrap(screrrors.CodeInternal, "sign envelope", err))
		return
	}

	// Body para el inbox: incluye el envelope firmado + metadata de audit
	// (citizen_id, purpose) que NO va en la firma porque puede variar por
	// recorrido (ej: re-routing). MVP simple: incluimos junto a la firma.
	inboxBody, _ := json.Marshal(map[string]any{
		"envelope":   env,
		"citizen_id": req.CitizenID,
		"purpose":    req.Purpose,
	})

	start := time.Now()
	respBody, status, auditID, err := h.postInbox(r.Context(), inboxURL, inboxBody, req.RequestID)
	if err != nil {
		httpx.Fail(w, r, screrrors.Wrap(screrrors.CodeUnavailable, "peer call failed", err))
		return
	}

	httpx.OK(w, r, map[string]any{
		"audit_id":   auditID,
		"latency_ms": time.Since(start).Milliseconds(),
		"status":     status,
		"request_id": req.RequestID,
		"body":       json.RawMessage(respBody),
	})
	_ = scrctx.RequestID(r.Context())
}

func (h *Handler) postInbox(ctx context.Context, url string, body []byte, requestID string) ([]byte, int, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, 0, "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-Id", requestID)
	resp, err := h.HTTP.Do(req)
	if err != nil {
		return nil, 0, "", err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, "", err
	}
	auditID := resp.Header.Get("X-CRI-Audit-Id")
	if resp.StatusCode >= 500 {
		return respBody, resp.StatusCode, auditID, fmt.Errorf("peer 5xx: %d", resp.StatusCode)
	}
	return respBody, resp.StatusCode, auditID, nil
}
