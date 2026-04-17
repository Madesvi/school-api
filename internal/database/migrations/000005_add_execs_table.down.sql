-- 1. Drop indexes first
DROP INDEX IF EXISTS idx_execs_username;
DROP INDEX IF EXISTS idx_execs_email;

-- 2. Drop the table
DROP TABLE IF EXISTS execs;
