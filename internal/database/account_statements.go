package database

import (
	"context"
	"go-fwallet/internal/models"
	"time"
)

func (d *Database) GetAccountsStatements(closingDate string, isClosingDateSpecify bool, c context.Context) ([]models.AccountStatement, error) {
	var ass []models.AccountStatement

	q := d.Client.NewSelect().Model(&ass)
	if isClosingDateSpecify {
		q.Where("closing_date = ?", closingDate)
	}

	err := q.Order("closing_date").Scan(c)
	if err != nil {
		return nil, err
	}

	return ass, nil
}

func (d *Database) CreateAccountsStatements(rp *models.AccountStatementRequestPayload, c context.Context) ([]models.AccountStatement, error) {
	var ass []models.AccountStatement

	err := d.Client.NewRaw(
		"SELECT account_id, sum(amount) FILTER (where transaction_type_code = 'W') as total_credit, sum(amount) FILTER (where transaction_type_code = 'D') as total_debit FROM public.transactions WHERE transaction_date >= ? and transaction_date <= ? group by account_id",
		rp.GteDate, rp.LteDate,
	).Scan(c, &ass)
	if err != nil {
		return nil, err
	}

	closingDate, err := time.Parse("2006-01-02", rp.LteDate)
	if err != nil {
		return nil, err
	}

	for i := range ass {
		ass[i].ClosingBalance = ass[i].TotalDebit - ass[i].TotalCredit
		ass[i].ClosingDate = closingDate.AddDate(0, 0, 1)
	}

	_, err = d.Client.NewInsert().Model(&ass).Exec(c)
	if err != nil {
		return nil, err
	}

	return ass, nil
}
