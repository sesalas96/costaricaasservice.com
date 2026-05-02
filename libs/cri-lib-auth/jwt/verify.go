// Package jwt verifica tokens RS256 emitidos por cri-svc-iduc-identity.
//
// Convenciones del JWT en saascr:
//   - alg: RS256
//   - sub: citizen_id u operator_id (string)
//   - jti: numérico (sequence) — usado en Roaring bitmap del gateway
//   - realm: slug del realm activo
//   - roles: array de strings (Role values)
//   - member: optional, member_id si aplica
package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/devsebas/saascr/libs/cri-lib-auth/auth"
)

// Claims son las pretensiones esperadas en cada token saascr.
type Claims struct {
	Realm  string   `json:"realm"`
	Roles  []string `json:"roles"`
	Member string   `json:"member,omitempty"`
	jwt.RegisteredClaims
}

// Verifier valida tokens RS256 con una clave pública conocida.
type Verifier struct {
	pub *rsa.PublicKey
}

// NewVerifier construye un Verifier con la clave pública dada.
func NewVerifier(pub *rsa.PublicKey) *Verifier { return &Verifier{pub: pub} }

// Verify parsea, valida la firma y caducidad, y retorna el Principal.
// El llamador (gateway) es responsable de comprobar la revocación por JTI
// contra su Roaring bitmap.
func (v *Verifier) Verify(tokenStr string) (*auth.Principal, error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected alg: %v", t.Header["alg"])
		}
		return v.pub, nil
	}, jwt.WithExpirationRequired(), jwt.WithLeeway(30*time.Second))
	if err != nil {
		return nil, err
	}
	claims, ok := tok.Claims.(*Claims)
	if !ok || !tok.Valid {
		return nil, errors.New("invalid token")
	}
	if claims.Realm == "" {
		return nil, errors.New("token missing realm claim")
	}
	roles := make([]auth.Role, 0, len(claims.Roles))
	for _, r := range claims.Roles {
		roles = append(roles, auth.Role(r))
	}
	return &auth.Principal{
		Sub:    claims.Subject,
		Roles:  roles,
		Realm:  claims.Realm,
		Member: claims.Member,
		JTI:    claims.ID,
	}, nil
}

// JTINum interpreta el JTI como número (para Roaring bitmap del gateway).
// Retorna 0 si no es numérico.
func JTINum(jti string) uint32 {
	n, err := strconv.ParseUint(jti, 10, 32)
	if err != nil {
		return 0
	}
	return uint32(n)
}
