-- +goose Up
CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TABLE IF NOT EXISTS shift_allocations (
    allocation_id BIGINT NOT NULL,
    type TEXT NOT NULL,
    duration_minutes INTEGER NOT NULL,
    date TIMESTAMPTZ NOT NULL,
    op_item_id BIGINT NOT NULL,
    zone_id BIGINT NOT NULL,
    market_id BIGINT NOT NULL
);

create index indx_allocation_id on shift_allocations (allocation_id);
SELECT create_hypertable('shift_allocations', 'date', if_not_exists => TRUE);

-- +goose Down
DROP TABLE IF EXISTS shift_allocations;
DROP EXTENSION IF EXISTS timescaledb;