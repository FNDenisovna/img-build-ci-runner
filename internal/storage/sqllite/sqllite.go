package sqllite

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var (
	dbpath = "~/.local/share/img-build-ci-runner"
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
	if len(args) > 0 && args[0] != "" {
		dbpath = args[0]
	}
	log.Printf("DB path path: %s\n", dbpath)

	dbpath, err = filepath.Abs(dbpath)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Can't parse DB directory. Path %s. Error: %v\n", dbpath, err))
	}
	log.Printf("Parsed DB path path: %s\n", dbpath)

	if _, err = os.Stat(dbpath); os.IsNotExist(err) {
		err = os.Mkdir(dbpath, 0750)
		if err != nil && !os.IsExist(err) {
			log.Fatalf(fmt.Sprintf("Can't create DB directory. Path %s. Error: %v\n", dbpath, err))
		}
	}

	dbpath = filepath.Join(dbpath, dbname)

	db, err = sql.Open("sqlite3", dbpath)
	if err != nil {
		return nil, err
	}
	if _, err = db.Exec(create); err != nil {
		return nil, err
	}
	return db, nil
}
