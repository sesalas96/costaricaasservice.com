// Package service contiene la lógica de orquestación del cri-bff-citizen.
//
// El BFF compone respuestas para MiCR llamando varios upstreams en paralelo:
//   - hacienda (declaración pre-llenada — que a su vez ya consulta registro-civil
//     vía interop)
//   - interop-audit (bitácora ciudadana de accesos)
package service

import (
	"context"
	"encoding/json"

	screrrors "github.com/devsebas/saascr/libs/cri-lib-shared/errors"

	"github.com/devsebas/saascr/gateway/cri-bff-citizen/internal/clients"
)

type Service struct {
	hacienda *clients.Hacienda
	ccss     *clients.CCSS
	audit    *clients.Audit
}

func New(hacienda *clients.Hacienda, ccss *clients.CCSS, audit *clients.Audit) *Service {
	return &Service{hacienda: hacienda, ccss: ccss, audit: audit}
}

// Dashboard agrupa toda la info que MiCR muestra al ciudadano en una sola
// llamada. Hace fan-out paralelo a los upstreams.
type Dashboard struct {
	Cedula        string          `json:"cedula"`
	Realm         string          `json:"realm"`
	Year          int             `json:"year"`
	TaxPrefilled  json.RawMessage `json:"taxPrefilled"`
	HealthProfile json.RawMessage `json:"healthProfile"`
	AccessLog     json.RawMessage `json:"accessLog"`
	UpstreamErrs  []string        `json:"upstreamErrors,omitempty"`
}

// GetDashboard compone el dashboard. Orden secuencial intencional:
// primero hacienda (que dispara el audit entry vía SS), luego audit (que ya
// debería incluirla). Llamarlas en paralelo race con el audit emit.
func (s *Service) GetDashboard(ctx context.Context, realm, cedula string, year int) (*Dashboard, error) {
	if cedula == "" {
		return nil, screrrors.New(screrrors.CodeBadRequest, "cedula required")
	}
	out := &Dashboard{Cedula: cedula, Realm: realm, Year: year}
	addErr := func(name string, err error) {
		if err != nil {
			out.UpstreamErrs = append(out.UpstreamErrs, name+": "+err.Error())
		}
	}

	if body, err := s.hacienda.PrefilledReturn(ctx, realm, cedula, year); err != nil {
		addErr("hacienda.prefilled", err)
	} else {
		out.TaxPrefilled = body
	}

	if body, err := s.ccss.HealthProfile(ctx, realm, cedula); err != nil {
		addErr("ccss.health-profile", err)
	} else {
		out.HealthProfile = body
	}

	// audit al final: ya están registradas las consultas hacienda y ccss.
	if body, err := s.audit.AccessLogByCitizen(ctx, realm, cedula); err != nil {
		addErr("audit.access-log", err)
	} else {
		out.AccessLog = body
	}

	return out, nil
}

// GetHealthProfile retorna solamente el perfil de salud.
func (s *Service) GetHealthProfile(ctx context.Context, realm, cedula string) (json.RawMessage, error) {
	if cedula == "" {
		return nil, screrrors.New(screrrors.CodeBadRequest, "cedula required")
	}
	body, err := s.ccss.HealthProfile(ctx, realm, cedula)
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeUnavailable, "ccss upstream", err)
	}
	return body, nil
}

// GetAccessLog retorna solamente la bitácora (consulta más liviana).
func (s *Service) GetAccessLog(ctx context.Context, realm, cedula string) (json.RawMessage, error) {
	if cedula == "" {
		return nil, screrrors.New(screrrors.CodeBadRequest, "cedula required")
	}
	body, err := s.audit.AccessLogByCitizen(ctx, realm, cedula)
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeUnavailable, "audit upstream", err)
	}
	return body, nil
}

// GetTaxPrefilled retorna solamente la declaración pre-llenada.
func (s *Service) GetTaxPrefilled(ctx context.Context, realm, cedula string, year int) (json.RawMessage, error) {
	if cedula == "" {
		return nil, screrrors.New(screrrors.CodeBadRequest, "cedula required")
	}
	body, err := s.hacienda.PrefilledReturn(ctx, realm, cedula, year)
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeUnavailable, "hacienda upstream", err)
	}
	return body, nil
}
