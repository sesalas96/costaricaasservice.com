// Package service del cri-svc-mtss.
//
// LaborProfile demuestra once-only: nombre/domicilio del trabajador se
// consultan a registro-civil vía interop SDK; el MTSS agrega los empleos
// activos desde su propio store.
package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	interop "github.com/devsebas/costaricaasservice/libs/cri-lib-interop-client/interop"
	screrrors "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/errors"

	"github.com/devsebas/costaricaasservice/members/cri-svc-mtss/internal/store"
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

type LaborProfile struct {
	Worker            PersonSummary       `json:"worker"`
	ActiveEmployments []*store.Employment `json:"activeEmployments"`
	OnceOnlyTrace     string              `json:"_onceOnlyTrace"`
}

func (s *Service) LaborProfile(ctx context.Context, realm, cedula string) (*LaborProfile, error) {
	if cedula == "" {
		return nil, screrrors.New(screrrors.CodeBadRequest, "cedula required")
	}

	person, traceID, err := s.fetchPerson(ctx, realm, cedula)
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeUnavailable, "registro-civil lookup failed", err)
	}

	emps := s.store.EmploymentsByCedula(realm, cedula)
	active := make([]*store.Employment, 0, len(emps))
	for _, e := range emps {
		if e.Status == "active" {
			active = append(active, e)
		}
	}

	return &LaborProfile{
		Worker:            *person,
		ActiveEmployments: active,
		OnceOnlyTrace:     "registro-civil/persons.get/v1 audit_id=" + traceID,
	}, nil
}

// EmploymentStatusForInterop expone empleos activos para que otros members
// (ej. BCR para verificar capacidad de pago) los consuman vía interop firmado.
func (s *Service) EmploymentStatusForInterop(realm, cedula string) []*store.Employment {
	all := s.store.EmploymentsByCedula(realm, cedula)
	out := make([]*store.Employment, 0, len(all))
	for _, e := range all {
		if e.Status == "active" {
			out = append(out, e)
		}
	}
	return out
}

func (s *Service) fetchPerson(ctx context.Context, realm, cedula string) (*PersonSummary, string, error) {
	resp, err := s.interopClient.Call(ctx, interop.CallRequest{
		TargetMember: "registro-civil",
		Service:      "persons.get",
		Version:      "v1",
		Body:         map[string]any{"cedula": cedula},
		CitizenID:    cedula,
		Purpose:      "labor_profile_lookup",
	})
	if err != nil {
		return nil, "", err
	}
	if resp.Status >= 400 {
		return nil, "", errors.New("registro-civil returned " + strconv.Itoa(resp.Status))
	}
	var ssEnv struct {
		Data struct {
			AuditID string          `json:"audit_id"`
			Body    json.RawMessage `json:"body"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &ssEnv); err != nil {
		return nil, "", err
	}
	var peerEnv struct {
		Data struct {
			Cedula   string `json:"cedula"`
			FullName string `json:"fullName"`
			Address  string `json:"address"`
		} `json:"data"`
	}
	if err := json.Unmarshal(ssEnv.Data.Body, &peerEnv); err != nil {
		return nil, "", err
	}
	return &PersonSummary{
		Cedula:   peerEnv.Data.Cedula,
		FullName: peerEnv.Data.FullName,
		Address:  peerEnv.Data.Address,
	}, ssEnv.Data.AuditID, nil
}
