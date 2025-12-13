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

CREATE TABLE IF NOT EXISTS stock_masters (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    issue_code VARCHAR(255) UNIQUE,
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

CREATE TABLE IF NOT EXISTS tick_rules (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    issue_code VARCHAR(255),
    tick_unit_number VARCHAR(255),
    applicable_date VARCHAR(255)
);

CREATE INDEX IF NOT EXISTS idx_tick_rules_deleted_at ON tick_rules(deleted_at);
CREATE INDEX IF NOT EXISTS idx_tick_rules_issue_code ON tick_rules(issue_code);

CREATE TABLE IF NOT EXISTS tick_levels (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    tick_rule_id BIGINT,
    lower_price DOUBLE PRECISION,
    upper_price DOUBLE PRECISION,
    tick_value DOUBLE PRECISION,
    CONSTRAINT fk_tick_levels_tick_rule
        FOREIGN KEY(tick_rule_id)
        REFERENCES tick_rules(id)
);

CREATE INDEX IF NOT EXISTS idx_tick_levels_deleted_at ON tick_levels(deleted_at);
