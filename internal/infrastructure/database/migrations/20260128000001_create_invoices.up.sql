CREATE SCHEMA IF NOT EXISTS billing;

-- Create sequence for invoice suffix (daily reset)
CREATE SEQUENCE IF NOT EXISTS billing.invoice_suffix_seq
    START WITH 1
    INCREMENT BY 1
    NO MAXVALUE
    CACHE 1;

CREATE TABLE IF NOT EXISTS billing.invoices (
    id BIGINT PRIMARY KEY,
    add_to TEXT NOT NULL,
    invoice_for TEXT NOT NULL,
    invoice_from TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS billing.invoice_items (
    id SERIAL PRIMARY KEY,
    invoice_id BIGINT NOT NULL REFERENCES billing.invoices(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price INTEGER NOT NULL,
    total INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_invoices_deleted_at ON billing.invoices(deleted_at);
CREATE INDEX idx_invoice_items_invoice_id ON billing.invoice_items(invoice_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_invoice_items_deleted_at ON billing.invoice_items(deleted_at);
