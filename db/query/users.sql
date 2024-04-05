-- name: CreateUser :one
INSERT INTO users (
  username,
  email,
  password,
  color
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users;
