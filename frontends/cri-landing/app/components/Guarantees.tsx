"use client";

import Link from "next/link";
import { useT } from "../i18n/LanguageProvider";

const GUARANTEE_ACCENTS = [
  "var(--color-area-iduc)",
  "var(--color-cr-blue-bright)",
  "var(--color-area-platform)",
  "var(--color-cr-red-bright)",
];

const AUDIENCE_ACCENTS = [
  "var(--color-cr-blue-bright)",
  "var(--color-cr-cream)",
  "var(--color-cr-red-bright)",
];

const CASE_ACCENTS = [
  "var(--color-area-members)",
  "var(--color-area-iduc)",
  "var(--color-area-platform)",
  "var(--color-area-interop)",
  "var(--color-cr-blue-bright)",
  "var(--color-area-gateway)",
  "var(--color-area-members)",
  "var(--color-area-platform)",
  "var(--color-cr-red-bright)",
];

const FEATURED_CASE_INDEXES = [0, 2, 3, 4, 8];

export default function Guarantees() {
  const { t } = useT();
  const g = t.guarantees;

  return (
    <section
      id="garantias"
      className="border-b border-[var(--color-border)] py-24"
    >
      <div className="mx-auto max-w-6xl px-6">
        {/* ── Cabezal con descripción ──────────────────────── */}
        <div className="mb-16 grid grid-cols-1 gap-10 lg:grid-cols-12 lg:gap-16">
          <div className="lg:col-span-5">
            <div className="text-xs font-mono uppercase tracking-widest text-[var(--color-muted)]">
              {g.eyebrow}
            </div>
            <h2 className="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">
              {g.heading}
            </h2>
          </div>
          <div className="lg:col-span-7">
            <p className="text-lg leading-relaxed text-[var(--color-muted)]">
              {g.lead}
            </p>
          </div>
        </div>

        {/* ── 3 audiencias ─────────────────────────────────── */}
        <div className="mb-20">
          <div className="mb-6 text-xs font-mono uppercase tracking-widest text-[var(--color-muted)]">
            {g.audiencesHeading}
          </div>
          <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
            {g.audiences.map((a, i) => {
              const accent = AUDIENCE_ACCENTS[i % AUDIENCE_ACCENTS.length];
              return (
                <div
                  key={a.title}
                  className="relative overflow-hidden rounded-xl border border-[var(--color-border)] bg-[var(--color-surface)] p-7"
                >
                  <div
                    aria-hidden
                    className="absolute inset-x-0 top-0 h-[3px]"
                    style={{ background: accent }}
                  />
                  <div
                    className="font-mono text-[10px] uppercase tracking-widest"
                    style={{ color: accent }}
                  >
                    {a.tag}
                  </div>
                  <h3 className="mt-2 text-lg font-semibold leading-snug">
                    {a.title}
                  </h3>
                  <p className="mt-3 text-[14px] leading-relaxed text-[var(--color-muted)]">
                    {a.body}
                  </p>
                </div>
              );
            })}
          </div>
        </div>

        {/* ── Casos de uso ─────────────────────────────────── */}
        <div className="mb-20">
          <div className="mb-8 grid grid-cols-1 gap-6 lg:grid-cols-12">
            <div className="lg:col-span-5">
              <div className="text-xs font-mono uppercase tracking-widest text-[var(--color-muted)]">
                {g.casesEyebrow}
              </div>
              <h3 className="mt-3 text-2xl font-semibold tracking-tight md:text-3xl">
                {g.casesHeading}
              </h3>
            </div>
            <div className="lg:col-span-7">
              <p className="pt-2 text-base leading-relaxed text-[var(--color-muted)]">
                {g.casesLead}
              </p>
            </div>
          </div>
          <div className="grid grid-cols-1 gap-px overflow-hidden rounded-xl border border-[var(--color-border)] bg-[var(--color-border)] sm:grid-cols-2 lg:grid-cols-3">
            {FEATURED_CASE_INDEXES.map((idx, position) => {
              const c = g.cases[idx];
              if (!c) return null;
              const accent = CASE_ACCENTS[idx % CASE_ACCENTS.length];
              const num = String(position + 1).padStart(2, "0");
              return (
                <article
                  key={c.title}
                  className="group flex h-full flex-col gap-4 bg-[var(--color-surface)] p-7 transition-colors hover:bg-[var(--color-surface-2)]"
                >
                  <div className="flex items-center justify-between">
                    <span
                      className="font-mono text-[10px] uppercase tracking-widest"
                      style={{ color: accent }}
                    >
                      {num} · {c.tag}
                    </span>
                    <span
                      className="h-1.5 w-1.5 rounded-full"
                      style={{ background: accent }}
                    />
                  </div>
                  <h4 className="text-lg font-semibold leading-snug">
                    {c.title}
                  </h4>
                  <p className="text-[13.5px] leading-relaxed text-[var(--color-muted)]">
                    {c.body}
                  </p>
                  <div className="mt-auto border-t border-[var(--color-border)] pt-3 font-mono text-[11px] text-[var(--color-muted)]">
                    {c.actors}
                  </div>
                </article>
              );
            })}
            <Link
              href="/casos"
              className="group relative flex h-full flex-col justify-between gap-4 p-7 transition-colors hover:bg-[var(--color-surface-2)]"
              style={{
                background:
                  "linear-gradient(135deg, var(--color-surface) 0%, var(--color-surface-2) 100%)",
              }}
            >
              <div
                aria-hidden
                className="absolute inset-x-0 top-0 h-[3px]"
                style={{
                  background:
                    "linear-gradient(90deg, var(--color-cr-blue-bright), var(--color-cr-red-bright))",
                }}
              />
              <div className="flex items-center justify-between">
                <span
                  className="font-mono text-[10px] uppercase tracking-widest"
                  style={{ color: "var(--color-cr-blue-bright)" }}
                >
                  {String(FEATURED_CASE_INDEXES.length + 1).padStart(2, "0")} ·{" "}
                  {g.casesCtaCard.tag}
                </span>
                <span
                  className="h-1.5 w-1.5 rounded-full"
                  style={{ background: "var(--color-cr-red-bright)" }}
                />
              </div>
              <div className="flex flex-col gap-3">
                <h4 className="text-lg font-semibold leading-snug">
                  {g.casesCtaCard.title}
                </h4>
                <p className="text-[13.5px] leading-relaxed text-[var(--color-muted)]">
                  {g.casesCtaCard.body}
                </p>
              </div>
              <div className="mt-auto flex items-center justify-between border-t border-[var(--color-border)] pt-3 font-mono text-[11px]">
                <span style={{ color: "var(--color-cr-blue-bright)" }}>
                  {g.casesCtaCard.action}
                </span>
                <span
                  aria-hidden
                  className="transition-transform group-hover:translate-x-1"
                  style={{ color: "var(--color-cr-blue-bright)" }}
                >
                  →
                </span>
              </div>
            </Link>
          </div>
        </div>

        {/* ── Garantías invariantes ─────────────────────────── */}
        <div className="border-t border-[var(--color-border)] pt-16">
          <div className="mb-8 flex flex-col gap-2">
            <div className="text-xs font-mono uppercase tracking-widest text-[var(--color-muted)]">
              {g.invariantsEyebrow}
            </div>
            <h3 className="text-2xl font-semibold tracking-tight md:text-3xl">
              {g.invariantsHeading}
            </h3>
          </div>
          <div className="grid grid-cols-1 gap-px overflow-hidden rounded-xl border border-[var(--color-border)] bg-[var(--color-border)] sm:grid-cols-2">
            {g.items.map((item, i) => {
              const accent =
                GUARANTEE_ACCENTS[i % GUARANTEE_ACCENTS.length];
              const num = String(i + 1).padStart(2, "0");
              return (
                <div
                  key={item.title}
                  className="group relative flex flex-col gap-4 bg-[var(--color-surface)] p-8 transition-colors hover:bg-[var(--color-surface-2)]"
                >
                  <div className="flex items-center justify-between">
                    <span
                      className="font-mono text-xs tracking-widest"
                      style={{ color: accent }}
                    >
                      {num}
                    </span>
                    <span
                      className="h-2 w-2 rounded-full"
                      style={{ background: accent }}
                    />
                  </div>
                  <h4 className="text-xl font-semibold">{item.title}</h4>
                  <p className="text-[15px] leading-relaxed text-[var(--color-muted)]">
                    {item.body}
                  </p>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </section>
  );
}
