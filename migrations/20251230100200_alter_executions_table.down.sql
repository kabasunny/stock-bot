ALTER TABLE executions DROP COLUMN symbol;
ALTER TABLE executions RENAME COLUMN executed_at TO execution_time;
ALTER TABLE executions RENAME COLUMN price TO execution_price;
ALTER TABLE executions RENAME COLUMN quantity TO execution_quantity;