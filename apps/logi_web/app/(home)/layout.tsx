import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "../globals.css";
import Header from "@/components/header";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "LogiApp - Platforma de Logística Inteligente",
  description: "Revoluciona tus operaciones logísticas con nuestra plataforma todo en uno.",
  openGraph: {
    title: "LogiApp - Plataforma de Logística Inteligente",
    description: "Revoluciona tus operaciones logísticas con nuestra plataforma todo en uno.",
    url: "https://logiapp-monorepo-logi-web.vercel.app",
    siteName: "LogiApp",
    images: [
      {
        url: "",
        width: 1200,
        height: 630,
        alt: "LogiApp - Plataforma de Logística Inteligente",
      },
    ],
    locale: "es_ES",
    type: "website",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="es">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <Header />
        {children}
      </body>
    </html>
  );
}
