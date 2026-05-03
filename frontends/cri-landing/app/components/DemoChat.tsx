"use client";

import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useT } from "../i18n/LanguageProvider";
import {
  ACTOR_LABELS,
  DEMO_CASES,
  localizedCase,
  type DemoActor,
  type DemoCaseDef,
  type DemoStep,
  type DemoStepKind,
} from "../data/demoCases";

const DAILY_QUOTA = 5;
const STORAGE_KEY = "cri-landing.demo-usage";

type ChatMessage =
  | { id: string; role: "citizen"; text: string }
  | { id: string; role: "llm"; text: string }
  | {
      id: string;
      role: "system";
      step: DemoStep;
    }
  | {
      id: string;
      role: "email";
      subject: string;
      preview: string;
      to: string;
    };

type Usage = { date: string; count: number };

const todayKey = () => {
  const d = new Date();
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`;
};

const readUsage = (): Usage => {
  if (typeof window === "undefined") return { date: todayKey(), count: 0 };
  try {
    const raw = window.localStorage.getItem(STORAGE_KEY);
    if (!raw) return { date: todayKey(), count: 0 };
    const parsed = JSON.parse(raw) as Usage;
    if (parsed.date !== todayKey()) return { date: todayKey(), count: 0 };
    return parsed;
  } catch {
    return { date: todayKey(), count: 0 };
  }
};

const writeUsage = (u: Usage) => {
  if (typeof window === "undefined") return;
  window.localStorage.setItem(STORAGE_KEY, JSON.stringify(u));
};

function actorBadgeColor(actor: DemoActor): string {
  return ACTOR_LABELS[actor]?.color ?? "var(--color-muted)";
}

function actorLabel(actor: DemoActor, lang: "es" | "en"): string {
  return ACTOR_LABELS[actor]?.[lang] ?? actor;
}

function actorTag(actor: DemoActor): string {
  return ACTOR_LABELS[actor]?.tag ?? actor.toUpperCase();
}

function useReducedMotion(): boolean {
  const [reduced, setReduced] = useState(false);
  useEffect(() => {
    if (typeof window === "undefined" || !window.matchMedia) return;
    const mq = window.matchMedia("(prefers-reduced-motion: reduce)");
    setReduced(mq.matches);
    const onChange = (e: MediaQueryListEvent) => setReduced(e.matches);
    mq.addEventListener("change", onChange);
    return () => mq.removeEventListener("change", onChange);
  }, []);
  return reduced;
}

function Typewriter({
  text,
  speedMs = 14,
  cursor = false,
  onTick,
}: {
  text: string;
  speedMs?: number;
  cursor?: boolean;
  onTick?: () => void;
}) {
  const reduced = useReducedMotion();
  const [count, setCount] = useState(0);
  const onTickRef = useRef(onTick);
  useEffect(() => {
    onTickRef.current = onTick;
  }, [onTick]);

  useEffect(() => {
    if (reduced) {
      setCount(text.length);
      return;
    }
    setCount(0);
    let i = 0;
    let cancelled = false;
    let timeout: ReturnType<typeof setTimeout>;
    const tick = () => {
      if (cancelled) return;
      i += 1;
      setCount(i);
      onTickRef.current?.();
      if (i < text.length) {
        timeout = setTimeout(tick, speedMs);
      }
    };
    timeout = setTimeout(tick, speedMs);
    return () => {
      cancelled = true;
      clearTimeout(timeout);
    };
  }, [text, speedMs, reduced]);

  const done = count >= text.length;
  return (
    <span>
      {text.slice(0, count)}
      {cursor && !done ? (
        <span
          aria-hidden
          className="ml-0.5 inline-block w-[2px] animate-pulse align-[-0.1em] bg-current"
          style={{ height: "0.95em" }}
        />
      ) : null}
    </span>
  );
}

function StepBadge({ kind, label }: { kind: DemoStepKind; label: string }) {
  const tone =
    kind === "request"
      ? "var(--color-cr-blue-bright)"
      : kind === "response"
        ? "var(--color-area-iduc)"
        : kind === "audit"
          ? "var(--color-area-interop)"
          : kind === "notify"
            ? "var(--color-area-platform)"
            : "var(--color-muted)";
  return (
    <span
      className="rounded-sm px-1.5 py-0.5 font-mono text-[9px] uppercase tracking-widest"
      style={{ background: `${tone}15`, color: tone, border: `1px solid ${tone}40` }}
    >
      {label}
    </span>
  );
}

export default function DemoChat() {
  const { t, lang } = useT();
  const d = t.demo;

  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [running, setRunning] = useState(false);
  const [activeCaseId, setActiveCaseId] = useState<string | null>(null);
  const [completed, setCompleted] = useState(false);
  const [usage, setUsage] = useState<Usage>({ date: todayKey(), count: 0 });
  const [mounted, setMounted] = useState(false);

  const timersRef = useRef<ReturnType<typeof setTimeout>[]>([]);
  const chatScrollRef = useRef<HTMLDivElement>(null);
  const chatPanelRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    setMounted(true);
    setUsage(readUsage());
  }, []);

  useEffect(() => {
    if (mounted) writeUsage(usage);
  }, [usage, mounted]);

  const scrollToBottom = useCallback(() => {
    const el = chatScrollRef.current;
    if (el) el.scrollTop = el.scrollHeight;
  }, []);

  useEffect(() => {
    scrollToBottom();
  }, [messages, scrollToBottom]);

  useEffect(
    () => () => {
      timersRef.current.forEach(clearTimeout);
    },
    [],
  );

  const cancelTimers = () => {
    timersRef.current.forEach(clearTimeout);
    timersRef.current = [];
  };

  const exceeded = mounted && usage.count >= DAILY_QUOTA;

  const runCase = (def: DemoCaseDef) => {
    if (running) return;
    if (exceeded) return;

    cancelTimers();
    const loc = localizedCase(def, lang);

    const seed: ChatMessage[] = [
      { id: `${def.id}-citizen-${Date.now()}`, role: "citizen", text: loc.citizenAsk },
      { id: `${def.id}-llm-intro-${Date.now()}`, role: "llm", text: loc.llmIntro },
    ];
    setMessages(seed);
    setRunning(true);
    setCompleted(false);
    setActiveCaseId(def.id);
    setUsage((u) => ({ date: todayKey(), count: u.count + 1 }));

    if (typeof window !== "undefined" && window.innerWidth < 1024) {
      requestAnimationFrame(() => {
        chatPanelRef.current?.scrollIntoView({ behavior: "smooth", block: "start" });
      });
    }

    let cumulative = 700;
    loc.steps.forEach((step, idx) => {
      const t = setTimeout(() => {
        setMessages((prev) => [
          ...prev,
          { id: `${def.id}-step-${idx}-${Date.now()}`, role: "system", step },
        ]);
      }, cumulative);
      timersRef.current.push(t);
      cumulative += step.delayMs;
    });

    const summaryDelay = cumulative + 400;
    timersRef.current.push(
      setTimeout(() => {
        setMessages((prev) => [
          ...prev,
          { id: `${def.id}-summary-${Date.now()}`, role: "llm", text: loc.finalSummary },
        ]);
      }, summaryDelay),
    );

    const emailDelay = summaryDelay + 600;
    timersRef.current.push(
      setTimeout(() => {
        setMessages((prev) => [
          ...prev,
          {
            id: `${def.id}-email-${Date.now()}`,
            role: "email",
            subject: loc.emailSubject,
            preview: loc.emailPreview,
            to: "tu@correo.cr",
          },
        ]);
        setRunning(false);
        setCompleted(true);
      }, emailDelay),
    );
  };

  const handleReset = () => {
    cancelTimers();
    setMessages([]);
    setRunning(false);
    setCompleted(false);
    setActiveCaseId(null);
    setUsage({ date: todayKey(), count: 0 });
  };

  const activeCase = useMemo(
    () => DEMO_CASES.find((c) => c.id === activeCaseId) ?? null,
    [activeCaseId],
  );

  const usagePct = Math.min(100, (usage.count / DAILY_QUOTA) * 100);

  return (
    <section
      id="demo"
      className="relative border-b border-[var(--color-border)] py-16 md:py-24"
    >
      <div aria-hidden className="pointer-events-none absolute inset-0 cr-glow opacity-60" />
      <div className="relative mx-auto max-w-6xl px-5 sm:px-6">
        {/* ── Cabezal ──────────────────────────────────────────── */}
        <div className="mb-10 grid grid-cols-1 gap-8 md:mb-14 lg:grid-cols-12 lg:gap-16">
          <div className="lg:col-span-5">
            <div className="text-[10px] font-mono uppercase tracking-widest text-[var(--color-muted)] sm:text-xs">
              {d.eyebrow}
            </div>
            <h2 className="mt-3 text-2xl font-semibold tracking-tight sm:text-3xl md:text-4xl">
              {d.heading}
            </h2>
          </div>
          <div className="lg:col-span-7">
            <p className="text-base leading-relaxed text-[var(--color-muted)] sm:text-lg">
              {d.lead}
            </p>
          </div>
        </div>

        {/* ── Cuota diaria ─────────────────────────────────────── */}
        <div className="mb-8 overflow-hidden rounded-xl border border-[var(--color-border)] bg-[var(--color-surface)]">
          <div className="flex flex-col gap-3 p-4 sm:flex-row sm:items-center sm:justify-between sm:p-5">
            <div className="flex flex-col gap-1">
              <div className="flex items-center gap-2 font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted)]">
                <span
                  className="inline-block h-2 w-2 rounded-full"
                  style={{
                    background: exceeded ? "var(--color-cr-red-bright)" : "var(--color-cr-blue-bright)",
                  }}
                />
                {d.quotaLabel}
              </div>
              <div className="flex items-baseline gap-2">
                <span className="text-2xl font-semibold tabular-nums sm:text-3xl">
                  {mounted ? usage.count : 0}
                </span>
                <span className="text-sm text-[var(--color-muted)]">/ {DAILY_QUOTA} {d.quotaUnit}</span>
              </div>
              <p className="text-[12px] text-[var(--color-muted)]">
                {exceeded ? d.quotaExceededHint : d.quotaResetHint}
              </p>
            </div>
            <div className="flex w-full flex-col gap-3 sm:w-1/2 sm:items-end">
              <div className="h-2 w-full overflow-hidden rounded-full bg-[var(--color-surface-2)]">
                <div
                  className="h-full transition-all duration-500"
                  style={{
                    width: `${usagePct}%`,
                    background: exceeded
                      ? "linear-gradient(90deg, var(--color-cr-red), var(--color-cr-red-bright))"
                      : "linear-gradient(90deg, var(--color-cr-blue), var(--color-cr-blue-bright))",
                  }}
                />
              </div>
              <div className="flex items-center gap-2">
                {exceeded ? (
                  <span className="font-mono text-[11px] uppercase tracking-widest text-[var(--color-cr-red-bright)]">
                    {d.quotaExceeded}
                  </span>
                ) : null}
                <button
                  type="button"
                  onClick={handleReset}
                  className="rounded-md border border-[var(--color-border)] bg-[var(--color-surface-2)] px-3 py-1.5 font-mono text-[11px] uppercase tracking-widest text-[var(--color-muted)] transition-colors hover:bg-[var(--color-surface)] hover:text-[var(--color-text)]"
                >
                  {d.reset}
                </button>
              </div>
            </div>
          </div>
        </div>

        {/* ── Cuerpo: casos + chat ─────────────────────────────── */}
        <div className="grid grid-cols-1 gap-5 lg:grid-cols-12">
          {/* Casos */}
          <div className="lg:col-span-5">
            <div className="mb-4 flex flex-col gap-1">
              <div className="text-xs font-mono uppercase tracking-widest text-[var(--color-muted)]">
                {d.casesHeading}
              </div>
              <p className="text-sm text-[var(--color-muted)]">{d.casesSub}</p>
            </div>
            <div className="grid grid-cols-1 gap-px overflow-hidden rounded-xl border border-[var(--color-border)] bg-[var(--color-border)]">
              {DEMO_CASES.map((c) => {
                const loc = localizedCase(c, lang);
                const isActive = activeCaseId === c.id;
                const disabled = running || exceeded;
                return (
                  <button
                    key={c.id}
                    type="button"
                    onClick={() => runCase(c)}
                    disabled={disabled}
                    className="group flex flex-col gap-2 bg-[var(--color-surface)] p-4 text-left transition-colors hover:bg-[var(--color-surface-2)] disabled:cursor-not-allowed disabled:opacity-50 sm:p-5"
                    style={
                      isActive
                        ? { boxShadow: `inset 3px 0 0 ${c.accent}` }
                        : undefined
                    }
                  >
                    <div className="flex items-center justify-between gap-2">
                      <span
                        className="font-mono text-[10px] uppercase tracking-widest"
                        style={{ color: c.accent }}
                      >
                        {loc.cardTag}
                      </span>
                      <span
                        className="h-1.5 w-1.5 rounded-full"
                        style={{ background: c.accent }}
                      />
                    </div>
                    <h3 className="text-base font-semibold leading-snug">
                      {loc.cardTitle}
                    </h3>
                    <p className="text-[13px] leading-relaxed text-[var(--color-muted)]">
                      {loc.cardBody}
                    </p>
                    <div className="mt-1 flex items-center justify-between font-mono text-[11px] text-[var(--color-muted)]">
                      <span className="italic">"{loc.citizenAsk}"</span>
                      <span
                        className="transition-transform group-hover:translate-x-0.5"
                        style={{ color: c.accent }}
                        aria-hidden
                      >
                        →
                      </span>
                    </div>
                  </button>
                );
              })}
            </div>
            <p className="mt-4 text-[11px] leading-relaxed text-[var(--color-muted)]">
              {d.disclaimer}
            </p>
          </div>

          {/* Chat */}
          <div ref={chatPanelRef} className="scroll-mt-20 lg:col-span-7">
            <div className="mb-4 flex flex-col gap-1">
              <div className="text-xs font-mono uppercase tracking-widest text-[var(--color-muted)]">
                {d.chatHeading}
              </div>
              <p className="text-sm text-[var(--color-muted)]">{d.chatSub}</p>
            </div>
            <div className="overflow-hidden rounded-xl border border-[var(--color-border)] bg-[var(--color-surface)]">
              {/* Header chat */}
              <div className="flex items-center justify-between gap-3 border-b border-[var(--color-border)] bg-[var(--color-surface-2)] px-4 py-3">
                <div className="flex items-center gap-2">
                  <span
                    className="inline-flex h-7 w-7 items-center justify-center rounded-full font-mono text-[10px] font-semibold text-white"
                    style={{
                      background:
                        "linear-gradient(120deg, var(--color-cr-blue) 0%, var(--color-cr-blue-bright) 100%)",
                    }}
                  >
                    LLM
                  </span>
                  <div className="flex flex-col leading-tight">
                    <span className="text-sm font-semibold">{d.llmLabel}</span>
                    <span className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted)]">
                      {running
                        ? d.runningLabel
                        : completed
                          ? d.completedLabel
                          : "ready"}
                    </span>
                  </div>
                </div>
                {activeCase ? (
                  <button
                    type="button"
                    disabled={running || exceeded}
                    onClick={() => runCase(activeCase)}
                    className="rounded-md border border-[var(--color-border)] bg-[var(--color-surface)] px-3 py-1.5 font-mono text-[11px] uppercase tracking-widest text-[var(--color-muted)] transition-colors hover:bg-[var(--color-surface-2)] hover:text-[var(--color-text)] disabled:cursor-not-allowed disabled:opacity-50"
                  >
                    {d.runAgain}
                  </button>
                ) : null}
              </div>

              {/* Mensajes */}
              <div
                ref={chatScrollRef}
                className="flex h-[520px] flex-col gap-3 overflow-y-auto p-4 sm:p-5"
              >
                {messages.length === 0 ? (
                  <div className="m-auto flex flex-col items-center gap-2 text-center">
                    <span className="font-mono text-xs uppercase tracking-widest text-[var(--color-muted)]">
                      {d.chatEmpty}
                    </span>
                    <span className="max-w-xs text-sm text-[var(--color-muted)]">
                      {d.chatEmptyHint}
                    </span>
                  </div>
                ) : (
                  messages.map((m) => {
                    if (m.role === "citizen") {
                      return (
                        <div key={m.id} className="flex justify-end">
                          <div className="flex max-w-[85%] flex-col items-end gap-1">
                            <span className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted)]">
                              {d.youLabel}
                            </span>
                            <div
                              className="rounded-2xl rounded-br-sm px-4 py-2.5 text-sm leading-relaxed text-white"
                              style={{
                                background:
                                  "linear-gradient(120deg, var(--color-cr-blue) 0%, var(--color-cr-blue-bright) 100%)",
                              }}
                            >
                              {m.text}
                            </div>
                          </div>
                        </div>
                      );
                    }
                    if (m.role === "llm") {
                      return (
                        <div key={m.id} className="flex justify-start">
                          <div className="flex max-w-[85%] flex-col items-start gap-1">
                            <span className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-cr-blue-bright)]">
                              {d.llmLabel}
                            </span>
                            <div className="rounded-2xl rounded-bl-sm border border-[var(--color-border)] bg-[var(--color-surface-2)] px-4 py-2.5 text-sm leading-relaxed text-[var(--color-text)]">
                              <Typewriter
                                text={m.text}
                                speedMs={14}
                                cursor
                                onTick={scrollToBottom}
                              />
                            </div>
                          </div>
                        </div>
                      );
                    }
                    if (m.role === "system") {
                      const { step } = m;
                      const color = actorBadgeColor(step.actor);
                      const label =
                        step.kind === "thinking"
                          ? d.statusBadges.thinking
                          : step.kind === "request"
                            ? d.statusBadges.request
                            : step.kind === "response"
                              ? d.statusBadges.response
                              : step.kind === "audit"
                                ? d.statusBadges.audit
                                : step.kind === "notify"
                                  ? d.statusBadges.notify
                                  : "";
                      return (
                        <div key={m.id} className="flex justify-start">
                          <div className="flex w-full max-w-[92%] gap-3">
                            <div
                              className="mt-0.5 inline-flex h-7 w-12 shrink-0 items-center justify-center rounded-md font-mono text-[10px] font-semibold uppercase tracking-widest"
                              style={{
                                background: `${color}18`,
                                color,
                                border: `1px solid ${color}40`,
                              }}
                            >
                              {actorTag(step.actor)}
                            </div>
                            <div className="flex flex-1 flex-col gap-1">
                              <div className="flex flex-wrap items-center gap-2">
                                <span className="text-sm font-semibold">
                                  {actorLabel(step.actor, lang)}
                                </span>
                                {label ? (
                                  <StepBadge kind={step.kind} label={label} />
                                ) : null}
                                {step.endpoint ? (
                                  <span className="font-mono text-[10px] text-[var(--color-muted)]">
                                    {step.endpoint}
                                  </span>
                                ) : null}
                              </div>
                              <p className="text-[13px] leading-relaxed text-[var(--color-muted)]">
                                <Typewriter
                                  text={step.text}
                                  speedMs={9}
                                  onTick={scrollToBottom}
                                />
                              </p>
                            </div>
                          </div>
                        </div>
                      );
                    }
                    if (m.role === "email") {
                      return (
                        <div key={m.id} className="flex justify-start">
                          <div className="flex w-full max-w-[92%] flex-col gap-2 rounded-xl border border-[var(--color-border)] bg-[var(--color-bg)] p-4">
                            <div className="flex items-center justify-between">
                              <span
                                className="font-mono text-[10px] uppercase tracking-widest"
                                style={{ color: "var(--color-area-platform)" }}
                              >
                                ✉ {d.emailHeading}
                              </span>
                              <span className="font-mono text-[10px] text-[var(--color-muted)]">
                                {d.emailTo}: {m.to}
                              </span>
                            </div>
                            <h4 className="text-sm font-semibold leading-snug">
                              <Typewriter
                                text={m.subject}
                                speedMs={16}
                                onTick={scrollToBottom}
                              />
                            </h4>
                            <p className="text-[13px] leading-relaxed text-[var(--color-muted)]">
                              <Typewriter
                                text={m.preview}
                                speedMs={9}
                                onTick={scrollToBottom}
                              />
                            </p>
                          </div>
                        </div>
                      );
                    }
                    return null;
                  })
                )}
                {running ? (
                  <div className="flex items-center gap-2 px-1 pt-1">
                    <span
                      className="inline-block h-1.5 w-1.5 animate-bounce rounded-full"
                      style={{ background: "var(--color-cr-blue-bright)", animationDelay: "0ms" }}
                    />
                    <span
                      className="inline-block h-1.5 w-1.5 animate-bounce rounded-full"
                      style={{ background: "var(--color-cr-blue-bright)", animationDelay: "120ms" }}
                    />
                    <span
                      className="inline-block h-1.5 w-1.5 animate-bounce rounded-full"
                      style={{ background: "var(--color-cr-blue-bright)", animationDelay: "240ms" }}
                    />
                  </div>
                ) : null}
              </div>

              {/* Composer (display only) */}
              <div className="border-t border-[var(--color-border)] bg-[var(--color-surface-2)] p-3">
                <div className="flex items-center gap-2 rounded-md border border-[var(--color-border)] bg-[var(--color-surface)] px-3 py-2">
                  <span aria-hidden className="text-[var(--color-muted)]">›</span>
                  <span className="flex-1 truncate text-sm text-[var(--color-muted)]">
                    {d.chatPlaceholder}
                  </span>
                  <span className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted)]">
                    demo
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
