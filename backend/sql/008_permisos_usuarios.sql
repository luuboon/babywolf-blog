-- ============================================
-- 008_permisos_usuarios.sql
-- Práctica 10: cierra huecos de permisos en public.users, la bitácora y
-- las funciones de bloqueo por intentos (004).
--
-- Antes de esto, con la llave pública (anon) y una sesión cualquiera:
--   · un usuario podía ponerse role = 'admin' editando su propia fila;
--   · cualquiera, sin sesión, leía rol, intentos fallidos y bloqueos;
--   · cualquiera podía llamar record_successful_login con el id de OTRA
--     cuenta y reiniciar su contador (anulaba el bloqueo por fuerza bruta);
--   · cualquiera podía escribir registros falsos en audit_log.
--
-- Qué NO cambia (revisado en el código):
--   · La web sólo edita username y avatar_url de su propio perfil
--     (profile.page.ts) y sigue pudiendo.
--   · La web y la app leen id, username, email y avatar_url de los autores
--     sin sesión (posts y comentarios): se conserva.
--   · Con sesión se sigue leyendo todo (perfil propio con select *, panel
--     de admin): ver "Pendiente" abajo.
--   · El backend Go entra como postgres (dueño de las tablas): las altas,
--     bajas y cambios de rol del admin siguen igual.
--   · check_account_lock y record_failed_login siguen abiertas: se usan
--     antes de iniciar sesión.
--
-- Pendiente (necesita cambiar la web primero): que con sesión un lector no
-- pueda leer el correo ni el estado de bloqueo de las demás cuentas, y que
-- sin sesión no se lea el correo.
--
-- Ejecutar UNA vez en Supabase Dashboard -> SQL Editor, después de 007.
-- Es idempotente. Pruebas: backend/pruebas/usuarios_permisos.sh
-- ============================================

-- ── 1. users: qué columnas puede escribir cada quien ───────
-- Sin permiso por columna, la regla "Users can update own profile" dejaba
-- cambiar CUALQUIER columna de la fila propia, incluido role.
REVOKE INSERT, UPDATE, DELETE, TRUNCATE ON public.users FROM anon, authenticated;
GRANT UPDATE (username, avatar_url) ON public.users TO authenticated;

-- La fila la crea el trigger handle_new_user al registrarse; nadie la
-- inserta a mano desde el cliente.
DROP POLICY IF EXISTS "Users can insert their own profile" ON public.users;

DROP POLICY IF EXISTS "Users can update own profile" ON public.users;
CREATE POLICY "Users can update own profile" ON public.users
  FOR UPDATE TO authenticated
  USING (id = (select auth.uid()))
  WITH CHECK (id = (select auth.uid()));

-- Segunda barrera: aunque alguien vuelva a dar permisos de más por error,
-- un cliente (anon / authenticated) no cambia columnas de seguridad. Las
-- funciones SECURITY DEFINER y el backend corren como el dueño y no entran
-- aquí.
CREATE OR REPLACE FUNCTION public.users_proteger_columnas()
RETURNS trigger LANGUAGE plpgsql SET search_path = public AS $$
BEGIN
  IF current_user IN ('anon', 'authenticated') AND (
       NEW.id IS DISTINCT FROM OLD.id
    OR NEW.supabase_uid IS DISTINCT FROM OLD.supabase_uid
    OR NEW.email IS DISTINCT FROM OLD.email
    OR NEW.role IS DISTINCT FROM OLD.role
    OR NEW.active IS DISTINCT FROM OLD.active
    OR NEW.failed_login_attempts IS DISTINCT FROM OLD.failed_login_attempts
    OR NEW.locked_until IS DISTINCT FROM OLD.locked_until
    OR NEW.must_change_password IS DISTINCT FROM OLD.must_change_password
  ) THEN
    RAISE EXCEPTION 'No autorizado: sólo se pueden cambiar username y avatar_url'
      USING ERRCODE = '42501';
  END IF;
  RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS users_proteger_columnas ON public.users;
CREATE TRIGGER users_proteger_columnas
  BEFORE UPDATE ON public.users
  FOR EACH ROW EXECUTE PROCEDURE public.users_proteger_columnas();

-- ── 2. users: qué puede leer quien no ha iniciado sesión ────
-- Sólo lo que se muestra de un autor. Rol, estado y bloqueos, no.
REVOKE SELECT ON public.users FROM anon;
GRANT SELECT (id, username, email, avatar_url) ON public.users TO anon;

-- ── 3. Funciones de 004 ─────────────────────────────────────
-- search_path fijo en todas las SECURITY DEFINER.
ALTER FUNCTION public.check_account_lock(TEXT) SET search_path = public;
ALTER FUNCTION public.record_failed_login(TEXT, TEXT) SET search_path = public;

-- Sólo la propia cuenta, con su sesión, registra su inicio, su salida o su
-- cambio de contraseña. Antes, sin sesión y con el id de otra cuenta, se
-- reiniciaba su contador de intentos.
CREATE OR REPLACE FUNCTION public.record_successful_login(p_user_id UUID, p_ip TEXT DEFAULT NULL)
RETURNS VOID LANGUAGE plpgsql SECURITY DEFINER SET search_path = public AS $$
BEGIN
  IF auth.uid() IS NULL OR auth.uid() <> p_user_id THEN
    RAISE EXCEPTION 'No autorizado' USING ERRCODE = '42501';
  END IF;

  UPDATE public.users
    SET failed_login_attempts = 0, locked_until = NULL
    WHERE id = p_user_id;

  INSERT INTO public.audit_log (user_id, email, action, ip_address)
    SELECT id, email, 'login_success', p_ip FROM public.users WHERE id = p_user_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.record_logout(p_user_id UUID)
RETURNS VOID LANGUAGE plpgsql SECURITY DEFINER SET search_path = public AS $$
BEGIN
  IF auth.uid() IS NULL OR auth.uid() <> p_user_id THEN
    RAISE EXCEPTION 'No autorizado' USING ERRCODE = '42501';
  END IF;

  INSERT INTO public.audit_log (user_id, email, action)
    SELECT id, email, 'logout' FROM public.users WHERE id = p_user_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.record_password_change(p_user_id UUID)
RETURNS VOID LANGUAGE plpgsql SECURITY DEFINER SET search_path = public AS $$
BEGIN
  IF auth.uid() IS NULL OR auth.uid() <> p_user_id THEN
    RAISE EXCEPTION 'No autorizado' USING ERRCODE = '42501';
  END IF;

  UPDATE public.users SET must_change_password = false WHERE id = p_user_id;

  INSERT INTO public.audit_log (user_id, email, action)
    SELECT id, email, 'password_change' FROM public.users WHERE id = p_user_id;
END;
$$;

ALTER FUNCTION public.record_admin_action(UUID, TEXT, TEXT, TEXT) SET search_path = public;

-- Postgres da EXECUTE a PUBLIC por omisión: se quita y se da sólo a quien
-- la necesita.
REVOKE ALL ON FUNCTION public.record_successful_login(UUID, TEXT) FROM PUBLIC, anon;
REVOKE ALL ON FUNCTION public.record_logout(UUID) FROM PUBLIC, anon;
REVOKE ALL ON FUNCTION public.record_password_change(UUID) FROM PUBLIC, anon;
GRANT EXECUTE ON FUNCTION public.record_successful_login(UUID, TEXT) TO authenticated;
GRANT EXECUTE ON FUNCTION public.record_logout(UUID) TO authenticated;
GRANT EXECUTE ON FUNCTION public.record_password_change(UUID) TO authenticated;

-- Ningún cliente la usa: el backend escribe la bitácora directamente.
REVOKE ALL ON FUNCTION public.record_admin_action(UUID, TEXT, TEXT, TEXT) FROM PUBLIC, anon, authenticated;

-- handle_new_user (supabase_auth_trigger.sql) crea la fila de users al
-- registrarse. Es función de trigger: no tiene que ser invocable desde la
-- API, y como es SECURITY DEFINER lleva search_path fijo.
DO $$
BEGIN
  IF to_regprocedure('public.handle_new_user()') IS NOT NULL THEN
    ALTER FUNCTION public.handle_new_user() SET search_path = public;
    REVOKE ALL ON FUNCTION public.handle_new_user() FROM PUBLIC, anon, authenticated;
  END IF;
END
$$;

-- ── 4. Bitácora ─────────────────────────────────────────────
-- Sólo escriben las funciones de arriba y el backend (como dueño). El
-- admin la sigue leyendo con "Admins can read audit log".
DROP POLICY IF EXISTS "Anyone can insert their own audit rows" ON public.audit_log;
REVOKE INSERT, UPDATE, DELETE, TRUNCATE ON public.audit_log FROM anon, authenticated;
REVOKE SELECT ON public.audit_log FROM anon;
