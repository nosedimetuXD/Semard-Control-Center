package handler

import (
	"encoding/json"
	"errors"
	"fmt"
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

type LoanHandler struct {
	db *pgxpool.Pool
}

func NewLoanHandler(db *pgxpool.Pool) *LoanHandler {
	return &LoanHandler{db: db}
}

// CreateRequest permite a un miembro solicitar el préstamo de un ítem de inventario
func (h *LoanHandler) CreateRequest(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		Error(w, http.StatusUnauthorized, "Sesión no válida")
		return
	}

	var req domain.CreateLoanRequestInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Payload JSON inválido")
		return
	}

	if req.ItemID == uuid.Nil {
		Error(w, http.StatusBadRequest, "item_id es obligatorio")
		return
	}
	if strings.TrimSpace(req.Reason) == "" {
		Error(w, http.StatusBadRequest, "El motivo del préstamo es obligatorio")
		return
	}
	if req.RequestedStartDate.IsZero() {
		req.RequestedStartDate = time.Now()
	}
	if req.RequestedEndDate.Before(req.RequestedStartDate) {
		Error(w, http.StatusBadRequest, "La fecha de finalización no puede ser anterior a la de inicio")
		return
	}

	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	// Verificar stock disponible del ítem
	var availableStock int
	err = h.db.QueryRow(r.Context(), "SELECT available_stock FROM inventory_items WHERE id = $1", req.ItemID).Scan(&availableStock)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			Error(w, http.StatusNotFound, "Ítem no encontrado en inventario")
			return
		}
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al verificar stock", err.Error())
		return
	}

	if availableStock <= 0 {
		Error(w, http.StatusBadRequest, "El ítem no cuenta con stock disponible para préstamo en este momento")
		return
	}

	query := `INSERT INTO loan_requests (item_id, requested_by, reason, requested_start_date, requested_end_date, status)
	          VALUES ($1, $2, $3, $4, $5, 'PENDING')
	          RETURNING id, item_id, requested_by, reason, requested_start_date, requested_end_date, approved_end_date, status, reviewer_feedback, reviewed_by, returned_at, created_at, updated_at`

	var loan domain.LoanRequest
	err = h.db.QueryRow(r.Context(), query, req.ItemID, userID, req.Reason, req.RequestedStartDate, req.RequestedEndDate).Scan(
		&loan.ID, &loan.ItemID, &loan.RequestedBy, &loan.Reason, &loan.RequestedStartDate, &loan.RequestedEndDate,
		&loan.ApprovedEndDate, &loan.Status, &loan.ReviewerFeedback, &loan.ReviewedBy, &loan.ReturnedAt, &loan.CreatedAt, &loan.UpdatedAt,
	)

	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al registrar solicitud de préstamo", err.Error())
		return
	}

	JSON(w, http.StatusCreated, loan)
}

// ListMyLoans lista los préstamos solicitados por el usuario autenticado
func (h *LoanHandler) ListMyLoans(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		Error(w, http.StatusUnauthorized, "Sesión no válida")
		return
	}

	query := `SELECT lr.id, lr.item_id, lr.requested_by, lr.reason, lr.requested_start_date, lr.requested_end_date,
	                 lr.approved_end_date, lr.status, lr.reviewer_feedback, lr.reviewed_by, lr.returned_at, lr.created_at, lr.updated_at,
	                 ii.name AS item_name, ii.code AS item_code
	          FROM loan_requests lr
	          JOIN inventory_items ii ON lr.item_id = ii.id
	          WHERE lr.requested_by = $1
	          ORDER BY lr.created_at DESC`

	rows, err := h.db.Query(r.Context(), query, userID)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar mis préstamos", err.Error())
		return
	}
	defer rows.Close()

	type MyLoanItem struct {
		domain.LoanRequest
		ItemName string `json:"item_name"`
		ItemCode string `json:"item_code"`
	}

	loans := make([]MyLoanItem, 0)
	for rows.Next() {
		var item MyLoanItem
		if err := rows.Scan(
			&item.ID, &item.ItemID, &item.RequestedBy, &item.Reason, &item.RequestedStartDate, &item.RequestedEndDate,
			&item.ApprovedEndDate, &item.Status, &item.ReviewerFeedback, &item.ReviewedBy, &item.ReturnedAt,
			&item.CreatedAt, &item.UpdatedAt, &item.ItemName, &item.ItemCode,
		); err != nil {
			ErrorWithDetails(w, http.StatusInternalServerError, "Error al leer datos de préstamo", err.Error())
			return
		}
		loans = append(loans, item)
	}

	JSON(w, http.StatusOK, loans)
}

// ListAllLoans lista todos los préstamos para Directores y Administradores
func (h *LoanHandler) ListAllLoans(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	statusFilter := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("status")))

	query := `SELECT lr.id, lr.item_id, lr.requested_by, lr.reason, lr.requested_start_date, lr.requested_end_date,
	                 lr.approved_end_date, lr.status, lr.reviewer_feedback, lr.reviewed_by, lr.returned_at, lr.created_at, lr.updated_at,
	                 ii.name AS item_name, ii.code AS item_code, u.full_name AS requester_name, u.email AS requester_email
	          FROM loan_requests lr
	          JOIN inventory_items ii ON lr.item_id = ii.id
	          JOIN users u ON lr.requested_by = u.id
	          WHERE 1=1`
	var args []interface{}

	if statusFilter != "" {
		query += " AND lr.status = $1"
		args = append(args, statusFilter)
	}
	query += " ORDER BY lr.created_at DESC"

	rows, err := h.db.Query(r.Context(), query, args...)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar préstamos", err.Error())
		return
	}
	defer rows.Close()

	type AdminLoanItem struct {
		domain.LoanRequest
		ItemName       string `json:"item_name"`
		ItemCode       string `json:"item_code"`
		RequesterName  string `json:"requester_name"`
		RequesterEmail string `json:"requester_email"`
	}

	loans := make([]AdminLoanItem, 0)
	for rows.Next() {
		var item AdminLoanItem
		if err := rows.Scan(
			&item.ID, &item.ItemID, &item.RequestedBy, &item.Reason, &item.RequestedStartDate, &item.RequestedEndDate,
			&item.ApprovedEndDate, &item.Status, &item.ReviewerFeedback, &item.ReviewedBy, &item.ReturnedAt,
			&item.CreatedAt, &item.UpdatedAt, &item.ItemName, &item.ItemCode, &item.RequesterName, &item.RequesterEmail,
		); err != nil {
			ErrorWithDetails(w, http.StatusInternalServerError, "Error al leer datos de préstamo", err.Error())
			return
		}
		loans = append(loans, item)
	}

	JSON(w, http.StatusOK, loans)
}

// ReviewLoan permite a un Director o Admin aprobar, aprobar modificando plazo o rechazar con feedback
func (h *LoanHandler) ReviewLoan(w http.ResponseWriter, r *http.Request) {
	reviewerID, err := middleware.GetUserID(r.Context())
	if err != nil {
		Error(w, http.StatusUnauthorized, "Sesión no válida")
		return
	}

	loanIDStr := chi.URLParam(r, "id")
	loanID, err := uuid.Parse(loanIDStr)
	if err != nil {
		Error(w, http.StatusBadRequest, "ID de solicitud de préstamo inválido")
		return
	}

	var req domain.ReviewLoanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Payload JSON inválido")
		return
	}

	statusUpper := domain.LoanStatus(strings.ToUpper(string(req.Status)))
	if statusUpper != domain.LoanApproved &&
		statusUpper != domain.LoanApprovedModified &&
		statusUpper != domain.LoanRejected {
		Error(w, http.StatusBadRequest, "Estado inválido. Debe ser: APPROVED, APPROVED_MODIFIED o REJECTED")
		return
	}

	// Validar plazo modificado si es APPROVED_MODIFIED
	if statusUpper == domain.LoanApprovedModified {
		if req.ApprovedEndDate == nil || req.ApprovedEndDate.IsZero() {
			Error(w, http.StatusBadRequest, "approved_end_date es obligatorio al aprobar con modificación de plazo")
			return
		}
	}

	// Validar feedback obligatorio si es REJECTED o APPROVED_MODIFIED
	if statusUpper == domain.LoanRejected || statusUpper == domain.LoanApprovedModified {
		if req.ReviewerFeedback == nil || strings.TrimSpace(*req.ReviewerFeedback) == "" {
			Error(w, http.StatusBadRequest, "El comentario o retroalimentación del revisor es obligatorio")
			return
		}
	}

	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al iniciar transacción", err.Error())
		return
	}
	defer tx.Rollback(r.Context())

	// Obtener el préstamo actual
	var currentStatus domain.LoanStatus
	var itemID uuid.UUID
	var requestedEndDate time.Time
	err = tx.QueryRow(r.Context(), "SELECT status, item_id, requested_end_date FROM loan_requests WHERE id = $1 FOR UPDATE", loanID).Scan(&currentStatus, &itemID, &requestedEndDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			Error(w, http.StatusNotFound, "Solicitud de préstamo no encontrada")
			return
		}
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar préstamo", err.Error())
		return
	}

	if currentStatus != domain.LoanPending {
		Error(w, http.StatusBadRequest, fmt.Sprintf("La solicitud ya no está pendiente (estado actual: %s)", currentStatus))
		return
	}

	finalApprovedEndDate := req.ApprovedEndDate
	if statusUpper == domain.LoanApproved {
		finalApprovedEndDate = &requestedEndDate
	}

	// Si se aprueba, descontar stock disponible
	if statusUpper == domain.LoanApproved || statusUpper == domain.LoanApprovedModified {
		var availableStock int
		err = tx.QueryRow(r.Context(), "SELECT available_stock FROM inventory_items WHERE id = $1 FOR UPDATE", itemID).Scan(&availableStock)
		if err != nil {
			ErrorWithDetails(w, http.StatusInternalServerError, "Error al verificar stock de ítem", err.Error())
			return
		}

		if availableStock <= 0 {
			Error(w, http.StatusBadRequest, "No hay stock disponible para aprobar este préstamo")
			return
		}

		_, err = tx.Exec(r.Context(), "UPDATE inventory_items SET available_stock = available_stock - 1, updated_at = NOW() WHERE id = $1", itemID)
		if err != nil {
			ErrorWithDetails(w, http.StatusInternalServerError, "Error al descontar stock", err.Error())
			return
		}
	}

	// Actualizar préstamo
	updateQuery := `UPDATE loan_requests SET
	                  status = $1,
	                  approved_end_date = $2,
	                  reviewer_feedback = $3,
	                  reviewed_by = $4,
	                  updated_at = NOW()
	                WHERE id = $5
	                RETURNING id, item_id, requested_by, reason, requested_start_date, requested_end_date, approved_end_date, status, reviewer_feedback, reviewed_by, returned_at, created_at, updated_at`

	var loan domain.LoanRequest
	err = tx.QueryRow(r.Context(), updateQuery, statusUpper, finalApprovedEndDate, req.ReviewerFeedback, reviewerID, loanID).Scan(
		&loan.ID, &loan.ItemID, &loan.RequestedBy, &loan.Reason, &loan.RequestedStartDate, &loan.RequestedEndDate,
		&loan.ApprovedEndDate, &loan.Status, &loan.ReviewerFeedback, &loan.ReviewedBy, &loan.ReturnedAt, &loan.CreatedAt, &loan.UpdatedAt,
	)

	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al actualizar solicitud de préstamo", err.Error())
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al confirmar transacción", err.Error())
		return
	}

	JSON(w, http.StatusOK, loan)
}

// ReturnLoan registra la devolución física del equipo y restaura el stock disponible
func (h *LoanHandler) ReturnLoan(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	loanIDStr := chi.URLParam(r, "id")
	loanID, err := uuid.Parse(loanIDStr)
	if err != nil {
		Error(w, http.StatusBadRequest, "ID de préstamo inválido")
		return
	}

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al iniciar transacción", err.Error())
		return
	}
	defer tx.Rollback(r.Context())

	var currentStatus domain.LoanStatus
	var itemID uuid.UUID
	err = tx.QueryRow(r.Context(), "SELECT status, item_id FROM loan_requests WHERE id = $1 FOR UPDATE", loanID).Scan(&currentStatus, &itemID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			Error(w, http.StatusNotFound, "Préstamo no encontrado")
			return
		}
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar préstamo", err.Error())
		return
	}

	if currentStatus != domain.LoanApproved && currentStatus != domain.LoanApprovedModified && currentStatus != domain.LoanOverdue {
		Error(w, http.StatusBadRequest, fmt.Sprintf("No se puede registrar devolución para un préstamo con estado: %s", currentStatus))
		return
	}

	// Restaurar stock disponible (sin superar total_stock)
	_, err = tx.Exec(r.Context(), `UPDATE inventory_items 
	                              SET available_stock = LEAST(total_stock, available_stock + 1), updated_at = NOW() 
	                              WHERE id = $1`, itemID)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al restaurar stock de inventario", err.Error())
		return
	}

	// Marcar como devuelto
	query := `UPDATE loan_requests SET 
	            status = 'RETURNED', 
	            returned_at = NOW(), 
	            updated_at = NOW() 
	          WHERE id = $1
	          RETURNING id, item_id, requested_by, reason, requested_start_date, requested_end_date, approved_end_date, status, reviewer_feedback, reviewed_by, returned_at, created_at, updated_at`

	var loan domain.LoanRequest
	err = tx.QueryRow(r.Context(), query, loanID).Scan(
		&loan.ID, &loan.ItemID, &loan.RequestedBy, &loan.Reason, &loan.RequestedStartDate, &loan.RequestedEndDate,
		&loan.ApprovedEndDate, &loan.Status, &loan.ReviewerFeedback, &loan.ReviewedBy, &loan.ReturnedAt, &loan.CreatedAt, &loan.UpdatedAt,
	)

	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al marcar devolución del préstamo", err.Error())
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al confirmar devolución", err.Error())
		return
	}

	JSON(w, http.StatusOK, map[string]interface{}{
		"message": "Devolución registrada exitosamente y stock de equipo restaurado",
		"loan":    loan,
	})
}
