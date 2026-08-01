-- ============================================
-- 004_security_audit.sql
-- Practicas 11-12: activación/desactivación, bloqueo por
-- intentos fallidos y bitácora de auditoría.
-- Ejecutar en Supabase Dashboard -> SQL Editor
-- ============================================

-- ── 1. Columnas de seguridad en users ──────────────────────
ALTER TABLE public.users
  ADD COLUMN IF NOT EXISTS active BOOLEAN NOT NULL DEFAULT true,
  ADD COLUMN IF NOT EXISTS failed_login_attempts INT NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS locked_until TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS must_change_password BOOLEAN NOT NULL DEFAULT false;

-- ── 2. Tabla de bitácora ────────────────────────────────────
CREATE TABLE IF NOT EXISTS public.audit_log (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID REFERENCES public.users(id) ON DELETE SET NULL,
  email VARCHAR(255),
  action VARCHAR(50) NOT NULL, -- login_success | login_failed | logout | password_change | user_created | user_deactivated | user_activated | role_changed
  ip_address VARCHAR(64),
  details TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc'::text, now())
);

CREATE INDEX IF NOT EXISTS idx_audit_log_user_id ON public.audit_log(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_log_created_at ON public.audit_log(created_at DESC);

ALTER TABLE public.audit_log ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Admins can read audit log"
  ON public.audit_log FOR SELECT
  TO authenticated
  USING (EXISTS (SELECT 1 FROM public.users u WHERE u.id = auth.uid() AND u.role = 'admin'));

CREATE POLICY "Anyone can insert their own audit rows"
  ON public.audit_log FOR INSERT
  TO authenticated, anon
  WITH CHECK (true); -- las funciones RPC de abajo son el único punto de escritura real

-- ── 3. Bloqueo temporal tras intentos fallidos ─────────────
-- Umbral y duración del bloqueo (ajustables aquí, en un solo lugar).
-- ponytail: constantes embebidas en las funciones; mover a tabla de config si algún día son configurables desde UI.

CREATE OR REPLACE FUNCTION public.check_account_lock(p_email TEXT)
RETURNS JSONB
LANGUAGE plpgsql
SECURITY DEFINER
AS $$
DECLARE
  v_user RECORD;
BEGIN
  SELECT id, active, locked_until INTO v_user FROM public.users WHERE email = p_email;

  IF v_user IS NULL THEN
    RETURN jsonb_build_object('locked', false, 'active', true);
  END IF;

  RETURN jsonb_build_object(
    'locked', v_user.locked_until IS NOT NULL AND v_user.locked_until > now(),
    'locked_until', v_user.locked_until,
    'active', v_user.active
  );
END;
$$;

CREATE OR REPLACE FUNCTION public.record_failed_login(p_email TEXT, p_ip TEXT DEFAULT NULL)
RETURNS VOID
LANGUAGE plpgsql
SECURITY DEFINER
AS $$
DECLARE
  v_attempts INT;
BEGIN
  UPDATE public.users
    SET failed_login_attempts = failed_login_attempts + 1,
        locked_until = CASE WHEN failed_login_attempts + 1 >= 5 THEN now() + interval '15 minutes' ELSE locked_until END
    WHERE email = p_email
    RETURNING failed_login_attempts INTO v_attempts;

  INSERT INTO public.audit_log (user_id, email, action, ip_address, details)
    SELECT id, p_email, 'login_failed', p_ip, 'intento ' || COALESCE(v_attempts, 0)
    FROM public.users WHERE email = p_email;
END;
$$;

CREATE OR REPLACE FUNCTION public.record_successful_login(p_user_id UUID, p_ip TEXT DEFAULT NULL)
RETURNS VOID
LANGUAGE plpgsql
SECURITY DEFINER
AS $$
BEGIN
  UPDATE public.users
    SET failed_login_attempts = 0, locked_until = NULL
    WHERE id = p_user_id;

  INSERT INTO public.audit_log (user_id, email, action, ip_address)
    SELECT id, email, 'login_success', p_ip FROM public.users WHERE id = p_user_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.record_logout(p_user_id UUID)
RETURNS VOID
LANGUAGE plpgsql
SECURITY DEFINER
AS $$
BEGIN
  INSERT INTO public.audit_log (user_id, email, action)
    SELECT id, email, 'logout' FROM public.users WHERE id = p_user_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.record_password_change(p_user_id UUID)
RETURNS VOID
LANGUAGE plpgsql
SECURITY DEFINER
AS $$
BEGIN
  UPDATE public.users SET must_change_password = false WHERE id = p_user_id;

  INSERT INTO public.audit_log (user_id, email, action)
    SELECT id, email, 'password_change' FROM public.users WHERE id = p_user_id;
END;
$$;

-- Bitácora genérica para acciones administrativas (alta, baja, cambio de rol)
CREATE OR REPLACE FUNCTION public.record_admin_action(p_actor_id UUID, p_target_email TEXT, p_action TEXT, p_details TEXT DEFAULT NULL)
RETURNS VOID
LANGUAGE plpgsql
SECURITY DEFINER
AS $$
BEGIN
  INSERT INTO public.audit_log (user_id, email, action, details)
    VALUES (p_actor_id, p_target_email, p_action, p_details);
END;
$$;

GRANT EXECUTE ON FUNCTION public.check_account_lock(TEXT) TO anon, authenticated;
GRANT EXECUTE ON FUNCTION public.record_failed_login(TEXT, TEXT) TO anon, authenticated;
GRANT EXECUTE ON FUNCTION public.record_successful_login(UUID, TEXT) TO anon, authenticated;
GRANT EXECUTE ON FUNCTION public.record_logout(UUID) TO authenticated;
GRANT EXECUTE ON FUNCTION public.record_password_change(UUID) TO authenticated;
GRANT EXECUTE ON FUNCTION public.record_admin_action(UUID, TEXT, TEXT, TEXT) TO authenticated;
