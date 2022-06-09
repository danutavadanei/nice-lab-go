package mysql

import (
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewConnection(cfg mysql.Config) (*sqlx.DB, error) {
	return sqlx.Open("mysql", cfg.FormatDSN())
}
