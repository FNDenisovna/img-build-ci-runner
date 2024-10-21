package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "root"
	password = "root"
	hostname = "127.0.0.1:3306"
	dbname   = "altpack"
)

func New() *sql.DB {
	connStr := fmt.Sprintf("%v:%v@tcp(%v)/", username, password, hostname)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbname)
	if err != nil {
		panic(err)
	}
	db.Close()

	db, err = sql.Open("mysql", connStr+dbname)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	return db
}

func (db *sql.DB) Close() {
	defer db.Close()
}
