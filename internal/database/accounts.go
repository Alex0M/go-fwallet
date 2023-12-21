package database

import (
	"context"
	"database/sql"
	"fmt"
	"go-fwallet/internal/helpers"
	"go-fwallet/internal/middleware"
	"go-fwallet/internal/models"
	"net/http"
)

func (d *Database) GetAccounts(c context.Context) ([]models.Account, error) {
	var accounts []models.Account
	err := d.Client.NewSelect().Model(&accounts).Scan(c)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (d *Database) GetSingleAccount(idStr string, c context.Context) (*models.Account, error) {
	id, err := helpers.StringToInt(idStr)
	if err != nil {
		return nil, middleware.NewHttpError("cannot conver account id to int", err.Error(), http.StatusBadRequest)
	}

	var a models.Account
	err = d.Client.NewSelect().Model(&a).Where("id = ?", id).Scan(c)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, middleware.NewHttpError("account not found", fmt.Sprintf("accountID:%d", id), http.StatusNotFound)
		}
		return nil, err
	}

	return &a, nil
}

func (d *Database) AddAccount(a *models.Account, c context.Context) error {
	_, err := d.Client.NewInsert().Model(a).Exec(c)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) EditAccount(idStr string, newAccount *models.Account, c context.Context) (*models.Account, error) {
	oldAcc, err := d.GetSingleAccount(idStr, c)
	if err != nil {
		return nil, err
	}

	_, err = d.Client.NewUpdate().Model(newAccount).OmitZero().Where("id = ?", oldAcc.ID).Returning("*").Exec(c)
	if err != nil {
		return nil, err
	}

	return newAccount, nil
}

func (d *Database) DeleteAccount(idStr string, c context.Context) error {
	a, err := d.GetSingleAccount(idStr, c)
	if err != nil {
		return err
	}

	_, err = d.Client.NewDelete().Model(a).WherePK().Exec(c)
	if err != nil {
		return err
	}

	return nil
}
