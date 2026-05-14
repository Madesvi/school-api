-- migrate:disable-transaction
DROP INDEX CONCURRENTLY IF EXISTS idx_teachers_email;

-- IMPORTANT:
-- --migrate:disable-transaction for CONCURRENTLY
-- Use only one command at a time with disable_transaction