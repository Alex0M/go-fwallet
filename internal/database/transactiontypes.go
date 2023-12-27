package database

import (
	"context"
	"database/sql"
	"fmt"
	"go-fwallet/internal/middleware"
	"go-fwallet/internal/models"
	"net/http"
)

func (d *Database) GetTransactionTypes(c context.Context) ([]models.TransactionType, error) {
	var tt []models.TransactionType
	err := d.Client.NewSelect().Model(&tt).Scan(c)
	if err != nil {
		return nil, err
	}

	return tt, nil
}

func (d *Database) GetTransactionType(code string, c context.Context) (*models.TransactionType, error) {
	var tt models.TransactionType
	err := d.Client.NewSelect().Model(&tt).Where("transaction_type_code = ?", code).Scan(c)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, middleware.NewHttpError("transaction type not found", fmt.Sprintf("transactionTypeCode:%s", code), http.StatusNotFound)
		}
		return nil, err
	}

	return &tt, nil
}

func (d *Database) AddTransactionType(tt *models.TransactionType, c context.Context) error {
	_, err := d.Client.NewInsert().Model(tt).Exec(c)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) EditTransactionType(code string, newTT *models.TransactionType, c context.Context) (*models.TransactionType, error) {
	oldTT, err := d.GetTransactionType(code, c)
	if err != nil {
		return nil, err
	}

	_, err = d.Client.NewUpdate().Model(newTT).OmitZero().Where("transaction_type_code = ?", oldTT.TransactionTypeCode).Returning("*").Exec(c)
	if err != nil {
		return nil, err
	}

	return newTT, nil
}

func (d *Database) DeleteTransactionType(code string, c context.Context) error {
	tt, err := d.GetTransactionType(code, c)
	if err != nil {
		return err
	}

	_, err = d.Client.NewDelete().Model(tt).WherePK().Exec(c)
	if err != nil {
		return err
	}

	return nil
}
