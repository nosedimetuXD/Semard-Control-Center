package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nosedimetuXD/Semard-Control-Center/backend/internal/domain"
	"github.com/nosedimetuXD/Semard-Control-Center/backend/internal/middleware"
)

type Print3DHandler struct {
	db          *pgxpool.Pool
	storagePath string
}

func NewPrint3DHandler(db *pgxpool.Pool, storagePath string) *Print3DHandler {
	cleanPath := filepath.Clean(storagePath)
	_ = os.MkdirAll(filepath.Join(cleanPath, "3d_prints"), 0755)
	return &Print3DHandler{
		db:          db,
		storagePath: cleanPath,
	}
}

// CreateRequest sube un archivo 3D (.stl, .obj, .3mf, .step) y registra la solicitud
func (h *Print3DHandler) CreateRequest(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		Error(w, http.StatusUnauthorized, "Sesión no válida")
		return
	}

	// Límite de subida: 50MB
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		Error(w, http.StatusBadRequest, "El archivo supera el tamaño máximo permitido de 50MB o el formulario es inválido")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		Error(w, http.StatusBadRequest, "El archivo 3D es obligatorio (campo 'file')")
		return
	}
	defer file.Close()

	// Validar extensión de manufactura aditiva
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".stl" && ext != ".obj" && ext != ".3mf" && ext != ".step" && ext != ".stp" {
		Error(w, http.StatusBadRequest, "Formato no permitido. Solo se aceptan archivos .stl, .obj, .3mf o .step")
		return
	}

	material := strings.ToUpper(strings.TrimSpace(r.FormValue("material")))
	if material == "" {
		material = "PLA"
	}

	color := strings.TrimSpace(r.FormValue("color"))
	var colorPtr *string
	if color != "" {
		colorPtr = &color
	}

	infill := 20
	if infillStr := r.FormValue("infill_percentage"); infillStr != "" {
		if val, err := strconv.Atoi(infillStr); err == nil {
			if val < 5 {
				val = 5
			} else if val > 100 {
				val = 100
			}
			infill = val
		}
	}

	notes := strings.TrimSpace(r.FormValue("notes"))
	var notesPtr *string
	if notes != "" {
		notesPtr = &notes
	}

	var projectIDPtr *uuid.UUID
	if projectIDStr := strings.TrimSpace(r.FormValue("project_id")); projectIDStr != "" {
		if pID, err := uuid.Parse(projectIDStr); err == nil {
			projectIDPtr = &pID
		}
	}

	// Guardar archivo físico en disco persistente
	uploadDir := filepath.Join(h.storagePath, "3d_prints")
	_ = os.MkdirAll(uploadDir, 0755)

	storedFileName := fmt.Sprintf("%s_%s", uuid.New().String(), filepath.Base(header.Filename))
	destinationPath := filepath.Join(uploadDir, storedFileName)

	destFile, err := os.Create(destinationPath)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al almacenar archivo en el servidor", err.Error())
		return
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, file); err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al escribir archivo en disco", err.Error())
		return
	}

	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	query := `INSERT INTO print3d_requests (requested_by, project_id, file_path, file_name, material, color, infill_percentage, notes, status)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'PENDING')
	          RETURNING id, requested_by, project_id, file_path, file_name, material, color, infill_percentage, notes, status,
	                    reviewer_feedback, approved_by, operator_id, approved_at, created_at, updated_at`

	var printReq domain.Print3DRequest
	err = h.db.QueryRow(r.Context(), query,
		userID, projectIDPtr, destinationPath, header.Filename, material, colorPtr, infill, notesPtr,
	).Scan(
		&printReq.ID, &printReq.RequestedBy, &printReq.ProjectID, &printReq.FilePath, &printReq.FileName,
		&printReq.Material, &printReq.Color, &printReq.InfillPercentage, &printReq.Notes, &printReq.Status,
		&printReq.ReviewerFeedback, &printReq.ApprovedBy, &printReq.OperatorID, &printReq.ApprovedAt,
		&printReq.CreatedAt, &printReq.UpdatedAt,
	)

	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al registrar solicitud en base de datos", err.Error())
		return
	}

	JSON(w, http.StatusCreated, printReq)
}

// ListMyRequests retorna las impresiones solicitadas por el usuario actual
func (h *Print3DHandler) ListMyRequests(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		Error(w, http.StatusUnauthorized, "Sesión no válida")
		return
	}

	query := `SELECT id, requested_by, project_id, file_path, file_name, material, color, infill_percentage, notes, status,
	                 reviewer_feedback, approved_by, operator_id, approved_at, created_at, updated_at
	          FROM print3d_requests
	          WHERE requested_by = $1
	          ORDER BY created_at DESC`

	rows, err := h.db.Query(r.Context(), query, userID)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar solicitudes de impresión", err.Error())
		return
	}
	defer rows.Close()

	requests := make([]domain.Print3DRequest, 0)
	for rows.Next() {
		var pr domain.Print3DRequest
		if err := rows.Scan(
			&pr.ID, &pr.RequestedBy, &pr.ProjectID, &pr.FilePath, &pr.FileName,
			&pr.Material, &pr.Color, &pr.InfillPercentage, &pr.Notes, &pr.Status,
			&pr.ReviewerFeedback, &pr.ApprovedBy, &pr.OperatorID, &pr.ApprovedAt,
			&pr.CreatedAt, &pr.UpdatedAt,
		); err != nil {
			ErrorWithDetails(w, http.StatusInternalServerError, "Error al leer datos de impresión", err.Error())
			return
		}
		requests = append(requests, pr)
	}

	JSON(w, http.StatusOK, requests)
}

// ListQueue lista la cola activa de impresión del laboratorio (visibilidad del taller)
func (h *Print3DHandler) ListQueue(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	query := `SELECT pr.id, pr.requested_by, pr.project_id, pr.file_path, pr.file_name, pr.material, pr.color,
	                 pr.infill_percentage, pr.notes, pr.status, pr.reviewer_feedback, pr.approved_by, pr.operator_id,
	                 pr.approved_at, pr.created_at, pr.updated_at,
	                 u.full_name AS requester_name, op.full_name AS operator_name
	          FROM print3d_requests pr
	          JOIN users u ON pr.requested_by = u.id
	          LEFT JOIN users op ON pr.operator_id = op.id
	          WHERE pr.status IN ('APPROVED', 'IN_PROGRESS', 'COMPLETED')
	          ORDER BY CASE pr.status 
	            WHEN 'IN_PROGRESS' THEN 1 
	            WHEN 'APPROVED' THEN 2 
	            WHEN 'COMPLETED' THEN 3 
	            ELSE 4 END, pr.approved_at ASC`

	rows, err := h.db.Query(r.Context(), query)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar cola de impresión", err.Error())
		return
	}
	defer rows.Close()

	type QueueItem struct {
		domain.Print3DRequest
		RequesterName string  `json:"requester_name"`
		OperatorName  *string `json:"operator_name,omitempty"`
	}

	queue := make([]QueueItem, 0)
	for rows.Next() {
		var item QueueItem
		if err := rows.Scan(
			&item.ID, &item.RequestedBy, &item.ProjectID, &item.FilePath, &item.FileName,
			&item.Material, &item.Color, &item.InfillPercentage, &item.Notes, &item.Status,
			&item.ReviewerFeedback, &item.ApprovedBy, &item.OperatorID, &item.ApprovedAt,
			&item.CreatedAt, &item.UpdatedAt, &item.RequesterName, &item.OperatorName,
		); err != nil {
			ErrorWithDetails(w, http.StatusInternalServerError, "Error al leer cola de impresión", err.Error())
			return
		}
		queue = append(queue, item)
	}

	JSON(w, http.StatusOK, queue)
}

// DownloadFile permite descargar el archivo 3D (para el dueño o el equipo del taller)
func (h *Print3DHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	claims, ok := middleware.GetUserClaims(r)
	if !ok || claims == nil {
		Error(w, http.StatusUnauthorized, "Sesión no válida")
		return
	}

	idStr := chi.URLParam(r, "id")
	reqID, err := uuid.Parse(idStr)
	if err != nil {
		Error(w, http.StatusBadRequest, "ID de impresión inválido")
		return
	}

	var filePath, fileName string
	var requestedBy uuid.UUID
	err = h.db.QueryRow(r.Context(), "SELECT file_path, file_name, requested_by FROM print3d_requests WHERE id = $1", reqID).Scan(&filePath, &fileName, &requestedBy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			Error(w, http.StatusNotFound, "Solicitud de impresión no encontrada")
			return
		}
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar archivo", err.Error())
		return
	}

	// Verificar permisos: dueño, operador 3D, admin o director
	isOwner := claims.UserID == requestedBy
	isAuthorizedTech := claims.Role == domain.RoleDirector || claims.Role == domain.RoleAdministrador || claims.CanOperate3D

	if !isOwner && !isAuthorizedTech {
		Error(w, http.StatusForbidden, "No tienes permiso para descargar este archivo")
		return
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		Error(w, http.StatusNotFound, "El archivo físico ya no se encuentra en el servidor")
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filePath)
}

// ReviewRequest aprueba o rechaza la solicitud de impresión (Admin o Director)
func (h *Print3DHandler) ReviewRequest(w http.ResponseWriter, r *http.Request) {
	reviewerID, err := middleware.GetUserID(r.Context())
	if err != nil {
		Error(w, http.StatusUnauthorized, "Sesión no válida")
		return
	}

	idStr := chi.URLParam(r, "id")
	reqID, err := uuid.Parse(idStr)
	if err != nil {
		Error(w, http.StatusBadRequest, "ID de solicitud inválido")
		return
	}

	var req domain.ReviewPrint3DRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Payload JSON inválido")
		return
	}

	statusUpper := domain.Print3DStatus(strings.ToUpper(string(req.Status)))
	if statusUpper != domain.Print3DApproved && statusUpper != domain.Print3DRejected {
		Error(w, http.StatusBadRequest, "Estado inválido. Debe ser APPROVED o REJECTED")
		return
	}

	if statusUpper == domain.Print3DRejected {
		if req.ReviewerFeedback == nil || strings.TrimSpace(*req.ReviewerFeedback) == "" {
			Error(w, http.StatusBadRequest, "El motivo de rechazo es obligatorio")
			return
		}
	}

	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	query := `UPDATE print3d_requests SET
	            status = $1,
	            reviewer_feedback = $2,
	            approved_by = $3,
	            approved_at = NOW(),
	            updated_at = NOW()
	          WHERE id = $4
	          RETURNING id, requested_by, project_id, file_path, file_name, material, color, infill_percentage, notes, status,
	                    reviewer_feedback, approved_by, operator_id, approved_at, created_at, updated_at`

	var printReq domain.Print3DRequest
	err = h.db.QueryRow(r.Context(), query, statusUpper, req.ReviewerFeedback, reviewerID, reqID).Scan(
		&printReq.ID, &printReq.RequestedBy, &printReq.ProjectID, &printReq.FilePath, &printReq.FileName,
		&printReq.Material, &printReq.Color, &printReq.InfillPercentage, &printReq.Notes, &printReq.Status,
		&printReq.ReviewerFeedback, &printReq.ApprovedBy, &printReq.OperatorID, &printReq.ApprovedAt,
		&printReq.CreatedAt, &printReq.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			Error(w, http.StatusNotFound, "Solicitud no encontrada")
			return
		}
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al evaluar solicitud", err.Error())
		return
	}

	JSON(w, http.StatusOK, printReq)
}

// UpdateStatus actualiza el estado de fabricación (Operador 3D, Admin o Director)
func (h *Print3DHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	operatorID, err := middleware.GetUserID(r.Context())
	if err != nil {
		Error(w, http.StatusUnauthorized, "Sesión no válida")
		return
	}

	idStr := chi.URLParam(r, "id")
	reqID, err := uuid.Parse(idStr)
	if err != nil {
		Error(w, http.StatusBadRequest, "ID de solicitud inválido")
		return
	}

	var req domain.UpdatePrint3DStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Payload JSON inválido")
		return
	}

	statusUpper := domain.Print3DStatus(strings.ToUpper(string(req.Status)))
	if statusUpper != domain.Print3DInProgress &&
		statusUpper != domain.Print3DCompleted &&
		statusUpper != domain.Print3DDelivered {
		Error(w, http.StatusBadRequest, "Estado operativo inválido. Debe ser: IN_PROGRESS, COMPLETED o DELIVERED")
		return
	}

	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	// Al pasar a IN_PROGRESS, se asigna formalmente el operador actual
	query := `UPDATE print3d_requests SET
	            status = $1,
	            notes = COALESCE($2, notes),
	            operator_id = COALESCE(operator_id, $3),
	            updated_at = NOW()
	          WHERE id = $4
	          RETURNING id, requested_by, project_id, file_path, file_name, material, color, infill_percentage, notes, status,
	                    reviewer_feedback, approved_by, operator_id, approved_at, created_at, updated_at`

	var printReq domain.Print3DRequest
	err = h.db.QueryRow(r.Context(), query, statusUpper, req.Notes, operatorID, reqID).Scan(
		&printReq.ID, &printReq.RequestedBy, &printReq.ProjectID, &printReq.FilePath, &printReq.FileName,
		&printReq.Material, &printReq.Color, &printReq.InfillPercentage, &printReq.Notes, &printReq.Status,
		&printReq.ReviewerFeedback, &printReq.ApprovedBy, &printReq.OperatorID, &printReq.ApprovedAt,
		&printReq.CreatedAt, &printReq.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			Error(w, http.StatusNotFound, "Solicitud no encontrada")
			return
		}
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al actualizar estado operativo", err.Error())
		return
	}

	JSON(w, http.StatusOK, printReq)
}
