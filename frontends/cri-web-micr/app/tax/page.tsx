"use client";

import { useEffect, useState } from "react";
import { api, formatCRC, type TaxPrefilled } from "@/lib/api";
import { useActiveCedula } from "@/lib/session";

export default function TaxPage() {
  const { cedula, loading: sessionLoading } = useActiveCedula();
  const [tax, setTax] = useState<TaxPrefilled | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [year] = useState(2025);

  useEffect(() => {
    if (!cedula) return;
    setLoading(true);
    api
      .taxPrefilled(cedula, year)
      .then((d) => {
        setTax(d);
        setError(null);
      })
      .catch((e: Error) => setError(e.message))
      .finally(() => setLoading(false));
  }, [cedula, year]);

  if (sessionLoading || loading) return <div className="text-slate-500">Cargando…</div>;
  if (error)
    return (
      <div className="bg-red-50 border border-red-200 text-red-800 rounded-lg p-4">
        <strong>Error: </strong>{error}
      </div>
    );
  if (!tax) return null;

  return (
    <div className="space-y-6">
      <header className="flex items-baseline justify-between">
        <div>
          <h1 className="text-2xl font-semibold">Mi declaración {year}</h1>
          <p className="text-sm text-slate-500">
            Pre-llenada por Hacienda. Sus datos personales fueron consultados al
            Registro Civil vía <strong>Conecta CR</strong> con propósito declarado
            (<code className="font-mono">prefill_tax_return</code>).
          </p>
        </div>
        <button
          disabled
          className="bg-blue-600 text-white px-4 py-2 rounded-lg text-sm font-medium opacity-60 cursor-not-allowed"
          title="MVP: la firma se conectará en el siguiente sprint"
        >
          Firmar y presentar
        </button>
      </header>

      <section className="bg-white border border-slate-200 rounded-xl shadow-sm">
        <div className="p-5 border-b border-slate-100">
          <h2 className="font-semibold mb-2">Identificación</h2>
          <dl className="grid sm:grid-cols-2 gap-x-8 gap-y-1.5 text-sm">
            <Row label="Nombre" value={tax.person.fullName} />
            <Row label="Cédula" value={<code className="font-mono">{tax.person.cedula}</code>} />
            <Row label="Domicilio" value={tax.person.address} />
            <Row label="Email" value={tax.person.email} />
          </dl>
        </div>

        <div className="p-5">
          <h2 className="font-semibold mb-2">Detalle del año {year}</h2>
          <dl className="grid sm:grid-cols-2 gap-x-8 gap-y-1.5 text-sm">
            <Row label="Ingreso bruto" value={formatCRC(tax.grossIncome)} />
            <Row label="Impuesto retenido" value={formatCRC(tax.withheldTax)} />
            <Row label="Deducciones" value={formatCRC(tax.deductions)} />
            <Row label="Dependientes" value={tax.hasDependents ? "Sí" : "No"} />
          </dl>
          <div className="mt-5 pt-4 border-t border-slate-100 flex items-baseline justify-between">
            <span className="text-sm text-slate-500">Estimado a pagar</span>
            <span className="text-2xl font-bold">{formatCRC(tax.estimatedDue)}</span>
          </div>
        </div>
      </section>

      <section className="bg-slate-50 border border-slate-200 rounded-xl p-5">
        <h3 className="font-semibold text-sm mb-2">Trazabilidad del pre-llenado</h3>
        <p className="text-xs text-slate-600 leading-relaxed">
          Para construir esta declaración, Hacienda consultó al Registro Civil sus
          datos básicos (nombre, domicilio). Esta consulta fue <strong>firmada
          digitalmente</strong> por el security-server de Hacienda y verificada por
          el security-server de Registro Civil antes de devolver los datos. Quedó
          registrada en su <a href="/access-log" className="text-blue-700 underline">bitácora</a>.
        </p>
        <code className="block mt-3 text-[11px] font-mono text-slate-500 break-all">
          {tax._onceOnlyTrace}
        </code>
      </section>
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
