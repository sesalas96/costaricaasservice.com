// Package store del cri-svc-registro-civil. Implementación in-memory para
// MVP demo (sin Postgres). Cuando se levante DB se reemplaza por un store
// pgx con el patrón schema-per-realm que ya usa iduc-identity.
package store

import (
	"sync"
	"time"
)

// Person es la fila del registro civil de un realm.
type Person struct {
	Cedula    string    `json:"cedula"`
	FullName  string    `json:"fullName"`
	BirthDate string    `json:"birthDate"` // YYYY-MM-DD
	Address   string    `json:"address"`
	Email     string    `json:"email"`
	Status    string    `json:"status"` // alive | deceased
	UpdatedAt time.Time `json:"updatedAt"`
}

// VitalEvent es un evento vital (nacimiento, matrimonio, defunción, etc).
type VitalEvent struct {
	ID        string    `json:"id"`
	Cedula    string    `json:"cedula"`
	Type      string    `json:"type"` // birth | marriage | divorce | death | name_change
	Date      string    `json:"date"`
	Notes     string    `json:"notes"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Store es un store in-memory thread-safe scopeado por realm.
type Store struct {
	mu      sync.RWMutex
	persons map[string]map[string]*Person      // realm → cedula → Person
	events  map[string]map[string][]*VitalEvent // realm → cedula → events
}

// New crea un Store vacío.
func New() *Store {
	return &Store{
		persons: make(map[string]map[string]*Person),
		events:  make(map[string]map[string][]*VitalEvent),
	}
}

// Seed inserta data mock en un realm. Idempotente.
func (s *Store) Seed(realm string, persons []Person, events []VitalEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.persons[realm]; !ok {
		s.persons[realm] = make(map[string]*Person)
		s.events[realm] = make(map[string][]*VitalEvent)
	}
	for i := range persons {
		p := persons[i]
		p.UpdatedAt = time.Now()
		s.persons[realm][p.Cedula] = &p
	}
	for i := range events {
		e := events[i]
		e.UpdatedAt = time.Now()
		s.events[realm][e.Cedula] = append(s.events[realm][e.Cedula], &e)
	}
}

// GetByCedula busca una persona en el realm.
func (s *Store) GetByCedula(realm, cedula string) *Person {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if r, ok := s.persons[realm]; ok {
		return r[cedula]
	}
	return nil
}

// EventsByCedula retorna los eventos vitales de una persona.
func (s *Store) EventsByCedula(realm, cedula string) []*VitalEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if r, ok := s.events[realm]; ok {
		out := make([]*VitalEvent, len(r[cedula]))
		copy(out, r[cedula])
		return out
	}
	return nil
}
