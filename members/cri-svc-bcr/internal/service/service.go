// Package service del cri-svc-bcr.
//
// BCR es un member PRIVADO (Tier 4) que consume el fabric costaricaasservice para
// onboarding de clientes (KYC + verificación de ingresos). No expone
// servicios — solo consume registro-civil y hacienda.
//
// Este es el caso de uso que diferencia el producto como SaaS B2G/B2B:
// los bancos privados pagan suscripción para consumir datos del Estado
// con consentimiento del ciudadano y todo queda auditado.
package service

import (
	"context"

	interop "github.com/devsebas/costaricaasservice/libs/cri-lib-interop-client/interop"
	screrrors "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/errors"
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

type personData struct {
	Cedula   string `json:"cedula"`
	FullName string `json:"fullName"`
	Address  string `json:"address"`
}

type taxData struct {
	Cedula      string  `json:"cedula"`
	Year        int     `json:"year"`
	GrossIncome float64 `json:"grossIncome"`
	Status      string  `json:"status"`
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

	person, personAuditID, err := interop.CallTyped[personData](ctx, s.interopClient, interop.CallRequest{
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

	tax, taxAuditID, err := interop.CallTyped[taxData](ctx, s.interopClient, interop.CallRequest{
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

	bracket := "low"
	approved := false
	reason := "ingreso bruto insuficiente"
	switch {
	case tax.GrossIncome >= 20_000_000:
		bracket = "high"
		approved = true
		reason = "ingreso comprobado, aprobado para producto premium"
	case tax.GrossIncome >= 10_000_000:
		bracket = "medium"
		approved = true
		reason = "ingreso comprobado, aprobado para producto estándar"
	}

	return &KYCResult{
		Cedula:        person.Cedula,
		FullName:      person.FullName,
		Address:       person.Address,
		IncomeBracket: bracket,
		GrossIncome:   tax.GrossIncome,
		Approved:      approved,
		Reason:        reason,
		Trace: []string{
			"registro-civil/persons.get/v1 audit_id=" + personAuditID,
			"hacienda/tax.status/v1 audit_id=" + taxAuditID,
		},
	}, nil
}
