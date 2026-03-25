-- ==============================================================================
-- 1. ESTRUCTURAS DE TABLAS (Asegurando schema)
-- ==============================================================================

-- Si las tablas no existen, créalas. (Ajustarlas si ya existen)
CREATE TABLE IF NOT EXISTS public.users (
  id uuid references auth.users not null primary key,
  email text,
  username text,
  role text default 'user',
  created_at timestamp with time zone default timezone('utc'::text, now()) not null
);

CREATE TABLE IF NOT EXISTS public.posts (
  id uuid default gen_random_uuid() primary key,
  author_id uuid references public.users(id) not null,
  title text not null,
  slug text not null unique,
  content text not null,
  category text default 'general',
  cover_image_url text,
  published boolean default false,
  created_at timestamp with time zone default timezone('utc'::text, now()) null,
  updated_at timestamp with time zone default timezone('utc'::text, now()) null
);

CREATE TABLE IF NOT EXISTS public.comments (
  id uuid default gen_random_uuid() primary key,
  post_id uuid references public.posts(id) on delete cascade not null,
  author_id uuid references public.users(id) not null,
  content text not null,
  created_at timestamp with time zone default timezone('utc'::text, now()) not null
);

-- ==============================================================================
-- 2. FASE D: ROW LEVEL SECURITY (RLS)
-- ==============================================================================

-- Habilitar RLS en todas las tablas
ALTER TABLE public.users ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.posts ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.comments ENABLE ROW LEVEL SECURITY;

-- -------------------------
-- POLÍTICAS PARA: users
-- -------------------------
-- 1. Cualquiera puede ver los usuarios
DROP POLICY IF EXISTS "Public users are viewable by everyone" ON public.users;
CREATE POLICY "Public users are viewable by everyone" 
  ON public.users FOR SELECT USING (true);

-- 2. Los usuarios solo pueden actualizar su propio perfil
DROP POLICY IF EXISTS "Users can update own profile" ON public.users;
CREATE POLICY "Users can update own profile" 
  ON public.users FOR UPDATE USING (auth.uid() = id);

-- 3. Inserción permitida al registrarse (o vía trigger)
DROP POLICY IF EXISTS "Users can insert their own profile" ON public.users;
CREATE POLICY "Users can insert their own profile"
  ON public.users FOR INSERT WITH CHECK (auth.uid() = id);

-- -------------------------
-- POLÍTICAS PARA: posts
-- -------------------------
-- 1. Cualquiera puede ver posts publicados
DROP POLICY IF EXISTS "Public posts are viewable by everyone" ON public.posts;
CREATE POLICY "Public posts are viewable by everyone" 
  ON public.posts FOR SELECT USING (published = true);

-- 2. Un autor puede ver sus propios posts no publicados
DROP POLICY IF EXISTS "Authors can view own unpublished posts" ON public.posts;
CREATE POLICY "Authors can view own unpublished posts" 
  ON public.posts FOR SELECT USING (auth.uid() = author_id);

-- 3. Solo autenticados pueden crear (el rol de admin se podría chequear aquí)
DROP POLICY IF EXISTS "Authenticated users can create posts" ON public.posts;
CREATE POLICY "Authenticated users can create posts" 
  ON public.posts FOR INSERT WITH CHECK (auth.uid() = author_id);

-- 4. El autor puede actualizar sus posts
DROP POLICY IF EXISTS "Authors can update their own posts" ON public.posts;
CREATE POLICY "Authors can update their own posts" 
  ON public.posts FOR UPDATE USING (auth.uid() = author_id);

-- 5. El autor puede borrar sus posts
DROP POLICY IF EXISTS "Authors can delete their own posts" ON public.posts;
CREATE POLICY "Authors can delete their own posts" 
  ON public.posts FOR DELETE USING (auth.uid() = author_id);

-- -------------------------
-- POLÍTICAS PARA: comments
-- -------------------------
-- 1. Cualquiera puede ver los comentarios
DROP POLICY IF EXISTS "Comments are viewable by everyone" ON public.comments;
CREATE POLICY "Comments are viewable by everyone" 
  ON public.comments FOR SELECT USING (true);

-- 2. Autenticados pueden crear comentarios
DROP POLICY IF EXISTS "Authenticated users can comment" ON public.comments;
CREATE POLICY "Authenticated users can comment" 
  ON public.comments FOR INSERT WITH CHECK (auth.uid() = author_id);

-- 3. Autores pueden borrar/editar sus comentarios
DROP POLICY IF EXISTS "Users can update own comments" ON public.comments;
CREATE POLICY "Users can update own comments" 
  ON public.comments FOR UPDATE USING (auth.uid() = author_id);
  
DROP POLICY IF EXISTS "Users can delete own comments" ON public.comments;
CREATE POLICY "Users can delete own comments" 
  ON public.comments FOR DELETE USING (auth.uid() = author_id);

-- ==============================================================================
-- 3. FASE A: SEED DE DATOS DE PRUEBA
-- ==============================================================================
-- INSTRUCCIÓN PARA EL USUARIO:
-- 1. Primero ve a Authentication -> Add User -> Create New User en Supabase.
-- 2. Crea el usuario: abrajam285@gmail.com
-- 3. Copia el User UID generado.
-- 4. Reemplaza 'uuid-del-admin-aqui' en todo el script con tu UID real.

DO $$
DECLARE
    admin_id uuid := 'af551782-2b8b-4b8e-8a5a-9aca63bc816a'; -- ¡Autocompletado por Antigravity!
    post1_id uuid;
    post2_id uuid;
BEGIN
    -- Asegurar el registro en public.users
    INSERT INTO public.users (id, supabase_uid, email, username, role)
    VALUES (admin_id, admin_id, 'abrajam285@gmail.com', 'admin_abrajam', 'admin')
    ON CONFLICT (id) DO UPDATE SET role = 'admin', supabase_uid = EXCLUDED.supabase_uid;

    -- Insertar Posts
    INSERT INTO public.posts (author_id, title, slug, content, category, cover_image_url, published)
    VALUES
    (admin_id, 'El Renacer de los 8-bits', 'renacer-8-bits', 'Los juegos de 8-bits están volviendo con fuerza. Desarrolladores indie están creando joyas...', 'retro', 'https://images.unsplash.com/photo-1550745165-9bc0b252726f?w=800&q=80', true),
    (admin_id, 'Review: The Legend of Zelda Tears of the Kingdom', 'zelda-totk-review', 'Una obra maestra absoluta que redefine la física y la creatividad en los videojuegos.', 'gaming', 'https://images.unsplash.com/photo-1588636545233-1ec91a1a5b82?w=800&q=80', true),
    (admin_id, 'Desarrollo con Angular 17', 'desarrollo-angular-17', 'Las nuevas Signals han cambiado para siempre la forma en que manejamos el estado en Angular...', 'tech', 'https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=800&q=80', true),
    (admin_id, 'El fin del E3: ¿Qué sigue para la industria?', 'fin-del-e3', 'Con la cancelación permanente del E3, el Summer Game Fest toma la corona.', 'gaming', 'https://images.unsplash.com/photo-1542751371-adc38448a05e?w=800&q=80', true),
    (admin_id, 'Supabase vs Firebase: La batalla definitiva', 'supabase-vs-firebase', 'PostgreSQL bajo el capó hace que Supabase sea una alternativa increíblemente potente...', 'tech', 'https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=800&q=80', true),
    (admin_id, 'Coleccionismo de Consolas Clásicas', 'coleccionismo-consolas', '¿Vale la pena invertir en consolas originales o la emulación es suficiente?...', 'retro', 'https://images.unsplash.com/photo-1605333555230-01967269b8fa?w=800&q=80', true)
    ON CONFLICT (slug) DO UPDATE SET title = EXCLUDED.title, content = EXCLUDED.content;

    -- Obtener IDs para comentar
    SELECT id INTO post1_id FROM public.posts WHERE slug = 'renacer-8-bits' LIMIT 1;
    SELECT id INTO post2_id FROM public.posts WHERE slug = 'zelda-totk-review' LIMIT 1;

    -- Insertar Comentarios
    INSERT INTO public.comments (post_id, author_id, content)
    VALUES
    (post1_id, admin_id, '¡Excelente artículo, totalmente de acuerdo!'),
    (post1_id, admin_id, 'Yo sigo prefiriendo jugar en hardware original, la latencia de emulación me molesta.'),
    (post2_id, admin_id, 'Definitivamente el mejor juego del año.');

    RAISE NOTICE '¡Seed y RLS completado con éxito!';
END $$;
