package database

import (
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sqlx.DB
}

func (d *Database) Migrate() error {
	createVersionsTable := `
	  CREATE TABLE IF NOT EXISTS MIGRATIONS (id INTEGER NOT NULL PRIMARY KEY, applied_at INTEGER NOT NULL)
	`
	_, err := d.db.Exec(createVersionsTable)
	if err != nil {
		return err
	}

	getMaxVersion := `
	  SELECT MAX(id) FROM migrations
	`
	row := d.db.QueryRow(getMaxVersion)
	maxMigration := len(migrations)
	maxVersion := 0
	row.Scan(&maxVersion)
	log.Printf("max version from db is %d, max migration is %d", maxVersion, maxMigration)
	if maxVersion < maxMigration {
		for i, query := range migrations {
			if i < maxVersion {
				log.Printf("skipping %d", i)
				continue
			}

			tx := d.db.MustBegin()
			log.Printf("applying %d %s", i+1, query)
			tx.MustExec(query)
			tx.MustExec("insert into migrations (id, applied_at) values ($1, datetime())", i+1)
			tx.Commit()
		}
	}

	return nil
}

func (d *Database) SyncItems() *SyncItems {
	return &SyncItems{
		db: d.db,
	}
}

func (d *Database) Close() error {
	return d.db.Close()
}

func OpenDatabase() (*Database, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(usr.HomeDir, ".config", "goplin")
	os.MkdirAll(path, 0755)
	log.Printf("opening %s", path)
	db, err := sqlx.Open("sqlite3", filepath.Join(path, "goplin.db"))
	if err != nil {
		return nil, err
	}

	return &Database{
		db: db,
	}, nil
}

func (d *Database) Begin() *sqlx.Tx {
	return d.db.MustBegin()
}

func (d *Database) HasSyncItem(id string) (bool, error) {

	hasID := false
	err := d.db.Get(&hasID, "select count(*) from sync_items where id = $1", id)
	if err != nil {
		return false, err
	}

	return hasID, nil
}
