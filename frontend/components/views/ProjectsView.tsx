"use client";

import { useState, useEffect } from "react";
import {
  FolderGit2,
  Plus,
  Users,
  FileText,
  BadgeDollarSign,
  Calendar,
  Sparkles,
  ExternalLink,
  CheckCircle2,
  X,
} from "lucide-react";
import { apiCall } from "@/lib/api";

interface ProjectsViewProps {
  user: any;
  onShowAlert: (msg: string, type: "info" | "success" | "error") => void;
  onTrackResponse: (res: any) => void;
}

export default function ProjectsView({
  user,
  onShowAlert,
  onTrackResponse,
}: ProjectsViewProps) {
  const [selectedLine, setSelectedLine] = useState("TODOS");
  const [projects, setProjects] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  // Modales
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showUpdateModal, setShowUpdateModal] = useState(false);
  const [selectedProjectId, setSelectedProjectId] = useState<string>("");

  // Form states
  const [newTitle, setNewTitle] = useState("");
  const [newDesc, setNewDesc] = useState("");
  const [newLine, setNewLine] = useState("Robótica Móvil");
  const [newStartDate, setNewStartDate] = useState("");

  const [updTitle, setUpdTitle] = useState("");
  const [updContent, setUpdContent] = useState("");
  const [updLinks, setUpdLinks] = useState("");

  const defaultProjects = [
    {
      id: "demo-1",
      title: "Vehículo Autónomo Todoterreno (Rover UdeC)",
      description: "Plataforma robótica móvil con tracción 4x4, navegación por SLAM y sensores LiDAR y ultrasónicos para mapeo de terrenos agrestes.",
      research_line: "Robótica Móvil",
      status: "IN_PROGRESS",
      lead: "Daniel Franco",
      members_count: 4,
      start_date: "2026-02-15",
    },
    {
      id: "demo-2",
      title: "Brazo Robótico Manipulador de 6 Grados de Libertad",
      description: "Desarrollo cinemático directo e inverso con servomotores de alto torque y cámara RGB-D para tareas de Pick & Place inteligente.",
      research_line: "Automatización Industrial",
      status: "IN_PROGRESS",
      lead: "Ing. Mecatrónica",
      members_count: 3,
      start_date: "2026-03-01",
    },
    {
      id: "demo-3",
      title: "Red de Monitoreo Ambiental de Calidad del Aire (IoT)",
      description: "Estaciones meteorológicas y sensores de material particulado distribuidos en la Universidad de Cartagena con telemetría LoRaWAN.",
      research_line: "Sistemas Embebidos & IoT",
      status: "IN_PROGRESS",
      lead: "Ing. de Sistemas",
      members_count: 5,
      start_date: "2026-01-20",
    },
    {
      id: "demo-4",
      title: "Dron Autónomo Hexacóptero para Inspección de Infraestructura",
      description: "Aeronave no tripulada con controlador de vuelo PX4 y visión artificial para detección de fisuras en puentes y edificaciones.",
      research_line: "Drones & Percepción",
      status: "IN_PROGRESS",
      lead: "Semillero SEMARD",
      members_count: 3,
      start_date: "2026-04-10",
    },
  ];

  useEffect(() => {
    loadProjects();
  }, []);

  async function loadProjects() {
    setLoading(true);
    const res = await apiCall("/api/v1/projects/showcase", { requireAuth: false });
    onTrackResponse(res);
    if (res.ok && Array.isArray(res.data) && res.data.length > 0) {
      setProjects(res.data);
    } else {
      setProjects(defaultProjects);
    }
    setLoading(false);
  }

  async function handleCreateProject() {
    if (!newTitle.trim()) {
      return onShowAlert("El título del proyecto es obligatorio", "error");
    }
    const res = await apiCall("/api/v1/projects", {
      method: "POST",
      body: {
        title: newTitle,
        description: newDesc,
        research_line: newLine,
        start_date: newStartDate ? new Date(newStartDate).toISOString() : new Date().toISOString(),
      },
      requireAuth: true,
    });
    onTrackResponse(res);
    if (res.ok) {
      onShowAlert("Proyecto de investigación creado con éxito", "success");
      setShowCreateModal(false);
      setNewTitle("");
      setNewDesc("");
      loadProjects();
    } else {
      onShowAlert(res.data?.error || "Error al crear proyecto", "error");
    }
  }

  async function handleSubmitUpdate() {
    if (!updTitle.trim() || !selectedProjectId) {
      return onShowAlert("El título del avance es obligatorio", "error");
    }
    const links = updLinks ? updLinks.split(",").map((s) => s.trim()) : [];
    const res = await apiCall(`/api/v1/projects/${selectedProjectId}/updates`, {
      method: "POST",
      body: {
        title: updTitle,
        content: updContent,
        attachments_url: links,
      },
      requireAuth: true,
    });
    onTrackResponse(res);
    if (res.ok) {
      onShowAlert("Avance de proyecto radicado con éxito para evaluación", "success");
      setShowUpdateModal(false);
      setUpdTitle("");
      setUpdContent("");
      setUpdLinks("");
    } else {
      onShowAlert(res.data?.error || "Error al radicar avance", "error");
    }
  }

  const lines = [
    "TODOS",
    "Robótica Móvil",
    "Automatización Industrial",
    "Sistemas Embebidos & IoT",
    "Drones & Percepción",
  ];

  const filteredProjects =
    selectedLine === "TODOS"
      ? projects
      : projects.filter((p) => p.research_line === selectedLine);

  return (
    <div className="space-y-6">
      {/* Cabecera del Módulo */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h2 className="text-xl font-bold text-white flex items-center gap-2">
            <FolderGit2 className="w-5 h-5 text-teal-400" />
            <span>Proyectos de Investigación</span>
          </h2>
          <p className="text-xs text-slate-400 mt-0.5">
            Líneas de desarrollo científico y tecnológico en robótica, automatización y mecatrónica.
          </p>
        </div>

        <button
          onClick={() => setShowCreateModal(true)}
          className="self-start md:self-auto text-xs font-semibold px-3.5 py-2 rounded-lg bg-teal-600 hover:bg-teal-500 text-white flex items-center gap-2 shadow-sm shadow-teal-500/20 transition"
        >
          <Plus className="w-4 h-4" />
          <span>Nuevo Proyecto</span>
        </button>
      </div>

      {/* Filtros de Líneas de Investigación */}
      <div className="flex items-center gap-1.5 overflow-x-auto pb-1 text-xs">
        {lines.map((line) => (
          <button
            key={line}
            onClick={() => setSelectedLine(line)}
            className={`px-3 py-1.5 rounded-lg font-medium transition whitespace-nowrap ${
              selectedLine === line
                ? "bg-teal-600 text-white shadow-sm"
                : "bg-slate-900 border border-slate-800 text-slate-400 hover:text-slate-200"
            }`}
          >
            {line}
          </button>
        ))}
      </div>

      {/* Grid de Tarjetas de Proyecto */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
        {filteredProjects.map((p) => (
          <div
            key={p.id}
            className="p-5 rounded-2xl bg-slate-900/90 border border-slate-800 hover:border-slate-700 transition space-y-4 flex flex-col justify-between shadow-lg"
          >
            <div className="space-y-2.5">
              <div className="flex items-start justify-between gap-2">
                <span className="text-[10px] font-bold uppercase tracking-wider px-2.5 py-0.5 rounded-full bg-teal-500/10 text-teal-400 border border-teal-500/20">
                  {p.research_line}
                </span>
                <span className="text-[10px] font-semibold px-2 py-0.5 rounded bg-emerald-500/10 text-emerald-300 border border-emerald-500/20">
                  En Desarrollo
                </span>
              </div>

              <h3 className="text-base font-bold text-white tracking-tight leading-snug">
                {p.title}
              </h3>
              <p className="text-xs text-slate-300 leading-relaxed line-clamp-3">
                {p.description}
              </p>
            </div>

            <div className="space-y-3 pt-3 border-t border-slate-800/80">
              <div className="flex items-center justify-between text-xs text-slate-400">
                <div className="flex items-center gap-1.5">
                  <Users className="w-3.5 h-3.5 text-slate-500" />
                  <span>Equipo: {p.lead || "SEMARD UdeC"}</span>
                </div>
                <div className="flex items-center gap-1">
                  <Calendar className="w-3.5 h-3.5 text-slate-500" />
                  <span className="text-[11px]">{p.start_date ? p.start_date.split("T")[0] : "2026"}</span>
                </div>
              </div>

              <div className="flex items-center gap-2 pt-1">
                <button
                  onClick={() => {
                    setSelectedProjectId(p.id);
                    setShowUpdateModal(true);
                  }}
                  className="flex-1 text-xs font-semibold py-2 px-3 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 border border-slate-700 transition flex items-center justify-center gap-1.5"
                >
                  <FileText className="w-3.5 h-3.5 text-teal-400" />
                  <span>Registrar Avance</span>
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Modal: Crear Proyecto */}
      {showCreateModal && (
        <div className="fixed inset-0 z-50 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-2xl max-w-lg w-full p-6 space-y-4 shadow-2xl animate-in fade-in zoom-in-95">
            <div className="flex items-center justify-between">
              <h3 className="text-base font-bold text-white">Nuevo Proyecto de Investigación</h3>
              <button
                onClick={() => setShowCreateModal(false)}
                className="text-slate-400 hover:text-white"
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            <div className="space-y-3 text-xs">
              <div>
                <label className="text-[11px] text-slate-400 block mb-1">Título del Proyecto</label>
                <input
                  type="text"
                  placeholder="Ej: Sistema de Inspección Submarina con ROV"
                  value={newTitle}
                  onChange={(e) => setNewTitle(e.target.value)}
                  className="input-field"
                />
              </div>

              <div>
                <label className="text-[11px] text-slate-400 block mb-1">Línea de Investigación</label>
                <select
                  value={newLine}
                  onChange={(e) => setNewLine(e.target.value)}
                  className="input-field"
                >
                  <option value="Robótica Móvil">Robótica Móvil</option>
                  <option value="Automatización Industrial">Automatización Industrial</option>
                  <option value="Sistemas Embebidos & IoT">Sistemas Embebidos & IoT</option>
                  <option value="Drones & Percepción">Drones & Percepción</option>
                  <option value="Inteligencia Artificial Aplicada">Inteligencia Artificial Aplicada</option>
                </select>
              </div>

              <div>
                <label className="text-[11px] text-slate-400 block mb-1">Descripción y Objetivos</label>
                <textarea
                  placeholder="Alcance, objetivos del proyecto y metodologías aplicadas..."
                  rows={3}
                  value={newDesc}
                  onChange={(e) => setNewDesc(e.target.value)}
                  className="input-field"
                />
              </div>

              <div>
                <label className="text-[11px] text-slate-400 block mb-1">Fecha de Inicio</label>
                <input
                  type="date"
                  value={newStartDate}
                  onChange={(e) => setNewStartDate(e.target.value)}
                  className="input-field"
                />
              </div>
            </div>

            <div className="flex gap-2 pt-2">
              <button
                onClick={() => setShowCreateModal(false)}
                className="btn-action flex-1 py-2"
              >
                Cancelar
              </button>
              <button
                onClick={handleCreateProject}
                className="btn-primary flex-1 py-2"
              >
                Crear Proyecto
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Modal: Registrar Avance */}
      {showUpdateModal && (
        <div className="fixed inset-0 z-50 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-2xl max-w-lg w-full p-6 space-y-4 shadow-2xl animate-in fade-in zoom-in-95">
            <div className="flex items-center justify-between">
              <h3 className="text-base font-bold text-white">Registrar Avance de Proyecto</h3>
              <button
                onClick={() => setShowUpdateModal(false)}
                className="text-slate-400 hover:text-white"
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            <div className="space-y-3 text-xs">
              <div>
                <label className="text-[11px] text-slate-400 block mb-1">Título del Avance o Hito</label>
                <input
                  type="text"
                  placeholder="Ej: Calibración de cinemática y pruebas de agarre"
                  value={updTitle}
                  onChange={(e) => setUpdTitle(e.target.value)}
                  className="input-field"
                />
              </div>

              <div>
                <label className="text-[11px] text-slate-400 block mb-1">Descripción y Resultados</label>
                <textarea
                  placeholder="Detalla qué se logró, dificultades encontradas y próximos pasos..."
                  rows={3}
                  value={updContent}
                  onChange={(e) => setUpdContent(e.target.value)}
                  className="input-field"
                />
              </div>

              <div>
                <label className="text-[11px] text-slate-400 block mb-1">Enlaces de Evidencia (GitHub, Drive, Video)</label>
                <input
                  type="text"
                  placeholder="https://github.com/..., https://drive.google.com/..."
                  value={updLinks}
                  onChange={(e) => setUpdLinks(e.target.value)}
                  className="input-field"
                />
              </div>
            </div>

            <div className="flex gap-2 pt-2">
              <button
                onClick={() => setShowUpdateModal(false)}
                className="btn-action flex-1 py-2"
              >
                Cancelar
              </button>
              <button
                onClick={handleSubmitUpdate}
                className="btn-primary flex-1 py-2"
              >
                Radicar Avance
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
