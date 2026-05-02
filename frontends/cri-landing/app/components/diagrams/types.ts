export type Audience = "citizen" | "technical";

// Convenios de ruteo entre handles. Los nodos exponen 4 lados, cada uno con
// `${side}-s` (source) y `${side}-t` (target).
export const HANDLES = {
  // forward: derecha → izquierda (flujo natural izq→der)
  fwd: { sourceHandle: "r-s", targetHandle: "l-t" },
  // back: izquierda → derecha (vuelta hacia atrás)
  back: { sourceHandle: "l-s", targetHandle: "r-t" },
  // down: bottom de origen → top de destino (rama hacia abajo)
  down: { sourceHandle: "b-s", targetHandle: "t-t" },
  // up: top de origen → bottom de destino (rama hacia arriba)
  up: { sourceHandle: "t-s", targetHandle: "b-t" },
  // loop por debajo (ambas patas en el bottom — para arcos de retorno)
  arcBelow: { sourceHandle: "b-s", targetHandle: "b-t" },
  // loop por arriba
  arcAbove: { sourceHandle: "t-s", targetHandle: "t-t" },
} as const;
