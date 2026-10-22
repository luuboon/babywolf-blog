-- ============================================
-- 005_soporte.sql
-- Mesa de ayuda BabyWolf (Práctica 7 de Desarrollo móvil, app del ecosistema).
-- Ejecutar UNA vez en Supabase Dashboard -> SQL Editor.
--
-- No modifica ninguna tabla existente del blog: users, posts, comments y
-- audit_log quedan igual, así que la web sigue funcionando sin cambios.
-- ============================================

-- Participantes: los 4 dispositivos (admin, agente y usuarios).
CREATE TABLE IF NOT EXISTS public.soporte_participantes (
  id UUID PRIMARY KEY,
  nombre VARCHAR(40) NOT NULL CHECK (length(trim(nombre)) > 0),
  rol VARCHAR(10) NOT NULL CHECK (rol IN ('admin', 'agente', 'usuario')),
  dispositivo VARCHAR(40) NOT NULL DEFAULT '',
  -- Sólo el admin tiene cuenta del blog; agente y usuarios entran por URL preasignada.
  user_id UUID REFERENCES public.users(id) ON DELETE SET NULL,
  activo BOOLEAN NOT NULL DEFAULT true,
  ultimo_latido TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Invitaciones de un solo uso (la URL preasignada de rol apunta a una de estas).
CREATE TABLE IF NOT EXISTS public.soporte_invitaciones (
  id VARCHAR(24) PRIMARY KEY,
  rol VARCHAR(10) NOT NULL CHECK (rol IN ('agente', 'usuario')),
  nombre VARCHAR(40) NOT NULL,
  expira_en TIMESTAMPTZ NOT NULL,
  usada_por UUID REFERENCES public.soporte_participantes(id) ON DELETE SET NULL,
  creada_por UUID NOT NULL REFERENCES public.soporte_participantes(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Reportes (fallas, procesos de atención, accesos perdidos).
CREATE TABLE IF NOT EXISTS public.soporte_reportes (
  id VARCHAR(12) PRIMARY KEY CHECK (id ~ '^REP[0-9]{3,6}$'),
  equipo VARCHAR(80) NOT NULL CHECK (length(trim(equipo)) > 0),
  descripcion VARCHAR(500) NOT NULL CHECK (length(trim(descripcion)) > 0),
  tipo VARCHAR(10) NOT NULL DEFAULT 'falla' CHECK (tipo IN ('falla', 'atencion', 'acceso')),
  prioridad VARCHAR(5) NOT NULL DEFAULT 'Media' CHECK (prioridad IN ('Baja', 'Media', 'Alta')),
  estado VARCHAR(12) NOT NULL DEFAULT 'abierto' CHECK (estado IN ('abierto', 'asignado', 'en_proceso', 'resuelto')),
  solicitante_id UUID NOT NULL REFERENCES public.soporte_participantes(id),
  agente_id UUID REFERENCES public.soporte_participantes(id),
  directo_admin BOOLEAN NOT NULL DEFAULT false,
  activo BOOLEAN NOT NULL DEFAULT true, -- baja lógica (TC05)
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Reutiliza la función del esquema inicial (001_initial_schema.sql).
DROP TRIGGER IF EXISTS update_soporte_reportes_modtime ON public.soporte_reportes;
CREATE TRIGGER update_soporte_reportes_modtime
  BEFORE UPDATE ON public.soporte_reportes FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- Mensajes entre dispositivos (con adjunto opcional en Storage).
CREATE TABLE IF NOT EXISTS public.soporte_mensajes (
  id BIGSERIAL PRIMARY KEY,
  reporte_id VARCHAR(12) REFERENCES public.soporte_reportes(id) ON DELETE SET NULL,
  emisor_id UUID NOT NULL REFERENCES public.soporte_participantes(id) ON DELETE CASCADE,
  emisor_nombre VARCHAR(40) NOT NULL,
  emisor_rol VARCHAR(15) NOT NULL,
  para_id UUID REFERENCES public.soporte_participantes(id) ON DELETE SET NULL,
  texto VARCHAR(500) NOT NULL DEFAULT '',
  archivo_path TEXT,
  archivo_nombre VARCHAR(120),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_soporte_reportes_activo ON public.soporte_reportes(activo, updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_soporte_mensajes_para ON public.soporte_mensajes(para_id);
CREATE INDEX IF NOT EXISTS idx_soporte_mensajes_archivo ON public.soporte_mensajes(archivo_path);

-- RLS activado SIN políticas: la API pública de Supabase (anon/authenticated)
-- no puede leer ni escribir estas tablas. Sólo el backend Go, que se conecta
-- como dueño de las tablas, accede a ellas y aplica las reglas de roles.
ALTER TABLE public.soporte_participantes ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.soporte_invitaciones ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.soporte_reportes ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.soporte_mensajes ENABLE ROW LEVEL SECURITY;

-- Bucket PRIVADO para los adjuntos: sólo se accede con URL preasignada
-- (firmada por el backend con la service role key). Máximo 5 MB.
INSERT INTO storage.buckets (id, name, public, file_size_limit, allowed_mime_types)
VALUES ('soporte_archivos', 'soporte_archivos', false, 5242880,
        ARRAY['image/png', 'image/jpeg', 'application/pdf', 'text/plain'])
ON CONFLICT (id) DO NOTHING;
