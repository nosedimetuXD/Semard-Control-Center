import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "SEMARD Control Center — Plataforma de Gestión",
  description: "Sistema integral de investigación, eventos, préstamos y taller de impresión 3D — Universidad de Cartagena",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="es" className="dark">
      <body className="bg-slate-950 text-slate-100 min-h-screen antialiased">
        {children}
      </body>
    </html>
  );
}
