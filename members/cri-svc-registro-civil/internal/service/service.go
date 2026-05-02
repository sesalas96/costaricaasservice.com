// Package service del cri-svc-registro-civil.
package service

import (
	screrrors "github.com/devsebas/saascr/libs/cri-lib-shared/errors"

	"github.com/devsebas/saascr/members/cri-svc-registro-civil/internal/store"
)

type Service struct {
	store *store.Store
}

func New(s *store.Store) *Service { return &Service{store: s} }

// GetPerson devuelve la persona o NotFound.
func (s *Service) GetPerson(realm, cedula string) (*store.Person, error) {
	if cedula == "" {
		return nil, screrrors.New(screrrors.CodeBadRequest, "cedula required")
	}
	p := s.store.GetByCedula(realm, cedula)
	if p == nil {
		return nil, screrrors.New(screrrors.CodeNotFound, "person not found")
	}
	return p, nil
}

// VitalEvents devuelve los eventos vitales de una persona.
func (s *Service) VitalEvents(realm, cedula string) ([]*store.VitalEvent, error) {
	if cedula == "" {
		return nil, screrrors.New(screrrors.CodeBadRequest, "cedula required")
	}
	if s.store.GetByCedula(realm, cedula) == nil {
		return nil, screrrors.New(screrrors.CodeNotFound, "person not found")
	}
	return s.store.EventsByCedula(realm, cedula), nil
}
