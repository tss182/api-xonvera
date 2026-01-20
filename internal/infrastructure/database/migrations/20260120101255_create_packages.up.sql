CREATE SCHEMA IF NOT EXISTS app;
CREATE TABLE IF NOT EXISTS app.packages (
    id serial PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    price INTEGER NOT NULL,
    discount_type VARCHAR(50) NOT NULL DEFAULT 'percentage' CHECK (discount_type IN ('percentage', 'amount')),
    discount INTEGER NOT NULL DEFAULT 0,
    duration VARCHAR(10) NOT NULL DEFAULT '1d',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_packages_deleted_at ON app.packages(deleted_at);
