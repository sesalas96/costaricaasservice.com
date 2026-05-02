"use client";

import { Handle, Position, type Node, type NodeProps } from "@xyflow/react";

export type AreaKey =
  | "frontend"
  | "gateway"
  | "iduc"
  | "interop"
  | "members"
  | "platform"
  | "external";

const AREA_COLOR: Record<AreaKey, string> = {
  frontend: "var(--color-area-frontend)",
  gateway: "var(--color-area-gateway)",
  iduc: "var(--color-area-iduc)",
  interop: "var(--color-area-interop)",
  members: "var(--color-area-members)",
  platform: "var(--color-area-platform)",
  external: "var(--color-muted)",
};

const AREA_LABEL: Record<AreaKey, string> = {
  frontend: "frontend",
  gateway: "gateway",
  iduc: "iduc",
  interop: "interop",
  members: "members",
  platform: "platform",
  external: "externo",
};

export type ServiceNode = Node<
  {
    title: string;
    subtitle?: string;
    area: AreaKey;
    emphasis?: boolean;
  },
  "service"
>;

// Cada lado expone source y target (handles superpuestos invisibles).
// IDs convenidos: `${side}-s` (source), `${side}-t` (target). Side ∈ {t,r,b,l}.
const SIDES = [
  { side: "t", pos: Position.Top },
  { side: "r", pos: Position.Right },
  { side: "b", pos: Position.Bottom },
  { side: "l", pos: Position.Left },
] as const;

const baseHandleStyle = (color: string) => ({
  background: color,
  border: "none",
  width: 6,
  height: 6,
});

export function ServiceNodeView({ data }: NodeProps<ServiceNode>) {
  const color = AREA_COLOR[data.area];
  const dotStyle = baseHandleStyle(color);
  return (
    <div
      className="rounded-lg border bg-[var(--color-surface)] px-4 py-3 shadow-sm transition-transform"
      style={{
        borderColor: data.emphasis ? color : "var(--color-border)",
        boxShadow: data.emphasis ? `0 0 0 1px ${color} inset` : undefined,
        minWidth: 180,
      }}
    >
      {SIDES.map(({ side, pos }) => (
        <span key={side}>
          {/* dot visible (source) */}
          <Handle
            id={`${side}-s`}
            type="source"
            position={pos}
            style={dotStyle}
          />
          {/* target invisible superpuesto */}
          <Handle
            id={`${side}-t`}
            type="target"
            position={pos}
            style={{ ...dotStyle, opacity: 0 }}
          />
        </span>
      ))}
      <div className="flex items-center gap-2">
        <span
          className="h-1.5 w-1.5 rounded-full"
          style={{ background: color }}
        />
        <span className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted)]">
          {AREA_LABEL[data.area]}
        </span>
      </div>
      <div className="mt-1 text-[13px] font-semibold text-[var(--color-text)]">
        {data.title}
      </div>
      {data.subtitle && (
        <div className="mt-0.5 text-[11px] leading-snug text-[var(--color-muted)]">
          {data.subtitle}
        </div>
      )}
    </div>
  );
}

export type GroupNode = Node<
  {
    title: string;
    area: AreaKey;
  },
  "group"
>;

export function GroupNodeView({ data }: NodeProps<GroupNode>) {
  const color = AREA_COLOR[data.area];
  return (
    <div
      className="h-full w-full rounded-xl border border-dashed bg-transparent"
      style={{ borderColor: color, opacity: 0.55 }}
    >
      <div
        className="absolute -top-3 left-3 rounded bg-[var(--color-bg)] px-2 font-mono text-[10px] uppercase tracking-widest"
        style={{ color }}
      >
        {data.title}
      </div>
    </div>
  );
}

export const nodeTypes = {
  service: ServiceNodeView,
  group: GroupNodeView,
};
