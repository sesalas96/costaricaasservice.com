// Package service del cri-svc-bcr.
//
// BCR es un member PRIVADO (Tier 4) que consume el fabric saascr para
// onboarding de clientes (KYC + verificación de ingresos). No expone
// servicios — solo consume registro-civil y hacienda.
//
// Este es el caso de uso que diferencia el producto como SaaS B2G/B2B:
// los bancos privados pagan suscripción para consumir datos del Estado
// con consentimiento del ciudadano y todo queda auditado.
package service

import (
	"context"
	"encoding/json"
	"strconv"

	interop "github.com/devsebas/saascr/libs/cri-lib-interop-client/interop"
	screrrors "github.com/devsebas/saascr/libs/cri-lib-shared/errors"
)

type Service struct {
	interopClient *interop.Client
}

func New(ic *interop.Client) *Service { return &Service{interopClient: ic} }

// KYCResult es el resumen del onboarding bancario.
type KYCResult struct {
	Cedula        string   `json:"cedula"`
	FullName      string   `json:"fullName"`
	Address       string   `json:"address"`
	IncomeBracket string   `json:"incomeBracket"` // "low" | "medium" | "high"
	GrossIncome   float64  `json:"grossIncome"`
	Approved      bool     `json:"approved"`
	Reason        string   `json:"reason"`
	Trace         []string `json:"trace"`
}

// KYCCheck dispara el flujo de onboarding: consulta registro-civil para
// nombre/domicilio + hacienda para ingreso bruto. Compone un resultado.
//
// Cada llamada genera un audit entry firmado; el ciudadano lo verá en su
// bitácora de MiCR como "bcr consultó tu cédula con propósito kyc_onboarding".
func (s *Service) KYCCheck(ctx context.Context, realm, cedula string) (*KYCResult, error) {
	if cedula == "" {
		return nil, screrrors.New(screrrors.CodeBadRequest, "cedula required")
	}

	// 1. Datos personales (Registro Civil)
	personResp, err := s.interopClient.Call(ctx, interop.CallRequest{
		TargetMember: "registro-civil",
		Service:      "persons.get",
		Version:      "v1",
		Body:         map[string]any{"cedula": cedula},
		CitizenID:    cedula,
		Purpose:      "kyc_onboarding",
	})
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeUnavailable, "registro-civil lookup", err)
	}
	if personResp.Status >= 400 {
		return nil, screrrors.New(screrrors.CodeNotFound, "person not in registro-civil: "+strconv.Itoa(personResp.Status))
	}

	type personData struct {
		Cedula   string `json:"cedula"`
		FullName string `json:"fullName"`
		Address  string `json:"address"`
	}
	personFields, personAuditID := unwrap[personData](personResp.Body)

	// 2. Ingresos (Hacienda)
	taxResp, err := s.interopClient.Call(ctx, interop.CallRequest{
		TargetMember: "hacienda",
		Service:      "tax.status",
		Version:      "v1",
		Body:         map[string]any{"cedula": cedula, "year": 2025},
		CitizenID:    cedula,
		Purpose:      "kyc_income_verification",
	})
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeUnavailable, "hacienda lookup", err)
	}

	type taxData struct {
		Cedula      string  `json:"cedula"`
		Year        int     `json:"year"`
		GrossIncome float64 `json:"grossIncome"`
		Status      string  `json:"status"`
	}
	taxFields, taxAuditID := unwrap[taxData](taxResp.Body)

	bracket := "low"
	approved := false
	reason := "ingreso bruto insuficiente"
	switch {
	case taxFields.GrossIncome >= 20_000_000:
		bracket = "high"
		approved = true
		reason = "ingreso comprobado, aprobado para producto premium"
	case taxFields.GrossIncome >= 10_000_000:
		bracket = "medium"
		approved = true
		reason = "ingreso comprobado, aprobado para producto estándar"
	}

	return &KYCResult{
		Cedula:        personFields.Cedula,
		FullName:      personFields.FullName,
		Address:       personFields.Address,
		IncomeBracket: bracket,
		GrossIncome:   taxFields.GrossIncome,
		Approved:      approved,
		Reason:        reason,
		Trace: []string{
			"registro-civil/persons.get/v1 audit_id=" + personAuditID,
			"hacienda/tax.status/v1 audit_id=" + taxAuditID,
		},
	}, nil
}

// unwrap desempaca el doble envelope SS → peer-dispatcher → data<T>.
// Si falla algún unmarshal devuelve zero values.
func unwrap[T any](raw []byte) (T, string) {
	var zero T
	var ssEnv struct {
		Data struct {
			AuditID string          `json:"audit_id"`
			Body    json.RawMessage `json:"body"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &ssEnv); err != nil {
		return zero, ""
	}
	var peerEnv struct {
		Data T `json:"data"`
	}
	if err := json.Unmarshal(ssEnv.Data.Body, &peerEnv); err != nil {
		return zero, ssEnv.Data.AuditID
	}
	return peerEnv.Data, ssEnv.Data.AuditID
}
