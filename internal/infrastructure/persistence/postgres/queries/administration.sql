-- name: CreateAdmin :one
INSERT INTO administration (
    id,
    full_name,
    username,
    password
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetAdminByID :one
SELECT *
FROM administration
WHERE id = $1
  AND deleted_at IS NULL;

-- name: GetAdminByUsername :one
SELECT *
FROM administration
WHERE username = $1
  AND deleted_at IS NULL;

-- name: UpdateAdmin :one
UPDATE administration
SET
    full_name  = COALESCE($2, full_name),
    username   = COALESCE($3, username),
    password   = COALESCE($4, password),
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteAdmin :exec
UPDATE administration
SET
    deleted_at = NOW(),
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL;

-- name: ListAdmins :many
SELECT *
FROM administration
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountAdmins :one
SELECT COUNT(*)
FROM administration
WHERE deleted_at IS NULL;
