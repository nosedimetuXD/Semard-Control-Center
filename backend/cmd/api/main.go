package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nosedimetuXD/Semard-Control-Center/backend/internal/config"
	"github.com/nosedimetuXD/Semard-Control-Center/backend/internal/database"
	"github.com/nosedimetuXD/Semard-Control-Center/backend/internal/handler"
	"github.com/nosedimetuXD/Semard-Control-Center/backend/internal/middleware"
)

func main() {
	log.Println("=== Iniciando SEMARD Control Center API ===")

	// 1. Cargar configuración de entorno
	cfg := config.Load()

	// 2. Conectar a PostgreSQL
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var db *pgxpool.Pool
	var err error

	db, err = database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Printf("Aviso de Conexión DB: %v (el servidor continuará para chequeos de salud)", err)
	} else {
		defer db.Close()
		// Ejecutar todas las migraciones en orden
		migrationsDir := "migrations"
		if err := database.RunMigrations(context.Background(), db, migrationsDir); err != nil {
			log.Printf("Aviso en migración automática: %v", err)
		}
	}

	// 3. Crear router Chi y middlewares
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.StructuredLogger)
	r.Use(middleware.NewCORS(cfg.AllowedOrigins))

	// 4. Inicializar handlers
	healthHandler := handler.NewHealthHandler(db)
	authHandler := handler.NewAuthHandler(cfg, db)
	userHandler := handler.NewUserHandler(db)
	eventHandler := handler.NewEventHandler(db)
	hubHandler := handler.NewHubHandler(db)

	// 5. Definir Rutas
	// 5.1 Monitoreo y Salud
	r.Get("/healthz", healthHandler.Check)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handler.JSON(w, http.StatusOK, map[string]string{
			"app":     "SEMARD Control Center API",
			"version": "1.0.0",
			"status":  "operational",
			"docs":    "/healthz",
		})
	})

	// 5.2 API v1
	r.Route("/api/v1", func(r chi.Router) {
		// Módulo Hub Institucional Público
		r.Route("/hub", func(r chi.Router) {
			r.Get("/info", hubHandler.GetHubInfo)
			r.Get("/directors", hubHandler.GetDirectors)
			r.Get("/projects", hubHandler.GetPublicProjects)
		})

		// Rutas públicas y privadas de Autenticación
		r.Route("/auth", func(r chi.Router) {
			r.Get("/google/login", authHandler.GoogleLogin)
			r.Get("/google/callback", authHandler.GoogleCallback)
			r.Post("/register-request", authHandler.RegisterRequest)

			// Ruta protegida para el usuario logueado
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireAuth(cfg.JWTSecret))
				r.Get("/me", authHandler.GetCurrentUser)
				r.Put("/me", userHandler.UpdateProfile)
			})
		})

		// Módulo de Directorio de Miembros
		r.Route("/users", func(r chi.Router) {
			r.Use(middleware.RequireAuth(cfg.JWTSecret))
			r.Get("/", userHandler.ListUsers)
			r.Get("/{id}", userHandler.GetUser)
		})

		// Módulo de Eventos (Públicos e Internos)
		r.Route("/events", func(r chi.Router) {
			// Consulta y registro público
			r.Get("/public", eventHandler.ListPublicEvents)
			r.Get("/{id}", eventHandler.GetEvent)
			r.Post("/{id}/register", eventHandler.RegisterAttendee)

			// Consulta interna para miembros del semillero
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireAuth(cfg.JWTSecret))
				r.Get("/internal", eventHandler.ListInternalEvents)

				// Ver inscritos (Admin o Director)
				r.With(middleware.RequireAdminOrDirector()).Get("/{id}/attendees", eventHandler.ListAttendees)

				// Gestión exclusiva para Directores
				r.Group(func(r chi.Router) {
					r.Use(middleware.RequireDirector())
					r.Post("/", eventHandler.CreateEvent)
					r.Put("/{id}", eventHandler.UpdateEvent)
					r.Delete("/{id}", eventHandler.DeleteEvent)
				})
			})
		})

		// Rutas protegidas exclusivas para Directores
		r.Route("/directors", func(r chi.Router) {
			r.Use(middleware.RequireAuth(cfg.JWTSecret))
			r.Use(middleware.RequireDirector())

			// Gestión de solicitudes de registro
			r.Get("/registration-requests", authHandler.GetPendingRequests)
			r.Post("/registration-requests/{id}/review", authHandler.ReviewRequest)

			// Gestión de integrantes y roles
			r.Patch("/users/{id}/role", userHandler.UpdateRole)
			r.Patch("/users/{id}/permissions", userHandler.UpdatePermissions)
			r.Patch("/users/{id}/status", userHandler.UpdateStatus)
		})
	})

	// 6. Iniciar servidor HTTP con Graceful Shutdown
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Servidor API escuchando en el puerto :%s", cfg.Port)
		serverErrors <- server.ListenAndServe()
	}()

	// Canal para escuchar señales de apagado (SIGINT, SIGTERM)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error fatal al iniciar servidor: %v", err)
		}
	case sig := <-shutdown:
		log.Printf("Señal de apagado recibida (%v), cerrando servidor de forma segura...", sig)
		ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelShutdown()

		if err := server.Shutdown(ctxShutdown); err != nil {
			log.Printf("Error al forzar apagado del servidor: %v", err)
			_ = server.Close()
		}
	}

	log.Println("=== Servidor detenido correctamente ===")
}
