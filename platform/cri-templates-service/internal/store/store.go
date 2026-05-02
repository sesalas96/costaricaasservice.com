// Package store es la capa de acceso a datos. Usar pgx/v5 con Pool por servicio.
// Convención: los métodos exportados aceptan context.Context y respetan el realm
// vía SET search_path al adquirir conexión (multi-tenancy schema-per-realm).
package store

// Store es el aggregator del data access. Sustituir cuando se modele dominio.
type Store struct{}

// New construye un Store nuevo.
func New() *Store { return &Store{} }
