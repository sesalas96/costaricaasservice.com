// Package service del cri-svc-ccss.
//
// HealthProfile demuestra once-only: nombre/domicilio se consultan a
// registro-civil vía interop SDK; CCSS agrega recetas activas + próxima cita
// desde su propio store.
package service

import (
	"context"

	interop "github.com/devsebas/costaricaasservice/libs/cri-lib-interop-client/interop"
	screrrors "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/errors"

	"github.com/devsebas/costaricaasservice/members/cri-svc-ccss/internal/store"
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

type HealthProfile struct {
	Patient             PersonSummary         `json:"patient"`
	ActivePrescriptions []*store.Prescription `json:"activePrescriptions"`
	NextAppointment     *store.Appointment    `json:"nextAppointment"`
	OnceOnlyTrace       string                `json:"_onceOnlyTrace"`
}

func (s *Service) HealthProfile(ctx context.Context, realm, cedula string) (*HealthProfile, error) {
	if cedula == "" {
		return nil, screrrors.New(screrrors.CodeBadRequest, "cedula required")
	}

	person, traceID, err := interop.CallTyped[PersonSummary](ctx, s.interopClient, interop.CallRequest{
		TargetMember: "registro-civil",
		Service:      "persons.get",
		Version:      "v1",
		Body:         map[string]any{"cedula": cedula},
		CitizenID:    cedula,
		Purpose:      "health_profile_lookup",
	})
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeUnavailable, "registro-civil lookup failed", err)
	}

	rxs := s.store.PrescriptionsByCedula(realm, cedula)
	active := make([]*store.Prescription, 0, len(rxs))
	for _, p := range rxs {
		if p.Status == "active" {
			active = append(active, p)
		}
	}
	apts := s.store.AppointmentsByCedula(realm, cedula)
	var next *store.Appointment
	for _, a := range apts {
		if a.Status == "scheduled" {
			if next == nil || a.When < next.When {
				next = a
			}
		}
	}

	return &HealthProfile{
		Patient:             person,
		ActivePrescriptions: active,
		NextAppointment:     next,
		OnceOnlyTrace:       "registro-civil/persons.get/v1 audit_id=" + traceID,
	}, nil
}

// PrescriptionsRawForInterop retorna las recetas activas para que otro member
// (p.ej. una farmacia privada) las consuma vía interop firmado.
func (s *Service) PrescriptionsRawForInterop(realm, cedula string) []*store.Prescription {
	return s.store.PrescriptionsByCedula(realm, cedula)
}
