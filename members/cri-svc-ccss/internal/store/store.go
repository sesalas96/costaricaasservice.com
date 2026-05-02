// Package store del cri-svc-ccss. In-memory para MVP.
//
// Once-only: la CCSS guarda SOLO datos médicos propios (recetas, citas,
// historial). El nombre/domicilio del paciente se piden a registro-civil
// vía interop, NO se duplican aquí.
package store

import (
	"sync"
	"time"
)

type Prescription struct {
	ID         string    `json:"id"`
	Cedula     string    `json:"cedula"`
	Drug       string    `json:"drug"`
	Dosage     string    `json:"dosage"`
	IssuedAt   string    `json:"issuedAt"`
	IssuedBy   string    `json:"issuedBy"`
	ValidUntil string    `json:"validUntil"`
	Status     string    `json:"status"` // active | dispensed | expired
	UpdatedAt  time.Time `json:"updatedAt"`
}

type Appointment struct {
	ID        string    `json:"id"`
	Cedula    string    `json:"cedula"`
	Specialty string    `json:"specialty"`
	Hospital  string    `json:"hospital"`
	When      string    `json:"when"` // YYYY-MM-DD HH:MM
	Status    string    `json:"status"` // scheduled | completed | cancelled
	UpdatedAt time.Time `json:"updatedAt"`
}

type Store struct {
	mu            sync.RWMutex
	prescriptions map[string]map[string][]*Prescription // realm → cedula → list
	appointments  map[string]map[string][]*Appointment
}

func New() *Store {
	return &Store{
		prescriptions: map[string]map[string][]*Prescription{},
		appointments:  map[string]map[string][]*Appointment{},
	}
}

func (s *Store) Seed(realm string, prescriptions []Prescription, appointments []Appointment) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.prescriptions[realm]; !ok {
		s.prescriptions[realm] = map[string][]*Prescription{}
		s.appointments[realm] = map[string][]*Appointment{}
	}
	for i := range prescriptions {
		p := prescriptions[i]
		p.UpdatedAt = time.Now()
		s.prescriptions[realm][p.Cedula] = append(s.prescriptions[realm][p.Cedula], &p)
	}
	for i := range appointments {
		a := appointments[i]
		a.UpdatedAt = time.Now()
		s.appointments[realm][a.Cedula] = append(s.appointments[realm][a.Cedula], &a)
	}
}

func (s *Store) PrescriptionsByCedula(realm, cedula string) []*Prescription {
	s.mu.RLock()
	defer s.mu.RUnlock()
	src := s.prescriptions[realm][cedula]
	out := make([]*Prescription, len(src))
	copy(out, src)
	return out
}

func (s *Store) AppointmentsByCedula(realm, cedula string) []*Appointment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	src := s.appointments[realm][cedula]
	out := make([]*Appointment, len(src))
	copy(out, src)
	return out
}
