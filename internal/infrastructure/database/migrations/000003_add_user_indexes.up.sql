-- Add indexes for frequently queried columns
-- Migration: 000003_add_user_indexes.up.sql

CREATE INDEX idx_users_email ON auth.users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_phone ON auth.users(phone) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_deleted_at ON auth.users(deleted_at);
