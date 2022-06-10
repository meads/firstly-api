// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: query.sql

package api

import (
	"context"
)

const createImage = `-- name: CreateImage :one
INSERT INTO images (
  data, created
) VALUES (
  $1, NOW()
)
RETURNING id, name, data, created, deleted
`

func (q *Queries) CreateImage(ctx context.Context, data string) (Image, error) {
	row := q.db.QueryRowContext(ctx, createImage, data)
	var i Image
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Data,
		&i.Created,
		&i.Deleted,
	)
	return i, err
}

const deleteImage = `-- name: DeleteImage :exec
DELETE FROM images
WHERE id = $1
`

func (q *Queries) DeleteImage(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteImage, id)
	return err
}

const getImage = `-- name: GetImage :one
SELECT id, name, data, created, deleted FROM images
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetImage(ctx context.Context, id int64) (Image, error) {
	row := q.db.QueryRowContext(ctx, getImage, id)
	var i Image
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Data,
		&i.Created,
		&i.Deleted,
	)
	return i, err
}

const listImages = `-- name: ListImages :many
SELECT id, name, data, created, deleted FROM images
ORDER BY name
`

func (q *Queries) ListImages(ctx context.Context) ([]Image, error) {
	rows, err := q.db.QueryContext(ctx, listImages)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Image
	for rows.Next() {
		var i Image
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Data,
			&i.Created,
			&i.Deleted,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const softDeleteImage = `-- name: SoftDeleteImage :exec
UPDATE images
SET deleted = TRUE
WHERE id = $1
`

func (q *Queries) SoftDeleteImage(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, softDeleteImage, id)
	return err
}
