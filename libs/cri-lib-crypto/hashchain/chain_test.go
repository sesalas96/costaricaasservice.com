package hashchain

import (
	"bytes"
	"testing"
)

func TestEntryHashDeterministic(t *testing.T) {
	prev := Genesis("demo")
	p1 := map[string]any{"id": "a", "ts": 1, "actor": "x"}
	p2 := map[string]any{"actor": "x", "ts": 1, "id": "a"} // misma carga, distinto orden

	h1, err := EntryHash(prev, p1)
	if err != nil {
		t.Fatal(err)
	}
	h2, err := EntryHash(prev, p2)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(h1, h2) {
		t.Errorf("EntryHash no es determinístico: %x vs %x", h1, h2)
	}
}

func TestVerifyChain(t *testing.T) {
	const realm = "demo"
	prev := Genesis(realm)

	mk := func(payload map[string]any) Entry {
		h, err := EntryHash(prev, payload)
		if err != nil {
			t.Fatal(err)
		}
		e := Entry{PrevHash: append([]byte(nil), prev...), EntryHash: h, Payload: payload}
		prev = h
		return e
	}

	chain := []Entry{
		mk(map[string]any{"id": "a"}),
		mk(map[string]any{"id": "b"}),
		mk(map[string]any{"id": "c"}),
	}

	if idx, err := Verify(realm, chain); idx != -1 || err != nil {
		t.Fatalf("Verify failed: idx=%d err=%v", idx, err)
	}

	// Tamper en el medio.
	chain[1].Payload["id"] = "tampered"
	if idx, err := Verify(realm, chain); idx != 1 || err == nil {
		t.Fatalf("expected tamper detected at idx=1, got idx=%d err=%v", idx, err)
	}
}

func TestMerkleRoot(t *testing.T) {
	leaves := [][]byte{[]byte("a"), []byte("b"), []byte("c")}
	root := MerkleRoot(leaves)
	if len(root) != 32 {
		t.Errorf("root size = %d, want 32", len(root))
	}
	if MerkleRoot(nil) != nil {
		t.Error("empty merkle should be nil")
	}
}
