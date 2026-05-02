// Package signing provee helpers Ed25519 y la interfaz Signer.
//
// IMPORTANTE: en producción la firma con clave de ciudadano se delega al KMS
// (Vault Transit / AWS KMS). LocalSigner es solo para tests y desarrollo.
// Ver ADR-0003 (KMS y custodia de claves).
package signing

import (
	"crypto/ed25519"
	"errors"
	"io"
)

// Signer es el contrato que implementan los firmadores: KMS, Vault, Local.
// Una implementación KMS no tiene acceso al privKey y usa el key_id en su lugar.
type Signer interface {
	Sign(msg []byte) ([]byte, error)
	PublicKey() ed25519.PublicKey
}

// Verify valida una firma Ed25519. Retorna nil si es válida.
func Verify(pub ed25519.PublicKey, msg, sig []byte) error {
	if len(pub) != ed25519.PublicKeySize {
		return errors.New("ed25519: invalid public key size")
	}
	if len(sig) != ed25519.SignatureSize {
		return errors.New("ed25519: invalid signature size")
	}
	if !ed25519.Verify(pub, msg, sig) {
		return errors.New("ed25519: invalid signature")
	}
	return nil
}

// LocalSigner es un Signer en memoria. SOLO para tests y dev.
type LocalSigner struct {
	priv ed25519.PrivateKey
	pub  ed25519.PublicKey
}

// NewLocalSigner genera un par nuevo.
func NewLocalSigner(rand io.Reader) (*LocalSigner, error) {
	pub, priv, err := ed25519.GenerateKey(rand)
	if err != nil {
		return nil, err
	}
	return &LocalSigner{priv: priv, pub: pub}, nil
}

// Sign firma el mensaje con la clave local.
func (s *LocalSigner) Sign(msg []byte) ([]byte, error) {
	return ed25519.Sign(s.priv, msg), nil
}

// PublicKey retorna la clave pública.
func (s *LocalSigner) PublicKey() ed25519.PublicKey { return s.pub }
