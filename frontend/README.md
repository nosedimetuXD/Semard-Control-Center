# SEMARD Control Center — Frontend (Next.js 15 + Vercel)

Frontend oficial y estudio de pruebas interactivas de **SEMARD Control Center** desarrollado en **Next.js 15 (App Router)**, **React 19**, **TypeScript** y **Tailwind CSS**.

---

## 🚀 Despliegue en Vercel en 3 Pasos

1. **Importar el Repositorio en Vercel:**
   - Ve a [vercel.com](https://vercel.com) e inicia sesión con tu cuenta de GitHub.
   - Haz clic en **"Add New..." ➔ "Project"** y selecciona `Semard-Control-Center`.

2. **Configurar el Directorio Raíz (Root Directory):**
   - En **Root Directory**, haz clic en **Edit** y selecciona la carpeta **`frontend`**.
   - Vercel detectará automáticamente el framework como **Next.js**.

3. **Configurar la Variable de Entorno:**
   - Despliega la sección **Environment Variables** y añade:
     - **Key:** `NEXT_PUBLIC_API_URL`
     - **Value:** `https://semardcontrolcenter.147.5.103.87.sslip.io`
   - Haz clic en **Deploy**.

---

## 💻 Desarrollo Local

Para correr el frontend localmente en tu máquina:

```powershell
# 1. Entrar a la carpeta frontend
cd frontend

# 2. Instalar dependencias
npm install

# 3. Iniciar servidor de desarrollo
npm run dev
```

Abre en tu navegador: [http://localhost:3000](http://localhost:3000)

---

## 🔑 Autenticación con Google OAuth

Cuando haces clic en **Google Login** desde tu despliegue en Vercel (o desde `localhost:3000`), el frontend envía automáticamente su origen en el parámetro `return_to`:
- `GET /api/v1/auth/google/login?return_to=https://tu-proyecto.vercel.app/auth/callback&redirect=true`
- Al iniciar sesión con tu cuenta institucional `@unicartagena.edu.co`, el backend en Go redirige a tu app en Vercel preservando el token JWT.
- La página `/auth/callback` almacena el token en `localStorage` y activa tu sesión inmediatamente.

---

## 📦 Módulos Incluidos

1. **Hub Informativo:** `/hub/info`, `/hub/directors`, `/hub/projects`, `/healthz`.
2. **Cartelera de Eventos:** Consulta pública e interna, creación de eventos y registro de asistentes.
3. **Solicitud de Ingreso Estudiantil:** `/auth/register-request`.
4. **Mi Perfil:** `/auth/me` con actualización de biografía y código estudiantil.
5. **Miembros & Roles:** Directorio, aprobación/rechazo de solicitudes, cambio de roles y permisos 3D.
6. **Proyectos:** Vitrina, proyectos personales, alta de proyectos y asignación de líderes.
7. **Avances de Proyecto:** Radicación de avances con evidencias y revisión formativa.
8. **Recursos:** Solicitudes de compra/insumos y evaluación de directores en 3 estados.
9. **Inventario:** Catálogo y stock de herramientas.
10. **Préstamos:** Solicitud, revisión con fecha modificada y devolución física con reposición de stock.
11. **Taller de Impresión 3D:** Subida multipart/form-data de `.stl`/`.obj`/`.3mf`/`.step`, cola de fabricación y transiciones de estado por el operario.
