package db

import (
	"database/sql"
	"payment/internal/config"
)

type Database struct {
	Db *sql.DB
}

func NewDatabase(config *config.Config) *Database {

}
