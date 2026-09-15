"use client";

import { useState, useEffect } from "react";
import { getApiBaseUrl, setApiBaseUrl, getAuthToken, setAuthToken, removeAuthToken, apiCall } from "@/lib/api";
import { LogIn, Key, LogOut, CheckCircle2, RefreshCw } from "lucide-react";

interface HeaderProps {
  onSessionChange?: () => void;
  onShowAlert?: (msg: string, type: "info" | "success" | "error") => void;
}

export default function Header({ onSessionChange, onShowAlert }: HeaderProps) {
  const [apiUrl, setApiUrl] = useState("");
  const [user, setUser] = useState<any>(null);
  const [token, setToken] = useState<string | null>(null);

  useEffect(() => {
    setApiUrl(getApiBaseUrl());
    const t = getAuthToken();
    setToken(t);
    if (t) {
      loadProfile(t);
    }
  }, []);

  async function loadProfile(t: string) {
    const res = await apiCall("/api/v1/auth/me", { requireAuth: true });
    if (res.ok && res.data) {
      setUser(res.data);
    } else {
      setUser(null);
    }
  }

  function handleSaveApiUrl(val: string) {
    setApiUrl(val);
    setApiBaseUrl(val);
  }

  function handleGoogleLogin() {
    const base = getApiBaseUrl();
    const returnTo = `${window.location.origin}/auth/callback`;
    window.location.href = `${base}/api/v1/auth/google/login?return_to=${encodeURIComponent(returnTo)}&redirect=true`;
  }

  function handlePromptToken() {
    const current = getAuthToken() || "";
    const input = prompt("Pega tu credencial JWT de sesión:", current);
    if (input !== null) {
      const clean = input.trim();
      if (clean) {
        setAuthToken(clean);
        setToken(clean);
        loadProfile(clean);
        onShowAlert?.("Token JWT guardado con éxito", "success");
      } else {
        removeAuthToken();
        setToken(null);
        setUser(null);
        onShowAlert?.("Token eliminado", "info");
      }
      onSessionChange?.();
    }
  }

  function handleLogout() {
    removeAuthToken();
    setToken(null);
    setUser(null);
    onShowAlert?.("Sesión cerrada", "info");
    onSessionChange?.();
  }

  return (
    <header className="border-b border-slate-800 bg-slate-900/90 backdrop-blur sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 py-3 flex flex-wrap items-center justify-between gap-4">
        {/* Marca SEMARD */}
        <div className="flex items-center gap-3">
          <div className="w-9 h-9 rounded-lg bg-teal-500/20 text-teal-400 flex items-center justify-center font-bold text-lg border border-teal-500/30">
            S
          </div>
          <div>
            <div className="flex items-center gap-2">
              <h1 className="text-sm font-bold tracking-tight text-white">
                SEMARD Control Center
              </h1>
              <span className="text-[10px] font-semibold px-2 py-0.5 rounded-full bg-teal-500/10 text-teal-400 border border-teal-500/30">
                Next.js / Vercel
              </span>
            </div>
            <p className="text-xs text-slate-400">Universidad de Cartagena — Facultad de Ingeniería</p>
          </div>
        </div>

        {/* Acciones de Sesión y Configuración */}
        <div className="flex items-center flex-wrap gap-2.5">
          {/* Selector de API URL */}
          <div className="flex items-center gap-1.5 bg-slate-950/80 px-2.5 py-1.5 rounded-md border border-slate-800 text-xs">
            <span className="text-slate-400 text-[11px]">API:</span>
            <input
              type="text"
              value={apiUrl}
              onChange={(e) => handleSaveApiUrl(e.target.value)}
              placeholder="https://..."
              className="bg-transparent text-slate-200 outline-none w-52 font-mono text-[11px] focus:ring-1 focus:ring-teal-500 rounded px-1"
              title="URL del servidor backend Go (por defecto Coolify VM)"
            />
            <button
              onClick={() => handleSaveApiUrl("https://semardcontrolcenter.147.5.103.87.sslip.io")}
              className="text-slate-500 hover:text-slate-300 ml-1"
              title="Restaurar backend de producción Coolify"
            >
              <RefreshCw className="w-3 h-3" />
            </button>
          </div>

          {/* Botón Google Login */}
          <button
            onClick={handleGoogleLogin}
            className="text-xs font-semibold px-3 py-1.5 rounded-md bg-white text-slate-900 hover:bg-slate-200 flex items-center gap-1.5 shadow-sm transition"
          >
            <svg className="w-3.5 h-3.5" viewBox="0 0 24 24">
              <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
              <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
              <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.06H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.94l2.85-2.22.81-.63z"/>
              <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.06l3.66 2.84c.87-2.6 3.3-4.52 6.16-4.52z"/>
            </svg>
            <span>Google Login</span>
          </button>

          {/* Pegar JWT */}
          <button
            onClick={handlePromptToken}
            className="text-xs px-2.5 py-1.5 rounded-md bg-slate-800 hover:bg-slate-700 text-slate-300 border border-slate-700 transition flex items-center gap-1.5"
          >
            <Key className="w-3.5 h-3.5" />
            <span>JWT</span>
          </button>

          {/* Usuario Autenticado */}
          {user && (
            <div className="flex items-center gap-2 bg-slate-800/90 px-2.5 py-1 rounded-md border border-slate-700 text-xs">
              <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
              <span className="font-medium text-slate-200">{user.full_name || user.email}</span>
              <span className="px-1.5 py-0.5 rounded text-[10px] font-bold bg-teal-500/20 text-teal-300">
                {user.role} {user.can_operate_3d ? "• 3D" : ""}
              </span>
              <button
                onClick={handleLogout}
                className="text-slate-400 hover:text-red-400 ml-1"
                title="Cerrar sesión"
              >
                <LogOut className="w-3.5 h-3.5" />
              </button>
            </div>
          )}
        </div>
      </div>
    </header>
  );
}
