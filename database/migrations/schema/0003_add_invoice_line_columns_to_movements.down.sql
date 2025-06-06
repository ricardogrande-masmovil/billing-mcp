-- Filename: 0003_add_invoice_line_columns_to_movements.down.sql
-- Description: Removes invoice line related columns from the movements table

-- Remove columns that support InvoiceLine representation
ALTER TABLE movements 
DROP COLUMN IF EXISTS amount_without_tax,
DROP COLUMN IF EXISTS amount_with_tax,
DROP COLUMN IF EXISTS tax_percentage,
DROP COLUMN IF EXISTS operation_type;
