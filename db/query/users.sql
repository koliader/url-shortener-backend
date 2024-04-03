-- name: CreateUser :one
INSERT INTO users (
  username,
  email
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users;
