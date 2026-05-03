// Package service del cri-svc-mopt.
//
// DriverProfile demuestra once-only: nombre/domicilio del conductor se
// consultan a registro-civil vía interop SDK; el MOPT agrega licencia
// y multas pendientes desde su propio store.
package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

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

	person, traceID, err := s.fetchPerson(ctx, realm, cedula)
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
		Driver:        *person,
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

func (s *Service) fetchPerson(ctx context.Context, realm, cedula string) (*PersonSummary, string, error) {
	resp, err := s.interopClient.Call(ctx, interop.CallRequest{
		TargetMember: "registro-civil",
		Service:      "persons.get",
		Version:      "v1",
		Body:         map[string]any{"cedula": cedula},
		CitizenID:    cedula,
		Purpose:      "driver_profile_lookup",
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
