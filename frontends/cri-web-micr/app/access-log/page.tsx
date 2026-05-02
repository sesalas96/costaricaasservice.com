"use client";

import { useEffect, useState } from "react";
import { api, formatDateTime, type AccessLog } from "@/lib/api";
import { useActiveCedula } from "@/lib/session";

export default function AccessLogPage() {
  const { cedula, loading: sessionLoading } = useActiveCedula();
  const [log, setLog] = useState<AccessLog | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!cedula) return;
    setLoading(true);
    api
      .accessLog(cedula)
      .then((d) => {
        setLog(d);
        setError(null);
      })
      .catch((e: Error) => setError(e.message))
      .finally(() => setLoading(false));
  }, [cedula]);

  if (sessionLoading || loading) return <div className="text-slate-500">Cargando…</div>;
  if (error)
    return (
      <div className="bg-red-50 border border-red-200 text-red-800 rounded-lg p-4">
        <strong>Error: </strong>{error}
      </div>
    );

  return (
    <div className="space-y-6">
      <header>
        <h1 className="text-2xl font-semibold">Mi bitácora</h1>
        <p className="text-sm text-slate-500">
          Accesos del Estado a sus datos personales — quién, cuándo y por qué.
          Cada entry está encadenada criptográficamente; modificar una línea pasada
          rompería todas las posteriores.
        </p>
      </header>

      <section className="bg-white border border-slate-200 rounded-xl shadow-sm overflow-hidden">
        <div className="px-5 py-3 border-b border-slate-200 flex items-baseline justify-between">
          <div>
            <span className="text-2xl font-bold">{log?.count ?? 0}</span>{" "}
            <span className="text-sm text-slate-500">acceso{log?.count === 1 ? "" : "s"} registrado{log?.count === 1 ? "" : "s"}</span>
          </div>
          <code className="text-xs text-slate-500 font-mono">cédula {cedula}</code>
        </div>

        {(log?.entries?.length ?? 0) === 0 ? (
          <div className="p-8 text-center text-slate-500">No hay accesos registrados aún.</div>
        ) : (
          <ul className="divide-y divide-slate-100">
            {log!.entries.map((e) => (
              <li key={e.id} className="px-5 py-4 hover:bg-slate-50">
                <div className="flex flex-wrap items-baseline justify-between gap-2">
                  <div className="text-sm flex items-center gap-2 flex-wrap">
                    <code className="font-mono font-semibold">{e.requesterMember}</code>
                    <SectorBadge member={e.requesterMember} />
                    <span>consultó</span>
                    <code className="font-mono">{e.service}/{e.version}</code>
                    <span>en</span>
                    <code className="font-mono">{e.targetMember}</code>
                  </div>
                  <div className="text-xs text-slate-500">{formatDateTime(e.ts)}</div>
                </div>
                <div className="mt-1 text-xs text-slate-600">
                  Propósito: <strong>{e.purpose}</strong>
                </div>
                <div className="mt-2 grid sm:grid-cols-2 gap-x-6 gap-y-1 text-[11px] font-mono text-slate-400">
                  <div>id: {e.id}</div>
                  <div>req: {e.requestId}</div>
                  <div>hash: {e.entryHash.slice(0, 24)}…</div>
                  <div>prev: {e.prevHash.slice(0, 24)}…</div>
                </div>
              </li>
            ))}
          </ul>
        )}
      </section>

      <p className="text-xs text-slate-500">
        El audit log es <strong>append-only</strong> y verificado por hash chain
        (SHA-256). Cualquier intento de borrado o modificación retroactiva sería
        detectable en menos de un segundo.
      </p>
    </div>
  );
}

// Members del sector privado: en producción esto se resuelve consultando el
// catálogo del cri-svc-interop-hub. Para MVP demo se hardcodea acá.
const PRIVATE_MEMBERS = new Set(["bcr", "bncr", "scotiabank", "fischel", "cima", "clinica-biblica"]);

function SectorBadge({ member }: { member: string }) {
  const isPrivate = PRIVATE_MEMBERS.has(member);
  if (isPrivate) {
    return (
      <span className="text-[10px] font-medium px-1.5 py-0.5 rounded bg-amber-100 text-amber-800">
        privado
      </span>
    );
  }
  return (
    <span className="text-[10px] font-medium px-1.5 py-0.5 rounded bg-emerald-100 text-emerald-800">
      Estado
    </span>
  );
}
