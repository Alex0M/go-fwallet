package database

import (
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"go.uber.org/zap"
)

// Need to delete after refactoing all handlers
var DB *bun.DB

var ErrValidateID = fmt.Errorf("cannot conver id to int")

type Database struct {
	Client *bun.DB
	Logger *zap.Logger
}

func Init(dsn string, l *zap.Logger) *Database {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	return &Database{
		Client: db,
		Logger: l,
	}
}

func (d *Database) Ping() error {
	return d.Client.Ping()
}

func (d *Database) GetErrValidateID() error {
	return ErrValidateID
}
