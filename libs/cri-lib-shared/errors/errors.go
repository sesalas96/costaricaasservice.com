// Package errors define el error estándar AppError de saascr.
// Todo error que cruce un boundary HTTP o de servicio se modela como AppError
// para que httpx.Fail pueda mapearlo a una respuesta envelope consistente.
package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// Code es un identificador semántico estable. Los clientes pueden ramificar
// lógica por Code; los Mensajes pueden cambiar.
type Code string

const (
	CodeBadRequest      Code = "BAD_REQUEST"
	CodeValidation      Code = "VALIDATION_ERROR"
	CodeUnauthorized    Code = "UNAUTHORIZED"
	CodeForbidden       Code = "FORBIDDEN"
	CodeNotFound        Code = "NOT_FOUND"
	CodeConflict        Code = "CONFLICT"
	CodeRateLimited     Code = "RATE_LIMITED"
	CodeInternal        Code = "INTERNAL"
	CodeUnavailable     Code = "UNAVAILABLE"
	CodeRealmRequired   Code = "REALM_REQUIRED"
	CodeRealmForbidden  Code = "REALM_FORBIDDEN"
	CodeInteropDenied   Code = "INTEROP_DENIED"
	CodeAuditMismatch   Code = "AUDIT_MISMATCH"
	CodeSigningRejected Code = "SIGNING_REJECTED"
)

// AppError es el error estándar del monorepo.
type AppError struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
	Cause   error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.Cause }

// New construye un AppError con status sugerido por el código.
func New(code Code, msg string) *AppError {
	return &AppError{Code: code, Message: msg, Status: statusFor(code)}
}

// Wrap envuelve un error nativo en un AppError.
func Wrap(code Code, msg string, cause error) *AppError {
	return &AppError{Code: code, Message: msg, Status: statusFor(code), Cause: cause}
}

// As intenta extraer un *AppError del error wrapping chain.
func As(err error) (*AppError, bool) {
	var e *AppError
	if errors.As(err, &e) {
		return e, true
	}
	return nil, false
}

func statusFor(c Code) int {
	switch c {
	case CodeBadRequest, CodeValidation, CodeRealmRequired:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden, CodeRealmForbidden, CodeInteropDenied, CodeSigningRejected:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeConflict:
		return http.StatusConflict
	case CodeRateLimited:
		return http.StatusTooManyRequests
	case CodeUnavailable:
		return http.StatusServiceUnavailable
	case CodeInternal, CodeAuditMismatch:
		fallthrough
	default:
		return http.StatusInternalServerError
	}
}
