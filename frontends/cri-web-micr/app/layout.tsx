import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import Link from "next/link";
import "./globals.css";

const geistSans = Geist({ variable: "--font-geist-sans", subsets: ["latin"] });
const geistMono = Geist_Mono({ variable: "--font-geist-mono", subsets: ["latin"] });

export const metadata: Metadata = {
  title: "MiCR — Carpeta Digital Ciudadana",
  description: "Acceso unificado del ciudadano al Estado digital de Costa Rica",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="es" className={`${geistSans.variable} ${geistMono.variable} h-full antialiased`}>
      <body className="min-h-full flex flex-col bg-slate-50 text-slate-900">
        <header className="bg-slate-900 text-white border-b-4 border-blue-600">
          <div className="max-w-6xl mx-auto px-6 py-4 flex items-center justify-between">
            <Link href="/" className="flex items-center gap-3">
              <div className="w-9 h-9 rounded-md bg-blue-600 grid place-items-center font-bold">CR</div>
              <div>
                <div className="text-lg font-semibold">MiCR</div>
                <div className="text-xs text-slate-300">Carpeta Digital Ciudadana</div>
              </div>
            </Link>
            <nav className="flex gap-6 text-sm">
              <Link href="/dashboard" className="hover:text-blue-300">Mi escritorio</Link>
              <Link href="/access-log" className="hover:text-blue-300">Mi bitácora</Link>
              <Link href="/tax" className="hover:text-blue-300">Mi declaración</Link>
            </nav>
          </div>
        </header>
        <main className="flex-1 max-w-6xl mx-auto px-6 py-8 w-full">{children}</main>
        <footer className="border-t border-slate-200 bg-white text-xs text-slate-500">
          <div className="max-w-6xl mx-auto px-6 py-4 flex justify-between">
            <span>costaricaasservice · realm <code className="font-mono">demo</code></span>
            <span>Modelo e-Estonia · Conecta CR</span>
          </div>
        </footer>
      </body>
    </html>
  );
}
