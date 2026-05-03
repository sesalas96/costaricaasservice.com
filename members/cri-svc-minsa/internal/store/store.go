// Package store del cri-svc-minsa. In-memory para MVP.
//
// Once-only: el MINSA guarda SOLO datos sanitarios (permisos de
// funcionamiento, autorizaciones a profesionales). Nombre/domicilio del
// titular se piden a registro-civil vía interop, NO se duplican aquí.
package store

import (
	"sync"
	"time"
)

type HealthPermit struct {
	Number       string    `json:"number"`
	HolderCedula string    `json:"holderCedula"`
	Kind         string    `json:"kind"` // farmacia | consultorio_medico | restaurante | salon_belleza | profesional_salud
	TradeName    string    `json:"tradeName,omitempty"`
	Address      string    `json:"address,omitempty"`
	IssuedAt     string    `json:"issuedAt"`
	ExpiresAt    string    `json:"expiresAt"`
	Status       string    `json:"status"` // active | expired | revoked
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Store struct {
	mu      sync.RWMutex
	permits map[string]map[string][]*HealthPermit // realm → cedula → list
}

func New() *Store {
	return &Store{permits: map[string]map[string][]*HealthPermit{}}
}

func (s *Store) Seed(realm string, permits []HealthPermit) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.permits[realm]; !ok {
		s.permits[realm] = map[string][]*HealthPermit{}
	}
	for i := range permits {
		p := permits[i]
		p.UpdatedAt = time.Now()
		s.permits[realm][p.HolderCedula] = append(s.permits[realm][p.HolderCedula], &p)
	}
}

func (s *Store) PermitsByCedula(realm, cedula string) []*HealthPermit {
	s.mu.RLock()
	defer s.mu.RUnlock()
	src := s.permits[realm][cedula]
	out := make([]*HealthPermit, len(src))
	copy(out, src)
	return out
}
