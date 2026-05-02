// Package auditstore mantiene la bitácora hash-chained de transacciones
// inter-member del realm. Implementación in-memory para MVP — se reemplaza
// por Postgres + Kafka events en Fase 2 (ver ADR-0002).
//
// Cada entry referencia el hash del entry anterior. Un escritor que intente
// modificar una entry pasada rompe la cadena de TODAS las posteriores.
package auditstore

import (
	"sync"
	"time"

	"github.com/devsebas/costaricaasservice/libs/cri-lib-crypto/hashchain"
	"github.com/devsebas/costaricaasservice/libs/cri-lib-shared/idgen"
)

// Entry es una transacción inter-member registrada en el log.
type Entry struct {
	ID              string    `json:"id"`               // ULID monotónico
	Realm           string    `json:"realm"`
	TS              time.Time `json:"ts"`
	RequesterMember string    `json:"requesterMember"`
	TargetMember    string    `json:"targetMember"`
	Service         string    `json:"service"`
	Version         string    `json:"version"`
	CitizenID       string    `json:"citizenId,omitempty"`
	Purpose         string    `json:"purpose"`
	RequestID       string    `json:"requestId"`
	Status          int       `json:"status"`
	PrevHash        []byte    `json:"prevHash"`         // hash del entry previo
	EntryHash       []byte    `json:"entryHash"`        // hash de este entry
}

// Store mantiene la cadena por realm en memoria.
type Store struct {
	mu      sync.RWMutex
	chains  map[string][]*Entry // realm → entries en orden
	prevTip map[string][]byte   // realm → último entry_hash (para encadenar el siguiente)
}

// New crea un Store vacío.
func New() *Store {
	return &Store{
		chains:  make(map[string][]*Entry),
		prevTip: make(map[string][]byte),
	}
}

// AppendInput agrupa los datos de un evento entrante.
type AppendInput struct {
	Realm           string
	RequesterMember string
	TargetMember    string
	Service         string
	Version         string
	CitizenID       string
	Purpose         string
	RequestID       string
	Status          int
}

// Append agrega un entry nuevo a la cadena del realm. Calcula prev_hash y
// entry_hash usando los helpers de cri-lib-crypto/hashchain.
func (s *Store) Append(in AppendInput) (*Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	prev, ok := s.prevTip[in.Realm]
	if !ok {
		prev = hashchain.Genesis(in.Realm)
	}

	now := time.Now().UTC()
	e := &Entry{
		ID:              idgen.New(),
		Realm:           in.Realm,
		TS:              now,
		RequesterMember: in.RequesterMember,
		TargetMember:    in.TargetMember,
		Service:         in.Service,
		Version:         in.Version,
		CitizenID:       in.CitizenID,
		Purpose:         in.Purpose,
		RequestID:       in.RequestID,
		Status:          in.Status,
		PrevHash:        append([]byte(nil), prev...),
	}
	payload := s.payloadOf(e)
	hash, err := hashchain.EntryHash(prev, payload)
	if err != nil {
		return nil, err
	}
	e.EntryHash = hash

	s.chains[in.Realm] = append(s.chains[in.Realm], e)
	s.prevTip[in.Realm] = hash
	return e, nil
}

// payloadOf canonicaliza el contenido del entry para el cálculo del hash.
// Excluye PrevHash y EntryHash (esos no se hashean dentro del payload).
func (s *Store) payloadOf(e *Entry) map[string]any {
	return map[string]any{
		"id":               e.ID,
		"realm":            e.Realm,
		"ts":               e.TS.UnixNano(),
		"requester_member": e.RequesterMember,
		"target_member":    e.TargetMember,
		"service":          e.Service,
		"version":          e.Version,
		"citizen_id":       e.CitizenID,
		"purpose":          e.Purpose,
		"request_id":       e.RequestID,
		"status":           e.Status,
	}
}

// AccessLogByCitizen retorna las entries cuyo CitizenID matchea, en orden de TS asc.
func (s *Store) AccessLogByCitizen(realm, citizenID string) []*Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Entry, 0, 8)
	for _, e := range s.chains[realm] {
		if e.CitizenID == citizenID {
			out = append(out, e)
		}
	}
	return out
}

// All devuelve toda la cadena de un realm (orden inserción). Útil para verify.
func (s *Store) All(realm string) []*Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	src := s.chains[realm]
	out := make([]*Entry, len(src))
	copy(out, src)
	return out
}

// Verify recorre la cadena entera y retorna nil si toda está íntegra.
// Caso de mismatch retorna (idx_invalido, error).
func (s *Store) Verify(realm string) (int, error) {
	s.mu.RLock()
	entries := make([]*Entry, len(s.chains[realm]))
	copy(entries, s.chains[realm])
	s.mu.RUnlock()

	chainEntries := make([]hashchain.Entry, len(entries))
	for i, e := range entries {
		chainEntries[i] = hashchain.Entry{
			PrevHash:  e.PrevHash,
			EntryHash: e.EntryHash,
			Payload:   s.payloadOf(e),
		}
	}
	return hashchain.Verify(realm, chainEntries)
}
