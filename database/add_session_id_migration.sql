-- Migration: Add session_id to orders table
-- Date: 2025-10-13
-- Purpose: Fix order creation for guest users

-- Add session_id column to orders.orders table
ALTER TABLE orders.orders 
ADD COLUMN session_id VARCHAR(255);

-- Add index for better performance
CREATE INDEX idx_orders_session ON orders.orders(session_id);

-- Verify the change
SELECT column_name, data_type, is_nullable 
FROM information_schema.columns 
WHERE table_schema = 'orders' 
  AND table_name = 'orders' 
  AND column_name = 'session_id';