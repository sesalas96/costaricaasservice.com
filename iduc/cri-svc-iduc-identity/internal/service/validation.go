package service

import (
	"regexp"
	"strings"

	screrrors "github.com/devsebas/saascr/libs/cri-lib-shared/errors"
)

// Cédula CR formato simple: 9 dígitos opcionales con guiones (1-1234-5678 o 112345678).
// Acepta también DNIs/passports en otros realms — la validación estricta por
// jurisdicción se hará por configuración del realm en una fase posterior.
var cedulaRe = regexp.MustCompile(`^[A-Z0-9-]{6,20}$`)

func validateCedula(c string) error {
	c = strings.ToUpper(strings.TrimSpace(c))
	if !cedulaRe.MatchString(c) {
		return screrrors.New(screrrors.CodeValidation, "invalid cedula format")
	}
	return nil
}

var emailRe = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

func validateEmail(e string) error {
	e = strings.TrimSpace(e)
	if !emailRe.MatchString(e) || len(e) > 254 {
		return screrrors.New(screrrors.CodeValidation, "invalid email")
	}
	return nil
}
