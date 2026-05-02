// Package service contiene la lógica de negocio del servicio.
// Recibe dependencias (store, eventos, clientes externos) por constructor;
// no toca HTTP ni JSON directamente — los handlers se encargan.
package service

// Service es el agregador de business logic del servicio.
// Sustituir cuando se modele el dominio real.
type Service struct{}

// New construye un Service nuevo.
func New() *Service { return &Service{} }
