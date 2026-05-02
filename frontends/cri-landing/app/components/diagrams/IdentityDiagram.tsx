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
    color?: string;
    dashed?: boolean;
    animated?: boolean;
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
    stroke: opts.color ?? "var(--color-area-iduc)",
    strokeWidth: opts.animated ? 1.8 : 1.5,
    ...(opts.dashed ? { strokeDasharray: "4 3" } : {}),
  },
  markerEnd: {
    type: MarkerType.ArrowClosed,
    color: opts.color ?? "var(--color-area-iduc)",
  },
  labelBgStyle: { fill: "var(--color-bg)" },
  labelStyle: {
    fill: "var(--color-muted)",
    fontSize: 11,
    fontFamily: "ui-monospace, SFMono-Regular, Menlo, monospace",
  },
});

export default function IdentityDiagram({ audience }: { audience: Audience }) {
  const { t } = useT();

  const technical = () => {
    const i = t.identity;
    // Fila horizontal principal con iduc-identity arriba y iduc-keys abajo
    // como ramas verticales de iduc-signing.
    //              [iduc-identity]
    //                    ↓
    //  [Citizen] → [MiCR] → [iduc-signing] → [doc] → [member]
    //                            ↓
    //                       [iduc-keys]
    const ROW = 240;
    const COLS = [0, 380, 760, 1140, 1520];
    const nodes: Node[] = [
      svc("citizen", COLS[0], ROW, "frontend", i.citizen.title, i.citizen.sub, true),
      svc("micr", COLS[1], ROW, "frontend", i.micr.title, i.micr.sub),
      svc("iduc", COLS[1], 0, "iduc", i.iduc.title, i.iduc.sub, true),
      svc("signing", COLS[2], ROW, "iduc", i.signing.title, i.signing.sub, true),
      svc("keys", COLS[2], 480, "iduc", i.keys.title, i.keys.sub),
      svc("doc", COLS[3], ROW, "platform", i.document.title, i.document.sub, true),
      svc("member", COLS[4], ROW, "members", i.member.title, i.member.sub),
    ];
    const edges: Edge[] = [
      mkEdge("e1", "citizen", "micr", {
        route: HANDLES.fwd,
        label: i.edges.login,
        animated: true,
      }),
      mkEdge("e2", "micr", "iduc", {
        route: HANDLES.up,
        animated: true,
      }),
      mkEdge("e3", "iduc", "micr", {
        route: HANDLES.down,
        label: i.edges.issueJwt,
        dashed: true,
      }),
      mkEdge("e4", "micr", "signing", {
        route: HANDLES.fwd,
        label: i.edges.sign,
        animated: true,
      }),
      mkEdge("e5", "signing", "keys", {
        route: HANDLES.down,
        label: i.edges.askKey,
        animated: true,
        color: "var(--color-area-platform)",
      }),
      mkEdge("e6", "keys", "signing", {
        route: HANDLES.up,
        dashed: true,
        color: "var(--color-area-platform)",
      }),
      mkEdge("e7", "signing", "doc", {
        route: HANDLES.fwd,
        label: i.edges.signedDoc,
        animated: true,
      }),
      mkEdge("e8", "doc", "member", {
        route: HANDLES.fwd,
        label: i.edges.verify,
        color: "var(--color-area-members)",
        animated: true,
      }),
    ];
    return { nodes, edges };
  };

  const citizen = () => {
    const c = t.identityCitizen;
    // 5 nodos en fila horizontal, todo forward, callout arriba.
    //   [Vos] → [Tu cédula digital] → [Caja fuerte] → [Documento firmado] → [Cualquier institución]
    const ROW = 240;
    const COLS = [0, 560, 1120, 1680, 2240];
    const nodes: Node[] = [
      svc("me", COLS[0], ROW, "frontend", c.me.title, c.me.sub, true),
      svc("eid", COLS[1], ROW, "iduc", c.eid.title, c.eid.sub, true),
      svc("safe", COLS[2], ROW, "iduc", c.safe.title, c.safe.sub),
      svc("doc", COLS[3], ROW, "platform", c.document.title, c.document.sub, true),
      svc("member", COLS[4], ROW, "members", c.member.title, c.member.sub),
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
      mkEdge("e1", "me", "eid", {
        route: HANDLES.fwd,
        label: c.edges.login,
        animated: true,
      }),
      mkEdge("e2", "eid", "safe", {
        route: HANDLES.fwd,
        label: c.edges.askSign,
        animated: true,
      }),
      mkEdge("e3", "safe", "doc", {
        route: HANDLES.fwd,
        label: c.edges.signs,
        animated: true,
        color: "var(--color-area-platform)",
      }),
      mkEdge("e4", "doc", "member", {
        route: HANDLES.fwd,
        label: c.edges.anyoneVerifies,
        color: "var(--color-area-members)",
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
      fitView
      fitViewOptions={{ padding: 0.16 }}
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
