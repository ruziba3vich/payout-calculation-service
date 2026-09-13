-- administration

CREATE TABLE IF NOT EXISTS administration (
    id          UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name   VARCHAR(255)    NOT NULL,
    username    VARCHAR(100)    NOT NULL,
    password    VARCHAR(255)    NOT NULL,          -- store bcrypt/argon2 hash
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ     NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_administration_username_active
    ON administration (username)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_administration_deleted_at
    ON administration (deleted_at)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_administration_created_at
    ON administration (created_at DESC);
