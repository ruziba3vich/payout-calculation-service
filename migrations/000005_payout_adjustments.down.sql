DROP INDEX IF EXISTS "idx_payout_adjustments_created_at";
DROP INDEX IF EXISTS "idx_payout_adjustments_order_id";
DROP INDEX IF EXISTS "idx_payout_adjustments_payout_id";
DROP TABLE IF EXISTS payout_adjustments;
DROP TYPE IF EXISTS "adjustment_type";
