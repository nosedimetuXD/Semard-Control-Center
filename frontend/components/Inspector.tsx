"use client";

import { useState } from "react";
import { ApiResponse } from "@/lib/api";
import { Copy, Check, Clock, Trash2 } from "lucide-react";

interface InspectorProps {
  response: ApiResponse | null;
  onClear?: () => void;
}

export default function Inspector({ response, onClear }: InspectorProps) {
  const [copied, setCopied] = useState(false);

  if (!response) {
    return (
      <div className="bg-slate-900 border border-slate-800 rounded-xl p-5 flex flex-col items-center justify-center text-center text-slate-500 h-64">
        <p className="text-xs">No se han ejecutado peticiones en esta sesión.</p>
        <p className="text-[11px] text-slate-600 mt-1">Haz clic en cualquier acción de los módulos para inspeccionar la respuesta en vivo.</p>
      </div>
    );
  }

  const isSuccess = response.ok;
  const jsonText =
    typeof response.data === "object"
      ? JSON.stringify(response.data, null, 2)
      : response.rawText;

  function handleCopy() {
    navigator.clipboard.writeText(jsonText);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }

  return (
    <div className="bg-slate-900 border border-slate-800 rounded-xl overflow-hidden flex flex-col">
      {/* Barra de cabecera de la respuesta */}
      <div className="bg-slate-950/80 px-4 py-3 border-b border-slate-800 flex flex-wrap items-center justify-between gap-2">
        <div className="flex items-center gap-2">
          <span className={`text-[11px] font-bold px-2 py-0.5 rounded font-mono ${
            response.method === "GET" ? "bg-blue-500/20 text-blue-300" :
            response.method === "POST" ? "bg-emerald-500/20 text-emerald-300" :
            response.method === "PUT" ? "bg-amber-500/20 text-amber-300" :
            response.method === "PATCH" ? "bg-purple-500/20 text-purple-300" :
            "bg-red-500/20 text-red-300"
          }`}>
            {response.method}
          </span>
          <span className="text-xs font-mono text-slate-300">{response.endpoint}</span>
        </div>

        <div className="flex items-center gap-2 text-xs">
          <span className="flex items-center gap-1 text-slate-400 font-mono text-[11px]">
            <Clock className="w-3 h-3 text-teal-400" />
            {response.durationMs}ms
          </span>

          <span className={`px-2 py-0.5 rounded font-mono font-bold text-[11px] ${
            isSuccess ? "bg-emerald-500/20 text-emerald-300" : "bg-red-500/20 text-red-300"
          }`}>
            {response.status} {response.statusText}
          </span>

          <button
            onClick={handleCopy}
            className="p-1.5 rounded hover:bg-slate-800 text-slate-400 hover:text-slate-200 transition"
            title="Copiar JSON"
          >
            {copied ? <Check className="w-3.5 h-3.5 text-emerald-400" /> : <Copy className="w-3.5 h-3.5" />}
          </button>

          {onClear && (
            <button
              onClick={onClear}
              className="p-1.5 rounded hover:bg-slate-800 text-slate-400 hover:text-red-400 transition"
              title="Limpiar consola"
            >
              <Trash2 className="w-3.5 h-3.5" />
            </button>
          )}
        </div>
      </div>

      {/* Visor de JSON / Texto */}
      <div className="p-4 bg-slate-950 overflow-auto max-h-96">
        <pre className="text-xs font-mono text-slate-300 whitespace-pre-wrap leading-relaxed">
          {jsonText}
        </pre>
      </div>
    </div>
  );
}
