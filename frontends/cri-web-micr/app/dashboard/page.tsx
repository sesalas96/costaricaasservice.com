"use client";

import { useEffect, useState } from "react";
import Link from "next/link";

import { api, formatCRC, formatDateTime, type Dashboard } from "@/lib/api";
import { useActiveCedula } from "@/lib/session";

export default function DashboardPage() {
  const { cedula, loading: sessionLoading } = useActiveCedula();
  const [data, setData] = useState<Dashboard | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!cedula) return;
    setLoading(true);
    api
      .dashboard(cedula, 2025)
      .then((d) => {
        setData(d);
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
  if (!data) return null;

  const tax = data.taxPrefilled?.data;
  const health = data.healthProfile?.data;
  const log = data.accessLog?.data;
  const lastEntry = log?.entries?.[0];

  return (
    <div className="space-y-6">
      <header className="flex items-baseline justify-between">
        <div>
          <h1 className="text-2xl font-semibold">Mi escritorio</h1>
          <p className="text-sm text-slate-500">
            Año fiscal {data.year} · Cédula <code className="font-mono">{data.cedula}</code>
          </p>
        </div>
      </header>

      <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-5">
        {/* Card 1: Mis datos (vienen de Registro Civil vía interop) */}
        <article className="bg-white border border-slate-200 rounded-xl p-5 shadow-sm">
          <div className="flex items-center justify-between mb-3">
            <h2 className="font-semibold">Mis datos</h2>
            <Badge color="emerald">Registro Civil</Badge>
          </div>
          {tax?.person ? (
            <dl className="space-y-1.5 text-sm">
              <Row label="Nombre" value={tax.person.fullName} />
              <Row label="Cédula" value={<code className="font-mono">{tax.person.cedula}</code>} />
              <Row label="Domicilio" value={tax.person.address} />
              <Row label="Email" value={tax.person.email} />
            </dl>
          ) : (
            <p className="text-slate-500 text-sm">Sin datos disponibles.</p>
          )}
          <p className="text-xs text-slate-400 mt-4 leading-relaxed">
            Estos datos viajan firmados desde el Registro Civil. Hacienda los consulta
            cada vez sin guardarlos (<em>once-only</em>).
          </p>
        </article>

        {/* Card 2: Mi declaración (datos de Hacienda + composite) */}
        <article className="bg-white border border-slate-200 rounded-xl p-5 shadow-sm">
          <div className="flex items-center justify-between mb-3">
            <h2 className="font-semibold">Mi declaración 2025</h2>
            <Badge color="blue">Hacienda</Badge>
          </div>
          {tax ? (
            <dl className="space-y-1.5 text-sm">
              <Row label="Ingreso bruto" value={formatCRC(tax.grossIncome)} />
              <Row label="Retenido" value={formatCRC(tax.withheldTax)} />
              <Row label="Deducciones" value={formatCRC(tax.deductions)} />
              <Row label="A pagar" value={<strong>{formatCRC(tax.estimatedDue)}</strong>} />
              <Row label="Dependientes" value={tax.hasDependents ? "Sí" : "No"} />
            </dl>
          ) : (
            <p className="text-slate-500 text-sm">Sin declaración.</p>
          )}
          <Link href="/tax" className="text-xs text-blue-700 hover:underline mt-4 inline-block">
            Ver detalle →
          </Link>
        </article>

        {/* Card 3: Mi salud (vienen de CCSS, que internamente consulta Registro Civil) */}
        <article className="bg-white border border-slate-200 rounded-xl p-5 shadow-sm">
          <div className="flex items-center justify-between mb-3">
            <h2 className="font-semibold">Mi salud</h2>
            <Badge color="rose">CCSS</Badge>
          </div>
          {health ? (
            <div className="space-y-3 text-sm">
              <div>
                <div className="text-3xl font-bold">{health.activePrescriptions.length}</div>
                <div className="text-xs text-slate-500">recetas activas</div>
              </div>
              {health.activePrescriptions.slice(0, 2).map((rx) => (
                <div key={rx.id} className="bg-slate-50 rounded-md p-2">
                  <div className="font-medium text-xs">{rx.drug}</div>
                  <div className="text-[11px] text-slate-500">{rx.dosage}</div>
                </div>
              ))}
              {health.nextAppointment && (
                <div className="border-t pt-2">
                  <div className="text-xs text-slate-500">Próxima cita</div>
                  <div className="font-medium text-xs">{health.nextAppointment.specialty}</div>
                  <div className="text-[11px] text-slate-500">
                    {health.nextAppointment.when} · {health.nextAppointment.hospital}
                  </div>
                </div>
              )}
            </div>
          ) : (
            <p className="text-slate-500 text-sm">Sin información de salud.</p>
          )}
        </article>

        {/* Card 4: Mi bitácora */}
        <article className="bg-white border border-slate-200 rounded-xl p-5 shadow-sm">
          <div className="flex items-center justify-between mb-3">
            <h2 className="font-semibold">Mi bitácora</h2>
            <Badge color="purple">Audit chain</Badge>
          </div>
          <div className="text-3xl font-bold mb-1">{log?.count ?? 0}</div>
          <p className="text-sm text-slate-500 mb-4">accesos a sus datos</p>
          {lastEntry && (
            <div className="bg-slate-50 rounded-lg p-3 text-xs">
              <div className="font-semibold">Último acceso</div>
              <div className="mt-1 text-slate-600">
                <code className="font-mono">{lastEntry.requesterMember}</code> consultó{" "}
                <code className="font-mono">{lastEntry.service}</code> en{" "}
                <code className="font-mono">{lastEntry.targetMember}</code>
              </div>
              <div className="mt-1 text-slate-500">{formatDateTime(lastEntry.ts)}</div>
              <div className="mt-1 text-slate-500">propósito: {lastEntry.purpose}</div>
            </div>
          )}
          <Link href="/access-log" className="text-xs text-blue-700 hover:underline mt-4 inline-block">
            Ver bitácora completa →
          </Link>
        </article>
      </div>

      {data.upstreamErrors && data.upstreamErrors.length > 0 && (
        <div className="bg-amber-50 border border-amber-200 text-amber-800 rounded-lg p-3 text-xs">
          <strong>Avisos: </strong>
          {data.upstreamErrors.join("; ")}
        </div>
      )}
    </div>
  );
}

function Row({ label, value }: { label: string; value: React.ReactNode }) {
  return (
    <div className="flex justify-between gap-4">
      <dt className="text-slate-500">{label}</dt>
      <dd className="text-slate-900 text-right">{value}</dd>
    </div>
  );
}

function Badge({ color, children }: { color: "emerald" | "blue" | "purple" | "rose"; children: React.ReactNode }) {
  const colors: Record<string, string> = {
    emerald: "bg-emerald-100 text-emerald-800",
    blue: "bg-blue-100 text-blue-800",
    purple: "bg-purple-100 text-purple-800",
    rose: "bg-rose-100 text-rose-800",
  };
  return <span className={`text-xs font-medium px-2 py-0.5 rounded ${colors[color]}`}>{children}</span>;
}
