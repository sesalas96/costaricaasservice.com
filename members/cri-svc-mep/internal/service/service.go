// Package service del cri-svc-mep.
//
// AcademicProfile demuestra once-only: nombre/domicilio del estudiante se
// consultan a registro-civil vía interop SDK; el MEP agrega matrícula
// activa y certificados desde su propio store.
package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	interop "github.com/devsebas/costaricaasservice/libs/cri-lib-interop-client/interop"
	screrrors "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/errors"

	"github.com/devsebas/costaricaasservice/members/cri-svc-mep/internal/store"
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

type AcademicProfile struct {
	Student          PersonSummary        `json:"student"`
	ActiveEnrollment *store.Enrollment    `json:"activeEnrollment,omitempty"`
	Certificates     []*store.Certificate `json:"certificates"`
	OnceOnlyTrace    string               `json:"_onceOnlyTrace"`
}

func (s *Service) AcademicProfile(ctx context.Context, realm, cedula string) (*AcademicProfile, error) {
	if cedula == "" {
		return nil, screrrors.New(screrrors.CodeBadRequest, "cedula required")
	}

	person, traceID, err := s.fetchPerson(ctx, realm, cedula)
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeUnavailable, "registro-civil lookup failed", err)
	}

	enrolls := s.store.EnrollmentsByCedula(realm, cedula)
	var active *store.Enrollment
	for _, e := range enrolls {
		if e.Status == "active" {
			active = e
			break
		}
	}
	certs := s.store.CertificatesByCedula(realm, cedula)

	return &AcademicProfile{
		Student:          *person,
		ActiveEnrollment: active,
		Certificates:     certs,
		OnceOnlyTrace:    "registro-civil/persons.get/v1 audit_id=" + traceID,
	}, nil
}

// CertificatesForInterop expone los certificados académicos para que otros
// members (ej. una universidad privada validando bachillerato) los consuman.
func (s *Service) CertificatesForInterop(realm, cedula string) []*store.Certificate {
	return s.store.CertificatesByCedula(realm, cedula)
}

func (s *Service) fetchPerson(ctx context.Context, realm, cedula string) (*PersonSummary, string, error) {
	resp, err := s.interopClient.Call(ctx, interop.CallRequest{
		TargetMember: "registro-civil",
		Service:      "persons.get",
		Version:      "v1",
		Body:         map[string]any{"cedula": cedula},
		CitizenID:    cedula,
		Purpose:      "academic_profile_lookup",
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
