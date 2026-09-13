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
)

type InventoryHandler struct {
	db *pgxpool.Pool
}

func NewInventoryHandler(db *pgxpool.Pool) *InventoryHandler {
	return &InventoryHandler{db: db}
}

// ListItems lista el catálogo de herramientas y equipos del semillero
func (h *InventoryHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	category := strings.TrimSpace(r.URL.Query().Get("category"))
	search := strings.TrimSpace(r.URL.Query().Get("search"))
	availableOnly := r.URL.Query().Get("available_only") == "true"

	query := `SELECT id, code, name, category, description, total_stock, available_stock, location, created_at, updated_at
	          FROM inventory_items WHERE 1=1`
	var args []interface{}
	argIdx := 1

	if category != "" {
		query += fmt.Sprintf(" AND category ILIKE $%d", argIdx)
		args = append(args, category)
		argIdx++
	}

	if availableOnly {
		query += " AND available_stock > 0"
	}

	if search != "" {
		pattern := "%" + search + "%"
		query += fmt.Sprintf(" AND (name ILIKE $%d OR code ILIKE $%d OR description ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, pattern)
		argIdx++
	}

	query += " ORDER BY category ASC, name ASC"

	rows, err := h.db.Query(r.Context(), query, args...)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar inventario", err.Error())
		return
	}
	defer rows.Close()

	items := make([]domain.InventoryItem, 0)
	for rows.Next() {
		var item domain.InventoryItem
		if err := rows.Scan(
			&item.ID, &item.Code, &item.Name, &item.Category, &item.Description,
			&item.TotalStock, &item.AvailableStock, &item.Location, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			ErrorWithDetails(w, http.StatusInternalServerError, "Error al leer ítem de inventario", err.Error())
			return
		}
		items = append(items, item)
	}

	JSON(w, http.StatusOK, items)
}

// GetItem retorna el detalle de un ítem por ID
func (h *InventoryHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	idStr := chi.URLParam(r, "id")
	itemID, err := uuid.Parse(idStr)
	if err != nil {
		Error(w, http.StatusBadRequest, "ID de ítem inválido")
		return
	}

	var item domain.InventoryItem
	query := `SELECT id, code, name, category, description, total_stock, available_stock, location, created_at, updated_at
	          FROM inventory_items WHERE id = $1 LIMIT 1`

	err = h.db.QueryRow(r.Context(), query, itemID).Scan(
		&item.ID, &item.Code, &item.Name, &item.Category, &item.Description,
		&item.TotalStock, &item.AvailableStock, &item.Location, &item.CreatedAt, &item.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			Error(w, http.StatusNotFound, "Ítem no encontrado en inventario")
			return
		}
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar ítem", err.Error())
		return
	}

	JSON(w, http.StatusOK, item)
}

// CreateItem registra un nuevo equipo o herramienta (Admin o Director)
func (h *InventoryHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateInventoryItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Payload JSON inválido")
		return
	}

	if strings.TrimSpace(req.Code) == "" {
		Error(w, http.StatusBadRequest, "El código único del ítem es obligatorio")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		Error(w, http.StatusBadRequest, "El nombre del ítem es obligatorio")
		return
	}
	if strings.TrimSpace(req.Category) == "" {
		Error(w, http.StatusBadRequest, "La categoría es obligatoria")
		return
	}
	if req.TotalStock <= 0 {
		req.TotalStock = 1
	}

	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	query := `INSERT INTO inventory_items (code, name, category, description, total_stock, available_stock, location)
	          VALUES ($1, $2, $3, $4, $5, $5, $6)
	          RETURNING id, code, name, category, description, total_stock, available_stock, location, created_at, updated_at`

	var item domain.InventoryItem
	err := h.db.QueryRow(r.Context(), query, req.Code, req.Name, req.Category, req.Description, req.TotalStock, req.Location).Scan(
		&item.ID, &item.Code, &item.Name, &item.Category, &item.Description,
		&item.TotalStock, &item.AvailableStock, &item.Location, &item.CreatedAt, &item.UpdatedAt,
	)

	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al registrar ítem de inventario", err.Error())
		return
	}

	JSON(w, http.StatusCreated, item)
}

// UpdateItem modifica los datos o ajusta el stock de un ítem (Admin o Director)
func (h *InventoryHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	itemID, err := uuid.Parse(idStr)
	if err != nil {
		Error(w, http.StatusBadRequest, "ID de ítem inválido")
		return
	}

	var req domain.UpdateInventoryItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Payload JSON inválido")
		return
	}

	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	query := `UPDATE inventory_items SET
	            name = COALESCE($1, name),
	            category = COALESCE($2, category),
	            description = COALESCE($3, description),
	            total_stock = COALESCE($4, total_stock),
	            available_stock = COALESCE($5, available_stock),
	            location = COALESCE($6, location),
	            updated_at = NOW()
	          WHERE id = $7
	          RETURNING id, code, name, category, description, total_stock, available_stock, location, created_at, updated_at`

	var item domain.InventoryItem
	err = h.db.QueryRow(r.Context(), query, req.Name, req.Category, req.Description, req.TotalStock, req.AvailableStock, req.Location, itemID).Scan(
		&item.ID, &item.Code, &item.Name, &item.Category, &item.Description,
		&item.TotalStock, &item.AvailableStock, &item.Location, &item.CreatedAt, &item.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			Error(w, http.StatusNotFound, "Ítem no encontrado")
			return
		}
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al actualizar ítem", err.Error())
		return
	}

	JSON(w, http.StatusOK, item)
}

// DeleteItem elimina un ítem si no tiene préstamos activos (Admin o Director)
func (h *InventoryHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	idStr := chi.URLParam(r, "id")
	itemID, err := uuid.Parse(idStr)
	if err != nil {
		Error(w, http.StatusBadRequest, "ID de ítem inválido")
		return
	}

	// Verificar si hay préstamos activos
	var activeLoans int
	err = h.db.QueryRow(r.Context(), "SELECT COUNT(*) FROM loan_requests WHERE item_id = $1 AND status IN ('PENDING', 'APPROVED', 'APPROVED_MODIFIED')", itemID).Scan(&activeLoans)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al verificar préstamos del ítem", err.Error())
		return
	}
	if activeLoans > 0 {
		Error(w, http.StatusBadRequest, "No es posible eliminar el ítem porque tiene préstamos activos o pendientes")
		return
	}

	result, err := h.db.Exec(r.Context(), "DELETE FROM inventory_items WHERE id = $1", itemID)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al eliminar ítem", err.Error())
		return
	}

	if result.RowsAffected() == 0 {
		Error(w, http.StatusNotFound, "Ítem no encontrado")
		return
	}

	JSON(w, http.StatusOK, map[string]string{
		"message": "Ítem de inventario eliminado correctamente",
	})
}
