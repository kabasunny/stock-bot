-- create_initial_tables.up.sql

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    order_id VARCHAR(255) UNIQUE,
    symbol VARCHAR(255),
    trade_type VARCHAR(255),
    order_type VARCHAR(255),
    quantity BIGINT NOT NULL,
    price DOUBLE PRECISION,
    trigger_price DOUBLE PRECISION,
    time_in_force VARCHAR(255) DEFAULT 'DAY',
    order_status VARCHAR(255),
    is_margin BOOLEAN NOT NULL DEFAULT false
);

CREATE INDEX IF NOT EXISTS idx_orders_deleted_at ON orders(deleted_at);
CREATE INDEX IF NOT EXISTS idx_orders_symbol ON orders(symbol);
CREATE INDEX IF NOT EXISTS idx_orders_trade_type ON orders(trade_type);
CREATE INDEX IF NOT EXISTS idx_orders_order_type ON orders(order_type);
CREATE INDEX IF NOT EXISTS idx_orders_time_in_force ON orders(time_in_force);
CREATE INDEX IF NOT EXISTS idx_orders_order_status ON orders(order_status);

CREATE TABLE IF NOT EXISTS executions (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    order_id VARCHAR(255),
    execution_id VARCHAR(255) UNIQUE,
    execution_time TIMESTAMPTZ,
    execution_price DOUBLE PRECISION,
    execution_quantity BIGINT,
    commission DOUBLE PRECISION,
    CONSTRAINT fk_executions_order
        FOREIGN KEY(order_id)
        REFERENCES orders(order_id)
);

CREATE INDEX IF NOT EXISTS idx_executions_deleted_at ON executions(deleted_at);
CREATE INDEX IF NOT EXISTS idx_executions_order_id ON executions(order_id);

CREATE TABLE IF NOT EXISTS positions (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    symbol VARCHAR(255),
    position_type VARCHAR(255),
    average_price DOUBLE PRECISION,
    quantity BIGINT
);

CREATE INDEX IF NOT EXISTS idx_positions_deleted_at ON positions(deleted_at);
CREATE INDEX IF NOT EXISTS idx_positions_symbol ON positions(symbol);
CREATE INDEX IF NOT EXISTS idx_positions_position_type ON positions(position_type);

CREATE TABLE IF NOT EXISTS signals (
    id BIGSERIAL PRIMARY KEY,
    symbol VARCHAR(255),
    signal_type VARCHAR(255),
    generated_at TIMESTAMPTZ,
    rationale TEXT,
    price DOUBLE PRECISION
);

CREATE INDEX IF NOT EXISTS idx_signals_symbol ON signals(symbol);

DROP TABLE IF EXISTS stock_masters;
CREATE TABLE IF NOT EXISTS stock_masters (
    issue_code VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    issue_name VARCHAR(255),
    trading_unit BIGINT,
    market_code VARCHAR(255),
    upper_limit DOUBLE PRECISION,
    lower_limit DOUBLE PRECISION
);
CREATE INDEX IF NOT EXISTS idx_stock_masters_deleted_at ON stock_masters(deleted_at);


CREATE TABLE IF NOT EXISTS stock_market_masters (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    issue_code VARCHAR(255),
    listing_market VARCHAR(255),
    previous_close DOUBLE PRECISION
);

CREATE INDEX IF NOT EXISTS idx_stock_market_masters_deleted_at ON stock_market_masters(deleted_at);
CREATE INDEX IF NOT EXISTS idx_stock_market_masters_issue_code ON stock_market_masters(issue_code);
CREATE INDEX IF NOT EXISTS idx_stock_market_masters_listing_market ON stock_market_masters(listing_market);

DROP TABLE IF EXISTS tick_levels;
DROP TABLE IF EXISTS tick_rules;

CREATE TABLE IF NOT EXISTS tick_rules (
    tick_unit_number VARCHAR(255) PRIMARY KEY,
    applicable_date VARCHAR(255),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS tick_levels (
    id BIGSERIAL PRIMARY KEY,
    tick_rule_unit_number VARCHAR(255) NOT NULL,
    lower_price DOUBLE PRECISION,
    upper_price DOUBLE PRECISION,
    tick_value DOUBLE PRECISION,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    CONSTRAINT fk_tick_levels_tick_rule
        FOREIGN KEY(tick_rule_unit_number)
        REFERENCES tick_rules(tick_unit_number)
);
CREATE INDEX IF NOT EXISTS idx_tick_levels_tick_rule_unit_number ON tick_levels(tick_rule_unit_number);
