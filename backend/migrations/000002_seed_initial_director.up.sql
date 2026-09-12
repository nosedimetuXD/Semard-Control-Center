-- SEMARD Control Center: Seed Director & Event Registrations
-- Migración 000002

-- 1. Tabla de Registro de Asistentes a Eventos Públicos
CREATE TABLE IF NOT EXISTS event_registrations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    institution VARCHAR(255),
    registered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(event_id, email)
);

CREATE INDEX IF NOT EXISTS idx_event_registrations_event ON event_registrations(event_id);

-- 2. Seed Primer Director: Daniel David Franco Benítez
INSERT INTO users (
    google_id, 
    email, 
    student_code, 
    full_name, 
    role, 
    can_operate_3d, 
    avatar_url, 
    bio, 
    is_active
) VALUES (
    '112422881009265436802',
    'dfrancob1@unicartagena.edu.co',
    '0222410043',
    'DANIEL DAVID FRANCO BENITEZ',
    'DIRECTOR',
    TRUE,
    'https://lh3.googleusercontent.com/a/ACg8ocKT3Q0QsBrRLREPPmel4Cg_6qcgttDKtV0CD0V7tEraC96F6g=s96-c',
    'Director e Investigador Principal - Semillero SEMARD, Universidad de Cartagena',
    TRUE
)
ON CONFLICT (email) DO UPDATE SET 
    role = 'DIRECTOR',
    google_id = EXCLUDED.google_id,
    student_code = EXCLUDED.student_code,
    full_name = EXCLUDED.full_name,
    can_operate_3d = TRUE,
    is_active = TRUE;
