"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

// Sesión-MVP: la cédula del ciudadano activo se guarda en localStorage.
// En producción esto vendría de un JWT validado por el gateway.

export function useActiveCedula(): { cedula: string | null; loading: boolean } {
  const router = useRouter();
  const [cedula, setCedula] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const c = window.localStorage.getItem("micr_cedula");
    if (!c) {
      router.replace("/");
      return;
    }
    setCedula(c);
    setLoading(false);
  }, [router]);

  return { cedula, loading };
}

export function clearSession() {
  if (typeof window !== "undefined") {
    window.localStorage.removeItem("micr_cedula");
  }
}
