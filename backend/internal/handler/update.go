package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nosedimetuXD/Semard-Control-Center/backend/internal/domain"
	"github.com/nosedimetuXD/Semard-Control-Center/backend/internal/middleware"
)

type UpdateHandler struct {
	db             *pgxpool.Pool
	projectHandler *ProjectHandler
}

func NewUpdateHandler(db *pgxpool.Pool, projectHandler *ProjectHandler) *UpdateHandler {
	return &UpdateHandler{db: db, projectHandler: projectHandler}
}

// ListUpdates lista los avances reportados para un proyecto
func (h *UpdateHandler) ListUpdates(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	idStr := chi.URLParam(r, "id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		Error(w, http.StatusBadRequest, "ID de proyecto inválido")
		return
	}

	query := `SELECT id, project_id, submitted_by, title, content, attachments_url, status, 
	                 director_feedback, reviewed_by, reviewed_at, created_at
	          FROM project_updates 
	          WHERE project_id = $1 
	          ORDER BY created_at DESC`

	rows, err := h.db.Query(r.Context(), query, projectID)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar avances del proyecto", err.Error())
		return
	}
	defer rows.Close()

	updates := make([]domain.ProjectUpdate, 0)
	for rows.Next() {
		var u domain.ProjectUpdate
		if err := rows.Scan(
			&u.ID, &u.ProjectID, &u.SubmittedBy, &u.Title, &u.Content, &u.AttachmentsURL,
			&u.Status, &u.DirectorFeedback, &u.ReviewedBy, &u.ReviewedAt, &u.CreatedAt,
		); err != nil {
			ErrorWithDetails(w, http.StatusInternalServerError, "Error al leer datos de avance", err.Error())
			return
		}
		updates = append(updates, u)
	}

	JSON(w, http.StatusOK, updates)
}

// CreateUpdate permite a un encargado del proyecto (o Director) subir un avance
func (h *UpdateHandler) CreateUpdate(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	idStr := chi.URLParam(r, "id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		Error(w, http.StatusBadRequest, "ID de proyecto inválido")
		return
	}

	claims, ok := middleware.GetUserClaims(r)
	if !ok || claims == nil {
		Error(w, http.StatusUnauthorized, "Sesión no válida")
		return
	}

	// Verificar si el usuario está autorizado para reportar avances en este proyecto
	authorized, err := h.projectHandler.IsUserAuthorizedForProject(r.Context(), projectID, claims.UserID, claims.Role)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al verificar permisos en el proyecto", err.Error())
		return
	}
	if !authorized {
		Error(w, http.StatusForbidden, "Solo los integrantes encargados asignados a este proyecto o un Director pueden enviar avances")
		return
	}

	var req domain.CreateUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Payload JSON inválido")
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		Error(w, http.StatusBadRequest, "El título del avance es obligatorio")
		return
	}
	if strings.TrimSpace(req.Content) == "" {
		Error(w, http.StatusBadRequest, "La descripción o contenido del avance es obligatorio")
		return
	}

	attachmentsJSON, err := json.Marshal(req.AttachmentsURL)
	if err != nil {
		attachmentsJSON = []byte("[]")
	}

	query := `INSERT INTO project_updates (project_id, submitted_by, title, content, attachments_url, status)
	          VALUES ($1, $2, $3, $4, $5, 'PENDING')
	          RETURNING id, project_id, submitted_by, title, content, attachments_url, status, director_feedback, reviewed_by, reviewed_at, created_at`

	var u domain.ProjectUpdate
	err = h.db.QueryRow(r.Context(), query, projectID, claims.UserID, req.Title, req.Content, attachmentsJSON).Scan(
		&u.ID, &u.ProjectID, &u.SubmittedBy, &u.Title, &u.Content, &u.AttachmentsURL,
		&u.Status, &u.DirectorFeedback, &u.ReviewedBy, &u.ReviewedAt, &u.CreatedAt,
	)

	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al registrar avance", err.Error())
		return
	}

	JSON(w, http.StatusCreated, u)
}

// ReviewUpdate permite a un Director evaluar un avance con feedback obligatorio
func (h *UpdateHandler) ReviewUpdate(w http.ResponseWriter, r *http.Request) {
	directorID, err := middleware.GetUserID(r.Context())
	if err != nil {
		Error(w, http.StatusUnauthorized, "Sesión no válida")
		return
	}

	updateIDStr := chi.URLParam(r, "updateId")
	updateID, err := uuid.Parse(updateIDStr)
	if err != nil {
		Error(w, http.StatusBadRequest, "ID de avance inválido")
		return
	}

	var req domain.ReviewUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Payload JSON inválido")
		return
	}

	// Validación estricta del estado
	statusUpper := domain.UpdateStatus(strings.ToUpper(string(req.Status)))
	if statusUpper != domain.UpdateApproved && statusUpper != domain.UpdateChangesRequested {
		Error(w, http.StatusBadRequest, "Estado inválido. Debe ser APPROVED o CHANGES_REQUESTED")
		return
	}

	// Validación ESTRICTA: Retroalimentación técnica obligatoria
	if strings.TrimSpace(req.DirectorFeedback) == "" {
		Error(w, http.StatusBadRequest, "La retroalimentación técnica del Director es estrictamente obligatoria para evaluar un avance")
		return
	}

	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	query := `UPDATE project_updates SET 
	            status = $1,
	            director_feedback = $2,
	            reviewed_by = $3,
	            reviewed_at = NOW()
	          WHERE id = $4
	          RETURNING id, project_id, submitted_by, title, content, attachments_url, status, director_feedback, reviewed_by, reviewed_at, created_at`

	var u domain.ProjectUpdate
	err = h.db.QueryRow(r.Context(), query, statusUpper, req.DirectorFeedback, directorID, updateID).Scan(
		&u.ID, &u.ProjectID, &u.SubmittedBy, &u.Title, &u.Content, &u.AttachmentsURL,
		&u.Status, &u.DirectorFeedback, &u.ReviewedBy, &u.ReviewedAt, &u.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			Error(w, http.StatusNotFound, "Avance no encontrado")
			return
		}
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al registrar revisión del avance", err.Error())
		return
	}

	JSON(w, http.StatusOK, u)
}
