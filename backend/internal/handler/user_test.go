package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nosedimetuXD/Semard-Control-Center/backend/internal/domain"
)

func TestUserHandler_DatabaseNilSafety(t *testing.T) {
	h := NewUserHandler(nil)

	// ListUsers with nil DB returns 503 Service Unavailable
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	w := httptest.NewRecorder()
	h.ListUsers(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("Esperaba status 503 sin DB, obtuvo %d", w.Code)
	}

	// UpdateRole with invalid JSON
	req = httptest.NewRequest(http.MethodPatch, "/api/v1/directors/users/test/role", bytes.NewBufferString("invalid json"))
	w = httptest.NewRecorder()
	h.UpdateRole(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("Esperaba status 503 sin DB, obtuvo %d", w.Code)
	}
}

func TestUserRoleValidation(t *testing.T) {
	validRoles := []domain.UserRole{domain.RoleDirector, domain.RoleAdministrador, domain.RoleMiembro}
	for _, r := range validRoles {
		if r != "DIRECTOR" && r != "ADMINISTRADOR" && r != "MIEMBRO" {
			t.Errorf("Rol no esperado: %s", r)
		}
	}
}
