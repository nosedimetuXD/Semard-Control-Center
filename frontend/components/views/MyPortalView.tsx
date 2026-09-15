"use client";

import { useState, useEffect } from "react";
import {
  User,
  FolderGit2,
  Package,
  Printer,
  Mail,
  GraduationCap,
  Sparkles,
  Clock,
  CheckCircle2,
  LogIn,
  Edit3,
} from "lucide-react";
import { apiCall, getApiBaseUrl } from "@/lib/api";

interface MyPortalViewProps {
  user: any;
  onShowAlert: (msg: string, type: "info" | "success" | "error") => void;
  onTrackResponse: (res: any) => void;
  onProfileUpdated: () => void;
}

export default function MyPortalView({
  user,
  onShowAlert,
  onTrackResponse,
  onProfileUpdated,
}: MyPortalViewProps) {
  const [myProjects, setMyProjects] = useState<any[]>([]);
  const [myLoans, setMyLoans] = useState<any[]>([]);
  const [my3D, setMy3D] = useState<any[]>([]);

  // Edición de Perfil
  const [bio, setBio] = useState("");
  const [studentCode, setStudentCode] = useState("");
  const [isEditing, setIsEditing] = useState(false);

  useEffect(() => {
    if (user) {
      setBio(user.bio || "");
      setStudentCode(user.student_code || "");
      loadUserData();
    }
  }, [user]);

  async function loadUserData() {
    // 1. Mis Proyectos
    const pRes = await apiCall("/api/v1/projects/my-projects", { requireAuth: true });
    if (pRes.ok && Array.isArray(pRes.data)) setMyProjects(pRes.data);

    // 2. Mis Préstamos
    const lRes = await apiCall("/api/v1/loans/my-loans", { requireAuth: true });
    if (lRes.ok && Array.isArray(lRes.data)) setMyLoans(lRes.data);

    // 3. Mis Piezas 3D
    const dRes = await apiCall("/api/v1/print3d/my-requests", { requireAuth: true });
    if (dRes.ok && Array.isArray(dRes.data)) setMy3D(dRes.data);
  }

  async function handleUpdateProfile() {
    const res = await apiCall("/api/v1/auth/me", {
      method: "PUT",
      body: { bio, student_code: studentCode },
      requireAuth: true,
    });
    onTrackResponse(res);
    if (res.ok) {
      onShowAlert("Perfil actualizado correctamente", "success");
      setIsEditing(false);
      onProfileUpdated();
    } else {
      onShowAlert(res.data?.error || "Error al actualizar perfil", "error");
    }
  }

  function handleGoogleLogin() {
    const base = getApiBaseUrl();
    const returnTo = `${window.location.origin}/auth/callback`;
    window.location.href = `${base}/api/v1/auth/google/login?return_to=${encodeURIComponent(
      returnTo
    )}&redirect=true`;
  }

  if (!user) {
    return (
      <div className="p-8 rounded-2xl bg-slate-900 border border-slate-800 text-center max-w-lg mx-auto space-y-4 my-10 shadow-xl">
        <div className="w-12 h-12 rounded-xl bg-teal-500/10 text-teal-400 flex items-center justify-center mx-auto border border-teal-500/20">
          <User className="w-6 h-6" />
        </div>
        <h3 className="text-base font-bold text-white">Acceso a Mi Portal de Miembro</h3>
        <p className="text-xs text-slate-400">
          Inicia sesión con tu correo institucional de la Universidad de Cartagena para consultar tus proyectos asignados, préstamos activos y piezas 3D.
        </p>
        <button
          onClick={handleGoogleLogin}
          className="btn-primary py-2 px-4 inline-flex items-center gap-2"
        >
          <LogIn className="w-4 h-4" />
          <span>Iniciar Sesión con Google UdeC</span>
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Tarjeta de Ficha de Usuario */}
      <div className="p-6 rounded-2xl bg-slate-900/90 border border-slate-800 shadow-xl space-y-4">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div className="flex items-center gap-4">
            {user.avatar_url ? (
              <img
                src={user.avatar_url}
                alt={user.full_name}
                className="w-16 h-16 rounded-2xl border-2 border-teal-500/40 object-cover shadow-lg"
              />
            ) : (
              <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-teal-500 to-emerald-700 text-white font-bold flex items-center justify-center text-xl border border-teal-400/30">
                {user.full_name?.charAt(0) || "U"}
              </div>
            )}
            <div>
              <div className="flex items-center gap-2 flex-wrap">
                <h2 className="text-lg font-bold text-white">{user.full_name}</h2>
                <span className="text-[10px] font-bold px-2 py-0.5 rounded-full bg-teal-500/10 text-teal-300 border border-teal-500/20">
                  {user.role}
                </span>
                {user.can_operate_3d && (
                  <span className="text-[10px] font-bold px-2 py-0.5 rounded-full bg-cyan-500/10 text-cyan-300 border border-cyan-500/20">
                    Operador 3D
                  </span>
                )}
              </div>
              <p className="text-xs text-slate-400 mt-1">{user.email}</p>
              <p className="text-[11px] text-slate-500 font-mono mt-0.5">
                Código Estudiantil: {user.student_code || "Sin registrar"}
              </p>
            </div>
          </div>

          <button
            onClick={() => setIsEditing(!isEditing)}
            className="self-start sm:self-auto text-xs font-semibold px-3 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 border border-slate-700 transition flex items-center gap-1.5"
          >
            <Edit3 className="w-3.5 h-3.5" />
            <span>{isEditing ? "Cerrar Edición" : "Editar Datos"}</span>
          </button>
        </div>

        {/* Panel de Edición Rápida */}
        {isEditing ? (
          <div className="pt-4 border-t border-slate-800 space-y-3 bg-slate-950/60 p-4 rounded-xl">
            <h4 className="text-xs font-bold text-teal-400 uppercase tracking-wider">
              Actualizar Biografía y Código
            </h4>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 text-xs">
              <input
                type="text"
                placeholder="Código Estudiantil"
                value={studentCode}
                onChange={(e) => setStudentCode(e.target.value)}
                className="input-field"
              />
              <input
                type="text"
                placeholder="Biografía breve / Enfoque de investigación"
                value={bio}
                onChange={(e) => setBio(e.target.value)}
                className="input-field"
              />
            </div>
            <button
              onClick={handleUpdateProfile}
              className="btn-primary text-xs py-2 px-4"
            >
              Guardar Cambios
            </button>
          </div>
        ) : (
          <div className="pt-3 border-t border-slate-800/80 text-xs text-slate-300">
            <span className="text-slate-500 block text-[11px] mb-0.5">Perfil Académico:</span>
            {user.bio || "Investigador en el Semillero de Robótica y Automatización (SEMARD) — Universidad de Cartagena."}
          </div>
        )}
      </div>

      {/* Grid de Tres Módulos del Miembro */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-5">
        {/* 1. Mis Proyectos */}
        <div className="p-5 rounded-2xl bg-slate-900/90 border border-slate-800 space-y-3 shadow-lg">
          <h3 className="text-xs font-bold text-white flex items-center gap-2 uppercase tracking-wider text-teal-400">
            <FolderGit2 className="w-4 h-4" />
            <span>Mis Proyectos Asignados</span>
          </h3>

          {myProjects.length === 0 ? (
            <div className="p-4 rounded-xl bg-slate-950/40 border border-slate-800/80 text-center text-xs text-slate-500">
              No tienes proyectos vinculados actualmente.
            </div>
          ) : (
            <div className="space-y-2">
              {myProjects.map((p) => (
                <div key={p.id} className="p-3 rounded-lg bg-slate-950/60 border border-slate-800 text-xs">
                  <p className="font-bold text-white">{p.title}</p>
                  <p className="text-[11px] text-teal-400">{p.research_line}</p>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* 2. Mis Préstamos de Equipos */}
        <div className="p-5 rounded-2xl bg-slate-900/90 border border-slate-800 space-y-3 shadow-lg">
          <h3 className="text-xs font-bold text-white flex items-center gap-2 uppercase tracking-wider text-indigo-400">
            <Package className="w-4 h-4" />
            <span>Mis Préstamos de Laboratorio</span>
          </h3>

          {myLoans.length === 0 ? (
            <div className="p-4 rounded-xl bg-slate-950/40 border border-slate-800/80 text-center text-xs text-slate-500">
              No tienes préstamos activos o pendientes.
            </div>
          ) : (
            <div className="space-y-2">
              {myLoans.map((l) => (
                <div key={l.id} className="p-3 rounded-lg bg-slate-950/60 border border-slate-800 text-xs flex justify-between items-center">
                  <div>
                    <p className="font-bold text-white">{l.item_name || "Equipo de Lab"}</p>
                    <p className="text-[11px] text-slate-400">Estado: {l.status}</p>
                  </div>
                  <span className="text-[10px] px-2 py-0.5 rounded bg-slate-800 text-slate-300">
                    {l.status}
                  </span>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* 3. Mis Piezas 3D */}
        <div className="p-5 rounded-2xl bg-slate-900/90 border border-slate-800 space-y-3 shadow-lg">
          <h3 className="text-xs font-bold text-white flex items-center gap-2 uppercase tracking-wider text-cyan-400">
            <Printer className="w-4 h-4" />
            <span>Mis Piezas en Taller 3D</span>
          </h3>

          {my3D.length === 0 ? (
            <div className="p-4 rounded-xl bg-slate-950/40 border border-slate-800/80 text-center text-xs text-slate-500">
              No has enviado piezas a manufacturar.
            </div>
          ) : (
            <div className="space-y-2">
              {my3D.map((d) => (
                <div key={d.id} className="p-3 rounded-lg bg-slate-950/60 border border-slate-800 text-xs flex justify-between items-center">
                  <div className="min-w-0">
                    <p className="font-bold text-white truncate">{d.filename || "Pieza 3D"}</p>
                    <p className="text-[11px] text-slate-400">{d.material} • {d.color}</p>
                  </div>
                  <span className="text-[10px] px-2 py-0.5 rounded bg-cyan-500/20 text-cyan-300 font-bold shrink-0">
                    {d.status}
                  </span>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
