// Package service del cri-svc-interop-audit.
package service

import (
	screrrors "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/errors"

	"github.com/devsebas/costaricaasservice/interop/cri-svc-interop-audit/internal/auditstore"
)

type Service struct {
	store *auditstore.Store
}

func New(s *auditstore.Store) *Service { return &Service{store: s} }

// AppendInput re-exporta para que los handlers no dependan del store.
type AppendInput = auditstore.AppendInput

// Append agrega un entry y devuelve la representación pública.
func (s *Service) Append(in AppendInput) (*EntryView, error) {
	if in.Realm == "" {
		return nil, screrrors.New(screrrors.CodeRealmRequired, "realm required")
	}
	e, err := s.store.Append(in)
	if err != nil {
		return nil, screrrors.Wrap(screrrors.CodeInternal, "append entry", err)
	}
	return entryView(e), nil
}

// AccessLog retorna la bitácora ciudadana de un cedula.
func (s *Service) AccessLog(realm, citizenID string) ([]*EntryView, error) {
	if citizenID == "" {
		return nil, screrrors.New(screrrors.CodeBadRequest, "citizenId required")
	}
	src := s.store.AccessLogByCitizen(realm, citizenID)
	out := make([]*EntryView, len(src))
	for i, e := range src {
		out[i] = entryView(e)
	}
	return out, nil
}

// VerifyResult devuelve el resultado de la verificación de la cadena.
type VerifyResult struct {
	Realm        string `json:"realm"`
	Total        int    `json:"total"`
	Valid        bool   `json:"valid"`
	FirstInvalid int    `json:"firstInvalid,omitempty"` // -1 si toda íntegra
	Error        string `json:"error,omitempty"`
}

// Verify corre el verificador completo de la cadena del realm.
func (s *Service) Verify(realm string) *VerifyResult {
	idx, err := s.store.Verify(realm)
	all := s.store.All(realm)
	r := &VerifyResult{Realm: realm, Total: len(all)}
	if err != nil {
		r.Valid = false
		r.FirstInvalid = idx
		r.Error = err.Error()
		return r
	}
	r.Valid = true
	r.FirstInvalid = -1
	return r
}

// EntryView es la representación pública en JSON (con hashes en hex).
type EntryView struct {
	ID              string `json:"id"`
	Realm           string `json:"realm"`
	TS              string `json:"ts"`
	RequesterMember string `json:"requesterMember"`
	TargetMember    string `json:"targetMember"`
	Service         string `json:"service"`
	Version         string `json:"version"`
	CitizenID       string `json:"citizenId,omitempty"`
	Purpose         string `json:"purpose"`
	RequestID       string `json:"requestId"`
	Status          int    `json:"status"`
	PrevHash        string `json:"prevHash"`
	EntryHash       string `json:"entryHash"`
}

func entryView(e *auditstore.Entry) *EntryView {
	return &EntryView{
		ID: e.ID, Realm: e.Realm, TS: e.TS.Format("2006-01-02T15:04:05.000Z07:00"),
		RequesterMember: e.RequesterMember, TargetMember: e.TargetMember,
		Service: e.Service, Version: e.Version,
		CitizenID: e.CitizenID, Purpose: e.Purpose,
		RequestID: e.RequestID, Status: e.Status,
		PrevHash: hex(e.PrevHash), EntryHash: hex(e.EntryHash),
	}
}

const hexChars = "0123456789abcdef"

func hex(b []byte) string {
	out := make([]byte, len(b)*2)
	for i, c := range b {
		out[i*2] = hexChars[c>>4]
		out[i*2+1] = hexChars[c&0x0f]
	}
	return string(out)
}
