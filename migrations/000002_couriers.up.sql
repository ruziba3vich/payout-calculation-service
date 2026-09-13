-- couriers

CREATE TABLE IF NOT EXISTS couriers (
    id UUID PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(32) NOT NULL,
    password VARCHAR(255) NOT NULL,
    hired_at DATE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_couriers_phone_active
    ON couriers (phone)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_couriers_is_active
    ON couriers (is_active)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_couriers_created_at
    ON couriers (created_at DESC);
