"use client";

import { useT } from "../i18n/LanguageProvider";

export default function Footer() {
  const { t } = useT();
  return (
    <footer className="border-t border-[var(--color-border)] py-10">
      <div className="mx-auto flex max-w-6xl items-center gap-3 px-6 text-sm text-[var(--color-muted)]">
        <span
          className="inline-block h-2 w-2 rounded-full"
          style={{ background: "var(--color-cr-blue-bright)" }}
        />
        <span className="font-mono text-xs uppercase tracking-widest">
          costaricaasservice
        </span>
        <span>· {t.footer.coreOf}</span>
      </div>
    </footer>
  );
}
