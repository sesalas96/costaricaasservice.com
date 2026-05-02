# ADR-0002 — Audit log inmutable: hash chain + Merkle por epoch

- **Estado:** Aceptado
- **Fecha:** 2026-05-02
- **Decisores:** equipo saascr

## Contexto

Una de las cuatro garantías estonias del producto es **trazabilidad ciudadana**: el ciudadano puede ver quién consultó sus datos, cuándo y por qué. Esto solo tiene valor si el log es **detectablemente inmutable** — un breach al servicio de audit no debe permitir reescribir registros pasados sin dejar huella.

Estonia usa KSI Blockchain (tecnología propietaria de Guardtime). Para un SaaS comercial necesitamos algo equivalente que no tenga vendor lock-in y sea verificable independientemente.

## Decisión

Cada transacción inter-member produce un `AuditEntry` que se escribe en `cri-svc-interop-audit` con la siguiente estructura:

```text
AuditEntry {
  id            ULID            -- ordenamiento total, monotónico
  realm         text            -- multi-tenancy
  ts            timestamptz     -- server time, UTC
  requester     {member, svc}   -- quién originó la consulta
  target        {member, svc}   -- a quién se le consultó
  service       string          -- servicio invocado en el catálogo
  citizen_id    nullable        -- si la consulta es sobre un ciudadano
  purpose       text            -- propósito declarado (ej: prefill_tax_return)
  request_hash  bytes           -- SHA-256 del request canonicalizado (sin payloads sensibles)
  prev_hash     bytes           -- hash del entry anterior (chain)
  entry_hash    bytes           -- SHA-256 de (prev_hash || canonical(this))
}
```

### Hash chain

- `prev_hash` del entry N+1 = `entry_hash` del entry N. Cualquier modificación retroactiva **rompe la cadena** de todos los registros posteriores.
- El primer entry de cada realm tiene `prev_hash = SHA-256("genesis::<realm>")`.
- Verificador automático corre en CI y como cron en prod (`cri-svc-interop-audit`).

### Merkle root por epoch

- Cada **epoch** (cada hora o cada 10k entries, lo que ocurra primero), se calcula un Merkle root sobre todos los `entry_hash` del epoch.
- El root se publica en:
  - **Siempre:** S3 con Object Lock (modo Compliance, retention 10 años).
  - **Premium tier (opcional por realm):** anclaje en una cadena pública (Bitcoin OP_RETURN o Ethereum L2) para garantía third-party.

### Visibilidad ciudadana

- `cri-svc-audit` lee de `cri-svc-interop-audit` y proyecta una vista por ciudadano: "consultas que tocaron tus datos en los últimos N días, con propósito y member".
- El ciudadano ve `entry_hash` y puede verificar la pertenencia al Merkle root vía un endpoint público de prueba.

## Consecuencias

**Positivas**
- Tampering retroactivo es **detectable en O(n) tiempo** desde el primer entry modificado hasta el final.
- Independiente de la honestidad del operador del servicio (si el operador modifica un registro, el verificador y los anchors externos lo delatan).
- Sin vendor lock-in (sin KSI / Guardtime).

**Negativas**
- Append-only es estricto — no se pueden corregir errores en datos pasados sin un mecanismo explícito de "compensating entry" (que también queda en la cadena).
- Verificación full-chain es O(n); aceptable hasta ~100M entries por realm; mitigado por Merkle por epoch.
- Hay que cuidar el **clock**: usar timestamps del servidor con fuente NTP estricta y un campo separado para "purpose-stated time" si aplica.

## Implementación de referencia

- `libs/cri-lib-crypto/hashchain` — helpers para canonicalize, hash, link.
- `cri-svc-interop-audit` — escritor, verificador, generador de Merkle roots, publicador a S3 Object Lock.
- Verificador en CI: `make test-audit-chain` debe correr en cada PR contra el log de staging.

## Alternativas descartadas

- **KSI / Guardtime**: vendor lock-in, costo, opacidad.
- **Solo append-only DB sin chain**: protege contra borrados accidentales, no contra un atacante con acceso a la DB.
- **Blockchain interno completo**: overkill, complejo, introduce un nuevo modelo de consenso para resolver un problema con una hash chain ya alcanza.

## Referencias

- Estonia KSI Blockchain (concepto).
- Certificate Transparency Merkle tree (RFC 6962).
