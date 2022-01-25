package db

import (
	"database/sql"
	"github.com/jackc/pgx"
)

type Connector interface {
	GetDb() *sql.DB
	Connect(user string, password string, dbAddress string, dbName string) error
	Close() error
	ExecTransaction(query string, args ...interface{}) error
}

type Transaction interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type TxFn func(tx *pgx.Tx) error

func WithTransaction(db *pgx.ConnPool, fn TxFn) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			err = tx.Rollback()
			panic(p)
		} else if err != nil {
			err = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
