package domain

import (
	"context"
	"time"
)

// SchedulerService manages job scheduling
type SchedulerService interface {
	// Scheduler control
	Start(ctx context.Context, configName string, intervalMin int) error
	Stop(ctx context.Context, configName string) error
	Restart(ctx context.Context, configName string) error

	// Status and monitoring
	GetStatus(ctx context.Context, configName string) (*SchedulerStatus, error)
	IsRunning(ctx context.Context, configName string) (bool, error)
	GetNextRunTime(ctx context.Context, configName string) (*time.Time, error)

	// Configuration
	UpdateInterval(ctx context.Context, configName string, intervalMin int) error
	UpdateStartFromTime(ctx context.Context, configName string, t time.Time) error

	// Manual execution
	TriggerNow(ctx context.Context, configName string) error
}

// ExecutorService handles job execution
type ExecutorService interface {
	// Execution
	Execute(ctx context.Context, jobName, jobGroup string) error

	// History tracking
	RecordStart(ctx context.Context, jobName, jobGroup string) (int64, error)
	RecordCompletion(ctx context.Context, executionID int64, recordsProcessed int64) error
	RecordFailure(ctx context.Context, executionID int64, err error) error

	// Processing time management
	GetLastProcessingTime(ctx context.Context, jobName string) (time.Time, error)
	UpdateProcessingTime(ctx context.Context, jobName string, t time.Time) error
}

// HistoryService manages job execution history
type HistoryService interface {
	// History retrieval
	GetExecutionHistory(ctx context.Context, jobName, jobGroup string, limit, offset int) ([]*JobExecution, error)
	GetLatestExecutions(ctx context.Context, limit int) ([]*JobExecution, error)
	GetFailedExecutions(ctx context.Context, limit int) ([]*JobExecution, error)

	// Statistics
	GetExecutionStats(ctx context.Context, jobName, jobGroup string) (*ExecutionStats, error)
	GetSuccessRate(ctx context.Context, jobName, jobGroup string, since time.Time) (float64, error)

	// Cleanup
	CleanupOldHistory(ctx context.Context, olderThan time.Duration) (int64, error)
}

// SchedulerStatus represents scheduler status
type SchedulerStatus struct {
	ConfigName           string
	Enabled              bool
	Running              bool
	IntervalMinutes      int
	LastStartTime        *time.Time
	LastStopTime         *time.Time
	LastRunTime          *time.Time
	NextRunTime          *time.Time
	TotalExecutions      int64
	SuccessfulExecutions int64
	FailedExecutions     int64
	LastError            string
}

// ExecutionStats contains execution statistics
type ExecutionStats struct {
	TotalExecutions       int64
	SuccessfulExecutions  int64
	FailedExecutions      int64
	AverageDurationMs     int64
	TotalRecordsProcessed int64
	LastSuccessfulRun     *time.Time
	LastFailedRun         *time.Time
}
