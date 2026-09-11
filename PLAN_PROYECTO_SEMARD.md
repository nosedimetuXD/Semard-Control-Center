# SEMARD Control Center

> **Semillero de Investigación SEMARD — Universidad de Cartagena**  
> **Backend:** Go (Golang) | **Base de Datos:** PostgreSQL (Coolify VM) | **Frontend:** PWA (Next.js / React)

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

## 7. Plan de Fases de Desarrollo (Roadmap por Etapas)

El desarrollo del proyecto se estructura en **6 etapas secuenciales**, garantizando entregables funcionales y verificables al final de cada hito:

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   Etapa 1    │ ──► │   Etapa 2    │ ──► │   Etapa 3    │
│ Infra & Base │     │ Auth & Users │     │ Hub & Eventos│
└──────────────┘     └──────────────┘     └──────────────┘
       │
       ▼
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   Etapa 4    │ ──► │   Etapa 5    │ ──► │   Etapa 6    │
│ Proyectos &  │     │ Inventario & │     │ Frontend PWA │
│ Recursos     │     │ Taller 3D    │     │ & Despliegue │
└──────────────┘     └──────────────┘     └──────────────┘
```

---

### Etapa 1: Cimientos de Infraestructura, Base de Datos y Arquitectura Go
* **Objetivo:** Establecer la base técnica del backend en Go bajo arquitectura limpia, orquestación con Docker y despliegue inicial en Coolify.
* **Tareas Clave:**
  - Inicialización del módulo de Go (`go mod init semard-api`) con Go 1.23+.
  - Estructuración modular limpia (`cmd/api`, `internal/config`, `internal/database`, `internal/domain`, `internal/repository`, `internal/service`, `internal/handler`, `internal/middleware`).
  - Configuración del sistema de migraciones versionadas (`golang-migrate`) con la base de datos PostgreSQL.
  - Creación del `Dockerfile` multi-stage optimizado (~20MB) y `docker-compose.yml` para despliegue local y en Coolify.
  - Middlewares base: Logging estructurado, CORS, Recover contra panics y endpoint de Healthcheck (`/healthz`).
* **Entregables:**
  - Repositorio backend funcional.
  - Base de datos PostgreSQL conectada y con migraciones ejecutadas automáticamente.
  - Pipeline de despliegue en Coolify operativo.

---

### Etapa 2: Módulo de Identidad, Autenticación Híbrida y Gestión de Usuarios
* **Objetivo:** Implementar el control de acceso, inicio de sesión institucional y el ciclo de admisión de miembros.
* **Tareas Clave:**
  - Integración de Google OAuth 2.0 con validación estricta de dominio (`@unicartagena.edu.co`).
  - Lógica de verificación contra el padrón de usuarios y generación de JWT con claims de rol (`DIRECTOR`, `ADMINISTRADOR`, `MIEMBRO`) y flag `can_operate_3d`.
  - Formulario y endpoint para solicitudes de nuevo registro (`registration_requests`) con código estudiantil, programa y motivación.
  - Panel de Directores: endpoints para listar solicitudes pendientes, aprobar (asignando rol) o rechazar con retroalimentación.
  - Middleware de autorización RBAC para proteger rutas según rol y permisos especiales.
* **Entregables:**
  - Autenticación con Google institucional 100% operativa.
  - Cola de solicitudes de membresía con flujo de aprobación por Directores.
  - Seguridad perimetral con tokens JWT y control de acceso por roles.

---

### Etapa 3: Módulo del Hub Institucional, Divulgación y Eventos
* **Objetivo:** Desarrollar los servicios de cara al público y la cartelera de eventos del semillero.
* **Tareas Clave:**
  - Endpoints públicos de información del semillero (misión, visión, líneas de investigación, contacto).
  - Endpoints públicos de perfiles directivos (fotos, biografías, redes académicas).
  - **Portafolio Público de Proyectos:** Consulta pública filtrada **exclusivamente para proyectos marcados como finalizados** (`COMPLETED`).
  - **Módulo de Eventos (Creación exclusiva por Directores):**
    - Eventos públicos: listado abierto, detalle de agenda y registro de asistentes externos.
    - Eventos internos: listado protegido para miembros del semillero con actas y cronograma interno.
* **Entregables:**
  - API completa del Hub para consulta abierta.
  - Sistema de gestión de eventos con control estricto de visibilidad (pública vs interna).

---

### Etapa 4: Módulo Logístico de Proyectos, Avances y Recursos
* **Objetivo:** Construir la plataforma para la gestión, seguimiento y financiamiento de proyectos de investigación.
* **Tareas Clave:**
  - **Creación de Proyectos:** Restringida exclusivamente a usuarios con rol `DIRECTOR`.
  - **Asignación de Encargados:** Asignación formal de Miembros, Administradores o Directores como encargados del proyecto.
  - **Avances y Entregables:** Carga de reportes y enlaces por los encargados (estado `PENDING`).
  - **Evaluación de Avances:** Aprobación exclusiva por Directores con **feedback técnico obligatorio**.
  - **Peticiones de Recursos:**
    - Solicitud por encargados en categorías (digital, económica, conocimiento, hardware).
    - Decisión por Directores en 3 estados: `Aprobada`, `Rechazada` o `Devuelta para Modificaciones` (con justificación obligatoria).
  - **Doble Vista de Proyectos:**
    - Vitrina general interna para todos los miembros del semillero (todos los proyectos y sus estados en tiempo real).
    - Panel operativo ("Mis Proyectos") donde los encargados solo interactúan con sus proyectos asignados.
* **Entregables:**
  - Flujo de vida completo de proyectos con trazabilidad de avances y feedbacks.
  - Motor de peticiones de recursos con ciclo de 3 estados.

---

### Etapa 5: Módulo Logístico de Inventario, Préstamos y Taller de Impresión 3D
* **Objetivo:** Automatizar el préstamo de instrumental de laboratorio y la producción de piezas 3D.
* **Tareas Clave:**
  - **Inventario:** CRUD de herramientas electrónicas, guías y consumibles con disponibilidad de stock.
  - **Préstamo de Recursos:**
    - Solicitud con justificación y tiempo de préstamo solicitado.
    - Aprobación por Directores o Administradores: capacidad de aprobar o **modificar el plazo concedido** (con notificación obligatoria al solicitante).
    - Rechazo con feedback obligatorio.
    - Control de entrega y devolución física.
  - **Taller de Impresión 3D:**
    - Carga y almacenamiento seguro de modelos 3D (.STL, .OBJ, .STEP) y especificaciones (material, color, infill).
    - Aprobación o rechazo técnico exclusivo por **Directores y Administradores** (con motivo técnico si se rechaza).
    - Cola de ejecución y cambio de estados (`Aprobada` ➔ `En Cola` ➔ `Imprimiendo` ➔ `Completada` ➔ `Entregada`) operada por **Directores, Administradores y Operadores 3D certificados**.
* **Entregables:**
  - Módulo de préstamos con control de plazos y notificaciones.
  - Cola de manufactura 3D con roles de operador autorizados.

---

### Etapa 6: Frontend PWA, Notificaciones, Pruebas y Despliegue en Producción
* **Objetivo:** Desarrollar la aplicación cliente PWA, integrar notificaciones y poner en marcha la solución completa.
* **Tareas Clave:**
  - Definición de UI/UX, componentes visuales e identidad gráfica de SEMARD.
  - Desarrollo de la interfaz en **Next.js / React** consumiendo la API de Go.
  - Implementación de capacidades PWA: Web App Manifest, Service Worker para cacheo offline de datos del Hub e instalabilidad en móviles y escritorio.
  - Centro de notificaciones in-app y alertas web push/correo (avances evaluados, recursos devueltos, préstamos modificados, impresión 3D lista).
  - Pruebas integradas de flujos, validación de permisos RBAC y auditoría de accesibilidad/rendimiento.
  - Despliegue final en la máquina virtual mediante Coolify y documentación técnica de entrega.
* **Entregables:**
  - PWA SEMARD instalable y totalmente conectada al backend en Go.
  - Plataforma desplegada y operativa en producción.
