CREATE SCHEMA IF NOT EXISTS app;

-- Create sequence for invoice suffix (daily reset)
CREATE SEQUENCE IF NOT EXISTS app.invoice_suffix_seq
    START WITH 1
    INCREMENT BY 1
    NO MAXVALUE
    CACHE 1;

CREATE TABLE IF NOT EXISTS app.invoices (
    id BIGINT PRIMARY KEY,
    issuer TEXT NOT NULL,
    customer TEXT NOT NULL,
    issue_date TEXT NOT NULL,
    due_date TIMESTAMP WITH TIME ZONE,
    note TEXT NOT NULL,
    author_id INT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS app.invoice_items (
    id SERIAL PRIMARY KEY,
    invoice_id BIGINT NOT NULL REFERENCES app.invoices(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    qty INTEGER NOT NULL DEFAULT 1,
    price INTEGER NOT NULL,
    total INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS app.invoice_user_daily_seq (
    user_id BIGINT NOT NULL,
    day DATE NOT NULL,
    counter INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, day)
);


CREATE INDEX idx_invoices_deleted_at ON app.invoices(deleted_at);
CREATE INDEX idx_invoice_items_invoice_id ON app.invoice_items(invoice_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_invoice_items_deleted_at ON app.invoice_items(deleted_at);
