"use client";

import { useState } from "react";
import { Terminal, X, ChevronUp, ChevronDown, Radio, Copy, Check } from "lucide-react";
import Inspector from "./Inspector";
import { getApiBaseUrl, setApiBaseUrl, getAuthToken, ApiResponse } from "@/lib/api";

interface DebugDrawerProps {
  lastResponse: ApiResponse | null;
  onClearResponse: () => void;
}

export default function DebugDrawer({
  lastResponse,
  onClearResponse,
}: DebugDrawerProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [copiedToken, setCopiedToken] = useState(false);

  const token = getAuthToken();

  function handleCopyToken() {
    if (token) {
      navigator.clipboard.writeText(token);
      setCopiedToken(true);
      setTimeout(() => setCopiedToken(false), 2000);
    }
  }

  return (
    <div className="fixed bottom-4 right-4 z-50 flex flex-col items-end">
      {/* Ventana de Consola Desplegable */}
      {isOpen && (
        <div className="w-[92vw] sm:w-[540px] max-h-[70vh] bg-slate-900 border border-slate-800 rounded-2xl shadow-2xl flex flex-col overflow-hidden mb-3 animate-in fade-in slide-in-from-bottom-5">
          {/* Header de la Consola */}
          <div className="p-3 bg-slate-950 border-b border-slate-800 flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Terminal className="w-4 h-4 text-teal-400" />
              <span className="text-xs font-bold text-white">Inspector de API & Red</span>
              <span className="text-[10px] px-1.5 py-0.2 rounded bg-teal-500/10 text-teal-400 border border-teal-500/20">
                Dev Mode
              </span>
            </div>
            <button
              onClick={() => setIsOpen(false)}
              className="p-1 rounded text-slate-400 hover:text-white"
            >
              <X className="w-4 h-4" />
            </button>
          </div>

          {/* Información de Conexión & Token */}
          <div className="p-3 bg-slate-950/60 border-b border-slate-800 text-[11px] space-y-1.5">
            <div className="flex items-center justify-between text-slate-400">
              <span>Backend Host:</span>
              <span className="font-mono text-slate-200 text-[10px] truncate max-w-[280px]">
                {getApiBaseUrl()}
              </span>
            </div>
            <div className="flex items-center justify-between text-slate-400">
              <span>Sesión JWT:</span>
              {token ? (
                <div className="flex items-center gap-1.5">
                  <span className="text-emerald-400 font-semibold text-[10px]">Activo</span>
                  <button
                    onClick={handleCopyToken}
                    className="text-slate-400 hover:text-white"
                    title="Copiar token"
                  >
                    {copiedToken ? <Check className="w-3 h-3 text-emerald-400" /> : <Copy className="w-3 h-3" />}
                  </button>
                </div>
              ) : (
                <span className="text-slate-500 text-[10px]">Sin token</span>
              )}
            </div>
          </div>

          {/* Visor de Respuestas */}
          <div className="flex-1 overflow-auto p-3">
            <Inspector response={lastResponse} onClear={onClearResponse} />
          </div>
        </div>
      )}

      {/* Botón Flotante */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="px-3.5 py-2 rounded-full bg-slate-900 hover:bg-slate-800 border border-slate-700 text-slate-200 shadow-xl flex items-center gap-2 text-xs font-semibold transition hover:scale-105"
      >
        <span className="w-2 h-2 rounded-full bg-teal-400 animate-pulse" />
        <Terminal className="w-3.5 h-3.5 text-teal-400" />
        <span className="hidden sm:inline">Inspector API</span>
        {isOpen ? <ChevronDown className="w-3.5 h-3.5" /> : <ChevronUp className="w-3.5 h-3.5" />}
      </button>
    </div>
  );
}
