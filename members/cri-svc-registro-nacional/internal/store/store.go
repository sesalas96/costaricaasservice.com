// Package store del cri-svc-registro-nacional. In-memory para MVP.
//
// Once-only: el Registro Nacional guarda SOLO datos de propiedades
// (inmuebles, vehículos, marcas, sociedades). Nombre/domicilio del
// propietario se piden a registro-civil vía interop.
package store

import (
	"sync"
	"time"
)

type RealEstate struct {
	Folio          string    `json:"folio"`            // folio real
	OwnerCedula    string    `json:"ownerCedula"`
	Province       string    `json:"province"`
	Canton         string    `json:"canton"`
	District       string    `json:"district"`
	AreaSquareMts  float64   `json:"areaSquareMts"`
	FiscalValue    int64     `json:"fiscalValueColons"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type Vehicle struct {
	Plate         string    `json:"plate"`
	OwnerCedula   string    `json:"ownerCedula"`
	Make          string    `json:"make"`
	Model         string    `json:"model"`
	Year          int       `json:"year"`
	FiscalValue   int64     `json:"fiscalValueColons"`
	Status        string    `json:"status"` // active | scrapped | exported
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Store struct {
	mu         sync.RWMutex
	realEstate map[string]map[string][]*RealEstate // realm → cedula → list
	vehicles   map[string]map[string][]*Vehicle    // realm → cedula → list
}

func New() *Store {
	return &Store{
		realEstate: map[string]map[string][]*RealEstate{},
		vehicles:   map[string]map[string][]*Vehicle{},
	}
}

func (s *Store) Seed(realm string, properties []RealEstate, vehicles []Vehicle) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.realEstate[realm]; !ok {
		s.realEstate[realm] = map[string][]*RealEstate{}
		s.vehicles[realm] = map[string][]*Vehicle{}
	}
	for i := range properties {
		p := properties[i]
		p.UpdatedAt = time.Now()
		s.realEstate[realm][p.OwnerCedula] = append(s.realEstate[realm][p.OwnerCedula], &p)
	}
	for i := range vehicles {
		v := vehicles[i]
		v.UpdatedAt = time.Now()
		s.vehicles[realm][v.OwnerCedula] = append(s.vehicles[realm][v.OwnerCedula], &v)
	}
}

func (s *Store) RealEstateByCedula(realm, cedula string) []*RealEstate {
	s.mu.RLock()
	defer s.mu.RUnlock()
	src := s.realEstate[realm][cedula]
	out := make([]*RealEstate, len(src))
	copy(out, src)
	return out
}

func (s *Store) VehiclesByCedula(realm, cedula string) []*Vehicle {
	s.mu.RLock()
	defer s.mu.RUnlock()
	src := s.vehicles[realm][cedula]
	out := make([]*Vehicle, len(src))
	copy(out, src)
	return out
}
