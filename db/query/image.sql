
-- name: Get :one
SELECT * FROM image
WHERE id = $1 LIMIT 1;

-- name: List :many
SELECT * FROM image LIMIT $1 OFFSET $2;

-- name: Create :one
INSERT INTO image (
  data, created
) VALUES (
  $1, NOW()
)
RETURNING *;

-- name: SoftDelete :exec
UPDATE image
SET deleted = 1
WHERE id = $1;

-- name: Delete :exec
DELETE FROM image
WHERE id = $1;

