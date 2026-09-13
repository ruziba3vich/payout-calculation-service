-- name: CreateCourier :one
INSERT INTO couriers (
    id,
    full_name,
    phone,
    password,
    hired_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetCourierByID :one
SELECT *
FROM couriers
WHERE id = $1
  AND deleted_at IS NULL;

-- name: UpdateCourier :one
UPDATE couriers
SET
    full_name  = COALESCE(sqlc.narg('full_name'), full_name),
    phone      = COALESCE(sqlc.narg('phone'), phone),
    password   = COALESCE(sqlc.narg('password'), password),
    hired_at   = COALESCE(sqlc.narg('hired_at'), hired_at),
    is_active  = COALESCE(sqlc.narg('is_active'), is_active),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
  AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteCourier :exec
UPDATE couriers
SET
    deleted_at = NOW(),
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL;

-- name: ListCouriers :many
SELECT *
FROM couriers
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountCouriers :one
SELECT COUNT(*)
FROM couriers
WHERE deleted_at IS NULL;

-- name: CourierExists :one
SELECT EXISTS (
    SELECT 1
    FROM couriers
    WHERE id = $1
      AND deleted_at IS NULL
);
