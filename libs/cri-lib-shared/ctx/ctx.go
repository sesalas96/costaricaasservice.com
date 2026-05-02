// Package ctx contiene helpers para inyectar y leer valores comunes
// del context.Context: requestId, realm, principal stub.
package ctx

import "context"

type ctxKey int

const (
	keyRequestID ctxKey = iota
	keyRealm
)

// WithRequestID coloca el requestId en el contexto.
func WithRequestID(parent context.Context, id string) context.Context {
	return context.WithValue(parent, keyRequestID, id)
}

// RequestID extrae el requestId. Devuelve "" si no existe.
func RequestID(c context.Context) string {
	if v, ok := c.Value(keyRequestID).(string); ok {
		return v
	}
	return ""
}

// WithRealm coloca el realm activo en el contexto.
func WithRealm(parent context.Context, realm string) context.Context {
	return context.WithValue(parent, keyRealm, realm)
}

// Realm extrae el realm. Devuelve "" si no existe.
func Realm(c context.Context) string {
	if v, ok := c.Value(keyRealm).(string); ok {
		return v
	}
	return ""
}
