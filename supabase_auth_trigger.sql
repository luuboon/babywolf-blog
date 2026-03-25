-- ==============================================================================
-- 4. TRIGGER PARA CREAR AUTOMÁTICAMENTE EL REGISTRO EN PUBLIC.USERS
-- ==============================================================================

-- Función que se ejecuta cada vez que se crea un usuario en auth.users
CREATE OR REPLACE FUNCTION public.handle_new_user()
RETURNS trigger AS $$
BEGIN
  INSERT INTO public.users (id, supabase_uid, email, username, role)
  VALUES (
    new.id, 
    new.id, 
    new.email, 
    COALESCE(new.raw_user_meta_data->>'username', split_part(new.email, '@', 1)), 
    'user'
  );
  RETURN new;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Disparador (Trigger) que vigila la tabla auth.users
DROP TRIGGER IF EXISTS on_auth_user_created ON auth.users;
CREATE TRIGGER on_auth_user_created
  AFTER INSERT ON auth.users
  FOR EACH ROW EXECUTE PROCEDURE public.handle_new_user();
