"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

const DEMO_CITIZENS = [
  { cedula: "1-1234-5678", name: "Sebastián Rojas Mora", hint: "ingreso 18M" },
  { cedula: "1-9876-5432", name: "María Castro Fernández", hint: "ingreso 24.5M, con dependientes" },
  { cedula: "2-0001-0001", name: "Luis Fernández Solís", hint: "ingreso 12M" },
];

export default function Home() {
  const router = useRouter();
  const [activeCedula, setActiveCedula] = useState<string | null>(null);

  useEffect(() => {
    if (typeof window !== "undefined") {
      setActiveCedula(window.localStorage.getItem("micr_cedula"));
    }
  }, []);

  const select = (cedula: string) => {
    window.localStorage.setItem("micr_cedula", cedula);
    router.push("/dashboard");
  };

  return (
    <div className="space-y-8">
      <section className="bg-white border border-slate-200 rounded-xl p-8 shadow-sm">
        <h1 className="text-2xl font-semibold mb-2">Bienvenido a su Carpeta Digital</h1>
        <p className="text-slate-600 max-w-2xl">
          Acá puede ver sus datos, su declaración pre-llenada de Hacienda, y la bitácora de
          accesos del Estado a su información personal — quién consultó sus datos, cuándo y
          con qué propósito.
        </p>
        {activeCedula && (
          <p className="mt-4 text-sm">
            Sesión activa como <code className="font-mono bg-slate-100 px-1.5 py-0.5 rounded">{activeCedula}</code>.{" "}
            <button
              onClick={() => router.push("/dashboard")}
              className="text-blue-700 hover:underline"
            >
              Ir al escritorio →
            </button>
          </p>
        )}
      </section>

      <section>
        <h2 className="text-lg font-semibold mb-4">Demo: seleccione un ciudadano</h2>
        <p className="text-sm text-slate-500 mb-4">
          Esta carpeta opera sobre el realm <code className="font-mono">demo</code>. En producción
          se accede con identidad digital (e-ID).
        </p>
        <div className="grid sm:grid-cols-3 gap-4">
          {DEMO_CITIZENS.map((c) => (
            <button
              key={c.cedula}
              onClick={() => select(c.cedula)}
              className="text-left bg-white border border-slate-200 rounded-xl p-5 shadow-sm hover:border-blue-500 hover:shadow-md transition"
            >
              <div className="font-semibold">{c.name}</div>
              <div className="text-xs text-slate-500 font-mono mt-1">{c.cedula}</div>
              <div className="text-xs text-slate-400 mt-3">{c.hint}</div>
            </button>
          ))}
        </div>
      </section>
    </div>
  );
}
