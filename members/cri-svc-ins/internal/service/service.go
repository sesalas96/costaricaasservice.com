// Package service del cri-svc-ins.
//
// InsuranceProfile demuestra once-only: nombre/domicilio del titular se
// consultan a registro-civil vía interop SDK; el INS agrega las pólizas
// activas desde su propio store.
package service

import (
	"context"

	interop "github.com/devsebas/costaricaasservice/libs/cri-lib-interop-client/interop"
	screrrors "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/errors"

	"github.com/devsebas/costaricaasservice/members/cri-svc-ins/internal/store"
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

type InsuranceProfile struct {
	Holder         PersonSummary   `json:"holder"`
	ActivePolicies []*store.Policy `json:"activePolicies"`
	OnceOnlyTrace  string          `json:"_onceOnlyTrace"`
}

func (s *Service) InsuranceProfile(ctx context.Context, realm, cedula string) (*InsuranceProfile, error) {
	if cedula == "" {
		return nil, screrrors.New(screrrors.CodeBadRequest, "cedula required")
	}

	person, traceID, err := interop.CallTyped[PersonSummary](ctx, s.interopClient, interop.CallRequest{
		TargetMember: "registro-civil",
		Service:      "persons.get",
		Version:      "v1",
		Body:         map[string]any{"cedula": cedula},
		CitizenID:    cedula,
		Purpose:      "insurance_profile_lookup",
	})
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeUnavailable, "registro-civil lookup failed", err)
	}

	pols := s.store.PoliciesByCedula(realm, cedula)
	active := make([]*store.Policy, 0, len(pols))
	for _, p := range pols {
		if p.Status == "active" {
			active = append(active, p)
		}
	}

	return &InsuranceProfile{
		Holder:         person,
		ActivePolicies: active,
		OnceOnlyTrace:  "registro-civil/persons.get/v1 audit_id=" + traceID,
	}, nil
}

// PoliciesRawForInterop expone las pólizas activas para que otro member
// (ej. BCR para validar que un préstamo de auto tenga SOA al día) las
// consuma vía interop firmado.
func (s *Service) PoliciesRawForInterop(realm, cedula string) []*store.Policy {
	return s.store.PoliciesByCedula(realm, cedula)
}
