// Package signing firma y verifica los payloads inter-member usando una
// firma "JWS detached"-style: el header y el signature se entregan junto
// al payload en un envelope; el payload no se base64-encodea ni viaja
// dentro de la firma. Esto permite que el body original cruce inalterado.
//
// Algoritmo: RS256 (RSA-3072+ SHA-256) o ES256 (ECDSA P-256 / P-384) según
// el tipo de clave del member. Para MVP soportamos solo RS256; ES256 se
// agregará cuando algún member exponga una clave EC.
//
// Formato del envelope:
//
//	{
//	  "header": {
//	    "alg": "RS256",
//	    "kid": "<member_slug>",
//	    "iat": 1714657200,
//	    "request_id": "01ULID...",
//	    "realm": "demo",
//	    "src_member": "hacienda",
//	    "dst_member": "registro-civil"
//	  },
//	  "payload": <opaque JSON, lo que el caller envió como body>,
//	  "signature": "<base64url>"
//	}
//
// Lo firmado es: SHA-256( canonical(header) || "." || canonical(payload) ).
package signing

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"sort"
)

const algRS256 = "RS256"

// Header son las pretensiones de un envelope firmado.
type Header struct {
	Alg       string `json:"alg"`
	Kid       string `json:"kid"`
	Iat       int64  `json:"iat"`
	RequestID string `json:"request_id"`
	Realm     string `json:"realm"`
	SrcMember string `json:"src_member"`
	DstMember string `json:"dst_member"`
	Service   string `json:"service"`
	Version   string `json:"version"`
}

// Envelope es lo que viaja en wire entre security-servers.
type Envelope struct {
	Header    Header          `json:"header"`
	Payload   json.RawMessage `json:"payload"`
	Signature string          `json:"signature"`
}

// Signer firma payloads con la clave privada del member local.
type Signer struct {
	priv  *rsa.PrivateKey
	kid   string // member slug
}

// NewSignerFromPEM carga la clave privada desde un archivo PEM PKCS#8/PKCS#1.
func NewSignerFromPEM(path, kid string) (*Signer, error) {
	pemBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read pem: %w", err)
	}
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("signing: invalid PEM")
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		k, err2 := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err2 != nil {
			return nil, fmt.Errorf("parse private key: %v / %v", err, err2)
		}
		parsed = k
	}
	priv, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("signing: only RSA keys supported in MVP")
	}
	if priv.N.BitLen() < 2048 {
		return nil, errors.New("signing: RSA key must be >= 2048 bits")
	}
	return &Signer{priv: priv, kid: kid}, nil
}

// Sign produce un Envelope firmado a partir de un header (sin alg/kid/iat)
// y un payload arbitrario (json.RawMessage).
func (s *Signer) Sign(header Header, payload json.RawMessage) (*Envelope, error) {
	header.Alg = algRS256
	header.Kid = s.kid
	if header.Iat == 0 {
		header.Iat = nowUnix()
	}
	digest, err := computeDigest(header, payload)
	if err != nil {
		return nil, err
	}
	sig, err := rsa.SignPKCS1v15(rand.Reader, s.priv, crypto.SHA256, digest)
	if err != nil {
		return nil, err
	}
	return &Envelope{
		Header:    header,
		Payload:   payload,
		Signature: base64.RawURLEncoding.EncodeToString(sig),
	}, nil
}

// Verify valida el Envelope con la clave pública dada.
func Verify(env *Envelope, pub any) error {
	if env.Header.Alg != algRS256 {
		return fmt.Errorf("signing: unsupported alg %q", env.Header.Alg)
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return errors.New("signing: not an RSA public key")
	}
	sig, err := base64.RawURLEncoding.DecodeString(env.Signature)
	if err != nil {
		return fmt.Errorf("decode signature: %w", err)
	}
	digest, err := computeDigest(env.Header, env.Payload)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, digest, sig)
}

// computeDigest = SHA-256( canonical(header) || "." || canonical(payload) ).
func computeDigest(header Header, payload json.RawMessage) ([]byte, error) {
	hb, err := canonicalJSON(header)
	if err != nil {
		return nil, err
	}
	pb := []byte(payload)
	if len(pb) == 0 {
		pb = []byte("null")
	}
	pb, err = canonicalRawJSON(pb)
	if err != nil {
		return nil, err
	}
	h := sha256.New()
	h.Write(hb)
	h.Write([]byte("."))
	h.Write(pb)
	return h.Sum(nil), nil
}

// canonicalJSON serializa una struct en forma determinística.
func canonicalJSON(v any) ([]byte, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return canonicalRawJSON(raw)
}

// canonicalRawJSON re-ordena alfabéticamente las claves de un objeto JSON
// (recursivamente) para producir una serialización determinística.
func canonicalRawJSON(raw []byte) ([]byte, error) {
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, err
	}
	return marshalSorted(v)
}

func marshalSorted(v any) ([]byte, error) {
	switch t := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		out := []byte{'{'}
		for i, k := range keys {
			if i > 0 {
				out = append(out, ',')
			}
			kb, err := json.Marshal(k)
			if err != nil {
				return nil, err
			}
			out = append(out, kb...)
			out = append(out, ':')
			vb, err := marshalSorted(t[k])
			if err != nil {
				return nil, err
			}
			out = append(out, vb...)
		}
		return append(out, '}'), nil
	case []any:
		out := []byte{'['}
		for i, e := range t {
			if i > 0 {
				out = append(out, ',')
			}
			eb, err := marshalSorted(e)
			if err != nil {
				return nil, err
			}
			out = append(out, eb...)
		}
		return append(out, ']'), nil
	default:
		return json.Marshal(t)
	}
}

func nowUnix() int64 { return timeNowFunc() }

// timeNowFunc se overridea en tests si hace falta.
var timeNowFunc = func() int64 { return now() }
