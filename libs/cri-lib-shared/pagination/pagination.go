// Package pagination parsea parámetros de paginación de URL y devuelve
// page/limit normalizados.
package pagination

import (
	"net/http"
	"strconv"
)

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

// Page representa una solicitud de paginación validada.
type Page struct {
	Page  int
	Limit int
}

// Offset devuelve el offset SQL equivalente.
func (p Page) Offset() int {
	if p.Page <= 1 {
		return 0
	}
	return (p.Page - 1) * p.Limit
}

// FromRequest extrae page=N&limit=M con defaults y clamping.
func FromRequest(r *http.Request) Page {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit <= 0 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	return Page{Page: page, Limit: limit}
}
