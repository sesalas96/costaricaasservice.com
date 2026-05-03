// Package service del cri-svc-minsa.
//
// HealthPermitProfile demuestra once-only: nombre/domicilio del titular se
// consultan a registro-civil vía interop SDK; el MINSA agrega los permisos
// sanitarios activos desde su propio store.
package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	interop "github.com/devsebas/costaricaasservice/libs/cri-lib-interop-client/interop"
	screrrors "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/errors"

	"github.com/devsebas/costaricaasservice/members/cri-svc-minsa/internal/store"
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

type HealthPermitProfile struct {
	Holder        PersonSummary         `json:"holder"`
	ActivePermits []*store.HealthPermit `json:"activePermits"`
	OnceOnlyTrace string                `json:"_onceOnlyTrace"`
}

func (s *Service) HealthPermitProfile(ctx context.Context, realm, cedula string) (*HealthPermitProfile, error) {
	if cedula == "" {
		return nil, screrrors.New(screrrors.CodeBadRequest, "cedula required")
	}

	person, traceID, err := s.fetchPerson(ctx, realm, cedula)
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeUnavailable, "registro-civil lookup failed", err)
	}

	all := s.store.PermitsByCedula(realm, cedula)
	active := make([]*store.HealthPermit, 0, len(all))
	for _, p := range all {
		if p.Status == "active" {
			active = append(active, p)
		}
	}

	return &HealthPermitProfile{
		Holder:        *person,
		ActivePermits: active,
		OnceOnlyTrace: "registro-civil/persons.get/v1 audit_id=" + traceID,
	}, nil
}

// PermitsRawForInterop expone permisos sanitarios para que otros members
// (ej. una municipalidad antes de otorgar patente comercial) los consuman.
func (s *Service) PermitsRawForInterop(realm, cedula string) []*store.HealthPermit {
	return s.store.PermitsByCedula(realm, cedula)
}

func (s *Service) fetchPerson(ctx context.Context, realm, cedula string) (*PersonSummary, string, error) {
	resp, err := s.interopClient.Call(ctx, interop.CallRequest{
		TargetMember: "registro-civil",
		Service:      "persons.get",
		Version:      "v1",
		Body:         map[string]any{"cedula": cedula},
		CitizenID:    cedula,
		Purpose:      "health_permit_lookup",
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
