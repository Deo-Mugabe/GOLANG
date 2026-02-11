package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Deo-Mugabe/GOLANG/internal/platform/database"
	"github.com/Deo-Mugabe/GOLANG/internal/scheduler/domain"
)

// configDB is the database model for scheduler_config table
type configDB struct {
	ConfigName    string       `db:"config_name"`
	Enabled       bool         `db:"enabled"`
	IntervalMin   int          `db:"interval_minutes"`
	LastStartTime sql.NullTime `db:"last_start_time"`
	LastStopTime  sql.NullTime `db:"last_stop_time"`
	LastRunTime   sql.NullTime `db:"last_run_time"`
	NextRunTime   sql.NullTime `db:"next_run_time"`
	StartFromTime time.Time    `db:"start_from_time"`
	CreatedAt     time.Time    `db:"created_at"`
	UpdatedAt     sql.NullTime `db:"updated_at"`
}

type configRepo struct {
	db *database.DB
}

// NewConfigRepository creates a new config repository
func NewConfigRepository(db *database.DB) domain.ConfigRepository {
	return &configRepo{db: db}
}

// Get retrieves scheduler config by name
func (r *configRepo) Get(ctx context.Context, name string) (*domain.SchedulerConfig, error) {
	query := `
        SELECT config_name, enabled, interval_minutes, last_start_time, 
               last_stop_time, last_run_time, next_run_time, start_from_time,
               created_at, updated_at
        FROM scheduler_config
        WHERE config_name = @p1
    `

	var dbModel configDB
	err := r.db.GetContext(ctx, &dbModel, query, name)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("scheduler config not found: %s", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get scheduler config: %w", err)
	}

	return toConfigDomain(&dbModel), nil
}

// List retrieves all scheduler configs
func (r *configRepo) List(ctx context.Context) ([]*domain.SchedulerConfig, error) {
	query := `
        SELECT config_name, enabled, interval_minutes, last_start_time, 
               last_stop_time, last_run_time, next_run_time, start_from_time,
               created_at, updated_at
        FROM scheduler_config
        ORDER BY config_name
    `

	var dbModels []configDB
	err := r.db.SelectContext(ctx, &dbModels, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list scheduler configs: %w", err)
	}

	configs := make([]*domain.SchedulerConfig, len(dbModels))
	for i, db := range dbModels {
		configs[i] = toConfigDomain(&db)
	}

	return configs, nil
}

// Create creates a new scheduler config
func (r *configRepo) Create(ctx context.Context, config *domain.SchedulerConfig) error {
	query := `
        INSERT INTO scheduler_config (
            config_name, enabled, interval_minutes, start_from_time, created_at
        ) VALUES (@p1, @p2, @p3, @p4, @p5)
    `

	_, err := r.db.ExecContext(ctx, query,
		config.ConfigName,
		config.Enabled,
		config.IntervalMin,
		config.StartFromTime,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to create scheduler config: %w", err)
	}

	return nil
}

// Update updates scheduler config
func (r *configRepo) Update(ctx context.Context, config *domain.SchedulerConfig) error {
	query := `
        UPDATE scheduler_config
        SET enabled = @p1,
            interval_minutes = @p2,
            last_start_time = @p3,
            last_stop_time = @p4,
            last_run_time = @p5,
            next_run_time = @p6,
            start_from_time = @p7,
            updated_at = @p8
        WHERE config_name = @p9
    `

	result, err := r.db.ExecContext(ctx, query,
		config.Enabled,
		config.IntervalMin,
		config.LastStartTime,
		config.LastStopTime,
		config.LastRunTime,
		config.NextRunTime,
		config.StartFromTime,
		time.Now(),
		config.ConfigName,
	)
	if err != nil {
		return fmt.Errorf("failed to update scheduler config: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("scheduler config not found: %s", config.ConfigName)
	}

	return nil
}

// Delete deletes scheduler config
func (r *configRepo) Delete(ctx context.Context, name string) error {
	result, err := r.db.ExecContext(ctx,
		"DELETE FROM scheduler_config WHERE config_name = @p1", name)
	if err != nil {
		return fmt.Errorf("failed to delete scheduler config: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("scheduler config not found: %s", name)
	}

	return nil
}

// UpdateEnabled updates only the enabled flag
func (r *configRepo) UpdateEnabled(ctx context.Context, name string, enabled bool) error {
	result, err := r.db.ExecContext(ctx,
		"UPDATE scheduler_config SET enabled = @p1, updated_at = @p2 WHERE config_name = @p3",
		enabled, time.Now(), name)
	if err != nil {
		return fmt.Errorf("failed to update enabled flag: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("scheduler config not found: %s", name)
	}

	return nil
}

// UpdateInterval updates only the interval
func (r *configRepo) UpdateInterval(ctx context.Context, name string, intervalMin int) error {
	result, err := r.db.ExecContext(ctx,
		"UPDATE scheduler_config SET interval_minutes = @p1, updated_at = @p2 WHERE config_name = @p3",
		intervalMin, time.Now(), name)
	if err != nil {
		return fmt.Errorf("failed to update interval: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("scheduler config not found: %s", name)
	}

	return nil
}

// UpdateLastRunTime updates last run time
func (r *configRepo) UpdateLastRunTime(ctx context.Context, name string, t time.Time) error {
	result, err := r.db.ExecContext(ctx,
		"UPDATE scheduler_config SET last_run_time = @p1, updated_at = @p2 WHERE config_name = @p3",
		t, time.Now(), name)
	if err != nil {
		return fmt.Errorf("failed to update last run time: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("scheduler config not found: %s", name)
	}

	return nil
}

// UpdateNextRunTime updates next run time
func (r *configRepo) UpdateNextRunTime(ctx context.Context, name string, t time.Time) error {
	result, err := r.db.ExecContext(ctx,
		"UPDATE scheduler_config SET next_run_time = @p1, updated_at = @p2 WHERE config_name = @p3",
		t, time.Now(), name)
	if err != nil {
		return fmt.Errorf("failed to update next run time: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("scheduler config not found: %s", name)
	}

	return nil
}

// toConfigDomain converts database model to domain model
func toConfigDomain(db *configDB) *domain.SchedulerConfig {
	config := &domain.SchedulerConfig{
		ConfigName:    db.ConfigName,
		Enabled:       db.Enabled,
		IntervalMin:   db.IntervalMin,
		StartFromTime: db.StartFromTime,
		CreatedAt:     db.CreatedAt,
	}

	if db.LastStartTime.Valid {
		config.LastStartTime = &db.LastStartTime.Time
	}
	if db.LastStopTime.Valid {
		config.LastStopTime = &db.LastStopTime.Time
	}
	if db.LastRunTime.Valid {
		config.LastRunTime = &db.LastRunTime.Time
	}
	if db.NextRunTime.Valid {
		config.NextRunTime = &db.NextRunTime.Time
	}
	if db.UpdatedAt.Valid {
		config.UpdatedAt = &db.UpdatedAt.Time
	}

	return config
}
