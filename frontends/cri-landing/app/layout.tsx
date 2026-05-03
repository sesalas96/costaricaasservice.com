import type { Metadata, Viewport } from "next";
import "./globals.css";
import Providers from "./components/Providers";

const SITE_URL = "https://costaricaasservice.netlify.app";
const TITLE = "costaricaasservice — Gobierno digital, as a product";
const DESCRIPTION =
  "Propuesta independiente inspirada en e-Estonia. Explora cómo se vería una capa de identidad ciudadana, interoperabilidad entre instituciones y firma digital pensada para que el ciudadano dé un dato una sola vez. Iniciativa privada, sin afiliación oficial.";

export const metadata: Metadata = {
  metadataBase: new URL(SITE_URL),
  title: {
    default: TITLE,
    template: "%s · costaricaasservice",
  },
  description: DESCRIPTION,
  applicationName: "costaricaasservice",
  authors: [{ name: "@devsebas", url: "https://github.com/sesalas96" }],
  creator: "@devsebas",
  publisher: "@devsebas",
  keywords: [
    "Costa Rica",
    "gobierno digital",
    "e-Estonia",
    "X-Road",
    "e-ID",
    "identidad digital",
    "interoperabilidad",
    "firma digital",
    "Conecta",
    "SaaS gobierno",
    "B2G",
  ],
  category: "technology",
  openGraph: {
    type: "website",
    url: SITE_URL,
    siteName: "costaricaasservice",
    title: "Gobierno digital, as a product",
    description: DESCRIPTION,
    locale: "es_CR",
    alternateLocale: ["en_US"],
    images: [
      {
        url: "/opengraph-image",
        width: 1200,
        height: 630,
        alt: "costaricaasservice — Gobierno digital, as a product",
      },
    ],
  },
  twitter: {
    card: "summary_large_image",
    title: TITLE,
    description: DESCRIPTION,
    images: ["/opengraph-image"],
  },
  robots: {
    index: true,
    follow: true,
    googleBot: {
      index: true,
      follow: true,
      "max-image-preview": "large",
      "max-snippet": -1,
    },
  },
  alternates: {
    canonical: SITE_URL,
  },
};

export const viewport: Viewport = {
  themeColor: "#07090d",
  width: "device-width",
  initialScale: 1,
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
