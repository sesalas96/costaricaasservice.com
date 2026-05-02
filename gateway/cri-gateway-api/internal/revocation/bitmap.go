// Package revocation mantiene un Roaring bitmap por realm con los JTIs
// (numéricos) revocados. El JWT verifier consulta el bitmap en O(log n)
// antes de aceptar un token; si el JTI está, el token se rechaza.
//
// El bitmap se rellena periódicamente desde el snapshot que expone
// cri-svc-iduc-identity en GET /internal/revoked-jti/snapshot?realm=<slug>.
//
// Diseño thread-safe: cada actualización reemplaza el bitmap completo
// detrás de un sync.RWMutex. Las consultas son lock-free después del
// snapshot del puntero (estilo copy-on-write del bitmap).
package revocation

import (
	"sync"

	"github.com/RoaringBitmap/roaring/v2"
)

// Set representa el conjunto de JTIs revocados de un realm.
type Set struct {
	mu  sync.RWMutex
	bm  *roaring.Bitmap
}

// NewSet crea un Set vacío.
func NewSet() *Set { return &Set{bm: roaring.New()} }

// Replace sustituye atómicamente el bitmap por uno construido desde la lista
// de JTIs dada. Usar tras cada poll de snapshot.
func (s *Set) Replace(jtis []int64) {
	bm := roaring.New()
	for _, j := range jtis {
		if j > 0 && j <= 1<<32-1 {
			bm.Add(uint32(j))
		}
	}
	s.mu.Lock()
	s.bm = bm
	s.mu.Unlock()
}

// Contains retorna true si el JTI está revocado.
func (s *Set) Contains(jti uint32) bool {
	s.mu.RLock()
	bm := s.bm
	s.mu.RUnlock()
	return bm.Contains(jti)
}

// Cardinality retorna la cantidad de elementos actualmente en el bitmap.
// Útil para métricas / observabilidad.
func (s *Set) Cardinality() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.bm.GetCardinality()
}

// Registry mantiene un Set por realm con concurrencia segura.
type Registry struct {
	mu   sync.RWMutex
	sets map[string]*Set
}

// NewRegistry crea un Registry vacío.
func NewRegistry() *Registry { return &Registry{sets: make(map[string]*Set)} }

// Get devuelve el Set del realm, creándolo en la primera consulta.
func (r *Registry) Get(realm string) *Set {
	r.mu.RLock()
	if s, ok := r.sets[realm]; ok {
		r.mu.RUnlock()
		return s
	}
	r.mu.RUnlock()

	r.mu.Lock()
	defer r.mu.Unlock()
	if s, ok := r.sets[realm]; ok {
		return s
	}
	s := NewSet()
	r.sets[realm] = s
	return s
}

// Realms retorna los realms conocidos por el registry.
func (r *Registry) Realms() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]string, 0, len(r.sets))
	for k := range r.sets {
		out = append(out, k)
	}
	return out
}
