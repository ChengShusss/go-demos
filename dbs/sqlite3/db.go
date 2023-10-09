package sqlite3

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"
)


func InitDb(dbPath string) (*gorp.DbMap, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	// add a table, setting the table name to 'posts' and
	// specifying that the Id property is an auto incrementing PK
	for name, model := range RegisterModels {
		dbmap.AddTableWithName(model, name).SetKeys(true, "Id")
	}

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		return nil, err
	}
	return dbmap, nil
}
