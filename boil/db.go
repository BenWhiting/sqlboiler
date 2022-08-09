package boil

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v4"
)

// Executor can perform SQL queries.
type Executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (pgx.Rows, error)
	QueryRow(query string, args ...interface{}) pgx.Row
}

// ContextExecutor can perform SQL queries with context
type ContextExecutor interface {
	Executor

	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) pgx.Row
}

// Transactor can commit and rollback, on top of being able to execute queries.
type Transactor interface {
	Commit() error
	Rollback() error

	Executor
}

type result struct {
	index int64
	rows  int64
}

func (r result) LastInsertId() (int64, error) {
	return r.index, nil
}

func (r result) RowsAffected() (int64, error) {
	return r.rows, nil
}

// Tx is an interface that describes a transaction.
type Tx interface {
	// Commit commits the transaction.
	Commit() error
	// Rollback rollsback the transaction.
	Rollback() error
}

// Beginner begins transactions.
type Beginner interface {
	Begin() (*BoilerPgxWrap, error)
}

// Begin a transaction with the current global database handle.
func Begin() (Transactor, error) {
	creator, ok := currentDB.(Beginner)
	if !ok {
		panic("database does not support transactions")
	}
	return creator.Begin()
}

// ContextTransactor can commit and rollback, on top of being able to execute
// context-aware queries.
type ContextTransactor interface {
	Commit() error
	Rollback() error

	ContextExecutor
}

// ContextBeginner allows creation of context aware transactions with options.
type ContextBeginner interface {
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
}

// BeginTx begins a transaction with the current global database handle.
func BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	creator, ok := currentDB.(ContextBeginner)
	if !ok {
		panic("database does not support context-aware transactions")
	}

	return creator.BeginTx(ctx, opts)
}

type BoilerPgxWrap struct {
	pgx.Tx
}

func NewPgxWrapper(tx pgx.Tx) *BoilerPgxWrap {
	return &BoilerPgxWrap{tx}
}

func (w *BoilerPgxWrap) Commit() error {
	return w.Tx.Commit(context.Background())
}

func (w *BoilerPgxWrap) Rollback() error {
	return w.Tx.Rollback(context.Background())
}

func (w *BoilerPgxWrap) Exec(query string, args ...interface{}) (sql.Result, error) {
	ct, err := w.Tx.Exec(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	return &result{rows: ct.RowsAffected(), index: 0}, nil
}

func (w *BoilerPgxWrap) Query(query string, args ...interface{}) (pgx.Rows, error) {
	return w.Tx.Query(context.Background(), query, args...)
}

func (w *BoilerPgxWrap) QueryRow(query string, args ...interface{}) pgx.Row {
	return w.Tx.QueryRow(context.Background(), query, args...)
}

func (w *BoilerPgxWrap) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	ct, err := w.Tx.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &result{rows: ct.RowsAffected(), index: 0}, nil
}

func (w *BoilerPgxWrap) QueryContext(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return w.Tx.Query(ctx, query, args...)
}

func (w *BoilerPgxWrap) QueryRowContext(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return w.Tx.QueryRow(ctx, query, args...)
}
