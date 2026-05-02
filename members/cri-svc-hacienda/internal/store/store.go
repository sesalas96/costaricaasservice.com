// Package store del cri-svc-hacienda. In-memory para MVP demo.
package store

import (
	"strconv"
	"sync"
	"time"
)

// TaxRecord — info tributaria propia de Hacienda (lo que NO viene de
// Registro Civil). Por once-only: nombre/domicilio se piden vía interop a
// registro-civil, no se duplican aquí.
type TaxRecord struct {
	Cedula        string    `json:"cedula"`
	Year          int       `json:"year"`
	GrossIncome   float64   `json:"grossIncome"`
	WithheldTax   float64   `json:"withheldTax"`
	Deductions    float64   `json:"deductions"`
	HasDependents bool      `json:"hasDependents"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Store struct {
	mu      sync.RWMutex
	records map[string]map[string]*TaxRecord // realm → "cedula:year" → TaxRecord
}

func New() *Store {
	return &Store{records: make(map[string]map[string]*TaxRecord)}
}

func (s *Store) Seed(realm string, recs []TaxRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.records[realm]; !ok {
		s.records[realm] = make(map[string]*TaxRecord)
	}
	for i := range recs {
		r := recs[i]
		r.UpdatedAt = time.Now()
		s.records[realm][key(r.Cedula, r.Year)] = &r
	}
}

func (s *Store) Get(realm, cedula string, year int) *TaxRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if r, ok := s.records[realm]; ok {
		return r[key(cedula, year)]
	}
	return nil
}

func key(cedula string, year int) string { return cedula + ":" + strconv.Itoa(year) }
