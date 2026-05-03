"use client";

import {
  createContext,
  useCallback,
  useContext,
  useState,
  type ReactNode,
} from "react";
import type { InstitutionId } from "../data/institutions";

export type InstitutionTarget = InstitutionId | "__list__";

type Selected =
  | { type: "list" }
  | { type: "detail"; id: InstitutionId; fromList: boolean }
  | null;

type Ctx = {
  selected: Selected;
  open: (target: InstitutionTarget) => void;
  openFromList: (id: InstitutionId) => void;
  goBackToList: () => void;
  close: () => void;
};

const InstitutionContext = createContext<Ctx>({
  selected: null,
  open: () => {},
  openFromList: () => {},
  goBackToList: () => {},
  close: () => {},
});

export function InstitutionProvider({ children }: { children: ReactNode }) {
  const [selected, setSelected] = useState<Selected>(null);

  const open = useCallback((target: InstitutionTarget) => {
    if (target === "__list__") {
      setSelected({ type: "list" });
    } else {
      setSelected({ type: "detail", id: target, fromList: false });
    }
  }, []);

  const openFromList = useCallback((id: InstitutionId) => {
    setSelected({ type: "detail", id, fromList: true });
  }, []);

  const goBackToList = useCallback(() => setSelected({ type: "list" }), []);
  const close = useCallback(() => setSelected(null), []);

  return (
    <InstitutionContext.Provider
      value={{ selected, open, openFromList, goBackToList, close }}
    >
      {children}
    </InstitutionContext.Provider>
  );
}

export function useInstitution() {
  return useContext(InstitutionContext);
}
