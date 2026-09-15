"use client";

import {
  FolderGit2,
  Printer,
  Package,
  Calendar,
  ArrowUpRight,
  Clock,
  CheckCircle2,
  AlertCircle,
  Cpu,
  Layers,
  Sparkles,
  ChevronRight,
  TrendingUp,
} from "lucide-react";
import { ViewKey } from "@/components/Sidebar";

interface DashboardViewProps {
  onNavigate: (view: ViewKey) => void;
  onOpenQuickAction: (actionType: string) => void;
}

export default function DashboardView({
  onNavigate,
  onOpenQuickAction,
}: DashboardViewProps) {
  const stats = [
    {
      title: "Proyectos en Marcha",
      val: "6 Activos",
      sub: "Robótica, IoT y Drones",
      icon: FolderGit2,
      color: "from-teal-500/20 to-teal-600/5 text-teal-400 border-teal-500/30",
      view: "projects" as ViewKey,
    },
    {
      title: "Cola de Impresión 3D",
      val: "4 en Espera",
      sub: "1 pieza en máquina activa",
      icon: Printer,
      color: "from-cyan-500/20 to-cyan-600/5 text-cyan-400 border-cyan-500/30",
      view: "print3d" as ViewKey,
    },
    {
      title: "Equipos de Laboratorio",
      val: "38 Ítems",
      sub: "32 disponibles para préstamo",
      icon: Package,
      color: "from-indigo-500/20 to-indigo-600/5 text-indigo-400 border-indigo-500/30",
      view: "inventory" as ViewKey,
    },
    {
      title: "Eventos & Talleres",
      val: "2 Próximos",
      sub: "Inscripciones abiertas",
      icon: Calendar,
      color: "from-amber-500/20 to-amber-600/5 text-amber-400 border-amber-500/30",
      view: "events" as ViewKey,
    },
  ];

  const featuredProjects = [
    {
      id: "1",
      title: "Vehículo Autónomo Todoterreno (Rover UdeC)",
      line: "Robótica Móvil",
      lead: "Daniel Franco",
      progress: 75,
      status: "Activo • Fase Pruebas",
    },
    {
      id: "2",
      title: "Brazo Robótico para Clasificación por Visión Artificial",
      line: "Automatización Industrial",
      lead: "Ing. Mecatrónica",
      progress: 60,
      status: "Activo • Integración",
    },
    {
      id: "3",
      title: "Red de Monitoreo Ambiental IoT para el Campus",
      line: "Sistemas Embebidos & IoT",
      lead: "Ing. de Sistemas",
      progress: 90,
      status: "Fase Final • Calibración",
    },
  ];

  const printers = [
    { name: "Ender-3 S1 Pro #01", status: "Imprimiendo", job: "Soporte Sensor LiDAR (PLA)", pct: 68, color: "text-cyan-400 bg-cyan-500/20 border-cyan-500/30" },
    { name: "Bambu Lab P1P #02", status: "Disponible", job: "En espera de cola técnica", pct: 0, color: "text-emerald-400 bg-emerald-500/20 border-emerald-500/30" },
    { name: "Anycubic Photon Mono", status: "Mantenimiento", job: "Cambio de resina UV", pct: 0, color: "text-amber-400 bg-amber-500/20 border-amber-500/30" },
  ];

  return (
    <div className="space-y-6">
      {/* Banner de Bienvenida Institucional */}
      <div className="relative overflow-hidden rounded-2xl bg-gradient-to-r from-slate-900 via-slate-900 to-teal-950/50 p-6 md:p-8 border border-slate-800 shadow-xl">
        <div className="absolute top-0 right-0 w-96 h-96 bg-teal-500/10 rounded-full blur-3xl pointer-events-none -mr-20 -mt-20" />
        <div className="relative z-10 max-w-2xl space-y-2.5">
          <div className="inline-flex items-center gap-2 px-2.5 py-1 rounded-full bg-teal-500/10 border border-teal-500/30 text-teal-400 text-xs font-semibold">
            <Sparkles className="w-3.5 h-3.5" />
            <span>Semillero de Investigación en Robótica y Automatización</span>
          </div>
          <h2 className="text-2xl md:text-3xl font-bold tracking-tight text-white leading-tight">
            Control Center y Gestión Operativa
          </h2>
          <p className="text-xs md:text-sm text-slate-300 leading-relaxed">
            Plataforma centralizada para la coordinación de proyectos científicos, fabricación de piezas en impresión 3D, préstamos de instrumentación y divulgación académica en la Universidad de Cartagena.
          </p>
          <div className="pt-2 flex flex-wrap gap-2.5">
            <button
              onClick={() => onOpenQuickAction("print3d")}
              className="px-3.5 py-2 rounded-lg bg-teal-600 hover:bg-teal-500 text-white font-semibold text-xs transition shadow-md shadow-teal-500/20 flex items-center gap-1.5"
            >
              <Printer className="w-3.5 h-3.5" />
              <span>Encargar Pieza 3D</span>
            </button>
            <button
              onClick={() => onOpenQuickAction("loan")}
              className="px-3.5 py-2 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 border border-slate-700 text-xs font-semibold transition flex items-center gap-1.5"
            >
              <Package className="w-3.5 h-3.5" />
              <span>Solicitar Préstamo de Equipo</span>
            </button>
          </div>
        </div>
      </div>

      {/* Grid de Métricas Principales (KPI Cards) */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((st, i) => {
          const Icon = st.icon;
          return (
            <button
              key={i}
              onClick={() => onNavigate(st.view)}
              className={`p-4 rounded-xl border bg-gradient-to-b ${st.color} text-left transition hover:scale-[1.02] hover:border-teal-400/50 flex flex-col justify-between`}
            >
              <div className="flex items-center justify-between mb-3">
                <span className="text-xs font-medium text-slate-300">{st.title}</span>
                <div className="p-2 rounded-lg bg-slate-900/60 border border-slate-700/50">
                  <Icon className="w-4 h-4" />
                </div>
              </div>
              <div>
                <div className="text-2xl font-bold text-white tracking-tight">{st.val}</div>
                <div className="text-[11px] text-slate-400 mt-0.5 flex items-center justify-between">
                  <span>{st.sub}</span>
                  <ChevronRight className="w-3 h-3 text-slate-500" />
                </div>
              </div>
            </button>
          );
        })}
      </div>

      {/* Cuerpo Dividido: Proyectos Destacados + Estado de Laboratorio */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Columna Izquierda: Proyectos Destacados (2/3) */}
        <div className="lg:col-span-2 space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-sm font-bold text-white flex items-center gap-2">
              <FolderGit2 className="w-4 h-4 text-teal-400" />
              <span>Proyectos de Investigación Activos</span>
            </h3>
            <button
              onClick={() => onNavigate("projects")}
              className="text-xs text-teal-400 hover:text-teal-300 font-semibold flex items-center gap-1"
            >
              <span>Ver todos los proyectos</span>
              <ArrowUpRight className="w-3.5 h-3.5" />
            </button>
          </div>

          <div className="space-y-3">
            {featuredProjects.map((p) => (
              <div
                key={p.id}
                className="p-4 rounded-xl bg-slate-900/80 border border-slate-800 hover:border-slate-700 transition space-y-3"
              >
                <div className="flex flex-wrap items-start justify-between gap-2">
                  <div>
                    <span className="text-[10px] font-semibold uppercase tracking-wider px-2 py-0.5 rounded bg-teal-500/10 text-teal-400 border border-teal-500/20">
                      {p.line}
                    </span>
                    <h4 className="text-sm font-semibold text-white mt-1.5">{p.title}</h4>
                    <p className="text-xs text-slate-400">Líder: {p.lead}</p>
                  </div>
                  <span className="text-[11px] font-medium px-2 py-0.5 rounded bg-slate-800 text-slate-300 border border-slate-700">
                    {p.status}
                  </span>
                </div>

                <div>
                  <div className="flex items-center justify-between text-[11px] text-slate-400 mb-1">
                    <span>Avance estimado</span>
                    <span className="font-semibold text-teal-400">{p.progress}%</span>
                  </div>
                  <div className="w-full h-1.5 bg-slate-800 rounded-full overflow-hidden">
                    <div
                      className="h-full bg-gradient-to-r from-teal-500 to-emerald-400 rounded-full"
                      style={{ width: `${p.progress}%` }}
                    />
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Columna Derecha: Taller 3D e Impresoras (1/3) */}
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-sm font-bold text-white flex items-center gap-2">
              <Printer className="w-4 h-4 text-cyan-400" />
              <span>Estado del Taller 3D</span>
            </h3>
            <button
              onClick={() => onNavigate("print3d")}
              className="text-xs text-cyan-400 hover:text-cyan-300 font-semibold"
            >
              Ver cola
            </button>
          </div>

          <div className="space-y-3">
            {printers.map((pr, idx) => (
              <div
                key={idx}
                className="p-3.5 rounded-xl bg-slate-900/80 border border-slate-800 space-y-2"
              >
                <div className="flex items-center justify-between">
                  <span className="text-xs font-semibold text-white">{pr.name}</span>
                  <span className={`text-[10px] font-bold px-2 py-0.5 rounded-full border ${pr.color}`}>
                    {pr.status}
                  </span>
                </div>
                <p className="text-[11px] text-slate-400 truncate">{pr.job}</p>
                {pr.pct > 0 && (
                  <div>
                    <div className="flex items-center justify-between text-[10px] text-slate-400 mb-1">
                      <span>Progreso de impresión</span>
                      <span className="font-mono text-cyan-300 font-bold">{pr.pct}%</span>
                    </div>
                    <div className="w-full h-1 bg-slate-800 rounded-full overflow-hidden">
                      <div
                        className="h-full bg-cyan-400 rounded-full"
                        style={{ width: `${pr.pct}%` }}
                      />
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>

          {/* Tarjeta de Próximo Taller SEMARD */}
          <div className="p-4 rounded-xl bg-gradient-to-br from-indigo-950/40 to-slate-900 border border-indigo-500/30 space-y-2">
            <span className="text-[10px] font-bold uppercase tracking-wider text-indigo-400">
              Próximo Taller Presencial
            </span>
            <h4 className="text-xs font-semibold text-white">
              Introducción al Control de Motores Brushless y ESC con ESP32
            </h4>
            <p className="text-[11px] text-slate-400">Viernes • Laboratorio de Robótica (Sede Piedra de Bolívar)</p>
            <button
              onClick={() => onNavigate("events")}
              className="w-full text-xs font-semibold py-1.5 rounded-lg bg-indigo-600 hover:bg-indigo-500 text-white transition mt-1"
            >
              Ver Convocatoria
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
