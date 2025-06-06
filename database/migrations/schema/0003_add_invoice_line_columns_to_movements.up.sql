-- Filename: 0003_add_invoice_line_columns_to_movements.up.sql
-- Description: Adds invoice line related columns to the movements table

-- Add columns for supporting InvoiceLine representation in the movements table
ALTER TABLE movements 
ADD COLUMN amount_without_tax DECIMAL(10, 2),
ADD COLUMN amount_with_tax DECIMAL(10, 2),
ADD COLUMN tax_percentage DECIMAL(5, 2),
ADD COLUMN operation_type VARCHAR(50);
