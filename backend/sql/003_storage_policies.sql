-- ============================================
-- Políticas de Storage para bucket: blog_assets
-- Ejecutar en Supabase Dashboard → SQL Editor
-- ============================================

-- Permitir que usuarios autenticados suban archivos
CREATE POLICY "Allow authenticated uploads"
  ON storage.objects FOR INSERT
  TO authenticated
  WITH CHECK (bucket_id = 'blog_assets');

-- Permitir que cualquier persona lea archivos (público)
CREATE POLICY "Allow public reads"
  ON storage.objects FOR SELECT
  TO anon, authenticated
  USING (bucket_id = 'blog_assets');

-- Permitir que usuarios autenticados actualicen archivos
CREATE POLICY "Allow authenticated updates"
  ON storage.objects FOR UPDATE
  TO authenticated
  USING (bucket_id = 'blog_assets');

-- Permitir que usuarios autenticados eliminen archivos
CREATE POLICY "Allow authenticated deletes"
  ON storage.objects FOR DELETE
  TO authenticated
  USING (bucket_id = 'blog_assets');
