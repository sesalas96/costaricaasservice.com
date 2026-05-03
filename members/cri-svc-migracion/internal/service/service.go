// Package service del cri-svc-migracion.
//
// ImmigrationProfile demuestra once-only: nombre/domicilio del titular se
// consultan a registro-civil vía interop SDK; Migración agrega el status
// migratorio y movimientos recientes desde su propio store.
package service

import (
	"context"

	interop "github.com/devsebas/costaricaasservice/libs/cri-lib-interop-client/interop"
	screrrors "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/errors"

	"github.com/devsebas/costaricaasservice/members/cri-svc-migracion/internal/store"
)

type Service struct {
	store         *store.Store
	interopClient *interop.Client
}

func New(s *store.Store, ic *interop.Client) *Service {
	return &Service{store: s, interopClient: ic}
}

type PersonSummary struct {
	Cedula   string `json:"cedula"`
	FullName string `json:"fullName"`
	Address  string `json:"address"`
}

type ImmigrationProfile struct {
	Person          PersonSummary            `json:"person"`
	Status          *store.ImmigrationStatus `json:"status,omitempty"`
	RecentMovements []*store.Movement        `json:"recentMovements"`
	OnceOnlyTrace   string                   `json:"_onceOnlyTrace"`
}

func (s *Service) ImmigrationProfile(ctx context.Context, realm, cedula string) (*ImmigrationProfile, error) {
	if cedula == "" {
		return nil, screrrors.New(screrrors.CodeBadRequest, "cedula required")
	}

	person, traceID, err := interop.CallTyped[PersonSummary](ctx, s.interopClient, interop.CallRequest{
		TargetMember: "registro-civil",
		Service:      "persons.get",
		Version:      "v1",
		Body:         map[string]any{"cedula": cedula},
		CitizenID:    cedula,
		Purpose:      "immigration_profile_lookup",
	})
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeUnavailable, "registro-civil lookup failed", err)
	}

	st := s.store.StatusByCedula(realm, cedula)
	movs := s.store.MovementsByCedula(realm, cedula)

	return &ImmigrationProfile{
		Person:          person,
		Status:          st,
		RecentMovements: movs,
		OnceOnlyTrace:   "registro-civil/persons.get/v1 audit_id=" + traceID,
	}, nil
}

// StatusForInterop expone el status migratorio para que otros members
// (ej. BCR antes de un préstamo, MEP al matricular) lo consuman.
func (s *Service) StatusForInterop(realm, cedula string) *store.ImmigrationStatus {
	return s.store.StatusByCedula(realm, cedula)
}
