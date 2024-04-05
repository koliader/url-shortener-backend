// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: urls.sql

package db

import (
	"context"
)

const createUrl = `-- name: CreateUrl :one
INSERT INTO urls (
  url,
  code,
  owner
) VALUES (
  $1, $2, $3
) RETURNING url, code, owner, clicks
`

type CreateUrlParams struct {
	Url   string  `json:"url"`
	Code  string  `json:"code"`
	Owner *string `json:"owner"`
}

func (q *Queries) CreateUrl(ctx context.Context, arg CreateUrlParams) (Url, error) {
	row := q.db.QueryRow(ctx, createUrl, arg.Url, arg.Code, arg.Owner)
	var i Url
	err := row.Scan(
		&i.Url,
		&i.Code,
		&i.Owner,
		&i.Clicks,
	)
	return i, err
}

const getUrlByCode = `-- name: GetUrlByCode :one
SELECT url, code, owner, clicks FROM urls
WHERE code = $1
LIMIT 1
`

func (q *Queries) GetUrlByCode(ctx context.Context, code string) (Url, error) {
	row := q.db.QueryRow(ctx, getUrlByCode, code)
	var i Url
	err := row.Scan(
		&i.Url,
		&i.Code,
		&i.Owner,
		&i.Clicks,
	)
	return i, err
}

const listUrlsByUser = `-- name: ListUrlsByUser :many
SELECT url, code, owner, clicks FROM urls
WHERE owner = $1
`

func (q *Queries) ListUrlsByUser(ctx context.Context, owner *string) ([]Url, error) {
	rows, err := q.db.Query(ctx, listUrlsByUser, owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Url{}
	for rows.Next() {
		var i Url
		if err := rows.Scan(
			&i.Url,
			&i.Code,
			&i.Owner,
			&i.Clicks,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateClicks = `-- name: UpdateClicks :one
UPDATE urls
SET clicks = clicks + 1
WHERE code = $1
RETURNING url, code, owner, clicks
`

func (q *Queries) UpdateClicks(ctx context.Context, code string) (Url, error) {
	row := q.db.QueryRow(ctx, updateClicks, code)
	var i Url
	err := row.Scan(
		&i.Url,
		&i.Code,
		&i.Owner,
		&i.Clicks,
	)
	return i, err
}
