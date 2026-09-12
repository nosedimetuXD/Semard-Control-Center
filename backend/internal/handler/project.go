package handler

import (
	"context"
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

type ProjectHandler struct {
	db *pgxpool.Pool
}

func NewProjectHandler(db *pgxpool.Pool) *ProjectHandler {
	return &ProjectHandler{db: db}
}

// ListShowcase retorna la vitrina general de proyectos para cualquier miembro autenticado
func (h *ProjectHandler) ListShowcase(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	query := `SELECT id, title, description, research_line, status, created_by, start_date, target_end_date, created_at, updated_at
	          FROM projects 
	          ORDER BY created_at DESC`

	rows, err := h.db.Query(r.Context(), query)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar proyectos", err.Error())
		return
	}
	defer rows.Close()

	projects := make([]domain.Project, 0)
	for rows.Next() {
		var p domain.Project
		if err := rows.Scan(
			&p.ID, &p.Title, &p.Description, &p.ResearchLine, &p.Status,
			&p.CreatedBy, &p.StartDate, &p.TargetEndDate, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			ErrorWithDetails(w, http.StatusInternalServerError, "Error al leer proyecto", err.Error())
			return
		}
		projects = append(projects, p)
	}

	// Enriquecer con miembros asignados
	for i := range projects {
		members, err := h.getProjectMembers(r.Context(), projects[i].ID)
		if err == nil {
			projects[i].Members = members
		}
	}

	JSON(w, http.StatusOK, projects)
}

// ListMyProjects retorna los proyectos donde el usuario autenticado es integrante o líder
func (h *ProjectHandler) ListMyProjects(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		Error(w, http.StatusUnauthorized, "Sesión no válida")
		return
	}

	query := `SELECT p.id, p.title, p.description, p.research_line, p.status, p.created_by, p.start_date, p.target_end_date, p.created_at, p.updated_at
	          FROM projects p
	          JOIN project_members pm ON p.id = pm.project_id
	          WHERE pm.user_id = $1
	          ORDER BY p.updated_at DESC`

	rows, err := h.db.Query(r.Context(), query, userID)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar proyectos del usuario", err.Error())
		return
	}
	defer rows.Close()

	projects := make([]domain.Project, 0)
	for rows.Next() {
		var p domain.Project
		if err := rows.Scan(
			&p.ID, &p.Title, &p.Description, &p.ResearchLine, &p.Status,
			&p.CreatedBy, &p.StartDate, &p.TargetEndDate, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			ErrorWithDetails(w, http.StatusInternalServerError, "Error al leer proyecto", err.Error())
			return
		}
		projects = append(projects, p)
	}

	for i := range projects {
		members, err := h.getProjectMembers(r.Context(), projects[i].ID)
		if err == nil {
			projects[i].Members = members
		}
	}

	JSON(w, http.StatusOK, projects)
}

// GetProject retorna el detalle completo de un proyecto
func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
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

	var p domain.Project
	query := `SELECT id, title, description, research_line, status, created_by, start_date, target_end_date, created_at, updated_at
	          FROM projects WHERE id = $1 LIMIT 1`

	err = h.db.QueryRow(r.Context(), query, projectID).Scan(
		&p.ID, &p.Title, &p.Description, &p.ResearchLine, &p.Status,
		&p.CreatedBy, &p.StartDate, &p.TargetEndDate, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			Error(w, http.StatusNotFound, "Proyecto no encontrado")
			return
		}
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar proyecto", err.Error())
		return
	}

	members, err := h.getProjectMembers(r.Context(), p.ID)
	if err == nil {
		p.Members = members
	}

	JSON(w, http.StatusOK, p)
}

// CreateProject permite a un Director crear un nuevo proyecto
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	directorID, err := middleware.GetUserID(r.Context())
	if err != nil {
		Error(w, http.StatusUnauthorized, "Sesión no válida")
		return
	}

	var req domain.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Payload JSON inválido")
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		Error(w, http.StatusBadRequest, "El título del proyecto es obligatorio")
		return
	}
	if strings.TrimSpace(req.Description) == "" {
		Error(w, http.StatusBadRequest, "La descripción del proyecto es obligatoria")
		return
	}
	if strings.TrimSpace(req.ResearchLine) == "" {
		Error(w, http.StatusBadRequest, "La línea de investigación es obligatoria")
		return
	}
	if req.StartDate.IsZero() {
		req.StartDate = time.Now()
	}

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al iniciar transacción", err.Error())
		return
	}
	defer tx.Rollback(r.Context())

	insertQuery := `INSERT INTO projects (title, description, research_line, status, created_by, start_date, target_end_date)
	                VALUES ($1, $2, $3, 'ACTIVE', $4, $5, $6)
	                RETURNING id, title, description, research_line, status, created_by, start_date, target_end_date, created_at, updated_at`

	var p domain.Project
	err = tx.QueryRow(r.Context(), insertQuery,
		req.Title, req.Description, req.ResearchLine, directorID, req.StartDate, req.TargetEndDate,
	).Scan(
		&p.ID, &p.Title, &p.Description, &p.ResearchLine, &p.Status,
		&p.CreatedBy, &p.StartDate, &p.TargetEndDate, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al crear proyecto", err.Error())
		return
	}

	// Asignar integrantes especificados
	for _, m := range req.Members {
		if m.UserID != uuid.Nil {
			_, err = tx.Exec(r.Context(), `INSERT INTO project_members (project_id, user_id, is_lead) 
			                              VALUES ($1, $2, $3) ON CONFLICT (project_id, user_id) DO NOTHING`,
				p.ID, m.UserID, m.IsLead)
			if err != nil {
				ErrorWithDetails(w, http.StatusInternalServerError, "Error al asignar integrantes iniciales", err.Error())
				return
			}
		}
	}

	if err := tx.Commit(r.Context()); err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al confirmar creación del proyecto", err.Error())
		return
	}

	// Cargar miembros creados
	p.Members, _ = h.getProjectMembers(r.Context(), p.ID)

	JSON(w, http.StatusCreated, p)
}

// UpdateProject actualiza metadatos y estado del proyecto (Exclusivo Directores)
func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
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

	var req domain.UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Payload JSON inválido")
		return
	}

	query := `UPDATE projects SET 
	            title = COALESCE($1, title),
	            description = COALESCE($2, description),
	            research_line = COALESCE($3, research_line),
	            status = COALESCE($4, status),
	            start_date = COALESCE($5, start_date),
	            target_end_date = COALESCE($6, target_end_date),
	            updated_at = NOW()
	          WHERE id = $7
	          RETURNING id, title, description, research_line, status, created_by, start_date, target_end_date, created_at, updated_at`

	var p domain.Project
	err = h.db.QueryRow(r.Context(), query,
		req.Title, req.Description, req.ResearchLine, req.Status, req.StartDate, req.TargetEndDate, projectID,
	).Scan(
		&p.ID, &p.Title, &p.Description, &p.ResearchLine, &p.Status,
		&p.CreatedBy, &p.StartDate, &p.TargetEndDate, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			Error(w, http.StatusNotFound, "Proyecto no encontrado")
			return
		}
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al actualizar proyecto", err.Error())
		return
	}

	p.Members, _ = h.getProjectMembers(r.Context(), p.ID)
	JSON(w, http.StatusOK, p)
}

// AssignMember asigna o actualiza el rol de líder de un integrante en el proyecto (Exclusivo Directores)
func (h *ProjectHandler) AssignMember(w http.ResponseWriter, r *http.Request) {
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

	var req domain.AssignMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Payload JSON inválido")
		return
	}

	if req.UserID == uuid.Nil {
		Error(w, http.StatusBadRequest, "user_id es obligatorio")
		return
	}

	query := `INSERT INTO project_members (project_id, user_id, is_lead, assigned_at)
	          VALUES ($1, $2, $3, NOW())
	          ON CONFLICT (project_id, user_id) DO UPDATE SET is_lead = EXCLUDED.is_lead`

	_, err = h.db.Exec(r.Context(), query, projectID, req.UserID, req.IsLead)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al asignar integrante al proyecto", err.Error())
		return
	}

	JSON(w, http.StatusOK, map[string]string{
		"message": "Integrante asignado exitosamente",
	})
}

// RemoveMember remueve a un integrante del proyecto (Exclusivo Directores)
func (h *ProjectHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
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

	userIDStr := chi.URLParam(r, "userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		Error(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}

	result, err := h.db.Exec(r.Context(), "DELETE FROM project_members WHERE project_id = $1 AND user_id = $2", projectID, userID)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al remover integrante", err.Error())
		return
	}

	if result.RowsAffected() == 0 {
		Error(w, http.StatusNotFound, "El integrante no estaba asignado a este proyecto")
		return
	}

	JSON(w, http.StatusOK, map[string]string{
		"message": "Integrante desasignado del proyecto correctamente",
	})
}

// Helper para obtener miembros con su perfil cargado
func (h *ProjectHandler) getProjectMembers(ctx context.Context, projectID uuid.UUID) ([]domain.ProjectMember, error) {
	query := `SELECT pm.project_id, pm.user_id, pm.is_lead, pm.assigned_at,
	                 u.full_name, u.email, u.role, u.avatar_url
	          FROM project_members pm
	          JOIN users u ON pm.user_id = u.id
	          WHERE pm.project_id = $1
	          ORDER BY pm.is_lead DESC, u.full_name ASC`

	rows, err := h.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make([]domain.ProjectMember, 0)
	for rows.Next() {
		var pm domain.ProjectMember
		var fullName, email string
		var role domain.UserRole
		var avatarURL *string

		if err := rows.Scan(
			&pm.ProjectID, &pm.UserID, &pm.IsLead, &pm.AssignedAt,
			&fullName, &email, &role, &avatarURL,
		); err != nil {
			return nil, err
		}

		pm.User = &domain.User{
			ID:        pm.UserID,
			FullName:  fullName,
			Email:     email,
			Role:      role,
			AvatarURL: avatarURL,
		}
		members = append(members, pm)
	}

	return members, nil
}

// Helper para verificar si un usuario es encargado asignado a un proyecto o Director
func (h *ProjectHandler) IsUserAuthorizedForProject(ctx context.Context, projectID, userID uuid.UUID, userRole domain.UserRole) (bool, error) {
	if userRole == domain.RoleDirector {
		return true, nil
	}

	var count int
	err := h.db.QueryRow(ctx, "SELECT COUNT(*) FROM project_members WHERE project_id = $1 AND user_id = $2", projectID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
