package database

import (
	"context"
	"database/sql"
	"fmt"
	"go-fwallet/internal/helpers"
	"go-fwallet/internal/middleware"
	"go-fwallet/internal/models"
	"net/http"
	"time"
)

func (d *Database) GetAccountsStatements(closingDate string, isClosingDateSpecify bool, c context.Context) ([]models.AccountStatementAPIResponce, error) {
	var ass []models.AccountStatementAPIResponce

	q := d.Client.NewSelect().Model(&ass).
		ColumnExpr("accss.closing_balance, accss.total_credit, accss.total_debit").
		ColumnExpr("a.name as account").
		Join("LEFT JOIN accounts as a ON accss.account_id = a.id")

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
	var a []int

	err := d.Client.NewSelect().ColumnExpr("account_id").TableExpr("transactions").Where("transaction_date >= ? and transaction_date <= ?", rp.GteDate, rp.LteDate).GroupExpr("account_id").Scan(c, &a)
	if err != nil {
		return nil, err
	}

	for _, v := range a {
		as, err := d.CreateAccountStatement(fmt.Sprint(v), rp, c)
		if err != nil {
			return nil, err
		}
		ass = append(ass, *as)
	}

	return ass, nil
}

func (d *Database) GetAccountStatement(idStr string, closingDate string, isClosingDateSpecify bool, c context.Context) (*models.AccountStatementAPIResponce, error) {
	id, err := helpers.StringToInt(idStr)
	if err != nil {
		return nil, middleware.NewHttpError("cannot conver account id to int", err.Error(), http.StatusBadRequest)
	}

	var as models.AccountStatementAPIResponce
	q := d.Client.NewSelect().Model(&as).
		ColumnExpr("accss.closing_balance, accss.total_credit, accss.total_debit").
		ColumnExpr("a.name as account").
		Join("LEFT JOIN accounts as a ON accss.account_id = a.id").
		Where("account_id = ?", id)

	if isClosingDateSpecify {
		q.Where("closing_date = ?", closingDate)
	}

	err = q.Scan(c)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, middleware.NewHttpError("account statment not found", fmt.Sprintf("accountID:%d", id), http.StatusNotFound)
		}
		return nil, err
	}

	return &as, nil
}

func (d *Database) CreateAccountStatement(idStr string, rp *models.AccountStatementRequestPayload, c context.Context) (*models.AccountStatement, error) {
	id, err := helpers.StringToInt(idStr)
	if err != nil {
		return nil, middleware.NewHttpError("cannot conver account id to int", err.Error(), http.StatusBadRequest)
	}

	var as models.AccountStatement
	err = d.Client.NewSelect().
		ColumnExpr("account_id").
		ColumnExpr("sum(amount) FILTER (where transaction_type_code = 'W') as total_credit").
		ColumnExpr("sum(amount) FILTER (where transaction_type_code = 'D') as total_debit").
		TableExpr("transactions").
		Where("transaction_date >= ? and transaction_date <= ? and account_id = ?", rp.GteDate, rp.LteDate, id).
		GroupExpr("account_id").
		Scan(c, &as)

	if err != nil {
		return nil, err
	}

	closingDate, err := time.Parse("2006-01-02", rp.LteDate)
	if err != nil {
		return nil, err
	}

	as.ClosingBalance = as.TotalDebit - as.TotalCredit
	as.ClosingDate = closingDate.AddDate(0, 0, 1)

	err = d.Client.NewRaw(`WITH upsert AS (
								UPDATE account_statements 
								SET updated_at = ?5
								WHERE (account_id = ?0 and closing_date = ?1)
								RETURNING *
							) 
							INSERT INTO account_statements (account_id, closing_date, closing_balance, total_credit, total_debit) 
							SELECT ?0, ?1, ?2, ?3, ?4 
							WHERE NOT EXISTS (SELECT 1 FROM upsert)`, id, as.ClosingDate.Format("2006-01-02"), as.ClosingBalance, as.TotalCredit, as.TotalDebit, time.Now()).Scan(c, &as)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &as, nil
}
