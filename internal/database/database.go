package database

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Database struct {
	Client *bun.DB
}

func Init(dsn string) *Database {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	return &Database{
		Client: db,
	}
}

func (d *Database) Ping() error {
	return d.Client.Ping()
}
