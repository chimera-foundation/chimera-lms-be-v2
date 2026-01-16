package platform

import (
	"database/sql"
	"github.com/uptrace/bun"
	pgdialect "github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	// "github.com/uptrace/bun/extra/bundebug" // You might need to go get this
)

func NewPostgresDB(dsn string) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	// uncomment to debug query
	// db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	return db
}