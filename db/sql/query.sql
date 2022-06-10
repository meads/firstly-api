-- name: GetImage :one
SELECT * FROM images
WHERE id = $1 LIMIT 1;

-- name: ListImages :many
SELECT * FROM images
ORDER BY name;

-- name: CreateImage :one
INSERT INTO images (
  name, data, created
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: SoftDeleteImage :exec
UPDATE images
SET deleted = TRUE
WHERE id = $1;

-- name: DeleteImage :exec
DELETE FROM images
WHERE id = $1;

