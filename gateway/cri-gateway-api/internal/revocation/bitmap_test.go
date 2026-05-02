package revocation

import (
	"sync"
	"testing"
)

func TestSetReplaceContains(t *testing.T) {
	s := NewSet()
	if s.Contains(42) {
		t.Error("empty set should not contain anything")
	}
	s.Replace([]int64{1, 2, 100, 999})
	for _, want := range []uint32{1, 2, 100, 999} {
		if !s.Contains(want) {
			t.Errorf("expected to contain %d", want)
		}
	}
	if s.Contains(3) {
		t.Error("should not contain 3")
	}
	if s.Cardinality() != 4 {
		t.Errorf("cardinality = %d, want 4", s.Cardinality())
	}

	// Replace shrinks the set.
	s.Replace([]int64{5})
	if s.Contains(1) {
		t.Error("after Replace, old entries should be gone")
	}
	if !s.Contains(5) {
		t.Error("after Replace, new entry should be present")
	}
}

func TestRegistryGet(t *testing.T) {
	r := NewRegistry()
	a := r.Get("cr-prod")
	b := r.Get("cr-prod")
	if a != b {
		t.Error("Get should return same instance per realm")
	}
	c := r.Get("sv-prod")
	if a == c {
		t.Error("different realms should have different Sets")
	}
}

func TestSetConcurrent(t *testing.T) {
	s := NewSet()
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				s.Contains(uint32(j))
			}
		}()
	}
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				s.Replace([]int64{int64(j + i*100)})
			}
		}(i)
	}
	wg.Wait()
}
