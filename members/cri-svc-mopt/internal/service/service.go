// Package service del cri-svc-mopt.
//
// DriverProfile demuestra once-only: nombre/domicilio del conductor se
// consultan a registro-civil vía interop SDK; el MOPT agrega licencia
// y multas pendientes desde su propio store.
package service

import (
	"context"

	interop "github.com/devsebas/costaricaasservice/libs/cri-lib-interop-client/interop"
	screrrors "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/errors"

	"github.com/devsebas/costaricaasservice/members/cri-svc-mopt/internal/store"
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

type DriverProfile struct {
	Driver        PersonSummary        `json:"driver"`
	License       *store.DriverLicense `json:"license,omitempty"`
	PendingFines  []*store.TrafficFine `json:"pendingFines"`
	OnceOnlyTrace string               `json:"_onceOnlyTrace"`
}

func (s *Service) DriverProfile(ctx context.Context, realm, cedula string) (*DriverProfile, error) {
	if cedula == "" {
		return nil, screrrors.New(screrrors.CodeBadRequest, "cedula required")
	}

	person, traceID, err := interop.CallTyped[PersonSummary](ctx, s.interopClient, interop.CallRequest{
		TargetMember: "registro-civil",
		Service:      "persons.get",
		Version:      "v1",
		Body:         map[string]any{"cedula": cedula},
		CitizenID:    cedula,
		Purpose:      "driver_profile_lookup",
	})
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeUnavailable, "registro-civil lookup failed", err)
	}

	lic := s.store.LicenseByCedula(realm, cedula)
	all := s.store.FinesByCedula(realm, cedula)
	pending := make([]*store.TrafficFine, 0, len(all))
	for _, f := range all {
		if f.Status == "pending" {
			pending = append(pending, f)
		}
	}

	return &DriverProfile{
		Driver:        person,
		License:       lic,
		PendingFines:  pending,
		OnceOnlyTrace: "registro-civil/persons.get/v1 audit_id=" + traceID,
	}, nil
}

// LicenseStatusForInterop devuelve la licencia para que otros members
// (ej. INS confirmando que un asegurado tiene licencia vigente al emitir SOA)
// la consuman vía interop firmado.
func (s *Service) LicenseStatusForInterop(realm, cedula string) *store.DriverLicense {
	return s.store.LicenseByCedula(realm, cedula)
}
