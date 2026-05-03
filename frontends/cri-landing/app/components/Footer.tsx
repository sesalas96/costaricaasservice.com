"use client";

import { useT } from "../i18n/LanguageProvider";

const PLATFORM_ROWS = [
  { key: "gateway", color: "var(--color-area-gateway)" },
  { key: "iduc", color: "var(--color-area-iduc)" },
  { key: "interop", color: "var(--color-area-interop)" },
  { key: "members", color: "var(--color-area-members)" },
  { key: "platform", color: "var(--color-area-platform)" },
  { key: "frontends", color: "var(--color-area-frontend)" },
] as const;

const SUBDOMAINS = [
  "app",
  "members",
  "admin",
  "api",
  "interop",
  "docs",
] as const;

const ECOSYSTEM_LINKS = [
  { key: "cases", href: "/casos" },
  { key: "diagrams", href: "#diagramas" },
  { key: "demo", href: "#demo" },
  { key: "deliverables", href: "#entregables" },
] as const;

export default function Footer() {
  const { t } = useT();
  const year = new Date().getFullYear();

  return (
    <footer className="relative overflow-hidden border-t border-[var(--color-border)] bg-[var(--color-bg)]">
      <div aria-hidden className="cr-flag-bar" />
      <div
        aria-hidden
        className="pointer-events-none absolute inset-0 cr-glow opacity-50"
      />
      <div className="relative mx-auto max-w-6xl px-5 pb-10 pt-14 sm:px-6 sm:pb-14 md:pt-20">
        <div className="grid gap-10 md:grid-cols-12 md:gap-14">
          <div className="md:col-span-5">
            <div className="flex items-center gap-3">
              <span
                className="inline-block h-2 w-2 rounded-full"
                style={{ background: "var(--color-cr-blue-bright)" }}
              />
              <span className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted)]">
                {t.footer.coreOf}
              </span>
            </div>
            <p className="mt-3 text-3xl font-semibold leading-[1.05] tracking-tight sm:text-4xl">
              <span className="cr-text-gradient">costaricaasservice</span>
            </p>
            <p className="mt-4 max-w-md text-sm leading-relaxed text-[var(--color-muted)] sm:text-base">
              {t.footer.tagline}
            </p>
            <div className="mt-6 inline-flex items-center gap-2 rounded-full border border-[var(--color-border)] bg-[var(--color-surface)] px-3 py-1 font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted)]">
              <span className="relative flex h-2 w-2">
                <span
                  className="absolute inline-flex h-full w-full animate-ping rounded-full opacity-60"
                  style={{ background: "var(--color-area-iduc)" }}
                />
                <span
                  className="relative inline-flex h-2 w-2 rounded-full"
                  style={{ background: "var(--color-area-iduc)" }}
                />
              </span>
              {t.footer.statusBadge}
            </div>
          </div>

          <div className="grid grid-cols-2 gap-8 md:col-span-7 md:grid-cols-3 md:gap-10">
            <div>
              <h3 className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted)]">
                {t.footer.cols.platform.heading}
              </h3>
              <ul className="mt-4 space-y-2.5 text-sm">
                {PLATFORM_ROWS.map((row) => (
                  <li key={row.key} className="flex items-center gap-2.5">
                    <span
                      aria-hidden
                      className="inline-block h-1.5 w-1.5 shrink-0 rounded-full"
                      style={{ background: row.color }}
                    />
                    <span className="text-[var(--color-text)]">
                      {t.footer.cols.platform[row.key]}
                    </span>
                  </li>
                ))}
              </ul>
            </div>

            <div>
              <h3 className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted)]">
                {t.footer.cols.audiences.heading}
              </h3>
              <ul className="mt-4 space-y-2.5 text-sm text-[var(--color-text)]">
                <li>{t.footer.cols.audiences.citizen}</li>
                <li>{t.footer.cols.audiences.member}</li>
                <li>{t.footer.cols.audiences.realm}</li>
                <li>{t.footer.cols.audiences.admin}</li>
              </ul>
            </div>

            <div className="col-span-2 md:col-span-1">
              <h3 className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted)]">
                {t.footer.cols.ecosystem.heading}
              </h3>
              <ul className="mt-4 space-y-2.5 text-sm">
                {ECOSYSTEM_LINKS.map((link) => (
                  <li key={link.key}>
                    <a
                      href={link.href}
                      className="text-[var(--color-text)] transition-colors hover:text-white"
                    >
                      {t.footer.cols.ecosystem[link.key]}
                    </a>
                  </li>
                ))}
              </ul>
            </div>
          </div>
        </div>

        <div className="mt-12 rounded-xl border border-[var(--color-border)] bg-[var(--color-surface)]/60 p-5 sm:p-6">
          <div className="flex flex-wrap items-baseline justify-between gap-x-4 gap-y-2">
            <div className="flex items-baseline gap-3">
              <span className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted)]">
                {t.footer.domains.heading}
              </span>
              <span className="font-mono text-sm text-[var(--color-text)]">
                costaricaasservice.com
              </span>
            </div>
            <span className="font-mono text-[11px] text-[var(--color-muted)]">
              {t.footer.domains.hint}
            </span>
          </div>
          <ul className="mt-4 grid grid-cols-2 gap-2 sm:grid-cols-3 md:grid-cols-6">
            {SUBDOMAINS.map((sub) => (
              <li
                key={sub}
                className="rounded-md border border-[var(--color-border)] bg-[var(--color-bg)] px-3 py-2"
              >
                <div className="font-mono text-xs">
                  <span className="text-[var(--color-text)]">{sub}</span>
                  <span className="text-[var(--color-muted)]">
                    .costaricaasservice.com
                  </span>
                </div>
                <div className="mt-1 flex items-center gap-1.5 font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted)]">
                  <span
                    aria-hidden
                    className="inline-block h-1.5 w-1.5 rounded-full"
                    style={{ background: "var(--color-cr-red-bright)" }}
                  />
                  {t.footer.domains.soon}
                </div>
              </li>
            ))}
          </ul>
        </div>

        <div className="mt-10 flex flex-col gap-3 border-t border-[var(--color-border)] pt-6 text-[var(--color-muted)] md:flex-row md:items-center md:justify-between">
          <div className="flex flex-wrap items-center gap-x-3 gap-y-1 font-mono text-[11px] uppercase tracking-widest">
            <span>© {year} costaricaasservice</span>
            <span aria-hidden>·</span>
            <span>{t.footer.bottom.madeIn}</span>
          </div>
          <div className="flex flex-wrap items-center gap-x-3 gap-y-1 font-mono text-[11px]">
            <a
              href="mailto:sebashian961@gmail.com"
              className="transition-colors hover:text-white"
            >
              sebashian961@gmail.com
            </a>
            <span aria-hidden>·</span>
            <span>by @devsebas</span>
          </div>
        </div>
        <p className="mt-4 max-w-3xl font-mono text-[10px] leading-relaxed text-[var(--color-muted)]/70">
          {t.footer.bottom.disclaimer}
        </p>
      </div>
    </footer>
  );
}
