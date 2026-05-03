// Package store del cri-svc-ins. In-memory para MVP.
//
// Once-only: el INS guarda SOLO datos de aseguramiento (pólizas, primas,
// vigencias). Nombre/domicilio del titular se piden a registro-civil vía
// interop, NO se duplican aquí.
package store

import (
	"sync"
	"time"
)

type Policy struct {
	Number        string    `json:"number"`
	HolderCedula  string    `json:"holderCedula"`
	Kind          string    `json:"kind"` // soa | riesgos_trabajo | voluntario_vida | voluntario_salud
	VehiclePlate  string    `json:"vehiclePlate,omitempty"`
	PremiumColons int64     `json:"premiumColons"` // anual
	ValidFrom     string    `json:"validFrom"`
	ValidUntil    string    `json:"validUntil"`
	Status        string    `json:"status"` // active | expired | cancelled
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Store struct {
	mu       sync.RWMutex
	policies map[string]map[string][]*Policy // realm → cedula → list
}

func New() *Store {
	return &Store{policies: map[string]map[string][]*Policy{}}
}

func (s *Store) Seed(realm string, policies []Policy) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.policies[realm]; !ok {
		s.policies[realm] = map[string][]*Policy{}
	}
	for i := range policies {
		p := policies[i]
		p.UpdatedAt = time.Now()
		s.policies[realm][p.HolderCedula] = append(s.policies[realm][p.HolderCedula], &p)
	}
}

func (s *Store) PoliciesByCedula(realm, cedula string) []*Policy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	src := s.policies[realm][cedula]
	out := make([]*Policy, len(src))
	copy(out, src)
	return out
}
