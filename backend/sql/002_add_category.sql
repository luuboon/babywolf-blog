-- 002_add_category.sql
-- Add category column to posts table for content categorization

ALTER TABLE public.posts
ADD COLUMN IF NOT EXISTS category VARCHAR(100) DEFAULT 'general';

-- Create index for category filtering
CREATE INDEX IF NOT EXISTS idx_posts_category ON public.posts(category);
