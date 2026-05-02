// Package passwords implementa hashing y verificación de contraseñas con
// argon2id, siguiendo el formato encoded estándar:
//
//	$argon2id$v=19$m=<memKiB>,t=<time>,p=<parallelism>$<saltB64>$<hashB64>
//
// Parámetros calibrados para servidores típicos en 2026 (~100ms en CPU
// moderna):  m=64MiB, t=3, p=2, saltLen=16, hashLen=32.
package passwords

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	memoryKiB   = 64 * 1024 // 64 MiB
	iterations  = 3
	parallelism = 2
	saltLen     = 16
	hashLen     = 32
)

var (
	ErrInvalidEncoded = errors.New("passwords: invalid encoded hash")
	ErrMismatch       = errors.New("passwords: password mismatch")
)

// Hash genera el encoded hash argon2id de la contraseña.
func Hash(password string) (string, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	h := argon2.IDKey([]byte(password), salt, iterations, memoryKiB, parallelism, hashLen)
	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, memoryKiB, iterations, parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(h),
	), nil
}

// Verify compara una contraseña contra un encoded hash. constant time.
func Verify(password, encoded string) error {
	parts := strings.Split(encoded, "$")
	// "" "argon2id" "v=19" "m=...,t=...,p=..." "<salt>" "<hash>"
	if len(parts) != 6 || parts[1] != "argon2id" {
		return ErrInvalidEncoded
	}
	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil || version != argon2.Version {
		return ErrInvalidEncoded
	}
	var mem, time, par uint32
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &mem, &time, &par); err != nil {
		return ErrInvalidEncoded
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return ErrInvalidEncoded
	}
	expected, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return ErrInvalidEncoded
	}
	got := argon2.IDKey([]byte(password), salt, time, mem, uint8(par), uint32(len(expected)))
	if subtle.ConstantTimeCompare(got, expected) != 1 {
		return ErrMismatch
	}
	return nil
}
