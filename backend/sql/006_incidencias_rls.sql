-- ============================================
-- 006_incidencias_rls.sql
-- Práctica 10 (Desarrollo móvil): seguridad, autenticación y control de
-- acceso a una base de datos compartida entre dispositivos.
--
-- Equivalente en Supabase de las "reglas de seguridad de Firestore": la app
-- habla DIRECTO con la base (API REST de Supabase) usando el JWT del usuario,
-- y Postgres decide con Row Level Security (RLS) qué filas puede ver o tocar.
-- Una operación no autorizada la rechaza la base, aunque la app se modifique.
--
-- Ejecutar UNA vez en Supabase Dashboard -> SQL Editor. Es idempotente.
-- No toca ninguna tabla del blog.
-- ============================================

CREATE TABLE IF NOT EXISTS public.incidencias (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  titulo      VARCHAR(80)  NOT NULL CHECK (length(trim(titulo)) > 0),
  descripcion VARCHAR(500) NOT NULL DEFAULT '',
  estado      VARCHAR(12)  NOT NULL DEFAULT 'abierta'
              CHECK (estado IN ('abierta', 'en_proceso', 'cerrada')),
  -- Propietario = UID de Supabase Auth (auth.users.id). La app lo manda
  -- explícito (Fase 2); si no lo manda, se toma del JWT.
  creado_por  UUID NOT NULL DEFAULT auth.uid() REFERENCES auth.users(id) ON DELETE CASCADE,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_incidencias_creado_por ON public.incidencias(creado_por);

-- updated_at automático (función creada en 001_initial_schema.sql).
DROP TRIGGER IF EXISTS update_incidencias_modtime ON public.incidencias;
CREATE TRIGGER update_incidencias_modtime
  BEFORE UPDATE ON public.incidencias FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- ── Permisos de tabla ────────────────────────────────────────
-- Sin sesión (rol anon) no hay NINGÚN permiso: la base responde
-- "permission denied" en vez de devolver una lista vacía (SEG-10).
REVOKE ALL ON public.incidencias FROM anon;
GRANT SELECT, INSERT, UPDATE, DELETE ON public.incidencias TO authenticated;

-- ── Reglas por fila (RLS) ────────────────────────────────────
ALTER TABLE public.incidencias ENABLE ROW LEVEL SECURITY;

-- (select auth.uid()) se evalúa una vez por consulta y no por fila
-- (recomendación de rendimiento de Supabase).

DROP POLICY IF EXISTS incidencias_leer_propias ON public.incidencias;
CREATE POLICY incidencias_leer_propias ON public.incidencias
  FOR SELECT TO authenticated
  USING (creado_por = (select auth.uid()));                 -- SEG-05

DROP POLICY IF EXISTS incidencias_crear_propias ON public.incidencias;
CREATE POLICY incidencias_crear_propias ON public.incidencias
  FOR INSERT TO authenticated
  WITH CHECK (creado_por = (select auth.uid()));            -- SEG-03, SEG-04

DROP POLICY IF EXISTS incidencias_editar_propias ON public.incidencias;
CREATE POLICY incidencias_editar_propias ON public.incidencias
  FOR UPDATE TO authenticated
  USING (creado_por = (select auth.uid()))                  -- SEG-06, SEG-08
  WITH CHECK (creado_por = (select auth.uid()));            -- SEG-09

DROP POLICY IF EXISTS incidencias_borrar_propias ON public.incidencias;
CREATE POLICY incidencias_borrar_propias ON public.incidencias
  FOR DELETE TO authenticated
  USING (creado_por = (select auth.uid()));                 -- SEG-07
