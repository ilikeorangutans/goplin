package database

import "github.com/jmoiron/sqlx"

type SyncItems struct {
	db *sqlx.DB
}
