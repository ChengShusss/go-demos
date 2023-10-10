package model

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"
)

// Register models used in Sqlite3
// Key 		: table name
// Value 	: Struct of data
var RegisterModels = map[string]interface{}{
	"infos": Info{},
}

// TODO fix dbMap init problem

var (
	dbPath = "./build/info.db"
)

var dbMap *gorp.DbMap

func init() {
	// should not initdb everytime
	GetDbMap(dbPath)
	InitDb()
}

func GetDbMap(path string) (*gorp.DbMap, error) {
	if dbMap != nil {
		return dbMap, nil
	}

	if path == "" {
		path = dbPath
	}
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	dbMap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	return dbMap, nil
}

// This should not be called everytime start in production env
func InitDb() error {

	dbmap, err := GetDbMap("")
	if err != nil {
		return err
	}

	for name, model := range RegisterModels {
		dbmap.AddTableWithName(model, name).SetKeys(true, "Id")
	}

	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		return err
	}
	return nil
}
