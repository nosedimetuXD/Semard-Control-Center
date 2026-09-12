package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func StructuredLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()

		defer func() {
			duration := time.Since(start)
			log.Printf(
				"[%s] %s | Status: %d | Latencia: %v | IP: %s",
				r.Method,
				r.URL.Path,
				ww.Status(),
				duration,
				r.RemoteAddr,
			)
		}()

		next.ServeHTTP(ww, r)
	})
}
