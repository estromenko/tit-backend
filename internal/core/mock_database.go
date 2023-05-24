package core

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func NewMockDatabase() (*Database, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	bunDB := bun.NewDB(db, pgdialect.New())

	return bunDB, mock
}
