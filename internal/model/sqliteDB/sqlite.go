package sqliteDB

import (
	"embed"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

func Connect(migrations embed.FS, dbName string) *sqlx.DB {
	db := sqlx.MustConnect("sqlite", dbName)

	// run goose migrations
	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		panic(err)
	}
	if err := goose.Up(db.DB, "migrations"); err != nil {
		panic(err)
	}

	return db
}
