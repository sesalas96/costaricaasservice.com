// Package service contiene la lógica de negocio del cri-svc-interop-hub.
//
// Responsabilidades:
//   - Mantener el catálogo de members (instituciones) y sus claves públicas.
//   - Mantener el catálogo de services publicados por cada member.
//   - Servir consultas de discovery y de pubkey al cri-svc-security-server
//     y al cri-svc-interop-router.
package service

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"regexp"
	"strings"
	"time"

	screrrors "github.com/devsebas/saascr/libs/cri-lib-shared/errors"

	"github.com/devsebas/saascr/interop/cri-svc-interop-hub/internal/store"
)

type Service struct {
	store *store.Store
}

func New(s *store.Store) *Service { return &Service{store: s} }

// MemberView es la representación pública de un member.
type MemberView struct {
	ID          string    `json:"id"`
	MemberID    string    `json:"memberId"`
	DisplayName string    `json:"displayName"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}

// MemberWithKey extiende MemberView con la clave pública (PEM). Solo se sirve
// a llamadores internos (security-server, router) que necesitan verificar firma.
type MemberWithKey struct {
	MemberView
	PublicKey string `json:"publicKey"`
}

// ServiceView es la representación pública de un servicio del catálogo.
type ServiceView struct {
	ID          string    `json:"id"`
	MemberID    string    `json:"memberId"`
	ServiceID   string    `json:"serviceId"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	SchemaURL   string    `json:"schemaUrl"`
	CreatedAt   time.Time `json:"createdAt"`
}

// RegisterMemberInput agrupa los datos de creación de member.
type RegisterMemberInput struct {
	Realm       string
	MemberID    string
	DisplayName string
	Description string
	PublicKey   string // PEM PKIX
}

// RegisterMember crea un member en el realm.
func (s *Service) RegisterMember(ctx context.Context, in RegisterMemberInput) (*MemberView, error) {
	if err := validateSlug(in.MemberID); err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.DisplayName) == "" {
		return nil, screrrors.New(screrrors.CodeValidation, "displayName required")
	}
	if err := validatePublicKeyPEM(in.PublicKey); err != nil {
		return nil, err
	}
	m, err := s.store.CreateMember(ctx, in.Realm, in.MemberID, in.DisplayName, in.Description, in.PublicKey)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, screrrors.New(screrrors.CodeConflict, "member_id already registered in this realm")
		}
		return nil, screrrors.Wrap(screrrors.CodeInternal, "create member", err)
	}
	return memberView(m), nil
}

// GetMember devuelve un member sin la clave pública (vista pública).
func (s *Service) GetMember(ctx context.Context, realm, memberSlug string) (*MemberView, error) {
	m, err := s.store.GetMemberBySlug(ctx, realm, memberSlug)
	if errors.Is(err, store.ErrNotFound) {
		return nil, screrrors.New(screrrors.CodeNotFound, "member not found")
	}
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeInternal, "get member", err)
	}
	return memberView(m), nil
}

// GetMemberWithKey devuelve un member CON la clave pública.
// Endpoint interno consumido por security-server y router.
func (s *Service) GetMemberWithKey(ctx context.Context, realm, memberSlug string) (*MemberWithKey, error) {
	m, err := s.store.GetMemberBySlug(ctx, realm, memberSlug)
	if errors.Is(err, store.ErrNotFound) {
		return nil, screrrors.New(screrrors.CodeNotFound, "member not found")
	}
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeInternal, "get member", err)
	}
	return &MemberWithKey{MemberView: *memberView(m), PublicKey: m.PublicKey}, nil
}

// ListMembers retorna los members del realm con paginación.
func (s *Service) ListMembers(ctx context.Context, realm string, limit, offset int) ([]MemberView, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	ms, err := s.store.ListMembers(ctx, realm, limit, offset)
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeInternal, "list members", err)
	}
	out := make([]MemberView, len(ms))
	for i := range ms {
		out[i] = *memberView(&ms[i])
	}
	return out, nil
}

// PublishServiceInput agrupa los datos de publicación.
type PublishServiceInput struct {
	Realm       string
	MemberSlug  string
	ServiceID   string
	Version     string
	Description string
	SchemaURL   string
}

// PublishService publica/actualiza un servicio en el catálogo.
func (s *Service) PublishService(ctx context.Context, in PublishServiceInput) (*ServiceView, error) {
	if err := validateSlug(in.MemberSlug); err != nil {
		return nil, err
	}
	if err := validateServiceID(in.ServiceID); err != nil {
		return nil, err
	}
	if in.Version == "" {
		in.Version = "v1"
	}
	sv, err := s.store.PublishService(ctx, in.Realm, in.MemberSlug, in.ServiceID, in.Version, in.Description, in.SchemaURL)
	if errors.Is(err, store.ErrNotFound) {
		return nil, screrrors.New(screrrors.CodeNotFound, "member not found")
	}
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeInternal, "publish service", err)
	}
	return serviceView(sv), nil
}

// LookupService busca un servicio específico en el catálogo.
func (s *Service) LookupService(ctx context.Context, realm, memberSlug, serviceID, version string) (*ServiceView, error) {
	if version == "" {
		version = "v1"
	}
	sv, err := s.store.LookupService(ctx, realm, memberSlug, serviceID, version)
	if errors.Is(err, store.ErrNotFound) {
		return nil, screrrors.New(screrrors.CodeNotFound, "service not in catalog")
	}
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeInternal, "lookup service", err)
	}
	return serviceView(sv), nil
}

// ListServicesByMember lista los servicios expuestos por un member.
func (s *Service) ListServicesByMember(ctx context.Context, realm, memberSlug string) ([]ServiceView, error) {
	svs, err := s.store.ListServicesByMember(ctx, realm, memberSlug)
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeInternal, "list services", err)
	}
	out := make([]ServiceView, len(svs))
	for i := range svs {
		out[i] = *serviceView(&svs[i])
	}
	return out, nil
}

// helpers ------------------------------------------------------------------

func memberView(m *store.Member) *MemberView {
	return &MemberView{
		ID: m.ID, MemberID: m.MemberID,
		DisplayName: m.DisplayName, Description: m.Description,
		Status: m.Status, CreatedAt: m.CreatedAt,
	}
}

func serviceView(sv *store.Service) *ServiceView {
	return &ServiceView{
		ID: sv.ID, MemberID: sv.MemberSlug,
		ServiceID: sv.ServiceID, Version: sv.Version,
		Description: sv.Description, SchemaURL: sv.SchemaURL,
		CreatedAt: sv.CreatedAt,
	}
}

var slugRe = regexp.MustCompile(`^[a-z][a-z0-9-]{1,49}$`)

func validateSlug(s string) error {
	if !slugRe.MatchString(s) {
		return screrrors.New(screrrors.CodeValidation, "invalid slug (lowercase, alphanumeric + dashes, 2-50 chars)")
	}
	return nil
}

var serviceIDRe = regexp.MustCompile(`^[a-z][a-z0-9._-]{1,49}$`)

func validateServiceID(s string) error {
	if !serviceIDRe.MatchString(s) {
		return screrrors.New(screrrors.CodeValidation, "invalid service id (lowercase, alphanumeric + dots/dashes, 2-50 chars)")
	}
	return nil
}

// validatePublicKeyPEM acepta una clave pública PKIX RSA o ECDSA en PEM.
func validatePublicKeyPEM(pemStr string) error {
	if strings.TrimSpace(pemStr) == "" {
		return screrrors.New(screrrors.CodeValidation, "publicKey required")
	}
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return screrrors.New(screrrors.CodeValidation, "publicKey is not valid PEM")
	}
	if _, err := x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		return screrrors.New(screrrors.CodeValidation, "publicKey is not valid PKIX")
	}
	return nil
}
