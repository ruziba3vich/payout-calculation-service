DROP INDEX IF EXISTS "idx_orders_created_at";
DROP INDEX IF EXISTS "idx_orders_courier_delivered";
DROP INDEX IF EXISTS "idx_orders_status";
DROP INDEX IF EXISTS "idx_orders_courier_id";
DROP TABLE IF EXISTS "orders";
DROP TYPE IF EXISTS "order_status";
