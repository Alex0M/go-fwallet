package database

import (
	"context"
	"database/sql"
	"go-fwallet/internal/helpers"
	"go-fwallet/internal/models"
	"time"

	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

func (d *Database) GetAllTransaction(tp *models.TransactionAPIRequestParams, c context.Context) ([]models.TransactionAPI, error) {
	d.Logger.Info("statring getting all transactions", zap.Any("params", tp))
	var t []models.TransactionAPI
	q := d.Client.NewSelect().Model(&t).ColumnExpr("t.id, t.transaction_date, t.amount, t.description, t.transaction_type_code").
		ColumnExpr("a.name AS account").
		ColumnExpr("c.name AS sub_category").
		ColumnExpr("cc.name AS category").
		Join("LEFT JOIN categories AS c ON t.category_id = c.id").
		Join("LEFT JOIN categories AS cc ON c.parent_id = cc.id").
		Join("LEFT JOIN accounts AS a ON t.account_id = a.id")

	q = queryByParams(q, tp)
	err := q.Order("t.transaction_date").Scan(c)
	if err != nil {
		d.Logger.Error("error while getting all transactions", zap.Error(err))
		return nil, err
	}

	d.Logger.Info("got all transactions")
	return t, nil
}

func (d *Database) GetSumOfTransactionsByCategory(tp *models.TransactionAPIRequestParams, c context.Context) ([]models.SumTransactionByCategory, error) {
	d.Logger.Info("statring getting sum of transactions by category", zap.Any("params", tp))
	var stbc []models.SumTransactionByCategory

	q := d.Client.NewSelect().Model(&stbc).ColumnExpr("sum(amount) as amount").
		ColumnExpr("c.name as category").
		Join("LEFT JOIN categories AS cc ON t.category_id = cc.id").
		Join("LEFT JOIN categories AS c ON cc.parent_id = c.id")

	q = queryByParams(q, tp)
	err := q.Group("category").Scan(c)
	if err != nil {
		d.Logger.Error("error while getting sum of transactions by category", zap.Error(err))
		return nil, err
	}

	d.Logger.Info("got sum of transactions by category")
	return stbc, nil
}

func (d *Database) GetTransaction(idStr string, c context.Context) (*models.Transaction, error) {
	id, err := helpers.StringToInt(idStr)
	if err != nil {
		d.Logger.Error(ErrValidateID.Error(), zap.Error(err))
		return nil, ErrValidateID
	}

	var t models.Transaction
	err = d.Client.NewSelect().Model(&t).Where("id = ?", id).Scan(c)
	if err != nil {
		if err == sql.ErrNoRows {
			d.Logger.Info("transaction not found", zap.Any("TransactionID", id))
			return nil, nil
		}
		d.Logger.Error("error while getting transaction", zap.Error(err))
		return nil, err
	}

	d.Logger.Info("get single transaction", zap.Any("transaction", t))
	return &t, nil
}

func (d *Database) AddTransaction(t *models.Transaction, c context.Context) error {
	d.Logger.Info("starting addding new transaction intoDB", zap.Any("transaction", t))
	_, err := d.Client.NewInsert().Model(t).Exec(c)
	if err != nil {
		d.Logger.Error("error while inserting new transaction into db", zap.Error(err))
		return err
	}
	d.Logger.Info("transaction was added successfully", zap.Any("account", t))
	return nil
}

func (d *Database) EditTransaction(idStr string, newTransaction *models.Transaction, c context.Context) (*models.Transaction, error) {
	d.Logger.Info("start updating transaction", zap.Any("TransactionID", idStr))
	oldT, err := d.GetTransaction(idStr, c)
	if err != nil {
		return nil, err
	}

	if oldT == nil {
		return oldT, nil
	}

	_, err = d.Client.NewUpdate().Model(newTransaction).OmitZero().Where("id = ?", oldT.ID).Returning("*").Exec(c)
	if err != nil {
		d.Logger.Error("error while updating account", zap.Error(err))
		return nil, err
	}

	d.Logger.Info("account was updated successfully", zap.Any("account", newTransaction))
	return newTransaction, nil
}

func (d *Database) DeleteTransaction(idStr string, c context.Context) (*models.Transaction, error) {
	d.Logger.Info("start deletting transaction", zap.Any("TransactionID", idStr))
	t, err := d.GetTransaction(idStr, c)
	if err != nil {
		return nil, err
	}

	if t != nil {
		_, err = d.Client.NewDelete().Model(t).WherePK().Exec(c)
		if err != nil {
			d.Logger.Error("error while deleting transaction", zap.Error(err))
			return nil, err
		}

		d.Logger.Info("transaction was deleted successfully", zap.Any("account", t))
	}

	return t, nil

}

func (d *Database) UpdateTransactionCategories(tap *models.TransactionAPIRequestParams, c context.Context) (int64, error) {
	d.Logger.Info("start updating transactions categories", zap.Any("Params", tap))
	res, err := d.Client.NewUpdate().
		Table("transactions", "bank_statement_descriptions").
		SetColumn("category_id", "bank_statement_descriptions.category_id").
		Set("updated_at = ?", time.Now()).
		Where("transactions.description = bank_statement_descriptions.bank_statement_description and transactions.transaction_date >= ? and transactions.transaction_date <= ? and transactions.transaction_type_code = ? and transactions.category_id = 1", tap.GteDate, tap.LteDate, tap.TransactionTypeCode).
		Exec(c)

	if err != nil {
		return 0, err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return ra, nil
}

func queryByParams(q *bun.SelectQuery, tp *models.TransactionAPIRequestParams) *bun.SelectQuery {
	if tp.IsTransactionTypeCodeSpecify {
		q.Where("t.transaction_type_code = ?", tp.TransactionTypeCode)
	}
	if tp.IsGteDateSpecify && tp.IsLteDateSpecify {
		q.Where("t.transaction_date >= ? and t.transaction_date <= ?", tp.GteDate, tp.LteDate)
	} else if tp.IsGteDateSpecify {
		q.Where("t.transaction_date >= ?", tp.GteDate)
	} else if tp.IsLteDateSpecify {
		q.Where("t.transaction_date <= ?", tp.LteDate)
	}

	return q
}
