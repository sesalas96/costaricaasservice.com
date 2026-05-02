// Package inbound expone POST /v1/interop/inbox: el endpoint que recibe
// requests de otros security-servers (peers).
//
// Flujo:
//
//  1. Decodificar el envelope firmado del peer.
//  2. Validar que dst_member == este SS member.
//  3. Resolver pubkey del src_member desde el hub (con cache).
//  4. Verificar firma JWS detached (signing.Verify).
//  5. POST al member local en su MemberDispatchURL con el payload de dispatch
//     ({service, version, request_id, realm, requester_member, citizen_id,
//      purpose, body}).
//  6. Devolver la respuesta del member al peer.
package inbound

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	screrrors "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/errors"
	"github.com/devsebas/costaricaasservice/libs/cri-lib-shared/httpx"

	"github.com/devsebas/costaricaasservice/interop/cri-svc-security-server/internal/config"
	"github.com/devsebas/costaricaasservice/interop/cri-svc-security-server/internal/hubclient"
	"github.com/devsebas/costaricaasservice/interop/cri-svc-security-server/internal/signing"
)

// Handler maneja /v1/interop/inbox.
type Handler struct {
	Cfg      *config.Config
	Hub      *hubclient.Client
	HTTP     *http.Client
	AuditURL string // URL del cri-svc-interop-audit; vacío = audit deshabilitado
}

// New construye el Handler.
func New(cfg *config.Config, hub *hubclient.Client) *Handler {
	return &Handler{
		Cfg:      cfg,
		Hub:      hub,
		HTTP:     &http.Client{Timeout: cfg.CallTimeout},
		AuditURL: cfg.Audit.URL,
	}
}

// inboxBody es el wire format que el peer SS envía a este SS.
type inboxBody struct {
	Envelope  *signing.Envelope `json:"envelope"`
	CitizenID string            `json:"citizen_id"`
	Purpose   string            `json:"purpose"`
}

// dispatchEnvelope es lo que esto SS envía al member local.
type dispatchEnvelope struct {
	Service         string          `json:"service"`
	Version         string          `json:"version"`
	RequestID       string          `json:"request_id"`
	Realm           string          `json:"realm"`
	RequesterMember string          `json:"requester_member"`
	CitizenID       string          `json:"citizen_id"`
	Purpose         string          `json:"purpose"`
	Body            json.RawMessage `json:"body"`
}

// ServeHTTP implementa http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ib inboxBody
	if err := json.NewDecoder(r.Body).Decode(&ib); err != nil || ib.Envelope == nil {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeBadRequest, "invalid inbox body"))
		return
	}
	env := ib.Envelope

	// Validaciones del header.
	if env.Header.DstMember != h.Cfg.Member.Slug {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeForbidden, "dst_member does not match"))
		return
	}
	if env.Header.Realm != h.Cfg.Member.Realm {
		httpx.Fail(w, r, screrrors.New(screrrors.CodeRealmForbidden, "realm mismatch"))
		return
	}
	// Skew defensivo: rechazar timestamps muy lejanos (5 min).
	if iat := env.Header.Iat; iat > 0 {
		now := time.Now().Unix()
		if iat > now+300 || iat < now-300 {
			httpx.Fail(w, r, screrrors.New(screrrors.CodeUnauthorized, "envelope timestamp out of window"))
			return
		}
	}

	// Resolver pubkey del src_member: primero del config local (anchor), luego
	// fallback al hub (descubrimiento dinámico).
	pub, err := h.resolvePubkey(r.Context(), env.Header.Realm, env.Header.SrcMember)
	if err != nil {
		httpx.Fail(w, r, screrrors.Wrap(screrrors.CodeInteropDenied, "lookup src pubkey", err))
		return
	}
	if err := signing.Verify(env, pub); err != nil {
		// Posible rotación de clave: invalidar cache del hub y reintentar.
		h.Hub.Invalidate(env.Header.Realm, env.Header.SrcMember)
		pub2, err2 := h.resolvePubkey(r.Context(), env.Header.Realm, env.Header.SrcMember)
		if err2 != nil || signing.Verify(env, pub2) != nil {
			httpx.Fail(w, r, screrrors.New(screrrors.CodeInteropDenied, "envelope signature invalid"))
			return
		}
	}

	// Dispatch al member local.
	disp := dispatchEnvelope{
		Service:         env.Header.Service,
		Version:         env.Header.Version,
		RequestID:       env.Header.RequestID,
		Realm:           env.Header.Realm,
		RequesterMember: env.Header.SrcMember,
		CitizenID:       ib.CitizenID,
		Purpose:         ib.Purpose,
		Body:            env.Payload,
	}
	respBody, status, err := h.dispatch(r.Context(), disp)
	if err != nil {
		httpx.Fail(w, r, screrrors.Wrap(screrrors.CodeInternal, "member dispatch failed", err))
		return
	}

	// Emit audit event (síncrono pero best-effort). Si el audit svc no está
	// configurado o falla, solo se loggea — no bloquea la respuesta al peer.
	auditID := h.emitAudit(env.Header, ib, status)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if auditID != "" {
		w.Header().Set("X-CRI-Audit-Id", auditID)
	}
	w.WriteHeader(status)
	_, _ = w.Write(respBody)
}

// emitAudit hace POST /internal/audit/entries al cri-svc-interop-audit.
// Best-effort: errores se loggean pero no propagan.
func (h *Handler) emitAudit(header signing.Header, ib inboxBody, status int) string {
	if h.AuditURL == "" {
		return ""
	}
	body, err := json.Marshal(map[string]any{
		"requesterMember": header.SrcMember,
		"targetMember":    header.DstMember,
		"service":         header.Service,
		"version":         header.Version,
		"citizenId":       ib.CitizenID,
		"purpose":         ib.Purpose,
		"requestId":       header.RequestID,
		"status":          status,
	})
	if err != nil {
		slog.Warn("audit emit: marshal", "err", err.Error())
		return ""
	}
	url := h.AuditURL + "/internal/audit/entries?realm=" + header.Realm
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		slog.Warn("audit emit: build req", "err", err.Error())
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-Id", header.RequestID)
	resp, err := h.HTTP.Do(req)
	if err != nil {
		slog.Warn("audit emit: post", "err", err.Error())
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		slog.Warn("audit emit: non-201", "status", resp.StatusCode)
		return ""
	}
	var env struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return ""
	}
	return env.Data.ID
}

// resolvePubkey intenta primero el anchor local del config (peers[*].public_key_path)
// y solo si no está hace fetch al hub.
func (h *Handler) resolvePubkey(ctx context.Context, realm, member string) (any, error) {
	if path := h.Cfg.PeerKeyPath(member); path != "" {
		pemBytes, err := os.ReadFile(path)
		if err == nil {
			block, _ := pem.Decode(pemBytes)
			if block != nil {
				if pub, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
					slog.Debug("resolved pubkey from local anchor", "member", member, "path", path)
					return pub, nil
				}
			}
		}
		slog.Debug("local pubkey path failed, falling back to hub", "member", member, "err", err)
	}
	pub, _, err := h.Hub.PublicKey(ctx, realm, member)
	return pub, err
}

func (h *Handler) dispatch(ctx context.Context, env dispatchEnvelope) ([]byte, int, error) {
	if h.Cfg.Member.MemberDispatchURL == "" {
		return nil, 0, errors.New("inbound: member_dispatch_url not configured")
	}
	body, err := json.Marshal(env)
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.Cfg.Member.MemberDispatchURL, bytes.NewReader(body))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-Id", env.RequestID)
	resp, err := h.HTTP.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("dispatch http: %w", err)
	}
	defer resp.Body.Close()
	rb, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return rb, resp.StatusCode, nil
}
