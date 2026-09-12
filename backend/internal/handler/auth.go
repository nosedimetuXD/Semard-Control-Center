package handler

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nosedimetuXD/Semard-Control-Center/backend/internal/config"
	"github.com/nosedimetuXD/Semard-Control-Center/backend/internal/domain"
	"github.com/nosedimetuXD/Semard-Control-Center/backend/internal/middleware"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthHandler struct {
	cfg         *config.Config
	db          *pgxpool.Pool
	oauthConfig *oauth2.Config
}

func NewAuthHandler(cfg *config.Config, db *pgxpool.Pool) *AuthHandler {
	oauthConf := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &AuthHandler{
		cfg:         cfg,
		db:          db,
		oauthConfig: oauthConf,
	}
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	HD            string `json:"hd"` // Dominio de Google Workspace (ej. unicartagena.edu.co)
}

type LoginResponse struct {
	Status  string       `json:"status"` // "AUTHENTICATED", "NEEDS_REGISTRATION", "PENDING_APPROVAL"
	Message string       `json:"message"`
	Token   string       `json:"token,omitempty"`
	User    *domain.User `json:"user,omitempty"`
	Google  *struct {
		GoogleID string `json:"google_id"`
		Email    string `json:"email"`
		FullName string `json:"full_name"`
		Avatar   string `json:"avatar"`
	} `json:"google_data,omitempty"`
}

// GoogleLogin inicia el flujo OAuth y devuelve la URL o redirige al navegador
func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	// Solicitar prompt select_account para permitir elegir cuenta institucional
	url := h.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "select_account"))

	// Si se accede desde un navegador web directamente o pide redirect=true, redirigir
	if r.URL.Query().Get("redirect") == "true" || strings.Contains(r.Header.Get("Accept"), "text/html") {
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}

	JSON(w, http.StatusOK, map[string]string{
		"auth_url": url,
		"state":    state,
	})
}

// GoogleCallback procesa el código devuelto por Google
func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		Error(w, http.StatusBadRequest, "Código de autorización no provisto")
		return
	}

	ctx := r.Context()
	token, err := h.oauthConfig.Exchange(ctx, code)
	if err != nil {
		ErrorWithDetails(w, http.StatusBadRequest, "Error al intercambiar código con Google", err.Error())
		return
	}

	// Obtener información del perfil desde Google
	client := h.oauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar perfil en Google", err.Error())
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var gUser GoogleUserInfo
	if err := json.Unmarshal(bodyBytes, &gUser); err != nil {
		Error(w, http.StatusInternalServerError, "Error al procesar datos del usuario de Google")
		return
	}

	// Validar dominio universitario estricto
	if !strings.HasSuffix(strings.ToLower(gUser.Email), strings.ToLower(h.cfg.AllowedEmailDomain)) {
		Error(w, http.StatusForbidden, fmt.Sprintf("Solo se permite acceso con correos institucionales (%s)", h.cfg.AllowedEmailDomain))
		return
	}

	// Fallback de desarrollo: Si la base de datos aún no está conectada localmente
	if h.db == nil {
		JSON(w, http.StatusOK, map[string]interface{}{
			"status":  "GOOGLE_AUTH_SUCCESS_TEST_MODE",
			"message": "Autenticación exitosa con Google y validación del dominio @unicartagena.edu.co aprobada.",
			"google_profile": map[string]interface{}{
				"google_id":      gUser.ID,
				"email":          gUser.Email,
				"full_name":      gUser.Name,
				"picture_url":    gUser.Picture,
				"verified_email": gUser.VerifiedEmail,
				"hosted_domain":  gUser.HD,
			},
			"database_status": "disconnected (conecta PostgreSQL para persistencia de usuarios)",
		})
		return
	}

	// 1. Verificar si el usuario ya existe en la base de datos
	var u domain.User
	query := `SELECT id, google_id, email, student_code, full_name, role, can_operate_3d, avatar_url, bio, is_active, created_at, updated_at 
	          FROM users WHERE email = $1 LIMIT 1`

	err = h.db.QueryRow(ctx, query, gUser.Email).Scan(
		&u.ID, &u.GoogleID, &u.Email, &u.StudentCode, &u.FullName, &u.Role,
		&u.CanOperate3D, &u.AvatarURL, &u.Bio, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)

	if err == nil {
		// Usuario encontrado
		if !u.IsActive {
			Error(w, http.StatusForbidden, "Tu cuenta se encuentra inactiva. Contacta a un Director.")
			return
		}

		// Actualizar google_id o avatar si estaban vacíos
		if u.GoogleID == nil || *u.GoogleID == "" {
			_, _ = h.db.Exec(ctx, "UPDATE users SET google_id = $1, avatar_url = COALESCE(avatar_url, $2), updated_at = NOW() WHERE id = $3", gUser.ID, gUser.Picture, u.ID)
		}

		jwtToken, err := h.generateJWT(&u)
		if err != nil {
			Error(w, http.StatusInternalServerError, "Error al generar credencial de sesión")
			return
		}

		JSON(w, http.StatusOK, LoginResponse{
			Status:  "AUTHENTICATED",
			Message: "Autenticación exitosa",
			Token:   jwtToken,
			User:    &u,
		})
		return
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar usuario en base de datos", err.Error())
		return
	}

	// 2. El usuario NO existe: verificar si ya tiene una solicitud pendiente
	var reqStatus domain.RegistrationStatus
	var feedback *string
	err = h.db.QueryRow(ctx, "SELECT status, director_feedback FROM registration_requests WHERE email = $1 ORDER BY created_at DESC LIMIT 1", gUser.Email).Scan(&reqStatus, &feedback)

	if err == nil {
		if reqStatus == domain.RegPending {
			JSON(w, http.StatusOK, LoginResponse{
				Status:  "PENDING_APPROVAL",
				Message: "Tu solicitud de membresía ya fue enviada y está en proceso de revisión por los Directores.",
			})
			return
		}
		if reqStatus == domain.RegRejected {
			msg := "Tu solicitud de ingreso previa fue rechazada por los Directores."
			if feedback != nil && *feedback != "" {
				msg += " Motivo: " + *feedback
			}
			JSON(w, http.StatusForbidden, LoginResponse{
				Status:  "REJECTED",
				Message: msg,
			})
			return
		}
	}

	// 3. No está registrado ni tiene solicitud pendiente: solicita completar datos
	JSON(w, http.StatusOK, LoginResponse{
		Status:  "NEEDS_REGISTRATION",
		Message: "Correo institucional verificado. Completa tu código estudiantil para solicitar ingreso al semillero.",
		Google: &struct {
			GoogleID string `json:"google_id"`
			Email    string `json:"email"`
			FullName string `json:"full_name"`
			Avatar   string `json:"avatar"`
		}{
			GoogleID: gUser.ID,
			Email:    gUser.Email,
			FullName: gUser.Name,
			Avatar:   gUser.Picture,
		},
	})
}

type RegisterRequestBody struct {
	GoogleID         string  `json:"google_id"`
	Email            string  `json:"email"`
	FullName         string  `json:"full_name"`
	StudentCode      string  `json:"student_code"`
	CareerProgram    *string `json:"career_program"`
	MotivationLetter *string `json:"motivation_letter"`
}

// RegisterRequest registra una solicitud de ingreso para ser evaluada por Directores
func (h *AuthHandler) RegisterRequest(w http.ResponseWriter, r *http.Request) {
	var body RegisterRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		Error(w, http.StatusBadRequest, "Datos de solicitud inválidos")
		return
	}

	if body.GoogleID == "" || body.Email == "" || body.FullName == "" || body.StudentCode == "" {
		Error(w, http.StatusBadRequest, "Código estudiantil, nombre, email y Google ID son obligatorios")
		return
	}

	if !strings.HasSuffix(strings.ToLower(body.Email), strings.ToLower(h.cfg.AllowedEmailDomain)) {
		Error(w, http.StatusForbidden, "El correo debe pertenecer a la Universidad de Cartagena")
		return
	}

	ctx := r.Context()
	query := `INSERT INTO registration_requests (google_id, email, full_name, student_code, career_program, motivation_letter, status)
	          VALUES ($1, $2, $3, $4, $5, $6, 'PENDING')
	          RETURNING id, created_at`

	var reqID uuid.UUID
	var createdAt time.Time
	err := h.db.QueryRow(ctx, query, body.GoogleID, body.Email, body.FullName, body.StudentCode, body.CareerProgram, body.MotivationLetter).Scan(&reqID, &createdAt)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "No se pudo registrar la solicitud de membresía", err.Error())
		return
	}

	JSON(w, http.StatusCreated, map[string]interface{}{
		"message":    "Solicitud de ingreso registrada con éxito. Será evaluada por un Director del semillero.",
		"request_id": reqID,
		"status":     "PENDING",
	})
}

// GetPendingRequests lista las solicitudes para el panel de Directores
func (h *AuthHandler) GetPendingRequests(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := `SELECT id, google_id, email, full_name, student_code, career_program, motivation_letter, status, created_at 
	          FROM registration_requests 
	          WHERE status = 'PENDING' 
	          ORDER BY created_at ASC`

	rows, err := h.db.Query(ctx, query)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar solicitudes", err.Error())
		return
	}
	defer rows.Close()

	requests := make([]domain.RegistrationRequest, 0)
	for rows.Next() {
		var req domain.RegistrationRequest
		if err := rows.Scan(&req.ID, &req.GoogleID, &req.Email, &req.FullName, &req.StudentCode, &req.CareerProgram, &req.MotivationLetter, &req.Status, &req.CreatedAt); err != nil {
			Error(w, http.StatusInternalServerError, "Error al leer datos de solicitud")
			return
		}
		requests = append(requests, req)
	}

	JSON(w, http.StatusOK, requests)
}

type ReviewRequestBody struct {
	Action   string          `json:"action"` // "APPROVE" o "REJECT"
	Role     domain.UserRole `json:"role,omitempty"` // Default: MIEMBRO
	Feedback *string         `json:"feedback,omitempty"`
}

// ReviewRequest permite a un Director aprobar o rechazar una solicitud
func (h *AuthHandler) ReviewRequest(w http.ResponseWriter, r *http.Request) {
	reqIDStr := chi.URLParam(r, "id")
	reqID, err := uuid.Parse(reqIDStr)
	if err != nil {
		Error(w, http.StatusBadRequest, "ID de solicitud no válido")
		return
	}

	var body ReviewRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		Error(w, http.StatusBadRequest, "Datos de revisión inválidos")
		return
	}

	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		Error(w, http.StatusUnauthorized, "No autorizado")
		return
	}

	ctx := r.Context()

	// Obtener datos de la solicitud
	var req domain.RegistrationRequest
	queryReq := `SELECT id, google_id, email, full_name, student_code, status 
	             FROM registration_requests WHERE id = $1`
	err = h.db.QueryRow(ctx, queryReq, reqID).Scan(&req.ID, &req.GoogleID, &req.Email, &req.FullName, &req.StudentCode, &req.Status)
	if err != nil {
		Error(w, http.StatusNotFound, "Solicitud no encontrada")
		return
	}

	if req.Status != domain.RegPending {
		Error(w, http.StatusBadRequest, "Esta solicitud ya fue resuelta previamente")
		return
	}

	now := time.Now()

	if body.Action == "APPROVE" {
		assignedRole := domain.RoleMiembro
		if body.Role == domain.RoleAdministrador || body.Role == domain.RoleDirector {
			assignedRole = body.Role
		}

		// Iniciar transacción: Crear usuario y actualizar solicitud
		tx, err := h.db.Begin(ctx)
		if err != nil {
			Error(w, http.StatusInternalServerError, "Error al iniciar transacción")
			return
		}
		defer tx.Rollback(ctx)

		// Insertar usuario
		insertUser := `INSERT INTO users (google_id, email, student_code, full_name, role, is_active)
		               VALUES ($1, $2, $3, $4, $5, TRUE)
		               ON CONFLICT (email) DO UPDATE SET is_active = TRUE, role = EXCLUDED.role
		               RETURNING id`
		var newUserID uuid.UUID
		err = tx.QueryRow(ctx, insertUser, req.GoogleID, req.Email, req.StudentCode, req.FullName, assignedRole).Scan(&newUserID)
		if err != nil {
			ErrorWithDetails(w, http.StatusInternalServerError, "Error al crear usuario aprobado", err.Error())
			return
		}

		// Actualizar solicitud
		updateReq := `UPDATE registration_requests 
		              SET status = 'APPROVED', reviewed_by = $1, reviewed_at = $2, director_feedback = $3
		              WHERE id = $4`
		_, err = tx.Exec(ctx, updateReq, claims.UserID, now, body.Feedback, reqID)
		if err != nil {
			ErrorWithDetails(w, http.StatusInternalServerError, "Error al actualizar estado de solicitud", err.Error())
			return
		}

		if err := tx.Commit(ctx); err != nil {
			Error(w, http.StatusInternalServerError, "Error al confirmar aprobación")
			return
		}

		JSON(w, http.StatusOK, map[string]interface{}{
			"message": "Solicitud aprobada con éxito. Usuario creado en el semillero.",
			"user_id": newUserID,
			"role":    assignedRole,
		})
		return
	}

	if body.Action == "REJECT" {
		updateReq := `UPDATE registration_requests 
		              SET status = 'REJECTED', reviewed_by = $1, reviewed_at = $2, director_feedback = $3
		              WHERE id = $4`
		_, err := h.db.Exec(ctx, updateReq, claims.UserID, now, body.Feedback, reqID)
		if err != nil {
			ErrorWithDetails(w, http.StatusInternalServerError, "Error al rechazar solicitud", err.Error())
			return
		}

		JSON(w, http.StatusOK, map[string]string{
			"message": "Solicitud rechazada con retroalimentación registrada.",
		})
		return
	}

	Error(w, http.StatusBadRequest, "Acción desconocida. Debe ser APPROVE o REJECT")
}

// GetCurrentUser devuelve el perfil del usuario autenticado actual
func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		Error(w, http.StatusUnauthorized, "No autorizado")
		return
	}

	var u domain.User
	query := `SELECT id, google_id, email, student_code, full_name, role, can_operate_3d, avatar_url, bio, is_active, created_at, updated_at 
	          FROM users WHERE id = $1 LIMIT 1`

	err := h.db.QueryRow(r.Context(), query, claims.UserID).Scan(
		&u.ID, &u.GoogleID, &u.Email, &u.StudentCode, &u.FullName, &u.Role,
		&u.CanOperate3D, &u.AvatarURL, &u.Bio, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		Error(w, http.StatusNotFound, "Usuario no encontrado")
		return
	}

	JSON(w, http.StatusOK, u)
}

func (h *AuthHandler) generateJWT(u *domain.User) (string, error) {
	claims := domain.JWTClaims{
		UserID:       u.ID,
		Email:        u.Email,
		FullName:     u.FullName,
		Role:         u.Role,
		CanOperate3D: u.CanOperate3D,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)), // 3 días
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "semard-control-center",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.cfg.JWTSecret))
}
