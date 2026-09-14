-- name: TryAdvisoryLock :one
SELECT pg_try_advisory_lock(sqlc.arg('key')::bigint);

-- name: AdvisoryUnlock :one
SELECT pg_advisory_unlock(sqlc.arg('key')::bigint);
