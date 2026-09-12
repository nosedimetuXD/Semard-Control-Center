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

type ResourceHandler struct {
	db             *pgxpool.Pool
	projectHandler *ProjectHandler
}

func NewResourceHandler(db *pgxpool.Pool, projectHandler *ProjectHandler) *ResourceHandler {
	return &ResourceHandler{db: db, projectHandler: projectHandler}
}

// ListProjectResources lista las solicitudes de recursos hechas para un proyecto
func (h *ResourceHandler) ListProjectResources(w http.ResponseWriter, r *http.Request) {
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

	query := `SELECT id, project_id, requested_by, resource_type, title, description, estimated_cost,
	                 status, director_feedback, reviewed_by, reviewed_at, created_at
	          FROM resource_requests
	          WHERE project_id = $1
	          ORDER BY created_at DESC`

	rows, err := h.db.Query(r.Context(), query, projectID)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar recursos del proyecto", err.Error())
		return
	}
	defer rows.Close()

	resources := make([]domain.ResourceRequest, 0)
	for rows.Next() {
		var res domain.ResourceRequest
		if err := rows.Scan(
			&res.ID, &res.ProjectID, &res.RequestedBy, &res.ResourceType, &res.Title, &res.Description,
			&res.EstimatedCost, &res.Status, &res.DirectorFeedback, &res.ReviewedBy, &res.ReviewedAt, &res.CreatedAt,
		); err != nil {
			ErrorWithDetails(w, http.StatusInternalServerError, "Error al leer datos de recurso", err.Error())
			return
		}
		resources = append(resources, res)
	}

	JSON(w, http.StatusOK, resources)
}

// CreateResourceRequest permite a un encargado del proyecto (o Director) solicitar un recurso
func (h *ResourceHandler) CreateResourceRequest(w http.ResponseWriter, r *http.Request) {
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

	// Verificar si el usuario está autorizado para este proyecto
	authorized, err := h.projectHandler.IsUserAuthorizedForProject(r.Context(), projectID, claims.UserID, claims.Role)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al verificar permisos", err.Error())
		return
	}
	if !authorized {
		Error(w, http.StatusForbidden, "Solo los integrantes encargados del proyecto o un Director pueden solicitar recursos")
		return
	}

	var req domain.CreateResourceRequestInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Payload JSON inválido")
		return
	}

	resType := domain.ResourceType(strings.ToUpper(string(req.ResourceType)))
	if resType != domain.ResourceDigital && resType != domain.ResourceEconomic &&
		resType != domain.ResourceKnowledge && resType != domain.ResourceHardware {
		Error(w, http.StatusBadRequest, "Tipo de recurso inválido. Debe ser: DIGITAL, ECONOMIC, KNOWLEDGE o HARDWARE")
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		Error(w, http.StatusBadRequest, "El título del recurso es obligatorio")
		return
	}
	if strings.TrimSpace(req.Description) == "" {
		Error(w, http.StatusBadRequest, "La descripción o justificación del recurso es obligatoria")
		return
	}
	if req.EstimatedCost < 0 {
		req.EstimatedCost = 0
	}

	query := `INSERT INTO resource_requests (project_id, requested_by, resource_type, title, description, estimated_cost, status)
	          VALUES ($1, $2, $3, $4, $5, $6, 'PENDING')
	          RETURNING id, project_id, requested_by, resource_type, title, description, estimated_cost, status, director_feedback, reviewed_by, reviewed_at, created_at`

	var res domain.ResourceRequest
	err = h.db.QueryRow(r.Context(), query, projectID, claims.UserID, resType, req.Title, req.Description, req.EstimatedCost).Scan(
		&res.ID, &res.ProjectID, &res.RequestedBy, &res.ResourceType, &res.Title, &res.Description,
		&res.EstimatedCost, &res.Status, &res.DirectorFeedback, &res.ReviewedBy, &res.ReviewedAt, &res.CreatedAt,
	)

	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al registrar solicitud de recurso", err.Error())
		return
	}

	JSON(w, http.StatusCreated, res)
}

// ListPendingResources retorna todas las solicitudes de recursos pendientes para evaluación de los Directores
func (h *ResourceHandler) ListPendingResources(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	query := `SELECT rr.id, rr.project_id, rr.requested_by, rr.resource_type, rr.title, rr.description,
	                 rr.estimated_cost, rr.status, rr.director_feedback, rr.reviewed_by, rr.reviewed_at, rr.created_at,
	                 p.title AS project_title, u.full_name AS requester_name
	          FROM resource_requests rr
	          JOIN projects p ON rr.project_id = p.id
	          JOIN users u ON rr.requested_by = u.id
	          WHERE rr.status = 'PENDING'
	          ORDER BY rr.created_at ASC`

	rows, err := h.db.Query(r.Context(), query)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar solicitudes de recursos pendientes", err.Error())
		return
	}
	defer rows.Close()

	type PendingResourceItem struct {
		domain.ResourceRequest
		ProjectTitle  string `json:"project_title"`
		RequesterName string `json:"requester_name"`
	}

	pendingList := make([]PendingResourceItem, 0)
	for rows.Next() {
		var item PendingResourceItem
		if err := rows.Scan(
			&item.ID, &item.ProjectID, &item.RequestedBy, &item.ResourceType, &item.Title, &item.Description,
			&item.EstimatedCost, &item.Status, &item.DirectorFeedback, &item.ReviewedBy, &item.ReviewedAt, &item.CreatedAt,
			&item.ProjectTitle, &item.RequesterName,
		); err != nil {
			ErrorWithDetails(w, http.StatusInternalServerError, "Error al leer datos de recursos", err.Error())
			return
		}
		pendingList = append(pendingList, item)
	}

	JSON(w, http.StatusOK, pendingList)
}

// ReviewResource permite a un Director resolver una solicitud en 3 estados: APPROVED, REJECTED, RETURNED_FOR_MODIFICATION
func (h *ResourceHandler) ReviewResource(w http.ResponseWriter, r *http.Request) {
	directorID, err := middleware.GetUserID(r.Context())
	if err != nil {
		Error(w, http.StatusUnauthorized, "Sesión no válida")
		return
	}

	resourceIDStr := chi.URLParam(r, "resourceId")
	resourceID, err := uuid.Parse(resourceIDStr)
	if err != nil {
		Error(w, http.StatusBadRequest, "ID de solicitud de recurso inválido")
		return
	}

	var req domain.ReviewResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Payload JSON inválido")
		return
	}

	// Validación estricta de los 3 estados permitidos para evaluación de recursos
	statusUpper := domain.ResourceStatus(strings.ToUpper(string(req.Status)))
	if statusUpper != domain.ResourceApproved &&
		statusUpper != domain.ResourceRejected &&
		statusUpper != domain.ResourceReturnedForModification {
		Error(w, http.StatusBadRequest, "Estado inválido. Debe ser: APPROVED, REJECTED o RETURNED_FOR_MODIFICATION")
		return
	}

	// Validación estricta: Retroalimentación o justificación del Director obligatoria
	if strings.TrimSpace(req.DirectorFeedback) == "" {
		Error(w, http.StatusBadRequest, "El feedback técnico del Director es estrictamente obligatorio para resolver la solicitud de recurso")
		return
	}

	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	query := `UPDATE resource_requests SET 
	            status = $1,
	            director_feedback = $2,
	            reviewed_by = $3,
	            reviewed_at = NOW()
	          WHERE id = $4
	          RETURNING id, project_id, requested_by, resource_type, title, description, estimated_cost, status, director_feedback, reviewed_by, reviewed_at, created_at`

	var res domain.ResourceRequest
	err = h.db.QueryRow(r.Context(), query, statusUpper, req.DirectorFeedback, directorID, resourceID).Scan(
		&res.ID, &res.ProjectID, &res.RequestedBy, &res.ResourceType, &res.Title, &res.Description,
		&res.EstimatedCost, &res.Status, &res.DirectorFeedback, &res.ReviewedBy, &res.ReviewedAt, &res.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			Error(w, http.StatusNotFound, "Solicitud de recurso no encontrada")
			return
		}
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al registrar evaluación de recurso", err.Error())
		return
	}

	JSON(w, http.StatusOK, res)
}
