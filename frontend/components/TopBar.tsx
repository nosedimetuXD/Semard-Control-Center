"use client";

import { useState } from "react";
import { getApiBaseUrl, getAuthToken, removeAuthToken } from "@/lib/api";
import {
  LogIn,
  LogOut,
  Radio,
  Search,
  Plus,
  ShieldCheck,
  UserCheck,
} from "lucide-react";

interface TopBarProps {
  user: any;
  onOpenNewAction?: () => void;
  onShowAlert: (msg: string, type: "info" | "success" | "error") => void;
  onLogout: () => void;
}

export default function TopBar({
  user,
  onOpenNewAction,
  onShowAlert,
  onLogout,
}: TopBarProps) {
  const [searchQuery, setSearchQuery] = useState("");

  function handleGoogleLogin() {
    const base = getApiBaseUrl();
    const returnTo = `${window.location.origin}/auth/callback`;
    window.location.href = `${base}/api/v1/auth/google/login?return_to=${encodeURIComponent(
      returnTo
    )}&redirect=true`;
  }

  return (
    <header className="h-16 bg-slate-900/90 backdrop-blur border-b border-slate-800 px-6 flex items-center justify-between gap-4 sticky top-0 z-40">
      {/* Search Input / Quick Filter */}
      <div className="flex-1 max-w-md hidden md:flex items-center gap-2 bg-slate-950/70 px-3 py-1.5 rounded-lg border border-slate-800 text-xs focus-within:border-teal-500/60 transition">
        <Search className="w-3.5 h-3.5 text-slate-400 shrink-0" />
        <input
          type="text"
          placeholder="Buscar proyectos, equipos, piezas 3D o eventos..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="bg-transparent text-slate-200 outline-none w-full placeholder:text-slate-500 text-xs"
        />
      </div>

      {/* Right Controls */}
      <div className="flex items-center gap-3 ml-auto">
        {/* Connection Status Pill */}
        <div className="hidden sm:flex items-center gap-2 px-2.5 py-1 rounded-full bg-slate-950/80 border border-slate-800 text-[11px] text-slate-300">
          <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
          <span className="font-mono text-slate-400">Coolify API</span>
          <span className="text-emerald-400 font-semibold">En Línea</span>
        </div>

        {/* Quick Action Button */}
        {onOpenNewAction && (
          <button
            onClick={onOpenNewAction}
            className="text-xs font-semibold px-3 py-1.5 rounded-lg bg-teal-600 hover:bg-teal-500 text-white flex items-center gap-1.5 shadow-sm shadow-teal-500/20 transition"
          >
            <Plus className="w-3.5 h-3.5" />
            <span>Acción Rápida</span>
          </button>
        )}

        {/* User Session Profile or Google Login */}
        {user ? (
          <div className="flex items-center gap-2.5 pl-2 border-l border-slate-800">
            {user.avatar_url ? (
              <img
                src={user.avatar_url}
                alt={user.full_name}
                className="w-8 h-8 rounded-full border border-teal-500/40 object-cover"
              />
            ) : (
              <div className="w-8 h-8 rounded-full bg-teal-600/30 text-teal-300 flex items-center justify-center font-bold text-xs border border-teal-500/40">
                {user.full_name?.charAt(0) || "U"}
              </div>
            )}
            <div className="hidden md:block text-left">
              <p className="text-xs font-semibold text-white leading-none">
                {user.full_name || "Investigador"}
              </p>
              <div className="flex items-center gap-1.5 mt-0.5">
                <span className="text-[10px] text-teal-400 font-medium">
                  {user.role}
                </span>
                {user.can_operate_3d && (
                  <span className="text-[9px] px-1 rounded bg-teal-500/20 text-teal-300 font-bold">
                    3D
                  </span>
                )}
              </div>
            </div>
            <button
              onClick={onLogout}
              className="p-1.5 text-slate-400 hover:text-red-400 rounded-md hover:bg-slate-800 transition"
              title="Cerrar sesión"
            >
              <LogOut className="w-4 h-4" />
            </button>
          </div>
        ) : (
          <button
            onClick={handleGoogleLogin}
            className="text-xs font-semibold px-3.5 py-1.5 rounded-lg bg-white text-slate-900 hover:bg-slate-100 flex items-center gap-2 shadow-sm transition"
          >
            <svg className="w-3.5 h-3.5" viewBox="0 0 24 24">
              <path
                fill="#4285F4"
                d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
              />
              <path
                fill="#34A853"
                d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
              />
              <path
                fill="#FBBC05"
                d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.06H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.94l2.85-2.22.81-.63z"
              />
              <path
                fill="#EA4335"
                d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.06l3.66 2.84c.87-2.6 3.3-4.52 6.16-4.52z"
              />
            </svg>
            <span>Acceso Institucional</span>
          </button>
        )}
      </div>
    </header>
  );
}
