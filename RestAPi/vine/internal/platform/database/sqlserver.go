package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/denisenkom/go-mssqldb" // SQL Server driver
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

// DB wraps sqlx.DB with additional functionality
type DB struct {
	*sqlx.DB
	config *Config
	logger zerolog.Logger
}

// New creates a new database connection
func New(config *Config, logger zerolog.Logger) (*DB, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid database config: %w", err)
	}

	logger.Info().
		Str("host", config.Host).
		Int("port", config.Port).
		Str("database", config.Database).
		Msg("connecting to SQL Server")

	// Connect with retry logic
	var db *sqlx.DB
	var err error

	for attempt := 1; attempt <= config.RetryAttempts; attempt++ {
		db, err = connect(config)
		if err == nil {
			break
		}

		if attempt < config.RetryAttempts {
			logger.Warn().
				Err(err).
				Int("attempt", attempt).
				Int("max_attempts", config.RetryAttempts).
				Msg("failed to connect, retrying...")
			time.Sleep(config.RetryInterval)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect after %d attempts: %w",
			config.RetryAttempts, err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	logger.Info().Msg("successfully connected to SQL Server")

	return &DB{
		DB:     db,
		config: config,
		logger: logger,
	}, nil
}

// connect establishes database connection
func connect(config *Config) (*sqlx.DB, error) {
	dsn := config.DSN()

	db, err := sqlx.Connect("sqlserver", dsn)
	if err != nil {
		return nil, err
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	db.logger.Info().Msg("closing database connection")
	return db.DB.Close()
}

// Ping checks database connectivity
func (db *DB) Ping(ctx context.Context) error {
	return db.DB.PingContext(ctx)
}

// Stats returns database statistics
func (db *DB) Stats() sql.DBStats {
	return db.DB.Stats()
}

// WithTimeout creates a context with query timeout
func (db *DB) WithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, db.config.QueryTimeout)
}

// ExecContext executes a query with logging
func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()

	result, err := db.DB.ExecContext(ctx, query, args...)

	duration := time.Since(start)
	db.logQuery(query, duration, err)

	return result, err
}

// QueryContext executes a query with logging
func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()

	rows, err := db.DB.QueryContext(ctx, query, args...)

	duration := time.Since(start)
	db.logQuery(query, duration, err)

	return rows, err
}

// QueryRowContext executes a query row with logging
func (db *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	start := time.Now()

	row := db.DB.QueryRowContext(ctx, query, args...)

	duration := time.Since(start)
	db.logQuery(query, duration, nil)

	return row
}

// logQuery logs query execution
func (db *DB) logQuery(query string, duration time.Duration, err error) {
	event := db.logger.Debug().
		Str("query", query).
		Dur("duration", duration)

	if err != nil {
		event.Err(err).Msg("query failed")
	} else {
		event.Msg("query executed")
	}
}
