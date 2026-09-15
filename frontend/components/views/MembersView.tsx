"use client";

import { useState, useEffect } from "react";
import {
  Users,
  ShieldCheck,
  Award,
  Printer,
  CheckCircle,
  XCircle,
  UserCheck,
  UserPlus,
  Mail,
  GraduationCap,
  X,
} from "lucide-react";
import { apiCall } from "@/lib/api";

interface MembersViewProps {
  user: any;
  onShowAlert: (msg: string, type: "info" | "success" | "error") => void;
  onTrackResponse: (res: any) => void;
}

export default function MembersView({
  user,
  onShowAlert,
  onTrackResponse,
}: MembersViewProps) {
  const [members, setMembers] = useState<any[]>([]);
  const [requests, setRequests] = useState<any[]>([]);
  const [activeTab, setActiveTab] = useState<"DIRECTORY" | "REQUESTS">("DIRECTORY");

  // Modal para evaluar solicitud
  const [selectedReq, setSelectedReq] = useState<any | null>(null);
  const [evalAction, setEvalAction] = useState<"APPROVE" | "REJECT">("APPROVE");
  const [evalRole, setEvalRole] = useState("MIEMBRO");
  const [evalFeedback, setEvalFeedback] = useState("");

  const defaultMembers = [
    {
      id: "m-1",
      full_name: "Daniel Franco",
      email: "dfrancob1@unicartagena.edu.co",
      student_code: "0222410043",
      role: "DIRECTOR",
      can_operate_3d: true,
      bio: "Director General del Semillero SEMARD. Enfoque en robótica móvil autónoma y sistemas mecatrónicos.",
    },
    {
      id: "m-2",
      full_name: "Investigador Principal Mecatrónica",
      email: "investigador.mecatronica@unicartagena.edu.co",
      student_code: "0222310012",
      role: "ADMINISTRADOR",
      can_operate_3d: true,
      bio: "Coordinador del área de cinemática y manipuladores robóticos industriales.",
    },
    {
      id: "m-3",
      full_name: "Estudiante Investigador IoT",
      email: "iot.estudiante@unicartagena.edu.co",
      student_code: "0222420088",
      role: "MIEMBRO",
      can_operate_3d: false,
      bio: "Desarrollo de firmware en ESP32 y sensores ambientales LoRaWAN.",
    },
    {
      id: "m-4",
      full_name: "Operador Especialista 3D",
      email: "operador3d@unicartagena.edu.co",
      student_code: "0222430055",
      role: "MIEMBRO",
      can_operate_3d: true,
      bio: "Optimización de parámetros de laminado y manufactura aditiva en PLA y PETG.",
    },
  ];

  const defaultRequests = [
    {
      id: "req-1",
      full_name: "Carlos Mendoza Pardo",
      email: "cmendozap@unicartagena.edu.co",
      student_code: "0222510034",
      career_program: "Ingeniería de Sistemas",
      motivation_letter: "Tengo experiencia en Python y visión por computador, deseo vincularme a los proyectos de SLAM para robótica móvil.",
      status: "PENDING",
      created_at: "Hoy",
    },
    {
      id: "req-2",
      full_name: "Laura Gómez Torres",
      email: "lgomezt@unicartagena.edu.co",
      student_code: "0222520091",
      career_program: "Ingeniería Mecánica",
      motivation_letter: "Interés en el diseño mecánico y simulación de esfuerzos en SolidWorks para chasis robóticos.",
      status: "PENDING",
      created_at: "Ayer",
    },
  ];

  useEffect(() => {
    loadMembers();
    if (user?.role === "DIRECTOR") {
      loadRequests();
    }
  }, [user]);

  async function loadMembers() {
    const res = await apiCall("/api/v1/users", { requireAuth: true });
    onTrackResponse(res);
    if (res.ok && Array.isArray(res.data) && res.data.length > 0) {
      setMembers(res.data);
    } else {
      setMembers(defaultMembers);
    }
  }

  async function loadRequests() {
    const res = await apiCall("/api/v1/directors/registration-requests", { requireAuth: true });
    onTrackResponse(res);
    if (res.ok && Array.isArray(res.data)) {
      setRequests(res.data);
    } else {
      setRequests(defaultRequests);
    }
  }

  async function handleReviewRequest() {
    if (!selectedReq) return;

    const res = await apiCall(`/api/v1/directors/registration-requests/${selectedReq.id}/review`, {
      method: "POST",
      body: {
        action: evalAction,
        role: evalRole,
        feedback: evalFeedback || null,
      },
      requireAuth: true,
    });
    onTrackResponse(res);

    if (res.ok) {
      onShowAlert(`Solicitud de ${selectedReq.full_name} evaluada con éxito (${evalAction})`, "success");
      setSelectedReq(null);
      setEvalFeedback("");
      loadRequests();
      loadMembers();
    } else {
      onShowAlert(res.data?.error || "Error al evaluar solicitud", "error");
    }
  }

  async function handleToggle3D(userId: string, canOperate: boolean) {
    const res = await apiCall(`/api/v1/directors/users/${userId}/permissions`, {
      method: "PATCH",
      body: { can_operate_3d: canOperate },
      requireAuth: true,
    });
    onTrackResponse(res);
    if (res.ok) {
      onShowAlert("Permisos de Operador 3D actualizados", "success");
      loadMembers();
    } else {
      onShowAlert(res.data?.error || "Error al actualizar permisos", "error");
    }
  }

  return (
    <div className="space-y-6">
      {/* Cabecera */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h2 className="text-xl font-bold text-white flex items-center gap-2">
            <Users className="w-5 h-5 text-teal-400" />
            <span>Comunidad & Directorio de Investigadores</span>
          </h2>
          <p className="text-xs text-slate-400 mt-0.5">
            Directores, administradores de laboratorio y miembros del semillero de investigación.
          </p>
        </div>

        {user?.role === "DIRECTOR" && (
          <div className="flex items-center gap-2 bg-slate-900 p-1 rounded-lg border border-slate-800 text-xs">
            <button
              onClick={() => setActiveTab("DIRECTORY")}
              className={`px-3 py-1.5 rounded-md font-medium transition ${
                activeTab === "DIRECTORY" ? "bg-teal-600 text-white shadow-sm" : "text-slate-400 hover:text-slate-200"
              }`}
            >
              Directorio General
            </button>
            <button
              onClick={() => setActiveTab("REQUESTS")}
              className={`px-3 py-1.5 rounded-md font-medium transition flex items-center gap-1.5 ${
                activeTab === "REQUESTS" ? "bg-teal-600 text-white shadow-sm" : "text-slate-400 hover:text-slate-200"
              }`}
            >
              <span>Solicitudes de Ingreso</span>
              <span className="w-2 h-2 rounded-full bg-amber-400" />
            </button>
          </div>
        )}
      </div>

      {/* VISTA 1: Directorio General de Miembros */}
      {activeTab === "DIRECTORY" && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
          {members.map((m) => (
            <div
              key={m.id}
              className="p-5 rounded-2xl bg-slate-900/90 border border-slate-800 hover:border-slate-700 transition space-y-4 flex flex-col justify-between shadow-lg"
            >
              <div className="space-y-3">
                <div className="flex items-start gap-3">
                  <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-teal-500/30 to-slate-800 text-teal-300 font-bold flex items-center justify-center text-base border border-teal-500/30 shrink-0">
                    {m.full_name?.charAt(0) || "U"}
                  </div>
                  <div className="min-w-0">
                    <h3 className="text-sm font-bold text-white truncate">{m.full_name}</h3>
                    <div className="flex items-center gap-1.5 mt-1 flex-wrap">
                      <span className={`text-[10px] font-bold px-2 py-0.5 rounded ${
                        m.role === "DIRECTOR"
                          ? "bg-amber-500/10 text-amber-300 border border-amber-500/20"
                          : m.role === "ADMINISTRADOR"
                          ? "bg-purple-500/10 text-purple-300 border border-purple-500/20"
                          : "bg-teal-500/10 text-teal-300 border border-teal-500/20"
                      }`}>
                        {m.role}
                      </span>
                      {m.can_operate_3d && (
                        <span className="text-[10px] font-bold px-1.5 py-0.5 rounded bg-cyan-500/10 text-cyan-300 border border-cyan-500/20 flex items-center gap-0.5">
                          <Printer className="w-2.5 h-2.5" />
                          <span>Taller 3D</span>
                        </span>
                      )}
                    </div>
                  </div>
                </div>

                <p className="text-xs text-slate-300 leading-relaxed line-clamp-2">
                  {m.bio || "Investigador en robótica y sistemas mecatrónicos de la Universidad de Cartagena."}
                </p>

                <div className="space-y-1 text-xs text-slate-400 border-t border-slate-800/80 pt-2.5">
                  <div className="flex items-center gap-1.5 truncate">
                    <Mail className="w-3.5 h-3.5 text-slate-500 shrink-0" />
                    <span className="truncate">{m.email}</span>
                  </div>
                  {m.student_code && (
                    <div className="flex items-center gap-1.5">
                      <GraduationCap className="w-3.5 h-3.5 text-slate-500 shrink-0" />
                      <span>Cód: {m.student_code}</span>
                    </div>
                  )}
                </div>
              </div>

              {/* Controles de Director */}
              {user?.role === "DIRECTOR" && m.id !== user.id && (
                <div className="pt-2 border-t border-slate-800 flex items-center justify-between text-[11px]">
                  <span className="text-slate-500 font-medium">Permiso 3D:</span>
                  <button
                    onClick={() => handleToggle3D(m.id, !m.can_operate_3d)}
                    className={`px-2 py-0.5 rounded font-semibold transition ${
                      m.can_operate_3d
                        ? "bg-red-500/10 text-red-400 border border-red-500/20 hover:bg-red-500/20"
                        : "bg-cyan-500/10 text-cyan-400 border border-cyan-500/20 hover:bg-cyan-500/20"
                    }`}
                  >
                    {m.can_operate_3d ? "Revocar 3D" : "+ Habilitar 3D"}
                  </button>
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      {/* VISTA 2: Bandeja de Solicitudes de Ingreso (Directores) */}
      {activeTab === "REQUESTS" && (
        <div className="space-y-4">
          <div className="bg-slate-900/90 border border-slate-800 p-5 rounded-2xl space-y-4 shadow-xl">
            <h3 className="text-sm font-bold text-white flex items-center gap-2">
              <UserPlus className="w-4 h-4 text-amber-400" />
              <span>Aspirantes Pendientes de Aprobación</span>
            </h3>

            <div className="space-y-3">
              {requests.map((req) => (
                <div
                  key={req.id}
                  className="p-4 rounded-xl bg-slate-950/60 border border-slate-800 space-y-3"
                >
                  <div className="flex flex-wrap items-start justify-between gap-2">
                    <div>
                      <h4 className="text-sm font-bold text-white">{req.full_name}</h4>
                      <p className="text-xs text-slate-400">
                        {req.email} • Cód: {req.student_code} • {req.career_program || "Ingeniería"}
                      </p>
                    </div>
                    <button
                      onClick={() => setSelectedReq(req)}
                      className="text-xs font-semibold px-3 py-1.5 rounded-lg bg-teal-600 hover:bg-teal-500 text-white transition shadow-sm"
                    >
                      Evaluar Solicitud
                    </button>
                  </div>

                  {req.motivation_letter && (
                    <p className="text-xs text-slate-300 bg-slate-900/80 p-2.5 rounded-lg border border-slate-800 italic">
                      "{req.motivation_letter}"
                    </p>
                  )}
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* Modal: Evaluar Solicitud */}
      {selectedReq && (
        <div className="fixed inset-0 z-50 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-2xl max-w-md w-full p-6 space-y-4 shadow-2xl animate-in fade-in zoom-in-95">
            <div className="flex items-center justify-between">
              <h3 className="text-base font-bold text-white">Resolución de Ingreso</h3>
              <button
                onClick={() => setSelectedReq(null)}
                className="text-slate-400 hover:text-white"
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            <div className="p-3 rounded-lg bg-slate-950/60 border border-slate-800 text-xs space-y-0.5">
              <p className="font-bold text-white">{selectedReq.full_name}</p>
              <p className="text-slate-400">{selectedReq.email} • {selectedReq.career_program}</p>
            </div>

            <div className="space-y-3 text-xs">
              <div>
                <label className="text-[11px] text-slate-400 block mb-1">Acción de Decisión</label>
                <select
                  value={evalAction}
                  onChange={(e) => setEvalAction(e.target.value as "APPROVE" | "REJECT")}
                  className="input-field"
                >
                  <option value="APPROVE">APPROVE (Aprobar e Incorporar al Semillero)</option>
                  <option value="REJECT">REJECT (Rechazar Solicitud)</option>
                </select>
              </div>

              {evalAction === "APPROVE" && (
                <div>
                  <label className="text-[11px] text-slate-400 block mb-1">Rol a Asignar</label>
                  <select
                    value={evalRole}
                    onChange={(e) => setEvalRole(e.target.value)}
                    className="input-field"
                  >
                    <option value="MIEMBRO">MIEMBRO (Estudiante Investigador)</option>
                    <option value="ADMINISTRADOR">ADMINISTRADOR (Gestor de Laboratorio)</option>
                    <option value="DIRECTOR">DIRECTOR (Liderazgo Científico)</option>
                  </select>
                </div>
              )}

              <div>
                <label className="text-[11px] text-slate-400 block mb-1">Observaciones / Feedback al Estudiante</label>
                <textarea
                  placeholder="Mensaje de bienvenida o justificación..."
                  rows={2}
                  value={evalFeedback}
                  onChange={(e) => setEvalFeedback(e.target.value)}
                  className="input-field"
                />
              </div>
            </div>

            <div className="flex gap-2 pt-2">
              <button
                onClick={() => setSelectedReq(null)}
                className="btn-action flex-1 py-2"
              >
                Cancelar
              </button>
              <button
                onClick={handleReviewRequest}
                className="btn-primary flex-1 py-2"
              >
                Guardar Decisión
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
