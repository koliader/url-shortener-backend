// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import ()

type Url struct {
	Url    string  `json:"url"`
	Code   string  `json:"code"`
	Owner  *string `json:"owner"`
	Clicks int32   `json:"clicks"`
}

type User struct {
	Username string  `json:"username"`
	Password *string `json:"password"`
	Color    string  `json:"color"`
}
