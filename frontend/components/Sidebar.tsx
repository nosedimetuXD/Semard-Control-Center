"use client";

import {
  LayoutDashboard,
  FolderGit2,
  Printer,
  Package,
  Calendar,
  Users,
  User,
  Cpu,
  Sparkles,
} from "lucide-react";

export type ViewKey =
  | "dashboard"
  | "projects"
  | "print3d"
  | "inventory"
  | "events"
  | "members"
  | "my-portal";

interface SidebarProps {
  currentView: ViewKey;
  onSelectView: (view: ViewKey) => void;
  user: any;
}

export default function Sidebar({ currentView, onSelectView, user }: SidebarProps) {
  const menuItems = [
    { key: "dashboard", label: "Dashboard General", icon: LayoutDashboard, badge: null },
    { key: "projects", label: "Proyectos de Investigación", icon: FolderGit2, badge: "6 activos" },
    { key: "print3d", label: "Taller de Impresión 3D", icon: Printer, badge: "Cola activa" },
    { key: "inventory", label: "Laboratorio & Equipos", icon: Package, badge: null },
    { key: "events", label: "Eventos & Convocatorias", icon: Calendar, badge: "Nuevo" },
    { key: "members", label: "Directorio & Miembros", icon: Users, badge: null },
    { key: "my-portal", label: "Mi Portal de Miembro", icon: User, badge: user ? "Activo" : null },
  ];

  return (
    <aside className="w-full lg:w-64 bg-slate-900/90 border-r border-slate-800 flex flex-col shrink-0 lg:min-h-screen">
      {/* Brand Header */}
      <div className="p-5 border-b border-slate-800 flex items-center gap-3">
        <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-teal-500 to-emerald-700 flex items-center justify-center text-white shadow-lg shadow-teal-500/10 border border-teal-400/30">
          <Cpu className="w-5 h-5 text-white" />
        </div>
        <div>
          <div className="flex items-center gap-1.5">
            <h1 className="text-sm font-bold tracking-tight text-white">SEMARD</h1>
            <span className="text-[10px] font-semibold px-1.5 py-0.2 rounded bg-teal-500/20 text-teal-300 border border-teal-500/30">
              UdeC
            </span>
          </div>
          <p className="text-[11px] text-slate-400">Control Center • Robótica</p>
        </div>
      </div>

      {/* Navigation */}
      <div className="p-3 flex-1 flex flex-row lg:flex-col gap-1 overflow-x-auto lg:overflow-visible">
        <div className="hidden lg:block px-3 py-2 text-[10px] font-bold uppercase tracking-wider text-slate-500">
          Módulos del Sistema
        </div>

        {menuItems.map((item) => {
          const Icon = item.icon;
          const active = currentView === item.key;
          return (
            <button
              key={item.key}
              onClick={() => onSelectView(item.key as ViewKey)}
              className={`w-full text-left px-3 py-2.5 rounded-lg font-medium transition flex items-center justify-between text-xs whitespace-nowrap ${
                active
                  ? "bg-teal-600/90 text-white shadow-sm shadow-teal-500/20 font-semibold"
                  : "text-slate-400 hover:bg-slate-800/80 hover:text-slate-200"
              }`}
            >
              <div className="flex items-center gap-2.5">
                <Icon className={`w-4 h-4 ${active ? "text-white" : "text-slate-400"}`} />
                <span>{item.label}</span>
              </div>
              {item.badge && (
                <span
                  className={`hidden lg:inline-block text-[10px] px-1.5 py-0.5 rounded-full font-semibold ${
                    active
                      ? "bg-white/20 text-white"
                      : "bg-slate-800 text-teal-400 border border-slate-700"
                  }`}
                >
                  {item.badge}
                </span>
              )}
            </button>
          );
        })}
      </div>

      {/* Footer del Semillero */}
      <div className="hidden lg:block p-4 border-t border-slate-800 bg-slate-950/40">
        <div className="bg-slate-950/80 p-3 rounded-lg border border-slate-800 text-[11px] space-y-1.5">
          <div className="flex items-center justify-between">
            <span className="text-slate-400">Plataforma</span>
            <span className="text-teal-400 font-semibold flex items-center gap-1">
              <Sparkles className="w-3 h-3" /> v1.0 Oficial
            </span>
          </div>
          <p className="text-[10px] text-slate-500 leading-tight">
            Universidad de Cartagena • Facultad de Ingeniería
          </p>
        </div>
      </div>
    </aside>
  );
}
