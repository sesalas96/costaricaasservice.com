package hashchain

import (
	"crypto/sha256"
)

// MerkleRoot calcula el root de un Merkle tree binario sobre los hashes
// dados (típicamente todos los entry_hash de un epoch). Usa duplicación
// del último elemento si el nivel tiene número impar (RFC 6962 simplificado).
//
// Devuelve nil si leaves está vacío.
func MerkleRoot(leaves [][]byte) []byte {
	if len(leaves) == 0 {
		return nil
	}
	level := make([][]byte, len(leaves))
	copy(level, leaves)
	for len(level) > 1 {
		next := make([][]byte, 0, (len(level)+1)/2)
		for i := 0; i < len(level); i += 2 {
			if i+1 == len(level) {
				next = append(next, hashPair(level[i], level[i]))
			} else {
				next = append(next, hashPair(level[i], level[i+1]))
			}
		}
		level = next
	}
	return level[0]
}

func hashPair(a, b []byte) []byte {
	h := sha256.New()
	h.Write(a)
	h.Write(b)
	return h.Sum(nil)
}
