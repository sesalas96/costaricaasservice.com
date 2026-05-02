# ADR-0003 — KMS y custodia de claves: las claves privadas nunca salen del KMS

- **Estado:** Aceptado
- **Fecha:** 2026-05-02
- **Decisores:** equipo costaricaasservice

## Contexto

Las claves privadas que firman documentos en nombre del ciudadano son el activo más sensible del sistema. Su compromiso significa:

- Capacidad de un atacante de firmar documentos legalmente vinculantes en nombre de cualquier ciudadano.
- Pérdida de la equivalencia jurídica con la firma manuscrita — toda la propuesta de valor del producto se cae.
- Riesgo legal directo para `costaricaasservice` y para el realm.

El ciberataque de Conti/Hive a Costa Rica (abril–mayo 2022) demostró que la asunción "la BD interna está segura" no se sostiene. El diseño debe asumir compromiso del entorno de servicio.

## Decisión

### 1. Las claves privadas nunca salen del KMS

- Generación, almacenamiento y operaciones criptográficas (firma, descifrado) ocurren **dentro del KMS**.
- Los servicios envían `(payload_hash, intent, key_id)` al KMS y reciben `(signature)`. **Jamás** ven la clave privada.
- En ningún lugar de la base de datos, logs, traces u objetos en memoria del proceso aparece material de clave privada.

### 2. Stack por ambiente

| Ambiente | KMS | Notas |
|---|---|---|
| `local` / `dev` | HashiCorp Vault (Transit secrets engine), modo `dev` | Levantado por `make up`. Token: `dev-root-token`. **No usar fuera de dev.** |
| `qa` / `staging` | HashiCorp Vault (Transit), HA, sealed con Shamir | Auto-unseal opcional con KMS cloud. |
| `prod` | AWS KMS (regional) o GCP KMS, dependiendo del realm | Customer-managed keys (CMK) por realm. Audit logging integrado con CloudTrail / Cloud Audit Logs. |

### 3. Tipos de claves

| Tipo | Algoritmo | Custodia | Uso |
|---|---|---|---|
| Ciudadano (firma) | Ed25519 | KMS, una clave por (realm, citizen_id) | Firmar documentos en nombre del ciudadano tras autenticación + intent explícito. |
| Member (interop) | RSA-3072 o ECDSA P-384 | KMS o HSM del member, X.509 cert firmado por la CA del realm | Firmar JWS detached de cada request inter-member desde el security-server local. |
| Realm CA | RSA-4096 o EdDSA Ed25519 | HSM (en prod), KMS root (en qa) | Firmar certs de members. Rotación lenta. |
| JWT signing (gateway) | RSA-3072 (RS256) | KMS por servicio | Firmar access/refresh tokens. Rotación trimestral. |

### 4. Patrón de envelope encryption (data-at-rest)

Para datos sensibles **no-clave** (PII, expedientes), `cri-svc-iduc-keys` y `cri-svc-files` usan envelope encryption: una DEK (data encryption key) por registro, cifrada con la KEK (key encryption key) del KMS.

### 5. Intent + autenticación obligatoria para firmar

Antes de pedir al KMS que firme con la clave de un ciudadano, `cri-svc-iduc-signing` requiere:

1. JWT del ciudadano válido (no revocado).
2. Confirmación reciente (< 5min) del segundo factor (passkey/WebAuthn).
3. Intent firmado por el ciudadano: `{citizen_id, document_hash, purpose, ts}` autenticado en cliente.
4. Audit entry escrita **antes** de invocar al KMS.

El KMS solo se invoca si los 4 pasos pasan. Cualquier short-circuit es bloqueado por el servicio.

### 6. Rotación

| Material | Cadencia | Trigger |
|---|---|---|
| Ciudadano (firma) | A demanda del ciudadano, o al renovar la cédula (5 años) | UI en MiCR. |
| Member (interop) | Cada 12 meses, o ante incidente | Auto-renovación coordinada con el hub. |
| Realm CA | Cada 5 años, o ante incidente crítico | Procedimiento manual con rotación coordinada de members. |
| JWT signing | Cada 90 días | Cron en gateway. |

## Consecuencias

**Positivas**
- Un compromiso de la BD del servicio **no** compromete las claves privadas.
- Cumplimiento más sencillo con futuras certificaciones (eIDAS-like, ISO 27001).
- Auditoría del KMS es independiente de la del servicio.

**Negativas**
- Latencia adicional por hop al KMS (mitigado por colocación de Vault en cluster local; en prod usar AWS KMS regional).
- Costo: AWS KMS cobra por operación (~$0.03 por 10k ops). A 1k req/s sostenidos, ~$3/hora. Aceptable.
- El KMS es ahora un punto crítico — necesita HA, monitoring agresivo, runbook de incidentes.

## Alternativas descartadas

- **Claves en BD cifradas con secret aplicativo**: el "secret aplicativo" termina en variables de entorno y memoria del proceso → vector de exfiltración trivial.
- **Custodia client-side (clave en el dispositivo del ciudadano)**: hace e-firmas legales imposibles si pierde el dispositivo, fricciona la UX a niveles que matan la adopción. Estonia lo evita usando server-side signing tras autenticación fuerte.
- **HSM físico en cada realm desde día 1**: prohibitivo en costo para arrancar; introducir cuando un realm lo demande regulatoriamente.

## Referencias

- Vault Transit secrets engine.
- AWS KMS asymmetric keys (Ed25519, ECDSA).
- eIDAS Qualified Electronic Signature requirements (referencia conceptual).
