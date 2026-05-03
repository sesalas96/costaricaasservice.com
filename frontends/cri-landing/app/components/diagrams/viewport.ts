"use client";

import type { ReactFlowInstance } from "@xyflow/react";

export const DIAGRAM_MOBILE_MAX = 768;

export function isMobileViewport(): boolean {
  if (typeof window === "undefined") return false;
  return window.matchMedia(`(max-width: ${DIAGRAM_MOBILE_MAX}px)`).matches;
}

// En mobile, en vez de centrar el diagrama (queda diminuto), anclamos el
// nodo más a la izquierda al borde izquierdo con un zoom legible. El usuario
// arrastra horizontalmente para seguir leyendo el flujo.
export function alignFirstNodeLeft(
  instance: ReactFlowInstance,
  opts: { zoom?: number; padLeft?: number; padTop?: number } = {},
) {
  const { zoom = 0.65, padLeft = 16, padTop = 24 } = opts;
  const nodes = instance.getNodes();
  if (!nodes.length) return;
  const minX = Math.min(...nodes.map((n) => n.position.x));
  const minY = Math.min(...nodes.map((n) => n.position.y));
  instance.setViewport({
    x: padLeft - minX * zoom,
    y: padTop - minY * zoom,
    zoom,
  });
}
