# SEMARD Control Center

> **Semillero de Investigación SEMARD — Universidad de Cartagena**  
> **Backend:** Go (Golang) | **Base de Datos:** PostgreSQL (Coolify VM) | **Frontend:** PWA (Next.js / React)  
> 📋 **Plan de Desarrollo del Equipo:** [PLAN_DE_DESARROLLO.md](./PLAN_DE_DESARROLLO.md) (2 Backend + 2 Frontend)

---

## 1. Definición de los 4 Roles del Sistema

El sistema cuenta con 4 roles principales:

1. **Directores**: Máxima autoridad del semillero. Creación de proyectos, creación y gestión de eventos (públicos e internos), evaluación de avances y recursos, decisiones institucionales y aprobación de nuevos registros.
2. **Administrador**: Gestión logística, inventario, aprobación de préstamos y de solicitudes 3D.
3. **Miembros**: Estudiantes integrantes del semillero con cuenta activa.
   * Pueden ser asignados a proyectos como **encargados**.
   * Pueden solicitar préstamos de insumos y solicitudes de impresión 3D.
   * Aquellos con el permiso especial de **Operador 3D** pueden gestionar y ejecutar impresiones aprobadas.
4. **Público (Visitantes / Comunidad Universitaria)**: Usuarios sin cuenta o sin sesión iniciada.

---

## 2. Matriz Completa de Roles y Permisos (Hub + App Logística)

| Módulo / Funcionalidad | 4. Público | 3. Miembros | 2. Administrador | 1. Directores | Permiso Especial: Operador 3D |
| :--- | :---: | :---: | :---: | :---: | :---: |
| **Ver Información Institucional & Líneas** | ✅ **Acceso Libre** | ✅ Acceso | ✅ Acceso | ✅ Acceso | Hereda Miembro |
| **Ver Perfil y Biografía de Directores** | ✅ **Acceso Libre** | ✅ Acceso | ✅ Acceso | ✅ Acceso | Hereda Miembro |
| **Ver y Registrarse en Eventos Abiertos** | ✅ **Acceso Libre** | ✅ Acceso | ✅ Acceso | ✅ Acceso | Hereda Miembro |
| **Ver Proyectos Finalizados (Portafolio Público)** | ✅ **Solo Finalizados** | ✅ Acceso | ✅ Acceso | ✅ Acceso | Hereda Miembro |
| **Ver Vitrina General de Proyectos y sus Estados** | ❌ **Solo Miembros** | ✅ **Acceso Total** | ✅ **Acceso Total** | ✅ **Acceso Total** | Hereda Miembro |
| **Crear Solicitud de Ingreso (Google Auth)** | ✅ **Puede solicitar** | Ya es miembro | Ya es miembro | Ya es miembro | - |
| **Ver Eventos Internos (Reuniones/Actas)** | ❌ Requiere cuenta | ✅ Acceso | ✅ Acceso | ✅ Acceso | Hereda Miembro |
| **Crear y Publicar Eventos (Públicos/Internos)** | ❌ | ❌ | ❌ | ✅ **Exclusivo Directores** | ❌ |
| **Solicitar Préstamos de Equipos** | ❌ | ✅ Puede solicitar | ✅ Puede solicitar | ✅ Puede solicitar | Hereda Miembro |
| **Aprobar / Modificar Plazo / Rechazar Préstamos** | ❌ | ❌ | ✅ **Aprueba / Rechaza** | ✅ **Aprueba / Rechaza** | ❌ |
| **Solicitar Impresión 3D** | ❌ | ✅ Puede solicitar | ✅ Puede solicitar | ✅ Puede solicitar | Hereda Miembro |
| **Aprobar / Rechazar Solicitudes 3D** | ❌ | ❌ | ✅ **Aprueba / Rechaza** | ✅ **Aprueba / Rechaza** | ❌ No aprueba |
| **Ejecutar y Gestionar Cola de Impresión 3D** | ❌ | ❌ | ✅ Ejecuta aprobadas | ✅ Ejecuta aprobadas | 🖨️ **Ejecuta aprobadas** |
| **Crear Proyectos del Semillero** | ❌ | ❌ | ❌ | ✅ **Exclusivo Directores** | ❌ |
| **Ser Asignado como Encargado de Proyecto** | ❌ | ✅ Asignable | ✅ Asignable | ✅ Asignable | - |
| **Actualizar Proyecto y Enviar Avances** | ❌ | 📝 Solo si es encargado | 📝 Solo si es encargado | 📝 Solo si es encargado | - |
| **Solicitar Recursos para Proyectos** | ❌ | 📝 Solo si es encargado | 📝 Solo si es encargado | 📝 Solo si es encargado | - |
| **Aprobar Avances con Feedback** | ❌ | ❌ | ❌ | ✅ **Exclusivo Directores** | ❌ |
| **Aprobar / Rechazar / Devolver Recursos** | ❌ | ❌ | ❌ | ✅ **Exclusivo Directores** | ❌ |
| **Aprobar / Rechazar Solicitudes de Registro** | ❌ | ❌ | ❌ | ✅ **Exclusivo Directores** | ❌ |

---

## 3. División Modular del Sistema

### Módulo 1: Hub Institucional SEMARD
* **1.1 Eventos Abiertos al Público:** Agenda, fechas, ponentes y formulario de registro para ferias, talleres abiertos y conferencias (creación exclusiva por Directores). **Accesible para el Público**.
* **1.2 Información del Semillero & Portafolio Público:**
  * Misión, visión, historia y líneas de investigación activas.
  * **Portafolio Público de Proyectos:** Exhibición para el público general **exclusivamente de los proyectos que se encuentren marcados como finalizados**.
* **1.3 Perfiles de Directores:** Tarjetas con fotografía profesional, biografía resumida, trayectoria investigativa y enlaces a perfiles académicos. **Accesible para el Público**.
* **1.4 Eventos Internos:** Cartelera protegida con agenda de reuniones semanales, entregas de avances, sustentaciones y actas (creación exclusiva por Directores). Visible únicamente para integrantes autenticados (**Miembros**, **Administradores**, **Directores**).

### Módulo 2: App de Gestión Logística y Operativa

#### 2.1 Gestión Integral de Proyectos
* **2.1.1 Creación y Asignación:**
  * Los proyectos **solo pueden ser creados por Directores**.
  * Pueden ser asignados a **Miembros, Administradores o Directores**. Estos usuarios asignados serán los **encargados** que podrán actualizar el proyecto.
* **2.1.2 Avances y Aprobación:**
  * Los avances realizados por los encargados del proyecto deberán ser evaluados y **aprobados por uno de los Directores**.
* **2.1.3 Feedback de Avances:**
  * Los encargados del proyecto reciben **feedback de los Directores** acerca de cada avance evaluado.
* **2.1.4 Solicitudes de Recursos del Proyecto:**
  * Los encargados del proyecto (sin importar su rol) pueden realizar peticiones de recursos para el proyecto (*digital, económico, de conocimiento, etc.*).
  * Esta solicitud puede ser **aprobada, rechazada, o devuelta para modificaciones**, cada una con su respectiva razón por uno de los Directores.

#### 2.2 Préstamos de Recursos y Herramientas
* Catálogo de inventario (multímetros, osciloscopios, fuentes, kits de desarrollo, guías técnicas).
* Solicitud: dispositivo/elemento a prestar, la razón y el tiempo de préstamo solicitado.
* Decisión por **Directores o Administradores**:
  * **Aprobar:** Plazo solicitado o **modificar el plazo** (el solicitante recibe una notificación con el ajuste).
  * **Rechazar:** Quien la rechaza debe dar un feedback de la razón del rechazo.
* Control de entrega, devolución y checklist de estado.

#### 2.3 Gestión de Solicitudes de Impresión 3D
* **Solicitud:** Carga de archivo 3D (.STL, .OBJ, .STEP), material, color, relleno (%) y vinculación a proyecto o justificación.
* **Aprobación o Rechazo:** Exclusivo de **Directores y Administradores**. Si es rechazada, se debe dar feedback de la razón del rechazo.
* **Gestión y Ejecución:** Una vez aprobada, puede ser gestionada y ejecutada por **Directores, Administradores y Estudiantes con permiso de Operador 3D**.

#### 2.4 Diferenciación de Vistas y Accesos a Proyectos (3 Niveles)
1. **Portafolio Público (Hub):** El público general únicamente tiene acceso a ver los proyectos marcados como **finalizados**, sirviendo como vitrina de logros y divulgación científica.
2. **Vitrina General de Proyectos (Interna):** Sección accesible para **todos los miembros del semillero** (Miembros, Administradores, Directores) donde pueden consultar **todos los proyectos y su estado en tiempo real** (`ACTIVE`, `PAUSED`, `COMPLETED`).
3. **Panel Operativo de Trabajo ("Mis Proyectos"):** Sección apartada de la creación y actualización general, donde los **encargados** (miembros, administradores o directores asignados) únicamente visualizan y gestionan los proyectos en los que participan directamente.

---

## 4. Arquitectura de Backend en Go & Despliegue

### 4.1 Tecnologías del Backend
* **Lenguaje:** Go (Golang) 1.23+
* **Enrutador & Middleware:** `go-chi/chi/v5`
* **Driver PostgreSQL:** `jackc/pgx/v5` con pool de conexiones nativo
* **Autenticación:** Google OAuth 2.0 (`golang.org/x/oauth2`) y tokens JWT (`golang-jwt/jwt/v5`)
* **Validaciones:** `go-playground/validator/v10`
* **Almacenamiento de Archivos:** Volumen montado persistente en Coolify para modelos 3D y anexos

### 4.2 Infraestructura (Coolify VM)
* **Orquestación:** Coolify sobre Máquina Virtual propia.
* **Base de Datos:** PostgreSQL en contenedor Docker administrado en Coolify.
* **Backend:** Contenedor Docker multi-stage en Go (peso estimado del binario final: ~20 MB).

---

## 5. Flujo de Autenticación y Registro Híbrido

1. **Autenticación con Google:**
   * El usuario inicia sesión utilizando su cuenta institucional con dominio `@unicartagena.edu.co`.
   * El backend valida el token de Google y verifica que el dominio sea estrictamente universitario.
2. **Cruce contra la Base de Datos:**
   * **Caso 1 (Pre-registrado):** Si el correo ya fue cargado previamente por un Director (con su código estudiantil y rol asignado), el usuario obtiene acceso instantáneo con sus credenciales y permisos correspondientes vía JWT.
   * **Caso 2 (No registrado):** Si el correo institucional no existe en el sistema:
     1. El sistema solicita al estudiante ingresar su **Código Estudiantil**, programa académico y motivación.
     2. Se genera un registro en la tabla `registration_requests` en estado `PENDING`.
     3. El usuario recibe un mensaje indicando que su solicitud está en evaluación.
     4. Los **Directores** visualizan las solicitudes pendientes en su panel y deciden **Aprobar** (asignándole rol de Miembro o Administrador) o **Rechazar** con feedback. Una vez aprobada, el usuario ya puede ingresar normalmente con su cuenta de Google.

---

## 6. Modelo de Datos y Entidades

La persistencia de datos se gestiona en **PostgreSQL** mediante migraciones versionadas y controladas, estructurando las siguientes entidades principales:
* **Usuarios y Membresías:** Cuentas institucionales (`@unicartagena.edu.co`), roles del sistema (`DIRECTOR`, `ADMINISTRADOR`, `MIEMBRO`), permiso especial de `Operador 3D` y solicitudes de registro.
* **Proyectos de Investigación:** Proyectos creados por directores, asignación de encargados, bitácora de avances, retroalimentación y solicitudes de recursos.
* **Inventario y Préstamos:** Catálogo de herramientas y equipos electrónicos, solicitudes de préstamo con control de tiempos, notificaciones por ajustes y trazabilidad de devoluciones.
* **Taller de Impresión 3D:** Solicitudes de manufactura aditiva con parámetros técnicos y almacenamiento de archivos 3D, aprobaciones y cola de ejecución.
* **Eventos:** Agenda y cartelera de eventos con visibilidad pública o interna protegida.

---

## 7. Plan de Desarrollo Dividido en Equipos (2 Backend + 2 Frontend)

Para optimizar el trabajo del equipo de **4 desarrolladores**, el proyecto se organiza en dos tracks de desarrollo paralelo sincronizados mediante **contratos de API (Swagger / OpenAPI)**:

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

---

### 7.1 Estructura del Repositorio (Monorepo)
```
Semard-Control-Center/
├── backend/            # Track Backend (Go 1.23+, chi, pgx, migraciones)
│   ├── cmd/api/
│   ├── internal/
│   ├── migrations/
│   ├── Dockerfile
│   └── docker-compose.yml
├── frontend/           # Track Frontend (Next.js 15, React, Tailwind, PWA)
│   ├── src/
│   │   ├── app/
│   │   ├── components/
│   │   ├── lib/
│   │   └── hooks/
│   ├── public/
│   └── package.json
└── docs/               # Documentación y Contratos de API (OpenAPI / Postman)
```

---

### 7.2 Track Backend (2 Desarrolladores)

* **Stack:** Go 1.23+, `go-chi/chi/v5`, `jackc/pgx/v5`, `golang-migrate`, Docker y Coolify.
* **Distribución de Roles:**
  - **Dev B1 (Infra, Auth & Proyectos):** Arquitectura base, conexión DB, autenticación Google OAuth + JWT, RBAC y módulo central de proyectos y recursos.
  - **Dev B2 (Hub, Inventario & Taller 3D):** Módulo de eventos, catálogo de herramientas/préstamos, almacenamiento de archivos STL/OBJ y cola de impresión 3D.

#### Fases del Backend:
* **Fase B1: Infraestructura, Migraciones y Auth Google (Semana 1)**
  - Configuración del módulo de Go, Clean Architecture y contenedor multi-stage.
  - Creación de migraciones SQL con todas las tablas del sistema.
  - Flujo Google OAuth 2.0 (`@unicartagena.edu.co`) con emisión de JWT y verificación de padrón vs `registration_requests`.
  - Endpoint de Healthcheck (`/healthz`) y despliegue inicial en Coolify.
* **Fase B2: Gestión de Miembros y Módulo del Hub (Semana 2)**
  - Endpoints para Directores: aprobación/rechazo de solicitudes de registro con feedback.
  - Endpoints informativos del semillero y equipo directivo.
  - Portafolio público: filtro de proyectos finalizados (`COMPLETED`).
  - Módulo de Eventos: creación exclusiva por Directores, listados públicos y eventos internos protegidos.
* **Fase B3: Módulo Logístico de Proyectos y Recursos (Semana 3)**
  - Creación de proyectos (exclusivo Directores) y asignación de encargados.
  - Subida de avances por encargados y evaluación con feedback por Directores.
  - Peticiones de recursos (digital, económico, conocimiento, hardware) con ciclo de 3 estados (`APPROVED`, `REJECTED`, `RETURNED_FOR_MODIFICATION`).
  - Endpoints diferenciados: Vitrina general interna vs "Mis Proyectos".
* **Fase B4: Inventario, Préstamos y Taller de Impresión 3D (Semana 4)**
  - CRUD de inventario y solicitudes de préstamo con ajuste de plazo y feedback.
  - Subida y almacenamiento seguro de modelos 3D (.STL, .OBJ).
  - Aprobación/rechazo técnico 3D por Directores/Admins y cola de ejecución para Operadores 3D.
* **Fase B5: Storage, Notificaciones, Pruebas y Despliegue VM (Semana 5)**
  - Sistema de notificaciones in-app y alertas de correo/web push.
  - Pruebas unitarias e integración de endpoints críticos.
  - Configuración final de producción en Coolify.

---

### 7.3 Track Frontend (2 Desarrolladores)

* **Stack:** Next.js 15 (App Router), React 19, TypeScript, Tailwind CSS, shadcn/ui y soporte PWA.
* **Distribución de Roles:**
  - **Dev F1 (Core PWA, Auth & App Logística de Proyectos):** Shell de la aplicación, estado de autenticación Google, dashboard administrativo, gestión de proyectos, avances y recursos.
  - **Dev F2 (Hub Institucional, UI/UX, Inventario & Taller 3D):** Vistas públicas del Hub, eventos, vitrina, módulo visual de inventario/préstamos, visor/subida de archivos 3D y manifiesto PWA.

#### Fases del Frontend:
* **Fase F1: Setup PWA, Design System y Flujo de Login (Semana 1)**
  - Inicialización del proyecto Next.js con Tailwind CSS, TypeScript y componentes base.
  - Configuración inicial de PWA (Web App Manifest, iconos e instalabilidad).
  - Pantalla de inicio de sesión con Google institucional.
  - Formulario modal de solicitud de ingreso para estudiantes no registrados previamente.
* **Fase F2: SEMARD Hub y Gestión de Eventos (Semana 2)**
  - Landing page institucional con áreas de investigación y portafolio de proyectos concluidos.
  - Perfiles y biografías de Directores.
  - Cartelera pública de eventos con formulario de inscripción.
  - Cartelera de eventos internos con actas (protegida para miembros autenticados).
  - Panel de creación de eventos para Directores.
* **Fase F3: Panel de Proyectos, Avances y Recursos (Semana 3)**
  - Vitrina interna de proyectos para consulta de todos los miembros.
  - Espacio de trabajo "Mis Proyectos" para los encargados asignados.
  - Formulario de entrega de avances y visor de feedback del Director.
  - Interfaz de solicitudes de recursos con badges de estado (Aprobada, Rechazada, Devuelta con observaciones).
  - Panel para Directores: aprobación de avances y revisión de recursos con retroalimentación.
* **Fase F4: Módulo de Inventario, Préstamos e Impresión 3D (Semana 4)**
  - Catálogo interactivo de herramientas electrónicas y formulario de solicitud de préstamo.
  - Panel de gestión de préstamos para Directores/Administradores (aprobar, modificar fecha o rechazar).
  - Formulario de solicitud de impresión 3D con carga de archivos STL/OBJ y selección de parámetros.
  - Tablero Kanban o cola de impresión para Directores, Administradores y Operadores 3D.
* **Fase F5: Cacheo Offline, Notificaciones y Auditoría PWA (Semana 5)**
  - Configuración del Service Worker para cacheo offline del Hub y eventos.
  - Centro de notificaciones in-app con campana y push notifications.
  - Pruebas responsivas en móvil y escritorio, auditoría Lighthouse (PWA > 90%).

---

### 7.4 Hitos de Integración y Sincronización (Milestones 2+2)

Para evitar bloqueos y asegurar que ambos equipos avancen coordinadamente:

| Hito | Semana | Objetivo de Integración | Verificación Conjunta |
| :--- | :---: | :--- | :--- |
| **M1: Auth & Setup** | Fin Semana 1 | Login de Google conectado a la API de Go y generación de sesión JWT en el cliente. | Inicio de sesión funcional en frontend contra el backend en Go. |
| **M2: Hub & Eventos** | Fin Semana 2 | Hub público alimentado dinámicamente desde la base de datos y eventos operativos. | Crear un evento como Director y visualizarlo en la vista pública e interna. |
| **M3: Proyectos & Flujo** | Fin Semana 3 | Ciclo completo de proyectos: creación, avances, feedback y recursos. | Un encargado envía avance/recurso y el Director lo evalúa con feedback en la UI. |
| **M4: Logística & 3D** | Fin Semana 4 | Préstamos de equipos e impresión 3D conectados de punta a punta. | Solicitar préstamo/impresión y gestionarlo como Admin/Operador 3D desde la interfaz. |
| **M5: Despliegue Final** | Fin Semana 5 | Integración completa, pruebas E2E y despliegue de Backend + Frontend en Coolify. | Auditoría PWA, instalación en smartphone y prueba integral en producción. |
