-- +goose Up
-- SQL in section 'Up' is executed when migrating up

CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TABLE IF NOT EXISTS shift_allocations (
    id BIGSERIAL,
    op_item_id BIGINT NOT NULL,
    market_id BIGINT NOT NULL,
    zone_id BIGINT NOT NULL,
    type TEXT NOT NULL,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    inserted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (id, starts_at)
);

SELECT create_hypertable('shift_allocations', 'starts_at', if_not_exists => TRUE);

-- +goose Down
-- SQL in section 'Down' is executed when rolling back

DROP TABLE IF EXISTS shift_allocations;
DROP EXTENSION IF EXISTS timescaledb;