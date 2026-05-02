import type { Metadata } from "next";
import "./globals.css";
import Providers from "./components/Providers";

export const metadata: Metadata = {
  title: "costaricaasservice — Núcleo de Conecta CR",
  description:
    "Plataforma SaaS B2G/B2B inspirada en e-Estonia. X-Road + e-ID + portal ciudadano como producto multi-tenant.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="es">
      <body>
        <div className="cr-flag-bar" aria-hidden />
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
