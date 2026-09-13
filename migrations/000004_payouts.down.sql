DROP INDEX IF EXISTS "idx_payouts_status";
DROP INDEX IF EXISTS "idx_payouts_period";
DROP INDEX IF EXISTS "idx_payouts_courier_id";
DROP TABLE IF EXISTS payouts;
DROP TYPE IF EXISTS "payout_status";
