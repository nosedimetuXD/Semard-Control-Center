package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
	projectHandler := handler.NewProjectHandler(db)
	updateHandler := handler.NewUpdateHandler(db, projectHandler)
	resourceHandler := handler.NewResourceHandler(db, projectHandler)
	inventoryHandler := handler.NewInventoryHandler(db)
	loanHandler := handler.NewLoanHandler(db)
	print3DHandler := handler.NewPrint3DHandler(db, cfg.StoragePath)

	// 5. Definir Rutas
	// 5.1 Monitoreo, Test Studio y Salud
	r.Get("/healthz", healthHandler.Check)
	r.Get("/test", handler.ServeTestUI)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept"), "text/html") {
			handler.ServeTestUI(w, r)
			return
		}
		handler.JSON(w, http.StatusOK, map[string]string{
			"app":     "SEMARD Control Center API",
			"version": "1.0.0",
			"status":  "operational",
			"test_ui": "/test",
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

		// Módulo Logístico de Proyectos, Avances y Recursos
		r.Route("/projects", func(r chi.Router) {
			r.Use(middleware.RequireAuth(cfg.JWTSecret))

			// Vistas de proyectos
			r.Get("/showcase", projectHandler.ListShowcase)
			r.Get("/my-projects", projectHandler.ListMyProjects)
			r.Get("/{id}", projectHandler.GetProject)

			// Avances de proyecto (Encargados o Directores)
			r.Get("/{id}/updates", updateHandler.ListUpdates)
			r.Post("/{id}/updates", updateHandler.CreateUpdate)
			r.With(middleware.RequireDirector()).Post("/updates/{updateId}/review", updateHandler.ReviewUpdate)

			// Solicitudes de recursos de proyecto (Encargados o Directores)
			r.Get("/{id}/resources", resourceHandler.ListProjectResources)
			r.Post("/{id}/resources", resourceHandler.CreateResourceRequest)
			r.With(middleware.RequireDirector()).Post("/resources/{resourceId}/review", resourceHandler.ReviewResource)

			// Gestión de proyectos y asignación de encargados (Exclusivo Directores)
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireDirector())
				r.Post("/", projectHandler.CreateProject)
				r.Put("/{id}", projectHandler.UpdateProject)
				r.Post("/{id}/members", projectHandler.AssignMember)
				r.Delete("/{id}/members/{userId}", projectHandler.RemoveMember)
			})
		})

		// Módulo de Inventario y Herramientas
		r.Route("/inventory", func(r chi.Router) {
			r.Use(middleware.RequireAuth(cfg.JWTSecret))

			// Consulta de catálogo
			r.Get("/", inventoryHandler.ListItems)
			r.Get("/{id}", inventoryHandler.GetItem)

			// Gestión exclusiva para Admins y Directores
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireAdminOrDirector())
				r.Post("/", inventoryHandler.CreateItem)
				r.Put("/{id}", inventoryHandler.UpdateItem)
				r.Delete("/{id}", inventoryHandler.DeleteItem)
			})
		})

		// Módulo de Préstamos de Equipos
		r.Route("/loans", func(r chi.Router) {
			r.Use(middleware.RequireAuth(cfg.JWTSecret))

			// Solicitudes del miembro
			r.Post("/requests", loanHandler.CreateRequest)
			r.Get("/my-loans", loanHandler.ListMyLoans)

			// Gestión exclusiva para Admins y Directores
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireAdminOrDirector())
				r.Get("/", loanHandler.ListAllLoans)
				r.Post("/{id}/review", loanHandler.ReviewLoan)
				r.Post("/{id}/return", loanHandler.ReturnLoan)
			})
		})

		// Módulo de Taller de Impresión 3D
		r.Route("/print3d", func(r chi.Router) {
			r.Use(middleware.RequireAuth(cfg.JWTSecret))

			// Solicitud y consulta de cola
			r.Post("/requests", print3DHandler.CreateRequest)
			r.Get("/my-requests", print3DHandler.ListMyRequests)
			r.Get("/queue", print3DHandler.ListQueue)
			r.Get("/requests/{id}/download", print3DHandler.DownloadFile)

			// Evaluación exclusiva para Admins y Directores
			r.With(middleware.RequireAdminOrDirector()).Post("/requests/{id}/review", print3DHandler.ReviewRequest)

			// Ejecución técnica de la cola (Operador 3D, Admin o Director)
			r.With(middleware.Require3DOperator()).Post("/requests/{id}/status", print3DHandler.UpdateStatus)
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

			// Bandeja global de recursos pendientes de aprobación
			r.Get("/resources/pending", resourceHandler.ListPendingResources)
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
