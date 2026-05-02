package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

// Session representa una sesión activa.
type Session struct {
	ID                 string
	CitizenID          string
	RefreshTokenHash   []byte
	JTI                int64
	IssuedAt           time.Time
	ExpiresAt          time.Time
	RevokedAt          *time.Time
	RotatedTo          *string
}

// NextJTI consume la sequence jti_seq y devuelve el siguiente JTI numérico.
func (s *Store) NextJTI(ctx context.Context, realm string) (int64, error) {
	pool, err := s.PoolFor(ctx, realm)
	if err != nil {
		return 0, err
	}
	var jti int64
	if err := pool.QueryRow(ctx, "SELECT nextval('jti_seq')").Scan(&jti); err != nil {
		return 0, err
	}
	return jti, nil
}

// CreateSession inserta una sesión nueva con su refresh_token_hash y JTI.
func (s *Store) CreateSession(ctx context.Context, realm, citizenID string, refreshHash []byte, jti int64, expiresAt time.Time) (*Session, error) {
	pool, err := s.PoolFor(ctx, realm)
	if err != nil {
		return nil, err
	}
	const q = `
		INSERT INTO sessions (citizen_id, refresh_token_hash, jti, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, citizen_id, refresh_token_hash, jti, issued_at, expires_at, revoked_at, rotated_to
	`
	var sess Session
	err = pool.QueryRow(ctx, q, citizenID, refreshHash, jti, expiresAt).Scan(
		&sess.ID, &sess.CitizenID, &sess.RefreshTokenHash, &sess.JTI,
		&sess.IssuedAt, &sess.ExpiresAt, &sess.RevokedAt, &sess.RotatedTo,
	)
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

// GetSessionByRefresh busca por refresh_token_hash. Retorna ErrNotFound si no existe.
func (s *Store) GetSessionByRefresh(ctx context.Context, realm string, refreshHash []byte) (*Session, error) {
	pool, err := s.PoolFor(ctx, realm)
	if err != nil {
		return nil, err
	}
	const q = `
		SELECT id, citizen_id, refresh_token_hash, jti, issued_at, expires_at, revoked_at, rotated_to
		FROM sessions WHERE refresh_token_hash = $1
	`
	var sess Session
	err = pool.QueryRow(ctx, q, refreshHash).Scan(
		&sess.ID, &sess.CitizenID, &sess.RefreshTokenHash, &sess.JTI,
		&sess.IssuedAt, &sess.ExpiresAt, &sess.RevokedAt, &sess.RotatedTo,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

// RotateSession marca la sesión actual como revocada y crea una nueva
// vinculada por rotated_to. Retorna la sesión nueva.
func (s *Store) RotateSession(ctx context.Context, realm, oldSessionID string, newRefreshHash []byte, newJTI int64, newExpiresAt time.Time, citizenID string) (*Session, error) {
	pool, err := s.PoolFor(ctx, realm)
	if err != nil {
		return nil, err
	}
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	const insertNew = `
		INSERT INTO sessions (citizen_id, refresh_token_hash, jti, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, citizen_id, refresh_token_hash, jti, issued_at, expires_at, revoked_at, rotated_to
	`
	var sess Session
	if err := tx.QueryRow(ctx, insertNew, citizenID, newRefreshHash, newJTI, newExpiresAt).Scan(
		&sess.ID, &sess.CitizenID, &sess.RefreshTokenHash, &sess.JTI,
		&sess.IssuedAt, &sess.ExpiresAt, &sess.RevokedAt, &sess.RotatedTo,
	); err != nil {
		return nil, err
	}

	const revokeOld = `
		UPDATE sessions SET revoked_at = NOW(), rotated_to = $1
		WHERE id = $2 AND revoked_at IS NULL
	`
	if _, err := tx.Exec(ctx, revokeOld, sess.ID, oldSessionID); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &sess, nil
}

// RevokeSession marca una sesión como revocada y agrega su JTI a revoked_jtis.
// Toma el JTI vigente desde la sesión.
func (s *Store) RevokeSession(ctx context.Context, realm, sessionID string) error {
	pool, err := s.PoolFor(ctx, realm)
	if err != nil {
		return err
	}
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var jti int64
	var expiresAt time.Time
	const fetch = `SELECT jti, expires_at FROM sessions WHERE id = $1`
	if err := tx.QueryRow(ctx, fetch, sessionID).Scan(&jti, &expiresAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	const revoke = `UPDATE sessions SET revoked_at = NOW() WHERE id = $1 AND revoked_at IS NULL`
	if _, err := tx.Exec(ctx, revoke, sessionID); err != nil {
		return err
	}

	const addRevoked = `
		INSERT INTO revoked_jtis (jti, expires_at)
		VALUES ($1, $2)
		ON CONFLICT (jti) DO NOTHING
	`
	if _, err := tx.Exec(ctx, addRevoked, jti, expiresAt); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// RevokedSnapshot retorna la lista de JTIs revocados aún no expirados.
// El gateway lo consume para reconstruir su Roaring bitmap.
func (s *Store) RevokedSnapshot(ctx context.Context, realm string) ([]int64, error) {
	pool, err := s.PoolFor(ctx, realm)
	if err != nil {
		return nil, err
	}
	const q = `SELECT jti FROM revoked_jtis WHERE expires_at > NOW() ORDER BY jti`
	rows, err := pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]int64, 0, 64)
	for rows.Next() {
		var jti int64
		if err := rows.Scan(&jti); err != nil {
			return nil, err
		}
		out = append(out, jti)
	}
	return out, rows.Err()
}
