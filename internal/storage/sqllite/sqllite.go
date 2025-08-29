package sqllite

import (
	"database/sql"
	"img-build-ci-runner/internal/resources"

	_ "github.com/mattn/go-sqlite3"
)

var (
	dbname = "packages.db"
)

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

// package_id int primary key auto_increment, package_name VARCHAR(128) NOT NULL, version VARCHAR(128) NOT NULL, release VARCHAR(128) NOT NULL, epoch int, changed datetime
// args[0] - dbpath
func New(args ...string) (db *sql.DB, err error) {
	var dbpath string
	if len(args) > 0 && args[0] != "" {
		dbpath = args[0]
	}

	dbpath = resources.ManageResources(dbpath, dbname)

	db, err = sql.Open("sqlite3", dbpath)
	if err != nil {
		return nil, err
	}
	if _, err = db.Exec(create); err != nil {
		return nil, err
	}
	return db, nil
}
