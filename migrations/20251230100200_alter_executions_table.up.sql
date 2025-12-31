ALTER TABLE executions ADD COLUMN symbol VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE executions RENAME COLUMN execution_time TO executed_at;
ALTER TABLE executions RENAME COLUMN execution_price TO price;
ALTER TABLE executions RENAME COLUMN execution_quantity TO quantity;