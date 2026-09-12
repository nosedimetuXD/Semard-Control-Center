# Plan de Desarrollo y Ejecución Técnica: SEMARD Control Center

> **Semillero de Investigación SEMARD — Universidad de Cartagena**  
> **Estrategia de Ejecución:** Desarrollo en Paralelo (4 Desarrolladores: 2 Backend + 2 Frontend)  
> **Duración Estimada:** 5 Semanas (5 Sprints)

---

## 1. Organización del Equipo de Trabajo

Para maximizar la eficiencia y evitar cuellos de botella, el equipo de **4 desarrolladores** se divide en dos tracks paralelos coordinados mediante contratos de interfaz (API Contracts):

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       ARQUITECTURA DEL EQUIPO (4 DEVS)                      │
├──────────────────────────────────────┬──────────────────────────────────────┤
│       TRACK BACKEND (2 Devs)         │       TRACK FRONTEND (2 Devs)        │
│          [ Go + PostgreSQL ]         │      [ Next.js + PWA + Tailwind ]    │
└──────────────────┬───────────────────┴──────────────────┬───────────────────┘
                   │                                      │
                   ▼                                      ▼
           [ Dev B1 & Dev B2 ]                    [ Dev F1 & Dev F2 ]
       Infra, Modelos, Lógica API             UI/UX, Estado, PWA, Vistas
                   │                                      │
                   └───────────────► ◄────────────────────┘
                          Hitos de Integración
```

### Roles y Responsabilidades:
* **Dev B1 (Backend Lead - Infra, Auth & Proyectos):**
  - Configuración del servidor en Go, conexión PostgreSQL con `pgx/v5`, pool y variables de entorno.
  - Flujo de autenticación con Google OAuth 2.0 (`@unicartagena.edu.co`), generación de JWT y middleware RBAC.
  - Lógica central de proyectos: creación (Directores), asignación de encargados, entregables/avances con feedback y ciclo de solicitudes de recursos.
* **Dev B2 (Backend - Hub, Inventario & Taller 3D):**
  - Módulo del Hub informativo, equipo directivo y portafolio público de proyectos finalizados.
  - Módulo de eventos públicos e internos (creación exclusiva por Directores).
  - Catálogo de inventario, solicitudes de préstamo con ajuste de plazo y feedback obligatorio.
  - Almacenamiento seguro de archivos 3D (.STL, .OBJ) y cola de ejecución para Operadores 3D.
* **Dev F1 (Frontend Lead - Core PWA, Auth & App Logística):**
  - Inicialización del proyecto en Next.js 15 (App Router) con TypeScript, Tailwind CSS y componentes base.
  - Gestión de sesión y estado global de autenticación institucional con Google.
  - Panel operativo de proyectos: vitrina interna, "Mis Proyectos" para encargados, formularios de avances con visor de feedback y módulo de recursos.
* **Dev F2 (Frontend - Hub Institucional, UI/UX, Inventario & 3D):**
  - Diseño visual e interfaz pública del Hub (landing page, líneas de investigación, biografía de Directores).
  - Cartelera pública de eventos, formulario de inscripción abierta y cartelera de eventos internos protegidos.
  - Catálogo interactivo de herramientas de laboratorio y formulario de préstamos.
  - Interfaz de carga de archivos 3D y tablero Kanban de la cola de impresión para Operadores 3D.

---

## 2. Estructura del Repositorio (Monorepo)

El repositorio se organizará en carpetas independientes para desacoplar el ciclo de vida del código:

```
Semard-Control-Center/
├── backend/                  # Track Backend
│   ├── cmd/api/main.go       # Entrada del servidor HTTP
│   ├── internal/
│   │   ├── config/           # Lectura de variables de entorno (.env)
│   │   ├── database/         # Pool pgx y migraciones automáticas
│   │   ├── domain/           # Modelos de entidades y contratos de repositorios
│   │   ├── repository/       # Consultas SQL nativas a PostgreSQL
│   │   ├── service/          # Lógica de negocio y reglas de validación
│   │   ├── handler/          # Controladores HTTP REST JSON (chi)
│   │   └── middleware/       # JWT Auth, RBAC, Logger, CORS, Recover
│   ├── migrations/           # Scripts SQL versionados (golang-migrate)
│   ├── storage/              # Volumen persistente para archivos 3D y adjuntos
│   ├── Dockerfile            # Multi-stage optimizado (Go Scratch / Alpine ~20MB)
│   └── docker-compose.yml    # Orquestación con PostgreSQL para Coolify
│
├── frontend/                 # Track Frontend
│   ├── src/
│   │   ├── app/              # Rutas Next.js 15 (App Router)
│   │   │   ├── (public)/     # Hub: Inicio, Directores, Eventos, Portafolio
│   │   │   ├── (auth)/       # Login Google, Solicitud de Membresía
│   │   │   └── (dashboard)/  # App Logística: Proyectos, Inventario, 3D, Eventos internos
│   │   ├── components/       # Componentes reutilizables (UI, Modales, Tablas, Forms)
│   │   ├── lib/              # Cliente API HTTP, utilidades y helpers
│   │   └── hooks/            # Custom hooks (useAuth, useProjects, useNotifications)
│   ├── public/               # Manifest PWA, iconos, service worker
│   └── package.json
│
└── docs/                     # Contratos de API
    └── openapi.yaml          # Especificación OpenAPI/Swagger para mocks de Frontend
```

---

## 3. Metodología de Trabajo: "API First" y Mocks

Para que los desarrolladores de Frontend no dependan de que el Backend esté finalizado:
1. **Día 1 y 2:** El equipo acuerda y documenta los contratos de API en `docs/openapi.yaml` (rutas, payloads JSON de entrada y respuestas).
2. **Desarrollo en Frontend con Mocks:** Los Devs F1 y F2 consumen datos simulados basados en el contrato mientras los Devs B1 y B2 implementan los endpoints reales.
3. **Conexión:** En los días finales de cada sprint se reemplazan los mocks por los endpoints del Backend.

---

## 4. Cronograma Detallado por Sprints (5 Semanas)

### Semana 1: Cimientos de Infraestructura, Setup y Autenticación Google

#### Track Backend (Dev B1 & Dev B2)
- [ ] Inicialización del módulo de Go (`go mod init semard-api`).
- [ ] Implementación de Clean Architecture (`cmd/`, `internal/`).
- [ ] Conexión a PostgreSQL con pool `pgx/v5` y motor de migraciones (`golang-migrate`).
- [ ] Configuración del `Dockerfile` multi-stage y `docker-compose.yml` para Coolify.
- [ ] Endpoints de autenticación:
  - `GET /api/v1/auth/google/login`: Genera URL de consentimiento Google.
  - `GET /api/v1/auth/google/callback`: Valida dominio `@unicartagena.edu.co`.
  - Verificación contra tabla de usuarios: si existe, genera JWT con rol y permisos.
  - Si no existe: permite registrar solicitud en `registration_requests` (`POST /api/v1/auth/register-request`).
- [ ] Middlewares: Logger estructurado, CORS, Recover y Healthcheck (`GET /healthz`).

#### Track Frontend (Dev F1 & Dev F2)
- [ ] Inicialización del proyecto Next.js 15 con TypeScript y Tailwind CSS.
- [ ] Configuración del manifiesto web (`manifest.json`) e iconos PWA iniciales.
- [ ] Implementación del Design System base (paleta de colores, tipografía, botones, inputs).
- [ ] Pantalla de Login con botón institucional "Iniciar sesión con Google".
- [ ] Modal de postulación para estudiantes no registrados previamente (captura de código estudiantil, programa académico y motivación).
- [ ] Cliente HTTP configurado con interceptores para inyección de token JWT.

#### Hito de Integración M1:
> **Verificación:** Inicio de sesión con cuenta institucional `@unicartagena.edu.co` desde el frontend contra la API de Go, guardando el token JWT y accediendo al estado autenticado.

---

### Semana 2: Gestión de Miembros, Hub Institucional y Eventos

#### Track Backend (Dev B1 & Dev B2)
- [ ] **Dev B1:**
  - Endpoints para Directores: listar solicitudes de membresía pendientes (`GET /api/v1/admin/registration-requests`).
  - Decisión de Directores: aprobar asignando rol o rechazar con retroalimentación (`POST /api/v1/admin/registration-requests/{id}/approve|reject`).
  - Middleware de protección RBAC por roles (`DIRECTOR`, `ADMINISTRADOR`, `MIEMBRO`).
- [ ] **Dev B2:**
  - Endpoints públicos del Hub: información del semillero, líneas de investigación y perfiles de Directores (`GET /api/v1/hub/info`, `GET /api/v1/hub/directors`).
  - Portafolio público: `GET /api/v1/hub/projects` (filtro estricto: **solo proyectos finalizados** `COMPLETED`).
  - Módulo de Eventos:
    - CRUD exclusivo para Directores (`POST`, `PUT`, `DELETE /api/v1/events`).
    - Consulta pública de eventos abiertos (`GET /api/v1/events/public`) y registro de asistentes externos.
    - Consulta privada de eventos internos protegida para miembros (`GET /api/v1/events/internal`).

#### Track Frontend (Dev F1 & Dev F2)
- [ ] **Dev F1:**
  - Panel administrativo para Directores: tabla de solicitudes de nuevo ingreso pendientes con botones de aprobar (asignar rol) o rechazar con feedback.
  - Guardias de navegación en frontend según el rol autenticado.
- [ ] **Dev F2:**
  - Landing page del SEMARD Hub: misión, visión, áreas de investigación y portafolio público de proyectos concluidos.
  - Vista de equipo directivo con fotografías, trayectorias y enlaces académicos.
  - Cartelera de eventos abiertos al público con formulario de inscripción para visitantes.
  - Cartelera protegida de eventos internos (reuniones semanales, actas y compromisos).
  - Formulario de creación/edición de eventos disponible exclusivamente para usuarios con rol Director.

#### Hito de Integración M2:
> **Verificación:** Un Director aprueba una solicitud de membresía real, crea un evento público y uno interno, y se valida la correcta visualización en el Hub público y en el dashboard privado.

---

### Semana 3: Módulo Logístico de Proyectos, Avances y Recursos

#### Track Backend (Dev B1 & Dev B2)
- [ ] **Dev B1 (Proyectos y Avances):**
  - Creación de proyectos (`POST /api/v1/projects`) restringida exclusivamente a usuarios con rol `DIRECTOR`.
  - Asignación de encargados del proyecto (Miembros, Administradores o Directores).
  - Endpoint "Mis Proyectos" (`GET /api/v1/projects/my-projects`) para los encargados asignados.
  - Envío de avances por encargados con enlaces/archivos adjuntos (`POST /api/v1/projects/{id}/updates`).
  - Evaluación de avances por Directores (`POST /api/v1/projects/updates/{id}/review`) con estado (`APPROVED` / `CHANGES_REQUESTED`) y **feedback técnico obligatorio**.
- [ ] **Dev B2 (Recursos y Vitrina):**
  - Solicitudes de recursos para proyectos por parte de encargados (`POST /api/v1/projects/{id}/resources`).
  - Evaluación de recursos por Directores en 3 estados: `Aprobada`, `Rechazada` o `Devuelta para Modificaciones` con feedback justificativo.
  - Vitrina general interna de proyectos (`GET /api/v1/projects/showcase`) que muestra todos los proyectos y sus estados para cualquier miembro autenticado.

#### Track Frontend (Dev F1 & Dev F2)
- [ ] **Dev F1 (Gestión de Proyectos para Encargados y Directores):**
  - Formulario de creación de proyectos con asignación de integrantes (solo visible para Directores).
  - Vista "Mis Proyectos": panel donde los encargados visualizan los proyectos que lideran.
  - Formulario para enviar entregables y avances con subida de comprobantes.
  - Visor de feedback del Director en cada avance con alertas de estado.
  - Panel de revisión de avances para Directores con editor de comentarios y retroalimentación.
- [ ] **Dev F2 (Recursos y Vitrina General):**
  - Vitrina general interna de proyectos: tarjetas con filtros por línea de investigación y estado.
  - Módulo de peticiones de recursos con categorías (digital, económica, conocimiento, hardware).
  - Visualización de estados de recursos con badges interactivos (`Aprobada`, `Rechazada`, `Requiere Modificación`).
  - Panel para Directores para resolver peticiones de recursos con campo obligatorio de justificación.

#### Hito de Integración M3:
> **Verificación:** Un Director crea un proyecto y asigna encargados; un encargado sube un avance y pide un recurso; el Director revisa el avance con feedback y devuelve el recurso para modificaciones con comentarios visibles para el encargado.

---

### Semana 4: Inventario, Préstamos de Equipos y Taller de Impresión 3D

#### Track Backend (Dev B1 & Dev B2)
- [ ] **Dev B1 (Préstamos e Inventario):**
  - CRUD de inventario de instrumental y guías con stock dinámico.
  - Solicitud de préstamo de insumos (`POST /api/v1/loans/requests`).
  - Resolución por Directores o Administradores:
    - Aprobar por tiempo solicitado o **aprobar modificando el plazo**.
    - Rechazar con feedback obligatorio.
  - Registro de devolución física y verificación de estado.
- [ ] **Dev B2 (Taller de Impresión 3D):**
  - Carga segura de modelos 3D (.STL, .OBJ, .STEP) a volumen persistente (`POST /api/v1/print3d/requests`).
  - Aprobación o rechazo técnico exclusivo por **Directores y Administradores** con feedback obligatorio en caso de rechazo.
  - Cola de trabajo operativa para **Directores, Administradores y Operadores 3D**:
    - Transición de estados: `APPROVED` ➔ `IN_PROGRESS` ➔ `COMPLETED` ➔ `DELIVERED`.

#### Track Frontend (Dev F1 & Dev F2)
- [ ] **Dev F1 (Módulo de Préstamos):**
  - Catálogo interactivo de equipos con disponibilidad en tiempo real.
  - Formulario de solicitud de préstamo (selección de equipo, fechas solicitadas y justificación).
  - Panel para Directores y Administradores para aprobar, modificar la fecha de entrega o rechazar con motivo.
  - Historial de préstamos del usuario con notificaciones en caso de que su plazo haya sido ajustado.
- [ ] **Dev F2 (Módulo de Impresión 3D):**
  - Formulario de solicitud de impresión 3D: carga de archivo .STL/.OBJ, selección de material, color, porcentaje de relleno y proyecto asociado.
  - Panel de evaluación técnica para Directores y Administradores.
  - Tablero de control de producción (Kanban / Cola de impresión) accesible para Directores, Administradores y Operadores 3D certificados para actualizar el progreso del trabajo.

#### Hito de Integración M4:
> **Verificación:** Un miembro solicita un multímetro y una pieza 3D; el Administrador aprueba el préstamo recortando el plazo (generando alerta) y aprueba la impresión 3D; un estudiante con rol de Operador 3D toma el trabajo y lo marca como completado.

---

### Semana 5: PWA Offline, Notificaciones, Pruebas y Despliegue en Coolify

#### Track Backend (Dev B1 & Dev B2)
- [ ] Servicio de notificaciones internas y alertas por correo electrónico.
- [ ] Pruebas unitarias en servicios críticos (validación de roles, cálculo de fechas de préstamos, transiciones de estado).
- [ ] Pruebas de integración de la API con PostgreSQL.
- [ ] Optimización de queries e índices en base de datos.
- [ ] Configuración del despliegue productivo en la máquina virtual mediante Coolify.

#### Track Frontend (Dev F1 & Dev F2)
- [ ] Configuración del Service Worker para cacheo offline de eventos y datos estáticos del Hub.
- [ ] Centro de notificaciones in-app con campana y soporte para Web Push Notifications.
- [ ] Pruebas de usabilidad responsiva en smartphones Android/iOS y navegadores de escritorio.
- [ ] Auditoría Lighthouse para garantizar cumplimiento PWA (> 90% en performance, PWA y accesibilidad).
- [ ] Compilación de producción y despliegue del frontend en Coolify.

#### Hito de Integración M5:
> **Verificación Final:** Plataforma completamente operativa en la VM con Coolify, instalable en dispositivos móviles como PWA, con navegación offline en el Hub y todos los flujos de autenticación, proyectos, préstamos y 3D verificados de punta a punta.

---

## 5. Matriz Resumen de Hitos de Integración

| Hito | Semana | Objetivo de Integración | Verificación Conjunta |
| :--- | :---: | :--- | :--- |
| **M1: Auth & Setup** | Fin Sem. 1 | Login Google institucional conectado entre Go y Next.js con JWT. | Inicio de sesión funcional en el frontend contra la API de Go. |
| **M2: Hub & Eventos** | Fin Sem. 2 | Hub y eventos alimentados dinámicamente desde PostgreSQL. | Crear un evento como Director y verlo en las vistas públicas/internas. |
| **M3: Proyectos & Flujo**| Fin Sem. 3 | Creación, asignación, avances con feedback y recursos. | Encargado envía avance/recurso y Director lo revisa y retroalimenta en la UI. |
| **M4: Logística & 3D** | Fin Sem. 4 | Préstamos de equipos e impresión 3D conectados de punta a punta. | Solicitar préstamo/impresión y gestionarlo como Admin/Operador 3D. |
| **M5: Despliegue Final** | Fin Sem. 5 | Pruebas integrales y puesta en producción en Coolify. | Instalación PWA en móvil y auditoría integral en la VM. |
