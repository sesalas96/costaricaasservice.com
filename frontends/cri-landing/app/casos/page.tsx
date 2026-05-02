"use client";

import Link from "next/link";
import { useT } from "../i18n/LanguageProvider";
import LanguageSwitcher from "../components/LanguageSwitcher";
import Footer from "../components/Footer";

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

const CATEGORY_ORDER = ["gobierno", "identidad", "bienestar", "justicia"] as const;

const CATEGORY_ACCENTS: Record<(typeof CATEGORY_ORDER)[number], string> = {
  gobierno: "var(--color-cr-blue-bright)",
  identidad: "var(--color-area-iduc)",
  bienestar: "var(--color-area-platform)",
  justicia: "var(--color-cr-red-bright)",
};

export default function CasosPage() {
  const { t } = useT();
  const p = t.casesPage;
  const cases = t.guarantees.cases;

  const grouped = CATEGORY_ORDER.map((key) => ({
    key,
    label: p.categories[key].label,
    description: p.categories[key].description,
    accent: CATEGORY_ACCENTS[key],
    items: cases
      .map((c, idx) => ({ c, idx }))
      .filter(({ c }) => c.category === key),
  }));

  const totalCases = cases.length;
  const totalDomains = grouped.filter((g) => g.items.length > 0).length;
  const countLabel = p.casesCount
    .replace("{n}", String(totalCases))
    .replace("{c}", String(totalDomains));

  return (
    <main className="min-h-screen">
      <section className="relative overflow-hidden border-b border-[var(--color-border)]">
        <div aria-hidden className="pointer-events-none absolute inset-0 cr-glow" />
        <div className="relative mx-auto max-w-6xl px-6 py-10 md:py-12">
          <div className="flex items-center justify-between gap-4">
            <Link
              href="/"
              className="inline-flex items-center gap-2 text-xs font-mono uppercase tracking-widest text-[var(--color-muted)] transition-colors hover:text-[var(--color-text)]"
            >
              <span aria-hidden>←</span>
              {p.backToHome}
            </Link>
            <LanguageSwitcher />
          </div>
        </div>
        <div className="relative mx-auto flex max-w-6xl flex-col gap-8 px-6 pb-16 pt-6 md:pb-20">
          <div className="flex items-center gap-3 text-xs font-mono uppercase tracking-widest text-[var(--color-muted)]">
            <span
              className="inline-block h-2 w-2 rounded-full"
              style={{ background: "var(--color-cr-blue-bright)" }}
            />
            {p.eyebrow}
          </div>
          <h1 className="max-w-4xl text-balance text-4xl font-semibold leading-[1.1] tracking-tight md:text-6xl">
            {p.heading}
          </h1>
          <p className="max-w-3xl text-lg leading-relaxed text-[var(--color-muted)]">
            {p.lead}
          </p>
          <div className="font-mono text-xs uppercase tracking-widest text-[var(--color-muted)]">
            {countLabel}
          </div>
        </div>
      </section>

      <section className="border-b border-[var(--color-border)] py-20">
        <div className="mx-auto flex max-w-6xl flex-col gap-20 px-6">
          {grouped.map((group) => {
            if (group.items.length === 0) return null;
            return (
              <div key={group.key}>
                <div className="mb-8 grid grid-cols-1 gap-6 lg:grid-cols-12">
                  <div className="lg:col-span-5">
                    <div
                      className="font-mono text-[10px] uppercase tracking-widest"
                      style={{ color: group.accent }}
                    >
                      {group.key}
                    </div>
                    <h2 className="mt-3 text-2xl font-semibold tracking-tight md:text-3xl">
                      {group.label}
                    </h2>
                  </div>
                  <div className="lg:col-span-7">
                    <p className="pt-2 text-base leading-relaxed text-[var(--color-muted)]">
                      {group.description}
                    </p>
                    <div className="mt-3 font-mono text-[11px] uppercase tracking-widest text-[var(--color-muted)]">
                      {group.items.length} ·{" "}
                      {group.items
                        .map(({ c }) => c.tag)
                        .join(" · ")}
                    </div>
                  </div>
                </div>
                <div className="grid grid-cols-1 gap-px overflow-hidden rounded-xl border border-[var(--color-border)] bg-[var(--color-border)] sm:grid-cols-2 lg:grid-cols-3">
                  {group.items.map(({ c, idx }, position) => {
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
                        <h3 className="text-lg font-semibold leading-snug">
                          {c.title}
                        </h3>
                        <p className="text-[13.5px] leading-relaxed text-[var(--color-muted)]">
                          {c.body}
                        </p>
                        <div className="mt-auto border-t border-[var(--color-border)] pt-3 font-mono text-[11px] text-[var(--color-muted)]">
                          {c.actors}
                        </div>
                      </article>
                    );
                  })}
                </div>
              </div>
            );
          })}
        </div>
      </section>

      <Footer />
    </main>
  );
}
