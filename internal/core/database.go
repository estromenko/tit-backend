package core

import (
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func NewDatabase(conf *Config) *bun.DB {
	sqlDB := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithAddr(fmt.Sprintf("%s:%d", conf.PgHost, conf.PgPort)),
		pgdriver.WithUser(conf.PgUser),
		pgdriver.WithPassword(conf.PgPassword),
		pgdriver.WithDatabase(conf.PgName),
		pgdriver.WithInsecure(true),
	))

	db := bun.NewDB(sqlDB, pgdialect.New())

	return db
}
