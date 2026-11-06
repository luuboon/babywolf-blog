-- ============================================
-- 007_app_lectores.sql
-- App BabyWolf del teléfono: lo que los LECTORES hacen con su cuenta.
--
--   · Reportes de noticias (Prácticas 9 y 10): la tabla incidencias de 006
--     pasa a ser el reporte que un lector hace sobre una noticia. Se agrega
--     la noticia, la prioridad y un número de versión para detectar
--     ediciones simultáneas desde dos dispositivos.
--   · El admin del blog ve y atiende TODOS los reportes; un lector sólo los
--     suyos (las reglas de 006 se conservan tal cual).
--   · Comentarios (Práctica 6): la tabla comments del esquema inicial, que la
--     web no usa, con sus reglas RLS y un límite de 280 caracteres.
--
-- Ejecutar UNA vez en Supabase Dashboard -> SQL Editor, después de 006.
-- Es idempotente. No cambia nada que use la web.
-- ============================================

-- ── ¿El usuario de la sesión es admin del blog? ─────────────
-- SECURITY DEFINER: la regla consulta public.users sin depender de los
-- permisos del que llama. search_path fijo para que no se pueda suplantar
-- la tabla con otra del mismo nombre.
CREATE OR REPLACE FUNCTION public.es_admin()
RETURNS boolean
LANGUAGE sql STABLE SECURITY DEFINER
SET search_path = public
AS $$
  SELECT EXISTS (
    SELECT 1 FROM public.users
    WHERE id = auth.uid() AND role = 'admin'
  );
$$;
REVOKE ALL ON FUNCTION public.es_admin() FROM PUBLIC, anon;
GRANT EXECUTE ON FUNCTION public.es_admin() TO authenticated;

-- ── Reportes de noticias ────────────────────────────────────
ALTER TABLE public.incidencias
  ADD COLUMN IF NOT EXISTS post_slug   VARCHAR(255),
  ADD COLUMN IF NOT EXISTS post_titulo VARCHAR(255),
  ADD COLUMN IF NOT EXISTS prioridad   VARCHAR(10) NOT NULL DEFAULT 'media',
  ADD COLUMN IF NOT EXISTS version     INTEGER     NOT NULL DEFAULT 1,
  ADD COLUMN IF NOT EXISTS actualizado_por UUID;

ALTER TABLE public.incidencias DROP CONSTRAINT IF EXISTS incidencias_prioridad_check;
ALTER TABLE public.incidencias ADD CONSTRAINT incidencias_prioridad_check
  CHECK (prioridad IN ('baja', 'media', 'alta', 'critica'));

CREATE INDEX IF NOT EXISTS idx_incidencias_post_slug ON public.incidencias(post_slug);

-- Cada cambio sube la versión y deja quién lo hizo. La app guarda con
-- "?version=eq.N": si otro dispositivo guardó antes, ya no hay fila con esa
-- versión, no se cambia nada y la app avisa del conflicto (Práctica 9, DB-07).
CREATE OR REPLACE FUNCTION public.incidencias_versionar()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
  -- Un lector ya no puede cambiar el dueño (lo rechaza la regla WITH CHECK
  -- de 006, SEG-09). La regla del admin no ve la fila anterior, así que
  -- para él se bloquea aquí: atiende reportes, no los reasigna a otro lector.
  IF NEW.creado_por IS DISTINCT FROM OLD.creado_por AND public.es_admin() THEN
    RAISE EXCEPTION 'El dueño de un reporte no se puede cambiar' USING ERRCODE = '42501';
  END IF;
  NEW.version := OLD.version + 1;
  NEW.actualizado_por := auth.uid();
  RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS incidencias_versionar ON public.incidencias;
CREATE TRIGGER incidencias_versionar
  BEFORE UPDATE ON public.incidencias
  FOR EACH ROW EXECUTE PROCEDURE public.incidencias_versionar();

-- El admin lee y atiende todos los reportes. Se suman a las reglas de 006
-- (las políticas permisivas se combinan con OR): un lector sigue viendo
-- sólo lo suyo y sigue sin poder pasar un reporte a otra persona.
DROP POLICY IF EXISTS incidencias_admin_leer ON public.incidencias;
CREATE POLICY incidencias_admin_leer ON public.incidencias
  FOR SELECT TO authenticated
  USING ((select public.es_admin()));

DROP POLICY IF EXISTS incidencias_admin_atender ON public.incidencias;
CREATE POLICY incidencias_admin_atender ON public.incidencias
  FOR UPDATE TO authenticated
  USING ((select public.es_admin()))
  WITH CHECK ((select public.es_admin()));

-- ── Comentarios ─────────────────────────────────────────────
-- users.id = auth.users.id (lo asegura el trigger handle_new_user), así que
-- author_id puede tomarse del JWT.
ALTER TABLE public.comments ALTER COLUMN author_id SET DEFAULT auth.uid();

ALTER TABLE public.comments DROP CONSTRAINT IF EXISTS comments_content_check;
ALTER TABLE public.comments ADD CONSTRAINT comments_content_check
  CHECK (char_length(btrim(content)) BETWEEN 1 AND 280) NOT VALID;

GRANT SELECT ON public.comments TO anon;
GRANT SELECT, INSERT, UPDATE, DELETE ON public.comments TO authenticated;
ALTER TABLE public.comments ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS "Comments are viewable by everyone" ON public.comments;
DROP POLICY IF EXISTS comments_leer ON public.comments;
CREATE POLICY comments_leer ON public.comments
  FOR SELECT USING (is_approved);

DROP POLICY IF EXISTS "Authenticated users can comment" ON public.comments;
DROP POLICY IF EXISTS comments_crear ON public.comments;
CREATE POLICY comments_crear ON public.comments
  FOR INSERT TO authenticated
  WITH CHECK (author_id = (select auth.uid()));

DROP POLICY IF EXISTS "Users can update own comments" ON public.comments;
DROP POLICY IF EXISTS comments_editar ON public.comments;
CREATE POLICY comments_editar ON public.comments
  FOR UPDATE TO authenticated
  USING (author_id = (select auth.uid()))
  WITH CHECK (author_id = (select auth.uid()));

-- Borra el autor, o el admin para moderar.
DROP POLICY IF EXISTS "Users can delete own comments" ON public.comments;
DROP POLICY IF EXISTS comments_borrar ON public.comments;
CREATE POLICY comments_borrar ON public.comments
  FOR DELETE TO authenticated
  USING (author_id = (select auth.uid()) OR (select public.es_admin()));
