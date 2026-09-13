-- name: CreatePayoutAdjustment :one
INSERT INTO payout_adjustments (
    id,
    payout_id,
    order_id,
    type,
    gross_delta,
    commission_delta,
    net_delta,
    reason
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetPayoutAdjustmentByID :one
SELECT *
FROM payout_adjustments
WHERE id = $1;

-- name: ListPayoutAdjustmentsByPayoutID :many
SELECT *
FROM payout_adjustments
WHERE payout_id = $1
ORDER BY created_at ASC;

-- name: ListPayoutAdjustments :many
SELECT *
FROM payout_adjustments
WHERE (sqlc.narg('payout_id')::uuid IS NULL OR payout_id = sqlc.narg('payout_id')::uuid)
  AND (sqlc.narg('order_id')::uuid  IS NULL OR order_id  = sqlc.narg('order_id')::uuid)
  AND (sqlc.narg('type')::adjustment_type IS NULL OR type = sqlc.narg('type')::adjustment_type)
  AND (sqlc.narg('created_from')::timestamptz IS NULL OR created_at >= sqlc.narg('created_from')::timestamptz)
  AND (sqlc.narg('created_to')::timestamptz   IS NULL OR created_at <  sqlc.narg('created_to')::timestamptz)
ORDER BY
    CASE WHEN sqlc.arg('sort_by')::text = 'created_at' AND sqlc.arg('sort_dir')::text = 'asc' THEN created_at END ASC,
    created_at DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: CountPayoutAdjustments :one
SELECT COUNT(*)
FROM payout_adjustments
WHERE (sqlc.narg('payout_id')::uuid IS NULL OR payout_id = sqlc.narg('payout_id')::uuid)
  AND (sqlc.narg('order_id')::uuid  IS NULL OR order_id  = sqlc.narg('order_id')::uuid)
  AND (sqlc.narg('type')::adjustment_type IS NULL OR type = sqlc.narg('type')::adjustment_type)
  AND (sqlc.narg('created_from')::timestamptz IS NULL OR created_at >= sqlc.narg('created_from')::timestamptz)
  AND (sqlc.narg('created_to')::timestamptz   IS NULL OR created_at <  sqlc.narg('created_to')::timestamptz);
