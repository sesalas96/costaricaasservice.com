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

const svc = (
  id: string,
  x: number,
  y: number,
  area: AreaKey,
  title: string,
  subtitle?: string,
  parent?: string,
  emphasis?: boolean,
): Node => ({
  id,
  type: "service",
  position: { x, y },
  data: { title, subtitle, area, emphasis },
  ...(parent
    ? { parentId: parent, extent: "parent" as const, draggable: false }
    : {}),
});

const group = (
  id: string,
  x: number,
  y: number,
  w: number,
  h: number,
  area: AreaKey,
  title: string,
): Node => ({
  id,
  type: "group",
  position: { x, y },
  data: { title, area },
  style: { width: w, height: h, background: "transparent", border: "none" },
  selectable: false,
  draggable: false,
});

type Route = { sourceHandle: string; targetHandle: string };

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
    stroke: opts.color ?? "var(--color-border)",
    strokeWidth: opts.animated ? 1.8 : 1.4,
    ...(opts.dashed ? { strokeDasharray: "4 3" } : {}),
  },
  markerEnd: {
    type: MarkerType.ArrowClosed,
    color: opts.color ?? "var(--color-muted)",
  },
  labelBgStyle: { fill: "var(--color-bg)" },
  labelStyle: {
    fill: "var(--color-muted)",
    fontSize: 11,
    fontFamily: "ui-monospace, SFMono-Regular, Menlo, monospace",
  },
});

export default function EcosystemDiagram({
  audience,
}: {
  audience: Audience;
}) {
  const { t } = useT();

  const technical = () => {
    const e = t.ecosystem;
    const GROUP_Y = 460;
    const GROUP_H = 560;
    const COL_W = 290;
    const COL_GAP = 80;
    const COL_X = (i: number) => i * (COL_W + COL_GAP);
    const CHILD_X = 18;
    const slot = (i: number) => 60 + i * 135;

    const nodes: Node[] = [
      svc("cli-citizen", 90, 0, "frontend", e.citizen.title, e.citizen.sub),
      svc("cli-operator", 600, 0, "frontend", e.operator.title, e.operator.sub),
      svc("cli-admin", 1110, 0, "frontend", e.admin.title, e.admin.sub),
      svc("gateway", 600, 140, "gateway", e.gateway.title, e.gateway.sub),
      svc("bff-citizen", 90, 280, "gateway", e.bffCitizen.title, e.bffCitizen.sub),
      svc("bff-member", 600, 280, "gateway", e.bffMember.title, e.bffMember.sub),
      svc("bff-admin", 1110, 280, "gateway", e.bffAdmin.title, e.bffAdmin.sub),

      group("g-iduc", COL_X(0), GROUP_Y, COL_W, GROUP_H, "iduc", e.iducGroup),
      svc("iduc-identity", CHILD_X, slot(0), "iduc", e.iducIdentity.title, e.iducIdentity.sub, "g-iduc"),
      svc("iduc-keys", CHILD_X, slot(1), "iduc", e.iducKeys.title, e.iducKeys.sub, "g-iduc"),
      svc("iduc-signing", CHILD_X, slot(2), "iduc", e.iducSigning.title, e.iducSigning.sub, "g-iduc"),

      group("g-interop", COL_X(1), GROUP_Y, COL_W, GROUP_H, "interop", e.interopGroup),
      svc("io-hub", CHILD_X, slot(0), "interop", e.interopHub.title, e.interopHub.sub, "g-interop"),
      svc("io-ss", CHILD_X, slot(1), "interop", e.securityServer.title, e.securityServer.sub, "g-interop"),
      svc("io-router", CHILD_X, slot(2), "interop", e.interopRouter.title, e.interopRouter.sub, "g-interop"),
      svc("io-audit", CHILD_X, slot(3), "interop", e.interopAudit.title, e.interopAudit.sub, "g-interop"),

      group("g-members", COL_X(2), GROUP_Y, COL_W, GROUP_H, "members", e.membersGroup),
      svc("m-rc", CHILD_X, slot(0), "members", e.registroCivil.title, e.registroCivil.sub, "g-members"),
      svc("m-hac", CHILD_X, slot(1), "members", e.hacienda.title, e.hacienda.sub, "g-members"),
      svc("m-future", CHILD_X, slot(2), "external", e.otherMembers.title, e.otherMembers.sub, "g-members"),

      group("g-platform", COL_X(3), GROUP_Y, COL_W, GROUP_H, "platform", e.platformGroup),
      svc("p-audit", CHILD_X, slot(0), "platform", e.platformAudit.title, e.platformAudit.sub, "g-platform"),
      svc("p-notif", CHILD_X, slot(1), "platform", e.notifications.title, e.notifications.sub, "g-platform"),
      svc("p-files", CHILD_X, slot(2), "platform", e.files.title, e.files.sub, "g-platform"),
      svc("p-flags", CHILD_X, slot(3), "platform", e.featureFlags.title, e.featureFlags.sub, "g-platform"),
    ];

    const edges: Edge[] = [
      // clientes → gateway: vertical down
      mkEdge("c1", "cli-citizen", "gateway", {
        route: HANDLES.down,
        animated: true,
        color: "var(--color-cr-blue-bright)",
      }),
      mkEdge("c2", "cli-operator", "gateway", { route: HANDLES.down }),
      mkEdge("c3", "cli-admin", "gateway", { route: HANDLES.down }),
      // gateway → BFFs: down
      mkEdge("g1", "gateway", "bff-citizen", {
        route: HANDLES.down,
        animated: true,
        color: "var(--color-cr-blue-bright)",
      }),
      mkEdge("g2", "gateway", "bff-member", { route: HANDLES.down }),
      mkEdge("g3", "gateway", "bff-admin", { route: HANDLES.down }),
      // BFF citizen → servicios (down hacia las columnas)
      mkEdge("bc1", "bff-citizen", "iduc-identity", {
        route: HANDLES.down,
        animated: true,
        color: "var(--color-area-iduc)",
      }),
      mkEdge("bc2", "bff-citizen", "p-audit", { route: HANDLES.down }),
      mkEdge("bc3", "bff-citizen", "iduc-signing", { route: HANDLES.down }),
      mkEdge("bm1", "bff-member", "io-hub", { route: HANDLES.down }),
      mkEdge("bm2", "bff-member", "m-hac", { route: HANDLES.down }),
      mkEdge("ba1", "bff-admin", "io-hub", { route: HANDLES.down }),
      mkEdge("ba2", "bff-admin", "p-flags", { route: HANDLES.down }),
      // dependencias internas (intra-columna y entre columnas)
      mkEdge("i1", "iduc-signing", "iduc-keys", {
        route: HANDLES.up,
        label: e.edges.kms,
        dashed: true,
      }),
      // m-rc/m-hac están a la derecha de io-ss, así que entran por la derecha
      mkEdge("i2", "m-rc", "io-ss", {
        route: HANDLES.back,
        label: e.edges.expose,
        dashed: true,
      }),
      mkEdge("i3", "m-hac", "io-ss", {
        route: HANDLES.back,
        label: e.edges.exposeConsume,
        dashed: true,
        animated: true,
        color: "var(--color-area-interop)",
      }),
      mkEdge("i4", "io-ss", "io-router", {
        route: HANDLES.down,
        label: e.edges.mtls,
        animated: true,
        color: "var(--color-area-interop)",
      }),
      mkEdge("i5", "io-router", "io-audit", {
        route: HANDLES.down,
        label: e.edges.event,
        animated: true,
        color: "var(--color-area-platform)",
      }),
    ];
    return { nodes, edges };
  };

  const citizen = () => {
    const c = t.ecosystemCitizen;
    // Layout horizontal: 4 nodos en fila + 2 ramas hacia abajo.
    //   [Vos] → [Tu app] → [Conecta] → [Instituciones]
    //              ↓             ↓
    //         [Tu firma]    [Tu bitácora]
    const ROW = 200;
    const BRANCH = 460;
    const COLS = [0, 600, 1200, 1800];
    const nodes: Node[] = [
      svc("me", COLS[0], ROW, "frontend", c.me.title, c.me.sub, undefined, true),
      svc("micr", COLS[1], ROW, "frontend", c.micr.title, c.micr.sub),
      svc("conecta", COLS[2], ROW, "interop", c.conecta.title, c.conecta.sub, undefined, true),
      svc("institutions", COLS[3], ROW, "members", c.institutions.title, c.institutions.sub, undefined, true),
      svc("signing", COLS[1], BRANCH, "iduc", c.signing.title, c.signing.sub),
      svc("audit", COLS[2], BRANCH, "platform", c.audit.title, c.audit.sub),
    ];
    const edges: Edge[] = [
      // fila principal — 3 pasos horizontales
      mkEdge("e1", "me", "micr", {
        route: HANDLES.fwd,
        label: c.edges.ingreso,
        animated: true,
        color: "var(--color-cr-blue-bright)",
      }),
      mkEdge("e2", "micr", "conecta", {
        route: HANDLES.fwd,
        label: c.edges.consulta,
        animated: true,
        color: "var(--color-cr-blue-bright)",
      }),
      mkEdge("e3", "conecta", "institutions", {
        route: HANDLES.fwd,
        label: c.edges.enruta,
        animated: true,
        color: "var(--color-area-interop)",
      }),
      // ramas hacia abajo
      mkEdge("e4", "micr", "signing", {
        route: HANDLES.down,
        label: c.edges.firmo,
        animated: true,
        color: "var(--color-area-iduc)",
      }),
      mkEdge("e5", "conecta", "audit", {
        route: HANDLES.down,
        label: c.edges.registra,
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
      fitViewOptions={{ padding: 0.12 }}
      minZoom={0.25}
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
