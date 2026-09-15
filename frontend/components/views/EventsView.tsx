"use client";

import { useState, useEffect } from "react";
import {
  Calendar,
  Clock,
  MapPin,
  Users,
  Plus,
  Sparkles,
  CheckCircle2,
  X,
  Share2,
} from "lucide-react";
import { apiCall } from "@/lib/api";

interface EventsViewProps {
  user: any;
  onShowAlert: (msg: string, type: "info" | "success" | "error") => void;
  onTrackResponse: (res: any) => void;
}

export default function EventsView({
  user,
  onShowAlert,
  onTrackResponse,
}: EventsViewProps) {
  const [events, setEvents] = useState<any[]>([]);
  const [filter, setFilter] = useState<"ALL" | "PUBLIC" | "INTERNAL">("ALL");

  // Modales
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showRegisterModal, setShowRegisterModal] = useState(false);
  const [selectedEvent, setSelectedEvent] = useState<any | null>(null);

  // Form states Evento
  const [title, setTitle] = useState("");
  const [desc, setDesc] = useState("");
  const [vis, setVis] = useState("PUBLIC");
  const [loc, setLoc] = useState("");
  const [startTime, setStartTime] = useState("");

  // Form states Registro
  const [regName, setRegName] = useState("");
  const [regEmail, setRegEmail] = useState("");
  const [regPhone, setRegPhone] = useState("");

  const defaultEvents = [
    {
      id: "ev-1",
      title: "Taller Práctico: Control de Motores Brushless y ESC con ESP32",
      description: "Aprende el principio de modulación PWM, calibración de variadores de velocidad electrónicos (ESC) y control de motores de drones.",
      visibility: "PUBLIC",
      location: "Laboratorio de Robótica (Sede Piedra de Bolívar)",
      start_time: "2026-09-25T14:00:00Z",
      attendees_count: 18,
    },
    {
      id: "ev-2",
      title: "Seminario Interno: Mapeo y Localización Simultánea (SLAM) con ROS 2",
      description: "Sesión técnica exclusiva para miembros del semillero sobre algoritmos de odometría láser y fusión sensorial en el Rover autónomo.",
      visibility: "INTERNAL",
      location: "Sala de Seminarios de Ingeniería / Discord SEMARD",
      start_time: "2026-10-02T16:00:00Z",
      attendees_count: 12,
    },
    {
      id: "ev-3",
      title: "Convocatoria Abierta: Proyectos de Grado y Semilleros 2026-II",
      description: "Presentación de líneas de investigación activas para estudiantes de Ingeniería de Sistemas, Mecánica, Electrónica y Afines.",
      visibility: "PUBLIC",
      location: "Auditorio de Ingeniería - Campus Piedra de Bolívar",
      start_time: "2026-10-15T10:00:00Z",
      attendees_count: 35,
    },
  ];

  useEffect(() => {
    loadEvents();
  }, []);

  async function loadEvents() {
    const res = await apiCall("/api/v1/events/public");
    onTrackResponse(res);
    if (res.ok && Array.isArray(res.data) && res.data.length > 0) {
      setEvents(res.data);
    } else {
      setEvents(defaultEvents);
    }
  }

  async function handleCreateEvent() {
    if (!title.trim()) return onShowAlert("El título del evento es obligatorio", "error");

    const res = await apiCall("/api/v1/events", {
      method: "POST",
      body: {
        title,
        description: desc,
        visibility: vis,
        location: loc,
        start_time: startTime ? new Date(startTime).toISOString() : new Date().toISOString(),
      },
      requireAuth: true,
    });
    onTrackResponse(res);

    if (res.ok) {
      onShowAlert("Evento publicado con éxito en la cartelera", "success");
      setShowCreateModal(false);
      setTitle("");
      setDesc("");
      loadEvents();
    } else {
      onShowAlert(res.data?.error || "Error al publicar evento (Solo Directores)", "error");
    }
  }

  async function handleRegisterAttendee() {
    if (!selectedEvent) return;
    if (!regName.trim() || !regEmail.trim()) {
      return onShowAlert("Nombre y correo electrónico son obligatorios", "error");
    }

    const res = await apiCall(`/api/v1/events/${selectedEvent.id}/register`, {
      method: "POST",
      body: {
        full_name: regName,
        email: regEmail,
        phone: regPhone,
      },
    });
    onTrackResponse(res);

    if (res.ok) {
      onShowAlert(`¡Inscripción confirmada para ${selectedEvent.title}!`, "success");
      setShowRegisterModal(false);
      setRegName("");
      setRegEmail("");
      setRegPhone("");
    } else {
      onShowAlert(res.data?.error || "Error al inscribirse al evento", "error");
    }
  }

  const filteredEvents = events.filter((ev) => {
    if (filter === "ALL") return true;
    return ev.visibility === filter;
  });

  return (
    <div className="space-y-6">
      {/* Cabecera */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h2 className="text-xl font-bold text-white flex items-center gap-2">
            <Calendar className="w-5 h-5 text-teal-400" />
            <span>Cartelera de Eventos & Talleres</span>
          </h2>
          <p className="text-xs text-slate-400 mt-0.5">
            Talleres presenciales, seminarios técnicos y convocatorias abiertas de robótica y automatización.
          </p>
        </div>

        <button
          onClick={() => setShowCreateModal(true)}
          className="self-start md:self-auto text-xs font-semibold px-3.5 py-2 rounded-lg bg-teal-600 hover:bg-teal-500 text-white flex items-center gap-2 shadow-sm shadow-teal-500/20 transition"
        >
          <Plus className="w-4 h-4" />
          <span>Publicar Evento</span>
        </button>
      </div>

      {/* Selector de Filtros */}
      <div className="flex items-center gap-2 text-xs">
        <button
          onClick={() => setFilter("ALL")}
          className={`px-3 py-1.5 rounded-lg font-medium transition ${
            filter === "ALL"
              ? "bg-teal-600 text-white"
              : "bg-slate-900 border border-slate-800 text-slate-400 hover:text-slate-200"
          }`}
        >
          Todos los Eventos
        </button>
        <button
          onClick={() => setFilter("PUBLIC")}
          className={`px-3 py-1.5 rounded-lg font-medium transition ${
            filter === "PUBLIC"
              ? "bg-teal-600 text-white"
              : "bg-slate-900 border border-slate-800 text-slate-400 hover:text-slate-200"
          }`}
        >
          Públicos (Abiertos)
        </button>
        <button
          onClick={() => setFilter("INTERNAL")}
          className={`px-3 py-1.5 rounded-lg font-medium transition ${
            filter === "INTERNAL"
              ? "bg-teal-600 text-white"
              : "bg-slate-900 border border-slate-800 text-slate-400 hover:text-slate-200"
          }`}
        >
          Internos (Semillero)
        </button>
      </div>

      {/* Grid de Tarjetas de Eventos */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
        {filteredEvents.map((ev) => (
          <div
            key={ev.id}
            className="p-5 rounded-2xl bg-slate-900/90 border border-slate-800 hover:border-slate-700 transition space-y-4 flex flex-col justify-between shadow-lg"
          >
            <div className="space-y-2.5">
              <div className="flex items-center justify-between">
                <span
                  className={`text-[10px] font-bold uppercase tracking-wider px-2.5 py-0.5 rounded-full border ${
                    ev.visibility === "PUBLIC"
                      ? "bg-emerald-500/10 text-emerald-400 border-emerald-500/20"
                      : "bg-indigo-500/10 text-indigo-400 border-indigo-500/20"
                  }`}
                >
                  {ev.visibility === "PUBLIC" ? "Abierto a Todo Público" : "Semillero Interno"}
                </span>
                <span className="text-[11px] text-slate-500 flex items-center gap-1">
                  <Users className="w-3 h-3" />
                  <span>{ev.attendees_count || 0} inscritos</span>
                </span>
              </div>

              <h3 className="text-base font-bold text-white tracking-tight leading-snug">
                {ev.title}
              </h3>
              <p className="text-xs text-slate-300 leading-relaxed line-clamp-3">
                {ev.description}
              </p>
            </div>

            <div className="space-y-3 pt-3 border-t border-slate-800/80">
              <div className="space-y-1.5 text-xs text-slate-400">
                <div className="flex items-center gap-2">
                  <Clock className="w-3.5 h-3.5 text-teal-400 shrink-0" />
                  <span className="truncate">
                    {ev.start_time
                      ? new Date(ev.start_time).toLocaleDateString("es-CO", {
                          weekday: "short",
                          year: "numeric",
                          month: "short",
                          day: "numeric",
                          hour: "2-digit",
                          minute: "2-digit",
                        })
                      : "Fecha por definir"}
                  </span>
                </div>
                <div className="flex items-center gap-2">
                  <MapPin className="w-3.5 h-3.5 text-indigo-400 shrink-0" />
                  <span className="truncate">{ev.location || "Universidad de Cartagena"}</span>
                </div>
              </div>

              <button
                onClick={() => {
                  setSelectedEvent(ev);
                  setShowRegisterModal(true);
                }}
                className="w-full text-xs font-semibold py-2 px-3 rounded-lg bg-teal-600 hover:bg-teal-500 text-white shadow-md shadow-teal-500/10 transition"
              >
                Inscribirme al Evento
              </button>
            </div>
          </div>
        ))}
      </div>

      {/* Modal: Inscribirse al Evento */}
      {showRegisterModal && selectedEvent && (
        <div className="fixed inset-0 z-50 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-2xl max-w-md w-full p-6 space-y-4 shadow-2xl animate-in fade-in zoom-in-95">
            <div className="flex items-center justify-between">
              <h3 className="text-base font-bold text-white">Inscripción al Evento</h3>
              <button
                onClick={() => setShowRegisterModal(false)}
                className="text-slate-400 hover:text-white"
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            <div className="p-3 rounded-lg bg-slate-950/60 border border-slate-800 text-xs space-y-1">
              <p className="font-bold text-white">{selectedEvent.title}</p>
              <p className="text-slate-400">{selectedEvent.location}</p>
            </div>

            <div className="space-y-3 text-xs">
              <div>
                <label className="text-[11px] text-slate-400 block mb-1">Nombre Completo</label>
                <input
                  type="text"
                  placeholder="Tu nombre completo"
                  value={regName}
                  onChange={(e) => setRegName(e.target.value)}
                  className="input-field"
                />
              </div>

              <div>
                <label className="text-[11px] text-slate-400 block mb-1">Correo Electrónico</label>
                <input
                  type="email"
                  placeholder="ejemplo@unicartagena.edu.co"
                  value={regEmail}
                  onChange={(e) => setRegEmail(e.target.value)}
                  className="input-field"
                />
              </div>

              <div>
                <label className="text-[11px] text-slate-400 block mb-1">Teléfono o WhatsApp (Opcional)</label>
                <input
                  type="tel"
                  placeholder="+57 300..."
                  value={regPhone}
                  onChange={(e) => setRegPhone(e.target.value)}
                  className="input-field"
                />
              </div>
            </div>

            <div className="flex gap-2 pt-2">
              <button
                onClick={() => setShowRegisterModal(false)}
                className="btn-action flex-1 py-2"
              >
                Cancelar
              </button>
              <button
                onClick={handleRegisterAttendee}
                className="btn-primary flex-1 py-2"
              >
                Confirmar Inscripción
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Modal: Publicar Evento (Director) */}
      {showCreateModal && (
        <div className="fixed inset-0 z-50 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-2xl max-w-md w-full p-6 space-y-4 shadow-2xl animate-in fade-in zoom-in-95">
            <div className="flex items-center justify-between">
              <h3 className="text-base font-bold text-white">Publicar Nuevo Evento</h3>
              <button
                onClick={() => setShowCreateModal(false)}
                className="text-slate-400 hover:text-white"
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            <div className="space-y-3 text-xs">
              <div>
                <label className="text-[11px] text-slate-400 block mb-1">Título del Evento</label>
                <input
                  type="text"
                  placeholder="Ej: Workshop de Visión por Computadora con OpenCV"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  className="input-field"
                />
              </div>

              <div className="grid grid-cols-2 gap-2">
                <div>
                  <label className="text-[11px] text-slate-400 block mb-1">Audiencia</label>
                  <select
                    value={vis}
                    onChange={(e) => setVis(e.target.value)}
                    className="input-field"
                  >
                    <option value="PUBLIC">PUBLIC (Abierto)</option>
                    <option value="INTERNAL">INTERNAL (Semillero)</option>
                  </select>
                </div>
                <div>
                  <label className="text-[11px] text-slate-400 block mb-1">Lugar</label>
                  <input
                    type="text"
                    placeholder="Lab Robótica / Sala"
                    value={loc}
                    onChange={(e) => setLoc(e.target.value)}
                    className="input-field"
                  />
                </div>
              </div>

              <div>
                <label className="text-[11px] text-slate-400 block mb-1">Descripción</label>
                <textarea
                  placeholder="Temas que se abordarán, requisitos previos..."
                  rows={2}
                  value={desc}
                  onChange={(e) => setDesc(e.target.value)}
                  className="input-field"
                />
              </div>

              <div>
                <label className="text-[11px] text-slate-400 block mb-1">Fecha y Hora</label>
                <input
                  type="datetime-local"
                  value={startTime}
                  onChange={(e) => setStartTime(e.target.value)}
                  className="input-field"
                />
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
                onClick={handleCreateEvent}
                className="btn-primary flex-1 py-2"
              >
                Publicar en Cartelera
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
