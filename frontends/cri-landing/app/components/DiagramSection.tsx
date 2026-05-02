"use client";

import { useState } from "react";
import { useT } from "../i18n/LanguageProvider";
import EcosystemDiagram from "./diagrams/EcosystemDiagram";
import InteropFlowDiagram from "./diagrams/InteropFlowDiagram";
import OnceOnlyDiagram from "./diagrams/OnceOnlyDiagram";
import IdentityDiagram from "./diagrams/IdentityDiagram";
import type { Audience } from "./diagrams/types";

type ModuleKey = "ecosistema" | "interop" | "onceOnly" | "identity";

const MODULE_ACCENT: Record<ModuleKey, string> = {
  ecosistema: "var(--color-cr-blue-bright)",
  interop: "var(--color-area-interop)",
  onceOnly: "var(--color-area-iduc)",
  identity: "var(--color-cr-red-bright)",
};

function AudienceSwitch({
  value,
  onChange,
  label,
}: {
  value: Audience;
  onChange: (a: Audience) => void;
  label: { citizen: string; technical: string };
}) {
  const opts: Audience[] = ["citizen", "technical"];
  return (
    <div
      role="radiogroup"
      aria-label="Audience"
      className="relative inline-flex rounded-full border border-[var(--color-border)] bg-[var(--color-surface)] p-1"
    >
      {/* Pastilla deslizante */}
      <span
        aria-hidden
        className="absolute top-1 bottom-1 w-1/2 rounded-full transition-transform duration-300 ease-out"
        style={{
          background:
            "linear-gradient(120deg, var(--color-cr-blue) 0%, var(--color-cr-blue-bright) 100%)",
          boxShadow: "0 4px 16px -4px rgba(29,79,196,0.55)",
          transform:
            value === "citizen" ? "translateX(0%)" : "translateX(100%)",
          left: 4,
          right: 4,
        }}
      />
      {opts.map((a) => {
        const active = a === value;
        return (
          <button
            key={a}
            role="radio"
            aria-checked={active}
            onClick={() => onChange(a)}
            className="relative z-10 px-5 py-2 text-sm font-medium transition-colors"
            style={{
              color: active ? "#fff" : "var(--color-muted)",
            }}
          >
            {a === "citizen" ? label.citizen : label.technical}
          </button>
        );
      })}
    </div>
  );
}

function ModuleNavItem({
  active,
  accent,
  tag,
  title,
  sub,
  onClick,
}: {
  active: boolean;
  accent: string;
  tag: string;
  title: string;
  sub: string;
  onClick: () => void;
}) {
  return (
    <button
      role="tab"
      aria-selected={active}
      onClick={onClick}
      className="group relative w-full overflow-hidden rounded-lg border p-4 text-left transition-colors"
      style={{
        background: active ? "var(--color-surface-2)" : "var(--color-surface)",
        borderColor: active ? accent : "var(--color-border)",
      }}
    >
      <div
        aria-hidden
        className="absolute left-0 top-0 h-full w-[3px] transition-opacity"
        style={{
          background: accent,
          opacity: active ? 1 : 0,
        }}
      />
      <div className="flex items-center gap-2">
        <span
          className="h-1.5 w-1.5 rounded-full"
          style={{ background: accent }}
        />
        <span
          className="font-mono text-[10px] uppercase tracking-widest"
          style={{ color: active ? accent : "var(--color-muted)" }}
        >
          {tag}
        </span>
      </div>
      <div className="mt-1 text-[14px] font-semibold text-[var(--color-text)]">
        {title}
      </div>
      <div className="mt-0.5 text-[12px] leading-snug text-[var(--color-muted)]">
        {sub}
      </div>
    </button>
  );
}

export default function DiagramSection() {
  const { t } = useT();
  const [active, setActive] = useState<ModuleKey>("ecosistema");
  const [audience, setAudience] = useState<Audience>("citizen");

  const modules: { key: ModuleKey; tag: string; title: string; sub: string }[] = [
    {
      key: "ecosistema",
      tag: t.diagrams.modules.ecosistema.tag,
      title: t.diagrams.modules.ecosistema.title,
      sub: t.diagrams.modules.ecosistema.sub,
    },
    {
      key: "interop",
      tag: t.diagrams.modules.interop.tag,
      title: t.diagrams.modules.interop.title,
      sub: t.diagrams.modules.interop.sub,
    },
    {
      key: "onceOnly",
      tag: t.diagrams.modules.onceOnly.tag,
      title: t.diagrams.modules.onceOnly.title,
      sub: t.diagrams.modules.onceOnly.sub,
    },
    {
      key: "identity",
      tag: t.diagrams.modules.identity.tag,
      title: t.diagrams.modules.identity.title,
      sub: t.diagrams.modules.identity.sub,
    },
  ];

  const activeModule = modules.find((m) => m.key === active)!;
  const accent = MODULE_ACCENT[active];
  const sub = t.diagrams.subs[active][audience];

  // key fuerza re-mount del diagrama para que fitView corra al cambiar
  const diagramKey = `${active}-${audience}`;

  return (
    <section id="diagramas" className="py-24">
      <div className="mx-auto max-w-6xl px-6">
        {/* ── Header con switch ───────────────────────────────── */}
        <div className="mb-10 flex flex-col gap-6 lg:flex-row lg:items-end lg:justify-between">
          <div className="flex flex-col gap-3">
            <div className="text-xs font-mono uppercase tracking-widest text-[var(--color-muted)]">
              {t.diagrams.eyebrow}
            </div>
            <h2 className="max-w-2xl text-3xl font-semibold tracking-tight md:text-4xl">
              {t.diagrams.heading}
            </h2>
          </div>
          <div className="flex flex-col items-start gap-2 lg:items-end">
            <AudienceSwitch
              value={audience}
              onChange={setAudience}
              label={{
                citizen: t.diagrams.audience.citizen,
                technical: t.diagrams.audience.technical,
              }}
            />
            <p className="font-mono text-[11px] text-[var(--color-muted)] lg:text-right">
              {t.diagrams.audience.hint}
            </p>
          </div>
        </div>

        {/* ── Hub modular ─────────────────────────────────────── */}
        <div className="grid grid-cols-1 gap-4 lg:grid-cols-12">
          {/* Sidebar de módulos */}
          <aside
            role="tablist"
            aria-orientation="vertical"
            className="lg:col-span-3"
          >
            <div className="mb-3 px-1 font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted)]">
              {t.diagrams.modulesHeading}
            </div>
            <div className="grid grid-cols-2 gap-2 sm:grid-cols-4 lg:grid-cols-1">
              {modules.map((m) => (
                <ModuleNavItem
                  key={m.key}
                  active={m.key === active}
                  accent={MODULE_ACCENT[m.key]}
                  tag={m.tag}
                  title={m.title}
                  sub={m.sub}
                  onClick={() => setActive(m.key)}
                />
              ))}
            </div>
          </aside>

          {/* Área principal */}
          <div className="lg:col-span-9">
            <div
              className="overflow-hidden rounded-xl border bg-[var(--color-surface)]"
              style={{ borderColor: "var(--color-border)" }}
            >
              {/* Encabezado del módulo activo */}
              <div
                className="flex flex-col gap-2 border-b border-[var(--color-border)] p-5"
                style={{
                  background:
                    "linear-gradient(180deg, rgba(15,18,24,0.6) 0%, transparent 100%)",
                }}
              >
                <div className="flex items-center gap-3">
                  <span
                    className="h-2 w-2 rounded-full"
                    style={{ background: accent }}
                  />
                  <span
                    className="font-mono text-[10px] uppercase tracking-widest"
                    style={{ color: accent }}
                  >
                    {activeModule.tag}
                  </span>
                  <span
                    className="rounded-full border border-[var(--color-border)] bg-[var(--color-bg)] px-2 py-0.5 font-mono text-[9px] uppercase tracking-widest text-[var(--color-muted)]"
                  >
                    {audience === "citizen"
                      ? t.diagrams.audience.citizen
                      : t.diagrams.audience.technical}
                    {" · "}
                    {t.diagrams.audience.activeBadge}
                  </span>
                </div>
                <h3 className="text-lg font-semibold">{activeModule.title}</h3>
                <p className="text-[14px] leading-relaxed text-[var(--color-muted)]">
                  {sub}
                </p>
              </div>

              {/* Canvas */}
              <div className="h-[620px] w-full">
                {active === "ecosistema" && (
                  <EcosystemDiagram key={diagramKey} audience={audience} />
                )}
                {active === "interop" && (
                  <InteropFlowDiagram key={diagramKey} audience={audience} />
                )}
                {active === "onceOnly" && (
                  <OnceOnlyDiagram key={diagramKey} audience={audience} />
                )}
                {active === "identity" && (
                  <IdentityDiagram key={diagramKey} audience={audience} />
                )}
              </div>
            </div>

            <div className="mt-3 flex items-center gap-3 font-mono text-[11px] text-[var(--color-muted)]">
              <span>tip:</span>
              <span>{t.diagrams.tip}</span>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
