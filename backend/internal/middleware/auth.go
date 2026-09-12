package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nosedimetuXD/Semard-Control-Center/backend/internal/domain"
)

type contextKey string

const UserContextKey contextKey = "currentUser"

// RequireAuth extrae y valida el token JWT del encabezado Authorization
func RequireAuth(jwtSecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"Encabezado Authorization requerido"}`, http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, `{"error":"Formato de token inválido. Debe ser: Bearer <token>"}`, http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]
			claims := &domain.JWTClaims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("método de firma no válido")
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, `{"error":"Token inválido o expirado"}`, http.StatusUnauthorized)
				return
			}

			// Inyectar claims en el contexto
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole verifica que el usuario autenticado posea uno de los roles autorizados
func RequireRole(allowedRoles ...domain.UserRole) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserContextKey).(*domain.JWTClaims)
			if !ok || claims == nil {
				http.Error(w, `{"error":"No autenticado"}`, http.StatusUnauthorized)
				return
			}

			for _, role := range allowedRoles {
				if claims.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, `{"error":"Acceso denegado: permisos insuficientes para esta operación"}`, http.StatusForbidden)
		})
	}
}

// RequireDirector asegura que solo usuarios con rol DIRECTOR accedan
func RequireDirector() func(next http.Handler) http.Handler {
	return RequireRole(domain.RoleDirector)
}

// RequireAdminOrDirector asegura que solo DIRECTORES o ADMINISTRADORES accedan
func RequireAdminOrDirector() func(next http.Handler) http.Handler {
	return RequireRole(domain.RoleDirector, domain.RoleAdministrador)
}

// Require3DOperator verifica si el usuario es Director, Administrador o Miembro con permiso Operador 3D
func Require3DOperator() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserContextKey).(*domain.JWTClaims)
			if !ok || claims == nil {
				http.Error(w, `{"error":"No autenticado"}`, http.StatusUnauthorized)
				return
			}

			if claims.Role == domain.RoleDirector || claims.Role == domain.RoleAdministrador || claims.CanOperate3D {
				next.ServeHTTP(w, r)
				return
			}

			http.Error(w, `{"error":"Acceso denegado: se requiere permiso de Operador 3D, Administrador o Director"}`, http.StatusForbidden)
		})
	}
}

// GetUserClaims helper para obtener los datos del usuario actual desde handlers
func GetUserClaims(r *http.Request) (*domain.JWTClaims, bool) {
	claims, ok := r.Context().Value(UserContextKey).(*domain.JWTClaims)
	return claims, ok
}
