// Package store del cri-svc-mep. In-memory para MVP.
//
// Once-only: el MEP guarda SOLO datos académicos (matrícula, certificados,
// historial de notas). Nombre/domicilio del estudiante se piden a
// registro-civil vía interop, NO se duplican aquí.
package store

import (
	"sync"
	"time"
)

type Enrollment struct {
	ID            string    `json:"id"`
	StudentCedula string    `json:"studentCedula"`
	Year          int       `json:"year"`
	Grade         string    `json:"grade"`         // 1ro_primaria | 6to_secundaria | etc.
	SchoolName    string    `json:"schoolName"`
	SchoolKind    string    `json:"schoolKind"`    // publico | privado | subvencionado
	Status        string    `json:"status"`        // active | completed | dropped
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Certificate struct {
	ID            string    `json:"id"`
	StudentCedula string    `json:"studentCedula"`
	Kind          string    `json:"kind"` // primaria_aprobada | bachillerato | titulo_tecnico
	IssuedAt      string    `json:"issuedAt"`
	IssuedBy      string    `json:"issuedBy"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Store struct {
	mu           sync.RWMutex
	enrollments  map[string]map[string][]*Enrollment  // realm → cedula → list
	certificates map[string]map[string][]*Certificate // realm → cedula → list
}

func New() *Store {
	return &Store{
		enrollments:  map[string]map[string][]*Enrollment{},
		certificates: map[string]map[string][]*Certificate{},
	}
}

func (s *Store) Seed(realm string, enrollments []Enrollment, certificates []Certificate) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.enrollments[realm]; !ok {
		s.enrollments[realm] = map[string][]*Enrollment{}
		s.certificates[realm] = map[string][]*Certificate{}
	}
	for i := range enrollments {
		e := enrollments[i]
		e.UpdatedAt = time.Now()
		s.enrollments[realm][e.StudentCedula] = append(s.enrollments[realm][e.StudentCedula], &e)
	}
	for i := range certificates {
		c := certificates[i]
		c.UpdatedAt = time.Now()
		s.certificates[realm][c.StudentCedula] = append(s.certificates[realm][c.StudentCedula], &c)
	}
}

func (s *Store) EnrollmentsByCedula(realm, cedula string) []*Enrollment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	src := s.enrollments[realm][cedula]
	out := make([]*Enrollment, len(src))
	copy(out, src)
	return out
}

func (s *Store) CertificatesByCedula(realm, cedula string) []*Certificate {
	s.mu.RLock()
	defer s.mu.RUnlock()
	src := s.certificates[realm][cedula]
	out := make([]*Certificate, len(src))
	copy(out, src)
	return out
}
