-- name: CreatePayout :one
INSERT INTO payouts (
    id,
    courier_id,
    period,
    delivered_count,
    gross_amount,
    commission_rate,
    commission_amount,
    net_amount
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetPayoutByID :one
SELECT *
FROM payouts
WHERE id = $1;

-- name: GetPayoutByCourierPeriod :one
SELECT *
FROM payouts
WHERE courier_id = $1
  AND period = $2;

-- name: UpdatePayoutStatus :one
UPDATE payouts
SET
    status     = sqlc.arg('status'),
    paid_at    = CASE
                     WHEN sqlc.arg('status')::payout_status = 'paid' THEN NOW()
                     ELSE paid_at
                 END,
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: ListPayouts :many
SELECT *
FROM payouts
WHERE (sqlc.narg('courier_id')::uuid IS NULL OR courier_id = sqlc.narg('courier_id')::uuid)
  AND (sqlc.narg('status')::payout_status IS NULL OR status = sqlc.narg('status')::payout_status)
  AND (sqlc.narg('period_from')::date IS NULL OR period >= sqlc.narg('period_from')::date)
  AND (sqlc.narg('period_to')::date   IS NULL OR period <= sqlc.narg('period_to')::date)
ORDER BY
    CASE WHEN sqlc.arg('sort_by')::text = 'period'     AND sqlc.arg('sort_dir')::text = 'asc'  THEN period     END ASC,
    CASE WHEN sqlc.arg('sort_by')::text = 'period'     AND sqlc.arg('sort_dir')::text = 'desc' THEN period     END DESC,
    CASE WHEN sqlc.arg('sort_by')::text = 'net_amount' AND sqlc.arg('sort_dir')::text = 'asc'  THEN net_amount END ASC,
    CASE WHEN sqlc.arg('sort_by')::text = 'net_amount' AND sqlc.arg('sort_dir')::text = 'desc' THEN net_amount END DESC,
    CASE WHEN sqlc.arg('sort_by')::text = 'created_at' AND sqlc.arg('sort_dir')::text = 'asc'  THEN created_at END ASC,
    created_at DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: CountPayouts :one
SELECT COUNT(*)
FROM payouts
WHERE (sqlc.narg('courier_id')::uuid IS NULL OR courier_id = sqlc.narg('courier_id')::uuid)
  AND (sqlc.narg('status')::payout_status IS NULL OR status = sqlc.narg('status')::payout_status)
  AND (sqlc.narg('period_from')::date IS NULL OR period >= sqlc.narg('period_from')::date)
  AND (sqlc.narg('period_to')::date   IS NULL OR period <= sqlc.narg('period_to')::date);
