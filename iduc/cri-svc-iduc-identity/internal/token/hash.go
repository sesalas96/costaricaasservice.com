package token

import "crypto/sha256"

// HashRefreshToken es lo que se guarda en sessions.refresh_token_hash.
// Nunca guardamos el refresh token en claro.
func HashRefreshToken(refresh string) []byte {
	h := sha256.Sum256([]byte(refresh))
	return h[:]
}
