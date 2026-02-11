package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Deo-Mugabe/GOLANG/internal/platform/database"
	"github.com/Deo-Mugabe/GOLANG/internal/scheduler/domain"
)

// executionDB is the database model for job_execution_history table
type executionDB struct {
	ExecutionID      int64        `db:"execution_id"`
	JobName          string       `db:"job_name"`
	JobGroup         string       `db:"job_group"`
	TriggerName      string       `db:"trigger_name"`
	TriggerGroup     string       `db:"trigger_group"`
	StartTime        time.Time    `db:"start_time"`
	EndTime          sql.NullTime `db:"end_time"`
	Status           string       `db:"status"`
	ErrorMessage     string       `db:"error_message"`
	RecordsProcessed int64        `db:"records_processed"`
	DurationMs       int64        `db:"duration_ms"`
	ProcessFromTime  sql.NullTime `db:"process_from_time"`
	ProcessToTime    sql.NullTime `db:"process_to_time"`
}

type executionRepo struct {
	db *database.DB
}

// NewExecutionRepository creates a new execution repository
func NewExecutionRepository(db *database.DB) domain.ExecutionRepository {
	return &executionRepo{db: db}
}

// GetByID retrieves an execution by ID
func (r *executionRepo) GetByID(ctx context.Context, id int64) (*domain.JobExecution, error) {
	query := `
        SELECT execution_id, job_name, job_group, trigger_name, trigger_group,
               start_time, end_time, status, error_message, records_processed,
               duration_ms, process_from_time, process_to_time
        FROM job_execution_history
        WHERE execution_id = @p1
    `

	var dbModel executionDB
	err := r.db.GetContext(ctx, &dbModel, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("execution not found: %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get execution: %w", err)
	}

	return toExecutionDomain(&dbModel), nil
}

// GetLatest retrieves the latest execution for a job
func (r *executionRepo) GetLatest(ctx context.Context, jobName, jobGroup string) (*domain.JobExecution, error) {
	query := `
        SELECT TOP 1 execution_id, job_name, job_group, trigger_name, trigger_group,
               start_time, end_time, status, error_message, records_processed,
               duration_ms, process_from_time, process_to_time
        FROM job_execution_history
        WHERE job_name = @p1 AND job_group = @p2
        ORDER BY start_time DESC
    `

	var dbModel executionDB
	err := r.db.GetContext(ctx, &dbModel, query, jobName, jobGroup)
	if err == sql.ErrNoRows {
		return nil, nil // No executions yet
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest execution: %w", err)
	}

	return toExecutionDomain(&dbModel), nil
}

// GetLastSuccessful retrieves the last successful execution
func (r *executionRepo) GetLastSuccessful(ctx context.Context, jobName, jobGroup string) (*domain.JobExecution, error) {
	query := `
        SELECT TOP 1 execution_id, job_name, job_group, trigger_name, trigger_group,
               start_time, end_time, status, error_message, records_processed,
               duration_ms, process_from_time, process_to_time
        FROM job_execution_history
        WHERE job_name = @p1 AND job_group = @p2 AND status = 'COMPLETED'
        ORDER BY start_time DESC
    `

	var dbModel executionDB
	err := r.db.GetContext(ctx, &dbModel, query, jobName, jobGroup)
	if err == sql.ErrNoRows {
		return nil, nil // No successful executions yet
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get last successful execution: %w", err)
	}

	return toExecutionDomain(&dbModel), nil
}

// List retrieves executions with pagination
func (r *executionRepo) List(ctx context.Context, jobName, jobGroup string, limit, offset int) ([]*domain.JobExecution, error) {
	query := `
        SELECT execution_id, job_name, job_group, trigger_name, trigger_group,
               start_time, end_time, status, error_message, records_processed,
               duration_ms, process_from_time, process_to_time
        FROM job_execution_history
        WHERE job_name = @p1 AND job_group = @p2
        ORDER BY start_time DESC
        OFFSET @p3 ROWS FETCH NEXT @p4 ROWS ONLY
    `

	var dbModels []executionDB
	err := r.db.SelectContext(ctx, &dbModels, query, jobName, jobGroup, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list executions: %w", err)
	}

	executions := make([]*domain.JobExecution, len(dbModels))
	for i, db := range dbModels {
		executions[i] = toExecutionDomain(&db)
	}

	return executions, nil
}

// ListByStatus retrieves executions by status
func (r *executionRepo) ListByStatus(ctx context.Context, status domain.ExecutionStatus, limit int) ([]*domain.JobExecution, error) {
	query := `
        SELECT TOP (@p2) execution_id, job_name, job_group, trigger_name, trigger_group,
               start_time, end_time, status, error_message, records_processed,
               duration_ms, process_from_time, process_to_time
        FROM job_execution_history
        WHERE status = @p1
        ORDER BY start_time DESC
    `

	var dbModels []executionDB
	err := r.db.SelectContext(ctx, &dbModels, query, string(status), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list executions by status: %w", err)
	}

	executions := make([]*domain.JobExecution, len(dbModels))
	for i, db := range dbModels {
		executions[i] = toExecutionDomain(&db)
	}

	return executions, nil
}

// ListSince retrieves executions since a time
func (r *executionRepo) ListSince(ctx context.Context, since time.Time, limit, offset int) ([]*domain.JobExecution, error) {
	query := `
        SELECT execution_id, job_name, job_group, trigger_name, trigger_group,
               start_time, end_time, status, error_message, records_processed,
               duration_ms, process_from_time, process_to_time
        FROM job_execution_history
        WHERE start_time >= @p1
        ORDER BY start_time DESC
        OFFSET @p2 ROWS FETCH NEXT @p3 ROWS ONLY
    `

	var dbModels []executionDB
	err := r.db.SelectContext(ctx, &dbModels, query, since, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list executions since: %w", err)
	}

	executions := make([]*domain.JobExecution, len(dbModels))
	for i, db := range dbModels {
		executions[i] = toExecutionDomain(&db)
	}

	return executions, nil
}

// Create creates a new execution record
func (r *executionRepo) Create(ctx context.Context, execution *domain.JobExecution) (int64, error) {
	query := `
        INSERT INTO job_execution_history (
            job_name, job_group, trigger_name, trigger_group, start_time,
            status, error_message, records_processed, duration_ms,
            process_from_time, process_to_time
        )
        OUTPUT INSERTED.execution_id
        VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11)
    `

	var executionID int64
	err := r.db.GetContext(ctx, &executionID, query,
		execution.JobName,
		execution.JobGroup,
		execution.TriggerName,
		execution.TriggerGroup,
		execution.StartTime,
		string(execution.Status),
		execution.ErrorMessage,
		execution.RecordsProcessed,
		execution.DurationMs,
		execution.ProcessFromTime,
		execution.ProcessToTime,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create execution: %w", err)
	}

	return executionID, nil
}

// Update updates an execution record
func (r *executionRepo) Update(ctx context.Context, execution *domain.JobExecution) error {
	query := `
        UPDATE job_execution_history
        SET end_time = @p1,
            status = @p2,
            error_message = @p3,
            records_processed = @p4,
            duration_ms = @p5,
            process_from_time = @p6,
            process_to_time = @p7
        WHERE execution_id = @p8
    `

	result, err := r.db.ExecContext(ctx, query,
		execution.EndTime,
		string(execution.Status),
		execution.ErrorMessage,
		execution.RecordsProcessed,
		execution.DurationMs,
		execution.ProcessFromTime,
		execution.ProcessToTime,
		execution.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update execution: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("execution not found: %d", execution.ID)
	}

	return nil
}

// UpdateStatus updates only the status
func (r *executionRepo) UpdateStatus(ctx context.Context, id int64, status domain.ExecutionStatus) error {
	result, err := r.db.ExecContext(ctx,
		"UPDATE job_execution_history SET status = @p1 WHERE execution_id = @p2",
		string(status), id)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("execution not found: %d", id)
	}

	return nil
}

// Complete marks execution as completed
func (r *executionRepo) Complete(ctx context.Context, id int64, endTime time.Time, recordsProcessed int64) error {
	query := `
        UPDATE job_execution_history
        SET end_time = @p1,
            status = 'COMPLETED',
            records_processed = @p2,
            duration_ms = DATEDIFF(MILLISECOND, start_time, @p1)
        WHERE execution_id = @p3
    `

	result, err := r.db.ExecContext(ctx, query, endTime, recordsProcessed, id)
	if err != nil {
		return fmt.Errorf("failed to complete execution: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("execution not found: %d", id)
	}

	return nil
}

// Fail marks execution as failed
func (r *executionRepo) Fail(ctx context.Context, id int64, endTime time.Time, errorMsg string) error {
	query := `
        UPDATE job_execution_history
        SET end_time = @p1,
            status = 'FAILED',
            error_message = @p2,
            duration_ms = DATEDIFF(MILLISECOND, start_time, @p1)
        WHERE execution_id = @p3
    `

	result, err := r.db.ExecContext(ctx, query, endTime, errorMsg, id)
	if err != nil {
		return fmt.Errorf("failed to mark execution as failed: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("execution not found: %d", id)
	}

	return nil
}

// CountByStatus counts executions by status
func (r *executionRepo) CountByStatus(ctx context.Context, status domain.ExecutionStatus) (int64, error) {
	var count int64
	err := r.db.GetContext(ctx, &count,
		"SELECT COUNT(*) FROM job_execution_history WHERE status = @p1",
		string(status))
	return count, err
}

// CountSuccessful counts successful executions for a job
func (r *executionRepo) CountSuccessful(ctx context.Context, jobName, jobGroup string) (int64, error) {
	var count int64
	err := r.db.GetContext(ctx, &count,
		"SELECT COUNT(*) FROM job_execution_history WHERE job_name = @p1 AND job_group = @p2 AND status = 'COMPLETED'",
		jobName, jobGroup)
	return count, err
}

// CountFailed counts failed executions for a job
func (r *executionRepo) CountFailed(ctx context.Context, jobName, jobGroup string) (int64, error) {
	var count int64
	err := r.db.GetContext(ctx, &count,
		"SELECT COUNT(*) FROM job_execution_history WHERE job_name = @p1 AND job_group = @p2 AND status = 'FAILED'",
		jobName, jobGroup)
	return count, err
}

// DeleteOlderThan deletes executions older than cutoff
func (r *executionRepo) DeleteOlderThan(ctx context.Context, cutoff time.Time) (int64, error) {
	result, err := r.db.ExecContext(ctx,
		"DELETE FROM job_execution_history WHERE start_time < @p1",
		cutoff)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old executions: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}

// toExecutionDomain converts database model to domain model
func toExecutionDomain(db *executionDB) *domain.JobExecution {
	exec := &domain.JobExecution{
		ID:               db.ExecutionID,
		JobName:          db.JobName,
		JobGroup:         db.JobGroup,
		TriggerName:      db.TriggerName,
		TriggerGroup:     db.TriggerGroup,
		StartTime:        db.StartTime,
		Status:           domain.ExecutionStatus(db.Status),
		ErrorMessage:     db.ErrorMessage,
		RecordsProcessed: db.RecordsProcessed,
		DurationMs:       db.DurationMs,
	}

	if db.EndTime.Valid {
		exec.EndTime = &db.EndTime.Time
	}
	if db.ProcessFromTime.Valid {
		exec.ProcessFromTime = &db.ProcessFromTime.Time
	}
	if db.ProcessToTime.Valid {
		exec.ProcessToTime = &db.ProcessToTime.Time
	}

	return exec
}
