-- 000002_add_unique_name_index.up.sql
-- Supprimer l'index classique devenu redondant
DROP INDEX IF EXISTS idx_users_name;

-- Créer l'index unique insensible à la casse (pour name non vide)
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_name_unique ON users (LOWER(name)) WHERE name != '';
