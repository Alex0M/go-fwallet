package database

import (
	"context"
	"database/sql"
	"go-fwallet/internal/helpers"
	"go-fwallet/internal/models"

	"go.uber.org/zap"
)

func (d *Database) GetAccounts(c context.Context) ([]models.Account, error) {
	var accounts []models.Account
	err := d.Client.NewSelect().Model(&accounts).Scan(c)
	if err != nil {
		d.Logger.Error("error while getting all accounts", zap.Error(err))
		return nil, err
	}

	return accounts, nil
}

func (d *Database) GetSingleAccount(idStr string, c context.Context) (*models.Account, error) {
	id, err := helpers.StringToInt(idStr)
	if err != nil {
		d.Logger.Error(ErrValidateID.Error(), zap.Error(err))
		return nil, ErrValidateID
	}

	var a models.Account
	err = d.Client.NewSelect().Model(&a).Where("id = ?", id).Scan(c)
	if err != nil {
		if err == sql.ErrNoRows {
			d.Logger.Info("account not found", zap.Any("accountID", id))
			return nil, nil
		}
		d.Logger.Error("error while getting account", zap.Error(err))
		return nil, err
	}

	d.Logger.Info("get single account", zap.Any("account", &a))
	return &a, nil
}

func (d *Database) AddAccount(a *models.Account, c context.Context) error {
	_, err := d.Client.NewInsert().Model(a).Exec(c)
	if err != nil {
		d.Logger.Error("error while inserting new account into db", zap.Error(err))
		return err
	}
	d.Logger.Info("account was added successfully", zap.Any("account", a))
	return nil
}

func (d *Database) EditAccount(idStr string, newAccount *models.Account, c context.Context) (*models.Account, error) {
	d.Logger.Info("start updating account", zap.Any("accounID", idStr))
	oldAcc, err := d.GetSingleAccount(idStr, c)
	if err != nil {
		return nil, err
	}

	if oldAcc == nil {
		return oldAcc, nil
	}

	_, err = d.Client.NewUpdate().Model(newAccount).OmitZero().Where("id = ?", oldAcc.ID).Returning("*").Exec(c)
	if err != nil {
		d.Logger.Error("error while updating account", zap.Error(err))
		return nil, err
	}

	d.Logger.Info("account was updated successfully", zap.Any("account", newAccount))
	return newAccount, nil
}

func (d *Database) DeleteAccount(idStr string, c context.Context) (*models.Account, error) {
	d.Logger.Info("start deleting account", zap.Any("accounID", idStr))
	a, err := d.GetSingleAccount(idStr, c)
	if err != nil {
		return nil, err
	}

	if a != nil {
		_, err = d.Client.NewDelete().Model(a).WherePK().Exec(c)
		if err != nil {
			d.Logger.Error("error while deleting accont", zap.Error(err))
			return nil, err
		}

		d.Logger.Info("account was deleted successfully", zap.Any("account", a))
	}

	return a, nil
}
