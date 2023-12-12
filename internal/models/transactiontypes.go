package models

import "time"

type TransactionType struct {
	TransactionTypeCode string    `json:"transactionTypeCode" bun:",pk"`
	Description         string    `json:"description"`
	CreatedAt           time.Time `json:"-" bun:"default:current_timestamp"`
	UpdatedAt           time.Time `json:"-" bun:"default:current_timestamp"`
}
