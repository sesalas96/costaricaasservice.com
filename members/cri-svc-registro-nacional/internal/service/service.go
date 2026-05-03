// Package service del cri-svc-registro-nacional.
//
// PropertyProfile demuestra once-only: nombre/domicilio del propietario se
// consultan a registro-civil vía interop SDK; el Registro Nacional agrega
// inmuebles y vehículos a nombre desde su propio store.
package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

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

	person, traceID, err := s.fetchPerson(ctx, realm, cedula)
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeUnavailable, "registro-civil lookup failed", err)
	}

	props := s.store.RealEstateByCedula(realm, cedula)
	vehs := s.store.VehiclesByCedula(realm, cedula)

	return &PropertyProfile{
		Owner:         *person,
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

func (s *Service) fetchPerson(ctx context.Context, realm, cedula string) (*PersonSummary, string, error) {
	resp, err := s.interopClient.Call(ctx, interop.CallRequest{
		TargetMember: "registro-civil",
		Service:      "persons.get",
		Version:      "v1",
		Body:         map[string]any{"cedula": cedula},
		CitizenID:    cedula,
		Purpose:      "property_profile_lookup",
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
