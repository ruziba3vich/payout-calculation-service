-- name: CreateOrder :one
INSERT INTO orders (
    id,
    courier_id,
    amount
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetOrderByID :one
SELECT *
FROM orders
WHERE id = $1;

-- name: UpdateOrderStatus :one
UPDATE orders
SET
    status       = sqlc.arg('status'),
    delivered_at = CASE
                       WHEN sqlc.arg('status')::order_status = 'delivered' THEN NOW()
                       ELSE delivered_at
                   END,
    updated_at   = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;

-- name: ListOrders :many
SELECT *
FROM orders
WHERE (sqlc.narg('courier_id')::uuid IS NULL OR courier_id = sqlc.narg('courier_id')::uuid)
    AND (sqlc.narg('status')::order_status IS NULL OR status = sqlc.narg('status')::order_status)
    AND (sqlc.narg('delivered_from')::timestamptz IS NULL OR delivered_at >= sqlc.narg('delivered_from')::timestamptz)
    AND (sqlc.narg('delivered_to')::timestamptz   IS NULL OR delivered_at <  sqlc.narg('delivered_to')::timestamptz)
    AND (sqlc.narg('created_from')::timestamptz   IS NULL OR created_at   >= sqlc.narg('created_from')::timestamptz)
    AND (sqlc.narg('created_to')::timestamptz     IS NULL OR created_at   <  sqlc.narg('created_to')::timestamptz)
    AND (sqlc.narg('amount_min')::numeric IS NULL OR amount >= sqlc.narg('amount_min')::numeric)
    AND (sqlc.narg('amount_max')::numeric IS NULL OR amount <= sqlc.narg('amount_max')::numeric)
ORDER BY
    CASE WHEN sqlc.arg('sort_by')::text = 'amount'       AND sqlc.arg('sort_dir')::text = 'asc'  THEN amount       END ASC,
    CASE WHEN sqlc.arg('sort_by')::text = 'amount'       AND sqlc.arg('sort_dir')::text = 'desc' THEN amount       END DESC,
    CASE WHEN sqlc.arg('sort_by')::text = 'delivered_at' AND sqlc.arg('sort_dir')::text = 'asc'  THEN delivered_at END ASC,
    CASE WHEN sqlc.arg('sort_by')::text = 'delivered_at' AND sqlc.arg('sort_dir')::text = 'desc' THEN delivered_at END DESC,
    CASE WHEN sqlc.arg('sort_by')::text = 'created_at'   AND sqlc.arg('sort_dir')::text = 'asc'  THEN created_at   END ASC,
    created_at DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: CountOrders :one
SELECT COUNT(*)
FROM orders
WHERE (sqlc.narg('courier_id')::uuid IS NULL OR courier_id = sqlc.narg('courier_id')::uuid)
    AND (sqlc.narg('status')::order_status IS NULL OR status = sqlc.narg('status')::order_status)
    AND (sqlc.narg('delivered_from')::timestamptz IS NULL OR delivered_at >= sqlc.narg('delivered_from')::timestamptz)
    AND (sqlc.narg('delivered_to')::timestamptz   IS NULL OR delivered_at <  sqlc.narg('delivered_to')::timestamptz)
    AND (sqlc.narg('created_from')::timestamptz   IS NULL OR created_at   >= sqlc.narg('created_from')::timestamptz)
    AND (sqlc.narg('created_to')::timestamptz     IS NULL OR created_at   <  sqlc.narg('created_to')::timestamptz)
    AND (sqlc.narg('amount_min')::numeric IS NULL OR amount >= sqlc.narg('amount_min')::numeric)
    AND (sqlc.narg('amount_max')::numeric IS NULL OR amount <= sqlc.narg('amount_max')::numeric);
