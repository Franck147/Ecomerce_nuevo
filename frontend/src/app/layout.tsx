"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import clsx from "clsx";
import {
  Store,
  Package,
  Layers,
  Percent,
  Menu,
  X,
  ChevronRight,
} from "lucide-react";
import "./globals.css";

const NAV_ITEMS = [
  {
    href: "/products",
    label: "Productos",
    icon: Package,
  },
  {
    href: "/pricing",
    label: "Precios",
    icon: Percent,
  },
];

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const pathname = usePathname();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  return (
    <html lang="es">
      <body className="min-h-screen bg-slate-50 text-slate-900 font-sans antialiased">
        <div className="flex min-h-screen">
          {/* ========================================================= */}
          {/* BARRA LATERAL IZQUIERDA                                   */}
          {/* ========================================================= */}
          <aside
            className={clsx(
              "fixed inset-y-0 left-0 z-50 w-64 bg-slate-900 text-slate-300 flex flex-col border-r border-slate-800 transition-transform duration-200 ease-in-out lg:translate-x-0 lg:static",
              mobileMenuOpen ? "translate-x-0" : "-translate-x-full"
            )}
          >
            {/* Logo de la Tienda */}
            <div className="h-16 flex items-center justify-between px-6 border-b border-slate-800/80 bg-slate-950/40">
              <Link href="/products" className="flex items-center gap-3">
                <div className="w-8 h-8 rounded-lg bg-indigo-600 flex items-center justify-center text-white font-bold text-base shadow-sm">
                  <Store className="w-4 h-4" />
                </div>
                <span className="font-bold text-base text-white tracking-tight">Ecommerce</span>
              </Link>
              <button
                onClick={() => setMobileMenuOpen(false)}
                className="lg:hidden text-slate-400 hover:text-white"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            {/* Menu de Navegacion (Solo los 2 CRUDs) */}
            <div className="flex-1 px-3 py-6 space-y-1.5 overflow-y-auto">
              <div className="px-3 pb-2 text-[11px] font-semibold tracking-wider text-slate-500 uppercase">
                Administración
              </div>

              {NAV_ITEMS.map((item) => {
                const isActive =
                  pathname === item.href ||
                  (item.href === "/products" && pathname === "/");

                const Icon = item.icon;

                return (
                  <Link
                    key={item.href}
                    href={item.href}
                    onClick={() => setMobileMenuOpen(false)}
                    className={clsx(
                      "group flex items-center px-3.5 py-3 rounded-xl text-sm font-medium transition-all duration-150",
                      isActive
                        ? "bg-indigo-600 text-white shadow-md shadow-indigo-600/30"
                        : "text-slate-400 hover:text-slate-100 hover:bg-slate-800/60"
                    )}
                  >
                    <div className="flex items-center gap-3">
                      <Icon
                        className={clsx(
                          "w-5 h-5 transition-colors",
                          isActive ? "text-white" : "text-slate-400 group-hover:text-slate-200"
                        )}
                      />
                      <span className="font-medium text-sm">{item.label}</span>
                    </div>
                  </Link>
                );
              })}
            </div>

            {/* Pie del Sidebar */}
            <div className="p-4 border-t border-slate-800/80 bg-slate-950/40">
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 rounded-full bg-slate-800 text-slate-300 border border-slate-700 flex items-center justify-center font-bold text-xs">
                  AD
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-xs font-semibold text-white truncate">Administrador</p>
                  <p className="text-[10px] text-slate-500 truncate">admin@ecommerce.local</p>
                </div>
              </div>
            </div>
          </aside>

          {/* ========================================================= */}
          {/* AREA PRINCIPAL                                            */}
          {/* ========================================================= */}
          <div className="flex-1 flex flex-col min-w-0">
            {/* Header Superior */}
            <header className="h-16 bg-white border-b border-slate-200 flex items-center justify-between px-4 sm:px-8 sticky top-0 z-30">
              <div className="flex items-center gap-3">
                <button
                  onClick={() => setMobileMenuOpen(true)}
                  className="lg:hidden p-2 rounded-lg text-slate-600 hover:bg-slate-100"
                >
                  <Menu className="w-5 h-5" />
                </button>
                <div className="flex items-center gap-2 text-xs text-slate-500">
                  <span>Tienda</span>
                  <ChevronRight className="w-3.5 h-3.5 text-slate-400" />
                  <span className="font-semibold text-slate-800">
                    {pathname.includes("pricing") ? "Precios" : "Productos"}
                  </span>
                </div>
              </div>
            </header>

            {/* Contenido */}
            <main className="flex-1 p-4 sm:p-8 max-w-7xl w-full mx-auto">{children}</main>
          </div>
        </div>
      </body>
    </html>
  );
}
