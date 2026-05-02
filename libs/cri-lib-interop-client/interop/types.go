// Package interop provee el SDK que un svc-member usa para:
//
//  1. Exponer servicios al hub (lado servidor del catálogo).
//  2. Consumir servicios remotos vía el cri-svc-security-server local.
//
// Inspirado en X-Road. Cada llamada inter-member produce un AuditEntry
// hash-chained según ADR-0002 y se firma con la clave del member origen
// según ADR-0003.
package interop

import "time"

// CallRequest representa una invocación inter-member desde el SDK.
// El SDK la envía al security-server local; este firma, enruta y vuelve.
type CallRequest struct {
	TargetMember string         // member destino (ej: "registro-civil")
	Service      string         // identificador en el catálogo (ej: "persons.get")
	Version      string         // semver del contrato (ej: "v1")
	Body         any            // payload JSON
	CitizenID    string         // si la consulta es sobre un ciudadano
	Purpose      string         // propósito declarado (ej: "prefill_tax_return")
	Timeout      time.Duration  // por defecto 5s
	Headers      map[string]any // metadata extra opcional
}

// CallResponse retorna el resultado del peer remoto y metadatos de auditoría.
type CallResponse struct {
	Status     int            // status HTTP del peer
	Body       []byte         // payload JSON crudo
	AuditID    string         // ULID del entry escrito en cri-svc-interop-audit
	EntryHash  []byte         // hash del entry (verificable contra Merkle root)
	Latency    time.Duration  // tiempo total
	Trace      map[string]any // metadata adicional opcional
}

// ServiceHandler es la firma que un member implementa para exponer un servicio.
// El SDK del lado servidor lo invoca tras validar firma + auditar.
type ServiceHandler func(ctx Context, body []byte) (response any, err error)

// Context se pasa al handler con metadatos del request entrante.
type Context interface {
	RequestID() string
	Realm() string
	RequesterMember() string
	CitizenID() string
	Purpose() string
}
