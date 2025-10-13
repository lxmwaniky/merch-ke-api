-- Migration: Add missing notes column  
-- Date: 2025-10-13
-- Purpose: Fix order creation - add notes column expected by backend

-- Add notes column to orders.orders table
ALTER TABLE orders.orders ADD COLUMN notes TEXT;

-- Verify the change
SELECT column_name, data_type, is_nullable 
FROM information_schema.columns 
WHERE table_schema = 'orders' 
  AND table_name = 'orders' 
  AND column_name = 'notes';