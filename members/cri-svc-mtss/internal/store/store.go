// Package store del cri-svc-mtss. In-memory para MVP.
//
// Once-only: el MTSS guarda SOLO datos laborales (contratos, planillas,
// salarios reportados). Nombre/domicilio del trabajador se piden a
// registro-civil vía interop, NO se duplican aquí.
package store

import (
	"sync"
	"time"
)

type Employment struct {
	ID            string    `json:"id"`
	WorkerCedula  string    `json:"workerCedula"`
	EmployerName  string    `json:"employerName"`
	EmployerCID   string    `json:"employerCID"` // cédula jurídica
	Position      string    `json:"position"`
	Contract      string    `json:"contract"` // indefinido | plazo_fijo | servicios_profesionales
	MonthlyColons int64     `json:"monthlyColons"`
	StartedAt     string    `json:"startedAt"`
	EndedAt       string    `json:"endedAt,omitempty"`
	Status        string    `json:"status"` // active | ended
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Store struct {
	mu          sync.RWMutex
	employments map[string]map[string][]*Employment // realm → cedula → list
}

func New() *Store {
	return &Store{employments: map[string]map[string][]*Employment{}}
}

func (s *Store) Seed(realm string, employments []Employment) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.employments[realm]; !ok {
		s.employments[realm] = map[string][]*Employment{}
	}
	for i := range employments {
		e := employments[i]
		e.UpdatedAt = time.Now()
		s.employments[realm][e.WorkerCedula] = append(s.employments[realm][e.WorkerCedula], &e)
	}
}

func (s *Store) EmploymentsByCedula(realm, cedula string) []*Employment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	src := s.employments[realm][cedula]
	out := make([]*Employment, len(src))
	copy(out, src)
	return out
}
