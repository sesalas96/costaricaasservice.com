package auth

import "context"

// Principal representa al sujeto autenticado tras el gateway: ciudadano,
// operador de member, admin de realm o admin saascr.
type Principal struct {
	Sub    string // identificador único (citizen_id, operator_id, etc.)
	Roles  []Role // roles asignados
	Realm  string // realm en el que opera este principal
	Member string // miembro asociado (si aplica). Vacío para citizen.
	JTI    string // JWT id (numeric, para revocación)
}

type ctxKey int

const principalKey ctxKey = 0

// WithPrincipal coloca el Principal en el context.
func WithPrincipal(parent context.Context, p *Principal) context.Context {
	return context.WithValue(parent, principalKey, p)
}

// FromContext extrae el Principal o nil.
func FromContext(c context.Context) *Principal {
	p, _ := c.Value(principalKey).(*Principal)
	return p
}
