package models

import "time"

type User struct {
	ID        int       `json:"id" bun:",pk,autoincrement"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"-" bun:"default:current_timestamp"`
	UpdatedAt time.Time `json:"-" bun:"default:current_timestamp"`
}
