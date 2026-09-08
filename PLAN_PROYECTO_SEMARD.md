# Plan Maestro y Ejecución del Proyecto: SEMARD PWA

> **Semillero de Investigación SEMARD — Universidad de Cartagena**  
> **Backend:** Go (Golang) | **Base de Datos:** PostgreSQL (Coolify VM) | **Frontend:** PWA (Next.js / React)

---

## 1. Definición Exacta de los 4 Roles del Sistema

El sistema cuenta estrictamente con **4 roles**:

1. **Directores**: Máxima autoridad del semillero. Creación de proyectos, evaluación de avances y recursos, decisiones institucionales y aprobación de nuevos registros.
2. **Administrador**: Gestión logística, inventario, aprobación de préstamos y de solicitudes 3D.
3. **Miembros**: Estudiantes integrantes del semillero con cuenta activa.
   * Pueden ser asignados a proyectos como **encargados**.
   * Pueden solicitar préstamos de insumos y solicitudes de impresión 3D.
   * Aquellos con el permiso especial de **Operador 3D** pueden gestionar y ejecutar impresiones aprobadas.
4. **Público (Visitantes / Comunidad Universitaria)**: Usuarios sin cuenta o sin sesión iniciada.

> [!NOTE]
> **No existe rol de Tutor en el sistema**. Cuando se habla de "encargados de proyectos", se refiere a cualquier usuario con rol de **Miembro**, **Administrador** o **Director** que haya sido asignado a un proyecto específico por un Director.

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
| **Crear y Publicar Eventos (Públicos/Internos)** | ❌ | ❌ | ✅ **Puede gestionar** | ✅ **Puede gestionar** | ❌ |
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
* **1.1 Eventos Abiertos al Público:** Agenda, fechas, ponentes y formulario de registro para ferias, talleres abiertos y conferencias. **Accesible para el Público**.
* **1.2 Información del Semillero & Portafolio Público:**
  * Misión, visión, historia y líneas de investigación activas.
  * **Portafolio Público de Proyectos:** Exhibición para el público general **exclusivamente de los proyectos que se encuentren marcados como finalizados**.
* **1.3 Perfiles de Directores:** Tarjetas con fotografía profesional, biografía resumida, trayectoria investigativa y enlaces a perfiles académicos. **Accesible para el Público**.
* **1.4 Eventos Internos:** Cartelera protegida (requiere inicio de sesión institucional con rol Miembro, Administrador o Director). Cronograma de reuniones semanales, entregas de avances, sustentaciones y actas.

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

## 6. Diseño de Base de Datos (PostgreSQL en Coolify)

```sql
-- Tipos Enumerados
CREATE TYPE user_role AS ENUM ('DIRECTOR', 'ADMINISTRADOR', 'MIEMBRO');
CREATE TYPE registration_status AS ENUM ('PENDING', 'APPROVED', 'REJECTED');
CREATE TYPE project_status AS ENUM ('DRAFT', 'ACTIVE', 'PAUSED', 'COMPLETED');
CREATE TYPE update_status AS ENUM ('PENDING', 'APPROVED', 'CHANGES_REQUESTED');
CREATE TYPE resource_type AS ENUM ('DIGITAL', 'ECONOMIC', 'KNOWLEDGE', 'HARDWARE');
CREATE TYPE resource_status AS ENUM ('PENDING', 'APPROVED', 'REJECTED', 'RETURNED_FOR_MODIFICATION');
CREATE TYPE loan_status AS ENUM ('PENDING', 'APPROVED', 'APPROVED_MODIFIED', 'REJECTED', 'RETURNED', 'OVERDUE');
CREATE TYPE print3d_status AS ENUM ('PENDING', 'APPROVED', 'IN_PROGRESS', 'COMPLETED', 'DELIVERED', 'REJECTED');
CREATE TYPE event_visibility AS ENUM ('PUBLIC', 'INTERNAL');

-- Usuarios
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    google_id VARCHAR(255) UNIQUE,
    email VARCHAR(255) UNIQUE NOT NULL, -- @unicartagena.edu.co
    student_code VARCHAR(50) UNIQUE,
    full_name VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'MIEMBRO',
    can_operate_3d BOOLEAN NOT NULL DEFAULT FALSE,
    avatar_url TEXT,
    bio TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Solicitudes de Registro de Nuevos Miembros
CREATE TABLE registration_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    google_id VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    student_code VARCHAR(50) NOT NULL,
    career_program VARCHAR(255),
    motivation_letter TEXT,
    status registration_status NOT NULL DEFAULT 'PENDING',
    director_feedback TEXT,
    reviewed_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reviewed_at TIMESTAMPTZ
);

-- Proyectos (Creados solo por Directores)
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    research_line VARCHAR(255) NOT NULL,
    status project_status NOT NULL DEFAULT 'ACTIVE',
    created_by UUID NOT NULL REFERENCES users(id),
    start_date DATE NOT NULL,
    target_end_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Encargados y Participantes (Miembros, Admins o Directores)
CREATE TABLE project_members (
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    is_lead BOOLEAN NOT NULL DEFAULT FALSE,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (project_id, user_id)
);

-- Avances de Proyecto (Aprobados por Directores con Feedback)
CREATE TABLE project_updates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    submitted_by UUID NOT NULL REFERENCES users(id), -- Encargado del proyecto
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    attachments_url JSONB DEFAULT '[]'::jsonb,
    status update_status NOT NULL DEFAULT 'PENDING',
    director_feedback TEXT,
    reviewed_by UUID REFERENCES users(id), -- Director que evalúa
    reviewed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Solicitudes de Recursos del Proyecto
CREATE TABLE resource_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    requested_by UUID NOT NULL REFERENCES users(id), -- Encargado del proyecto
    resource_type resource_type NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    estimated_cost NUMERIC(12, 2) DEFAULT 0.00,
    status resource_status NOT NULL DEFAULT 'PENDING',
    director_feedback TEXT, -- Feedback obligatorio si es rechazada o devuelta
    reviewed_by UUID REFERENCES users(id), -- Director que evalúa
    reviewed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Inventario y Préstamos
CREATE TABLE inventory_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    description TEXT,
    total_stock INT NOT NULL DEFAULT 1,
    available_stock INT NOT NULL DEFAULT 1,
    location VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE loan_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID NOT NULL REFERENCES inventory_items(id),
    requested_by UUID NOT NULL REFERENCES users(id),
    reason TEXT NOT NULL,
    requested_start_date DATE NOT NULL,
    requested_end_date DATE NOT NULL,
    approved_end_date DATE,
    status loan_status NOT NULL DEFAULT 'PENDING',
    reviewer_feedback TEXT,
    reviewed_by UUID REFERENCES users(id), -- Director o Administrador
    returned_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Solicitudes de Impresión 3D
CREATE TABLE print3d_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    requested_by UUID NOT NULL REFERENCES users(id),
    project_id UUID REFERENCES projects(id),
    file_path TEXT NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    material VARCHAR(50) NOT NULL,
    color VARCHAR(50),
    infill_percentage INT NOT NULL DEFAULT 20,
    notes TEXT,
    status print3d_status NOT NULL DEFAULT 'PENDING',
    reviewer_feedback TEXT,
    approved_by UUID REFERENCES users(id), -- Exclusivo Director o Administrador
    operator_id UUID REFERENCES users(id), -- Director, Administrador u Operador 3D que ejecuta
    approved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Eventos
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    visibility event_visibility NOT NULL DEFAULT 'PUBLIC',
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    location VARCHAR(255) NOT NULL,
    banner_url TEXT,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## 7. Plan de Ejecución y Hoja de Ruta

```
[ Fase 1: Backend Go & PostgreSQL ] ➔ [ Fase 2: Auth & Usuarios ] ➔ [ Fase 3: Proyectos & Recursos ] ➔ [ Fase 4: Préstamos & 3D ] ➔ [ Fase 5: Frontend PWA ]
```

### Fase 1: Arquitectura Base e Infraestructura (Backend en Go)
* Inicialización del módulo de Go (`go mod init semard-api`).
* Configuración de Dockerfile multi-stage y `docker-compose.yml` para despliegue en Coolify.
* Conexión con PostgreSQL mediante `pgx/v5` y configuración del sistema de migraciones SQL.
* Configuración de variables de entorno y middlewares base (Logger, CORS, Recover).

### Fase 2: Autenticación Híbrida y Gestión de Miembros
* Integración de Google OAuth 2.0 restringido a correos `@unicartagena.edu.co`.
* Lógica de verificación contra el padrón de usuarios (`users`) y generación de tokens JWT.
* Flujo de solicitud de registro para usuarios no listados (`registration_requests`).
* Endpoints para que los Directores listen, aprueben (asignando rol y código) o rechacen solicitudes con observaciones.

### Fase 3: Módulo Logístico de Proyectos
* Endpoints de creación de proyectos (exclusivo para Directores).
* Endpoints de asignación de miembros (Miembros, Administradores o Directores como encargados).
* Endpoints para que los encargados envíen avances con archivos adjuntos.
* Flujo de revisión de avances por parte de Directores (Aprobación / Solicitud de cambios con feedback).
* Endpoints para solicitudes de recursos de proyectos y su ciclo de vida (`Aprobada`, `Rechazada`, `Devuelta para Modificaciones`).

### Fase 4: Inventario, Préstamos de Equipos e Impresión 3D
* CRUD de inventario de herramientas y equipos.
* Solicitudes de préstamo de insumos con flujo de decisión por Directores o Administradores (Aprobar, Modificar plazo con notificación, o Rechazar con feedback).
* Carga segura y almacenamiento de modelos 3D (.STL, .OBJ).
* Flujo de decisión de impresión 3D: Aprobación/Rechazo por Directores y Administradores.
* Cola de ejecución y gestión de estado 3D para Directores, Administradores y Operadores 3D certificados.

### Fase 5: Eventos, Frontend PWA e Identidad Visual
* Definición de identidad visual (paleta de colores, tipografía, logo).
* Construcción del Hub público (eventos, directores, líneas de investigación).
* Construcción de los paneles administrativos y vistas protegidas.
* Configuración de Service Worker PWA, manifiesto web, caché offline y notificaciones push.
