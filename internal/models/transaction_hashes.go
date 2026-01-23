package models

import "time"

type TransactionHash struct {
	TransactionID int       `json:"transactionID"`
	Hash          string    `json:"hash"`
	CreatedAt     time.Time `json:"-" bun:"default:current_timestamp"`
	UpdatedAt     time.Time `json:"-" bun:"default:current_timestamp"`
}
