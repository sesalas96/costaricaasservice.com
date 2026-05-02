// Package hashchain implementa la hash chain de ADR-0002 para el audit log
// inter-member. Cada entry encadena el SHA-256 de su predecesor; cualquier
// modificación retroactiva rompe la cadena de todos los registros posteriores.
package hashchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sort"
)

// Genesis devuelve el hash inicial de un realm: SHA-256("genesis::<realm>").
// Cada realm tiene su propia cadena, anclada en este valor determinístico.
func Genesis(realm string) []byte {
	h := sha256.Sum256([]byte("genesis::" + realm))
	return h[:]
}

// Canonicalize serializa un mapa en JSON canonical (claves ordenadas, sin
// espacios, escapes mínimos) para producir un input determinístico al hash.
func Canonicalize(payload map[string]any) ([]byte, error) {
	keys := make([]string, 0, len(payload))
	for k := range payload {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, k := range keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		kb, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}
		buf.Write(kb)
		buf.WriteByte(':')
		vb, err := json.Marshal(payload[k])
		if err != nil {
			return nil, err
		}
		buf.Write(vb)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// EntryHash calcula SHA-256(prevHash || canonical(payload)).
// payload no debe incluir los campos prev_hash ni entry_hash.
func EntryHash(prevHash []byte, payload map[string]any) ([]byte, error) {
	canon, err := Canonicalize(payload)
	if err != nil {
		return nil, err
	}
	h := sha256.New()
	h.Write(prevHash)
	h.Write(canon)
	return h.Sum(nil), nil
}

// Verify recorre una secuencia ordenada de entries y valida que cada uno
// referencie al hash de su predecesor.
//
// Retorna el índice del primer entry inválido, o -1 si toda la cadena es
// válida.
type Entry struct {
	PrevHash  []byte
	EntryHash []byte
	Payload   map[string]any
}

// Verify valida una cadena entera contra el genesis del realm.
func Verify(realm string, entries []Entry) (firstInvalidIdx int, err error) {
	expectedPrev := Genesis(realm)
	for i, e := range entries {
		if !bytes.Equal(e.PrevHash, expectedPrev) {
			return i, errors.New("hashchain: prev_hash mismatch")
		}
		got, err := EntryHash(e.PrevHash, e.Payload)
		if err != nil {
			return i, err
		}
		if !bytes.Equal(got, e.EntryHash) {
			return i, errors.New("hashchain: entry_hash mismatch")
		}
		expectedPrev = e.EntryHash
	}
	return -1, nil
}

// Hex devuelve la representación hexadecimal de un hash.
func Hex(h []byte) string { return hex.EncodeToString(h) }
