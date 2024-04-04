-- name: CreateUrl :one
INSERT INTO urls (
  url,
  code,
  owner
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetUrlByCode :one
SELECT * FROM urls
WHERE code = $1
LIMIT 1;

-- name: ListUrlsByUser :many
SELECT * FROM urls
WHERE owner = $1;