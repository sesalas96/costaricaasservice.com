import type { Edge, Node } from "@xyflow/react";

const NODE_STAGGER_MS = 55;
const EDGE_BASE_DELAY_MS = 480;
const EDGE_STAGGER_MS = 35;

// Aplica delays de animación a los nodos (cascada de aparición) y a los edges
// (entran después que los nodos). xyflow expone `style` sobre el wrapper del
// nodo y sobre el path del edge — usamos `animationDelay` en ambos casos.
export function withDiagramAnimations<N extends Node, E extends Edge>(
  nodes: N[],
  edges: E[],
): { nodes: N[]; edges: E[] } {
  const animatedNodes = nodes.map((n, i) => ({
    ...n,
    style: {
      ...(n.style ?? {}),
      animationDelay: `${i * NODE_STAGGER_MS}ms`,
    },
  }));

  const animatedEdges = edges.map((e, i) => ({
    ...e,
    style: {
      ...(e.style ?? {}),
      animationDelay: `${EDGE_BASE_DELAY_MS + i * EDGE_STAGGER_MS}ms`,
    },
  }));

  return { nodes: animatedNodes, edges: animatedEdges };
}
