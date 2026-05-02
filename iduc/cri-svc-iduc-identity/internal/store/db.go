// Package store es la capa de acceso a datos de cri-svc-iduc-identity.
//
// Multi-tenancy schema-per-realm: el Store mantiene un Pool por realm.
// Cada conexión adquirida del pool ejecuta `SET search_path` al schema del realm
// vía un `BeforeAcquire` hook configurado en pgxpool.Config. Esto evita
// fugas accidentales entre tenants — una conexión no puede ver el schema
// de un realm distinto.
package store

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/devsebas/costaricaasservice/libs/cri-lib-shared/realm"
)

// Store mantiene un pool por realm.
type Store struct {
	mu     sync.RWMutex
	dsn    string
	pools  map[string]*pgxpool.Pool // key: realm slug
	maxC   int32
	minC   int32
}

// New construye un Store.
func New(dsn string, maxConns, minConns int32) *Store {
	return &Store{
		dsn:   dsn,
		pools: make(map[string]*pgxpool.Pool),
		maxC:  maxConns,
		minC:  minConns,
	}
}

// PoolFor devuelve el pool del realm, creándolo si es la primera vez.
// El BeforeAcquire de cada conexión hace SET search_path al schema del realm.
func (s *Store) PoolFor(ctx context.Context, realmSlug string) (*pgxpool.Pool, error) {
	if err := realm.Validate(realmSlug); err != nil {
		return nil, err
	}
	schema := realm.SchemaName(realmSlug)

	s.mu.RLock()
	if pool, ok := s.pools[realmSlug]; ok {
		s.mu.RUnlock()
		return pool, nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()
	if pool, ok := s.pools[realmSlug]; ok {
		return pool, nil
	}

	cfg, err := pgxpool.ParseConfig(s.dsn)
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}
	cfg.MaxConns = s.maxC
	cfg.MinConns = s.minC

	cfg.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		// SET search_path explícito en cada acquire — defensivo, evita que
		// una sesión reusada apunte a otro realm.
		if _, err := conn.Exec(ctx, "SET search_path TO "+pgx.Identifier{schema}.Sanitize()); err != nil {
			return false
		}
		return true
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect pool: %w", err)
	}
	s.pools[realmSlug] = pool
	return pool, nil
}

// Close cierra todos los pools.
func (s *Store) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, p := range s.pools {
		p.Close()
	}
	s.pools = nil
}

// ErrNotFound es retornado cuando una entidad esperada no existe.
var ErrNotFound = errors.New("store: not found")
