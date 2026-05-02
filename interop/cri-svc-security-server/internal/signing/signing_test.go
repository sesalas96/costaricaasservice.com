package signing

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
)

// makeTestSigner genera un par RSA-3072 in-memory y devuelve un Signer y la pubkey.
func makeTestSigner(t *testing.T, kid string) (*Signer, *rsa.PublicKey) {
	t.Helper()
	priv, err := rsa.GenerateKey(rand.Reader, 3072)
	if err != nil {
		t.Fatal(err)
	}
	tmp := filepath.Join(t.TempDir(), "priv.pem")
	der, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(tmp, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}), 0o600); err != nil {
		t.Fatal(err)
	}
	s, err := NewSignerFromPEM(tmp, kid)
	if err != nil {
		t.Fatal(err)
	}
	return s, &priv.PublicKey
}

func TestSignVerifyRoundtrip(t *testing.T) {
	signer, pub := makeTestSigner(t, "hacienda")

	header := Header{
		RequestID: "01ULID",
		Realm:     "demo",
		SrcMember: "hacienda",
		DstMember: "registro-civil",
		Service:   "persons.get",
		Version:   "v1",
	}
	payload := json.RawMessage(`{"cedula":"112345678"}`)

	env, err := signer.Sign(header, payload)
	if err != nil {
		t.Fatal(err)
	}
	if env.Header.Alg != "RS256" {
		t.Errorf("alg = %s", env.Header.Alg)
	}
	if env.Header.Kid != "hacienda" {
		t.Errorf("kid = %s", env.Header.Kid)
	}
	if env.Header.Iat == 0 {
		t.Error("iat not set")
	}
	if env.Signature == "" {
		t.Error("signature empty")
	}

	if err := Verify(env, pub); err != nil {
		t.Errorf("Verify roundtrip failed: %v", err)
	}
}

func TestVerifyTamperPayload(t *testing.T) {
	signer, pub := makeTestSigner(t, "x")
	env, _ := signer.Sign(Header{Realm: "demo"}, json.RawMessage(`{"a":1}`))

	env.Payload = json.RawMessage(`{"a":2}`) // tamper
	if err := Verify(env, pub); err == nil {
		t.Error("Verify must fail after tamper")
	}
}

func TestVerifyTamperHeader(t *testing.T) {
	signer, pub := makeTestSigner(t, "x")
	env, _ := signer.Sign(Header{Realm: "demo", SrcMember: "x"}, json.RawMessage(`{}`))

	env.Header.SrcMember = "y" // tamper
	if err := Verify(env, pub); err == nil {
		t.Error("Verify must fail after header tamper")
	}
}

func TestVerifyWrongKey(t *testing.T) {
	signer, _ := makeTestSigner(t, "x")
	env, _ := signer.Sign(Header{Realm: "demo"}, json.RawMessage(`{}`))

	otherPriv, _ := rsa.GenerateKey(rand.Reader, 3072)
	if err := Verify(env, &otherPriv.PublicKey); err == nil {
		t.Error("Verify must fail with wrong key")
	}
}
