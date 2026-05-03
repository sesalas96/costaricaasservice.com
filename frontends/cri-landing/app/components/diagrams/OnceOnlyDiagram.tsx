"use client";

import {
  Background,
  BackgroundVariant,
  Controls,
  MarkerType,
  ReactFlow,
  type Edge,
  type Node,
} from "@xyflow/react";
import "@xyflow/react/dist/style.css";

import { useT } from "../../i18n/LanguageProvider";
import { nodeTypes, type AreaKey } from "./nodes";
import { HANDLES, type Audience } from "./types";
import { alignFirstNodeLeft, isMobileViewport } from "./viewport";

type Route = { sourceHandle: string; targetHandle: string };

const svc = (
  id: string,
  x: number,
  y: number,
  area: AreaKey,
  title: string,
  subtitle?: string,
  emphasis?: boolean,
): Node => ({
  id,
  type: "service",
  position: { x, y },
  data: { title, subtitle, area, emphasis },
});

const mkEdge = (
  id: string,
  source: string,
  target: string,
  opts: {
    route?: Route;
    label?: string;
    color?: string;
    animated?: boolean;
    dashed?: boolean;
  } = {},
): Edge => ({
  id,
  source,
  target,
  ...(opts.route
    ? {
        sourceHandle: opts.route.sourceHandle,
        targetHandle: opts.route.targetHandle,
      }
    : {}),
  type: "smoothstep",
  label: opts.label,
  className: opts.animated ? "flow-pulse" : undefined,
  style: {
    stroke: opts.color ?? "var(--color-area-interop)",
    strokeWidth: opts.animated ? 1.8 : 1.5,
    ...(opts.dashed ? { strokeDasharray: "4 3" } : {}),
  },
  markerEnd: {
    type: MarkerType.ArrowClosed,
    color: opts.color ?? "var(--color-area-interop)",
  },
  labelBgStyle: { fill: "var(--color-bg)" },
  labelStyle: {
    fill: "var(--color-muted)",
    fontSize: 11,
    fontFamily: "ui-monospace, SFMono-Regular, Menlo, monospace",
  },
});

export default function OnceOnlyDiagram({
  audience,
}: {
  audience: Audience;
}) {
  const { t } = useT();
  const GREEN = "var(--color-area-iduc)";
  const AMBER = "var(--color-area-members)";

  const technical = () => {
    const o = t.onceOnly;
    // Fila top forward; X-Road abajo; retorno como arco único por debajo.
    const TOP_Y = 220;
    const BOTTOM_Y = 460;
    const COLS = [0, 380, 760, 1140, 1520, 1900];
    const nodes: Node[] = [
      svc("citizen", COLS[0], TOP_Y, "frontend", o.citizen.title, o.citizen.sub, true),
      svc("micr", COLS[1], TOP_Y, "frontend", o.micr.title, o.micr.sub),
      svc("bff", COLS[2], TOP_Y, "gateway", o.bff.title, o.bff.sub),
      svc("hac", COLS[3], TOP_Y, "members", o.hacienda.title, o.hacienda.sub, true),
      svc("rc", COLS[5], TOP_Y, "members", o.rc.title, o.rc.sub, true),
      svc("ss-hac", COLS[3], BOTTOM_Y, "interop", o.ssHac.title),
      svc("router", COLS[4], BOTTOM_Y, "interop", o.router.title),
      svc("ss-rc", COLS[5], BOTTOM_Y, "interop", o.ssRc.title),
      {
        id: "callout",
        type: "service",
        position: { x: COLS[1], y: 0 },
        data: {
          title: o.callout.title,
          subtitle: o.callout.sub,
          area: "iduc",
          emphasis: true,
        },
      },
    ];
    const edges: Edge[] = [
      // ida del ciudadano (1, 2, 3) — forward horizontal limpio
      mkEdge("u1", "citizen", "micr", {
        route: HANDLES.fwd,
        label: o.edges.request,
        color: GREEN,
        animated: true,
      }),
      mkEdge("u2", "micr", "bff", {
        route: HANDLES.fwd,
        label: o.edges.step2,
        color: GREEN,
        animated: true,
      }),
      mkEdge("u3", "bff", "hac", {
        route: HANDLES.fwd,
        label: o.edges.getPrefilled,
        color: GREEN,
        animated: true,
      }),
      // X-Road (4 → bottom → 5)
      mkEdge("x1", "hac", "ss-hac", {
        route: HANDLES.down,
        label: o.edges.needsPerson,
        animated: true,
      }),
      mkEdge("x2", "ss-hac", "router", { route: HANDLES.fwd, animated: true }),
      mkEdge("x3", "router", "ss-rc", { route: HANDLES.fwd, animated: true }),
      mkEdge("x4", "ss-rc", "rc", {
        route: HANDLES.up,
        label: o.edges.getPerson,
        animated: true,
      }),
      // datos del ciudadano vuelven a Hacienda — arco por debajo (no choca con la fila top)
      mkEdge("d1", "rc", "hac", {
        route: HANDLES.arcBelow,
        label: o.edges.citizenData,
        color: AMBER,
        animated: true,
      }),
      // declaración pre-llenada vuelve al ciudadano — arco grande por debajo
      mkEdge("ret", "hac", "citizen", {
        route: HANDLES.arcBelow,
        label: `${o.edges.prefilled} → ${o.edges.signClick}`,
        color: GREEN,
        animated: true,
      }),
    ];
    return { nodes, edges };
  };

  const citizen = () => {
    const c = t.onceOnlyCitizen;
    // 5 nodos en fila horizontal + arco único de retorno por debajo + callout arriba.
    //   [Vos] → [Tu app] → [Hacienda] → [Conecta] → [Reg.Civil]
    //                                                       │
    //                          ←─── te llega lista ─────────┘ (arco bottom)
    const ROW = 220;
    const COLS = [0, 620, 1240, 1860, 2480];
    const nodes: Node[] = [
      svc("me", COLS[0], ROW, "frontend", c.me.title, c.me.sub, true),
      svc("micr", COLS[1], ROW, "frontend", c.micr.title, c.micr.sub),
      svc("hac", COLS[2], ROW, "members", c.hacienda.title, c.hacienda.sub, true),
      svc("bridge", COLS[3], ROW, "interop", c.bridge.title, c.bridge.sub),
      svc("rc", COLS[4], ROW, "members", c.rc.title, c.rc.sub, true),
      {
        id: "callout",
        type: "service",
        position: { x: COLS[1], y: 0 },
        data: {
          title: c.callout.title,
          subtitle: c.callout.sub,
          area: "iduc",
          emphasis: true,
        },
      },
    ];
    const edges: Edge[] = [
      // 4 pasos forward, todos en la fila horizontal
      mkEdge("u1", "me", "micr", {
        route: HANDLES.fwd,
        label: c.edges.request,
        color: GREEN,
        animated: true,
      }),
      mkEdge("u2", "micr", "hac", {
        route: HANDLES.fwd,
        label: c.edges.micrAsks,
        color: GREEN,
        animated: true,
      }),
      mkEdge("u3", "hac", "bridge", {
        route: HANDLES.fwd,
        label: c.edges.hacAsksRc,
        animated: true,
      }),
      mkEdge("u4", "bridge", "rc", {
        route: HANDLES.fwd,
        label: c.edges.rcAnswers,
        color: AMBER,
        animated: true,
      }),
      // arco único de retorno por debajo: Reg.Civil → Vos
      // (lleva el resultado de vuelta al ciudadano sin chocar con la fila top)
      mkEdge("back", "rc", "me", {
        route: HANDLES.arcBelow,
        label: `${c.edges.readyToSign} → ${c.edges.signed}`,
        color: GREEN,
        animated: true,
      }),
    ];
    return { nodes, edges };
  };

  const { nodes, edges } = audience === "technical" ? technical() : citizen();

  return (
    <ReactFlow
      nodes={nodes}
      edges={edges}
      nodeTypes={nodeTypes}
      onInit={(instance) => {
        if (isMobileViewport()) {
          alignFirstNodeLeft(instance);
        } else {
          instance.fitView({ padding: 0.16 });
        }
      }}
      minZoom={0.3}
      maxZoom={1.4}
      proOptions={{ hideAttribution: true }}
      nodesDraggable={false}
      nodesConnectable={false}
      elementsSelectable={false}
    >
      <Background
        variant={BackgroundVariant.Dots}
        gap={20}
        size={1}
        color="#1a1f2b"
      />
      <Controls showInteractive={false} />
    </ReactFlow>
  );
}
