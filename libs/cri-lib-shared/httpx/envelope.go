// Package httpx contiene helpers para escribir respuestas HTTP en el
// formato envelope estándar de costaricaasservice: {data, meta} en éxito, {error, meta}
// en error. Siempre con requestId.
package httpx

import (
	"encoding/json"
	"log/slog"
	"net/http"

	scrctx "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/ctx"
	screrrors "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/errors"
)

// Meta contiene metadatos comunes a toda respuesta.
type Meta struct {
	RequestID string `json:"requestId"`
	Realm     string `json:"realm,omitempty"`
}

// Success es el envelope de éxito.
type Success struct {
	Data any  `json:"data"`
	Meta Meta `json:"meta"`
}

// ErrorBody es el cuerpo de error dentro del envelope.
type ErrorBody struct {
	Code    screrrors.Code `json:"code"`
	Message string         `json:"message"`
}

// Failure es el envelope de error.
type Failure struct {
	Error ErrorBody `json:"error"`
	Meta  Meta      `json:"meta"`
}

// OK escribe 200 con el envelope {data, meta}.
func OK(w http.ResponseWriter, r *http.Request, data any) {
	write(w, r, http.StatusOK, Success{Data: data, Meta: metaFrom(r)})
}

// Created escribe 201 con el envelope {data, meta}.
func Created(w http.ResponseWriter, r *http.Request, data any) {
	write(w, r, http.StatusCreated, Success{Data: data, Meta: metaFrom(r)})
}

// NoContent escribe 204.
func NoContent(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// Fail escribe el error mapeando AppError → status; cualquier otro error → 500.
func Fail(w http.ResponseWriter, r *http.Request, err error) {
	if appErr, ok := screrrors.As(err); ok {
		write(w, r, appErr.Status, Failure{
			Error: ErrorBody{Code: appErr.Code, Message: appErr.Message},
			Meta:  metaFrom(r),
		})
		return
	}
	slog.ErrorContext(r.Context(), "unhandled error", "err", err.Error(), "requestId", scrctx.RequestID(r.Context()))
	write(w, r, http.StatusInternalServerError, Failure{
		Error: ErrorBody{Code: screrrors.CodeInternal, Message: "internal server error"},
		Meta:  metaFrom(r),
	})
}

func metaFrom(r *http.Request) Meta {
	return Meta{
		RequestID: scrctx.RequestID(r.Context()),
		Realm:     scrctx.Realm(r.Context()),
	}
}

func write(w http.ResponseWriter, _ *http.Request, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("encode response", "err", err.Error())
	}
}
