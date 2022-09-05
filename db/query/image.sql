
-- name: GetImage :one
SELECT * FROM image
WHERE id = $1 LIMIT 1;

-- name: ListImages :many
SELECT * FROM image;

-- name: CreateImage :one
INSERT INTO image (
  data, created
) VALUES (
  $1, NOW()
)
RETURNING *;

-- name: SoftDeleteImage :exec
UPDATE image
SET deleted = 1
WHERE id = $1;

-- name: DeleteImage :exec
DELETE FROM image
WHERE id = $1;

