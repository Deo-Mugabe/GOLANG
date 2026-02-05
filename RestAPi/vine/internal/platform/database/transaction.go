package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// TxFunc is a function that runs within a transaction
type TxFunc func(*sqlx.Tx) error

// WithTransaction executes a function within a transaction
func (db *DB) WithTransaction(ctx context.Context, fn TxFunc) error {
	return db.withTransaction(ctx, nil, fn)
}

// WithTransactionIsolation executes a function within a transaction with specific isolation level
func (db *DB) WithTransactionIsolation(ctx context.Context, level sql.IsolationLevel, fn TxFunc) error {
	opts := &sql.TxOptions{
		Isolation: level,
	}
	return db.withTransaction(ctx, opts, fn)
}

// withTransaction is the internal transaction executor
func (db *DB) withTransaction(ctx context.Context, opts *sql.TxOptions, fn TxFunc) error {
	tx, err := db.DB.BeginTxx(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer rollback (will be no-op if committed)
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after rollback
		}
	}()

	// Execute function
	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %w", err, rbErr)
		}
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Transaction isolation level helpers
const (
	IsolationReadUncommitted = sql.LevelReadUncommitted
	IsolationReadCommitted   = sql.LevelReadCommitted
	IsolationRepeatableRead  = sql.LevelRepeatableRead
	IsolationSerializable    = sql.LevelSerializable
)
