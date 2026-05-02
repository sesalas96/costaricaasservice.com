// Package upstream describe los servicios internos a los que el gateway
// hace reverse proxy y la regla de matching por path.
package upstream

import (
	"errors"
	"net/url"
	"strings"
)

// Route describe una regla: cualquier request cuyo path empiece con Prefix
// se enruta al Target. El Prefix se conserva en la URL upstream a menos que
// StripPrefix sea true.
type Route struct {
	Prefix      string
	Target      *url.URL
	StripPrefix bool
}

// Routes mantiene una lista ordenada (la primera coincidencia gana).
type Routes struct {
	rules []Route
}

// New construye Routes a partir de una lista de specs.
//
// Cada spec es del tipo: {prefix, targetURL, strip}.
type Spec struct {
	Prefix      string `mapstructure:"prefix"`
	TargetURL   string `mapstructure:"target_url"`
	StripPrefix bool   `mapstructure:"strip_prefix"`
}

// New parsea un slice de Spec a Routes.
func New(specs []Spec) (*Routes, error) {
	rules := make([]Route, 0, len(specs))
	for _, s := range specs {
		if s.Prefix == "" || s.TargetURL == "" {
			return nil, errors.New("upstream: prefix and target_url are required")
		}
		u, err := url.Parse(s.TargetURL)
		if err != nil {
			return nil, err
		}
		rules = append(rules, Route{Prefix: s.Prefix, Target: u, StripPrefix: s.StripPrefix})
	}
	return &Routes{rules: rules}, nil
}

// Match retorna la Route que coincide con el path o nil.
func (r *Routes) Match(path string) *Route {
	for i := range r.rules {
		if strings.HasPrefix(path, r.rules[i].Prefix) {
			return &r.rules[i]
		}
	}
	return nil
}

// Rules expone las reglas (solo lectura).
func (r *Routes) Rules() []Route { return append([]Route(nil), r.rules...) }
