// Package store del cri-svc-migracion. In-memory para MVP.
//
// Once-only: Migración guarda SOLO datos migratorios (status, nacionalidad
// origen, vigencia DIMEX, movimientos). Nombre/domicilio del titular se
// piden a registro-civil vía interop.
package store

import (
	"sync"
	"time"
)

type ImmigrationStatus struct {
	Cedula       string    `json:"cedula"` // o DIMEX para extranjeros
	Nationality  string    `json:"nationality"`
	Category     string    `json:"category"` // ciudadano | residente_permanente | residente_temporal | refugiado | turista
	DocumentType string    `json:"documentType"` // cedula_nacional | dimex | pasaporte
	IssuedAt     string    `json:"issuedAt"`
	ExpiresAt    string    `json:"expiresAt,omitempty"`
	Status       string    `json:"status"` // active | expired | cancelled
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Movement struct {
	ID         string    `json:"id"`
	Cedula     string    `json:"cedula"`
	BorderPost string    `json:"borderPost"` // SJO | LIR | PEN | PASO_CANOAS | PEÑAS_BLANCAS
	Direction  string    `json:"direction"`  // entry | exit
	When       string    `json:"when"`       // ISO 8601
	UpdatedAt  time.Time `json:"updatedAt"`
}

type Store struct {
	mu        sync.RWMutex
	statuses  map[string]map[string]*ImmigrationStatus // realm → cedula → status
	movements map[string]map[string][]*Movement        // realm → cedula → list
}

func New() *Store {
	return &Store{
		statuses:  map[string]map[string]*ImmigrationStatus{},
		movements: map[string]map[string][]*Movement{},
	}
}

func (s *Store) Seed(realm string, statuses []ImmigrationStatus, movements []Movement) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.statuses[realm]; !ok {
		s.statuses[realm] = map[string]*ImmigrationStatus{}
		s.movements[realm] = map[string][]*Movement{}
	}
	for i := range statuses {
		st := statuses[i]
		st.UpdatedAt = time.Now()
		s.statuses[realm][st.Cedula] = &st
	}
	for i := range movements {
		m := movements[i]
		m.UpdatedAt = time.Now()
		s.movements[realm][m.Cedula] = append(s.movements[realm][m.Cedula], &m)
	}
}

func (s *Store) StatusByCedula(realm, cedula string) *ImmigrationStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.statuses[realm][cedula]
}

func (s *Store) MovementsByCedula(realm, cedula string) []*Movement {
	s.mu.RLock()
	defer s.mu.RUnlock()
	src := s.movements[realm][cedula]
	out := make([]*Movement, len(src))
	copy(out, src)
	return out
}
