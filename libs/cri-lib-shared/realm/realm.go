// Package realm define el formato y validación del identificador de realm
// (jurisdicción soberana). Un realm es un slug ASCII en minúsculas con guiones,
// 3 a 32 chars (e.g., cr-prod, sv-prod, demo).
package realm

import (
	"regexp"

	screrrors "github.com/devsebas/saascr/libs/cri-lib-shared/errors"
)

var validRe = regexp.MustCompile(`^[a-z][a-z0-9-]{2,31}$`)

// Validate retorna nil si el slug de realm es válido.
func Validate(r string) error {
	if r == "" {
		return screrrors.New(screrrors.CodeRealmRequired, "realm is required")
	}
	if !validRe.MatchString(r) {
		return screrrors.New(screrrors.CodeRealmRequired, "realm has invalid format")
	}
	return nil
}

// SchemaName traduce un realm a su nombre de schema en Postgres.
// Reemplaza guiones por underscores: cr-prod → cr_prod.
func SchemaName(r string) string {
	out := make([]byte, 0, len(r))
	for i := 0; i < len(r); i++ {
		c := r[i]
		if c == '-' {
			c = '_'
		}
		out = append(out, c)
	}
	return string(out)
}
