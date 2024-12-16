package asyncdataloader

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func NewDB() (*sql.DB, error) {
	if db != nil {
		return db, nil
	}
	conf := mysql.NewConfig()
	conf.User = "root"
	conf.DBName = "asyncmap_test"
	var err error
	db, err = sql.Open("mysql", conf.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return db, nil
}
