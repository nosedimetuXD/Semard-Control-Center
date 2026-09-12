package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nosedimetuXD/Semard-Control-Center/backend/internal/domain"
	"github.com/nosedimetuXD/Semard-Control-Center/backend/internal/middleware"
)

type UserHandler struct {
	db *pgxpool.Pool
}

func NewUserHandler(db *pgxpool.Pool) *UserHandler {
	return &UserHandler{db: db}
}

// ListUsers retorna todos los miembros con soporte de filtros opcionales (role, status, search)
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	roleFilter := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("role")))
	statusFilter := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("status")))
	search := strings.TrimSpace(r.URL.Query().Get("search"))

	query := `SELECT id, google_id, email, student_code, full_name, role, can_operate_3d, avatar_url, bio, is_active, created_at, updated_at 
	          FROM users WHERE 1=1`
	var args []interface{}
	argIdx := 1

	if roleFilter != "" {
		query += fmt.Sprintf(" AND role = $%d", argIdx)
		args = append(args, roleFilter)
		argIdx++
	}

	if statusFilter == "active" {
		query += fmt.Sprintf(" AND is_active = $%d", argIdx)
		args = append(args, true)
		argIdx++
	} else if statusFilter == "inactive" {
		query += fmt.Sprintf(" AND is_active = $%d", argIdx)
		args = append(args, false)
		argIdx++
	}

	if search != "" {
		searchPattern := "%" + search + "%"
		query += fmt.Sprintf(" AND (full_name ILIKE $%d OR email ILIKE $%d OR student_code ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, searchPattern)
		argIdx++
	}

	query += " ORDER BY role ASC, full_name ASC"

	rows, err := h.db.Query(r.Context(), query, args...)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar usuarios", err.Error())
		return
	}
	defer rows.Close()

	users := make([]domain.User, 0)
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(
			&u.ID, &u.GoogleID, &u.Email, &u.StudentCode, &u.FullName, &u.Role,
			&u.CanOperate3D, &u.AvatarURL, &u.Bio, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			ErrorWithDetails(w, http.StatusInternalServerError, "Error al leer datos de usuario", err.Error())
			return
		}
		users = append(users, u)
	}

	JSON(w, http.StatusOK, users)
}

// GetUser retorna el perfil de un usuario por su ID
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	idStr := chi.URLParam(r, "id")
	targetID, err := uuid.Parse(idStr)
	if err != nil {
		Error(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}

	var u domain.User
	query := `SELECT id, google_id, email, student_code, full_name, role, can_operate_3d, avatar_url, bio, is_active, created_at, updated_at 
	          FROM users WHERE id = $1 LIMIT 1`

	err = h.db.QueryRow(r.Context(), query, targetID).Scan(
		&u.ID, &u.GoogleID, &u.Email, &u.StudentCode, &u.FullName, &u.Role,
		&u.CanOperate3D, &u.AvatarURL, &u.Bio, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			Error(w, http.StatusNotFound, "Usuario no encontrado")
			return
		}
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar usuario", err.Error())
		return
	}

	JSON(w, http.StatusOK, u)
}

// UpdateProfile actualiza el perfil propio del usuario autenticado
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		Error(w, http.StatusUnauthorized, "Sesión no válida")
		return
	}

	var req domain.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Payload JSON inválido")
		return
	}

	query := `UPDATE users SET 
	            bio = COALESCE($1, bio),
	            student_code = COALESCE($2, student_code),
	            avatar_url = COALESCE($3, avatar_url),
	            updated_at = NOW()
	          WHERE id = $4
	          RETURNING id, google_id, email, student_code, full_name, role, can_operate_3d, avatar_url, bio, is_active, created_at, updated_at`

	var u domain.User
	err = h.db.QueryRow(r.Context(), query, req.Bio, req.StudentCode, req.AvatarURL, userID).Scan(
		&u.ID, &u.GoogleID, &u.Email, &u.StudentCode, &u.FullName, &u.Role,
		&u.CanOperate3D, &u.AvatarURL, &u.Bio, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al actualizar perfil", err.Error())
		return
	}

	JSON(w, http.StatusOK, u)
}

// UpdateRole permite a un Director cambiar el rol de un miembro
func (h *UserHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	idStr := chi.URLParam(r, "id")
	targetID, err := uuid.Parse(idStr)
	if err != nil {
		Error(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}

	var req domain.UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Payload JSON inválido")
		return
	}

	roleUpper := domain.UserRole(strings.ToUpper(string(req.Role)))
	if roleUpper != domain.RoleDirector && roleUpper != domain.RoleAdministrador && roleUpper != domain.RoleMiembro {
		Error(w, http.StatusBadRequest, "Rol inválido. Debe ser DIRECTOR, ADMINISTRADOR o MIEMBRO")
		return
	}

	query := `UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2 
	          RETURNING id, google_id, email, student_code, full_name, role, can_operate_3d, avatar_url, bio, is_active, created_at, updated_at`

	var u domain.User
	err = h.db.QueryRow(r.Context(), query, roleUpper, targetID).Scan(
		&u.ID, &u.GoogleID, &u.Email, &u.StudentCode, &u.FullName, &u.Role,
		&u.CanOperate3D, &u.AvatarURL, &u.Bio, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			Error(w, http.StatusNotFound, "Usuario no encontrado")
			return
		}
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al actualizar rol", err.Error())
		return
	}

	JSON(w, http.StatusOK, u)
}

// UpdatePermissions permite a un Director habilitar/deshabilitar can_operate_3d
func (h *UserHandler) UpdatePermissions(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	idStr := chi.URLParam(r, "id")
	targetID, err := uuid.Parse(idStr)
	if err != nil {
		Error(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}

	var req domain.UpdatePermissionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Payload JSON inválido")
		return
	}

	query := `UPDATE users SET can_operate_3d = $1, updated_at = NOW() WHERE id = $2 
	          RETURNING id, google_id, email, student_code, full_name, role, can_operate_3d, avatar_url, bio, is_active, created_at, updated_at`

	var u domain.User
	err = h.db.QueryRow(r.Context(), query, req.CanOperate3D, targetID).Scan(
		&u.ID, &u.GoogleID, &u.Email, &u.StudentCode, &u.FullName, &u.Role,
		&u.CanOperate3D, &u.AvatarURL, &u.Bio, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			Error(w, http.StatusNotFound, "Usuario no encontrado")
			return
		}
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al actualizar permisos", err.Error())
		return
	}

	JSON(w, http.StatusOK, u)
}

// UpdateStatus permite a un Director activar o suspender a un usuario
func (h *UserHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	idStr := chi.URLParam(r, "id")
	targetID, err := uuid.Parse(idStr)
	if err != nil {
		Error(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}

	var req domain.UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Payload JSON inválido")
		return
	}

	query := `UPDATE users SET is_active = $1, updated_at = NOW() WHERE id = $2 
	          RETURNING id, google_id, email, student_code, full_name, role, can_operate_3d, avatar_url, bio, is_active, created_at, updated_at`

	var u domain.User
	err = h.db.QueryRow(r.Context(), query, req.IsActive, targetID).Scan(
		&u.ID, &u.GoogleID, &u.Email, &u.StudentCode, &u.FullName, &u.Role,
		&u.CanOperate3D, &u.AvatarURL, &u.Bio, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			Error(w, http.StatusNotFound, "Usuario no encontrado")
			return
		}
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al actualizar estado", err.Error())
		return
	}

	JSON(w, http.StatusOK, u)
}
