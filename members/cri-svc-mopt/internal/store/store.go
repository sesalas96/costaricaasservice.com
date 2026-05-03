// Package store del cri-svc-mopt. In-memory para MVP.
//
// Once-only: el MOPT guarda SOLO datos viales (licencias, multas, puntos).
// Nombre/domicilio del conductor se piden a registro-civil vía interop.
package store

import (
	"sync"
	"time"
)

type DriverLicense struct {
	Number     string    `json:"number"`
	Cedula     string    `json:"cedula"`
	Categories []string  `json:"categories"` // B1, A2, C, etc.
	IssuedAt   string    `json:"issuedAt"`
	ExpiresAt  string    `json:"expiresAt"`
	Status     string    `json:"status"` // active | suspended | expired
	Points     int       `json:"points"` // licencia por puntos
	UpdatedAt  time.Time `json:"updatedAt"`
}

type TrafficFine struct {
	ID            string    `json:"id"`
	Cedula        string    `json:"cedula"`
	VehiclePlate  string    `json:"vehiclePlate"`
	Reason        string    `json:"reason"`
	AmountColons  int64     `json:"amountColons"`
	IssuedAt      string    `json:"issuedAt"`
	Status        string    `json:"status"` // pending | paid | appealed
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Store struct {
	mu       sync.RWMutex
	licenses map[string]map[string]*DriverLicense   // realm → cedula → license
	fines    map[string]map[string][]*TrafficFine   // realm → cedula → fines
}

func New() *Store {
	return &Store{
		licenses: map[string]map[string]*DriverLicense{},
		fines:    map[string]map[string][]*TrafficFine{},
	}
}

func (s *Store) Seed(realm string, licenses []DriverLicense, fines []TrafficFine) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.licenses[realm]; !ok {
		s.licenses[realm] = map[string]*DriverLicense{}
		s.fines[realm] = map[string][]*TrafficFine{}
	}
	for i := range licenses {
		l := licenses[i]
		l.UpdatedAt = time.Now()
		s.licenses[realm][l.Cedula] = &l
	}
	for i := range fines {
		f := fines[i]
		f.UpdatedAt = time.Now()
		s.fines[realm][f.Cedula] = append(s.fines[realm][f.Cedula], &f)
	}
}

func (s *Store) LicenseByCedula(realm, cedula string) *DriverLicense {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.licenses[realm][cedula]
}

func (s *Store) FinesByCedula(realm, cedula string) []*TrafficFine {
	s.mu.RLock()
	defer s.mu.RUnlock()
	src := s.fines[realm][cedula]
	out := make([]*TrafficFine, len(src))
	copy(out, src)
	return out
}
