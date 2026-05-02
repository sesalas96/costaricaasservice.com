// Package service contiene la lógica de negocio de cri-svc-iduc-identity.
//
// Recibe sus dependencias por constructor (store, signer) y no toca HTTP/JSON.
// Los handlers son responsables de traducir AppError → respuesta HTTP envelope.
package service

import (
	"context"
	"errors"
	"strings"
	"time"

	screrrors "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/errors"

	"github.com/devsebas/costaricaasservice/iduc/cri-svc-iduc-identity/internal/passwords"
	"github.com/devsebas/costaricaasservice/iduc/cri-svc-iduc-identity/internal/store"
	"github.com/devsebas/costaricaasservice/iduc/cri-svc-iduc-identity/internal/token"
)

// Service agrega las capacidades de identidad.
type Service struct {
	store      *store.Store
	signer     *token.Signer
	refreshTTL time.Duration
}

// New construye un Service.
func New(s *store.Store, sg *token.Signer, refreshTTL time.Duration) *Service {
	return &Service{store: s, signer: sg, refreshTTL: refreshTTL}
}

// RegisterInput contiene los datos del registro.
type RegisterInput struct {
	Realm    string
	Cedula   string
	Email    string
	Password string
}

// CitizenView es la representación pública del ciudadano (sin password_hash).
type CitizenView struct {
	ID        string    `json:"id"`
	Cedula    string    `json:"cedula"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

// Register crea un ciudadano.
func (s *Service) Register(ctx context.Context, in RegisterInput) (*CitizenView, error) {
	if err := validateCedula(in.Cedula); err != nil {
		return nil, err
	}
	if err := validateEmail(in.Email); err != nil {
		return nil, err
	}
	if len(in.Password) < 10 {
		return nil, screrrors.New(screrrors.CodeValidation, "password too short (min 10)")
	}
	hash, err := passwords.Hash(in.Password)
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeInternal, "hash password", err)
	}
	c, err := s.store.CreateCitizen(ctx, in.Realm, in.Cedula, strings.ToLower(in.Email), hash)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, screrrors.New(screrrors.CodeConflict, "cedula or email already registered")
		}
		return nil, screrrors.Wrap(screrrors.CodeInternal, "create citizen", err)
	}
	return citizenView(c), nil
}

// LoginInput es la entrada de login.
type LoginInput struct {
	Realm    string
	Email    string
	Password string
}

// TokenPair es la salida de login y refresh.
type TokenPair struct {
	AccessToken      string      `json:"accessToken"`
	AccessExpiresAt  time.Time   `json:"accessExpiresAt"`
	RefreshToken     string      `json:"refreshToken"`
	RefreshExpiresAt time.Time   `json:"refreshExpiresAt"`
	Citizen          CitizenView `json:"citizen"`
}

// Login autentica y emite un par access+refresh.
func (s *Service) Login(ctx context.Context, in LoginInput) (*TokenPair, error) {
	c, err := s.store.GetCitizenByEmail(ctx, in.Realm, strings.ToLower(in.Email))
	if errors.Is(err, store.ErrNotFound) {
		return nil, screrrors.New(screrrors.CodeUnauthorized, "invalid credentials")
	}
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeInternal, "lookup citizen", err)
	}
	if c.Status != "active" {
		return nil, screrrors.New(screrrors.CodeForbidden, "citizen not active")
	}
	if err := passwords.Verify(in.Password, c.PasswordHash); err != nil {
		return nil, screrrors.New(screrrors.CodeUnauthorized, "invalid credentials")
	}
	return s.issueTokenPair(ctx, in.Realm, c)
}

// RefreshInput es la entrada de refresh.
type RefreshInput struct {
	Realm        string
	RefreshToken string
}

// Refresh rota el refresh token (detecta reuse) y emite un nuevo par.
func (s *Service) Refresh(ctx context.Context, in RefreshInput) (*TokenPair, error) {
	hash := token.HashRefreshToken(in.RefreshToken)
	sess, err := s.store.GetSessionByRefresh(ctx, in.Realm, hash)
	if errors.Is(err, store.ErrNotFound) {
		return nil, screrrors.New(screrrors.CodeUnauthorized, "invalid refresh token")
	}
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeInternal, "lookup session", err)
	}
	if sess.RevokedAt != nil {
		// Reuse: refresh ya rotado. TODO: revocar la cadena hacia adelante.
		return nil, screrrors.New(screrrors.CodeUnauthorized, "refresh token already rotated")
	}
	if time.Now().After(sess.ExpiresAt) {
		return nil, screrrors.New(screrrors.CodeUnauthorized, "refresh token expired")
	}

	c, err := s.store.GetCitizenByID(ctx, in.Realm, sess.CitizenID)
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeInternal, "lookup citizen", err)
	}

	jti, err := s.store.NextJTI(ctx, in.Realm)
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeInternal, "next jti", err)
	}
	access, accessExp, err := s.signer.IssueAccess(c.ID, in.Realm, "", []string{"CITIZEN"}, jti)
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeInternal, "sign access", err)
	}
	newRefresh, err := token.NewRefreshToken()
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeInternal, "gen refresh", err)
	}
	newExp := time.Now().Add(s.refreshTTL)
	if _, err := s.store.RotateSession(ctx, in.Realm, sess.ID, token.HashRefreshToken(newRefresh), jti, newExp, c.ID); err != nil {
		return nil, screrrors.Wrap(screrrors.CodeInternal, "rotate session", err)
	}
	return &TokenPair{
		AccessToken: access, AccessExpiresAt: accessExp,
		RefreshToken: newRefresh, RefreshExpiresAt: newExp,
		Citizen: *citizenView(c),
	}, nil
}

// Logout revoca la sesión asociada al refresh token y agrega el JTI a
// revoked_jtis para que el gateway lo bloquee en su Roaring bitmap.
func (s *Service) Logout(ctx context.Context, realm, refreshToken string) error {
	hash := token.HashRefreshToken(refreshToken)
	sess, err := s.store.GetSessionByRefresh(ctx, realm, hash)
	if errors.Is(err, store.ErrNotFound) {
		return nil // idempotente
	}
	if err != nil {
		return screrrors.Wrap(screrrors.CodeInternal, "lookup session", err)
	}
	return s.store.RevokeSession(ctx, realm, sess.ID)
}

// RevokedSnapshot devuelve los JTIs revocados aún vigentes — consumido por el
// gateway para reconstruir su Roaring bitmap.
func (s *Service) RevokedSnapshot(ctx context.Context, realm string) ([]int64, error) {
	return s.store.RevokedSnapshot(ctx, realm)
}

func (s *Service) issueTokenPair(ctx context.Context, realm string, c *store.Citizen) (*TokenPair, error) {
	jti, err := s.store.NextJTI(ctx, realm)
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeInternal, "next jti", err)
	}
	access, accessExp, err := s.signer.IssueAccess(c.ID, realm, "", []string{"CITIZEN"}, jti)
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeInternal, "sign access", err)
	}
	refresh, err := token.NewRefreshToken()
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeInternal, "gen refresh", err)
	}
	exp := time.Now().Add(s.refreshTTL)
	if _, err := s.store.CreateSession(ctx, realm, c.ID, token.HashRefreshToken(refresh), jti, exp); err != nil {
		return nil, screrrors.Wrap(screrrors.CodeInternal, "create session", err)
	}
	return &TokenPair{
		AccessToken: access, AccessExpiresAt: accessExp,
		RefreshToken: refresh, RefreshExpiresAt: exp,
		Citizen: *citizenView(c),
	}, nil
}

func citizenView(c *store.Citizen) *CitizenView {
	return &CitizenView{
		ID: c.ID, Cedula: c.Cedula, Email: c.Email, Status: c.Status, CreatedAt: c.CreatedAt,
	}
}
