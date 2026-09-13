-- payout_adjustments

CREATE TYPE "adjustment_type" AS ENUM (
    'order_returned',
    'order_cancelled',
    'order_added',
    'manual'
);

CREATE TABLE IF NOT EXISTS payout_adjustments (
    "id" UUID PRIMARY KEY,
    "payout_id" UUID NOT NULL REFERENCES payouts (id),
    "order_id" UUID NULL REFERENCES orders (id),
    "type" adjustment_type NOT NULL,
    "gross_delta" NUMERIC(14, 2) NOT NULL DEFAULT 0,
    "commission_delta" NUMERIC(14, 2) NOT NULL DEFAULT 0,
    "net_delta" NUMERIC(14, 2) NOT NULL DEFAULT 0,
    "reason" TEXT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT "uq_payout_adjustments_payout_order_type"
        UNIQUE ("payout_id", "order_id", "type")
);

CREATE INDEX IF NOT EXISTS "idx_payout_adjustments_payout_id"
    ON "payout_adjustments" ("payout_id");

CREATE INDEX IF NOT EXISTS "idx_payout_adjustments_order_id"
    ON "payout_adjustments" ("order_id");

CREATE INDEX IF NOT EXISTS "idx_payout_adjustments_created_at"
    ON "payout_adjustments" ("created_at" DESC);
