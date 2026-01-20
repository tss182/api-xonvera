-- Rollback indexes
-- Migration: 000003_add_user_indexes.down.sql

DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_phone;
DROP INDEX IF EXISTS idx_users_email;
