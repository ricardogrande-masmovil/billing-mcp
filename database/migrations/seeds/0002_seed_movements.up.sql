-- Filename: 0002_seed_movements.up.sql
-- Description: Inserts seed data into the movements table

-- Insert movements for existing invoices
-- We're using predefined UUIDs for easier referencing and testing
INSERT INTO movements (
    id, 
    invoice_id, 
    amount, 
    movement_type, 
    description, 
    transaction_date, 
    status, 
    amount_without_tax, 
    amount_with_tax, 
    tax_percentage, 
    operation_type,
    created_at, 
    updated_at
) VALUES
-- Movements for invoice 123e4567-e89b-12d3-a456-426614174001 (SENT)
('233e4567-e89b-12d3-a456-426614174001', '123e4567-e89b-12d3-a456-426614174001', 50.25, 'CREDIT', 'Base service charge', '2025-01-15T10:30:00Z', 'PENDING', 45.68, 50.25, 10.00, 'SERVICE', NOW(), NOW()),
('233e4567-e89b-12d3-a456-426614174002', '123e4567-e89b-12d3-a456-426614174001', 50.25, 'CREDIT', 'Additional features', '2025-01-15T10:35:00Z', 'PENDING', 45.68, 50.25, 10.00, 'FEATURE', NOW(), NOW()),

-- Movements for invoice 123e4567-e89b-12d3-a456-426614174002 (PAID)
('233e4567-e89b-12d3-a456-426614174003', '123e4567-e89b-12d3-a456-426614174002', 150.75, 'CREDIT', 'Monthly subscription', '2025-01-20T12:00:00Z', 'INVOICED', 125.63, 150.75, 20.00, 'SUBSCRIPTION', NOW(), NOW()),
('233e4567-e89b-12d3-a456-426614174004', '123e4567-e89b-12d3-a456-426614174002', 100.00, 'CREDIT', 'Premium support', '2025-01-20T12:05:00Z', 'INVOICED', 83.33, 100.00, 20.00, 'SUPPORT', NOW(), NOW()),

-- Movements for invoice 123e4567-e89b-12d3-a456-426614174003 (OVERDUE)
('233e4567-e89b-12d3-a456-426614174005', '123e4567-e89b-12d3-a456-426614174003', 75.00, 'CREDIT', 'One-time service fee', '2025-02-01T10:00:00Z', 'PENDING', 62.50, 75.00, 20.00, 'SERVICE', NOW(), NOW()),

-- Movements for invoice 123e4567-e89b-12d3-a456-426614174004 (DRAFT)
('233e4567-e89b-12d3-a456-426614174006', '123e4567-e89b-12d3-a456-426614174004', 300.00, 'CREDIT', 'Professional consulting', '2025-03-10T15:00:00Z', 'PENDING', 250.00, 300.00, 20.00, 'CONSULTING', NOW(), NOW()),
('233e4567-e89b-12d3-a456-426614174007', '123e4567-e89b-12d3-a456-426614174004', 200.00, 'CREDIT', 'Implementation fee', '2025-03-10T15:10:00Z', 'PENDING', 166.67, 200.00, 20.00, 'IMPLEMENTATION', NOW(), NOW()),

-- Movements for invoice 123e4567-e89b-12d3-a456-426614174005 (SENT)
('233e4567-e89b-12d3-a456-426614174008', '123e4567-e89b-12d3-a456-426614174005', 120.25, 'CREDIT', 'Software license', '2025-03-15T17:00:00Z', 'PENDING', 100.21, 120.25, 20.00, 'LICENSE', NOW(), NOW());
