// Package service del cri-svc-registro-nacional.
//
// PropertyProfile demuestra once-only: nombre/domicilio del propietario se
// consultan a registro-civil vía interop SDK; el Registro Nacional agrega
// inmuebles y vehículos a nombre desde su propio store.
package service

import (
	"context"

	interop "github.com/devsebas/costaricaasservice/libs/cri-lib-interop-client/interop"
	screrrors "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/errors"

	"github.com/devsebas/costaricaasservice/members/cri-svc-registro-nacional/internal/store"
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

type PropertyProfile struct {
	Owner         PersonSummary       `json:"owner"`
	RealEstate    []*store.RealEstate `json:"realEstate"`
	Vehicles      []*store.Vehicle    `json:"vehicles"`
	OnceOnlyTrace string              `json:"_onceOnlyTrace"`
}

func (s *Service) PropertyProfile(ctx context.Context, realm, cedula string) (*PropertyProfile, error) {
	if cedula == "" {
		return nil, screrrors.New(screrrors.CodeBadRequest, "cedula required")
	}

	person, traceID, err := interop.CallTyped[PersonSummary](ctx, s.interopClient, interop.CallRequest{
		TargetMember: "registro-civil",
		Service:      "persons.get",
		Version:      "v1",
		Body:         map[string]any{"cedula": cedula},
		CitizenID:    cedula,
		Purpose:      "property_profile_lookup",
	})
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeUnavailable, "registro-civil lookup failed", err)
	}

	props := s.store.RealEstateByCedula(realm, cedula)
	vehs := s.store.VehiclesByCedula(realm, cedula)

	return &PropertyProfile{
		Owner:         person,
		RealEstate:    props,
		Vehicles:      vehs,
		OnceOnlyTrace: "registro-civil/persons.get/v1 audit_id=" + traceID,
	}, nil
}

// AssetsForInterop expone propiedades + vehículos para que otros members
// (ej. BCR como garantía hipotecaria, Hacienda para impuesto solidario)
// los consuman vía interop firmado.
func (s *Service) AssetsForInterop(realm, cedula string) map[string]any {
	return map[string]any{
		"realEstate": s.store.RealEstateByCedula(realm, cedula),
		"vehicles":   s.store.VehiclesByCedula(realm, cedula),
	}
}
