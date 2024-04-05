-- name: CreateUser :one
INSERT INTO users (
  username,
  password,
  color
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1
LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users;
