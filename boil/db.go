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

// Tx is an interface that describes a transaction.
type Tx interface {
	// Commit commits the transaction.
	Commit() error
	// Rollback rollsback the transaction.
	Rollback() error
}

// Beginner begins transactions.
type Beginner interface {
	Begin() (*PgxWrapper, error)
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

type Result struct {
	index int64
	rows  int64
}

func (r Result) LastInsertId() (int64, error) {
	return r.index, nil
}

func (r Result) RowsAffected() (int64, error) {
	return r.rows, nil
}

type PgxWrapper struct {
	tx pgx.Tx
}

func NewPgxWrapper(tx pgx.Tx) *PgxWrapper {
	return &PgxWrapper{tx: tx}
}

func (w *PgxWrapper) Commit() error {
	return w.tx.Commit(context.Background())
}

func (w *PgxWrapper) Rollback() error {
	return w.tx.Rollback(context.Background())
}

func (w *PgxWrapper) Exec(query string, args ...interface{}) (sql.Result, error) {
	ct, err := w.tx.Exec(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	return Result{index: 0, rows: ct.RowsAffected()}, nil
}

func (w *PgxWrapper) QueryRow(query string, args ...interface{}) pgx.Row {
	return w.tx.QueryRow(context.Background(), query, args...)
}

func (w *PgxWrapper) Query(query string, args ...interface{}) (pgx.Rows, error) {
	return w.tx.Query(context.Background(), query, args...)
}

func (w *PgxWrapper) ExecContext(_ context.Context, query string, args ...interface{}) (sql.Result, error) {
	return w.Exec(query, args...)
}

func (w *PgxWrapper) QueryContext(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return w.tx.Query(ctx, query, args...)
}

func (w *PgxWrapper) QueryRowContext(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return w.tx.QueryRow(ctx, query, args...)
}
