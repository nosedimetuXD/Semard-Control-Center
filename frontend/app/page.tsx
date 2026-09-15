"use client";

import { useState, useEffect, Suspense } from "react";
import { useSearchParams } from "next/navigation";
import Header from "@/components/Header";
import Inspector from "@/components/Inspector";
import { apiCall, ApiResponse, getAuthToken } from "@/lib/api";
import {
  Globe,
  Calendar,
  UserPlus,
  User,
  Users,
  FolderGit2,
  FileCheck2,
  BadgeDollarSign,
  Package,
  ArrowLeftRight,
  Printer,
  CheckCircle,
  AlertCircle,
  X,
} from "lucide-react";

type TabKey =
  | "hub"
  | "events"
  | "registration"
  | "profile"
  | "members"
  | "projects"
  | "updates"
  | "resources"
  | "inventory"
  | "loans"
  | "print3d";

function DashboardContent() {
  const searchParams = useSearchParams();
  const [activeTab, setActiveTab] = useState<TabKey>("hub");
  const [lastResponse, setLastResponse] = useState<ApiResponse | null>(null);
  const [alert, setAlert] = useState<{ msg: string; type: "info" | "success" | "error" } | null>(null);

  // Form states
  // Eventos
  const [evTitle, setEvTitle] = useState("");
  const [evDesc, setEvDesc] = useState("");
  const [evVis, setEvVis] = useState("PUBLIC");
  const [evLoc, setEvLoc] = useState("");
  const [evStart, setEvStart] = useState("");
  const [regEvId, setRegEvId] = useState("");
  const [regName, setRegName] = useState("");
  const [regEmail, setRegEmail] = useState("");
  const [regPhone, setRegPhone] = useState("");

  // Solicitud Ingreso
  const [newReqGoogleId, setNewReqGoogleId] = useState("");
  const [newReqEmail, setNewReqEmail] = useState("");
  const [newReqName, setNewReqName] = useState("");
  const [newReqCode, setNewReqCode] = useState("");
  const [newReqProgram, setNewReqProgram] = useState("");
  const [newReqLetter, setNewReqLetter] = useState("");

  // Perfil
  const [profBio, setProfBio] = useState("");
  const [profCode, setProfCode] = useState("");

  // Miembros & Roles
  const [targetUserId, setTargetUserId] = useState("");
  const [newRole, setNewRole] = useState("MIEMBRO");
  const [revReqId, setRevReqId] = useState("");
  const [revReqAction, setRevReqAction] = useState("APPROVE");
  const [revReqRole, setRevReqRole] = useState("MIEMBRO");
  const [revReqFeedback, setRevReqFeedback] = useState("");

  // Proyectos
  const [projTitle, setProjTitle] = useState("");
  const [projDesc, setProjDesc] = useState("");
  const [projLine, setProjLine] = useState("");
  const [projStart, setProjStart] = useState("");
  const [assignProjId, setAssignProjId] = useState("");
  const [assignUserId, setAssignUserId] = useState("");
  const [assignIsLead, setAssignIsLead] = useState(false);

  // Avances
  const [updProjId, setUpdProjId] = useState("");
  const [updTitle, setUpdTitle] = useState("");
  const [updContent, setUpdContent] = useState("");
  const [updLinks, setUpdLinks] = useState("");
  const [revUpdId, setRevUpdId] = useState("");
  const [revUpdStatus, setRevUpdStatus] = useState("APPROVED");
  const [revUpdFeedback, setRevUpdFeedback] = useState("");

  // Recursos
  const [resProjId, setResProjId] = useState("");
  const [resType, setResType] = useState("MATERIAL");
  const [resTitle, setResTitle] = useState("");
  const [resDesc, setResDesc] = useState("");
  const [resCost, setResCost] = useState("0");
  const [evalResId, setEvalResId] = useState("");
  const [evalResStatus, setEvalResStatus] = useState("APPROVED");
  const [evalResFeedback, setEvalResFeedback] = useState("");

  // Inventario
  const [invCode, setInvCode] = useState("");
  const [invName, setInvName] = useState("");
  const [invCat, setInvCat] = useState("EQUIPAMIENTO");
  const [invStock, setInvStock] = useState("1");
  const [invLoc, setInvLoc] = useState("");

  // Préstamos
  const [loanItemId, setLoanItemId] = useState("");
  const [loanReason, setLoanReason] = useState("");
  const [loanStart, setLoanStart] = useState("");
  const [loanEnd, setLoanEnd] = useState("");
  const [revLoanId, setRevLoanId] = useState("");
  const [revLoanStatus, setRevLoanStatus] = useState("APPROVED");
  const [revLoanModEnd, setRevLoanModEnd] = useState("");
  const [revLoanFeedback, setRevLoanFeedback] = useState("");

  // Taller 3D
  const [file3d, setFile3d] = useState<File | null>(null);
  const [mat3d, setMat3d] = useState("PLA");
  const [color3d, setColor3d] = useState("Negro");
  const [infill3d, setInfill3d] = useState("20");
  const [notes3d, setNotes3d] = useState("");
  const [rev3dId, setRev3dId] = useState("");
  const [rev3dStatus, setRev3dStatus] = useState("APPROVED");
  const [op3dId, setOp3dId] = useState("");
  const [op3dStatus, setOp3dStatus] = useState("IN_PROGRESS");

  useEffect(() => {
    if (searchParams.get("auth_success") === "true") {
      showAlert("¡Autenticación con Google exitosa! Credencial JWT activa.", "success");
      window.history.replaceState({}, document.title, window.location.pathname);
    } else if (searchParams.get("status") === "NEEDS_REGISTRATION") {
      setActiveTab("registration");
      setNewReqEmail(searchParams.get("email") || "");
      setNewReqName(searchParams.get("name") || "");
      setNewReqGoogleId(searchParams.get("google_id") || "");
      showAlert("Correo institucional verificado. Completa tu código estudiantil para solicitar membresía.", "info");
      window.history.replaceState({}, document.title, window.location.pathname);
    } else if (searchParams.get("status") === "PENDING_APPROVAL") {
      showAlert("Tu solicitud de ingreso ya está radicada y en proceso de revisión por los Directores.", "info");
      window.history.replaceState({}, document.title, window.location.pathname);
    } else if (searchParams.get("status") === "REJECTED") {
      const msg = searchParams.get("msg") || "Tu solicitud de ingreso previa fue rechazada.";
      showAlert(msg, "error");
      window.history.replaceState({}, document.title, window.location.pathname);
    }
  }, [searchParams]);

  function showAlert(msg: string, type: "info" | "success" | "error") {
    setAlert({ msg, type });
  }

  async function executeCall(
    endpoint: string,
    method: "GET" | "POST" | "PUT" | "PATCH" | "DELETE" = "GET",
    body?: any,
    requireAuth: boolean = false,
    isFormData: boolean = false
  ) {
    if (requireAuth && !getAuthToken()) {
      showAlert("Esta ruta requiere iniciar sesión con Google o pegar un JWT primero.", "error");
      return;
    }

    const res = await apiCall(endpoint, { method, body, requireAuth, isFormData });
    setLastResponse(res);

    if (res.ok) {
      showAlert(`Petición ${method} ${endpoint} ejecutada con éxito (${res.status})`, "success");
    } else {
      const errMsg = res.data?.error || res.rawText || "Error en la petición";
      showAlert(`Error ${res.status}: ${errMsg}`, "error");
    }
  }

  const tabs: { key: TabKey; label: string; icon: any }[] = [
    { key: "hub", label: "Hub Informativo", icon: Globe },
    { key: "events", label: "Eventos", icon: Calendar },
    { key: "registration", label: "Solicitud Ingreso", icon: UserPlus },
    { key: "profile", label: "Mi Perfil (/me)", icon: User },
    { key: "members", label: "Miembros & Roles", icon: Users },
    { key: "projects", label: "Proyectos", icon: FolderGit2 },
    { key: "updates", label: "Avances", icon: FileCheck2 },
    { key: "resources", label: "Recursos (3 Estados)", icon: BadgeDollarSign },
    { key: "inventory", label: "Inventario", icon: Package },
    { key: "loans", label: "Préstamos", icon: ArrowLeftRight },
    { key: "print3d", label: "Taller Impresión 3D", icon: Printer },
  ];

  return (
    <div className="min-h-screen flex flex-col bg-slate-950 text-slate-100">
      <Header
        onSessionChange={() => {}}
        onShowAlert={(msg, type) => showAlert(msg, type)}
      />

      {/* Notificación Alert */}
      {alert && (
        <div className="max-w-7xl mx-auto w-full px-4 pt-4">
          <div
            className={`p-3 rounded-lg text-xs border flex items-center justify-between ${
              alert.type === "error"
                ? "bg-red-500/10 text-red-300 border-red-500/30"
                : alert.type === "success"
                ? "bg-emerald-500/10 text-emerald-300 border-emerald-500/30"
                : "bg-slate-800 text-slate-200 border-slate-700"
            }`}
          >
            <div className="flex items-center gap-2">
              {alert.type === "success" ? (
                <CheckCircle className="w-4 h-4 text-emerald-400 shrink-0" />
              ) : alert.type === "error" ? (
                <AlertCircle className="w-4 h-4 text-red-400 shrink-0" />
              ) : (
                <AlertCircle className="w-4 h-4 text-teal-400 shrink-0" />
              )}
              <span>{alert.msg}</span>
            </div>
            <button
              onClick={() => setAlert(null)}
              className="text-slate-400 hover:text-slate-200 ml-2"
            >
              <X className="w-3.5 h-3.5" />
            </button>
          </div>
        </div>
      )}

      {/* Contenedor Principal: Tabs + Paneles + Consola */}
      <div className="max-w-7xl mx-auto w-full px-4 py-6 flex flex-col lg:flex-row gap-6 flex-1">
        {/* Navegación Lateral de Módulos */}
        <nav className="w-full lg:w-60 shrink-0 flex flex-row lg:flex-col gap-1 overflow-x-auto lg:overflow-visible pb-2 lg:pb-0 text-sm">
          {tabs.map((t) => {
            const Icon = t.icon;
            const active = activeTab === t.key;
            return (
              <button
                key={t.key}
                onClick={() => setActiveTab(t.key)}
                className={`text-left px-3 py-2 rounded-lg font-medium transition flex items-center gap-2.5 whitespace-nowrap text-xs ${
                  active
                    ? "bg-teal-600 text-white shadow-sm"
                    : "text-slate-400 hover:bg-slate-900 hover:text-slate-200"
                }`}
              >
                <Icon className={`w-4 h-4 ${active ? "text-white" : "text-slate-400"}`} />
                <span>{t.label}</span>
              </button>
            );
          })}
        </nav>

        {/* Panel de Trabajo + Consola Inspector */}
        <main className="flex-1 min-w-0 flex flex-col gap-6">
          {/* SECCIÓN 1: HUB INFORMATIVO */}
          {activeTab === "hub" && (
            <div className="bg-slate-900 p-5 rounded-xl border border-slate-800 space-y-4">
              <div>
                <h2 className="text-base font-bold text-white flex items-center gap-2">
                  <Globe className="w-4 h-4 text-teal-400" />
                  <span>Hub Institucional (Público)</span>
                </h2>
                <p className="text-xs text-slate-400 mt-1">Endpoints abiertos para divulgación científica y consulta institucional.</p>
              </div>

              <div className="flex flex-wrap gap-2">
                <button onClick={() => executeCall("/api/v1/hub/info")} className="btn-action">GET /hub/info</button>
                <button onClick={() => executeCall("/api/v1/hub/directors")} className="btn-action">GET /hub/directors</button>
                <button onClick={() => executeCall("/api/v1/hub/projects")} className="btn-action">GET /hub/projects (Portafolio Concluido)</button>
                <button onClick={() => executeCall("/healthz")} className="btn-action">GET /healthz</button>
              </div>
            </div>
          )}

          {/* SECCIÓN 2: EVENTOS */}
          {activeTab === "events" && (
            <div className="bg-slate-900 p-5 rounded-xl border border-slate-800 space-y-4">
              <h2 className="text-base font-bold text-white flex items-center gap-2">
                <Calendar className="w-4 h-4 text-teal-400" />
                <span>Cartelera de Eventos</span>
              </h2>

              <div className="flex flex-wrap gap-2">
                <button onClick={() => executeCall("/api/v1/events/public")} className="btn-action">GET /events/public</button>
                <button onClick={() => executeCall("/api/v1/events/internal", "GET", null, true)} className="btn-action">GET /events/internal (Protegido)</button>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pt-4 border-t border-slate-800">
                {/* Crear Evento */}
                <div className="space-y-3 bg-slate-950/60 p-4 rounded-lg border border-slate-800/80">
                  <h3 className="text-xs font-bold uppercase tracking-wider text-teal-400">Crear Evento (Directores)</h3>
                  <input
                    type="text"
                    placeholder="Título del evento"
                    value={evTitle}
                    onChange={(e) => setEvTitle(e.target.value)}
                    className="input-field"
                  />
                  <textarea
                    placeholder="Descripción detallada"
                    rows={2}
                    value={evDesc}
                    onChange={(e) => setEvDesc(e.target.value)}
                    className="input-field"
                  />
                  <div className="grid grid-cols-2 gap-2">
                    <select value={evVis} onChange={(e) => setEvVis(e.target.value)} className="input-field">
                      <option value="PUBLIC">PUBLIC (Abierto)</option>
                      <option value="INTERNAL">INTERNAL (Semillero)</option>
                    </select>
                    <input
                      type="text"
                      placeholder="Lugar / Enlace"
                      value={evLoc}
                      onChange={(e) => setEvLoc(e.target.value)}
                      className="input-field"
                    />
                  </div>
                  <input
                    type="datetime-local"
                    value={evStart}
                    onChange={(e) => setEvStart(e.target.value)}
                    className="input-field"
                  />
                  <button
                    onClick={() => executeCall("/api/v1/events", "POST", {
                      title: evTitle,
                      description: evDesc,
                      visibility: evVis,
                      location: evLoc,
                      start_time: evStart ? new Date(evStart).toISOString() : new Date().toISOString(),
                    }, true)}
                    className="btn-primary w-full"
                  >
                    Publicar Evento
                  </button>
                </div>

                {/* Inscripción a Evento */}
                <div className="space-y-3 bg-slate-950/60 p-4 rounded-lg border border-slate-800/80">
                  <h3 className="text-xs font-bold uppercase tracking-wider text-indigo-400">Inscripción Pública a Evento</h3>
                  <input
                    type="text"
                    placeholder="UUID del Evento"
                    value={regEvId}
                    onChange={(e) => setRegEvId(e.target.value)}
                    className="input-field"
                  />
                  <input
                    type="text"
                    placeholder="Nombre Completo"
                    value={regName}
                    onChange={(e) => setRegName(e.target.value)}
                    className="input-field"
                  />
                  <input
                    type="email"
                    placeholder="Correo Electrónico"
                    value={regEmail}
                    onChange={(e) => setRegEmail(e.target.value)}
                    className="input-field"
                  />
                  <input
                    type="text"
                    placeholder="Teléfono (opcional)"
                    value={regPhone}
                    onChange={(e) => setRegPhone(e.target.value)}
                    className="input-field"
                  />
                  <button
                    onClick={() => {
                      if (!regEvId) return showAlert("Indica el UUID del evento", "error");
                      executeCall(`/api/v1/events/${regEvId}/register`, "POST", {
                        full_name: regName,
                        email: regEmail,
                        phone: regPhone,
                      });
                    }}
                    className="btn-primary w-full bg-indigo-600 hover:bg-indigo-500"
                  >
                    Inscribirse al Evento
                  </button>
                </div>
              </div>
            </div>
          )}

          {/* SECCIÓN 3: SOLICITUD DE INGRESO */}
          {activeTab === "registration" && (
            <div className="bg-slate-900 p-5 rounded-xl border border-slate-800 space-y-4">
              <div>
                <h2 className="text-base font-bold text-white flex items-center gap-2">
                  <UserPlus className="w-4 h-4 text-teal-400" />
                  <span>Solicitud de Ingreso al Semillero</span>
                </h2>
                <p className="text-xs text-slate-400 mt-1">
                  Vinculación de nuevos estudiantes de la Universidad de Cartagena (POST /api/v1/auth/register-request).
                </p>
              </div>

              <div className="space-y-3 bg-slate-950/60 p-4 rounded-lg border border-slate-800/80 max-w-lg">
                <input
                  type="text"
                  placeholder="Google ID (rellenado tras login)"
                  value={newReqGoogleId}
                  onChange={(e) => setNewReqGoogleId(e.target.value)}
                  className="input-field"
                />
                <input
                  type="email"
                  placeholder="Correo institucional (@unicartagena.edu.co)"
                  value={newReqEmail}
                  onChange={(e) => setNewReqEmail(e.target.value)}
                  className="input-field"
                />
                <input
                  type="text"
                  placeholder="Nombre completo"
                  value={newReqName}
                  onChange={(e) => setNewReqName(e.target.value)}
                  className="input-field"
                />
                <input
                  type="text"
                  placeholder="Código Estudiantil (ej. 0222410043)"
                  value={newReqCode}
                  onChange={(e) => setNewReqCode(e.target.value)}
                  className="input-field"
                />
                <input
                  type="text"
                  placeholder="Programa Académico (ej. Ing. de Sistemas)"
                  value={newReqProgram}
                  onChange={(e) => setNewReqProgram(e.target.value)}
                  className="input-field"
                />
                <textarea
                  placeholder="Carta breve de motivos de ingreso al semillero"
                  rows={2}
                  value={newReqLetter}
                  onChange={(e) => setNewReqLetter(e.target.value)}
                  className="input-field"
                />
                <button
                  onClick={() => {
                    if (!newReqEmail || !newReqCode || !newReqName) {
                      return showAlert("Correo, nombre y código estudiantil son requeridos", "error");
                    }
                    executeCall("/api/v1/auth/register-request", "POST", {
                      google_id: newReqGoogleId,
                      email: newReqEmail,
                      full_name: newReqName,
                      student_code: newReqCode,
                      career_program: newReqProgram || null,
                      motivation_letter: newReqLetter || null,
                    });
                  }}
                  className="btn-primary w-full"
                >
                  Enviar Solicitud de Membresía
                </button>
              </div>
            </div>
          )}

          {/* SECCIÓN 4: MI PERFIL */}
          {activeTab === "profile" && (
            <div className="bg-slate-900 p-5 rounded-xl border border-slate-800 space-y-4">
              <div className="flex items-center justify-between">
                <h2 className="text-base font-bold text-white flex items-center gap-2">
                  <User className="w-4 h-4 text-teal-400" />
                  <span>Mi Perfil (/api/v1/auth/me)</span>
                </h2>
                <button onClick={() => executeCall("/api/v1/auth/me", "GET", null, true)} className="btn-action">
                  Consultar Perfil
                </button>
              </div>

              <div className="space-y-3 bg-slate-950/60 p-4 rounded-lg border border-slate-800/80 max-w-md">
                <h3 className="text-xs font-bold text-teal-400 uppercase">Actualizar Mis Datos</h3>
                <input
                  type="text"
                  placeholder="Biografía / Área de enfoque"
                  value={profBio}
                  onChange={(e) => setProfBio(e.target.value)}
                  className="input-field"
                />
                <input
                  type="text"
                  placeholder="Código Estudiantil"
                  value={profCode}
                  onChange={(e) => setProfCode(e.target.value)}
                  className="input-field"
                />
                <button
                  onClick={() => executeCall("/api/v1/auth/me", "PUT", { bio: profBio, student_code: profCode }, true)}
                  className="btn-primary w-full"
                >
                  Guardar Cambios
                </button>
              </div>
            </div>
          )}

          {/* SECCIÓN 5: MIEMBROS & ROLES */}
          {activeTab === "members" && (
            <div className="bg-slate-900 p-5 rounded-xl border border-slate-800 space-y-4">
              <h2 className="text-base font-bold text-white flex items-center gap-2">
                <Users className="w-4 h-4 text-teal-400" />
                <span>Directorio de Miembros y Gestión de Roles</span>
              </h2>

              <div className="flex flex-wrap gap-2">
                <button onClick={() => executeCall("/api/v1/users", "GET", null, true)} className="btn-action">GET /users (Todos)</button>
                <button onClick={() => executeCall("/api/v1/users?role=DIRECTOR", "GET", null, true)} className="btn-action">Directores</button>
                <button onClick={() => executeCall("/api/v1/users?role=MIEMBRO", "GET", null, true)} className="btn-action">Miembros</button>
                <button onClick={() => executeCall("/api/v1/directors/registration-requests", "GET", null, true)} className="btn-action">Solicitudes Nuevas</button>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pt-2">
                {/* Modificar Rol o Permisos */}
                <div className="bg-slate-950/60 p-4 rounded-lg border border-slate-800/80 space-y-3">
                  <h3 className="text-xs font-bold uppercase tracking-wider text-teal-400">Modificar Rol o Permisos (Directores)</h3>
                  <input
                    type="text"
                    placeholder="UUID del Usuario objetivo"
                    value={targetUserId}
                    onChange={(e) => setTargetUserId(e.target.value)}
                    className="input-field"
                  />
                  <div className="flex gap-2">
                    <select value={newRole} onChange={(e) => setNewRole(e.target.value)} className="input-field">
                      <option value="MIEMBRO">MIEMBRO</option>
                      <option value="ADMINISTRADOR">ADMINISTRADOR</option>
                      <option value="DIRECTOR">DIRECTOR</option>
                    </select>
                    <button
                      onClick={() => {
                        if (!targetUserId) return showAlert("Indica el UUID del usuario", "error");
                        executeCall(`/api/v1/directors/users/${targetUserId}/role`, "PATCH", { role: newRole }, true);
                      }}
                      className="btn-primary"
                    >
                      Cambiar
                    </button>
                  </div>
                  <div className="flex items-center gap-2 pt-1">
                    <button
                      onClick={() => {
                        if (!targetUserId) return showAlert("Indica el UUID del usuario", "error");
                        executeCall(`/api/v1/directors/users/${targetUserId}/permissions`, "PATCH", { can_operate_3d: true }, true);
                      }}
                      className="text-[11px] px-2.5 py-1 rounded bg-emerald-600/20 text-emerald-400 border border-emerald-500/30 hover:bg-emerald-600/30"
                    >
                      + Habilitar Operador 3D
                    </button>
                    <button
                      onClick={() => {
                        if (!targetUserId) return showAlert("Indica el UUID del usuario", "error");
                        executeCall(`/api/v1/directors/users/${targetUserId}/permissions`, "PATCH", { can_operate_3d: false }, true);
                      }}
                      className="text-[11px] px-2.5 py-1 rounded bg-red-600/20 text-red-400 border border-red-500/30 hover:bg-red-600/30"
                    >
                      - Revocar 3D
                    </button>
                  </div>
                </div>

                {/* Evaluar Solicitud de Ingreso */}
                <div className="bg-slate-950/60 p-4 rounded-lg border border-slate-800/80 space-y-3">
                  <h3 className="text-xs font-bold uppercase tracking-wider text-teal-400">Evaluar Solicitud de Ingreso</h3>
                  <input
                    type="text"
                    placeholder="UUID de la Solicitud"
                    value={revReqId}
                    onChange={(e) => setRevReqId(e.target.value)}
                    className="input-field"
                  />
                  <div className="grid grid-cols-2 gap-2">
                    <select value={revReqAction} onChange={(e) => setRevReqAction(e.target.value)} className="input-field">
                      <option value="APPROVE">APPROVE (Aprobar)</option>
                      <option value="REJECT">REJECT (Rechazar)</option>
                    </select>
                    <select value={revReqRole} onChange={(e) => setRevReqRole(e.target.value)} className="input-field">
                      <option value="MIEMBRO">MIEMBRO</option>
                      <option value="ADMINISTRADOR">ADMINISTRADOR</option>
                      <option value="DIRECTOR">DIRECTOR</option>
                    </select>
                  </div>
                  <input
                    type="text"
                    placeholder="Feedback / Observaciones"
                    value={revReqFeedback}
                    onChange={(e) => setRevReqFeedback(e.target.value)}
                    className="input-field"
                  />
                  <button
                    onClick={() => {
                      if (!revReqId) return showAlert("Indica el UUID de la solicitud", "error");
                      executeCall(`/api/v1/directors/registration-requests/${revReqId}/review`, "POST", {
                        action: revReqAction,
                        role: revReqRole,
                        feedback: revReqFeedback || null,
                      }, true);
                    }}
                    className="btn-primary w-full"
                  >
                    Procesar Evaluación
                  </button>
                </div>
              </div>
            </div>
          )}

          {/* SECCIÓN 6: PROYECTOS */}
          {activeTab === "projects" && (
            <div className="bg-slate-900 p-5 rounded-xl border border-slate-800 space-y-4">
              <h2 className="text-base font-bold text-white flex items-center gap-2">
                <FolderGit2 className="w-4 h-4 text-teal-400" />
                <span>Módulo de Proyectos</span>
              </h2>

              <div className="flex flex-wrap gap-2">
                <button onClick={() => executeCall("/api/v1/projects/showcase", "GET", null, true)} className="btn-action">GET /projects/showcase</button>
                <button onClick={() => executeCall("/api/v1/projects/my-projects", "GET", null, true)} className="btn-action">GET /projects/my-projects</button>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pt-2">
                {/* Crear Proyecto */}
                <div className="space-y-3 bg-slate-950/60 p-4 rounded-lg border border-slate-800/80">
                  <h3 className="text-xs font-bold uppercase tracking-wider text-teal-400">Crear Proyecto (Directores)</h3>
                  <input
                    type="text"
                    placeholder="Título del proyecto"
                    value={projTitle}
                    onChange={(e) => setProjTitle(e.target.value)}
                    className="input-field"
                  />
                  <textarea
                    placeholder="Descripción y objetivos"
                    rows={2}
                    value={projDesc}
                    onChange={(e) => setProjDesc(e.target.value)}
                    className="input-field"
                  />
                  <input
                    type="text"
                    placeholder="Línea de investigación"
                    value={projLine}
                    onChange={(e) => setProjLine(e.target.value)}
                    className="input-field"
                  />
                  <input
                    type="date"
                    value={projStart}
                    onChange={(e) => setProjStart(e.target.value)}
                    className="input-field"
                  />
                  <button
                    onClick={() => executeCall("/api/v1/projects", "POST", {
                      title: projTitle,
                      description: projDesc,
                      research_line: projLine,
                      start_date: projStart ? new Date(projStart).toISOString() : new Date().toISOString(),
                    }, true)}
                    className="btn-primary w-full"
                  >
                    Crear Proyecto
                  </button>
                </div>

                {/* Asignar Miembro */}
                <div className="space-y-3 bg-slate-950/60 p-4 rounded-lg border border-slate-800/80">
                  <h3 className="text-xs font-bold uppercase tracking-wider text-teal-400">Asignar Integrante</h3>
                  <input
                    type="text"
                    placeholder="UUID del Proyecto"
                    value={assignProjId}
                    onChange={(e) => setAssignProjId(e.target.value)}
                    className="input-field"
                  />
                  <input
                    type="text"
                    placeholder="UUID del Usuario"
                    value={assignUserId}
                    onChange={(e) => setAssignUserId(e.target.value)}
                    className="input-field"
                  />
                  <label className="flex items-center gap-2 text-xs text-slate-300">
                    <input
                      type="checkbox"
                      checked={assignIsLead}
                      onChange={(e) => setAssignIsLead(e.target.checked)}
                      className="rounded bg-slate-900 border-slate-700 text-teal-600 focus:ring-teal-500"
                    />
                    <span>Designar como Encargado Líder (is_lead)</span>
                  </label>
                  <button
                    onClick={() => {
                      if (!assignProjId || !assignUserId) return showAlert("Indica IDs de proyecto y usuario", "error");
                      executeCall(`/api/v1/projects/${assignProjId}/members`, "POST", {
                        user_id: assignUserId,
                        is_lead: assignIsLead,
                      }, true);
                    }}
                    className="btn-primary w-full"
                  >
                    Vincular al Proyecto
                  </button>
                </div>
              </div>
            </div>
          )}

          {/* SECCIÓN 7: AVANCES & FEEDBACK */}
          {activeTab === "updates" && (
            <div className="bg-slate-900 p-5 rounded-xl border border-slate-800 space-y-4">
              <h2 className="text-base font-bold text-white flex items-center gap-2">
                <FileCheck2 className="w-4 h-4 text-teal-400" />
                <span>Avances de Proyecto & Feedback Formativo</span>
              </h2>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {/* Enviar Avance */}
                <div className="space-y-3 bg-slate-950/60 p-4 rounded-lg border border-slate-800/80">
                  <h3 className="text-xs font-bold uppercase tracking-wider text-teal-400">Registrar Nuevo Avance</h3>
                  <input
                    type="text"
                    placeholder="UUID del Proyecto"
                    value={updProjId}
                    onChange={(e) => setUpdProjId(e.target.value)}
                    className="input-field"
                  />
                  <input
                    type="text"
                    placeholder="Título del avance"
                    value={updTitle}
                    onChange={(e) => setUpdTitle(e.target.value)}
                    className="input-field"
                  />
                  <textarea
                    placeholder="Contenido descriptivo del avance"
                    rows={2}
                    value={updContent}
                    onChange={(e) => setUpdContent(e.target.value)}
                    className="input-field"
                  />
                  <input
                    type="text"
                    placeholder="Enlaces separados por coma (GitHub, Drive)"
                    value={updLinks}
                    onChange={(e) => setUpdLinks(e.target.value)}
                    className="input-field"
                  />
                  <div className="flex gap-2">
                    <button
                      onClick={() => {
                        if (!updProjId) return showAlert("Indica el UUID del proyecto", "error");
                        executeCall(`/api/v1/projects/${updProjId}/updates`, "GET", null, true);
                      }}
                      className="btn-action flex-1"
                    >
                      Ver Avances
                    </button>
                    <button
                      onClick={() => {
                        if (!updProjId) return showAlert("Indica el UUID del proyecto", "error");
                        const links = updLinks ? updLinks.split(",").map((s) => s.trim()) : [];
                        executeCall(`/api/v1/projects/${updProjId}/updates`, "POST", {
                          title: updTitle,
                          content: updContent,
                          attachments_url: links,
                        }, true);
                      }}
                      className="btn-primary flex-1"
                    >
                      Enviar Avance
                    </button>
                  </div>
                </div>

                {/* Evaluar Avance */}
                <div className="space-y-3 bg-slate-950/60 p-4 rounded-lg border border-slate-800/80">
                  <h3 className="text-xs font-bold uppercase tracking-wider text-teal-400">Evaluar Avance (Solo Directores)</h3>
                  <input
                    type="text"
                    placeholder="UUID del Avance"
                    value={revUpdId}
                    onChange={(e) => setRevUpdId(e.target.value)}
                    className="input-field"
                  />
                  <select value={revUpdStatus} onChange={(e) => setRevUpdStatus(e.target.value)} className="input-field">
                    <option value="APPROVED">APPROVED (Aprobado)</option>
                    <option value="CHANGES_REQUESTED">CHANGES_REQUESTED (Solicitar Corrección)</option>
                  </select>
                  <textarea
                    placeholder="Feedback formativo (Obligatorio en correcciones)"
                    rows={2}
                    value={revUpdFeedback}
                    onChange={(e) => setRevUpdFeedback(e.target.value)}
                    className="input-field"
                  />
                  <button
                    onClick={() => {
                      if (!revUpdId) return showAlert("Indica el UUID del avance", "error");
                      executeCall(`/api/v1/projects/updates/${revUpdId}/review`, "POST", {
                        status: revUpdStatus,
                        director_feedback: revUpdFeedback,
                      }, true);
                    }}
                    className="btn-primary w-full"
                  >
                    Emitir Evaluación
                  </button>
                </div>
              </div>
            </div>
          )}

          {/* SECCIÓN 8: RECURSOS (3 ESTADOS) */}
          {activeTab === "resources" && (
            <div className="bg-slate-900 p-5 rounded-xl border border-slate-800 space-y-4">
              <div className="flex items-center justify-between">
                <h2 className="text-base font-bold text-white flex items-center gap-2">
                  <BadgeDollarSign className="w-4 h-4 text-teal-400" />
                  <span>Solicitudes de Recursos & Compras (3 Estados)</span>
                </h2>
                <button
                  onClick={() => executeCall("/api/v1/directors/resources/pending", "GET", null, true)}
                  className="btn-action"
                >
                  Bandeja Global Pendientes
                </button>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {/* Solicitar Recurso */}
                <div className="space-y-3 bg-slate-950/60 p-4 rounded-lg border border-slate-800/80">
                  <h3 className="text-xs font-bold uppercase tracking-wider text-teal-400">Solicitar Recurso para Proyecto</h3>
                  <input
                    type="text"
                    placeholder="UUID del Proyecto"
                    value={resProjId}
                    onChange={(e) => setResProjId(e.target.value)}
                    className="input-field"
                  />
                  <div className="grid grid-cols-2 gap-2">
                    <select value={resType} onChange={(e) => setResType(e.target.value)} className="input-field">
                      <option value="MATERIAL">MATERIAL</option>
                      <option value="HARDWARE">HARDWARE</option>
                      <option value="SOFTWARE">SOFTWARE</option>
                      <option value="OTRO">OTRO</option>
                    </select>
                    <input
                      type="number"
                      placeholder="Costo Estimado COP"
                      value={resCost}
                      onChange={(e) => setResCost(e.target.value)}
                      className="input-field"
                    />
                  </div>
                  <input
                    type="text"
                    placeholder="Título del recurso"
                    value={resTitle}
                    onChange={(e) => setResTitle(e.target.value)}
                    className="input-field"
                  />
                  <textarea
                    placeholder="Justificación y especificaciones"
                    rows={2}
                    value={resDesc}
                    onChange={(e) => setResDesc(e.target.value)}
                    className="input-field"
                  />
                  <div className="flex gap-2">
                    <button
                      onClick={() => {
                        if (!resProjId) return showAlert("Indica el UUID del proyecto", "error");
                        executeCall(`/api/v1/projects/${resProjId}/resources`, "GET", null, true);
                      }}
                      className="btn-action flex-1"
                    >
                      Ver Recursos
                    </button>
                    <button
                      onClick={() => {
                        if (!resProjId) return showAlert("Indica el UUID del proyecto", "error");
                        executeCall(`/api/v1/projects/${resProjId}/resources`, "POST", {
                          resource_type: resType,
                          title: resTitle,
                          description: resDesc,
                          estimated_cost: parseFloat(resCost) || 0,
                        }, true);
                      }}
                      className="btn-primary flex-1"
                    >
                      Radicar Recurso
                    </button>
                  </div>
                </div>

                {/* Evaluar Recurso */}
                <div className="space-y-3 bg-slate-950/60 p-4 rounded-lg border border-slate-800/80">
                  <h3 className="text-xs font-bold uppercase tracking-wider text-teal-400">Evaluar Recurso (Solo Directores)</h3>
                  <input
                    type="text"
                    placeholder="UUID del Recurso"
                    value={evalResId}
                    onChange={(e) => setEvalResId(e.target.value)}
                    className="input-field"
                  />
                  <select value={evalResStatus} onChange={(e) => setEvalResStatus(e.target.value)} className="input-field">
                    <option value="APPROVED">APPROVED (Aprobado para compra)</option>
                    <option value="REJECTED">REJECTED (Rechazado)</option>
                    <option value="PENDING_QUOTE">PENDING_QUOTE (Pendiente Cotización)</option>
                  </select>
                  <textarea
                    placeholder="Feedback o resolución justificada"
                    rows={2}
                    value={evalResFeedback}
                    onChange={(e) => setEvalResFeedback(e.target.value)}
                    className="input-field"
                  />
                  <button
                    onClick={() => {
                      if (!evalResId) return showAlert("Indica el UUID del recurso", "error");
                      executeCall(`/api/v1/projects/resources/${evalResId}/review`, "POST", {
                        status: evalResStatus,
                        director_feedback: evalResFeedback,
                      }, true);
                    }}
                    className="btn-primary w-full"
                  >
                    Guardar Resolución
                  </button>
                </div>
              </div>
            </div>
          )}

          {/* SECCIÓN 9: INVENTARIO */}
          {activeTab === "inventory" && (
            <div className="bg-slate-900 p-5 rounded-xl border border-slate-800 space-y-4">
              <h2 className="text-base font-bold text-white flex items-center gap-2">
                <Package className="w-4 h-4 text-teal-400" />
                <span>Inventario y Equipamiento</span>
              </h2>

              <div className="flex flex-wrap gap-2">
                <button onClick={() => executeCall("/api/v1/inventory", "GET", null, true)} className="btn-action">GET /inventory (Catálogo)</button>
                <button onClick={() => executeCall("/api/v1/inventory?available_only=true", "GET", null, true)} className="btn-action">Solo Disponibles</button>
              </div>

              <div className="space-y-3 bg-slate-950/60 p-4 rounded-lg border border-slate-800/80 max-w-lg">
                <h3 className="text-xs font-bold uppercase tracking-wider text-teal-400">Crear Ítem en Inventario (Admin / Director)</h3>
                <div className="grid grid-cols-2 gap-2">
                  <input
                    type="text"
                    placeholder="Código (ej. LAB-OSC-01)"
                    value={invCode}
                    onChange={(e) => setInvCode(e.target.value)}
                    className="input-field"
                  />
                  <select value={invCat} onChange={(e) => setInvCat(e.target.value)} className="input-field">
                    <option value="EQUIPAMIENTO">EQUIPAMIENTO</option>
                    <option value="HERRAMIENTA">HERRAMIENTA</option>
                    <option value="CONSUMIBLE">CONSUMIBLE</option>
                    <option value="MODULO_ELECTRONICO">MODULO_ELECTRONICO</option>
                  </select>
                </div>
                <input
                  type="text"
                  placeholder="Nombre del ítem (ej. Osciloscopio Digital)"
                  value={invName}
                  onChange={(e) => setInvName(e.target.value)}
                  className="input-field"
                />
                <div className="grid grid-cols-2 gap-2">
                  <input
                    type="number"
                    placeholder="Stock Total"
                    value={invStock}
                    onChange={(e) => setInvStock(e.target.value)}
                    className="input-field"
                  />
                  <input
                    type="text"
                    placeholder="Ubicación (ej. Armario B)"
                    value={invLoc}
                    onChange={(e) => setInvLoc(e.target.value)}
                    className="input-field"
                  />
                </div>
                <button
                  onClick={() => executeCall("/api/v1/inventory", "POST", {
                    code: invCode,
                    name: invName,
                    category: invCat,
                    total_stock: parseInt(invStock) || 1,
                    location: invLoc,
                  }, true)}
                  className="btn-primary w-full"
                >
                  Registrar en Inventario
                </button>
              </div>
            </div>
          )}

          {/* SECCIÓN 10: PRÉSTAMOS */}
          {activeTab === "loans" && (
            <div className="bg-slate-900 p-5 rounded-xl border border-slate-800 space-y-4">
              <h2 className="text-base font-bold text-white flex items-center gap-2">
                <ArrowLeftRight className="w-4 h-4 text-teal-400" />
                <span>Gestión de Préstamos de Equipos</span>
              </h2>

              <div className="flex flex-wrap gap-2">
                <button onClick={() => executeCall("/api/v1/loans/my-loans", "GET", null, true)} className="btn-action">Mis Préstamos (/my-loans)</button>
                <button onClick={() => executeCall("/api/v1/loans", "GET", null, true)} className="btn-action">Panel General (/loans)</button>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {/* Solicitar Préstamo */}
                <div className="space-y-3 bg-slate-950/60 p-4 rounded-lg border border-slate-800/80">
                  <h3 className="text-xs font-bold uppercase tracking-wider text-teal-400">Solicitar Préstamo de Equipo</h3>
                  <input
                    type="text"
                    placeholder="UUID del Ítem de Inventario"
                    value={loanItemId}
                    onChange={(e) => setLoanItemId(e.target.value)}
                    className="input-field"
                  />
                  <textarea
                    placeholder="Motivo del préstamo"
                    rows={2}
                    value={loanReason}
                    onChange={(e) => setLoanReason(e.target.value)}
                    className="input-field"
                  />
                  <div className="grid grid-cols-2 gap-2">
                    <div>
                      <label className="text-[10px] text-slate-400 block mb-1">Fecha Inicio</label>
                      <input
                        type="date"
                        value={loanStart}
                        onChange={(e) => setLoanStart(e.target.value)}
                        className="input-field"
                      />
                    </div>
                    <div>
                      <label className="text-[10px] text-slate-400 block mb-1">Fecha Devolución</label>
                      <input
                        type="date"
                        value={loanEnd}
                        onChange={(e) => setLoanEnd(e.target.value)}
                        className="input-field"
                      />
                    </div>
                  </div>
                  <button
                    onClick={() => {
                      if (!loanItemId) return showAlert("Indica el ID del ítem", "error");
                      executeCall("/api/v1/loans/requests", "POST", {
                        item_id: loanItemId,
                        reason: loanReason,
                        requested_start_date: loanStart ? new Date(loanStart).toISOString() : new Date().toISOString(),
                        requested_end_date: loanEnd ? new Date(loanEnd).toISOString() : new Date().toISOString(),
                      }, true);
                    }}
                    className="btn-primary w-full"
                  >
                    Enviar Solicitud
                  </button>
                </div>

                {/* Evaluar / Devolver */}
                <div className="space-y-3 bg-slate-950/60 p-4 rounded-lg border border-slate-800/80">
                  <h3 className="text-xs font-bold uppercase tracking-wider text-teal-400">Evaluar o Devolver Préstamo</h3>
                  <input
                    type="text"
                    placeholder="UUID del Préstamo"
                    value={revLoanId}
                    onChange={(e) => setRevLoanId(e.target.value)}
                    className="input-field"
                  />
                  <select value={revLoanStatus} onChange={(e) => setRevLoanStatus(e.target.value)} className="input-field">
                    <option value="APPROVED">APPROVED (Aprobar con fecha solicitada)</option>
                    <option value="APPROVED_MODIFIED">APPROVED_MODIFIED (Aprobar con fecha ajustada)</option>
                    <option value="REJECTED">REJECTED (Rechazar)</option>
                  </select>
                  {revLoanStatus === "APPROVED_MODIFIED" && (
                    <div>
                      <label className="text-[10px] text-slate-400 block mb-1">Nueva Fecha Límite Aprobada</label>
                      <input
                        type="date"
                        value={revLoanModEnd}
                        onChange={(e) => setRevLoanModEnd(e.target.value)}
                        className="input-field"
                      />
                    </div>
                  )}
                  <input
                    type="text"
                    placeholder="Feedback o justificación"
                    value={revLoanFeedback}
                    onChange={(e) => setRevLoanFeedback(e.target.value)}
                    className="input-field"
                  />
                  <div className="flex gap-2">
                    <button
                      onClick={() => {
                        if (!revLoanId) return showAlert("Indica el UUID del préstamo", "error");
                        executeCall(`/api/v1/loans/${revLoanId}/review`, "POST", {
                          status: revLoanStatus,
                          reviewer_feedback: revLoanFeedback || null,
                          approved_end_date:
                            revLoanStatus === "APPROVED_MODIFIED" && revLoanModEnd
                              ? new Date(revLoanModEnd).toISOString()
                              : null,
                        }, true);
                      }}
                      className="btn-primary flex-1"
                    >
                      Evaluar
                    </button>
                    <button
                      onClick={() => {
                        if (!revLoanId) return showAlert("Indica el UUID del préstamo", "error");
                        executeCall(`/api/v1/loans/${revLoanId}/return`, "POST", null, true);
                      }}
                      className="btn-action flex-1 bg-amber-600/20 text-amber-300 border-amber-500/30 hover:bg-amber-600/30"
                    >
                      Registrar Devolución
                    </button>
                  </div>
                </div>
              </div>
            </div>
          )}

          {/* SECCIÓN 11: TALLER IMPRESIÓN 3D */}
          {activeTab === "print3d" && (
            <div className="bg-slate-900 p-5 rounded-xl border border-slate-800 space-y-4">
              <h2 className="text-base font-bold text-white flex items-center gap-2">
                <Printer className="w-4 h-4 text-teal-400" />
                <span>Taller de Impresión 3D</span>
              </h2>

              <div className="flex flex-wrap gap-2">
                <button onClick={() => executeCall("/api/v1/print3d/my-requests", "GET", null, true)} className="btn-action">Mis Piezas (/my-requests)</button>
                <button onClick={() => executeCall("/api/v1/print3d/queue", "GET", null, true)} className="btn-action">Cola del Laboratorio (/queue)</button>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {/* Subir Pieza */}
                <div className="space-y-3 bg-slate-950/60 p-4 rounded-lg border border-slate-800/80">
                  <h3 className="text-xs font-bold uppercase tracking-wider text-teal-400">Subir Archivo 3D (.stl, .obj, .3mf, .step)</h3>
                  <input
                    type="file"
                    accept=".stl,.obj,.3mf,.step"
                    onChange={(e) => setFile3d(e.target.files?.[0] || null)}
                    className="input-field text-slate-400 file:mr-2 file:py-1 file:px-2 file:rounded file:border-0 file:text-xs file:bg-teal-600 file:text-white"
                  />
                  <div className="grid grid-cols-3 gap-2">
                    <select value={mat3d} onChange={(e) => setMat3d(e.target.value)} className="input-field">
                      <option value="PLA">PLA</option>
                      <option value="PETG">PETG</option>
                      <option value="ABS">ABS</option>
                      <option value="RESINA">RESINA</option>
                    </select>
                    <input
                      type="text"
                      placeholder="Color"
                      value={color3d}
                      onChange={(e) => setColor3d(e.target.value)}
                      className="input-field"
                    />
                    <input
                      type="number"
                      placeholder="% Relleno"
                      value={infill3d}
                      onChange={(e) => setInfill3d(e.target.value)}
                      className="input-field"
                    />
                  </div>
                  <input
                    type="text"
                    placeholder="Notas o parámetros de impresión"
                    value={notes3d}
                    onChange={(e) => setNotes3d(e.target.value)}
                    className="input-field"
                  />
                  <button
                    onClick={() => {
                      if (!file3d) return showAlert("Selecciona un archivo 3D", "error");
                      const fd = new FormData();
                      fd.append("file", file3d);
                      fd.append("material", mat3d);
                      fd.append("color", color3d);
                      fd.append("infill_percentage", infill3d);
                      fd.append("notes", notes3d);
                      executeCall("/api/v1/print3d/requests", "POST", fd, true, true);
                    }}
                    className="btn-primary w-full"
                  >
                    Enviar a Cola de Impresión
                  </button>
                </div>

                {/* Ciclo de Operación 3D */}
                <div className="space-y-4">
                  {/* Evaluación */}
                  <div className="space-y-3 bg-slate-950/60 p-4 rounded-lg border border-slate-800/80">
                    <h3 className="text-xs font-bold uppercase tracking-wider text-teal-400">Evaluar Solicitud (Admin / Director)</h3>
                    <div className="flex gap-2">
                      <input
                        type="text"
                        placeholder="UUID Solicitud 3D"
                        value={rev3dId}
                        onChange={(e) => setRev3dId(e.target.value)}
                        className="input-field"
                      />
                      <select value={rev3dStatus} onChange={(e) => setRev3dStatus(e.target.value)} className="input-field">
                        <option value="APPROVED">APPROVED</option>
                        <option value="REJECTED">REJECTED</option>
                      </select>
                      <button
                        onClick={() => {
                          if (!rev3dId) return showAlert("Indica el ID de solicitud", "error");
                          executeCall(`/api/v1/print3d/requests/${rev3dId}/review`, "POST", { status: rev3dStatus }, true);
                        }}
                        className="btn-primary"
                      >
                        Evaluar
                      </button>
                    </div>
                  </div>

                  {/* Operación Técnica */}
                  <div className="space-y-3 bg-slate-950/60 p-4 rounded-lg border border-slate-800/80">
                    <h3 className="text-xs font-bold uppercase tracking-wider text-teal-400">Control Operativo (Operador 3D)</h3>
                    <div className="flex gap-2">
                      <input
                        type="text"
                        placeholder="UUID Solicitud 3D"
                        value={op3dId}
                        onChange={(e) => setOp3dId(e.target.value)}
                        className="input-field"
                      />
                      <select value={op3dStatus} onChange={(e) => setOp3dStatus(e.target.value)} className="input-field">
                        <option value="IN_PROGRESS">IN_PROGRESS (En Máquina)</option>
                        <option value="COMPLETED">COMPLETED (Impresión Finalizada)</option>
                        <option value="DELIVERED">DELIVERED (Entregada a Alumno)</option>
                      </select>
                      <button
                        onClick={() => {
                          if (!op3dId) return showAlert("Indica el ID de solicitud", "error");
                          executeCall(`/api/v1/print3d/requests/${op3dId}/status`, "POST", { status: op3dStatus }, true);
                        }}
                        className="btn-primary"
                      >
                        Actualizar
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          )}

          {/* Consola de Inspección de Respuestas en Vivo */}
          <div>
            <div className="flex items-center justify-between mb-2">
              <h3 className="text-xs font-bold uppercase tracking-wider text-slate-400">
                Consola de Inspección HTTP en Vivo
              </h3>
            </div>
            <Inspector
              response={lastResponse}
              onClear={() => setLastResponse(null)}
            />
          </div>
        </main>
      </div>
    </div>
  );
}

export default function Home() {
  return (
    <Suspense fallback={
      <div className="min-h-screen flex items-center justify-center bg-slate-950 text-slate-200">
        <div className="w-8 h-8 border-4 border-teal-500/20 border-t-teal-400 rounded-full animate-spin" />
      </div>
    }>
      <DashboardContent />
    </Suspense>
  );
}
