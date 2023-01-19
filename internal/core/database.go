package core

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func NewDatabase(conf *Config, log *zerolog.Logger) *bun.DB {
	sqlDB := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithAddr(fmt.Sprintf("%s:%d", conf.PgHost, conf.PgPort)),
		pgdriver.WithUser(conf.PgUser),
		pgdriver.WithPassword(conf.PgPassword),
		pgdriver.WithDatabase(conf.PgName),
		pgdriver.WithInsecure(true),
	))

	db := bun.NewDB(sqlDB, pgdialect.New())

	if err := db.Ping(); err != nil {
		log.Err(err).Msg("Database ping failed")
		os.Exit(1)
	}

	return db
}
