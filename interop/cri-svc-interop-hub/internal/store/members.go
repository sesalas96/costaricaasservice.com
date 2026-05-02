package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type Member struct {
	ID          string
	MemberID    string
	DisplayName string
	Description string
	PublicKey   string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CreateMember inserta un member nuevo en el realm.
func (s *Store) CreateMember(ctx context.Context, realm, memberID, displayName, description, publicKey string) (*Member, error) {
	pool, err := s.PoolFor(ctx, realm)
	if err != nil {
		return nil, err
	}
	const q = `
		INSERT INTO members (member_id, display_name, description, public_key)
		VALUES ($1, $2, $3, $4)
		RETURNING id, member_id, display_name, description, public_key, status, created_at, updated_at
	`
	var m Member
	if err := pool.QueryRow(ctx, q, memberID, displayName, description, publicKey).Scan(
		&m.ID, &m.MemberID, &m.DisplayName, &m.Description, &m.PublicKey, &m.Status, &m.CreatedAt, &m.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &m, nil
}

// GetMemberBySlug busca un member por su member_id (slug).
func (s *Store) GetMemberBySlug(ctx context.Context, realm, memberID string) (*Member, error) {
	pool, err := s.PoolFor(ctx, realm)
	if err != nil {
		return nil, err
	}
	const q = `
		SELECT id, member_id, display_name, description, public_key, status, created_at, updated_at
		FROM members WHERE member_id = $1
	`
	var m Member
	err = pool.QueryRow(ctx, q, memberID).Scan(
		&m.ID, &m.MemberID, &m.DisplayName, &m.Description, &m.PublicKey, &m.Status, &m.CreatedAt, &m.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// ListMembers retorna los members del realm con paginación opcional.
func (s *Store) ListMembers(ctx context.Context, realm string, limit, offset int) ([]Member, error) {
	pool, err := s.PoolFor(ctx, realm)
	if err != nil {
		return nil, err
	}
	const q = `
		SELECT id, member_id, display_name, description, public_key, status, created_at, updated_at
		FROM members
		ORDER BY member_id
		LIMIT $1 OFFSET $2
	`
	rows, err := pool.Query(ctx, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Member, 0, limit)
	for rows.Next() {
		var m Member
		if err := rows.Scan(&m.ID, &m.MemberID, &m.DisplayName, &m.Description, &m.PublicKey, &m.Status, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

// UpdateMemberStatus cambia el estado de un member.
func (s *Store) UpdateMemberStatus(ctx context.Context, realm, memberID, status string) error {
	pool, err := s.PoolFor(ctx, realm)
	if err != nil {
		return err
	}
	const q = `UPDATE members SET status = $1, updated_at = NOW() WHERE member_id = $2`
	tag, err := pool.Exec(ctx, q, status, memberID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
