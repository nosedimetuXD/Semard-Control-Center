package handler

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nosedimetuXD/Semard-Control-Center/backend/internal/domain"
)

type HubHandler struct {
	db *pgxpool.Pool
}

func NewHubHandler(db *pgxpool.Pool) *HubHandler {
	return &HubHandler{db: db}
}

// GetHubInfo retorna información institucional del semillero y sus líneas de investigación
func (h *HubHandler) GetHubInfo(w http.ResponseWriter, r *http.Request) {
	info := map[string]interface{}{
		"name":         "SEMARD — Semillero de Investigación en Robótica y Sistemas Autónomos",
		"acronym":      "SEMARD",
		"institution":  "Universidad de Cartagena",
		"faculty":      "Facultad de Ingeniería",
		"program":      "Ingeniería de Sistemas / Ingeniería Mecatrónica",
		"lab_location": "Laboratorio de Robótica y Control, Campus Piedra de Bolívar, Cartagena de Indias",
		"contact": map[string]string{
			"email":     "semard@unicartagena.edu.co",
			"instagram": "@semard_unicartagena",
		},
		"mission": "Promover la investigación formativa, desarrollo tecnológico y divulgación científica en estudiantes de pregrado, mediante el diseño e implementación de proyectos de robótica aplicada, manufactura aditiva y sistemas inteligentes.",
		"vision":  "Ser reconocidos como el semillero de investigación líder en el Caribe colombiano en robótica y automatización, con participación destacada en competencias nacionales y proyectos con impacto social.",
		"research_lines": []map[string]string{
			{
				"name":        "Robótica Móvil y Sistemas Autónomos",
				"description": "Diseño, navegación, cinemática y sensado de robots terrestres y aéreos.",
			},
			{
				"name":        "Manufactura Aditiva y Prototipado 3D",
				"description": "Modelado CAD, optimización de impresión FDM/SLA y desarrollo de piezas mecánicas.",
			},
			{
				"name":        "Internet de las Cosas (IoT) y Sistemas Embebidos",
				"description": "Telemetría, microcontroladores, protocolos industriales y computación en el borde.",
			},
			{
				"name":        "Inteligencia Artificial y Visión por Computador",
				"description": "Reconocimiento de patrones, visión artificial para inspección y control guiado por imagen.",
			},
		},
	}

	JSON(w, http.StatusOK, info)
}

// GetDirectors retorna los perfiles públicos del equipo directivo
func (h *HubHandler) GetDirectors(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	query := `SELECT full_name, email, avatar_url, bio
	          FROM users 
	          WHERE role = 'DIRECTOR' AND is_active = TRUE
	          ORDER BY full_name ASC`

	rows, err := h.db.Query(r.Context(), query)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar directores", err.Error())
		return
	}
	defer rows.Close()

	directors := make([]map[string]interface{}, 0)
	for rows.Next() {
		var fullName, email string
		var avatarURL, bio *string

		if err := rows.Scan(&fullName, &email, &avatarURL, &bio); err != nil {
			ErrorWithDetails(w, http.StatusInternalServerError, "Error al leer datos de directores", err.Error())
			return
		}

		directors = append(directors, map[string]interface{}{
			"full_name":  fullName,
			"email":      email,
			"avatar_url": avatarURL,
			"bio":        bio,
			"role_title": "Director de Semillero",
		})
	}

	JSON(w, http.StatusOK, directors)
}

// GetPublicProjects retorna el portafolio público: estrictamente proyectos concluidos (COMPLETED)
func (h *HubHandler) GetPublicProjects(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Base de datos no disponible")
		return
	}

	query := `SELECT p.id, p.title, p.description, p.research_line, p.status, p.created_by, p.start_date, p.target_end_date, p.created_at, p.updated_at,
	                 u.full_name AS creator_name
	          FROM projects p
	          JOIN users u ON p.created_by = u.id
	          WHERE p.status = 'COMPLETED'
	          ORDER BY p.updated_at DESC`

	rows, err := h.db.Query(r.Context(), query)
	if err != nil {
		ErrorWithDetails(w, http.StatusInternalServerError, "Error al consultar portafolio público", err.Error())
		return
	}
	defer rows.Close()

	type PublicProjectItem struct {
		domain.Project
		CreatorName string `json:"creator_name"`
	}

	projects := make([]PublicProjectItem, 0)
	for rows.Next() {
		var item PublicProjectItem
		if err := rows.Scan(
			&item.ID, &item.Title, &item.Description, &item.ResearchLine, &item.Status,
			&item.CreatedBy, &item.StartDate, &item.TargetEndDate, &item.CreatedAt, &item.UpdatedAt,
			&item.CreatorName,
		); err != nil {
			ErrorWithDetails(w, http.StatusInternalServerError, "Error al leer proyecto", err.Error())
			return
		}
		projects = append(projects, item)
	}

	JSON(w, http.StatusOK, projects)
}
