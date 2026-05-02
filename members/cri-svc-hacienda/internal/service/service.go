// Package service del cri-svc-hacienda.
//
// El método clave que demuestra ONCE-ONLY es PrefilledReturn: para construir
// la declaración de renta pre-llenada, Hacienda consume datos básicos del
// ciudadano desde Registro Civil vía interop, en lugar de pedir al ciudadano
// volver a digitarlos.
package service

import (
	"context"
	"encoding/json"
	"errors"

	screrrors "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/errors"
	interop "github.com/devsebas/costaricaasservice/libs/cri-lib-interop-client/interop"

	"github.com/devsebas/costaricaasservice/members/cri-svc-hacienda/internal/store"
)

type Service struct {
	store  *store.Store
	interopClient *interop.Client // local SS endpoint
}

// New construye el service. interopClient apunta al SS local de hacienda.
func New(s *store.Store, ic *interop.Client) *Service {
	return &Service{store: s, interopClient: ic}
}

// TaxStatus es la respuesta de un status simple.
type TaxStatus struct {
	Cedula      string  `json:"cedula"`
	Year        int     `json:"year"`
	GrossIncome float64 `json:"grossIncome"`
	WithheldTax float64 `json:"withheldTax"`
	Status      string  `json:"status"` // "filed", "pending", "no-record"
}

// GetTaxStatus retorna el status tributario sin consultar Registro Civil
// (datos puramente de Hacienda).
func (s *Service) GetTaxStatus(realm, cedula string, year int) (*TaxStatus, error) {
	if cedula == "" {
		return nil, screrrors.New(screrrors.CodeBadRequest, "cedula required")
	}
	rec := s.store.Get(realm, cedula, year)
	if rec == nil {
		return &TaxStatus{Cedula: cedula, Year: year, Status: "no-record"}, nil
	}
	return &TaxStatus{
		Cedula:      rec.Cedula,
		Year:        rec.Year,
		GrossIncome: rec.GrossIncome,
		WithheldTax: rec.WithheldTax,
		Status:      "pending",
	}, nil
}

// PersonSummary son los datos básicos del ciudadano (vienen de Registro Civil
// vía interop). Hacienda NO los persiste — los consulta cada vez.
type PersonSummary struct {
	Cedula   string `json:"cedula"`
	FullName string `json:"fullName"`
	Address  string `json:"address"`
	Email    string `json:"email"`
}

// PrefilledReturn es la declaración de renta pre-llenada que ve el ciudadano.
type PrefilledReturn struct {
	Person          PersonSummary `json:"person"`
	Year            int           `json:"year"`
	GrossIncome     float64       `json:"grossIncome"`
	WithheldTax     float64       `json:"withheldTax"`
	Deductions      float64       `json:"deductions"`
	EstimatedDue    float64       `json:"estimatedDue"`
	HasDependents   bool          `json:"hasDependents"`
	OnceOnlyTrace   string        `json:"_onceOnlyTrace"` // demuestra a quién se consultó vía interop
}

// PrefilledReturn arma la declaración pre-llenada. Para construir la sección
// "person" hace una llamada inter-member a registro-civil vía el SS local.
func (s *Service) PrefilledReturn(ctx context.Context, realm, cedula string, year int) (*PrefilledReturn, error) {
	if cedula == "" {
		return nil, screrrors.New(screrrors.CodeBadRequest, "cedula required")
	}
	rec := s.store.Get(realm, cedula, year)
	if rec == nil {
		return nil, screrrors.New(screrrors.CodeNotFound, "no tax record for cedula/year")
	}

	person, traceID, err := s.fetchPerson(ctx, realm, cedula)
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeUnavailable, "registro-civil lookup failed", err)
	}

	estimated := computeEstimatedDue(rec.GrossIncome, rec.WithheldTax, rec.Deductions, rec.HasDependents)

	return &PrefilledReturn{
		Person:        *person,
		Year:          rec.Year,
		GrossIncome:   rec.GrossIncome,
		WithheldTax:   rec.WithheldTax,
		Deductions:    rec.Deductions,
		EstimatedDue:  estimated,
		HasDependents: rec.HasDependents,
		OnceOnlyTrace: "registro-civil/persons.get/v1 audit_id=" + traceID,
	}, nil
}

// fetchPerson llama vía interop SDK al SS local; SS firma y enruta a registro-civil.
func (s *Service) fetchPerson(ctx context.Context, realm, cedula string) (*PersonSummary, string, error) {
	resp, err := s.interopClient.Call(ctx, interop.CallRequest{
		TargetMember: "registro-civil",
		Service:      "persons.get",
		Version:      "v1",
		Body:         map[string]any{"cedula": cedula},
		CitizenID:    cedula,
		Purpose:      "prefill_tax_return",
	})
	if err != nil {
		return nil, "", err
	}
	if resp.Status >= 400 {
		return nil, "", errors.New("registro-civil returned " + http2string(resp.Status))
	}

	// resp.Body viene en doble envelope:
	//   nivel 1 — SS local devuelve {data: {audit_id, body, ...}, meta:{}}
	//   nivel 2 — dentro de .data.body está el envelope del peer dispatcher:
	//             {data: <Person>, meta:{}}
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
			Email    string `json:"email"`
		} `json:"data"`
	}
	if err := json.Unmarshal(ssEnv.Data.Body, &peerEnv); err != nil {
		return nil, "", err
	}
	return &PersonSummary{
		Cedula:   peerEnv.Data.Cedula,
		FullName: peerEnv.Data.FullName,
		Address:  peerEnv.Data.Address,
		Email:    peerEnv.Data.Email,
	}, ssEnv.Data.AuditID, nil
}

// http2string evita importar strconv solo para esto.
func http2string(s int) string {
	d := []byte{0, 0, 0}
	i := 2
	if s == 0 {
		return "0"
	}
	for s > 0 && i >= 0 {
		d[i] = byte('0' + s%10)
		s /= 10
		i--
	}
	return string(d[i+1:])
}

// computeEstimatedDue calcula una estimación tonta del impuesto restante.
// Tax bracket simple para el demo: 13% bruto - retenciones - deducciones.
func computeEstimatedDue(gross, withheld, deductions float64, hasDependents bool) float64 {
	rate := 0.13
	if hasDependents {
		rate = 0.10
	}
	due := gross*rate - withheld - deductions
	if due < 0 {
		return 0
	}
	return due
}
