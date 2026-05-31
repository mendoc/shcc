-- 000002_add_unique_name_index.down.sql
DROP INDEX IF EXISTS idx_users_name_unique;

-- Recréer l'index classique au cas où on reviendrait en arrière
CREATE INDEX IF NOT EXISTS idx_users_name ON users(name);
