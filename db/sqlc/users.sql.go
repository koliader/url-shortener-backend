// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: users.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  username,
  email,
  password
) VALUES (
  $1, $2, $3
) RETURNING username, email, password
`

type CreateUserParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Username, arg.Email, arg.Password)
	var i User
	err := row.Scan(&i.Username, &i.Email, &i.Password)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT username, email, password FROM users
WHERE email = $1
LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(&i.Username, &i.Email, &i.Password)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT username, email, password FROM users
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(&i.Username, &i.Email, &i.Password); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
