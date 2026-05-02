"use client";

import { LanguageProvider } from "../i18n/LanguageProvider";

export default function Providers({ children }: { children: React.ReactNode }) {
  return <LanguageProvider>{children}</LanguageProvider>;
}
