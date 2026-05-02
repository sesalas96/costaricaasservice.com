package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

// Citizen es la representación en memoria de la fila citizens.
type Citizen struct {
	ID           string
	Cedula       string
	Email        string
	PasswordHash string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// CreateCitizen inserta un ciudadano nuevo. Retorna error si la cédula o
// email ya existen (unique constraint).
func (s *Store) CreateCitizen(ctx context.Context, realm, cedula, email, passwordHash string) (*Citizen, error) {
	pool, err := s.PoolFor(ctx, realm)
	if err != nil {
		return nil, err
	}
	const q = `
		INSERT INTO citizens (cedula, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, cedula, email, password_hash, status, created_at, updated_at
	`
	var c Citizen
	err = pool.QueryRow(ctx, q, cedula, email, passwordHash).Scan(
		&c.ID, &c.Cedula, &c.Email, &c.PasswordHash, &c.Status, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// GetCitizenByEmail busca un ciudadano por email dentro del realm.
func (s *Store) GetCitizenByEmail(ctx context.Context, realm, email string) (*Citizen, error) {
	pool, err := s.PoolFor(ctx, realm)
	if err != nil {
		return nil, err
	}
	const q = `
		SELECT id, cedula, email, password_hash, status, created_at, updated_at
		FROM citizens WHERE email = $1
	`
	var c Citizen
	err = pool.QueryRow(ctx, q, email).Scan(
		&c.ID, &c.Cedula, &c.Email, &c.PasswordHash, &c.Status, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// GetCitizenByID busca un ciudadano por id.
func (s *Store) GetCitizenByID(ctx context.Context, realm, id string) (*Citizen, error) {
	pool, err := s.PoolFor(ctx, realm)
	if err != nil {
		return nil, err
	}
	const q = `
		SELECT id, cedula, email, password_hash, status, created_at, updated_at
		FROM citizens WHERE id = $1
	`
	var c Citizen
	err = pool.QueryRow(ctx, q, id).Scan(
		&c.ID, &c.Cedula, &c.Email, &c.PasswordHash, &c.Status, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}
