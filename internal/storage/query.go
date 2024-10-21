package storage

import (
	"database/sql"
	"fmt"

	model "altpack-vers-checker/internal/integration/model"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (c *Storage) Close() {
	defer c.db.Close()
}

func (c *Storage) GetPackage(name, branch string) (*model.SqlPack, error) {
	row := c.db.QueryRow(fmt.Sprintf("SELECT id, name, version, release, epoch, changed, branch FROM packages WHERE name='%s' and branch='%s'", name, branch))

	pack := &model.SqlPack{}
	if err := row.Scan(&pack.Id, &pack.Name, &pack.Version, &pack.Release, &pack.Epoch, &pack.Changed, &pack.Branch); err != nil {
		return pack, err
	}

	return pack, nil
}

func (c *Storage) GetPackages(branch string, limit int) ([]model.SqlPack, error) {
	limitq := ""
	if limit > 0 {
		limitq = fmt.Sprintf(" LIMIT %s", limit)
	}
	rows, err := c.db.Query(fmt.Sprintf("SELECT * FROM packages where branch='%s' ORDER BY id DESC%s", branch, limitq))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	packs := []model.SqlPack{}
	for rows.Next() {
		i := model.SqlPack{}
		if err = rows.Scan(&i.Id, &i.Name, &i.Version, &i.Release, &i.Epoch, &i.Changed, &i.Branch); err != nil {
			return nil, err
		}
		packs = append(packs, i)
	}
	return packs, nil
}

func (c *Storage) packExists(name, branch string) (int, error) {
	pack, err := c.GetPackage(name, branch)
	if err != nil {
		return -1, err
	}

	if pack != nil && pack.Id > 0 {
		return pack.Id, nil
	}
	return 0, nil
}

func (c *Storage) InsertPackage(pack *model.SqlPack) (int, error) {
	//INSERT INTO table_name (column1, column2, column3, ...) VALUES (value1, value2, value3, ...);
	exid, err := c.packExists(pack.Name, pack.Branch)
	if exid > 0 {
		return exid, nil
	}

	res, err := c.db.Exec(`INSERT INTO packages
		(name, version, release, epoch, changed, branch) 
		VALUES (?,?,?,?,?,?);`,
		&pack.Name, &pack.Version, &pack.Release, &pack.Epoch, &pack.Changed, pack.Branch)
	if err != nil {
		return 0, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), nil
}

func (c *Storage) DeletePackageById(id int) error {
	_, err := c.db.Exec("DELETE FROM packages WHERE id=$1;", id)

	if err != nil {
		return err
	}

	return nil
}

func (c *Storage) DeletePackageByName(name, branch string) error {
	_, err := c.db.Exec(fmt.Sprint("DELETE FROM packages WHERE name='%s' and branch='%s';", name, branch))

	if err != nil {
		return err
	}

	return nil
}
