package middleware

import (
	"net/http"
	"strings"

	"github.com/go-chi/cors"
)

func NewCORS(allowedOrigins []string) func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			if origin == "" {
				return true
			}
			// Permitir cualquier subdominio o preview de Vercel
			if strings.HasSuffix(origin, ".vercel.app") || strings.Contains(origin, "vercel.app") {
				return true
			}
			// Permitir desarrollo local
			if strings.HasPrefix(origin, "http://localhost:") || strings.HasPrefix(origin, "http://127.0.0.1:") {
				return true
			}
			// Verificar orígenes permitidos en configuración
			for _, o := range allowedOrigins {
				if o == "*" || strings.EqualFold(o, origin) {
					return true
				}
			}
			return false
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
