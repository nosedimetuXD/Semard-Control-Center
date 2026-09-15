"use client";

import { useState, useEffect } from "react";
import {
  Printer,
  Upload,
  Layers,
  Clock,
  CheckCircle2,
  AlertTriangle,
  FileCheck,
  Cpu,
  Download,
  Play,
  RotateCcw,
} from "lucide-react";
import { apiCall } from "@/lib/api";

interface Print3DViewProps {
  user: any;
  onShowAlert: (msg: string, type: "info" | "success" | "error") => void;
  onTrackResponse: (res: any) => void;
}

export default function Print3DView({
  user,
  onShowAlert,
  onTrackResponse,
}: Print3DViewProps) {
  const [file, setFile] = useState<File | null>(null);
  const [material, setMaterial] = useState("PLA");
  const [color, setColor] = useState("Negro");
  const [infill, setInfill] = useState("20");
  const [notes, setNotes] = useState("");
  const [queue, setQueue] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);

  const defaultQueue = [
    {
      id: "p3d-1",
      filename: "Soporte_LiDAR_Sensor_v2.stl",
      material: "PLA",
      color: "Negro",
      infill_percentage: 20,
      user_name: "Daniel Franco",
      status: "IN_PROGRESS",
      estimated_time: "2h 45m",
      created_at: "Hoy, 10:30 AM",
    },
    {
      id: "p3d-2",
      filename: "Engranaje_Reductor_Eje_6mm.3mf",
      material: "PETG",
      color: "Naranja",
      infill_percentage: 50,
      user_name: "Investigador Mecatrónica",
      status: "APPROVED",
      estimated_time: "4h 10m",
      created_at: "Hoy, 09:15 AM",
    },
    {
      id: "p3d-3",
      filename: "Carcasa_ESP32_Antena_Externa.step",
      material: "PLA",
      color: "Gris",
      infill_percentage: 25,
      user_name: "Miembro SEMARD",
      status: "PENDING",
      estimated_time: "1h 30m",
      created_at: "Ayer, 04:20 PM",
    },
    {
      id: "p3d-4",
      filename: "Pinza_Robotica_Gripper.stl",
      material: "PETG",
      color: "Negro",
      infill_percentage: 40,
      user_name: "Semillero Robótica",
      status: "COMPLETED",
      estimated_time: "Finalizado",
      created_at: "12 Sept",
    },
  ];

  useEffect(() => {
    loadQueue();
  }, []);

  async function loadQueue() {
    const res = await apiCall("/api/v1/print3d/queue", { requireAuth: true });
    onTrackResponse(res);
    if (res.ok && Array.isArray(res.data) && res.data.length > 0) {
      setQueue(res.data);
    } else {
      setQueue(defaultQueue);
    }
  }

  async function handleUploadPrint() {
    if (!file) {
      return onShowAlert("Por favor selecciona un archivo 3D (.stl, .obj, .3mf, .step)", "error");
    }

    setLoading(true);
    const fd = new FormData();
    fd.append("file", file);
    fd.append("material", material);
    fd.append("color", color);
    fd.append("infill_percentage", infill);
    fd.append("notes", notes);

    const res = await apiCall("/api/v1/print3d/requests", {
      method: "POST",
      body: fd,
      requireAuth: true,
      isFormData: true,
    });
    onTrackResponse(res);
    setLoading(false);

    if (res.ok) {
      onShowAlert("Modelo 3D cargado e ingresado a la cola de fabricación", "success");
      setFile(null);
      setNotes("");
      loadQueue();
    } else {
      onShowAlert(res.data?.error || "Error al subir modelo 3D", "error");
    }
  }

  async function handleOperatorAdvance(id: string, nextStatus: string) {
    const res = await apiCall(`/api/v1/print3d/requests/${id}/status`, {
      method: "POST",
      body: { status: nextStatus },
      requireAuth: true,
    });
    onTrackResponse(res);
    if (res.ok) {
      onShowAlert(`Estado de la pieza actualizado a ${nextStatus}`, "success");
      loadQueue();
    } else {
      onShowAlert(res.data?.error || "No tienes permisos de Operador 3D", "error");
    }
  }

  const getStatusBadge = (status: string) => {
    switch (status) {
      case "IN_PROGRESS":
        return <span className="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-cyan-500/20 text-cyan-300 border border-cyan-500/30 animate-pulse">En Impresión</span>;
      case "APPROVED":
        return <span className="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-emerald-500/20 text-emerald-300 border border-emerald-500/30">Aprobado</span>;
      case "COMPLETED":
        return <span className="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-teal-500/20 text-teal-300 border border-teal-500/30">Completada</span>;
      case "DELIVERED":
        return <span className="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-slate-800 text-slate-300 border border-slate-700">Entregada</span>;
      default:
        return <span className="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-amber-500/20 text-amber-300 border border-amber-500/30">Pendiente Revisión</span>;
    }
  };

  return (
    <div className="space-y-6">
      {/* Encabezado */}
      <div>
        <h2 className="text-xl font-bold text-white flex items-center gap-2">
          <Printer className="w-5 h-5 text-teal-400" />
          <span>Taller de Fabricación & Impresión 3D</span>
        </h2>
        <p className="text-xs text-slate-400 mt-0.5">
          Prototipado rápido de componentes mecánicos, soportes de sensores y carcasas para proyectos de investigación.
        </p>
      </div>

      {/* Grid Principal: Formulario de Carga (Izquierda) + Cola en Vivo (Derecha) */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        {/* Formulario de Encargo 3D (5 Cols) */}
        <div className="lg:col-span-5 space-y-4">
          <div className="p-5 rounded-2xl bg-slate-900/90 border border-slate-800 space-y-4 shadow-xl">
            <div className="flex items-center justify-between border-b border-slate-800 pb-3">
              <h3 className="text-sm font-bold text-white flex items-center gap-2">
                <Upload className="w-4 h-4 text-teal-400" />
                <span>Encargar Fabricación de Pieza</span>
              </h3>
              <span className="text-[10px] px-2 py-0.5 rounded bg-teal-500/10 text-teal-400 border border-teal-500/20">
                Máx. 50MB
              </span>
            </div>

            {/* Zona de Arrastre / Selección de Archivo */}
            <label className="border-2 border-dashed border-slate-800 hover:border-teal-500/60 rounded-xl p-5 flex flex-col items-center justify-center text-center cursor-pointer transition bg-slate-950/40 hover:bg-slate-950/80">
              <input
                type="file"
                accept=".stl,.obj,.3mf,.step"
                onChange={(e) => setFile(e.target.files?.[0] || null)}
                className="hidden"
              />
              <Printer className="w-8 h-8 text-slate-500 mb-2" />
              <span className="text-xs font-semibold text-slate-200">
                {file ? file.name : "Selecciona tu archivo 3D"}
              </span>
              <span className="text-[10px] text-slate-500 mt-1">
                Formatos permitidos: .stl, .obj, .3mf, .step
              </span>
            </label>

            {/* Parámetros de Impresión */}
            <div className="space-y-3 text-xs">
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="text-[11px] text-slate-400 block mb-1">Material</label>
                  <select
                    value={material}
                    onChange={(e) => setMaterial(e.target.value)}
                    className="input-field"
                  >
                    <option value="PLA">PLA (Estándar)</option>
                    <option value="PETG">PETG (Alta Resistencia)</option>
                    <option value="ABS">ABS (Térmico)</option>
                    <option value="RESINA">Resina UV (Detalle Fino)</option>
                  </select>
                </div>
                <div>
                  <label className="text-[11px] text-slate-400 block mb-1">Color</label>
                  <select
                    value={color}
                    onChange={(e) => setColor(e.target.value)}
                    className="input-field"
                  >
                    <option value="Negro">Negro Mate</option>
                    <option value="Blanco">Blanco</option>
                    <option value="Gris">Gris Técnico</option>
                    <option value="Naranja">Naranja UdeC</option>
                    <option value="Azul">Azul Eléctrico</option>
                  </select>
                </div>
              </div>

              <div>
                <div className="flex items-center justify-between text-[11px] text-slate-400 mb-1">
                  <span>Densidad de Relleno (Infill)</span>
                  <span className="font-mono text-teal-400 font-bold">{infill}%</span>
                </div>
                <input
                  type="range"
                  min="10"
                  max="100"
                  step="5"
                  value={infill}
                  onChange={(e) => setInfill(e.target.value)}
                  className="w-full accent-teal-500 cursor-pointer"
                />
                <div className="flex justify-between text-[9px] text-slate-500 mt-0.5">
                  <span>10% (Ligero)</span>
                  <span>20% (Normal)</span>
                  <span>50% (Fuerte)</span>
                  <span>100% (Sólido)</span>
                </div>
              </div>

              <div>
                <label className="text-[11px] text-slate-400 block mb-1">Notas de Impresión o Tolerancia</label>
                <textarea
                  placeholder="Ej: Requiere soportes en voladizos, orientación específica para roscas..."
                  rows={2}
                  value={notes}
                  onChange={(e) => setNotes(e.target.value)}
                  className="input-field"
                />
              </div>

              <button
                disabled={loading}
                onClick={handleUploadPrint}
                className="btn-primary w-full py-2.5 shadow-md shadow-teal-500/20"
              >
                {loading ? "Cargando archivo 3D..." : "Enviar a Cola de Fabricación"}
              </button>
            </div>
          </div>
        </div>

        {/* Cola de Fabricación en Tiempo Real (7 Cols) */}
        <div className="lg:col-span-7 space-y-4">
          <div className="p-5 rounded-2xl bg-slate-900/90 border border-slate-800 space-y-4 shadow-xl">
            <div className="flex items-center justify-between border-b border-slate-800 pb-3">
              <div>
                <h3 className="text-sm font-bold text-white flex items-center gap-2">
                  <Layers className="w-4 h-4 text-cyan-400" />
                  <span>Cola de Fabricación en Vivo</span>
                </h3>
                <p className="text-[11px] text-slate-400">Estado de solicitudes enviadas por los integrantes</p>
              </div>
              <button
                onClick={loadQueue}
                className="text-xs text-teal-400 hover:text-teal-300 font-semibold"
              >
                Actualizar
              </button>
            </div>

            <div className="space-y-3">
              {queue.map((item) => (
                <div
                  key={item.id}
                  className="p-4 rounded-xl bg-slate-950/60 border border-slate-800 hover:border-slate-700 transition space-y-2.5"
                >
                  <div className="flex flex-wrap items-center justify-between gap-2">
                    <div className="flex items-center gap-2 min-w-0">
                      <div className="w-8 h-8 rounded-lg bg-cyan-500/10 text-cyan-400 flex items-center justify-center font-bold text-xs border border-cyan-500/20 shrink-0">
                        3D
                      </div>
                      <div className="min-w-0">
                        <h4 className="text-xs font-bold text-white truncate">
                          {item.filename || item.file_name || "Pieza_SEMARD.stl"}
                        </h4>
                        <p className="text-[11px] text-slate-400">
                          {item.user_name || "Estudiante SEMARD"} • {item.material} • {item.color} • {item.infill_percentage}% infill
                        </p>
                      </div>
                    </div>
                    <div>{getStatusBadge(item.status)}</div>
                  </div>

                  {/* Acciones de Operador Técnico */}
                  <div className="flex items-center justify-between pt-2 border-t border-slate-800/80 text-[11px] text-slate-400">
                    <div className="flex items-center gap-1.5">
                      <Clock className="w-3.5 h-3.5 text-slate-500" />
                      <span>{item.created_at || "Reciente"}</span>
                    </div>

                    <div className="flex items-center gap-2">
                      {item.status === "APPROVED" && (
                        <button
                          onClick={() => handleOperatorAdvance(item.id, "IN_PROGRESS")}
                          className="px-2.5 py-1 rounded bg-cyan-600/20 text-cyan-300 border border-cyan-500/30 hover:bg-cyan-600/30 font-semibold flex items-center gap-1"
                        >
                          <Play className="w-3 h-3" />
                          <span>Poner en Máquina</span>
                        </button>
                      )}
                      {item.status === "IN_PROGRESS" && (
                        <button
                          onClick={() => handleOperatorAdvance(item.id, "COMPLETED")}
                          className="px-2.5 py-1 rounded bg-emerald-600/20 text-emerald-300 border border-emerald-500/30 hover:bg-emerald-600/30 font-semibold flex items-center gap-1"
                        >
                          <CheckCircle2 className="w-3 h-3" />
                          <span>Marcar Finalizada</span>
                        </button>
                      )}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
