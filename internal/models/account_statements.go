package models

import "time"

type AccountStatementRequestPayload struct {
	GteDate string `json:"gteDate"`
	LteDate string `json:"lteDate"`
}

type AccountStatement struct {
	AccountID      int       `json:"accountID"`
	ClosingDate    time.Time `json:"closingDate"`
	ClosingBalance int       `json:"closingBalance"`
	TotalCredit    int       `json:"totalCredit"`
	TotalDebit     int       `json:"totalDebit"`
	CreatedAt      time.Time `json:"-" bun:"default:current_timestamp"`
	UpdatedAt      time.Time `json:"-" bun:"default:current_timestamp"`
}
