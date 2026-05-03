"use client";

import { useState } from "react";
import { useT } from "../i18n/LanguageProvider";

type View = "business" | "technical";

const BUSINESS_ACCENTS = [
  "var(--color-cr-blue-bright)",
  "var(--color-area-iduc)",
  "var(--color-cr-red-bright)",
  "var(--color-area-platform)",
];

const TECHNICAL_ACCENTS = [
  "var(--color-area-gateway)",
  "var(--color-area-iduc)",
  "var(--color-area-interop)",
  "var(--color-area-members)",
  "var(--color-area-platform)",
  "var(--color-area-frontend)",
];

function ViewSwitch({
  value,
  onChange,
  label,
}: {
  value: View;
  onChange: (v: View) => void;
  label: { business: string; technical: string; aria: string };
}) {
  const opts: View[] = ["business", "technical"];
  return (
    <div
      role="radiogroup"
      aria-label={label.aria}
      className="relative inline-flex rounded-full border border-[var(--color-border)] bg-[var(--color-surface)] p-1"
    >
      <span
        aria-hidden
        className="absolute top-1 bottom-1 w-1/2 rounded-full transition-transform duration-300 ease-out"
        style={{
          background:
            "linear-gradient(120deg, var(--color-cr-blue) 0%, var(--color-cr-blue-bright) 100%)",
          boxShadow: "0 4px 16px -4px rgba(29,79,196,0.55)",
          transform:
            value === "business" ? "translateX(0%)" : "translateX(100%)",
          left: 4,
          right: 4,
        }}
      />
      {opts.map((v) => {
        const active = v === value;
        return (
          <button
            key={v}
            role="radio"
            aria-checked={active}
            onClick={() => onChange(v)}
            className="relative z-10 px-5 py-2 text-sm font-medium transition-colors"
            style={{ color: active ? "#fff" : "var(--color-muted)" }}
          >
            {v === "business" ? label.business : label.technical}
          </button>
        );
      })}
    </div>
  );
}

export default function Deliverables() {
  const { t } = useT();
  const d = t.deliverables;
  const [view, setView] = useState<View>("business");

  const data = view === "business" ? d.business : d.technical;
  const accents = view === "business" ? BUSINESS_ACCENTS : TECHNICAL_ACCENTS;
  const gridCols =
    view === "business"
      ? "sm:grid-cols-2 lg:grid-cols-2"
      : "sm:grid-cols-2 lg:grid-cols-3";

  return (
    <section
      id="entregables"
      className="border-b border-t border-[var(--color-border)] py-16 md:py-24"
    >
      <div className="mx-auto max-w-6xl px-5 sm:px-6">
        {/* ── Header con switch ───────────────────────────────── */}
        <div className="mb-8 flex flex-col gap-6 md:mb-10 lg:flex-row lg:items-end lg:justify-between">
          <div className="flex flex-col gap-3">
            <div className="text-[10px] font-mono uppercase tracking-widest text-[var(--color-muted)] sm:text-xs">
              {d.eyebrow}
            </div>
            <h2 className="max-w-2xl text-2xl font-semibold tracking-tight sm:text-3xl md:text-4xl">
              {d.heading}
            </h2>
          </div>
          <div className="flex flex-col items-start gap-2 lg:items-end">
            <ViewSwitch
              value={view}
              onChange={setView}
              label={{
                business: d.view.business,
                technical: d.view.technical,
                aria: d.view.aria,
              }}
            />
            <p className="font-mono text-[11px] text-[var(--color-muted)] lg:text-right">
              {d.view.hint}
            </p>
          </div>
        </div>

        {/* ── Lead + intro ────────────────────────────────────── */}
        <div className="mb-8 grid grid-cols-1 gap-6 md:mb-10 lg:grid-cols-12 lg:gap-10">
          <p className="text-base leading-relaxed text-[var(--color-muted)] sm:text-lg lg:col-span-6">
            {d.lead}
          </p>
          <p className="text-base leading-relaxed text-[var(--color-muted)] lg:col-span-6">
            {data.intro}
          </p>
        </div>

        {/* ── Grid de entregables ─────────────────────────────── */}
        <div
          className={`grid grid-cols-1 gap-px overflow-hidden rounded-xl border border-[var(--color-border)] bg-[var(--color-border)] ${gridCols}`}
        >
          {data.cards.map((card, i) => {
            const accent = accents[i % accents.length];
            const num = String(i + 1).padStart(2, "0");
            return (
              <article
                key={card.title}
                className="group relative flex h-full flex-col gap-4 bg-[var(--color-surface)] p-5 transition-colors hover:bg-[var(--color-surface-2)] sm:p-7"
              >
                <div
                  aria-hidden
                  className="absolute inset-x-0 top-0 h-[3px]"
                  style={{ background: accent }}
                />
                <div className="flex items-center justify-between">
                  <span
                    className="font-mono text-[10px] uppercase tracking-widest"
                    style={{ color: accent }}
                  >
                    {num} · {card.tag}
                  </span>
                  <span
                    className="h-1.5 w-1.5 rounded-full"
                    style={{ background: accent }}
                  />
                </div>
                <div className="flex flex-col gap-1">
                  <h3 className="text-lg font-semibold leading-snug">
                    {card.title}
                  </h3>
                  <p className="text-[13px] leading-relaxed text-[var(--color-muted)]">
                    {card.sub}
                  </p>
                </div>
                <ul className="mt-1 flex flex-col gap-2">
                  {card.items.map((item) => (
                    <li
                      key={item}
                      className="flex gap-2 text-[13.5px] leading-relaxed text-[var(--color-text)]"
                    >
                      <span
                        aria-hidden
                        className="mt-[7px] h-1 w-1 shrink-0 rounded-full"
                        style={{ background: accent }}
                      />
                      <span>{item}</span>
                    </li>
                  ))}
                </ul>
                <div className="mt-auto border-t border-[var(--color-border)] pt-3 font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted)]">
                  {view === "business"
                    ? d.view.business
                    : d.view.technical}
                  {" · "}
                  {d.view.activeBadge}
                </div>
              </article>
            );
          })}
        </div>
      </div>
    </section>
  );
}
