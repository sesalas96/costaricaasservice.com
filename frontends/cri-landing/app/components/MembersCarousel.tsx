"use client";

import { useT } from "../i18n/LanguageProvider";
import {
  INSTITUTIONS,
  localized,
  type InstitutionDef,
} from "../data/institutions";
import { useInstitution } from "./InstitutionContext";

function MemberCard({ inst }: { inst: InstitutionDef }) {
  const { t, lang } = useT();
  const { open } = useInstitution();
  const loc = localized(inst, lang);
  const isLive = inst.status === "live";
  const statusLabel = isLive
    ? t.membersCarousel.liveLabel
    : t.membersCarousel.plannedLabel;

  return (
    <button
      onClick={() => open(inst.id)}
      aria-label={`${loc.name} — ${loc.short}`}
      className="group mx-2 flex h-[180px] w-[280px] shrink-0 flex-col rounded-xl border border-[var(--color-border)] bg-[var(--color-surface)] p-4 text-left transition-colors hover:bg-[var(--color-surface-2)] focus:outline-none focus:ring-2 focus:ring-[var(--color-cr-blue-bright)] sm:w-[300px]"
      style={{
        boxShadow: "0 1px 0 rgba(255,255,255,0.02) inset",
      }}
    >
      <div
        aria-hidden
        className="-mx-4 -mt-4 mb-3 h-[3px] rounded-t-xl"
        style={{ background: inst.accent }}
      />
      <div className="flex items-start justify-between gap-2">
        <div className="flex items-center gap-2">
          <span
            className="h-7 w-7 shrink-0 rounded-md"
            style={{
              background: `linear-gradient(135deg, ${inst.accent}, transparent)`,
              border: `1px solid ${inst.accent}`,
            }}
          />
          <span
            className="font-mono text-[10px] uppercase tracking-widest"
            style={{ color: inst.accent }}
          >
            {inst.tag}
          </span>
        </div>
        <span
          className="inline-flex items-center gap-1 rounded-full px-2 py-0.5 font-mono text-[9px] uppercase tracking-widest"
          style={{
            background: isLive
              ? "rgba(16,185,129,0.12)"
              : "rgba(139,147,163,0.10)",
            color: isLive ? "var(--color-area-iduc)" : "var(--color-muted)",
            border: `1px solid ${
              isLive ? "var(--color-area-iduc)" : "var(--color-border)"
            }`,
          }}
        >
          <span
            className="h-1 w-1 rounded-full"
            style={{
              background: isLive
                ? "var(--color-area-iduc)"
                : "var(--color-muted)",
            }}
          />
          {statusLabel}
        </span>
      </div>
      <h3 className="mt-3 line-clamp-2 text-[14px] font-semibold leading-snug text-[var(--color-text)]">
        {loc.name}
      </h3>
      <p className="mt-1 line-clamp-2 text-[12px] leading-snug text-[var(--color-muted)]">
        {loc.short}
      </p>
      <div className="mt-auto flex items-center justify-between border-t border-[var(--color-border)] pt-2 font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted)]">
        <span>
          {loc.expone.length}{" "}
          {loc.expone.length === 1 ? "endpoint" : "endpoints"}
        </span>
        <span
          aria-hidden
          className="transition-transform group-hover:translate-x-0.5"
          style={{ color: inst.accent }}
        >
          →
        </span>
      </div>
    </button>
  );
}

export default function MembersCarousel() {
  const { t } = useT();
  const { open } = useInstitution();
  const c = t.membersCarousel;

  // Duplicamos la lista para que el loop sea continuo (translateX(-50%)).
  const items = [...INSTITUTIONS, ...INSTITUTIONS];

  return (
    <section
      id="members"
      className="border-b border-[var(--color-border)] py-16 md:py-24"
    >
      <div className="mx-auto max-w-6xl px-5 sm:px-6">
        <div className="mb-10 grid grid-cols-1 gap-6 md:mb-12 lg:grid-cols-12 lg:gap-12">
          <div className="lg:col-span-5">
            <div className="text-[10px] font-mono uppercase tracking-widest text-[var(--color-muted)] sm:text-xs">
              {c.eyebrow}
            </div>
            <h2 className="mt-3 text-2xl font-semibold tracking-tight sm:text-3xl md:text-4xl">
              {c.heading}
            </h2>
          </div>
          <div className="lg:col-span-7">
            <p className="text-base leading-relaxed text-[var(--color-muted)] sm:text-lg">
              {c.lead}
            </p>
          </div>
        </div>
      </div>

      {/* Marquee a ancho completo (rompe el max-w para que se sienta infinito) */}
      <div
        className="members-marquee-viewport py-2"
        aria-roledescription="carousel"
      >
        <ul className="members-marquee-track" role="list">
          {items.map((inst, i) => (
            <li key={`${inst.id}-${i}`} aria-hidden={i >= INSTITUTIONS.length}>
              <MemberCard inst={inst} />
            </li>
          ))}
        </ul>
      </div>

      <div className="mx-auto mt-6 flex max-w-6xl flex-col items-start gap-3 px-5 sm:flex-row sm:items-center sm:justify-between sm:px-6">
        <p className="font-mono text-[11px] text-[var(--color-muted)]">
          {c.hint}
        </p>
        <button
          onClick={() => open("__list__")}
          className="font-mono text-[11px] uppercase tracking-widest text-[var(--color-cr-blue-bright)] transition-colors hover:text-[var(--color-text)]"
        >
          {c.cta} →
        </button>
      </div>
    </section>
  );
}
