package models

import "time"

type Account struct {
	ID          int       `json:"id" bun:",pk,autoincrement"`
	Name        string    `json:"name"`
	AccountType string    `json:"accountType"`
	Currency    string    `json:"currency"`
	UserID      int       `json:"userID"`
	UpdatedAt   time.Time `json:"-" bun:"default:current_timestamp"`
	CreatedAt   time.Time `json:"-" bun:"default:current_timestamp"`
}
