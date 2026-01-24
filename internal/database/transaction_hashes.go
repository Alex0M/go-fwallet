package database

import (
	"context"
	"database/sql"
	"fmt"
	"go-fwallet/internal/middleware"
	"go-fwallet/internal/models"
	"net/http"
)

func (d *Database) GetSingleTransactionHash(hash string, c context.Context) (*models.TransactionHash, error) {
	var th models.TransactionHash
	err := d.Client.NewSelect().Model(&th).Where("hash = ?", hash).Scan(c)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, middleware.NewHttpError("transaction hash not found", fmt.Sprintf("transactionHash:%s", hash), http.StatusNotFound)
		}
		return nil, err
	}

	return &th, nil
}

func (d *Database) AddTransactionHash(th *models.TransactionHash, c context.Context) error {
	_, err := d.Client.NewInsert().Model(th).Exec(c)
	if err != nil {
		return err
	}
	return nil
}
