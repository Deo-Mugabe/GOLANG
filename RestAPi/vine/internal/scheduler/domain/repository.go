package domain

import (
	"context"
	"time"
)

// ConfigRepository handles scheduler configuration
type ConfigRepository interface {
	// Read operations
	Get(ctx context.Context, name string) (*SchedulerConfig, error)
	List(ctx context.Context) ([]*SchedulerConfig, error)

	// Write operations
	Create(ctx context.Context, config *SchedulerConfig) error
	Update(ctx context.Context, config *SchedulerConfig) error
	Delete(ctx context.Context, name string) error

	// Specific updates
	UpdateEnabled(ctx context.Context, name string, enabled bool) error
	UpdateInterval(ctx context.Context, name string, intervalMin int) error
	UpdateLastRunTime(ctx context.Context, name string, t time.Time) error
	UpdateNextRunTime(ctx context.Context, name string, t time.Time) error
}

// ExecutionRepository handles job execution history
type ExecutionRepository interface {
	// Read operations
	GetByID(ctx context.Context, id int64) (*JobExecution, error)
	GetLatest(ctx context.Context, jobName, jobGroup string) (*JobExecution, error)
	GetLastSuccessful(ctx context.Context, jobName, jobGroup string) (*JobExecution, error)
	List(ctx context.Context, jobName, jobGroup string, limit, offset int) ([]*JobExecution, error)
	ListByStatus(ctx context.Context, status ExecutionStatus, limit int) ([]*JobExecution, error)
	ListSince(ctx context.Context, since time.Time, limit, offset int) ([]*JobExecution, error)

	// Write operations
	Create(ctx context.Context, execution *JobExecution) (int64, error)
	Update(ctx context.Context, execution *JobExecution) error

	// Specific updates
	UpdateStatus(ctx context.Context, id int64, status ExecutionStatus) error
	Complete(ctx context.Context, id int64, endTime time.Time, recordsProcessed int64) error
	Fail(ctx context.Context, id int64, endTime time.Time, errorMsg string) error

	// Statistics
	CountByStatus(ctx context.Context, status ExecutionStatus) (int64, error)
	CountSuccessful(ctx context.Context, jobName, jobGroup string) (int64, error)
	CountFailed(ctx context.Context, jobName, jobGroup string) (int64, error)

	// Cleanup
	DeleteOlderThan(ctx context.Context, cutoff time.Time) (int64, error)
}
