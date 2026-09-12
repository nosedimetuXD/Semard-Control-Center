package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nosedimetuXD/Semard-Control-Center/backend/internal/domain"
	"github.com/nosedimetuXD/Semard-Control-Center/backend/internal/middleware"
)

func TestProjectNilSafety(t *testing.T) {
	ph := NewProjectHandler(nil)

	// ListShowcase with nil DB
	req := httptest.NewRequest(http.MethodGet, "/api/v1/projects/showcase", nil)
	w := httptest.NewRecorder()
	ph.ListShowcase(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("Esperaba status 503 sin DB, obtuvo %d", w.Code)
	}
}

func TestUpdateReview_MandatoryFeedback(t *testing.T) {
	ph := NewProjectHandler(nil)
	uh := NewUpdateHandler(nil, ph)

	// Contexto con director
	directorID := uuid.New()
	claims := &domain.JWTClaims{
		UserID: directorID,
		Role:   domain.RoleDirector,
	}
	ctx := context.WithValue(context.Background(), middleware.UserContextKey, claims)

	// Ruta con chi URLParam updateId
	rCtx := chi.NewRouteContext()
	rCtx.URLParams.Add("updateId", uuid.New().String())
	ctx = context.WithValue(ctx, chi.RouteCtxKey, rCtx)

	// 1. Envío con feedback vacío (debe fallar con 400 Bad Request)
	payloadWithoutFeedback := `{"status": "APPROVED", "director_feedback": "   "}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/projects/updates/test/review", bytes.NewBufferString(payloadWithoutFeedback)).WithContext(ctx)
	w := httptest.NewRecorder()

	uh.ReviewUpdate(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Esperaba status 400 por feedback vacío, obtuvo %d", w.Code)
	}

	// 2. Envío con status inválido
	payloadInvalidStatus := `{"status": "UNKNOWN_STATUS", "director_feedback": "Buen trabajo"}`
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/projects/updates/test/review", bytes.NewBufferString(payloadInvalidStatus)).WithContext(ctx)
	w2 := httptest.NewRecorder()

	uh.ReviewUpdate(w2, req2)

	if w2.Code != http.StatusBadRequest {
		t.Fatalf("Esperaba status 400 por status inválido, obtuvo %d", w2.Code)
	}
}

func TestResourceReview_3StatesAndMandatoryFeedback(t *testing.T) {
	ph := NewProjectHandler(nil)
	rh := NewResourceHandler(nil, ph)

	directorID := uuid.New()
	claims := &domain.JWTClaims{
		UserID: directorID,
		Role:   domain.RoleDirector,
	}
	ctx := context.WithValue(context.Background(), middleware.UserContextKey, claims)

	rCtx := chi.NewRouteContext()
	rCtx.URLParams.Add("resourceId", uuid.New().String())
	ctx = context.WithValue(ctx, chi.RouteCtxKey, rCtx)

	// 1. Validar que rechaza feedback vacío
	payloadWithoutFeedback := `{"status": "RETURNED_FOR_MODIFICATION", "director_feedback": ""}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/projects/resources/test/review", bytes.NewBufferString(payloadWithoutFeedback)).WithContext(ctx)
	w := httptest.NewRecorder()

	rh.ReviewResource(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Esperaba status 400 por feedback vacío, obtuvo %d", w.Code)
	}

	// 2. Validar que rechaza estado inválido
	payloadInvalidStatus := `{"status": "PENDING", "director_feedback": "Faltan cotizaciones"}`
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/projects/resources/test/review", bytes.NewBufferString(payloadInvalidStatus)).WithContext(ctx)
	w2 := httptest.NewRecorder()

	rh.ReviewResource(w2, req2)

	if w2.Code != http.StatusBadRequest {
		t.Fatalf("Esperaba status 400 al intentar asignar status no resolutivo, obtuvo %d", w2.Code)
	}
}
