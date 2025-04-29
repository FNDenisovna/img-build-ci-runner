package sqllite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const dbpath = "/var/lib/img-build-ci-runner/packages.db"

const create string = `
  CREATE TABLE IF NOT EXISTS packages (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  name VARCHAR(128) NOT NULL,
  version VARCHAR(128) NOT NULL,
  release VARCHAR(128) NOT NULL,
  epoch INTEGER,
  changed DATETIME NOT NULL,
  branch VARCHAR(128) NOT NULL
  );`

//package_id int primary key auto_increment, package_name VARCHAR(128) NOT NULL, version VARCHAR(128) NOT NULL, release VARCHAR(128) NOT NULL, epoch int, changed datetime

func New() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", dbpath)
	if err != nil {
		return nil, err
	}
	if _, err = db.Exec(create); err != nil {
		return nil, err
	}
	return db, nil
}
