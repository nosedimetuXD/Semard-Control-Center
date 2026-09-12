package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nosedimetuXD/Semard-Control-Center/backend/internal/domain"
)

func TestRequireRole(t *testing.T) {
	tests := []struct {
		name           string
		userRole       domain.UserRole
		allowedRoles   []domain.UserRole
		expectedStatus int
	}{
		{
			name:           "Director tiene acceso a ruta de Director",
			userRole:       domain.RoleDirector,
			allowedRoles:   []domain.UserRole{domain.RoleDirector},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Miembro denegado en ruta exclusiva de Director",
			userRole:       domain.RoleMiembro,
			allowedRoles:   []domain.UserRole{domain.RoleDirector},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Admin tiene acceso a ruta de Admin o Director",
			userRole:       domain.RoleAdministrador,
			allowedRoles:   []domain.UserRole{domain.RoleDirector, domain.RoleAdministrador},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := RequireRole(tc.allowedRoles...)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			claims := &domain.JWTClaims{
				UserID:   uuid.New(),
				Email:    "test@unicartagena.edu.co",
				FullName: "Test User",
				Role:     tc.userRole,
			}

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			ctx := context.WithValue(req.Context(), UserContextKey, claims)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("Se esperaba código %d, se obtuvo %d", tc.expectedStatus, rr.Code)
			}
		})
	}
}

func TestRequireAuth(t *testing.T) {
	secret := "test-secret"
	claims := domain.JWTClaims{
		UserID:       uuid.New(),
		Email:        "estudiante@unicartagena.edu.co",
		FullName:     "Estudiante Prueba",
		Role:         domain.RoleMiembro,
		CanOperate3D: false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Error al firmar token: %v", err)
	}

	handler := RequireAuth(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userClaims, ok := GetUserClaims(r)
		if !ok || userClaims.Email != "estudiante@unicartagena.edu.co" {
			t.Errorf("No se recuperaron los claims esperados del contexto")
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Se esperaba código 200, se obtuvo %d", rr.Code)
	}
}
