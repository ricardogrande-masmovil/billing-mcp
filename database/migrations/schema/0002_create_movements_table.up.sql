-- Filename: 0002_create_movements_table.up.sql
-- Description: Creates the movements table to store invoice movement data.

CREATE TABLE IF NOT EXISTS movements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    invoice_id UUID NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    movement_type VARCHAR(50) NOT NULL,
    description TEXT,
    transaction_date TIMESTAMPTZ NOT NULL,
    status VARCHAR(50) NOT NULL,
    -- Fields for invoice line representation
    amount_without_tax DECIMAL(10, 2),
    amount_with_tax DECIMAL(10, 2),
    tax_percentage DECIMAL(5, 2),
    operation_type VARCHAR(50),

    CONSTRAINT fk_movements_invoice_id FOREIGN KEY (invoice_id)
        REFERENCES invoices (id)
        ON DELETE CASCADE
);

-- Create indexes for frequently queried columns
CREATE INDEX IF NOT EXISTS idx_movements_invoice_id ON movements (invoice_id);
CREATE INDEX IF NOT EXISTS idx_movements_status ON movements (status);
CREATE INDEX IF NOT EXISTS idx_movements_deleted_at ON movements (deleted_at);
