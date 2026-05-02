package interop

import (
	"errors"
	"sync"
)

// Registry mantiene los handlers que el svc-member expone vía interop.
// El binary del member llama Register en main(), y luego el security-server
// (cuando recibe un request entrante) consulta este registry vía un endpoint
// HTTP local del member.
type Registry struct {
	mu       sync.RWMutex
	handlers map[string]ServiceHandler // key = "<service>:<version>"
}

// NewRegistry crea un registry vacío.
func NewRegistry() *Registry {
	return &Registry{handlers: make(map[string]ServiceHandler)}
}

// Register vincula un handler a un servicio+versión específicos.
// Si ya existe, sobreescribe.
func (r *Registry) Register(service, version string, h ServiceHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[service+":"+version] = h
}

// Lookup retorna el handler o error si no existe.
func (r *Registry) Lookup(service, version string) (ServiceHandler, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	h, ok := r.handlers[service+":"+version]
	if !ok {
		return nil, errors.New("interop: service not registered")
	}
	return h, nil
}

// Services lista todos los pares (service,version) registrados.
func (r *Registry) Services() [][2]string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([][2]string, 0, len(r.handlers))
	for k := range r.handlers {
		var sv [2]string
		for i, c := range k {
			if c == ':' {
				sv[0] = k[:i]
				sv[1] = k[i+1:]
				break
			}
		}
		out = append(out, sv)
	}
	return out
}
