package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nosedimetuXD/Semard-Control-Center/backend/internal/domain"
	"github.com/nosedimetuXD/Semard-Control-Center/backend/internal/middleware"
)

func TestPrint3DReviewValidation(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "semard_3d_test")
	if err != nil {
		t.Fatalf("Error al crear temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	ph := NewPrint3DHandler(nil, tempDir)

	directorID := uuid.New()
	claims := &domain.JWTClaims{
		UserID: directorID,
		Role:   domain.RoleDirector,
	}
	ctx := context.WithValue(context.Background(), middleware.UserContextKey, claims)

	rCtx := chi.NewRouteContext()
	rCtx.URLParams.Add("id", uuid.New().String())
	ctx = context.WithValue(ctx, chi.RouteCtxKey, rCtx)

	// 1. REJECTED sin feedback
	payloadNoFeedback := `{"status": "REJECTED"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/print3d/requests/test/review", bytes.NewBufferString(payloadNoFeedback)).WithContext(ctx)
	w := httptest.NewRecorder()

	ph.ReviewRequest(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("Esperaba status 400 por falta de motivo de rechazo, obtuvo %d", w.Code)
	}

	// 2. Status inválido
	payloadBadStatus := `{"status": "IN_PROGRESS"}`
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/print3d/requests/test/review", bytes.NewBufferString(payloadBadStatus)).WithContext(ctx)
	w2 := httptest.NewRecorder()

	ph.ReviewRequest(w2, req2)
	if w2.Code != http.StatusBadRequest {
		t.Fatalf("Esperaba status 400 al pasar status operativo en review, obtuvo %d", w2.Code)
	}
}

func TestPrint3DStatusValidation(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "semard_3d_test2")
	if err != nil {
		t.Fatalf("Error al crear temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	ph := NewPrint3DHandler(nil, tempDir)

	operatorID := uuid.New()
	claims := &domain.JWTClaims{
		UserID:       operatorID,
		Role:         domain.RoleMiembro,
		CanOperate3D: true,
	}
	ctx := context.WithValue(context.Background(), middleware.UserContextKey, claims)

	rCtx := chi.NewRouteContext()
	rCtx.URLParams.Add("id", uuid.New().String())
	ctx = context.WithValue(ctx, chi.RouteCtxKey, rCtx)

	// Status operativo no válido
	payloadBadStatus := `{"status": "APPROVED"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/print3d/requests/test/status", bytes.NewBufferString(payloadBadStatus)).WithContext(ctx)
	w := httptest.NewRecorder()

	ph.UpdateStatus(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("Esperaba status 400 por status no operativo, obtuvo %d", w.Code)
	}
}
