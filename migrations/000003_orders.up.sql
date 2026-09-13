-- orders

CREATE TYPE "order_status" AS ENUM (
    'pending',
    'delivered',
    'cancelled',
    'returned'
);

CREATE TABLE IF NOT EXISTS orders (
    "id" UUID PRIMARY KEY,
    "courier_id" UUID NOT NULL REFERENCES couriers (id),
    "amount" NUMERIC(14, 2) NOT NULL CHECK (amount >= 0),
    "status" order_status NOT NULL DEFAULT 'pending',
    "delivered_at" TIMESTAMPTZ NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS "idx_orders_courier_id"
    ON "orders" ("courier_id");

CREATE INDEX IF NOT EXISTS "idx_orders_status"
    ON "orders" ("status");

CREATE INDEX IF NOT EXISTS "idx_orders_courier_delivered"
    ON "orders" ("courier_id", "delivered_at")
    WHERE "status" = 'delivered';

CREATE INDEX IF NOT EXISTS "idx_orders_created_at"
    ON "orders" ("created_at" DESC);
