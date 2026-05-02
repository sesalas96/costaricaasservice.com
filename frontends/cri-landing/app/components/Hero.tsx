"use client";

import { useT } from "../i18n/LanguageProvider";
import LanguageSwitcher from "./LanguageSwitcher";

export default function Hero() {
  const { t } = useT();
  return (
    <section className="relative overflow-hidden border-b border-[var(--color-border)]">
      <div aria-hidden className="pointer-events-none absolute inset-0 cr-glow" />
      <div className="relative mx-auto max-w-6xl px-6 py-10 md:py-12">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3 text-xs font-mono uppercase tracking-widest text-[var(--color-muted)]">
            <span
              className="inline-block h-2 w-2 rounded-full"
              style={{ background: "var(--color-cr-blue-bright)" }}
            />
            {t.hero.eyebrow}
          </div>
          <LanguageSwitcher />
        </div>
      </div>
      <div className="relative mx-auto flex max-w-6xl flex-col gap-10 px-6 pb-24 pt-10 md:pb-32">
        <h1 className="max-w-4xl text-balance text-5xl font-semibold leading-[1.05] tracking-tight md:text-7xl">
          {t.hero.titleStart}{" "}
          <span className="cr-text-gradient">{t.hero.titleHighlight}</span>
        </h1>
        <p className="max-w-2xl text-lg leading-relaxed text-[var(--color-muted)] md:text-xl">
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
        <div className="flex flex-wrap items-center gap-3">
          <a
            href="#diagramas"
            className="inline-flex items-center gap-2 rounded-md px-5 py-3 text-sm font-medium text-white transition-transform hover:-translate-y-0.5"
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
            className="inline-flex items-center gap-2 rounded-md border border-[var(--color-border)] bg-[var(--color-surface)] px-5 py-3 text-sm font-medium text-[var(--color-text)] transition-colors hover:bg-[var(--color-surface-2)]"
          >
            {t.hero.ctaSecondary}
          </a>
        </div>
        <div className="mt-6 flex flex-wrap items-center gap-x-8 gap-y-2 font-mono text-xs text-[var(--color-muted)]">
          <span>By @devsebas ~ sebashian961@gmail.com</span>
        </div>
      </div>
    </section>
  );
}
