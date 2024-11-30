package sqlutil

import (
	"database/sql"
	"fmt"
)

type SQLHandle interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type DatabaseTransaction interface {
	Rollback() error
	Commit() error
	SQLHandle // Embedding SQLHandle interface
}

type databaseTransaction struct {
	tx *sql.Tx
}

func (t *databaseTransaction) Rollback() error {
	return t.tx.Rollback()
}

func (t *databaseTransaction) Commit() error {
	return t.tx.Commit()
}

func (t *databaseTransaction) QueryRow(query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRow(query, args...)
}

func (t *databaseTransaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.tx.Exec(query, args...)
}

type databaseTransactionCreator struct {
	db *sql.DB
}

type DatabaseTransactionCreator interface {
	Begin() (*databaseTransaction, error)
}

func (s *databaseTransactionCreator) Begin() (*databaseTransaction, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	return &databaseTransaction{tx: tx}, nil
}

func NewDatabaseTransactionCreator(db *sql.DB) *databaseTransactionCreator {
	return &databaseTransactionCreator{db}
}

func GetExecer(db *sql.DB, tx DatabaseTransaction) (SQLHandle, error) {
	if tx == nil {
		return db, nil
	}

	if tx, ok := tx.(*databaseTransaction); ok {
		return tx, nil
	}
	return nil, fmt.Errorf("DatabaseTransaction not supported")
}
