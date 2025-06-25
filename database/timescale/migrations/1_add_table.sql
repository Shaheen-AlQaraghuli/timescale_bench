-- +goose Up

CREATE TABLE IF NOT EXISTS shift_allocations (
    id bigserial primary key,
    op_item_id BIGINT NOT NULL,
    market_id BIGINT NOT NULL,
    zone_id BIGINT NOT NULL,
    type TEXT NOT NULL,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    inserted_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Separate indexes instead of composite PK
CREATE INDEX IF NOT EXISTS idx_starts_at ON shift_allocations (starts_at);

-- +goose Down
DROP TABLE IF EXISTS shift_allocations;
