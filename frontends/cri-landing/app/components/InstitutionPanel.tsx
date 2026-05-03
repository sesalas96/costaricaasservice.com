"use client";

import { useEffect, useState } from "react";
import { createPortal } from "react-dom";
import { useT } from "../i18n/LanguageProvider";
import {
  INSTITUTIONS,
  getInstitution,
  localized,
  type InstitutionDef,
} from "../data/institutions";
import { useInstitution } from "./InstitutionContext";

function StatusBadge({ status }: { status: InstitutionDef["status"] }) {
  const { t } = useT();
  const isLive = status === "live";
  return (
    <span
      className="inline-flex items-center gap-1 rounded-full px-2 py-0.5 font-mono text-[10px] uppercase tracking-widest"
      style={{
        background: isLive ? "rgba(16,185,129,0.12)" : "rgba(139,147,163,0.12)",
        color: isLive ? "var(--color-area-iduc)" : "var(--color-muted)",
        border: `1px solid ${
          isLive ? "var(--color-area-iduc)" : "var(--color-border)"
        }`,
      }}
    >
      <span
        className="h-1.5 w-1.5 rounded-full"
        style={{
          background: isLive ? "var(--color-area-iduc)" : "var(--color-muted)",
        }}
      />
      {isLive ? t.institutionPanel.badgeLive : t.institutionPanel.badgePlanned}
    </span>
  );
}

function ListView() {
  const { t, lang } = useT();
  const { openFromList } = useInstitution();
  const liveCount = INSTITUTIONS.filter((i) => i.status === "live").length;
  const plannedCount = INSTITUTIONS.length - liveCount;
  const counter = t.institutionPanel.listCount
    .replace("{n}", String(INSTITUTIONS.length))
    .replace("{l}", String(liveCount))
    .replace("{p}", String(plannedCount));

  return (
    <>
      <header className="border-b border-[var(--color-border)] p-6">
        <h2 className="text-xl font-semibold tracking-tight">
          {t.institutionPanel.listHeading}
        </h2>
        <p className="mt-2 text-[13px] leading-relaxed text-[var(--color-muted)]">
          {t.institutionPanel.listSub}
        </p>
        <div className="mt-3 font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted)]">
          {counter}
        </div>
      </header>
      <ul className="divide-y divide-[var(--color-border)]">
        {INSTITUTIONS.map((inst) => {
          const loc = localized(inst, lang);
          return (
            <li key={inst.id}>
              <button
                onClick={() => openFromList(inst.id)}
                className="flex w-full items-start gap-3 p-5 text-left transition-colors hover:bg-[var(--color-surface-2)]"
              >
                <span
                  className="mt-1 h-8 w-8 shrink-0 rounded-md"
                  style={{
                    background: `linear-gradient(135deg, ${inst.accent}, transparent)`,
                    border: `1px solid ${inst.accent}`,
                  }}
                />
                <div className="min-w-0 flex-1">
                  <div className="flex flex-wrap items-center gap-2">
                    <span
                      className="font-mono text-[10px] uppercase tracking-widest"
                      style={{ color: inst.accent }}
                    >
                      {inst.tag}
                    </span>
                    <StatusBadge status={inst.status} />
                  </div>
                  <div className="mt-1 truncate text-[14px] font-semibold text-[var(--color-text)]">
                    {loc.name}
                  </div>
                  <div className="mt-0.5 text-[12px] leading-snug text-[var(--color-muted)]">
                    {loc.short}
                  </div>
                </div>
                <span
                  aria-hidden
                  className="shrink-0 self-center text-[var(--color-muted)]"
                >
                  →
                </span>
              </button>
            </li>
          );
        })}
      </ul>
    </>
  );
}

function Section({
  title,
  items,
  empty,
}: {
  title: string;
  items: string[];
  empty?: string;
}) {
  if (items.length === 0 && !empty) return null;
  return (
    <section>
      <h3 className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted)]">
        {title}
      </h3>
      {items.length === 0 ? (
        <p className="mt-2 text-[13px] italic text-[var(--color-muted)]">
          {empty}
        </p>
      ) : (
        <ul className="mt-2 space-y-1.5">
          {items.map((item) => (
            <li
              key={item}
              className="rounded-md border border-[var(--color-border)] bg-[var(--color-bg)] px-3 py-2 font-mono text-[12px] leading-snug text-[var(--color-text)]"
            >
              {item}
            </li>
          ))}
        </ul>
      )}
    </section>
  );
}

function DetailView({
  inst,
  fromList,
}: {
  inst: InstitutionDef;
  fromList: boolean;
}) {
  const { t, lang } = useT();
  const { goBackToList } = useInstitution();
  const loc = localized(inst, lang);

  return (
    <>
      <header className="border-b border-[var(--color-border)] p-6">
        {fromList && (
          <button
            onClick={goBackToList}
            className="mb-3 font-mono text-[11px] uppercase tracking-widest text-[var(--color-muted)] transition-colors hover:text-[var(--color-text)]"
          >
            {t.institutionPanel.backToList}
          </button>
        )}
        <div className="flex items-start gap-3">
          <span
            className="mt-1 h-10 w-10 shrink-0 rounded-md"
            style={{
              background: `linear-gradient(135deg, ${inst.accent}, transparent)`,
              border: `1px solid ${inst.accent}`,
            }}
          />
          <div className="min-w-0 flex-1">
            <div className="flex flex-wrap items-center gap-2">
              <span
                className="font-mono text-[10px] uppercase tracking-widest"
                style={{ color: inst.accent }}
              >
                {inst.tag}
              </span>
              <StatusBadge status={inst.status} />
            </div>
            <h2 className="mt-1 text-xl font-semibold tracking-tight">
              {loc.name}
            </h2>
            <p className="mt-1 text-[13px] leading-relaxed text-[var(--color-muted)]">
              {loc.short}
            </p>
          </div>
        </div>
      </header>
      <div className="space-y-6 p-6">
        <Section
          title={t.institutionPanel.sectionCustodia}
          items={loc.custodia}
        />
        <Section
          title={t.institutionPanel.sectionExpone}
          items={loc.expone}
        />
        <Section
          title={t.institutionPanel.sectionConsume}
          items={loc.consume}
          empty={t.institutionPanel.emptyConsume}
        />
      </div>
    </>
  );
}

export default function InstitutionPanel() {
  const { selected, close } = useInstitution();
  const { t } = useT();
  const [mounted, setMounted] = useState(false);
  const isOpen = selected !== null;

  useEffect(() => setMounted(true), []);

  useEffect(() => {
    if (!isOpen) return;
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") close();
    };
    document.addEventListener("keydown", onKey);
    const original = document.body.style.overflow;
    document.body.style.overflow = "hidden";
    return () => {
      document.removeEventListener("keydown", onKey);
      document.body.style.overflow = original;
    };
  }, [isOpen, close]);

  if (!mounted) return null;

  const detailInst =
    selected?.type === "detail" ? getInstitution(selected.id) : undefined;

  return createPortal(
    <div
      className="institution-panel-overlay"
      data-open={isOpen}
      aria-hidden={!isOpen}
    >
      <div
        className="institution-panel-backdrop"
        onClick={close}
        aria-hidden
      />
      <aside
        role="dialog"
        aria-modal="true"
        aria-label={
          selected?.type === "detail"
            ? detailInst
              ? localized(detailInst, "es").name
              : t.institutionPanel.listHeading
            : t.institutionPanel.listHeading
        }
        className="institution-panel-aside"
      >
        <button
          onClick={close}
          aria-label={t.institutionPanel.closeAria}
          className="absolute right-3 top-3 z-10 flex h-10 w-10 items-center justify-center rounded-md border border-[var(--color-border)] bg-[var(--color-surface)] text-base text-[var(--color-muted)] transition-colors hover:bg-[var(--color-surface-2)] hover:text-[var(--color-text)] sm:right-4 sm:top-4"
        >
          ×
        </button>
        <div className="h-full overflow-y-auto">
          {selected?.type === "list" && <ListView />}
          {selected?.type === "detail" && detailInst && (
            <DetailView inst={detailInst} fromList={selected.fromList} />
          )}
        </div>
      </aside>
    </div>,
    document.body,
  );
}
