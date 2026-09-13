-- payouts

CREATE TYPE "payout_status" AS ENUM (
    'calculated',
    'paid',
    'cancelled'
);

CREATE TABLE IF NOT EXISTS payouts (
    "id" UUID PRIMARY KEY,
    "courier_id" UUID NOT NULL REFERENCES couriers (id),
    "period" DATE NOT NULL CHECK ("period" = DATE_TRUNC('month', "period")::date),
    "delivered_count" INTEGER NOT NULL DEFAULT 0 CHECK ("delivered_count" >= 0),
    "gross_amount" NUMERIC(14, 2) NOT NULL DEFAULT 0 CHECK ("gross_amount" >= 0),
    "commission_rate" NUMERIC(5, 4) NOT NULL CHECK ("commission_rate" >= 0 AND "commission_rate" <= 1),
    "commission_amount" NUMERIC(14, 2) NOT NULL DEFAULT 0 CHECK ("commission_amount" >= 0),
    "net_amount" NUMERIC(14, 2) NOT NULL DEFAULT 0,
    "status" payout_status NOT NULL DEFAULT 'calculated',
    "calculated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "paid_at" TIMESTAMPTZ NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT "uq_payouts_courier_period" UNIQUE ("courier_id", "period")
);

CREATE INDEX IF NOT EXISTS "idx_payouts_courier_id"
    ON "payouts" ("courier_id");

CREATE INDEX IF NOT EXISTS "idx_payouts_period"
    ON "payouts" ("period" DESC);

CREATE INDEX IF NOT EXISTS "idx_payouts_status"
    ON "payouts" ("status");
