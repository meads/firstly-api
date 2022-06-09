-- name: GetImage :one
SELECT * FROM image
WHERE id = $1 LIMIT 1;

-- name: ListImages :many
SELECT * FROM image
ORDER BY name;

-- name: CreateImage :one
INSERT INTO image (
  name, data
) VALUES (
  $1, $2
)
RETURNING *;

-- name: DeleteImage :exec
DELETE FROM image
WHERE id = $1;
