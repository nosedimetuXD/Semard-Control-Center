package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nosedimetuXD/Semard-Control-Center/backend/internal/domain"
	"github.com/nosedimetuXD/Semard-Control-Center/backend/internal/middleware"
)

func TestLoanReviewValidation(t *testing.T) {
	lh := NewLoanHandler(nil)

	directorID := uuid.New()
	claims := &domain.JWTClaims{
		UserID: directorID,
		Role:   domain.RoleDirector,
	}
	ctx := context.WithValue(context.Background(), middleware.UserContextKey, claims)

	rCtx := chi.NewRouteContext()
	rCtx.URLParams.Add("id", uuid.New().String())
	ctx = context.WithValue(ctx, chi.RouteCtxKey, rCtx)

	// 1. APPROVED_MODIFIED sin approved_end_date debe fallar con 400
	payloadNoDate := `{"status": "APPROVED_MODIFIED", "reviewer_feedback": "Ajuste de plazo"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/loans/test/review", bytes.NewBufferString(payloadNoDate)).WithContext(ctx)
	w := httptest.NewRecorder()

	lh.ReviewLoan(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("Esperaba status 400 por falta de approved_end_date, obtuvo %d", w.Code)
	}

	// 2. REJECTED sin reviewer_feedback debe fallar con 400
	payloadNoFeedback := `{"status": "REJECTED"}`
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/loans/test/review", bytes.NewBufferString(payloadNoFeedback)).WithContext(ctx)
	w2 := httptest.NewRecorder()

	lh.ReviewLoan(w2, req2)
	if w2.Code != http.StatusBadRequest {
		t.Fatalf("Esperaba status 400 por falta de feedback en rechazo, obtuvo %d", w2.Code)
	}

	// 3. Status inválido
	payloadBadStatus := `{"status": "INVALID_STATUS"}`
	req3 := httptest.NewRequest(http.MethodPost, "/api/v1/loans/test/review", bytes.NewBufferString(payloadBadStatus)).WithContext(ctx)
	w3 := httptest.NewRecorder()

	lh.ReviewLoan(w3, req3)
	if w3.Code != http.StatusBadRequest {
		t.Fatalf("Esperaba status 400 por status inválido, obtuvo %d", w3.Code)
	}
}

func TestLoanRequestDateValidation(t *testing.T) {
	lh := NewLoanHandler(nil)

	userID := uuid.New()
	claims := &domain.JWTClaims{
		UserID: userID,
		Role:   domain.RoleMiembro,
	}
	ctx := context.WithValue(context.Background(), middleware.UserContextKey, claims)

	// Fecha fin anterior a fecha inicio
	start := time.Now()
	end := start.Add(-24 * time.Hour)

	payload := `{"item_id": "` + uuid.New().String() + `", "reason": "Proyecto de robótica", "requested_start_date": "` + start.Format(time.RFC3339) + `", "requested_end_date": "` + end.Format(time.RFC3339) + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/loans/requests", bytes.NewBufferString(payload)).WithContext(ctx)
	w := httptest.NewRecorder()

	lh.CreateRequest(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("Esperaba status 400 por fechas inconsistentes, obtuvo %d", w.Code)
	}
}
