-- migrate:disable-transaction
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_teachers_email;