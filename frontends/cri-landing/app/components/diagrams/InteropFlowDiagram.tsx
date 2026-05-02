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
    dashed?: boolean;
    animated?: boolean;
    color?: string;
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

export default function InteropFlowDiagram({
  audience,
}: {
  audience: Audience;
}) {
  const { t } = useT();

  const technical = () => {
    const i = t.interop;
    // 5 nodos en fila (forward), audit arriba, hub abajo.
    // Ida arriba (right→left, animated). Vuelta por debajo (bottom→bottom dashed).
    const ROW = 240;
    const COLS = [0, 320, 640, 960, 1280];
    const nodes: Node[] = [
      svc("src", COLS[0], ROW, "members", i.haciendaSvc.title, i.haciendaSvc.sub, true),
      svc("ss-src", COLS[1], ROW, "interop", i.ssHac.title, i.ssHac.sub, true),
      svc("router", COLS[2], ROW, "interop", i.router.title, i.router.sub, true),
      svc("ss-dst", COLS[3], ROW, "interop", i.ssRc.title, i.ssRc.sub, true),
      svc("dst", COLS[4], ROW, "members", i.rcSvc.title, i.rcSvc.sub, true),
      svc("audit", COLS[2], 0, "interop", i.audit.title, i.audit.sub),
      svc("hub", COLS[2], 480, "interop", i.hub.title, i.hub.sub),
    ];
    const edges: Edge[] = [
      // ida — fila horizontal animada
      mkEdge("f1", "src", "ss-src", {
        route: HANDLES.fwd,
        label: i.edges.request,
        animated: true,
      }),
      mkEdge("f2", "ss-src", "router", {
        route: HANDLES.fwd,
        label: i.edges.jwsSigned,
        animated: true,
      }),
      mkEdge("f3", "router", "ss-dst", {
        route: HANDLES.fwd,
        label: i.edges.routeMtls,
        animated: true,
      }),
      mkEdge("f4", "ss-dst", "dst", {
        route: HANDLES.fwd,
        label: i.edges.verifyDeliver,
        animated: true,
      }),
      // vuelta — un solo arco por debajo, de dst a src
      mkEdge("r-back", "dst", "src", {
        route: HANDLES.arcBelow,
        label: i.edges.response,
        dashed: true,
        color: "var(--color-muted)",
      }),
      // audit — sale del router hacia arriba
      mkEdge("a1", "router", "audit", {
        route: HANDLES.up,
        label: i.edges.kafkaEvent,
        animated: true,
        color: "var(--color-area-platform)",
      }),
      // hub — alimenta a ambos security-servers
      mkEdge("h1", "hub", "ss-src", {
        route: { sourceHandle: "l-s", targetHandle: "b-t" },
        label: i.edges.catalogCa,
        dashed: true,
        color: "var(--color-muted)",
      }),
      mkEdge("h2", "hub", "ss-dst", {
        route: { sourceHandle: "r-s", targetHandle: "b-t" },
        dashed: true,
        color: "var(--color-muted)",
      }),
    ];
    return { nodes, edges };
  };

  const citizen = () => {
    const c = t.interopCitizen;
    // 3 nodos horizontales + audit abajo. Sin return arrows.
    //   [Inst A] → [Conecta] → [Inst B]
    //                  ↓
    //              [Tu bitácora]
    const ROW = 200;
    const COLS = [0, 380, 760];
    const nodes: Node[] = [
      svc("a", COLS[0], ROW, "members", c.institutionA.title, c.institutionA.sub, true),
      svc("bridge", COLS[1], ROW, "interop", c.bridge.title, c.bridge.sub, true),
      svc("b", COLS[2], ROW, "members", c.institutionB.title, c.institutionB.sub, true),
      svc("audit", COLS[1], 460, "platform", c.audit.title, c.audit.sub),
    ];
    const edges: Edge[] = [
      mkEdge("e1", "a", "bridge", {
        route: HANDLES.fwd,
        label: c.edges.needs,
        animated: true,
      }),
      mkEdge("e2", "bridge", "b", {
        route: HANDLES.fwd,
        label: c.edges.sealed,
        animated: true,
      }),
      // arco de respuesta por debajo (dashed, sin label para no chocar)
      mkEdge("r1", "b", "a", {
        route: HANDLES.arcBelow,
        label: c.edges.response,
        dashed: true,
        color: "var(--color-muted)",
      }),
      // audit hacia abajo desde el puente
      mkEdge("au", "bridge", "audit", {
        route: HANDLES.down,
        label: c.edges.logged,
        animated: true,
        color: "var(--color-area-platform)",
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
      fitView
      fitViewOptions={{ padding: 0.18 }}
      minZoom={0.35}
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
