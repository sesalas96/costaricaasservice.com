// Package token firma JWT RS256 con la clave privada del servicio.
//
// Convención de claims (ver libs/cri-lib-auth/jwt para verificación):
//   - alg: RS256
//   - sub:    citizen_id (string UUID)
//   - jti:    sequence numérico (string) — usado en Roaring bitmap del gateway
//   - realm:  slug del realm activo
//   - roles:  array de strings (RoleCitizen por defecto)
//   - member: optional, member_id si aplica (vacío para citizen)
//   - exp/iat según política (access 15min / refresh 30d)
package token

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Signer firma access y refresh tokens con la clave privada RS256.
type Signer struct {
	priv          *rsa.PrivateKey
	issuer        string
	accessTTL     time.Duration
	refreshTTL    time.Duration
	keyID         string
}

// Config son los parámetros para el Signer.
type Config struct {
	PrivatePEMPath string        // ruta del PEM PKCS#8 con la llave privada RSA
	Issuer         string        // sub-iss del JWT (e.g., "cri-svc-iduc-identity")
	AccessTTL      time.Duration // 15min recomendado
	RefreshTTL     time.Duration // 30d recomendado
	KeyID          string        // kid para rotación; opcional
}

// NewSigner carga la PEM y construye el Signer.
func NewSigner(cfg Config) (*Signer, error) {
	pemBytes, err := os.ReadFile(cfg.PrivatePEMPath)
	if err != nil {
		return nil, fmt.Errorf("read pem: %w", err)
	}
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("token: invalid PEM")
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// fallback PKCS#1
		k, err2 := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err2 != nil {
			return nil, fmt.Errorf("parse private key: %v / %v", err, err2)
		}
		parsed = k
	}
	priv, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("token: not an RSA private key")
	}
	if priv.N.BitLen() < 2048 {
		return nil, errors.New("token: RSA key must be >= 2048 bits")
	}
	return &Signer{
		priv:       priv,
		issuer:     cfg.Issuer,
		accessTTL:  cfg.AccessTTL,
		refreshTTL: cfg.RefreshTTL,
		keyID:      cfg.KeyID,
	}, nil
}

// IssueAccess firma un access token.
func (s *Signer) IssueAccess(citizenID, realm, member string, roles []string, jtiNum int64) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(s.accessTTL)
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims{
		Realm:  realm,
		Roles:  roles,
		Member: member,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   citizenID,
			ID:        strconv.FormatInt(jtiNum, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	})
	if s.keyID != "" {
		tok.Header["kid"] = s.keyID
	}
	signed, err := tok.SignedString(s.priv)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, exp, nil
}

// NewRefreshToken genera un refresh token aleatorio (no es JWT — es un opaque
// token cuyo SHA-256 se guarda en sessions.refresh_token_hash).
func NewRefreshToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64URLEncode(buf), nil
}

// PublicKey expone la pubkey para que el gateway pueda verificar.
func (s *Signer) PublicKey() *rsa.PublicKey { return &s.priv.PublicKey }

type accessClaims struct {
	Realm  string   `json:"realm"`
	Roles  []string `json:"roles"`
	Member string   `json:"member,omitempty"`
	jwt.RegisteredClaims
}

// base64URLEncode evita importar encoding/base64 si solo se usa una vez.
func base64URLEncode(b []byte) string {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	out := make([]byte, 0, (len(b)*4+2)/3)
	for i := 0; i < len(b); i += 3 {
		var n uint32
		switch {
		case i+2 < len(b):
			n = uint32(b[i])<<16 | uint32(b[i+1])<<8 | uint32(b[i+2])
			out = append(out, alphabet[(n>>18)&0x3F], alphabet[(n>>12)&0x3F], alphabet[(n>>6)&0x3F], alphabet[n&0x3F])
		case i+1 < len(b):
			n = uint32(b[i])<<16 | uint32(b[i+1])<<8
			out = append(out, alphabet[(n>>18)&0x3F], alphabet[(n>>12)&0x3F], alphabet[(n>>6)&0x3F])
		default:
			n = uint32(b[i]) << 16
			out = append(out, alphabet[(n>>18)&0x3F], alphabet[(n>>12)&0x3F])
		}
	}
	return string(out)
}
