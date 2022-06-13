-- name: GetImage :one
SELECT * FROM images
WHERE id = $1 LIMIT 1;

-- name: ListImages :many
SELECT * FROM images;

-- name: CreateImage :one
INSERT INTO images (
  data, created
) VALUES (
  $1, NOW()
)
RETURNING *;

-- name: SoftDeleteImage :exec
UPDATE images
SET deleted = 1
WHERE id = $1;

-- name: DeleteImage :exec
DELETE FROM images
WHERE id = $1;

