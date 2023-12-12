package database

import (
	"context"
	"database/sql"
	"go-fwallet/internal/models"

	"go.uber.org/zap"
)

func (d *Database) GetTransactionTypes(c context.Context) ([]models.TransactionType, error) {
	var tt []models.TransactionType
	err := d.Client.NewSelect().Model(&tt).Scan(c)
	if err != nil {
		d.Logger.Error("error while getting all transaction types", zap.Error(err))
		return nil, err
	}

	return tt, nil
}

func (d *Database) GetTransactionType(code string, c context.Context) (*models.TransactionType, error) {
	d.Logger.Info("starting getting a transaction type by id", zap.Any("transactionTypeCode", code))
	var tt models.TransactionType
	err := d.Client.NewSelect().Model(&tt).Where("transaction_type_code = ?", code).Scan(c)
	if err != nil {
		if err == sql.ErrNoRows {
			d.Logger.Info("transaction type not found", zap.Any("transactionTypeCode", code))
			return nil, nil
		}
		d.Logger.Error("error while getting transactionTypeCode", zap.Error(err))
		return nil, err
	}

	d.Logger.Info("get single transaction type", zap.Any("transactionType", &tt))
	return &tt, nil
}

func (d *Database) AddTransactionType(tt *models.TransactionType, c context.Context) error {
	d.Logger.Info("starting addding new transaction type into DB", zap.Any("transactionType", tt))
	_, err := d.Client.NewInsert().Model(tt).Exec(c)
	if err != nil {
		d.Logger.Error("error while inserting new transaction type into db", zap.Error(err))
		return err
	}
	d.Logger.Info("transaction type was added successfully", zap.Any("transactionType", tt))
	return nil
}

func (d *Database) EditTransactionType(code string, newTT *models.TransactionType, c context.Context) (*models.TransactionType, error) {
	d.Logger.Info("start updating transaction type", zap.Any("transactionTypeCode", code))
	oldTT, err := d.GetTransactionType(code, c)
	if err != nil {
		return nil, err
	}

	if oldTT == nil {
		return oldTT, nil
	}

	_, err = d.Client.NewUpdate().Model(newTT).OmitZero().Where("transaction_type_code = ?", oldTT.TransactionTypeCode).Returning("*").Exec(c)
	if err != nil {
		d.Logger.Error("error while updating transaction type", zap.Error(err))
		return nil, err
	}

	d.Logger.Info("transaction type was updated successfully", zap.Any("transactionType", newTT))
	return newTT, nil
}

func (d *Database) DeleteTransactionType(code string, c context.Context) (*models.TransactionType, error) {
	d.Logger.Info("start deleting transaction type", zap.Any("transactionTypeCode", code))
	tt, err := d.GetTransactionType(code, c)
	if err != nil {
		return nil, err
	}

	if tt != nil {
		_, err = d.Client.NewDelete().Model(tt).WherePK().Exec(c)
		if err != nil {
			d.Logger.Error("error while deleting transaction type", zap.Error(err))
			return nil, err
		}

		d.Logger.Info("tranaction type was deleted successfully", zap.Any("transaction type", tt))
	}

	return tt, nil
}
