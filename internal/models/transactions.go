package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Transaction struct {
	ID                  int       `json:"id" bun:",pk,autoincrement"`
	AccountID           int       `json:"accountID"`
	CategoryID          int       `json:"categoryID"`
	TransactionTypeCode string    `json:"transactionTypeCode"`
	Amount              int       `json:"amount"`
	Description         string    `json:"description"`
	TransactionDate     time.Time `json:"transactionDate"`
	UpdatedAt           time.Time `json:"-" bun:"default:current_timestamp"`
	CreatedAt           time.Time `json:"-" bun:"default:current_timestamp"`
}

type TransactionAPI struct {
	bun.BaseModel `bun:"transactions,alias:t"`

	ID                  int       `json:"id"`
	TransactionDate     time.Time `json:"transactionDate"`
	Category            string    `json:"category"`
	SubCategory         string    `json:"subCategory"`
	Amount              int       `json:"amount"`
	Account             string    `json:"account"`
	Description         string    `json:"description"`
	TransactionTypeCode string    `json:"transactionTypeCode"`
}

type TransactionAPIRequestParams struct {
	GteDate                      string `json:"gteDate"`
	IsGteDateSpecify             bool   `json:"isGteDateSpecify,omitempty"`
	LteDate                      string `json:"lteDate"`
	IsLteDateSpecify             bool   `json:"isLteDateSpecify,omitempty"`
	TransactionTypeCode          string `json:"transactionTypeCode"`
	IsTransactionTypeCodeSpecify bool   `json:"isTransactionTypeCodeSpecify,omitempty"`
}

type SumTransactionByCategory struct {
	bun.BaseModel `bun:"transactions,alias:t"`

	Category string `json:"category"`
	Amount   int    `json:"amount"`
}
