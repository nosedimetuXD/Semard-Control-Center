"use client";

import { useState, useEffect, Suspense } from "react";
import { useSearchParams } from "next/navigation";
import Sidebar, { ViewKey } from "@/components/Sidebar";
import TopBar from "@/components/TopBar";
import DebugDrawer from "@/components/DebugDrawer";
import DashboardView from "@/components/views/DashboardView";
import ProjectsView from "@/components/views/ProjectsView";
import Print3DView from "@/components/views/Print3DView";
import InventoryView from "@/components/views/InventoryView";
import EventsView from "@/components/views/EventsView";
import MembersView from "@/components/views/MembersView";
import MyPortalView from "@/components/views/MyPortalView";
import { apiCall, getAuthToken, removeAuthToken, ApiResponse } from "@/lib/api";
import { CheckCircle, AlertCircle, X, Printer, Package, FolderGit2 } from "lucide-react";

function AppContent() {
  const searchParams = useSearchParams();
  const [currentView, setCurrentView] = useState<ViewKey>("dashboard");
  const [user, setUser] = useState<any | null>(null);
  const [lastResponse, setLastResponse] = useState<ApiResponse | null>(null);
  const [alert, setAlert] = useState<{ msg: string; type: "info" | "success" | "error" } | null>(null);

  // Modal de Acción Rápida Global
  const [quickActionModal, setQuickActionModal] = useState(false);

  useEffect(() => {
    loadUserProfile();

    const authSuccess = searchParams.get("auth_success");
    const status = searchParams.get("status");

    if (authSuccess === "true") {
      showAlert("¡Sesión iniciada con éxito mediante Google UdeC!", "success");
      loadUserProfile();
      window.history.replaceState({}, document.title, window.location.pathname);
    } else if (status === "NEEDS_REGISTRATION") {
      setCurrentView("members");
      showAlert("Cuenta institucional verificada. Por favor radica tu solicitud de membresía en el formulario de la pestaña Miembros.", "info");
      window.history.replaceState({}, document.title, window.location.pathname);
    } else if (status === "PENDING_APPROVAL") {
      showAlert("Tu solicitud de ingreso ya está radicada y en proceso de revisión por los Directores.", "info");
      window.history.replaceState({}, document.title, window.location.pathname);
    } else if (status === "REJECTED") {
      const msg = searchParams.get("msg") || "Tu solicitud de ingreso previa fue rechazada.";
      showAlert(msg, "error");
      window.history.replaceState({}, document.title, window.location.pathname);
    }
  }, [searchParams]);

  async function loadUserProfile() {
    const token = getAuthToken();
    if (!token) {
      setUser(null);
      return;
    }
    const res = await apiCall("/api/v1/auth/me", { requireAuth: true });
    setLastResponse(res);
    if (res.ok && res.data) {
      setUser(res.data);
    } else {
      setUser(null);
    }
  }

  function showAlert(msg: string, type: "info" | "success" | "error") {
    setAlert({ msg, type });
  }

  function handleLogout() {
    removeAuthToken();
    setUser(null);
    showAlert("Sesión cerrada correctamente", "info");
  }

  return (
    <div className="min-h-screen flex flex-col lg:flex-row bg-slate-950 text-slate-100 antialiased">
      {/* Sidebar Lateral */}
      <Sidebar
        currentView={currentView}
        onSelectView={(v) => setCurrentView(v)}
        user={user}
      />

      {/* Contenido Principal */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Barra Superior */}
        <TopBar
          user={user}
          onOpenNewAction={() => setQuickActionModal(true)}
          onShowAlert={showAlert}
          onLogout={handleLogout}
        />

        {/* Notificación Toast */}
        {alert && (
          <div className="max-w-7xl mx-auto w-full px-6 pt-4">
            <div
              className={`p-3 rounded-xl text-xs border flex items-center justify-between shadow-lg ${
                alert.type === "error"
                  ? "bg-red-500/10 text-red-300 border-red-500/30"
                  : alert.type === "success"
                  ? "bg-emerald-500/10 text-emerald-300 border-emerald-500/30"
                  : "bg-slate-900 text-teal-300 border-teal-500/30"
              }`}
            >
              <div className="flex items-center gap-2">
                {alert.type === "success" ? (
                  <CheckCircle className="w-4 h-4 text-emerald-400 shrink-0" />
                ) : alert.type === "error" ? (
                  <AlertCircle className="w-4 h-4 text-red-400 shrink-0" />
                ) : (
                  <CheckCircle className="w-4 h-4 text-teal-400 shrink-0" />
                )}
                <span className="font-medium">{alert.msg}</span>
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

        {/* Vistas Dinámicas del Producto */}
        <main className="flex-1 p-6 max-w-7xl mx-auto w-full">
          {currentView === "dashboard" && (
            <DashboardView
              onNavigate={(v) => setCurrentView(v)}
              onOpenQuickAction={(actionType) => {
                if (actionType === "print3d") setCurrentView("print3d");
                if (actionType === "loan") setCurrentView("inventory");
              }}
            />
          )}

          {currentView === "projects" && (
            <ProjectsView
              user={user}
              onShowAlert={showAlert}
              onTrackResponse={(res) => setLastResponse(res)}
            />
          )}

          {currentView === "print3d" && (
            <Print3DView
              user={user}
              onShowAlert={showAlert}
              onTrackResponse={(res) => setLastResponse(res)}
            />
          )}

          {currentView === "inventory" && (
            <InventoryView
              user={user}
              onShowAlert={showAlert}
              onTrackResponse={(res) => setLastResponse(res)}
            />
          )}

          {currentView === "events" && (
            <EventsView
              user={user}
              onShowAlert={showAlert}
              onTrackResponse={(res) => setLastResponse(res)}
            />
          )}

          {currentView === "members" && (
            <MembersView
              user={user}
              onShowAlert={showAlert}
              onTrackResponse={(res) => setLastResponse(res)}
            />
          )}

          {currentView === "my-portal" && (
            <MyPortalView
              user={user}
              onShowAlert={showAlert}
              onTrackResponse={(res) => setLastResponse(res)}
              onProfileUpdated={loadUserProfile}
            />
          )}
        </main>
      </div>

      {/* Cajón de Depuración Flotante (Discreto, no interfiere con el diseño) */}
      <DebugDrawer
        lastResponse={lastResponse}
        onClearResponse={() => setLastResponse(null)}
      />

      {/* Modal de Acción Rápida */}
      {quickActionModal && (
        <div className="fixed inset-0 z-50 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-2xl max-w-sm w-full p-6 space-y-4 shadow-2xl animate-in fade-in zoom-in-95">
            <div className="flex items-center justify-between">
              <h3 className="text-base font-bold text-white">Acceso Rápido SEMARD</h3>
              <button
                onClick={() => setQuickActionModal(false)}
                className="text-slate-400 hover:text-white"
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            <p className="text-xs text-slate-400">
              ¿Qué trámite deseas realizar en el laboratorio?
            </p>

            <div className="space-y-2 text-xs">
              <button
                onClick={() => {
                  setQuickActionModal(false);
                  setCurrentView("print3d");
                }}
                className="w-full p-3 rounded-xl bg-slate-950/60 border border-slate-800 hover:border-teal-500/50 text-left flex items-center gap-3 transition"
              >
                <div className="w-8 h-8 rounded-lg bg-cyan-500/10 text-cyan-400 flex items-center justify-center">
                  <Printer className="w-4 h-4" />
                </div>
                <div>
                  <div className="font-bold text-white">Encargar Impresión 3D</div>
                  <div className="text-[11px] text-slate-400">Subir modelo .stl, .step o .obj</div>
                </div>
              </button>

              <button
                onClick={() => {
                  setQuickActionModal(false);
                  setCurrentView("inventory");
                }}
                className="w-full p-3 rounded-xl bg-slate-950/60 border border-slate-800 hover:border-teal-500/50 text-left flex items-center gap-3 transition"
              >
                <div className="w-8 h-8 rounded-lg bg-indigo-500/10 text-indigo-400 flex items-center justify-center">
                  <Package className="w-4 h-4" />
                </div>
                <div>
                  <div className="font-bold text-white">Solicitar Préstamo de Equipo</div>
                  <div className="text-[11px] text-slate-400">Osciloscopios, sensores, kits</div>
                </div>
              </button>

              <button
                onClick={() => {
                  setQuickActionModal(false);
                  setCurrentView("projects");
                }}
                className="w-full p-3 rounded-xl bg-slate-950/60 border border-slate-800 hover:border-teal-500/50 text-left flex items-center gap-3 transition"
              >
                <div className="w-8 h-8 rounded-lg bg-teal-500/10 text-teal-400 flex items-center justify-center">
                  <FolderGit2 className="w-4 h-4" />
                </div>
                <div>
                  <div className="font-bold text-white">Registrar Avance de Proyecto</div>
                  <div className="text-[11px] text-slate-400">Subir bitácora y enlaces</div>
                </div>
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default function Home() {
  return (
    <Suspense
      fallback={
        <div className="min-h-screen flex items-center justify-center bg-slate-950 text-slate-200">
          <div className="w-8 h-8 border-4 border-teal-500/20 border-t-teal-400 rounded-full animate-spin" />
        </div>
      }
    >
      <AppContent />
    </Suspense>
  );
}
