package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nosedimetuXD/Semard-Control-Center/backend/internal/domain"
	"github.com/nosedimetuXD/Semard-Control-Center/backend/internal/middleware"
)

type EventHandler struct {
	db *pgxpool.Pool
}

func NewEventHandler(db *pgxpool.Pool) *EventHandler {
	return &EventHandler{db: db}
}

// ListPublicEvents retorna todos los eventos de visibilidad PUBLIC
func (h *EventHandler) ListPublicEvents(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	upcoming := r.URL.Query().Get("upcoming") == "true"
	query := `SELECT id, title, description, visibility, start_time, end_time, location, banner_url, created_by, created_at, updated_at
	          FROM events WHERE visibility = 'PUBLIC'`
	var args []interface{}

	if upcoming {
		query += " AND start_time >= NOW()"
	}
	query += " ORDER BY start_time ASC"

	rows, err := h.db.Query(r.Context(), query, args...)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar eventos públicos", err.Error())
		return
	}
	defer rows.Close()

	events := make([]domain.Event, 0)
	for rows.Next() {
		var e domain.Event
		if err := rows.Scan(
			&e.ID, &e.Title, &e.Description, &e.Visibility, &e.StartTime, &e.EndTime,
			&e.Location, &e.BannerURL, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			ErrorWithDetails(w, http.StatusInternalServerError, "Error al leer evento", err.Error())
			return
		}
		events = append(events, e)
	}

	JSON(w, http.StatusOK, events)
}

// ListInternalEvents retorna eventos para miembros del semillero (públicos e internos)
func (h *EventHandler) ListInternalEvents(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	query := `SELECT id, title, description, visibility, start_time, end_time, location, banner_url, created_by, created_at, updated_at
	          FROM events ORDER BY start_time DESC`

	rows, err := h.db.Query(r.Context(), query)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar eventos", err.Error())
		return
	}
	defer rows.Close()

	events := make([]domain.Event, 0)
	for rows.Next() {
		var e domain.Event
		if err := rows.Scan(
			&e.ID, &e.Title, &e.Description, &e.Visibility, &e.StartTime, &e.EndTime,
			&e.Location, &e.BannerURL, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			ErrorWithDetails(w, http.StatusInternalServerError, "Error al leer evento", err.Error())
			return
		}
		events = append(events, e)
	}

	JSON(w, http.StatusOK, events)
}

// GetEvent retorna el detalle de un evento específico
func (h *EventHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	idStr := chi.URLParam(r, "id")
	eventID, err := uuid.Parse(idStr)
	if err != nil {
		Error(w, http.StatusBadRequest, "ID de evento inválido")
		return
	}

	var e domain.Event
	query := `SELECT id, title, description, visibility, start_time, end_time, location, banner_url, created_by, created_at, updated_at
	          FROM events WHERE id = $1 LIMIT 1`

	err = h.db.QueryRow(r.Context(), query, eventID).Scan(
		&e.ID, &e.Title, &e.Description, &e.Visibility, &e.StartTime, &e.EndTime,
		&e.Location, &e.BannerURL, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			Error(w, http.StatusNotFound, "Evento no encontrado")
			return
		}
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar evento", err.Error())
		return
	}

	JSON(w, http.StatusOK, e)
}

// CreateEvent crea un nuevo evento (Exclusivo Directores)
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	directorID, err := middleware.GetUserID(r.Context())
	if err != nil {
		Error(w, http.StatusUnauthorized, "Sesión no válida")
		return
	}

	var req domain.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Payload JSON inválido")
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		Error(w, http.StatusBadRequest, "El título del evento es obligatorio")
		return
	}
	if strings.TrimSpace(req.Description) == "" {
		Error(w, http.StatusBadRequest, "La descripción del evento es obligatoria")
		return
	}
	if strings.TrimSpace(req.Location) == "" {
		Error(w, http.StatusBadRequest, "El lugar o enlace del evento es obligatorio")
		return
	}
	if req.StartTime.IsZero() {
		req.StartTime = time.Now()
	}

	visibility := domain.EventVisibility(strings.ToUpper(string(req.Visibility)))
	if visibility != domain.VisibilityPublic && visibility != domain.VisibilityInternal {
		visibility = domain.VisibilityPublic
	}

	query := `INSERT INTO events (title, description, visibility, start_time, end_time, location, banner_url, created_by)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	          RETURNING id, title, description, visibility, start_time, end_time, location, banner_url, created_by, created_at, updated_at`

	var e domain.Event
	err = h.db.QueryRow(r.Context(), query,
		req.Title, req.Description, visibility, req.StartTime, req.EndTime, req.Location, req.BannerURL, directorID,
	).Scan(
		&e.ID, &e.Title, &e.Description, &e.Visibility, &e.StartTime, &e.EndTime,
		&e.Location, &e.BannerURL, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt,
	)

	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al crear evento", err.Error())
		return
	}

	JSON(w, http.StatusCreated, e)
}

// UpdateEvent actualiza los datos de un evento (Exclusivo Directores)
func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	idStr := chi.URLParam(r, "id")
	eventID, err := uuid.Parse(idStr)
	if err != nil {
		Error(w, http.StatusBadRequest, "ID de evento inválido")
		return
	}

	var req domain.UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Payload JSON inválido")
		return
	}

	query := `UPDATE events SET 
	            title = COALESCE($1, title),
	            description = COALESCE($2, description),
	            visibility = COALESCE($3, visibility),
	            start_time = COALESCE($4, start_time),
	            end_time = COALESCE($5, end_time),
	            location = COALESCE($6, location),
	            banner_url = COALESCE($7, banner_url),
	            updated_at = NOW()
	          WHERE id = $8
	          RETURNING id, title, description, visibility, start_time, end_time, location, banner_url, created_by, created_at, updated_at`

	var e domain.Event
	err = h.db.QueryRow(r.Context(), query,
		req.Title, req.Description, req.Visibility, req.StartTime, req.EndTime, req.Location, req.BannerURL, eventID,
	).Scan(
		&e.ID, &e.Title, &e.Description, &e.Visibility, &e.StartTime, &e.EndTime,
		&e.Location, &e.BannerURL, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			Error(w, http.StatusNotFound, "Evento no encontrado")
			return
		}
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al actualizar evento", err.Error())
		return
	}

	JSON(w, http.StatusOK, e)
}

// DeleteEvent elimina un evento (Exclusivo Directores)
func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	idStr := chi.URLParam(r, "id")
	eventID, err := uuid.Parse(idStr)
	if err != nil {
		Error(w, http.StatusBadRequest, "ID de evento inválido")
		return
	}

	result, err := h.db.Exec(r.Context(), "DELETE FROM events WHERE id = $1", eventID)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al eliminar evento", err.Error())
		return
	}

	if result.RowsAffected() == 0 {
		Error(w, http.StatusNotFound, "Evento no encontrado")
		return
	}

	JSON(w, http.StatusOK, map[string]string{
		"message": "Evento eliminado correctamente",
	})
}

// RegisterAttendee registra a un asistente en un evento público
func (h *EventHandler) RegisterAttendee(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	idStr := chi.URLParam(r, "id")
	eventID, err := uuid.Parse(idStr)
	if err != nil {
		Error(w, http.StatusBadRequest, "ID de evento inválido")
		return
	}

	// Verificar que el evento exista y sea público
	var visibility domain.EventVisibility
	err = h.db.QueryRow(r.Context(), "SELECT visibility FROM events WHERE id = $1", eventID).Scan(&visibility)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			Error(w, http.StatusNotFound, "Evento no encontrado")
			return
		}
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al verificar evento", err.Error())
		return
	}

	if visibility != domain.VisibilityPublic {
		Error(w, http.StatusForbidden, "Solo es posible registrar asistentes en eventos públicos")
		return
	}

	var req domain.RegisterAttendeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Payload JSON inválido")
		return
	}

	if strings.TrimSpace(req.FullName) == "" {
		Error(w, http.StatusBadRequest, "El nombre completo es obligatorio")
		return
	}
	if strings.TrimSpace(req.Email) == "" || !strings.Contains(req.Email, "@") {
		Error(w, http.StatusBadRequest, "Correo electrónico válido es obligatorio")
		return
	}

	query := `INSERT INTO event_registrations (event_id, full_name, email, phone, institution)
	          VALUES ($1, $2, $3, $4, $5)
	          ON CONFLICT (event_id, email) DO UPDATE SET 
	            full_name = EXCLUDED.full_name,
	            phone = COALESCE(EXCLUDED.phone, event_registrations.phone),
	            institution = COALESCE(EXCLUDED.institution, event_registrations.institution)
	          RETURNING id, event_id, full_name, email, phone, institution, registered_at`

	var reg domain.EventRegistration
	err = h.db.QueryRow(r.Context(), query, eventID, req.FullName, req.Email, req.Phone, req.Institution).Scan(
		&reg.ID, &reg.EventID, &reg.FullName, &reg.Email, &reg.Phone, &reg.Institution, &reg.RegisteredAt,
	)

	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al registrar asistente", err.Error())
		return
	}

	JSON(w, http.StatusCreated, map[string]interface{}{
		"message":      "Inscripción al evento confirmada exitosamente",
		"registration": reg,
	})
}

// ListAttendees lista los inscritos a un evento (Directores y Administradores)
func (h *EventHandler) ListAttendees(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	idStr := chi.URLParam(r, "id")
	eventID, err := uuid.Parse(idStr)
	if err != nil {
		Error(w, http.StatusBadRequest, "ID de evento inválido")
		return
	}

	query := `SELECT id, event_id, full_name, email, phone, institution, registered_at
	          FROM event_registrations WHERE event_id = $1 ORDER BY registered_at DESC`

	rows, err := h.db.Query(r.Context(), query, eventID)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar inscritos", err.Error())
		return
	}
	defer rows.Close()

	registrations := make([]domain.EventRegistration, 0)
	for rows.Next() {
		var reg domain.EventRegistration
		if err := rows.Scan(
			&reg.ID, &reg.EventID, &reg.FullName, &reg.Email, &reg.Phone, &reg.Institution, &reg.RegisteredAt,
		); err != nil {
			ErrorWithDetails(w, http.StatusInternalServerError, "Error al leer registro de asistente", err.Error())
			return
		}
		registrations = append(registrations, reg)
	}

	JSON(w, http.StatusOK, map[string]interface{}{
		"event_id":  eventID,
		"total":     len(registrations),
		"attendees": registrations,
	})
}
