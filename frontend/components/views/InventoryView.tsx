"use client";

import { useState, useEffect } from "react";
import {
  Package,
  Plus,
  Search,
  CheckCircle,
  AlertCircle,
  Clock,
  Calendar,
  Layers,
  Cpu,
  X,
  ArrowRight,
} from "lucide-react";
import { apiCall } from "@/lib/api";

interface InventoryViewProps {
  user: any;
  onShowAlert: (msg: string, type: "info" | "success" | "error") => void;
  onTrackResponse: (res: any) => void;
}

export default function InventoryView({
  user,
  onShowAlert,
  onTrackResponse,
}: InventoryViewProps) {
  const [items, setItems] = useState<any[]>([]);
  const [categoryFilter, setCategoryFilter] = useState("TODOS");
  const [search, setSearch] = useState("");
  const [selectedItem, setSelectedItem] = useState<any | null>(null);

  // Modales
  const [showLoanModal, setShowLoanModal] = useState(false);
  const [showCreateModal, setShowCreateModal] = useState(false);

  // Form states Préstamo
  const [loanReason, setLoanReason] = useState("");
  const [loanStartDate, setLoanStartDate] = useState("");
  const [loanEndDate, setLoanEndDate] = useState("");

  // Form states Crear Item
  const [newCode, setNewCode] = useState("");
  const [newName, setNewName] = useState("");
  const [newCategory, setNewCategory] = useState("EQUIPAMIENTO");
  const [newStock, setNewStock] = useState("1");
  const [newLocation, setNewLocation] = useState("");

  const defaultItems = [
    {
      id: "inv-1",
      code: "LAB-OSC-01",
      name: "Osciloscopio Digital Rigol DS1054Z 50MHz (4 Canales)",
      category: "EQUIPAMIENTO",
      total_stock: 3,
      available_stock: 2,
      location: "Mesa Principal Lab - Piedra de Bolívar",
    },
    {
      id: "inv-2",
      code: "LAB-PWR-02",
      name: "Fuente de Poder Regulable DC 30V 5A con Display Digital",
      category: "EQUIPAMIENTO",
      total_stock: 4,
      available_stock: 3,
      location: "Armario 1 - Nivel A",
    },
    {
      id: "inv-3",
      code: "LAB-SEN-LIDAR",
      name: "Sensor LiDAR 360° RPLIDAR A1M8 con Interfaz USB",
      category: "MODULO_ELECTRONICO",
      total_stock: 2,
      available_stock: 1,
      location: "Caja Robótica Móvil #3",
    },
    {
      id: "inv-4",
      code: "LAB-MCU-ESP32",
      name: "Kit de Desarrollo ESP32-WROOM-32 Wi-Fi & Bluetooth (Pack 5)",
      category: "MODULO_ELECTRONICO",
      total_stock: 10,
      available_stock: 8,
      location: "Gaveta IoT",
    },
    {
      id: "inv-5",
      code: "LAB-TOOL-SOLD",
      name: "Estación de Soldadura TS100 con Control Térmico y Soporte",
      category: "HERRAMIENTA",
      total_stock: 4,
      available_stock: 4,
      location: "Banco de Ensamble",
    },
    {
      id: "inv-6",
      code: "LAB-SEN-RGBD",
      name: "Cámara de Profundidad Intel RealSense D435i",
      category: "MODULO_ELECTRONICO",
      total_stock: 1,
      available_stock: 0,
      location: "En Proyecto Rover Autónomo",
    },
  ];

  useEffect(() => {
    loadInventory();
  }, []);

  async function loadInventory() {
    const res = await apiCall("/api/v1/inventory", { requireAuth: true });
    onTrackResponse(res);
    if (res.ok && Array.isArray(res.data) && res.data.length > 0) {
      setItems(res.data);
    } else {
      setItems(defaultItems);
    }
  }

  async function handleCreateItem() {
    if (!newCode.trim() || !newName.trim()) {
      return onShowAlert("Código y nombre del equipo son requeridos", "error");
    }

    const res = await apiCall("/api/v1/inventory", {
      method: "POST",
      body: {
        code: newCode,
        name: newName,
        category: newCategory,
        total_stock: parseInt(newStock) || 1,
        location: newLocation,
      },
      requireAuth: true,
    });
    onTrackResponse(res);

    if (res.ok) {
      onShowAlert("Equipo dado de alta en inventario correctamente", "success");
      setShowCreateModal(false);
      setNewCode("");
      setNewName("");
      loadInventory();
    } else {
      onShowAlert(res.data?.error || "Error al registrar ítem", "error");
    }
  }

  async function handleSubmitLoan() {
    if (!selectedItem) return;
    if (!loanReason.trim()) {
      return onShowAlert("Por favor indica el motivo del préstamo", "error");
    }

    const res = await apiCall("/api/v1/loans/requests", {
      method: "POST",
      body: {
        item_id: selectedItem.id,
        reason: loanReason,
        requested_start_date: loanStartDate ? new Date(loanStartDate).toISOString() : new Date().toISOString(),
        requested_end_date: loanEndDate ? new Date(loanEndDate).toISOString() : new Date(Date.now() + 7 * 86400000).toISOString(),
      },
      requireAuth: true,
    });
    onTrackResponse(res);

    if (res.ok) {
      onShowAlert(`Solicitud de préstamo para ${selectedItem.name} radicada con éxito`, "success");
      setShowLoanModal(false);
      setLoanReason("");
      loadInventory();
    } else {
      onShowAlert(res.data?.error || "Error al solicitar préstamo", "error");
    }
  }

  const categories = [
    "TODOS",
    "EQUIPAMIENTO",
    "HERRAMIENTA",
    "MODULO_ELECTRONICO",
    "CONSUMIBLE",
  ];

  const filteredItems = items.filter((it) => {
    const matchCat = categoryFilter === "TODOS" || it.category === categoryFilter;
    const matchSearch =
      search === "" ||
      it.name?.toLowerCase().includes(search.toLowerCase()) ||
      it.code?.toLowerCase().includes(search.toLowerCase());
    return matchCat && matchSearch;
  });

  return (
    <div className="space-y-6">
      {/* Cabecera */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h2 className="text-xl font-bold text-white flex items-center gap-2">
            <Package className="w-5 h-5 text-teal-400" />
            <span>Laboratorio & Inventario de Equipos</span>
          </h2>
          <p className="text-xs text-slate-400 mt-0.5">
            Instrumentación científica, osciloscopios, sensores y herramientas disponibles para préstamo.
          </p>
        </div>

        <button
          onClick={() => setShowCreateModal(true)}
          className="self-start md:self-auto text-xs font-semibold px-3.5 py-2 rounded-lg bg-teal-600 hover:bg-teal-500 text-white flex items-center gap-2 shadow-sm shadow-teal-500/20 transition"
        >
          <Plus className="w-4 h-4" />
          <span>Dar de Alta Equipo</span>
        </button>
      </div>

      {/* Barra de Filtros y Búsqueda */}
      <div className="flex flex-wrap items-center justify-between gap-3 bg-slate-900/80 p-3 rounded-xl border border-slate-800">
        <div className="flex items-center gap-1.5 overflow-x-auto text-xs">
          {categories.map((c) => (
            <button
              key={c}
              onClick={() => setCategoryFilter(c)}
              className={`px-3 py-1.5 rounded-lg font-medium transition whitespace-nowrap ${
                categoryFilter === c
                  ? "bg-teal-600 text-white shadow-sm"
                  : "bg-slate-950/60 border border-slate-800 text-slate-400 hover:text-slate-200"
              }`}
            >
              {c}
            </button>
          ))}
        </div>

        <div className="w-full md:w-64 relative">
          <Search className="w-3.5 h-3.5 text-slate-400 absolute left-3 top-2.5" />
          <input
            type="text"
            placeholder="Buscar por código o nombre..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="input-field pl-8 py-1.5 text-xs bg-slate-950/80"
          />
        </div>
      </div>

      {/* Grid de Ítems de Inventario */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {filteredItems.map((it) => {
          const isAvailable = (it.available_stock ?? it.total_stock) > 0;
          return (
            <div
              key={it.id}
              className="p-4 rounded-xl bg-slate-900/90 border border-slate-800 hover:border-slate-700 transition space-y-3 flex flex-col justify-between shadow-md"
            >
              <div className="space-y-2">
                <div className="flex items-center justify-between">
                  <span className="font-mono text-[10px] font-bold px-2 py-0.5 rounded bg-slate-800 text-teal-400 border border-slate-700">
                    {it.code}
                  </span>
                  <span
                    className={`text-[10px] font-bold px-2 py-0.5 rounded-full border ${
                      isAvailable
                        ? "bg-emerald-500/10 text-emerald-400 border-emerald-500/20"
                        : "bg-red-500/10 text-red-400 border-red-500/20"
                    }`}
                  >
                    {isAvailable ? `${it.available_stock ?? it.total_stock} disponibles` : "Todos prestados"}
                  </span>
                </div>

                <h3 className="text-xs font-bold text-white leading-snug">
                  {it.name}
                </h3>
                <p className="text-[11px] text-slate-400">
                  Ubicación: {it.location || "Laboratorio SEMARD"}
                </p>
              </div>

              <div className="pt-2 border-t border-slate-800/80 flex items-center justify-between">
                <span className="text-[10px] text-slate-500 font-medium">
                  Total en Lab: {it.total_stock}
                </span>
                <button
                  disabled={!isAvailable}
                  onClick={() => {
                    setSelectedItem(it);
                    setShowLoanModal(true);
                  }}
                  className={`text-xs font-semibold px-3 py-1.5 rounded-lg transition flex items-center gap-1 ${
                    isAvailable
                      ? "bg-teal-600/90 hover:bg-teal-500 text-white shadow-sm shadow-teal-500/20"
                      : "bg-slate-800 text-slate-500 cursor-not-allowed"
                  }`}
                >
                  <span>Pedir Préstamo</span>
                  <ArrowRight className="w-3 h-3" />
                </button>
              </div>
            </div>
          );
        })}
      </div>

      {/* Modal: Solicitar Préstamo */}
      {showLoanModal && selectedItem && (
        <div className="fixed inset-0 z-50 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-2xl max-w-md w-full p-6 space-y-4 shadow-2xl animate-in fade-in zoom-in-95">
            <div className="flex items-center justify-between">
              <h3 className="text-base font-bold text-white">Solicitud de Préstamo de Equipo</h3>
              <button
                onClick={() => setShowLoanModal(false)}
                className="text-slate-400 hover:text-white"
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            <div className="p-3 rounded-lg bg-slate-950/60 border border-slate-800 text-xs space-y-1">
              <p className="font-bold text-teal-400">{selectedItem.code}</p>
              <p className="text-white font-medium">{selectedItem.name}</p>
            </div>

            <div className="space-y-3 text-xs">
              <div>
                <label className="text-[11px] text-slate-400 block mb-1">Motivo o Proyecto</label>
                <textarea
                  placeholder="Ej: Pruebas de adquisición de señales para proyecto Rover..."
                  rows={2}
                  value={loanReason}
                  onChange={(e) => setLoanReason(e.target.value)}
                  className="input-field"
                />
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="text-[11px] text-slate-400 block mb-1">Fecha de Inicio</label>
                  <input
                    type="date"
                    value={loanStartDate}
                    onChange={(e) => setLoanStartDate(e.target.value)}
                    className="input-field"
                  />
                </div>
                <div>
                  <label className="text-[11px] text-slate-400 block mb-1">Fecha de Devolución</label>
                  <input
                    type="date"
                    value={loanEndDate}
                    onChange={(e) => setLoanEndDate(e.target.value)}
                    className="input-field"
                  />
                </div>
              </div>
            </div>

            <div className="flex gap-2 pt-2">
              <button
                onClick={() => setShowLoanModal(false)}
                className="btn-action flex-1 py-2"
              >
                Cancelar
              </button>
              <button
                onClick={handleSubmitLoan}
                className="btn-primary flex-1 py-2"
              >
                Radicar Préstamo
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Modal: Crear Ítem de Inventario */}
      {showCreateModal && (
        <div className="fixed inset-0 z-50 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-2xl max-w-md w-full p-6 space-y-4 shadow-2xl animate-in fade-in zoom-in-95">
            <div className="flex items-center justify-between">
              <h3 className="text-base font-bold text-white">Alta de Nuevo Equipo / Insumo</h3>
              <button
                onClick={() => setShowCreateModal(false)}
                className="text-slate-400 hover:text-white"
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            <div className="space-y-3 text-xs">
              <div className="grid grid-cols-2 gap-2">
                <div>
                  <label className="text-[11px] text-slate-400 block mb-1">Código de Inventario</label>
                  <input
                    type="text"
                    placeholder="LAB-XXX-01"
                    value={newCode}
                    onChange={(e) => setNewCode(e.target.value)}
                    className="input-field font-mono"
                  />
                </div>
                <div>
                  <label className="text-[11px] text-slate-400 block mb-1">Categoría</label>
                  <select
                    value={newCategory}
                    onChange={(e) => setNewCategory(e.target.value)}
                    className="input-field"
                  >
                    <option value="EQUIPAMIENTO">EQUIPAMIENTO</option>
                    <option value="HERRAMIENTA">HERRAMIENTA</option>
                    <option value="MODULO_ELECTRONICO">MODULO_ELECTRONICO</option>
                    <option value="CONSUMIBLE">CONSUMIBLE</option>
                  </select>
                </div>
              </div>

              <div>
                <label className="text-[11px] text-slate-400 block mb-1">Nombre Descriptivo</label>
                <input
                  type="text"
                  placeholder="Ej: Multímetro Digital Fluke 117"
                  value={newName}
                  onChange={(e) => setNewName(e.target.value)}
                  className="input-field"
                />
              </div>

              <div className="grid grid-cols-2 gap-2">
                <div>
                  <label className="text-[11px] text-slate-400 block mb-1">Stock Total Inicial</label>
                  <input
                    type="number"
                    value={newStock}
                    onChange={(e) => setNewStock(e.target.value)}
                    className="input-field"
                  />
                </div>
                <div>
                  <label className="text-[11px] text-slate-400 block mb-1">Ubicación en Laboratorio</label>
                  <input
                    type="text"
                    placeholder="Estante B - Repisa 1"
                    value={newLocation}
                    onChange={(e) => setNewLocation(e.target.value)}
                    className="input-field"
                  />
                </div>
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
                onClick={handleCreateItem}
                className="btn-primary flex-1 py-2"
              >
                Guardar en Inventario
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
