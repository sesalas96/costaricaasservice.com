// Package store es la capa de acceso a datos del cri-svc-interop-hub.
// Multi-tenancy schema-per-realm (mismo patrón que iduc-identity).
package store

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/devsebas/saascr/libs/cri-lib-shared/realm"
)

type Store struct {
	mu    sync.RWMutex
	dsn   string
	pools map[string]*pgxpool.Pool
	maxC  int32
	minC  int32
}

func New(dsn string, maxConns, minConns int32) *Store {
	return &Store{dsn: dsn, pools: make(map[string]*pgxpool.Pool), maxC: maxConns, minC: minConns}
}

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
		_, err := conn.Exec(ctx, "SET search_path TO "+pgx.Identifier{schema}.Sanitize())
		return err == nil
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect pool: %w", err)
	}
	s.pools[realmSlug] = pool
	return pool, nil
}

func (s *Store) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, p := range s.pools {
		p.Close()
	}
	s.pools = nil
}

var ErrNotFound = errors.New("store: not found")
