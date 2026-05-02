package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type Service struct {
	ID          string
	MemberID    string // FK al members.id (UUID)
	MemberSlug  string // member_id slug (denormalizado en queries con JOIN)
	ServiceID   string
	Version     string
	Description string
	SchemaURL   string
	Exposed     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// PublishService crea/actualiza un servicio publicado por un member.
// Usa upsert para que (member_id, service_id, version) sea idempotente.
func (s *Store) PublishService(ctx context.Context, realm, memberSlug, serviceID, version, description, schemaURL string) (*Service, error) {
	pool, err := s.PoolFor(ctx, realm)
	if err != nil {
		return nil, err
	}
	const q = `
		WITH m AS (SELECT id FROM members WHERE member_id = $1)
		INSERT INTO services (member_id, service_id, version, description, schema_url, exposed)
		SELECT m.id, $2, $3, $4, $5, TRUE FROM m
		ON CONFLICT (member_id, service_id, version) DO UPDATE
			SET description = EXCLUDED.description,
			    schema_url  = EXCLUDED.schema_url,
			    exposed     = TRUE,
			    updated_at  = NOW()
		RETURNING id, member_id, service_id, version, description, schema_url, exposed, created_at, updated_at
	`
	var sv Service
	err = pool.QueryRow(ctx, q, memberSlug, serviceID, version, description, schemaURL).Scan(
		&sv.ID, &sv.MemberID, &sv.ServiceID, &sv.Version, &sv.Description, &sv.SchemaURL, &sv.Exposed, &sv.CreatedAt, &sv.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound // member no existe
	}
	if err != nil {
		return nil, err
	}
	sv.MemberSlug = memberSlug
	return &sv, nil
}

// LookupService busca un servicio en el catálogo.
func (s *Store) LookupService(ctx context.Context, realm, memberSlug, serviceID, version string) (*Service, error) {
	pool, err := s.PoolFor(ctx, realm)
	if err != nil {
		return nil, err
	}
	const q = `
		SELECT s.id, s.member_id, m.member_id, s.service_id, s.version,
		       s.description, s.schema_url, s.exposed, s.created_at, s.updated_at
		FROM services s
		JOIN members m ON m.id = s.member_id
		WHERE m.member_id = $1 AND s.service_id = $2 AND s.version = $3 AND s.exposed = TRUE
	`
	var sv Service
	err = pool.QueryRow(ctx, q, memberSlug, serviceID, version).Scan(
		&sv.ID, &sv.MemberID, &sv.MemberSlug, &sv.ServiceID, &sv.Version,
		&sv.Description, &sv.SchemaURL, &sv.Exposed, &sv.CreatedAt, &sv.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &sv, nil
}

// ListServicesByMember lista los servicios expuestos por un member.
func (s *Store) ListServicesByMember(ctx context.Context, realm, memberSlug string) ([]Service, error) {
	pool, err := s.PoolFor(ctx, realm)
	if err != nil {
		return nil, err
	}
	const q = `
		SELECT s.id, s.member_id, m.member_id, s.service_id, s.version,
		       s.description, s.schema_url, s.exposed, s.created_at, s.updated_at
		FROM services s
		JOIN members m ON m.id = s.member_id
		WHERE m.member_id = $1 AND s.exposed = TRUE
		ORDER BY s.service_id, s.version
	`
	rows, err := pool.Query(ctx, q, memberSlug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Service, 0, 8)
	for rows.Next() {
		var sv Service
		if err := rows.Scan(
			&sv.ID, &sv.MemberID, &sv.MemberSlug, &sv.ServiceID, &sv.Version,
			&sv.Description, &sv.SchemaURL, &sv.Exposed, &sv.CreatedAt, &sv.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, sv)
	}
	return out, rows.Err()
}

// UnpublishService marca un servicio como no expuesto (soft delete).
func (s *Store) UnpublishService(ctx context.Context, realm, memberSlug, serviceID, version string) error {
	pool, err := s.PoolFor(ctx, realm)
	if err != nil {
		return err
	}
	const q = `
		UPDATE services SET exposed = FALSE, updated_at = NOW()
		WHERE member_id = (SELECT id FROM members WHERE member_id = $1)
		  AND service_id = $2 AND version = $3
	`
	tag, err := pool.Exec(ctx, q, memberSlug, serviceID, version)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
