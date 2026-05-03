"use client";

import { useT } from "../i18n/LanguageProvider";
import LanguageSwitcher from "./LanguageSwitcher";

export default function Hero() {
  const { t } = useT();
  return (
    <section className="relative overflow-hidden border-b border-[var(--color-border)]">
      <div aria-hidden className="pointer-events-none absolute inset-0 cr-glow" />
      <div className="relative mx-auto max-w-6xl px-5 py-6 sm:px-6 md:py-12">
        <div className="flex items-center justify-between gap-3">
          <div className="flex min-w-0 items-center gap-2 font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted)] sm:gap-3 sm:text-xs">
            <span
              className="inline-block h-2 w-2 shrink-0 rounded-full"
              style={{ background: "var(--color-cr-blue-bright)" }}
            />
            <span className="truncate">{t.hero.eyebrow}</span>
          </div>
          <LanguageSwitcher />
        </div>
      </div>
      <div className="relative mx-auto flex max-w-6xl flex-col gap-6 px-5 pb-16 pt-4 sm:px-6 md:gap-10 md:pb-32 md:pt-10">
        <h1 className="max-w-4xl text-balance text-4xl font-semibold leading-[1.08] tracking-tight sm:text-5xl md:text-7xl">
          {t.hero.titleStart}{" "}
          <span className="cr-text-gradient">{t.hero.titleHighlight}</span>
        </h1>
        <p className="max-w-2xl text-base leading-relaxed text-[var(--color-muted)] sm:text-lg md:text-xl">
          {t.hero.descBefore}
          <strong className="text-[var(--color-text)]">
            {t.hero.descBold1}
          </strong>
          {t.hero.descMid1}
          <strong className="text-[var(--color-text)]">
            {t.hero.descBold2}
          </strong>
          {t.hero.descMid2}
          <strong className="text-[var(--color-text)]">
            {t.hero.descBold3}
          </strong>
          {t.hero.descAfter}
        </p>
        <div className="flex flex-col gap-3 sm:flex-row sm:flex-wrap sm:items-center">
          <a
            href="#diagramas"
            className="inline-flex items-center justify-center gap-2 rounded-md px-5 py-3 text-sm font-medium text-white transition-transform hover:-translate-y-0.5"
            style={{
              background:
                "linear-gradient(120deg, var(--color-cr-blue) 0%, var(--color-cr-blue-bright) 100%)",
              boxShadow: "0 8px 24px -8px rgba(29,79,196,0.6)",
            }}
          >
            {t.hero.ctaPrimary}
            <span aria-hidden>→</span>
          </a>
          <a
            href="#garantias"
            className="inline-flex items-center justify-center gap-2 rounded-md border border-[var(--color-border)] bg-[var(--color-surface)] px-5 py-3 text-sm font-medium text-[var(--color-text)] transition-colors hover:bg-[var(--color-surface-2)]"
          >
            {t.hero.ctaSecondary}
          </a>
        </div>
        <div className="mt-2 flex flex-wrap items-center gap-x-8 gap-y-2 font-mono text-[11px] text-[var(--color-muted)] sm:mt-6 sm:text-xs">
          <span className="break-all">By @devsebas ~ sebashian961@gmail.com</span>
        </div>
      </div>
    </section>
  );
}
