-- Práctica 12 · Optimización de GET /api/posts.
-- La consulta es: WHERE published = true ORDER BY created_at DESC LIMIT n.
-- Un índice parcial ordenado evita escanear y ordenar toda la tabla.
-- Idempotente y sin riesgo para los datos: sólo agrega un índice.
CREATE INDEX IF NOT EXISTS idx_posts_publicados_recientes
    ON public.posts (created_at DESC, id)
    WHERE published;
