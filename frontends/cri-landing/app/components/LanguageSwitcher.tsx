"use client";

import { useT } from "../i18n/LanguageProvider";
import type { Lang } from "../i18n/translations";

const LABEL: Record<Lang, string> = { es: "ES", en: "EN" };

export default function LanguageSwitcher() {
  const { lang, setLang, t } = useT();
  const langs: Lang[] = ["es", "en"];

  return (
    <div
      role="group"
      aria-label={t.langSwitch.aria}
      className="inline-flex overflow-hidden rounded-md border border-[var(--color-border)] bg-[var(--color-surface)] font-mono text-[11px]"
    >
      {langs.map((l) => {
        const active = l === lang;
        return (
          <button
            key={l}
            onClick={() => setLang(l)}
            className="px-3 py-1.5 transition-colors"
            style={{
              background: active ? "var(--color-cr-blue)" : "transparent",
              color: active ? "#fff" : "var(--color-muted)",
            }}
          >
            {LABEL[l]}
          </button>
        );
      })}
    </div>
  );
}
